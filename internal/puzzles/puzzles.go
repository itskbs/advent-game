// Package puzzles loads the advent calendar's daily quizzes. All 24
// days are baked into the binary at compile time via go:embed, so the
// tool has zero runtime dependencies and works fully offline.
package puzzles

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed data/*.json
var dataFS embed.FS

// FirstDay and LastDay bound the advent calendar (Dec 1 - Dec 24).
const (
	FirstDay = 1
	LastDay  = 24
)

// TypeChoice is a standard four-option multiple-choice puzzle. It's
// the default when a puzzle's "type" field is left empty, so existing
// puzzle files don't need to be touched to keep working.
const TypeChoice = "choice"

// TypeLookup is a puzzle with no given choices: the player has to
// research the answer themselves (e.g. on the internet) and type it
// in, and it's checked against Answer (case-insensitively, trimmed).
const TypeLookup = "lookup"

// Puzzle is a single day's quiz: a story beat continuing the advent
// calendar's narrative, and a question of a given difficulty and
// type. A "choice" puzzle offers four choices and an AnswerIndex; a
// "lookup" puzzle instead has a free-text Answer the player must
// find and type in themselves. Either way, a correct answer earns
// Letter.
type Puzzle struct {
	Day         int      `json:"day"`
	Title       string   `json:"title"`
	Story       string   `json:"story"`
	Difficulty  string   `json:"difficulty"`
	Type        string   `json:"type"`
	Question    string   `json:"question"`
	Choices     []string `json:"choices,omitempty"`
	AnswerIndex int      `json:"answerIndex,omitempty"`
	Answer      string   `json:"answer,omitempty"`
	Letter      string   `json:"letter"`
}

// Load reads and validates the puzzle for the given day (1-24).
func Load(day int) (Puzzle, error) {
	if day < FirstDay || day > LastDay {
		return Puzzle{}, fmt.Errorf("day must be between %d and %d, got %d", FirstDay, LastDay, day)
	}

	name := fmt.Sprintf("data/day%02d.json", day)
	raw, err := dataFS.ReadFile(name)
	if err != nil {
		return Puzzle{}, fmt.Errorf("no puzzle found for day %d: %w", day, err)
	}

	var p Puzzle
	if err := json.Unmarshal(raw, &p); err != nil {
		return Puzzle{}, fmt.Errorf("puzzle for day %d is malformed: %w", day, err)
	}
	if p.Type == "" {
		p.Type = TypeChoice
	}

	switch p.Type {
	case TypeChoice:
		if len(p.Choices) != 4 {
			return Puzzle{}, fmt.Errorf("puzzle for day %d must have exactly 4 choices, got %d", day, len(p.Choices))
		}
		if p.AnswerIndex < 0 || p.AnswerIndex > 3 {
			return Puzzle{}, fmt.Errorf("puzzle for day %d has an invalid answerIndex %d", day, p.AnswerIndex)
		}
	case TypeLookup:
		if strings.TrimSpace(p.Answer) == "" {
			return Puzzle{}, fmt.Errorf("puzzle for day %d is type %q but has no answer", day, TypeLookup)
		}
	default:
		return Puzzle{}, fmt.Errorf("puzzle for day %d has unknown type %q; want %q or %q", day, p.Type, TypeChoice, TypeLookup)
	}
	return p, nil
}
