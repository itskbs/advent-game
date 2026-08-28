// Package quiz handles presenting a puzzle's question and choices to
// the player and reading back their answer.
package quiz

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/itskbs/advent-game/internal/puzzles"
)

// Ask prints the puzzle's question and numbered choices to out, reads
// the player's pick from in, and reports whether it was correct.
func Ask(p puzzles.Puzzle, in io.Reader, out io.Writer) (bool, error) {
	fmt.Fprintf(out, "\n🎄 Day %d: %s\n\n", p.Day, p.Title)
	fmt.Fprintln(out, p.Question)
	for i, choice := range p.Choices {
		fmt.Fprintf(out, "  %d) %s\n", i+1, choice)
	}

	reader := bufio.NewReader(in)
	for {
		fmt.Fprint(out, "\nYour answer (1-4): ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("reading answer: %w", err)
		}

		choice, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil || choice < 1 || choice > len(p.Choices) {
			fmt.Fprintf(out, "Please enter a number between 1 and %d\n", len(p.Choices))
			continue
		}
		return choice-1 == p.AnswerIndex, nil
	}
}
