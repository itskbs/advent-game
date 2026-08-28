// Package progress persists which days the player has solved and the
// letters they've earned, in a small JSON file under the user's
// config directory.
package progress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry records the outcome of one solved day.
type Entry struct {
	Solved   bool      `json:"solved"`
	Letter   string    `json:"letter"`
	SolvedAt time.Time `json:"solvedAt"`
}

// Progress is the full save file: one entry per solved day, keyed by
// day number.
type Progress struct {
	Days map[int]Entry `json:"days"`
}

func empty() Progress {
	return Progress{Days: map[int]Entry{}}
}

// Path returns the location of the progress file. ADVENT_GAME_HOME
// overrides the directory (used in tests); otherwise it lives under
// the user's standard config directory.
func Path() (string, error) {
	if v := os.Getenv("ADVENT_GAME_HOME"); v != "" {
		return filepath.Join(v, "progress.json"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolving config dir: %w", err)
	}
	return filepath.Join(dir, "advent-game", "progress.json"), nil
}

// Load reads the progress file, returning an empty Progress if none
// exists yet (i.e. the player hasn't solved anything).
func Load() (Progress, error) {
	path, err := Path()
	if err != nil {
		return Progress{}, err
	}

	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return empty(), nil
	}
	if err != nil {
		return Progress{}, fmt.Errorf("reading progress file: %w", err)
	}

	var p Progress
	if err := json.Unmarshal(raw, &p); err != nil {
		return Progress{}, fmt.Errorf("parsing progress file: %w", err)
	}
	if p.Days == nil {
		p.Days = map[int]Entry{}
	}
	return p, nil
}

// Save writes the progress file, creating its parent directory if needed.
func (p Progress) Save() error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	raw, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding progress file: %w", err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return fmt.Errorf("writing progress file: %w", err)
	}
	return nil
}

// Reset wipes all saved progress, overwriting the save file with an
// empty result. Used by the CLI's --debug --reset flag for local testing.
func Reset() error {
	return empty().Save()
}

// MarkSolved records day as solved with the given reward letter,
// overwriting any previous result for that day (replaying a day is
// allowed and simply re-confirms the letter).
func (p *Progress) MarkSolved(day int, letter string, at time.Time) {
	p.Days[day] = Entry{Solved: true, Letter: letter, SolvedAt: at}
}

// SolvedCount returns how many days have been solved so far.
func (p Progress) SolvedCount() int {
	n := 0
	for _, e := range p.Days {
		if e.Solved {
			n++
		}
	}
	return n
}
