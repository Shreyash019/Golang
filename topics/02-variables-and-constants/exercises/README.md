# Exercises — 02 · Variables & constants

Run everything from the topic folder (`topics/02-variables-and-constants/`, the one with `go.mod`). **Don't open `solutions/` until you've tried.**

| # | Exercise | Trains | Run |
|---|---|---|---|
| E1 | [`01-predict-scope-and-iota/`](01-predict-scope-and-iota/): write `prediction.go` **before** running anything | iota and implicit repetition, where a variable's scope starts, shadowing, `:=` redeclaration, untyped constants and default types (notes §5.2–5.7) | `make ex E=01-predict-scope-and-iota` then `go run ./exercises/01-predict-scope-and-iota` |
| E2 | [`02-fix-shadowing/`](02-fix-shadowing/) | The `:=`-in-an-inner-block shadowing bug that compiles and passes `go vet` (notes §7.1) | `make ex E=02-fix-shadowing` |
| E3 | [`03-implement-bytesize/`](03-implement-bytesize/) | Unit constants with iota, untyped constants in int64 and float64 arithmetic, narrow scopes (notes §5.3, §5.5, §6.4) | `make ex E=03-implement-bytesize` |
| E4 | [`04-refactor-backoff/`](04-refactor-backoff/): behaviour tests pass already and must keep passing; `idiom_test.go` fails until the habits are gone | Java/Python "constants class" habits → `const` blocks, MixedCaps, `time.Duration` constants, no getters, no global state, no explicit zero initializers (notes §8) | `make ex E=04-refactor-backoff` then `go run ./exercises/04-refactor-backoff 3` |
| E5 | [`05-project-envcheck/`](05-project-envcheck/) (`cmd/envcheck` + `internal/config`) | A real config loader: default constants, an iota enum whose zero value means "unset", parsing with `:=` and `=` without shadowing | `make ex E=05-project-envcheck` then `APP_PORT=9090 go run ./exercises/05-project-envcheck/cmd/envcheck` |

The loop: read the failing test → make it pass → `make check E=<folder>` (gofmt + vet + race) → ask the tutor to `check En`.

Stuck? Say `I'm stuck on E3` (or `hint`): the tutor gives hints one step at a time, never the answer unless you ask for it.
