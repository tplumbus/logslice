package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestJSONPathFilter_WithPipeline_NestedMatch(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:00Z","meta":{"env":"production","region":"us-east"}}`,
		`{"ts":"2024-01-01T00:00:01Z","meta":{"env":"staging","region":"eu-west"}}`,
		`{"ts":"2024-01-01T00:00:02Z","meta":{"env":"production","region":"ap-south"}}`,
	}

	jpf, err := NewJSONPathFilter("meta.env", "production")
	if err != nil {
		t.Fatalf("NewJSONPathFilter: %v", err)
	}

	pipeline := NewPipeline(jpf)
	var buf bytes.Buffer
	if err := pipeline.Run(strings.NewReader(strings.Join(lines, "\n")), &buf); err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	got := strings.Split(out, "\n")
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
	for _, line := range got {
		if !strings.Contains(line, "production") {
			t.Errorf("unexpected line without 'production': %s", line)
		}
	}
}

func TestJSONPathFilter_WithPipeline_ParsedFromFlag(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:00Z","request":{"method":"GET","path":"/health"}}`,
		`{"ts":"2024-01-01T00:00:01Z","request":{"method":"POST","path":"/data"}}`,
		`{"ts":"2024-01-01T00:00:02Z","request":{"method":"GET","path":"/metrics"}}`,
	}

	jpf, err := ParseJSONPathFlag("request.method=GET")
	if err != nil {
		t.Fatalf("ParseJSONPathFlag: %v", err)
	}

	pipeline := NewPipeline(jpf)
	var buf bytes.Buffer
	if err := pipeline.Run(strings.NewReader(strings.Join(lines, "\n")), &buf); err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	got := strings.Split(out, "\n")
	if len(got) != 2 {
		t.Fatalf("expected 2 GET lines, got %d", len(got))
	}
	for _, line := range got {
		if !strings.Contains(line, "GET") {
			t.Errorf("line missing GET method: %s", line)
		}
	}
}
