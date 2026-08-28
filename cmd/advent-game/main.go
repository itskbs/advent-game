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
	debug, all, reset, debugDay, rest, err := parseFlags(args)
	if err != nil {
		return err
	}

	switch {
	case reset && !debug:
		return fmt.Errorf("--reset requires --debug")
	case reset && (all || debugDay != 0):
		return fmt.Errorf("--reset can't be combined with --day<N> or --all")
	case reset:
		if len(rest) > 0 {
			return fmt.Errorf("--debug --reset doesn't take a separate command")
		}
		return runReset()
	case debugDay != 0 && !debug:
		return fmt.Errorf("--day<N> requires --debug")
	case debugDay != 0:
		if len(rest) > 0 {
			return fmt.Errorf("--debug --day<N> plays that day directly; it doesn't take a separate command")
		}
		return runDay(debugDay, true)
	case all && !debug:
		return fmt.Errorf("--all requires --debug")
	case debug && !all:
		return fmt.Errorf("--debug requires --day<N>, --all, or --reset")
	}

	if len(rest) == 0 {
		printUsage()
		return nil
	}

	cmd := rest[0]
	switch {
	case cmd == "status":
		return runStatus(all)
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
		return runDay(day, all)
	default:
		return fmt.Errorf("unknown command %q; run 'advent-game help'", cmd)
	}
}

// parseFlags pulls the debug-mode flags out of args, returning whatever's
// left as rest. --day<N> and --reset are self-contained (each runs on its
// own), while --all is a modifier that must accompany a command in rest,
// e.g. "--debug --all status" or "--debug --all day12".
func parseFlags(args []string) (debug, all, reset bool, debugDay int, rest []string, err error) {
	for _, a := range args {
		switch {
		case a == "--debug":
			debug = true
		case a == "--all":
			all = true
		case a == "--reset":
			reset = true
		case strings.HasPrefix(a, "--day"):
			n, convErr := strconv.Atoi(strings.TrimPrefix(a, "--day"))
			if convErr != nil {
				return false, false, false, 0, nil, fmt.Errorf("%q is not a valid flag; expected --day<N>", a)
			}
			debugDay = n
		default:
			rest = append(rest, a)
		}
	}
	return debug, all, reset, debugDay, rest, nil
}

func printUsage() {
	fmt.Println(`advent-game — a CLI advent calendar

Usage:
  advent-game day<N>    Play the puzzle for day N (1-24), e.g. "advent-game day1"
  advent-game status    Show which days you've solved so far
  advent-game solution  Reveal the solution word from your collected letters
  advent-game help      Show this message

Debug (for local testing):
  advent-game --debug --day<N>      Play day N immediately, regardless of date
  advent-game --debug --all day<N>  Same, via the day<N> command
  advent-game --debug --all status  Show status with every day treated as unlocked
  advent-game --debug --reset       Wipe all saved progress (unsolves every day)`)
}
