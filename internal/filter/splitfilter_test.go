package filter

import (
	"encoding/json"
	"testing"
)

func TestNewSplitFilter_Valid(t *testing.T) {
	f, err := NewSplitFilter("tags", ",", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.outField != "tags" {
		t.Errorf("expected outField 'tags', got %q", f.outField)
	}
}

func TestNewSplitFilter_EmptyField(t *testing.T) {
	_, err := NewSplitFilter("", ",", "")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewSplitFilter_EmptyDelimiter(t *testing.T) {
	_, err := NewSplitFilter("tags", "", "")
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestSplitFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	f, _ := NewSplitFilter("tags", ",", "")
	if !f.MatchesLine(`{"tags":"a,b"}`) {
		t.Error("expected MatchesLine to always return true")
	}
}

func TestSplitFilter_TransformLine_SplitsField(t *testing.T) {
	f, _ := NewSplitFilter("tags", ",", "")
	out, err := f.TransformLine(`{"tags":"a,b,c"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	arr, ok := obj["tags"].([]interface{})
	if !ok || len(arr) != 3 {
		t.Errorf("expected 3-element array, got %v", obj["tags"])
	}
}

func TestSplitFilter_TransformLine_CustomOutField(t *testing.T) {
	f, _ := NewSplitFilter("tags", ",", "tag_list")
	out, _ := f.TransformLine(`{"tags":"x,y"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if _, ok := obj["tag_list"]; !ok {
		t.Error("expected 'tag_list' field in output")
	}
}

func TestSplitFilter_TransformLine_MissingField(t *testing.T) {
	f, _ := NewSplitFilter("tags", ",", "")
	line := `{"level":"info"}`
	out, _ := f.TransformLine(line)
	if out != line {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestSplitFilter_TransformLine_NonStringField(t *testing.T) {
	f, _ := NewSplitFilter("count", ",", "")
	line := `{"count":42}`
	out, _ := f.TransformLine(line)
	if out != line {
		t.Errorf("expected unchanged line for non-string field, got %q", out)
	}
}

func TestParseSplitFlag_Valid(t *testing.T) {
	f, err := ParseSplitFlag("tags:,")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.field != "tags" || f.delimiter != "," {
		t.Errorf("unexpected filter state: %+v", f)
	}
}

func TestParseSplitFlag_WithOutField(t *testing.T) {
	f, err := ParseSplitFlag("tags:,:tag_list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.outField != "tag_list" {
		t.Errorf("expected outField 'tag_list', got %q", f.outField)
	}
}

func TestParseSplitFlag_TabDelimiter(t *testing.T) {
	f, err := ParseSplitFlag(`cols:\t`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.delimiter != "\t" {
		t.Errorf("expected tab delimiter, got %q", f.delimiter)
	}
}

func TestParseSplitFlag_MissingDelimiter(t *testing.T) {
	_, err := ParseSplitFlag("tags")
	if err == nil {
		t.Fatal("expected error for missing delimiter")
	}
}
