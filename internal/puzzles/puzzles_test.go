package puzzles

import "testing"

func TestLoadAllDaysAreWellFormed(t *testing.T) {
	for day := FirstDay; day <= LastDay; day++ {
		p, err := Load(day)
		if err != nil {
			t.Fatalf("day %d: %v", day, err)
		}
		if p.Day != day {
			t.Errorf("day %d: JSON \"day\" field is %d, want %d", day, p.Day, day)
		}
		if p.Letter == "" {
			t.Errorf("day %d: missing reward letter", day)
		}
	}
}

func TestLoadRejectsOutOfRangeDay(t *testing.T) {
	for _, day := range []int{0, -1, 25} {
		if _, err := Load(day); err == nil {
			t.Errorf("Load(%d) should have returned an error", day)
		}
	}
}
