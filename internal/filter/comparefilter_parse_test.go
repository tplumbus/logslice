package filter

import (
	"testing"
)

func TestParseCompareFlag_Valid(t *testing.T) {
	cases := []struct {
		input string
		op    CompareOp
		value float64
	}{
		{"latency:gte:100", OpGte, 100},
		{"code:eq:200", OpEq, 200},
		{"score:lt:0.5", OpLt, 0.5},
		{"count:ne:0", OpNe, 0},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			f, err := ParseCompareFlag(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f.op != tc.op {
				t.Errorf("expected op %q, got %q", tc.op, f.op)
			}
			if f.value != tc.value {
				t.Errorf("expected value %v, got %v", tc.value, f.value)
			}
		})
	}
}

func TestParseCompareFlag_MissingParts(t *testing.T) {
	_, err := ParseCompareFlag("latency:gte")
	if err == nil {
		t.Fatal("expected error for missing value part")
	}
}

func TestParseCompareFlag_BadValue(t *testing.T) {
	_, err := ParseCompareFlag("latency:gte:notanumber")
	if err == nil {
		t.Fatal("expected error for non-numeric value")
	}
}

func TestParseCompareFlag_EmptyField(t *testing.T) {
	_, err := ParseCompareFlag(":gte:100")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseCompareFlag_UnknownOp(t *testing.T) {
	_, err := ParseCompareFlag("x:between:5")
	if err == nil {
		t.Fatal("expected error for unknown operator")
	}
}
