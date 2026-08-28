package main

import (
	"fmt"

	"github.com/itskbs/advent-game/internal/progress"
)

func runReset() error {
	if err := progress.Reset(); err != nil {
		return err
	}
	fmt.Println("🔄 Progress reset — every day is unsolved again.")
	return nil
}
