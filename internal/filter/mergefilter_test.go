package filter

import (
	"encoding/json"
	"testing"
)

func TestNewMergeFilter_Valid(t *testing.T) {
	_, err := NewMergeFilter([]string{"env=prod", "region=us-east"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewMergeFilter_Empty(t *testing.T) {
	_, err := NewMergeFilter([]string{}, false)
	if err == nil {
		t.Fatal("expected error for empty pairs")
	}
}

func TestNewMergeFilter_InvalidPair(t *testing.T) {
	_, err := NewMergeFilter([]string{"noequalssign"}, false)
	if err == nil {
		t.Fatal("expected error for invalid pair")
	}
}

func TestMergeFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	f, _ := NewMergeFilter([]string{"k=v"}, false)
	if !f.MatchesLine(`{"msg":"hello"}`) {
		t.Fatal("expected MatchesLine to return true")
	}
}

func TestMergeFilter_TransformLine_AddsFields(t *testing.T) {
	f, _ := NewMergeFilter([]string{"env=prod", "region=us-east"}, false)
	out, err := f.TransformLine(`{"msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if obj["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", obj["env"])
	}
	if obj["region"] != "us-east" {
		t.Errorf("expected region=us-east, got %v", obj["region"])
	}
}

func TestMergeFilter_TransformLine_NoOverwrite(t *testing.T) {
	f, _ := NewMergeFilter([]string{"env=prod"}, false)
	out, _ := f.TransformLine(`{"env":"dev","msg":"hi"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["env"] != "dev" {
		t.Errorf("expected original env=dev to be preserved, got %v", obj["env"])
	}
}

func TestMergeFilter_TransformLine_WithOverwrite(t *testing.T) {
	f, _ := NewMergeFilter([]string{"env=prod"}, true)
	out, _ := f.TransformLine(`{"env":"dev","msg":"hi"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["env"] != "prod" {
		t.Errorf("expected env to be overwritten to prod, got %v", obj["env"])
	}
}

func TestMergeFilter_TransformLine_InvalidJSON(t *testing.T) {
	f, _ := NewMergeFilter([]string{"env=prod"}, false)
	out, _ := f.TransformLine(`not-json`)
	if out != `not-json` {
		t.Errorf("expected passthrough for invalid JSON, got %q", out)
	}
}
