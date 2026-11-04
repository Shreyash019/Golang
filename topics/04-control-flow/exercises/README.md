# Exercises — 04 · Control Flow

Run everything from the topic folder (`topics/04-control-flow/`, the one with `go.mod`). **Don't open `solutions/` until you've tried.**

| # | Exercise | Trains | Run |
|---|---|---|---|
| E1 | [`01-predict-loops/`](01-predict-loops/): write `prediction.go` **before** running anything | Changing the loop variable (three-clause vs range), `break`/`continue` in a `switch`, `fallthrough`, labeled `continue`, what `range` copies and evaluates once, ranging over a string, per-iteration variables (notes §5.3–5.8, §7.1) | `make ex E=01-predict-loops` then `go run ./exercises/01-predict-loops` |
| E2 | [`02-fix-header/`](02-fix-header/) | `break` inside a `switch` leaves the switch, not the loop (notes §7.1, §5.8) | `make ex E=02-fix-header` |
| E3 | [`03-implement-loops/`](03-implement-loops/) | A tagless `switch` over ranges, the condition-only `for`, a labeled `continue`, ranging over runes vs byte offsets, a loop that carries state (notes §5.2, §5.5, §5.6, §5.8) | `make ex E=03-implement-loops` |
| E4 | [`04-refactor-shipping/`](04-refactor-shipping/): behaviour tests pass already and must keep passing; `idiom_test.go` fails until the habits are gone | Guard clauses instead of arrow code and `else` after `return`, `switch` instead of an if/else-if ladder, `range` instead of index loops, early `return` instead of flag variables, no do-while emulation, no `== true` (notes §8) | `make ex E=04-refactor-shipping` |
| E5 | [`05-project-logstat/`](05-project-logstat/) (`cmd/logstat` + `internal/logstat`) | A log filter CLI: `for sc.Scan()` as a while loop, `continue` for skipped lines, a `switch` with several values per case, `break` on a limit, errors from the scanner and the writer | `make ex E=05-project-logstat` then `go run ./exercises/05-project-logstat/cmd/logstat -level warn < exercises/05-project-logstat/testdata/app.log` |

The loop: read the failing test → make it pass → `make check E=<folder>` (gofmt + vet + race) → ask the tutor to `check En`.

Stuck? Say `I'm stuck on E3` (or `hint`): the tutor gives hints one step at a time, never the answer unless you ask for it.
