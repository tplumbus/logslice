package filter

import (
	"testing"
)

func TestParseFieldQuery_Valid(t *testing.T) {
	tests := []struct {
		input       string
		wantField   string
		wantValue   string
	}{
		{"level=info", "level", "info"},
		{"service=auth", "service", "auth"},
		{"code=200", "code", "200"},
		{"msg=hello world", "msg", "hello world"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseFieldQuery(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Field != tc.wantField || got.Value != tc.wantValue {
				t.Errorf("got (%q, %q), want (%q, %q)", got.Field, got.Value, tc.wantField, tc.wantValue)
			}
		})
	}
}

func TestParseFieldQuery_Invalid(t *testing.T) {
	inputs := []string{"", "noequalssign", "=value"}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			_, err := ParseFieldQuery(input)
			if err == nil {
				t.Errorf("expected error for input %q, got nil", input)
			}
		})
	}
}

func TestFieldQuery_MatchesLine(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		line    string
		want    bool
	}{
		{"string match", "level=info", `{"level":"info","msg":"started"}`, true},
		{"string no match", "level=error", `{"level":"info","msg":"started"}`, false},
		{"numeric match", "code=200", `{"code":200,"path":"/health"}`, true},
		{"numeric no match", "code=404", `{"code":200,"path":"/health"}`, false},
		{"missing field", "host=localhost", `{"level":"info"}`, false},
		{"boolean match", "ok=true", `{"ok":true}`, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fq, err := ParseFieldQuery(tc.query)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			got, err := fq.MatchesLine([]byte(tc.line))
			if err != nil {
				t.Fatalf("match error: %v", err)
			}
			if got != tc.want {
				t.Errorf("MatchesLine(%q) = %v, want %v", tc.line, got, tc.want)
			}
		})
	}
}

func TestFieldQuery_MatchesLine_InvalidJSON(t *testing.T) {
	fq := FieldQuery{Field: "level", Value: "info"}
	_, err := fq.MatchesLine([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
