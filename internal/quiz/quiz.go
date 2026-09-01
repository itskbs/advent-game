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

// Ask prints the puzzle's question to out, reads the player's answer
// from in, and reports whether it was correct. A "choice" puzzle
// shows numbered options and expects a pick from 1-4; a "lookup"
// puzzle has no options — the player has to find the answer
// themselves (e.g. by searching online) and type it in.
func Ask(p puzzles.Puzzle, in io.Reader, out io.Writer) (bool, error) {
	fmt.Fprintf(out, "\n🎄 Day %d: %s [%s]\n\n", p.Day, p.Title, p.Difficulty)
	if p.Story != "" {
		fmt.Fprintln(out, p.Story)
		fmt.Fprintln(out)
	}
	fmt.Fprintln(out, p.Question)

	reader := bufio.NewReader(in)

	if p.Type == puzzles.TypeLookup {
		fmt.Fprint(out, "\n🔍 Your answer: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("reading answer: %w", err)
		}
		return strings.EqualFold(strings.TrimSpace(line), strings.TrimSpace(p.Answer)), nil
	}

	for i, choice := range p.Choices {
		fmt.Fprintf(out, "  %d) %s\n", i+1, choice)
	}

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
