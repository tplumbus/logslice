package filter

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSelectFilter_WithPipeline_ProjectsFields(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:01Z","level":"info","msg":"started","pid":1}`,
		`{"ts":"2024-01-01T00:00:02Z","level":"error","msg":"failed","pid":2}`,
		`{"ts":"2024-01-01T00:00:03Z","level":"info","msg":"done","pid":3}`,
	}

	sf, err := NewSelectFilter([]string{"level", "msg"})
	if err != nil {
		t.Fatalf("NewSelectFilter: %v", err)
	}

	p := NewPipeline(sf)
	var results []string
	err = p.Run(lines, func(line string) error {
		results = append(results, line)
		return nil
	})
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(r), &m); err != nil {
			t.Fatalf("result not valid JSON: %v", err)
		}
		if _, ok := m["ts"]; ok {
			t.Error("field 'ts' should be absent")
		}
		if _, ok := m["pid"]; ok {
			t.Error("field 'pid' should be absent")
		}
		if _, ok := m["level"]; !ok {
			t.Error("field 'level' should be present")
		}
		if _, ok := m["msg"]; !ok {
			t.Error("field 'msg' should be present")
		}
	}
}

func TestSelectFilter_WithPipeline_CombinedWithFieldFilter(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"ok","code":200}`,
		`{"level":"error","msg":"bad","code":500}`,
		`{"level":"info","msg":"also ok","code":201}`,
	}

	fq, err := ParseFieldQuery("level=error")
	if err != nil {
		t.Fatalf("ParseFieldQuery: %v", err)
	}
	sf, err := NewSelectFilter([]string{"msg", "code"})
	if err != nil {
		t.Fatalf("NewSelectFilter: %v", err)
	}

	p := NewPipeline(fq, sf)
	var results []string
	_ = p.Run(lines, func(line string) error {
		results = append(results, line)
		return nil
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.Contains(results[0], `"bad"`) {
		t.Errorf("expected msg=bad in result, got %s", results[0])
	}
	if strings.Contains(results[0], `"level"`) {
		t.Errorf("field 'level' should be projected out, got %s", results[0])
	}
}
