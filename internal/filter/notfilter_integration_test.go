package filter

import (
	"strings"
	"testing"
)

func TestNotFilter_WithPipeline_ExcludesMatchingLines(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"startup"}`,
		`{"level":"debug","msg":"verbose"}`,
		`{"level":"warn","msg":"slow query"}`,
		`{"level":"debug","msg":"trace"}`,
		`{"level":"error","msg":"crash"}`,
	}

	fq, _ := ParseFieldQuery("level=debug")
	nf, err := NewNotFilter(fq)
	if err != nil {
		t.Fatalf("NewNotFilter error: %v", err)
	}

	p := NewPipeline(nf)
	var out []string
	for _, l := range lines {
		result, pass, _ := p.Run(l)
		if pass {
			out = append(out, result)
		}
	}

	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(out), out)
	}
	for _, line := range out {
		if strings.Contains(line, `"debug"`) {
			t.Errorf("debug line should have been excluded: %s", line)
		}
	}
}

func TestNotFilter_WithPipeline_ParsedFromFlag(t *testing.T) {
	lines := []string{
		`{"level":"info","env":"prod"}`,
		`{"level":"debug","env":"prod"}`,
		`{"level":"info","env":"staging"}`,
		`{"level":"warn","env":"prod"}`,
	}

	nf, err := ParseNotFlag("level=debug,env=staging")
	if err != nil {
		t.Fatalf("ParseNotFlag error: %v", err)
	}

	p := NewPipeline(nf)
	var out []string
	for _, l := range lines {
		result, pass, _ := p.Run(l)
		if pass {
			out = append(out, result)
		}
	}

	// Only {"level":"info","env":"prod"} and {"level":"warn","env":"prod"} should pass
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(out), out)
	}
	for _, line := range out {
		if strings.Contains(line, `"debug"`) || strings.Contains(line, `"staging"`) {
			t.Errorf("excluded line leaked through: %s", line)
		}
	}
}
