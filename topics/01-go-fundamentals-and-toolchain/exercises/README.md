# Exercises — 01 · Go fundamentals & toolchain

Run everything from the topic folder (`topics/01-go-fundamentals-and-toolchain/`, the one with `go.mod`). **Don't open `solutions/` until you've tried.**

| # | Exercise | Trains | Run |
|---|---|---|---|
| E1 | [`01-predict-init-order/`](01-predict-init-order/): write `prediction.go` **before** running anything | Program initialization order: variables by dependency, files by name, then `init`s, then `main` (notes §5.3) | `make ex E=01-predict-init-order` then `go run ./exercises/01-predict-init-order` |
| E2 | [`02-fix-init-trap/`](02-fix-init-trap/) | The "variables are initialized before `init()` runs" gotcha (notes §7.6) | `make ex E=02-fix-init-trap` |
| E3 | [`03-implement-gomod/`](03-implement-gomod/) | Reading `go.mod` (notes §5.5); a table-tested parser with line-numbered errors | `make ex E=03-implement-gomod` |
| E4 | [`04-refactor-config/`](04-refactor-config/): behaviour tests pass already and must keep passing; `idiom_test.go` fails until the Java habits are gone | `main` + `run()` structure, no getters/`Impl`/`utils`/panics (notes §8) | `make ex E=04-refactor-config` |
| E5 | [`05-project-pkgscan/`](05-project-pkgscan/) (`cmd/pkgscan` + `internal/scan` + `internal/report`) | Packages, `main`, imports, doc comments and the go command's directory rules, in a real CLI | `make ex E=05-project-pkgscan` then `go run ./exercises/05-project-pkgscan/cmd/pkgscan .` |

The loop: read the failing test → make it pass → `make check E=<folder>` (gofmt + vet + race) → ask the tutor to `check En`.

Stuck? Say `I'm stuck on E3` (or `hint`): the tutor gives hints one step at a time, never the answer unless you ask for it.
