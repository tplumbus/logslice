package filter_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func makeSliceLine(i int) string {
	b, _ := json.Marshal(map[string]interface{}{"index": i, "msg": fmt.Sprintf("line-%d", i)})
	return string(b)
}

func TestSliceFilter_WithPipeline_WindowSlice(t *testing.T) {
	var input []string
	for i := 0; i < 8; i++ {
		input = append(input, makeSliceLine(i))
	}

	sf, err := filter.NewSliceFilter(2, 5)
	if err != nil {
		t.Fatalf("NewSliceFilter: %v", err)
	}

	p := filter.NewPipeline(strings.NewReader(strings.Join(input, "\n")))
	p.AddFilter(sf)

	var buf bytes.Buffer
	if err := p.Run(&buf); err != nil {
		t.Fatalf("pipeline run: %v", err)
	}

	outLines := nonEmptyLines(buf.String())
	if len(outLines) != 3 {
		t.Fatalf("expected 3 lines (indices 2,3,4), got %d: %v", len(outLines), outLines)
	}
	for _, l := range outLines {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(l), &obj); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		idx := int(obj["index"].(float64))
		if idx < 2 || idx >= 5 {
			t.Errorf("unexpected index %d in output", idx)
		}
	}
}

func TestSliceFilter_WithPipeline_OpenEndSlice(t *testing.T) {
	var input []string
	for i := 0; i < 6; i++ {
		input = append(input, makeSliceLine(i))
	}

	sf, err := filter.NewSliceFilter(4, -1)
	if err != nil {
		t.Fatalf("NewSliceFilter: %v", err)
	}

	p := filter.NewPipeline(strings.NewReader(strings.Join(input, "\n")))
	p.AddFilter(sf)

	var buf bytes.Buffer
	if err := p.Run(&buf); err != nil {
		t.Fatalf("pipeline run: %v", err)
	}

	outLines := nonEmptyLines(buf.String())
	if len(outLines) != 2 {
		t.Fatalf("expected 2 lines (indices 4,5), got %d", len(outLines))
	}
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}
