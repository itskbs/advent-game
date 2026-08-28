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
  "question": "The question text",
  "choices": ["A", "B", "C", "D"],
  "answerIndex": 0,
  "letter": "A"
}
```

Only `day01.json` has real content so far — days 2-24 are `TODO`
placeholders (`letter: "?"`). Fill them in, and choose the 24 letters
so that, read in order (day 1 → day 24), they spell your final
solution word or phrase.

### Testing the date lock without waiting for December

Set `ADVENT_GAME_NOW` to a fixed RFC3339 timestamp to pretend it's a
different date, and `ADVENT_GAME_HOME` to sandbox progress storage:

```sh
ADVENT_GAME_NOW=2026-12-05T09:00:00Z ADVENT_GAME_HOME=/tmp/advent-game-dev \
  go run ./cmd/advent-game day5
```

### Releasing

Tagging a version (`git tag v0.1.0 && git push --tags`) triggers
`.github/workflows/release.yml`, which runs
[GoReleaser](https://goreleaser.com) to:

1. Cross-compile binaries for macOS/Linux (amd64 + arm64)
2. Publish a GitHub release with those binaries attached
3. Push an updated formula to the `itskbs/homebrew-advent-game`
   tap repo

One-time setup before the first release:

- Create the `itskbs/homebrew-advent-game` repo (empty is fine
  — GoReleaser creates `Formula/advent-game.rb` in it).
- Add a repo secret `HOMEBREW_TAP_GITHUB_TOKEN`: a GitHub PAT with
  write access to that tap repo.
