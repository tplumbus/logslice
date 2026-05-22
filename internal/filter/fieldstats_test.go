package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestFieldStats_Observe_CountsValues(t *testing.T) {
	fs := NewFieldStats("level")
	lines := []string{
		`{"level":"info","msg":"started"}`,
		`{"level":"warn","msg":"slow"}`,
		`{"level":"info","msg":"done"}`,
		`{"level":"error","msg":"fail"}`,
		`{"level":"info","msg":"retry"}`,
	}
	for _, l := range lines {
		fs.Observe(l)
	}
	if fs.Total != 5 {
		t.Fatalf("expected total 5, got %d", fs.Total)
	}
	if fs.Counts["info"] != 3 {
		t.Errorf("expected info=3, got %d", fs.Counts["info"])
	}
	if fs.Counts["warn"] != 1 {
		t.Errorf("expected warn=1, got %d", fs.Counts["warn"])
	}
	if fs.Missed != 0 {
		t.Errorf("expected no missed, got %d", fs.Missed)
	}
}

func TestFieldStats_Observe_MissingField(t *testing.T) {
	fs := NewFieldStats("level")
	fs.Observe(`{"msg":"no level here"}`)
	if fs.Missed != 1 {
		t.Errorf("expected missed=1, got %d", fs.Missed)
	}
}

func TestFieldStats_Observe_InvalidJSON(t *testing.T) {
	fs := NewFieldStats("level")
	fs.Observe(`not json`)
	if fs.Missed != 1 {
		t.Errorf("expected missed=1 for invalid json, got %d", fs.Missed)
	}
	if fs.Total != 1 {
		t.Errorf("expected total=1, got %d", fs.Total)
	}
}

func TestFieldStats_TopN_Order(t *testing.T) {
	fs := NewFieldStats("status")
	for i := 0; i < 5; i++ {
		fs.Observe(`{"status":"200"}`)
	}
	for i := 0; i < 2; i++ {
		fs.Observe(`{"status":"500"}`)
	}
	fs.Observe(`{"status":"404"}`)

	top := fs.TopN(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
	if top[0].Value != "200" || top[0].Count != 5 {
		t.Errorf("expected top entry 200:5, got %s:%d", top[0].Value, top[0].Count)
	}
	if top[1].Value != "500" || top[1].Count != 2 {
		t.Errorf("expected second entry 500:2, got %s:%d", top[1].Value, top[1].Count)
	}
}

func TestFieldStats_TopN_AllWhenZero(t *testing.T) {
	fs := NewFieldStats("level")
	fs.Observe(`{"level":"info"}`)
	fs.Observe(`{"level":"warn"}`)
	top := fs.TopN(0)
	if len(top) != 2 {
		t.Errorf("expected 2 entries, got %d", len(top))
	}
}

func TestFieldStats_WriteSummary(t *testing.T) {
	fs := NewFieldStats("level")
	fs.Observe(`{"level":"info"}`)
	fs.Observe(`{"level":"info"}`)
	fs.Observe(`{"level":"warn"}`)

	var buf bytes.Buffer
	fs.WriteSummary(&buf)
	out := buf.String()

	if !strings.Contains(out, "field: level") {
		t.Errorf("summary missing field name, got: %s", out)
	}
	if !strings.Contains(out, "info: 2") {
		t.Errorf("summary missing info count, got: %s", out)
	}
	if !strings.Contains(out, "total lines: 3") {
		t.Errorf("summary missing total, got: %s", out)
	}
}
