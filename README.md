# advent-game

A CLI advent calendar: one quiz per day, four choices, one correct
answer, one letter earned. Collect all 24 and reveal the solution word
on Dec 24.

## Playing

```sh
advent-game day1        # play today's (or an earlier) puzzle
advent-game status      # see which days you've solved
advent-game solution    # reveal the word from your collected letters
advent-game help
```

Each day unlocks on its own date (`day5` opens on Dec 5) and stays
open — and replayable — afterwards. Progress is stored locally at
`~/Library/Application Support/advent-game/progress.json` (or the
platform equivalent of Go's `os.UserConfigDir`).

## Installing

```sh
brew tap itskbs/advent-game
brew install advent-game
```

## Development

Requires Go 1.27+. No external dependencies — everything is stdlib.

```sh
go run ./cmd/advent-game day1
go test ./...
```

### Puzzle content

Each day's quiz lives in `internal/puzzles/data/dayNN.json` and is
baked into the binary at build time via `go:embed`. Format:

```json
{
  "day": 1,
  "title": "Short title shown above the question",
  "story": "A sentence or two continuing that day's story beat",
  "difficulty": "easy | medium | hard",
  "type": "choice",
  "question": "The question text",
  "choices": ["A", "B", "C", "D"],
  "answerIndex": 0,
  "letter": "A"
}
```

`letter` can be a space (`" "`) instead of an actual letter, to spell
a multi-word solution phrase — `day.go` shows a friendlier message
("you earned a space in the solution phrase") for that case instead
of printing a quoted blank.

`type` is optional and defaults to `"choice"` (a normal four-option
quiz) if omitted, so existing puzzle files don't need to change. Set
it to `"lookup"` for a puzzle with no given choices — the player has
to research the answer themselves (e.g. search it up online) and
type it in, checked case-insensitively against `answer`:

```json
{
  "day": 15,
  "title": "...",
  "story": "...",
  "difficulty": "hard",
  "type": "lookup",
  "question": "A question whose answer has to be looked up, not guessed",
  "answer": "1933",
  "letter": "L"
}
```

A `"lookup"` puzzle has no `choices`/`answerIndex` — just `answer`.
Day 15 is a real example of this.

#### The story so far

**"The Case of the Silent Sleigh Bell."** The enchanted Sleigh Bell —
the one that wakes itself every Christmas Eve and gives the sleigh
its lift — has gone missing from the workshop vault. You're a
brand-new junior elf sent to investigate. The trail leads to Jingle,
a young reindeer calf who accidentally shattered the Bell mid-jump
while practicing tricks, scattering its chimes everywhere someone
nearby had just felt real Christmas cheer. Each day's trial earns
enough trust/clues to keep the investigation moving, until the last
chime turns up on Christmas Eve and the Bell rings again.

Days 1, 8, 15, 20, and 24 are fully written and mark the story's
current beats (intro → first clue → confession → penultimate clue →
finale) — use them as the template and tone for the remaining days.
Difficulty is deliberately mixed rather than escalating, so pick
whatever fits each day's question, not its position in the calendar.
Day 15 is also the one `"lookup"`-type example — mix in more of
those wherever a question is better answered by research than by
picking from four options.

The solution phrase is **"THE SLEIGH BELL IS BACK!"** (24 characters,
including 4 spaces and a trailing `!`) — every day's `letter` field
already has the correct character set below, even for still-`TODO`
days, so the mechanics work end-to-end today; only `title`,
`story`, `difficulty`, `question`, and `choices` are left to write in:

| Day | Letter | Day | Letter | Day | Letter | Day | Letter |
|----:|:------:|----:|:------:|----:|:------:|----:|:------:|
| 1 | `T` | 7 | `E` | 13 | `E` | 19 | `" "` |
| 2 | `H` | 8 | `I` ✅ | 14 | `L` | 20 | `B` ✅ |
| 3 | `E` | 9 | `G` | 15 | `L` ✅ | 21 | `A` |
| 4 | `" "` | 10 | `H` | 16 | `" "` | 22 | `C` |
| 5 | `S` | 11 | `" "` | 17 | `I` | 23 | `K` |
| 6 | `L` | 12 | `B` | 18 | `S` | 24 | `!` ✅ |

(✅ = fully written; the rest still have `TODO` title/story/question/
choices, but the correct letter already in place.)

### Testing the date lock without waiting for December

Set `ADVENT_GAME_NOW` to a fixed RFC3339 timestamp to pretend it's a
different date, and `ADVENT_GAME_HOME` to sandbox progress storage:

```sh
ADVENT_GAME_NOW=2026-12-05T09:00:00Z ADVENT_GAME_HOME=/tmp/advent-game-dev \
  go run ./cmd/advent-game day5
```

For a quicker one-off check of a specific day (or every day), skip the
env vars and use the debug flags instead:

```sh
go run ./cmd/advent-game --debug --day12       # play day 12 right now, lock or no lock
go run ./cmd/advent-game --debug --all status  # status view with every day shown unlocked
go run ./cmd/advent-game --debug --all day20   # same idea, via the day<N> command
```

### Releasing

Tagging a version (`git tag v0.1.0 && git push --tags`) triggers
`.github/workflows/release.yml`, which runs
[GoReleaser](https://goreleaser.com) to:

1. Cross-compile binaries for macOS/Linux (amd64 + arm64)
2. Publish a GitHub release with those binaries attached
3. Push an updated formula to the `itskbs/homebrew-advent-game`
   tap repo
4. Tag and create a matching release on the tap repo itself, so its
   version history mirrors this repo's

One-time setup before the first release:

- Create the `itskbs/homebrew-advent-game` repo (empty is fine
  — GoReleaser creates `Formula/advent-game.rb` in it).
- Add a repo secret `HOMEBREW_TAP_GITHUB_TOKEN`: a GitHub PAT with
  write access to that tap repo.
