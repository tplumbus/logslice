package filter

import (
	"encoding/json"
	"fmt"
	"testing"
)

func makeLogLine(t *testing.T, fields map[string]interface{}) string {
	t.Helper()
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(b)
}

func TestNewSortFilter_Valid(t *testing.T) {
	_, err := NewSortFilter("level", SortAsc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSortFilter_EmptyField(t *testing.T) {
	_, err := NewSortFilter("", SortAsc)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestSortFilter_Asc(t *testing.T) {
	sf, _ := NewSortFilter("level", SortAsc)
	levels := []string{"warn", "error", "debug", "info"}
	for _, l := range levels {
		sf.Observe(fmt.Sprintf(`{"level":%q}`, l))
	}
	got := sf.Sorted()
	want := []string{"debug", "error", "info", "warn"}
	for i, w := range want {
		var m map[string]interface{}
		json.Unmarshal([]byte(got[i]), &m)
		if m["level"] != w {
			t.Errorf("pos %d: got %v want %v", i, m["level"], w)
		}
	}
}

func TestSortFilter_Desc(t *testing.T) {
	sf, _ := NewSortFilter("level", SortDesc)
	levels := []string{"warn", "error", "debug", "info"}
	for _, l := range levels {
		sf.Observe(fmt.Sprintf(`{"level":%q}`, l))
	}
	got := sf.Sorted()
	want := []string{"warn", "info", "error", "debug"}
	for i, w := range want {
		var m map[string]interface{}
		json.Unmarshal([]byte(got[i]), &m)
		if m["level"] != w {
			t.Errorf("pos %d: got %v want %v", i, m["level"], w)
		}
	}
}

func TestSortFilter_MissingField_SortsToFront(t *testing.T) {
	sf, _ := NewSortFilter("level", SortAsc)
	sf.Observe(`{"msg":"no level"}`)
	sf.Observe(`{"level":"warn","msg":"b"}`)
	got := sf.Sorted()
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
	var m map[string]interface{}
	json.Unmarshal([]byte(got[0]), &m)
	if _, ok := m["level"]; ok {
		t.Error("expected missing-field line first in asc order")
	}
}

func TestSortFilter_Reset(t *testing.T) {
	sf, _ := NewSortFilter("level", SortAsc)
	sf.Observe(`{"level":"info"}`)
	sf.Reset()
	if got := sf.Sorted(); len(got) != 0 {
		t.Errorf("expected empty after reset, got %d lines", len(got))
	}
}
