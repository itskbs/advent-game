// Package puzzles loads the advent calendar's daily quizzes. All 24
// days are baked into the binary at compile time via go:embed, so the
// tool has zero runtime dependencies and works fully offline.
package puzzles

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed data/*.json
var dataFS embed.FS

// FirstDay and LastDay bound the advent calendar (Dec 1 - Dec 24).
const (
	FirstDay = 1
	LastDay  = 24
)

// Puzzle is a single day's quiz: a question, four choices, which one
// is correct, and the letter awarded for a correct answer.
type Puzzle struct {
	Day         int      `json:"day"`
	Title       string   `json:"title"`
	Question    string   `json:"question"`
	Choices     []string `json:"choices"`
	AnswerIndex int      `json:"answerIndex"`
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
	if len(p.Choices) != 4 {
		return Puzzle{}, fmt.Errorf("puzzle for day %d must have exactly 4 choices, got %d", day, len(p.Choices))
	}
	if p.AnswerIndex < 0 || p.AnswerIndex > 3 {
		return Puzzle{}, fmt.Errorf("puzzle for day %d has an invalid answerIndex %d", day, p.AnswerIndex)
	}
	return p, nil
}
