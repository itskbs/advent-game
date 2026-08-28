package main

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/itskbs/advent-game/internal/progress"
)

func TestParseFlags(t *testing.T) {
	cases := []struct {
		name          string
		args          []string
		wantDebug     bool
		wantAll       bool
		wantReset     bool
		wantDebugDay  int
		wantRest      []string
		wantErrSubstr string
	}{
		{
			name:     "plain command, no flags",
			args:     []string{"status"},
			wantRest: []string{"status"},
		},
		{
			name:         "debug day",
			args:         []string{"--debug", "--day5"},
			wantDebug:    true,
			wantDebugDay: 5,
		},
		{
			name:      "debug all with a command",
			args:      []string{"--debug", "--all", "status"},
			wantDebug: true,
			wantAll:   true,
			wantRest:  []string{"status"},
		},
		{
			name:      "debug reset",
			args:      []string{"--debug", "--reset"},
			wantDebug: true,
			wantReset: true,
		},
		{
			name:          "invalid day flag",
			args:          []string{"--daynope"},
			wantErrSubstr: `"--daynope" is not a valid flag`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			debug, all, reset, debugDay, rest, err := parseFlags(tc.args)
			if tc.wantErrSubstr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrSubstr) {
					t.Fatalf("parseFlags(%v) error = %v, want substring %q", tc.args, err, tc.wantErrSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseFlags(%v) unexpected error: %v", tc.args, err)
			}
			if debug != tc.wantDebug || all != tc.wantAll || reset != tc.wantReset || debugDay != tc.wantDebugDay {
				t.Errorf("parseFlags(%v) = (debug=%v, all=%v, reset=%v, debugDay=%v), want (%v, %v, %v, %v)",
					tc.args, debug, all, reset, debugDay, tc.wantDebug, tc.wantAll, tc.wantReset, tc.wantDebugDay)
			}
			if !reflect.DeepEqual(rest, tc.wantRest) {
				t.Errorf("parseFlags(%v) rest = %v, want %v", tc.args, rest, tc.wantRest)
			}
		})
	}
}

func TestRunDebugModeGuards(t *testing.T) {
	cases := []struct {
		name          string
		args          []string
		wantErrSubstr string
	}{
		{"day flag without --debug", []string{"--day5"}, "requires --debug"},
		{"--all without --debug", []string{"--all", "status"}, "requires --debug"},
		{"--reset without --debug", []string{"--reset"}, "requires --debug"},
		{"--debug without --day, --all, or --reset", []string{"--debug"}, "requires --day<N>, --all, or --reset"},
		{"--debug --day<N> with a trailing command", []string{"--debug", "--day5", "status"}, "doesn't take a separate command"},
		{"--reset combined with --day<N>", []string{"--debug", "--reset", "--day5"}, "can't be combined"},
		{"--reset combined with --all", []string{"--debug", "--reset", "--all"}, "can't be combined"},
		{"--reset with a trailing command", []string{"--debug", "--reset", "status"}, "doesn't take a separate command"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := run(tc.args)
			if err == nil || !strings.Contains(err.Error(), tc.wantErrSubstr) {
				t.Fatalf("run(%v) error = %v, want substring %q", tc.args, err, tc.wantErrSubstr)
			}
		})
	}
}

func TestRunReset(t *testing.T) {
	t.Setenv("ADVENT_GAME_HOME", t.TempDir())

	prog, err := progress.Load()
	if err != nil {
		t.Fatalf("progress.Load() failed: %v", err)
	}
	prog.MarkSolved(1, "A", time.Now())
	if err := prog.Save(); err != nil {
		t.Fatalf("progress.Save() failed: %v", err)
	}

	if err := run([]string{"--debug", "--reset"}); err != nil {
		t.Fatalf("run([--debug --reset]) failed: %v", err)
	}

	after, err := progress.Load()
	if err != nil {
		t.Fatalf("progress.Load() after reset failed: %v", err)
	}
	if got := after.SolvedCount(); got != 0 {
		t.Errorf("SolvedCount() after reset = %d, want 0", got)
	}
}
