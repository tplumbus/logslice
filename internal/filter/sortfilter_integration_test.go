package filter_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestSortFilter_WithPipeline_SortsOutput(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-03T00:00:00Z","level":"warn","msg":"c"}`,
		`{"ts":"2024-01-01T00:00:00Z","level":"info","msg":"a"}`,
		`{"ts":"2024-01-02T00:00:00Z","level":"error","msg":"b"}`,
	}

	sf, err := filter.NewSortFilter("msg", filter.SortAsc)
	if err != nil {
		t.Fatalf("NewSortFilter: %v", err)
	}

	for _, l := range lines {
		sf.Observe(l)
	}
	got := sf.Sorted()

	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	wantOrder := []string{"a", "b", "c"}
	for i, w := range wantOrder {
		if !strings.Contains(got[i], fmt.Sprintf(`"msg":%q`, w)) {
			t.Errorf("line %d: want msg=%q, got: %s", i, w, got[i])
		}
	}
}

func TestSortFilter_WithPipeline_DescOrder(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:00Z","code":"100"}`,
		`{"ts":"2024-01-02T00:00:00Z","code":"300"}`,
		`{"ts":"2024-01-03T00:00:00Z","code":"200"}`,
	}

	sf, err := filter.NewSortFilter("code", filter.SortDesc)
	if err != nil {
		t.Fatalf("NewSortFilter: %v", err)
	}

	for _, l := range lines {
		sf.Observe(l)
	}
	got := sf.Sorted()

	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	wantOrder := []string{"300", "200", "100"}
	for i, w := range wantOrder {
		if !strings.Contains(got[i], fmt.Sprintf(`"code":%q`, w)) {
			t.Errorf("line %d: want code=%q, got: %s", i, w, got[i])
		}
	}
}
