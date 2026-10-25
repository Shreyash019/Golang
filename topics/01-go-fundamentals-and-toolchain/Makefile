# Run `make` to list the targets. Every command runs inside this topic's module.
#   N = a worked example name   (make example N=slice-aliasing)
#   E = an exercise folder      (make ex E=03-implement-stack)

.DEFAULT_GOAL := help
.PHONY: help examples example ex check test fmt vet lint verify

help: ## list the targets
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "} {printf "  %-9s %s\n", $$1, $$2}'

examples: ## list the worked examples from the notes
	go run ./cmd/examples

example: ## run one worked example: make example N=<name> (N=all runs them all)
	@test -n "$(N)" || { echo "usage: make example N=<name>"; exit 2; }
	go run ./cmd/examples $(N)

ex: ## run one exercise's tests verbosely: make ex E=<exercise-folder>
	@test -n "$(E)" || { echo "usage: make ex E=<exercise-folder>"; exit 2; }
	go test -race -v ./exercises/$(E)/...

check: ## gofmt + vet + race tests on one exercise (what "check En" runs): make check E=<exercise-folder>
	@test -n "$(E)" || { echo "usage: make check E=<exercise-folder>"; exit 2; }
	@test -z "$$(gofmt -l exercises/$(E))" || { gofmt -l exercises/$(E); echo "gofmt: the files above need formatting (make fmt)"; exit 1; }
	go vet ./exercises/$(E)/...
	go test -race ./exercises/$(E)/...

test: ## test the worked examples and reference solutions (not your exercise attempts)
	go test -race ./cmd/... ./internal/... $(if $(wildcard exercises/solutions),./exercises/solutions/...)

fmt: ## format all code
	gofmt -w .

vet: ## go vet the whole module, starter code included
	go vet ./...

lint: ## golangci-lint on the examples and solutions (if installed)
	golangci-lint run ./cmd/... ./internal/... $(if $(wildcard exercises/solutions),./exercises/solutions/...)

verify: ## the tutor's checks before declaring the topic ready (AGENTS.md §6.4)
	@test -z "$$(gofmt -l .)" || { gofmt -l .; echo "gofmt: the files above need formatting"; exit 1; }
	go vet ./...
	$(MAKE) --no-print-directory test
