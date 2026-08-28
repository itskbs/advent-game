// Command advent-game is a small CLI advent calendar: one quiz per
// day, each correct answer earning a letter towards a final solution
// word revealed on Dec 24.
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	cmd := args[0]
	switch {
	case cmd == "status":
		return runStatus()
	case cmd == "solution":
		return runSolution()
	case cmd == "help", cmd == "--help", cmd == "-h":
		printUsage()
		return nil
	case strings.HasPrefix(cmd, "day"):
		day, err := strconv.Atoi(strings.TrimPrefix(cmd, "day"))
		if err != nil {
			return fmt.Errorf("%q is not a valid command; try 'day1'..'day24'", cmd)
		}
		return runDay(day)
	default:
		return fmt.Errorf("unknown command %q; run 'advent-game help'", cmd)
	}
}

func printUsage() {
	fmt.Println(`advent-game — a CLI advent calendar

Usage:
  advent-game day<N>    Play the puzzle for day N (1-24), e.g. "advent-game day1"
  advent-game status    Show which days you've solved so far
  advent-game solution  Reveal the solution word from your collected letters
  advent-game help      Show this message`)
}
