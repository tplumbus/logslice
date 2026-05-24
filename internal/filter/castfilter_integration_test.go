package filter

import (
	"strings"
	"testing"
)

func TestCastFilter_WithPipeline_CastsField(t *testing.T) {
	lines := []string{
		`{"level":"info","status":"200"}`,
		`{"level":"error","status":"500"}`,
		`{"level":"warn","status":"404"}`,
	}

	cf, err := NewCastFilter("status", "int")
	if err != nil {
		t.Fatalf("NewCastFilter: %v", err)
	}

	p := NewPipeline(cf)
	var buf strings.Builder
	if err := p.Run(lines, &buf); err != nil {
		t.Fatalf("pipeline run: %v", err)
	}

	result := buf.String()
	if strings.Contains(result, `"status":"200"`) {
		t.Error("expected status to be cast from string to number")
	}
	if !strings.Contains(result, `"status":200`) {
		t.Errorf("expected numeric status 200 in output, got: %s", result)
	}
}

func TestCastFilter_WithPipeline_CombinedWithFieldFilter(t *testing.T) {
	lines := []string{
		`{"level":"info","retries":"3"}`,
		`{"level":"debug","retries":"0"}`,
		`{"level":"info","retries":"1"}`,
	}

	fq, err := ParseFieldQuery("level=info")
	if err != nil {
		t.Fatalf("ParseFieldQuery: %v", err)
	}
	cf, err := NewCastFilter("retries", "int")
	if err != nil {
		t.Fatalf("NewCastFilter: %v", err)
	}

	p := NewPipeline(fq, cf)
	var buf strings.Builder
	if err := p.Run(lines, &buf); err != nil {
		t.Fatalf("pipeline run: %v", err)
	}

	result := buf.String()
	outLines := strings.Split(strings.TrimSpace(result), "\n")
	if len(outLines) != 2 {
		t.Errorf("expected 2 lines (info only), got %d: %s", len(outLines), result)
	}
	for _, l := range outLines {
		if strings.Contains(l, `"retries":"`) {
			t.Errorf("expected retries to be numeric, got string in: %s", l)
		}
	}
}
