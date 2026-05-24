package filter

import (
	"testing"
)

func TestNewTemplateFilter_Valid(t *testing.T) {
	_, err := NewTemplateFilter(`{{.level}} {{.message}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewTemplateFilter_Empty(t *testing.T) {
	_, err := NewTemplateFilter("")
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestNewTemplateFilter_InvalidSyntax(t *testing.T) {
	_, err := NewTemplateFilter(`{{.level`)
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}

func TestTemplateFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	f, _ := NewTemplateFilter(`{{.level}}`)
	if !f.MatchesLine(`{"level":"info"}`) {
		t.Fatal("MatchesLine should always return true")
	}
	if !f.MatchesLine(`not json at all`) {
		t.Fatal("MatchesLine should return true even for non-JSON")
	}
}

func TestTemplateFilter_TransformLine_Basic(t *testing.T) {
	f, _ := NewTemplateFilter(`[{{.level}}] {{.message}}`)
	got := f.TransformLine(`{"level":"warn","message":"disk full"}`)
	want := "[warn] disk full"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestTemplateFilter_TransformLine_MissingKey(t *testing.T) {
	f, _ := NewTemplateFilter(`{{.level}} {{.missing}}`)
	// missingkey=zero means missing fields render as <no value> or zero
	got := f.TransformLine(`{"level":"info"}`)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestTemplateFilter_TransformLine_InvalidJSON(t *testing.T) {
	f, _ := NewTemplateFilter(`{{.level}}`)
	input := `not valid json`
	got := f.TransformLine(input)
	if got != input {
		t.Fatalf("expected original line back, got %q", got)
	}
}

func TestTemplateFilter_TransformLine_NumericField(t *testing.T) {
	f, _ := NewTemplateFilter(`status={{.status}}`)
	got := f.TransformLine(`{"status":200}`)
	want := "status=200"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
