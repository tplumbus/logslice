package filter

import (
	"strings"
	"testing"
)

func TestTypeFilter_WithPipeline_FiltersStringFields(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"start"}`,
		`{"level":42,"msg":"bad level type"}`,
		`{"level":"warn","msg":"low disk"}`,
		`{"level":null,"msg":"null level"}`,
		`{"msg":"no level field"}`,
	}

	tf, err := NewTypeFilter("level", "string")
	if err != nil {
		t.Fatalf("NewTypeFilter: %v", err)
	}

	pipe := NewPipeline(tf)
	var got []string
	err = pipe.Run(lines, func(line string) error {
		got = append(got, line)
		return nil
	})
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
	for _, line := range got {
		if !strings.Contains(line, `"level":"`) {
			t.Errorf("expected string level field in: %s", line)
		}
	}
}

func TestTypeFilter_WithPipeline_ParsedFromFlag(t *testing.T) {
	lines := []string{
		`{"count":1,"msg":"one"}`,
		`{"count":"two","msg":"two"}`,
		`{"count":3,"msg":"three"}`,
	}

	tf, err := ParseTypeFlag("count:number")
	if err != nil {
		t.Fatalf("ParseTypeFlag: %v", err)
	}

	pipe := NewPipeline(tf)
	var got []string
	err = pipe.Run(lines, func(line string) error {
		got = append(got, line)
		return nil
	})
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
	if !strings.Contains(got[0], `"count":1`) {
		t.Errorf("unexpected first line: %s", got[0])
	}
	if !strings.Contains(got[1], `"count":3`) {
		t.Errorf("unexpected second line: %s", got[1])
	}
}
