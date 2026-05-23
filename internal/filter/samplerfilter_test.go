package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestNewSamplerFilter_Valid(t *testing.T) {
	f, err := filter.NewSamplerFilter(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewSamplerFilter_ZeroInvalid(t *testing.T) {
	_, err := filter.NewSamplerFilter(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestNewSamplerFilter_NegativeInvalid(t *testing.T) {
	_, err := filter.NewSamplerFilter(-5)
	if err == nil {
		t.Fatal("expected error for negative n")
	}
}

func TestSamplerFilter_EveryLine(t *testing.T) {
	f, _ := filter.NewSamplerFilter(1)
	line := []byte(`{"msg":"hello"}`)
	for i := 0; i < 5; i++ {
		ok, err := f.MatchesLine(line)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Errorf("expected line %d to match with n=1", i+1)
		}
	}
}

func TestSamplerFilter_EveryThird(t *testing.T) {
	f, _ := filter.NewSamplerFilter(3)
	line := []byte(`{"msg":"hello"}`)
	expected := []bool{true, false, false, true, false, false, true}
	for i, want := range expected {
		got, err := f.MatchesLine(line)
		if err != nil {
			t.Fatalf("line %d: unexpected error: %v", i+1, err)
		}
		if got != want {
			t.Errorf("line %d: got %v, want %v", i+1, got, want)
		}
	}
}

func TestSamplerFilter_Reset(t *testing.T) {
	f, _ := filter.NewSamplerFilter(2)
	line := []byte(`{"msg":"hello"}`)

	// First call should match (count=1)
	ok, _ := f.MatchesLine(line)
	if !ok {
		t.Error("expected first call to match")
	}
	// Second call should not match (count=2)
	ok, _ = f.MatchesLine(line)
	if ok {
		t.Error("expected second call to not match")
	}

	f.Reset()

	// After reset, first call should match again
	ok, _ = f.MatchesLine(line)
	if !ok {
		t.Error("expected first call after reset to match")
	}
}

func TestSamplerFilter_WithPipeline(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"one"}`,
		`{"level":"info","msg":"two"}`,
		`{"level":"info","msg":"three"}`,
		`{"level":"info","msg":"four"}`,
		`{"level":"info","msg":"five"}`,
		`{"level":"info","msg":"six"}`,
	}

	sampler, _ := filter.NewSamplerFilter(2)
	p := filter.NewPipeline(sampler)

	var results []string
	err := p.Run(lines, func(line string) error {
		results = append(results, line)
		return nil
	})
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results (every 2nd), got %d", len(results))
	}
}
