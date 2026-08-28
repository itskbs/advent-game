package main

import (
	"fmt"

	"github.com/kevschroeder99/advent-game/internal/lock"
	"github.com/kevschroeder99/advent-game/internal/progress"
	"github.com/kevschroeder99/advent-game/internal/puzzles"
)

func runStatus() error {
	prog, err := progress.Load()
	if err != nil {
		return err
	}

	now := lock.Now()
	for day := puzzles.FirstDay; day <= puzzles.LastDay; day++ {
		switch {
		case prog.Days[day].Solved:
			fmt.Printf("Day %2d: ✅ %s\n", day, prog.Days[day].Letter)
		case lock.IsUnlocked(day, now):
			fmt.Printf("Day %2d: ⬜ not solved yet\n", day)
		default:
			fmt.Printf("Day %2d: 🔒 locked until %s\n", day, lock.UnlockDate(day, now).Format("Jan 2"))
		}
	}
	fmt.Printf("\n%d/%d solved\n", prog.SolvedCount(), puzzles.LastDay)
	return nil
}
