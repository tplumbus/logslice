package filter

import (
	"testing"
	"time"
)

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestParseTimeRange_Valid(t *testing.T) {
	tr, err := ParseTimeRange("2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.From.IsZero() || tr.To.IsZero() {
		t.Fatal("expected non-zero bounds")
	}
}

func TestParseTimeRange_OpenBounds(t *testing.T) {
	tr, err := ParseTimeRange("", "2024-06-01T12:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tr.From.IsZero() {
		t.Error("expected From to be zero for open lower bound")
	}
}

func TestParseTimeRange_InvalidFrom(t *testing.T) {
	_, err := ParseTimeRange("not-a-time", "")
	if err == nil {
		t.Fatal("expected error for invalid from timestamp")
	}
}

func TestParseTimeRange_ToBeforeFrom(t *testing.T) {
	_, err := ParseTimeRange("2024-01-02T00:00:00Z", "2024-01-01T00:00:00Z")
	if err == nil {
		t.Fatal("expected error when to is before from")
	}
}

func TestTimeRange_Contains(t *testing.T) {
	tr, _ := ParseTimeRange("2024-03-01T00:00:00Z", "2024-03-31T23:59:59Z")

	cases := []struct {
		ts   time.Time
		want bool
	}{
		{mustTime("2024-03-15T10:00:00Z"), true},
		{mustTime("2024-02-28T23:59:59Z"), false},
		{mustTime("2024-04-01T00:00:00Z"), false},
		{mustTime("2024-03-01T00:00:00Z"), true},
		{mustTime("2024-03-31T23:59:59Z"), true},
	}

	for _, c := range cases {
		got := tr.Contains(c.ts)
		if got != c.want {
			t.Errorf("Contains(%s) = %v, want %v", c.ts.Format(time.RFC3339), got, c.want)
		}
	}
}

func TestTimeRange_IsZero(t *testing.T) {
	var tr TimeRange
	if !tr.IsZero() {
		t.Error("empty TimeRange should be zero")
	}
	tr.From = mustTime("2024-01-01T00:00:00Z")
	if tr.IsZero() {
		t.Error("TimeRange with From set should not be zero")
	}
}
