package filter

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewFormatFilter_Valid(t *testing.T) {
	for _, fmt := range []string{"json", "pretty", "text", "JSON", "PRETTY", "TEXT"} {
		_, err := NewFormatFilter(fmt)
		if err != nil {
			t.Errorf("NewFormatFilter(%q) unexpected error: %v", fmt, err)
		}
	}
}

func TestNewFormatFilter_Invalid(t *testing.T) {
	_, err := NewFormatFilter("csv")
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestFormatFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	f, _ := NewFormatFilter("json")
	if !f.MatchesLine(`{"a":1}`) {
		t.Error("expected MatchesLine to return true")
	}
}

func TestFormatFilter_TransformLine_CompactJSON(t *testing.T) {
	f, _ := NewFormatFilter("json")
	input := `{ "level": "info", "msg": "hello" }`
	out, err := f.TransformLine(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "\n") {
		t.Error("compact JSON should not contain newlines")
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}
}

func TestFormatFilter_TransformLine_PrettyJSON(t *testing.T) {
	f, _ := NewFormatFilter("pretty")
	input := `{"level":"info","msg":"hello"}`
	out, err := f.TransformLine(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\n") {
		t.Error("pretty JSON should contain newlines")
	}
	if !strings.Contains(out, "  ") {
		t.Error("pretty JSON should contain indentation")
	}
}

func TestFormatFilter_TransformLine_Text(t *testing.T) {
	f, _ := NewFormatFilter("text")
	input := `{"level":"info"}`
	out, err := f.TransformLine(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "level=info") {
		t.Errorf("expected key=value pair in output, got: %s", out)
	}
}

func TestFormatFilter_TransformLine_NonJSON_PassThrough(t *testing.T) {
	f, _ := NewFormatFilter("pretty")
	input := "not json at all"
	out, err := f.TransformLine(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("expected pass-through for non-JSON, got: %s", out)
	}
}

func TestParseFormatFlag_Valid(t *testing.T) {
	f, err := ParseFormatFlag("pretty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestParseFormatFlag_Empty(t *testing.T) {
	_, err := ParseFormatFlag("")
	if err == nil {
		t.Fatal("expected error for empty flag")
	}
}
