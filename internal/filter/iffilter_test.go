package filter

import (
	"encoding/json"
	"testing"
)

func TestNewIfFilter_Valid(t *testing.T) {
	cond, _ := ParseFieldQuery("level=error")
	trans, _ := NewAddFieldsFilter(map[string]string{"flagged": "true"})
	f, err := NewIfFilter(cond, trans)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewIfFilter_NilCondition(t *testing.T) {
	trans, _ := NewAddFieldsFilter(map[string]string{"flagged": "true"})
	_, err := NewIfFilter(nil, trans)
	if err == nil {
		t.Fatal("expected error for nil condition")
	}
}

func TestNewIfFilter_NilTransform(t *testing.T) {
	cond, _ := ParseFieldQuery("level=error")
	_, err := NewIfFilter(cond, nil)
	if err == nil {
		t.Fatal("expected error for nil transform")
	}
}

func TestIfFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	cond, _ := ParseFieldQuery("level=error")
	trans, _ := NewAddFieldsFilter(map[string]string{"flagged": "true"})
	f, _ := NewIfFilter(cond, trans)

	lines := []string{
		`{"level":"error","msg":"boom"}`,
		`{"level":"info","msg":"ok"}`,
		`not json at all`,
	}
	for _, l := range lines {
		if !f.MatchesLine(l) {
			t.Errorf("MatchesLine(%q) = false, want true", l)
		}
	}
}

func TestIfFilter_TransformLine_AppliesWhenConditionMatches(t *testing.T) {
	cond, _ := ParseFieldQuery("level=error")
	trans, _ := NewAddFieldsFilter(map[string]string{"flagged": "true"})
	f, _ := NewIfFilter(cond, trans)

	line := `{"level":"error","msg":"boom"}`
	out := f.TransformLine(line)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["flagged"] != "true" {
		t.Errorf("expected flagged=true, got %v", m["flagged"])
	}
}

func TestIfFilter_TransformLine_PassthroughWhenNoMatch(t *testing.T) {
	cond, _ := ParseFieldQuery("level=error")
	trans, _ := NewAddFieldsFilter(map[string]string{"flagged": "true"})
	f, _ := NewIfFilter(cond, trans)

	line := `{"level":"info","msg":"all good"}`
	out := f.TransformLine(line)
	if out != line {
		t.Errorf("expected passthrough, got %q", out)
	}
}

func TestParseIfFlag_Valid(t *testing.T) {
	f, err := ParseIfFlag("level=error:flagged=true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"level":"error","msg":"boom"}`
	out := f.TransformLine(line)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if m["flagged"] != "true" {
		t.Errorf("expected flagged=true, got %v", m["flagged"])
	}
}

func TestParseIfFlag_MissingColon(t *testing.T) {
	_, err := ParseIfFlag("level=error")
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseIfFlag_EmptyCondition(t *testing.T) {
	_, err := ParseIfFlag(":flagged=true")
	if err == nil {
		t.Fatal("expected error for empty condition")
	}
}

func TestParseIfFlag_EmptyTransform(t *testing.T) {
	_, err := ParseIfFlag("level=error:")
	if err == nil {
		t.Fatal("expected error for empty transform")
	}
}
