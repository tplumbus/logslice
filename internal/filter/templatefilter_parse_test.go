package filter

import (
	"testing"
)

func TestParseTemplateFlag_Valid(t *testing.T) {
	f, err := ParseTemplateFlag(`{{.level}} {{.message}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestParseTemplateFlag_Empty(t *testing.T) {
	_, err := ParseTemplateFlag("")
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestParseTemplateFlag_WhitespaceOnly(t *testing.T) {
	_, err := ParseTemplateFlag("   ")
	if err == nil {
		t.Fatal("expected error for whitespace-only value")
	}
}

func TestParseTemplateFlag_UnescapesNewline(t *testing.T) {
	f, err := ParseTemplateFlag(`{{.level}}\n{{.message}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := f.TransformLine(`{"level":"info","message":"hello"}`)
	if got != "info\nhello" {
		t.Fatalf("got %q, want %q", got, "info\nhello")
	}
}

func TestParseTemplateFlag_UnescapesTab(t *testing.T) {
	f, err := ParseTemplateFlag(`{{.level}}\t{{.message}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := f.TransformLine(`{"level":"debug","message":"trace"}`)
	if got != "debug\ttrace" {
		t.Fatalf("got %q, want %q", got, "debug\ttrace")
	}
}

func TestParseTemplateFlag_InvalidTemplate(t *testing.T) {
	_, err := ParseTemplateFlag(`{{.level`)
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}
