# Exercises — 03 · Basic types & operators

Run everything from the topic folder (`topics/03-basic-types-and-operators/`, the one with `go.mod`). **Don't open `solutions/` until you've tried.**

| # | Exercise | Trains | Run |
|---|---|---|---|
| E1 | [`01-predict-int-math/`](01-predict-int-math/): write `prediction.go` **before** running anything | Wraparound, `/` `%` `>>` on negatives, exact vs inexact floats, conversions, precedence, `^` and `&^`, strings as bytes (notes §5, §6.1, §6.5) | `make ex E=01-predict-int-math` then `go run ./exercises/01-predict-int-math` |
| E2 | [`02-fix-percentage/`](02-fix-percentage/) | Integer division truncates before the multiply (notes §5.3, §7.6) | `make ex E=02-fix-percentage` |
| E3 | [`03-implement-checked/`](03-implement-checked/) | Overflow-checked `Add`/`Sub`/`Mul`/`Abs`, range-checked conversions, float special values (notes §5.2, §5.4, §5.5, §7.11) | `make ex E=03-implement-checked` |
| E4 | [`04-refactor-invoice/`](04-refactor-invoice/): behaviour tests pass already and must keep passing; `idiom_test.go` fails until the habits are gone | Money as `float64` → integer cents, `math.Ceil` → integer ceiling division, `== true`, getters (notes §8) | `make ex E=04-refactor-invoice` then `go run ./exercises/04-refactor-invoice Keyboard=49.99x2 Mouse=19.95` |
| E5 | [`05-project-ipcalc/`](05-project-ipcalc/) (`cmd/ipcalc` + `internal/ipv4`) | A subnet calculator: `uint32` addresses, masks with shifts, `&` `|` `^`, conversions between sizes, the defined 32-bit shift | `make ex E=05-project-ipcalc` then `go run ./exercises/05-project-ipcalc/cmd/ipcalc 192.168.1.10/24` |

The loop: read the failing test → make it pass → `make check E=<folder>` (gofmt + vet + race) → ask the tutor to `check En`.

Stuck? Say `I'm stuck on E3` (or `hint`): the tutor gives hints one step at a time, never the answer unless you ask for it.
