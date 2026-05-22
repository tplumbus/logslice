package filter

import (
	"testing"
)

// stubFilter is a simple Filter implementation for testing.
type stubFilter struct {
	result bool
}

func (s *stubFilter) MatchesLine(_ string) bool { return s.result }

func TestNewInvertFilter_NilInner(t *testing.T) {
	_, err := NewInvertFilter(nil)
	if err == nil {
		t.Fatal("expected error for nil inner filter, got nil")
	}
}

func TestNewInvertFilter_Valid(t *testing.T) {
	f, err := NewInvertFilter(&stubFilter{result: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil InvertFilter")
	}
}

func TestInvertFilter_MatchesLine_InnerTrue(t *testing.T) {
	f, _ := NewInvertFilter(&stubFilter{result: true})
	if f.MatchesLine(`{"level":"info"}`) {
		t.Error("expected false when inner returns true")
	}
}

func TestInvertFilter_MatchesLine_InnerFalse(t *testing.T) {
	f, _ := NewInvertFilter(&stubFilter{result: false})
	if !f.MatchesLine(`{"level":"debug"}`) {
		t.Error("expected true when inner returns false")
	}
}

func TestInvertFilter_WithFieldQuery(t *testing.T) {
	q, err := ParseFieldQuery(`level=error`)
	if err != nil {
		t.Fatalf("ParseFieldQuery: %v", err)
	}
	inv, err := NewInvertFilter(q)
	if err != nil {
		t.Fatalf("NewInvertFilter: %v", err)
	}

	tests := []struct {
		line string
		want bool
	}{
		{`{"level":"error","msg":"boom"}`, false},
		{`{"level":"info","msg":"ok"}`, true},
		{`{"level":"debug","msg":"trace"}`, true},
	}

	for _, tc := range tests {
		got := inv.MatchesLine(tc.line)
		if got != tc.want {
			t.Errorf("MatchesLine(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestInvertFilter_WithRegexFilter(t *testing.T) {
	rx, err := NewRegexFilter("msg", "fail")
	if err != nil {
		t.Fatalf("NewRegexFilter: %v", err)
	}
	inv, _ := NewInvertFilter(rx)

	if inv.MatchesLine(`{"msg":"connection failed"}`) {
		t.Error("expected false: regex matches, invert should return false")
	}
	if !inv.MatchesLine(`{"msg":"all good"}`) {
		t.Error("expected true: regex does not match, invert should return true")
	}
}
