# Exercises — 05 · Functions, Closures & Defer

Run everything from the topic folder (`topics/05-functions-closures-and-defer/`, the one with `go.mod`). **Don't open `solutions/` until you've tried.**

| # | Exercise | Trains | Run |
|---|---|---|---|
| E1 | [`01-predict-defer/`](01-predict-defer/): write `prediction.go` **before** running anything | `defer` order and argument evaluation, deferred closures and named results, closures capturing variables, variadic aliasing, `append` in a callee, composing func values (notes §5.1–5.7) | `make ex E=01-predict-defer` then `go run ./exercises/01-predict-defer` |
| E2 | [`02-fix-linecount/`](02-fix-linecount/) | `defer` in a loop holds every file open until the function returns (notes §5.7, §7.4) | `make ex E=02-fix-linecount` |
| E3 | [`03-implement-funcs/`](03-implement-funcs/) | Multiple results and a sentinel error, variadics, returning a closure, closures sharing state, a higher-order retry with `%w`, recursion over a tree (notes §5.2–5.8) | `make ex E=03-implement-funcs` |
| E4 | [`04-refactor-checkout/`](04-refactor-checkout/): behaviour tests pass already and must keep passing; `idiom_test.go` fails until the habits are gone | A strategy interface + `Impl` types + factory → func values and closures; `panic`/`recover` as exceptions → returned errors; an out-parameter → `(value, error)` (notes §8) | `make ex E=04-refactor-checkout` then `go run ./exercises/04-refactor-checkout -coupon SAVE10 1999 500` |
| E5 | [`05-project-calc/`](05-project-calc/) (`cmd/calc` + `internal/expr`) | An expression calculator: a recursive-descent parser from mutually recursive closures, errors with positions, a sentinel error, and a `-trace` mode built on `defer` + named results | `make ex E=05-project-calc` then `go run ./exercises/05-project-calc/cmd/calc -trace '2 * (3 + 4)'` |

The loop: read the failing test → make it pass → `make check E=<folder>` (gofmt + vet + race) → ask the tutor to `check En`.

Stuck? Say `I'm stuck on E3` (or `hint`): the tutor gives hints one step at a time, never the answer unless you ask for it.
