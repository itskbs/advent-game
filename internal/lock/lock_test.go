package lock

import (
	"testing"
	"time"
)

func TestIsUnlocked(t *testing.T) {
	cases := []struct {
		name string
		now  string
		day  int
		want bool
	}{
		{"well before December", "2026-08-28T00:00:00Z", 1, false},
		{"day in the future this December", "2026-12-01T00:00:00Z", 5, false},
		{"day is today", "2026-12-05T09:00:00Z", 5, true},
		{"day already passed this December", "2026-12-10T00:00:00Z", 5, true},
		{"last day of the calendar, on the day", "2026-12-24T23:59:00Z", 24, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			now, err := time.Parse(time.RFC3339, tc.now)
			if err != nil {
				t.Fatalf("bad fixture time: %v", err)
			}
			if got := IsUnlocked(tc.day, now); got != tc.want {
				t.Errorf("IsUnlocked(%d, %s) = %v, want %v", tc.day, tc.now, got, tc.want)
			}
		})
	}
}
