package filter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestDedupFilter_WithPipeline_RemovesDuplicates(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:01Z","request_id":"abc","msg":"request start"}`,
		`{"ts":"2024-01-01T00:00:02Z","request_id":"def","msg":"request start"}`,
		`{"ts":"2024-01-01T00:00:03Z","request_id":"abc","msg":"duplicate"}`,
		`{"ts":"2024-01-01T00:00:04Z","request_id":"ghi","msg":"unique"}`,
	}

	dedup, err := filter.NewDedupFilter("request_id", 0)
	if err != nil {
		t.Fatalf("failed to create dedup filter: %v", err)
	}

	p := filter.NewPipeline(lines, dedup)
	results, runErr := p.Run()
	if runErr != nil {
		t.Fatalf("pipeline run failed: %v", runErr)
	}

	if len(results) != 3 {
		t.Errorf("expected 3 unique results, got %d", len(results))
	}

	for _, r := range results {
		if strings.Contains(r, `"duplicate"`) {
			t.Errorf("duplicate line should have been filtered: %s", r)
		}
	}
}

func TestDedupFilter_WithPipeline_MaxSeenAllowsN(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:01Z","level":"warn","msg":"a"}`,
		`{"ts":"2024-01-01T00:00:02Z","level":"warn","msg":"b"}`,
		`{"ts":"2024-01-01T00:00:03Z","level":"warn","msg":"c"}`,
		`{"ts":"2024-01-01T00:00:04Z","level":"error","msg":"d"}`,
	}

	dedup, err := filter.NewDedupFilter("level", 2)
	if err != nil {
		t.Fatalf("failed to create dedup filter: %v", err)
	}

	p := filter.NewPipeline(lines, dedup)
	results, runErr := p.Run()
	if runErr != nil {
		t.Fatalf("pipeline run failed: %v", runErr)
	}

	// "warn" appears 3 times but maxSeen=2 so only 2 pass; "error" appears once
	if len(results) != 3 {
		t.Errorf("expected 3 results (2 warn + 1 error), got %d: %v", len(results), results)
	}
}
