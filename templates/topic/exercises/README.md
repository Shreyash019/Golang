# Exercises

Run everything from the topic folder (the one with `go.mod`). **Don't open `solutions/` until you've tried.**

| # | Exercise | Trains | Run |
|---|---|---|---|
| E1 | `01-predict-<name>/`: write `prediction.go` **before** running anything | | `make ex E=01-predict-<name>` then `go run ./exercises/01-predict-<name>` |
| E2 | `02-fix-<name>/` | | `make ex E=02-fix-<name>` |
| E3 | `03-implement-<name>/` | | `make ex E=03-implement-<name>` |
| E4 | `04-refactor-<name>/`: tests pass already and must keep passing | | `make ex E=04-refactor-<name>` |
| E5 | `05-project-<name>/` (`cmd/` + `internal/`) | | `make ex E=05-project-<name>` |

The loop: read the failing test → make it pass → `make check E=<folder>` (gofmt + vet + race) → ask the tutor to `check En`.
