package main

import (
	"fmt"
	"os"
	"time"

	"github.com/itskbs/advent-game/internal/lock"
	"github.com/itskbs/advent-game/internal/progress"
	"github.com/itskbs/advent-game/internal/puzzles"
	"github.com/itskbs/advent-game/internal/quiz"
)

func runDay(day int) error {
	if day < puzzles.FirstDay || day > puzzles.LastDay {
		return fmt.Errorf("day must be between %d and %d, got %d", puzzles.FirstDay, puzzles.LastDay, day)
	}

	now := lock.Now()
	if !lock.IsUnlocked(day, now) {
		fmt.Printf("🔒 Day %d isn't open yet. Come back on %s.\n", day, lock.UnlockDate(day, now).Format("Jan 2"))
		return nil
	}

	p, err := puzzles.Load(day)
	if err != nil {
		return err
	}

	correct, err := quiz.Ask(p, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	if !correct {
		fmt.Printf("❌ Not quite — try again anytime with `advent-game day%d`.\n", day)
		return nil
	}

	prog, err := progress.Load()
	if err != nil {
		return err
	}
	prog.MarkSolved(day, p.Letter, time.Now())
	if err := prog.Save(); err != nil {
		return err
	}

	fmt.Printf("✅ Correct! You earned the letter %q. (%d/%d days solved)\n", p.Letter, prog.SolvedCount(), puzzles.LastDay)
	return nil
}
