// Package lock decides which days of the advent calendar are open to
// play, so the game stays true to a real advent calendar: a day
// unlocks on its date and stays open (replayable) forever after.
package lock

import (
	"os"
	"time"
)

// Now returns the current time, or a fake time from ADVENT_GAME_NOW
// (RFC3339, e.g. "2026-12-05T00:00:00Z") when set. This exists purely
// so the lock logic can be exercised in tests and during development
// without waiting for December.
func Now() time.Time {
	if v := os.Getenv("ADVENT_GAME_NOW"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}
	}
	return time.Now()
}

// UnlockDate returns the date on which the given day unlocks:
// December <day> of the current year, at local midnight.
func UnlockDate(day int, now time.Time) time.Time {
	return time.Date(now.Year(), time.December, day, 0, 0, 0, 0, now.Location())
}

// IsUnlocked reports whether the given day can be played right now.
// Today and any past day are always playable (and replayable); days
// later in December stay locked until their date arrives.
func IsUnlocked(day int, now time.Time) bool {
	return !now.Before(UnlockDate(day, now))
}
