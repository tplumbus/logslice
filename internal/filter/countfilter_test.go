package filter

import (
	"fmt"
	"testing"
)

func TestNewCountFilter_NoField(t *testing.T) {
	cf, err := NewCountFilter("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cf == nil {
		t.Fatal("expected non-nil CountFilter")
	}
}

func TestNewCountFilter_WithField(t *testing.T) {
	cf, err := NewCountFilter("level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cf == nil {
		t.Fatal("expected non-nil CountFilter")
	}
}

func TestCountFilter_Total_NoField(t *testing.T) {
	cf, _ := NewCountFilter("")
	lines := []string{
		`{"msg":"a"}`,
		`{"msg":"b"}`,
		`{"msg":"c"}`,
	}
	for _, l := range lines {
		match, err := cf.MatchesLine(l)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !match {
			t.Error("expected CountFilter to always return true")
		}
	}
	if cf.Total() != 3 {
		t.Errorf("expected total 3, got %d", cf.Total())
	}
	if cf.Counts() != nil {
		t.Error("expected nil counts when no field set")
	}
}

func TestCountFilter_CountsByField(t *testing.T) {
	cf, _ := NewCountFilter("level")
	input := []struct {
		line  string
		want  string
	}{
		{`{"level":"info","msg":"a"}`, "info"},
		{`{"level":"warn","msg":"b"}`, "warn"},
		{`{"level":"info","msg":"c"}`, "info"},
		{`{"level":"error","msg":"d"}`, "error"},
	}
	for _, tc := range input {
		cf.MatchesLine(tc.line)
	}
	counts := cf.Counts()
	if counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", counts["info"])
	}
	if counts["warn"] != 1 {
		t.Errorf("expected warn=1, got %d", counts["warn"])
	}
	if counts["error"] != 1 {
		t.Errorf("expected error=1, got %d", counts["error"])
	}
	if cf.Total() != 4 {
		t.Errorf("expected total 4, got %d", cf.Total())
	}
}

func TestCountFilter_MissingField(t *testing.T) {
	cf, _ := NewCountFilter("level")
	cf.MatchesLine(`{"msg":"no level here"}`)
	counts := cf.Counts()
	if counts["<missing>"] != 1 {
		t.Errorf("expected <missing>=1, got %d", counts["<missing>"])
	}
}

func TestCountFilter_InvalidJSON(t *testing.T) {
	cf, _ := NewCountFilter("level")
	cf.MatchesLine(`not-json`)
	counts := cf.Counts()
	if counts["<invalid>"] != 1 {
		t.Errorf("expected <invalid>=1, got %d", counts["<invalid>"])
	}
}

func TestCountFilter_Summary(t *testing.T) {
	cf, _ := NewCountFilter("")
	cf.MatchesLine(`{"msg":"a"}`)
	cf.MatchesLine(`{"msg":"b"}`)
	summary := cf.Summary()
	expected := fmt.Sprintf("total lines: %d", 2)
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
