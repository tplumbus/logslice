package filter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/input"
	"github.com/yourorg/logslice/internal/output"
)

func TestIfFilter_WithPipeline_TransformsOnlyMatchingLines(t *testing.T) {
	lines := strings.Join([]string{
		`{"level":"error","msg":"disk full"}`,
		`{"level":"info","msg":"started"}`,
		`{"level":"error","msg":"timeout"}`,
		`{"level":"warn","msg":"slow"}`,
	}, "\n")

	cond, _ := filter.ParseFieldQuery("level=error")
	trans, _ := filter.NewAddFieldsFilter(map[string]string{"alert": "1"})
	ifF, err := filter.NewIfFilter(cond, trans)
	if err != nil {
		t.Fatalf("NewIfFilter: %v", err)
	}

	reader, _ := input.NewLineReader(nil, strings.NewReader(lines))
	var buf bytes.Buffer
	writer, _ := output.NewWriterFromIO(&buf)
	pipe := filter.NewPipeline(reader, writer, []filter.Filter{ifF})
	if err := pipe.Run(); err != nil {
		t.Fatalf("pipeline run: %v", err)
	}

	outLines := nonBlankLines(buf.String())
	if len(outLines) != 4 {
		t.Fatalf("expected 4 output lines, got %d", len(outLines))
	}

	for _, raw := range outLines {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if m["level"] == "error" {
			if m["alert"] != "1" {
				t.Errorf("error line missing alert field: %s", raw)
			}
		} else {
			if _, ok := m["alert"]; ok {
				t.Errorf("non-error line should not have alert field: %s", raw)
			}
		}
	}
}

func TestIfFilter_WithPipeline_ParsedFromFlag(t *testing.T) {
	lines := strings.Join([]string{
		`{"status":"500","path":"/api"}`,
		`{"status":"200","path":"/health"}`,
	}, "\n")

	ifF, err := filter.ParseIfFlag("status=500:critical=yes")
	if err != nil {
		t.Fatalf("ParseIfFlag: %v", err)
	}

	reader, _ := input.NewLineReader(nil, strings.NewReader(lines))
	var buf bytes.Buffer
	writer, _ := output.NewWriterFromIO(&buf)
	pipe := filter.NewPipeline(reader, writer, []filter.Filter{ifF})
	if err := pipe.Run(); err != nil {
		t.Fatalf("pipeline run: %v", err)
	}

	outLines := nonBlankLines(buf.String())
	for _, raw := range outLines {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if m["status"] == "500" && m["critical"] != "yes" {
			t.Errorf("expected critical=yes on 500 line: %s", raw)
		}
		if m["status"] == "200" {
			if _, ok := m["critical"]; ok {
				t.Errorf("200 line should not have critical field: %s", raw)
			}
		}
	}
}

func nonBlankLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}
