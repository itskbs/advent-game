package main

import (
	"fmt"
	"strings"

	"github.com/kevschroeder99/advent-game/internal/progress"
	"github.com/kevschroeder99/advent-game/internal/puzzles"
)

func runSolution() error {
	prog, err := progress.Load()
	if err != nil {
		return err
	}

	var b strings.Builder
	missing := 0
	for day := puzzles.FirstDay; day <= puzzles.LastDay; day++ {
		if e := prog.Days[day]; e.Solved {
			b.WriteString(e.Letter)
		} else {
			b.WriteString("_")
			missing++
		}
	}

	if missing > 0 {
		fmt.Printf("%s\n\n%d day(s) still unsolved — keep going!\n", b.String(), missing)
		return nil
	}
	fmt.Printf("🎉 %s\n", b.String())
	return nil
}
