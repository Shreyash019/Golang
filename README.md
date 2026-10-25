# Go — The Complete, Test-Driven Guide (for Python & Java Engineers)

A guide to writing **production-quality, idiomatic Go** and understanding **why Go works the way it does**: the runtime, memory and concurrency, not just the syntax. It's written for engineers who already know Python or Java, so every topic contrasts Go with them and calls out the habits that don't translate (inheritance thinking, exceptions, getters/setters, big interfaces, over-abstraction).

Each topic is a folder with:
- **deep-dive notes**: internals with diagrams, a runnable example for every concept with its real output, gotchas, idiomatic vs non-idiomatic code, and production and interview angles;
- **test-driven exercises** that build up from *predict the output* to *a mini-project*, with reference solutions.

> This repo is also my own learning log: the progress columns below and the 🗣️ Q&A section at the end of each topic file record my sessions with an AI tutor.

---

## 📚 How to use this repo

### Suggested order
1. **Phase 0: diagnose.** If you've already been learning Go for a while, start with `diagnose me` rather than at topic 01. It finds out *why* you get stuck and which topics to skim or go deep on.
2. Then go through the phases in order:

| Phase | What you get |
|---|---|
| 1. Getting Started | The toolchain, values, types, control flow and functions, with the right mental models from day one |
| 2. Data & Memory Basics | Arrays, slices, maps, strings and pointers: where most "Go is confusing" moments come from |
| 3. Types & Abstraction | Go's answer to OOP: structs, methods, embedding, interfaces, the type system, generics, packages, API design |
| 4. Errors, Testing & Debugging | Errors as values, then the feedback loop: tests, test doubles, Delve |
| 5. Concurrency | Goroutines, channels, `sync`, data races, `context`, patterns |
| 6. Runtime & Performance | Stack/heap, escape analysis, GC, profiling, reflection, `unsafe` |
| 7. Standard Library & I/O | Everyday packages, files and streams, JSON, HTTP |
| 8. Idioms, Patterns & Production | Functional and design patterns the Go way, idiomatic Go, production services, a capstone, mock interviews |

3. Within a phase, the topic number is the suggested order. Each topic file lists its *Prerequisites* at the top.

### Why topics are grouped
Many Go concepts are tightly linked: closures and `defer`, arrays and slices, interfaces and polymorphism, packages and encapsulation. Learning them as separate items is a common way to get confused, so related items share one topic, and the **Covers** column shows which items of the original 50-item syllabus (listed in [Syllabus coverage](#-syllabus-coverage)) each topic includes.

Every concept also has **exactly one home**, where it's taught in depth. Other topics give a short recap and link back, so you never get two different explanations of the same thing:

| Concept | Home topic | Also used in (recap + link only) |
|---|---|---|
| `go build/run/test/vet/doc`, gofmt (the basics) | 01 | 18 (testing in depth), 28 (benchmarks/pprof) |
| Scope & shadowing | 02 | 04, 17 (the `err` shadowing trap) |
| Defined types vs aliases | 13 (03 gives a first look) | 10, 14 |
| Closures, `defer` | 05 | 17 (`recover`), 23 (`defer mu.Unlock()`), 26, 34 |
| Value vs pointer semantics | 09 | 06, 10, 27 |
| Method sets | 10 | 12 |
| Interface values (type + data words), the nil-interface trap | 12 | 13, 17 |
| Exported vs unexported, `internal/` | 15 | 16, 37 |
| Interfaces as dependency boundaries / dependency injection | 16 | 19, 35, 37 |
| Benchmarks: writing them | 18 | 28 (reading them, optimising) |
| Stack vs heap, escape analysis, GC | 27 (09 gives an intro) | 28 |
| "Memory model" ⚠️ two meanings | **24**: the *Go memory model* (happens-before, races) · **27**: memory *management* (stack, heap, GC) | — |
| Worker pools, pipelines | 26 | 37 |
| Serialization | 32 | 31 |
| Functional options | 34 | 35, 37 |
| Logging | 37 (`log/slog` in production) | 20 |

### Where things are
```
golang/
├── README.md                       ← you are here: how to use + topic index & progress
├── AGENTS.md                       ← rules the AI tutors follow (the single source of truth)
├── templates/
│   ├── topic-template.md           ← structure every topic's notes follow
│   └── topic/                      ← scaffold module every topic starts from
├── scripts/new-topic.sh            ← creates a topic module from the scaffold
└── topics/
    └── NN-topic-slug/              ← each topic is its own Go module: runs on its own
        ├── go.mod                  ← github.com/Shreyash019/Golang/topics/NN-topic-slug
        ├── Makefile                ← `make` lists every command
        ├── topic-slug.md           ← notes (top) + Q&A log (bottom)
        ├── cmd/examples/           ← runs the worked examples from the notes
        ├── internal/examples/      ← the worked examples, output verified by tests
        └── exercises/
            ├── README.md
            ├── 01-predict-…/       ← E1 predict the output
            ├── 02-fix-…/           ← E2 fix the bug (a real gotcha)
            ├── 03-implement-…/     ← E3 implement to failing table-driven tests
            ├── 04-refactor-…/      ← E4 rewrite Java/Python-style code as idiomatic Go
            ├── 05-project-…/       ← E5 mini-project, laid out like a Go backend:
            │   ├── cmd/<app>/      ←   main.go: wiring + graceful shutdown only
            │   └── internal/       ←   config/, <domain>/, storage/, httpapi/, apperr/ (only what it needs)
            └── solutions/          ← reference solutions: open only after trying
```

There's no `go.mod` at the repo root: every topic builds, runs and tests on its own, using only the standard library (`net/http`, `encoding/json`, `log/slog`, `context`, …) unless a topic is about a specific package.

Folder names come from the topic's **bold short name**: **Arrays & slices** (06) → `topics/06-arrays-and-slices/arrays-and-slices.md`.

### How each topic file is organised
| Section | Read it for |
|---|---|
| 1. TL;DR | The short version. Read it first, and reread it for revision |
| 2. Why It Exists | The problem Go's design solves, and what Python/Java do instead |
| 3. Mental Model | The one picture to keep in mind, and where it breaks |
| 4. Core Concepts | Precise terms and common confusions |
| 5. How It Works Inside | The runtime/compiler mechanism, with diagrams and "predict the output" |
| 6. Worked Examples | A runnable example per concept, with its real output (`make example N=<name>`) |
| 7. Gotchas | The classic traps: reproduction → why → fix |
| 8. Idiomatic vs Non-Idiomatic | Side-by-side code, with sources |
| 9. Coming from Python / Java | Mapping table + habits to unlearn |
| 10. Production Lens | Performance, testing, observability, API design, a senior review |
| 11. Interview Lens | Common questions and what a strong answer shows |
| 12. Exercises | E1–E5 links, commands and status/scores |
| 13. Cheat Sheet | 10 lines to reread |
| 🗣️ Q&A Discussion | Every session: answers, scores, exercise reviews, where I got stuck |

**First read:** sections 1–4, then predict your way through 5–7 (cover the `Output` blocks), then do E1–E2, read 8–11, do E3–E5, then the Q&A.
**Revision:** the TL;DR, the Cheat Sheet, then *Session summary → gaps* and the *Stuck* entries in the Q&A log.

### Studying with an AI tutor
The same commands work in all three tools:

| Command | What happens |
|---|---|
| `diagnose me` | Placement session: reads my code, asks questions, writes `topics/00-diagnostic/diagnostic.md` with why I get stuck + a learning plan |
| `/learn-go <topic>` or `learn <topic>` | Creates the topic folder with full notes + verified exercises, then starts the Q&A |
| `check E3` (or `check <topic> E3`) | Runs `make check` (gofmt, vet, `go test -race`) on my solution, reviews it like a senior Go reviewer, scores it /5 |
| `I'm stuck on E3` / `hint` | Hint ladder: what did you try → where to look → smaller example → shape of the solution. Never the answer unless I ask |
| `show solution E3` | The full solution, explained. The exercise is then marked 🔁 redo until I redo it unaided |
| `quiz me on <topic>` / `continue <topic>` | New Q&A session focused on past gaps and stuck patterns |
| `revise <topic>` | Quick recap + 3 questions on weak spots |
| `review my code <path>` | A senior-style review of any Go code in the repo |
| `mock interview [topic]` | 45–60 min Go technical interview, logged in topic 39 |

| Tool | Reads rules from | Slash command from |
|---|---|---|
| Claude (VS Code) | `CLAUDE.md` → `AGENTS.md` | `.claude/skills/learn-go/SKILL.md` (`/learn-go`) |
| Copilot (VS Code, agent mode) | `.github/copilot-instructions.md` → `AGENTS.md` | `.github/prompts/learn-go.prompt.md` (`/learn-go`) |
| Gemini (Antigravity) | `.agent/rules/go-tutor.md`, `GEMINI.md` → `AGENTS.md` | `.agent/workflows/learn-go.md` (`/learn-go`) |

To change how the tutor behaves, edit **only `AGENTS.md`**. All the other files point to it.

### Running a topic
Each topic is its own Go module (`github.com/Shreyash019/Golang/topics/NN-slug`). `cd` into it and use `make` (run `make` alone to list the targets):

```bash
cd topics/06-arrays-and-slices

make examples                       # list the worked examples from the notes
make example N=slice-aliasing       # run one (N=all runs them all)

make ex E=01-predict-append         # E1: fill in prediction.go FIRST, then run
go run ./exercises/01-predict-append
make ex E=03-implement-stack        # run one exercise's tests (verbose, race detector)
make check E=03-implement-stack     # gofmt + vet + race tests: run before asking the tutor to "check E3"

make test                           # examples + reference solutions must pass
make verify                         # everything the tutor checks before a topic is ready
```

Plain `go` commands work too (`go test ./exercises/03-implement-stack/ -run TestPush -v`).

The loop for every exercise: **read the failing test → make it pass → `make check E=…` → ask the tutor to `check En`**. Starter tests fail on purpose. E4's behaviour tests pass from the start and must keep passing while you refactor.

A new topic is created with `scripts/new-topic.sh NN-slug` (the AI tutor runs it for you).

| Toolchain | Version | Install (macOS) |
|---|---|---|
| Go | 1.22+ (topics are created with the installed `go version`) | `brew install go` |
| make | any | comes with Xcode Command Line Tools (`xcode-select --install`) |
| golangci-lint (optional) | latest | `brew install golangci-lint` |
| VS Code Go extension | latest | installs `gopls`, which picks up each topic module automatically |

---

## 🗺️ Topic Index & Progress

Status: ⬜ Not started · 📖 In progress · ✅ Done · 🔁 Needs revision
Priority (from the syllabus): 🔴 must master · 🟠 important · 🟡 good to know
*(The AI tutor updates these tables after every session.)*

### Phase 0 — Diagnosis
| # | Topic | Pri | Covers | Status | Link | Last session | Notes / gaps |
|---|---|---|---|---|---|---|---|
| 00 | **Diagnostic**: where I get stuck and why (placement + learning plan). Say `diagnose me` | 🔴 | all | ⬜ | | | |

### Phase 1 — Getting Started
| # | Topic | Pri | Covers | Status | Link | Last session | Notes / gaps |
|---|---|---|---|---|---|---|---|
| 01 | **Go fundamentals & toolchain**: Go's philosophy, installation, the `go` command (build, run, test, vet, install, doc), gofmt, project structure, `main`, packages, imports, comments & doc comments | 🔴 | 1, 42 | ⬜ | | | |
| 02 | **Variables & constants**: `var`, `:=`, `const`, `iota`, untyped constants, zero values, type inference, scope, shadowing | 🔴 | 2 | ⬜ | | | |
| 03 | **Basic types & operators**: bool, integers, floats, complex, string, byte, rune, conversions, overflow, a first look at aliases vs defined types; arithmetic, comparison, logical, bitwise, assignment, `++`/`--`, precedence | 🔴 | 3, 4 | ⬜ | | | |
| 04 | **Control flow**: `if` (with init), `for` (all forms), `range` incl. range-over-int, `switch`, `fallthrough`, `break`/`continue`, labels, `goto` | 🔴 | 5 | ⬜ | | | |
| 05 | **Functions, closures & defer**: declaration, parameters, multiple & named returns, variadics, anonymous functions, first-class & higher-order functions, closures, recursion, `defer` | 🔴 | 6, 7 | ⬜ | | | |

### Phase 2 — Data & Memory Basics
| # | Topic | Pri | Covers | Status | Link | Last session | Notes / gaps |
|---|---|---|---|---|---|---|---|
| 06 | **Arrays & slices**: arrays and their value semantics, the slice header, length & capacity, `append` growth, `copy`, slicing & full slice expressions, nil vs empty, shared backing arrays, the `slices` package | 🔴 | 8, 9 | ⬜ | | | |
| 07 | **Maps**: creation, lookup, comma-ok, insert, delete, iteration order, nil maps, internals, the `maps` package | 🔴 | 10 | ⬜ | | | |
| 08 | **Strings, bytes & runes**: representation, UTF-8, indexing vs ranging, conversions, `strings`, `strings.Builder`, `bytes`, `unicode/utf8` | 🔴 | 11 | ⬜ | | | |
| 09 | **Pointers & value semantics**: addresses, `&`, `*`, dereferencing, nil pointers, copying values, when to use pointers, a first look at stack vs heap | 🔴 | 12, 27 (intro) | ⬜ | | | |

### Phase 3 — Types & Abstraction
| # | Topic | Pri | Covers | Status | Link | Last session | Notes / gaps |
|---|---|---|---|---|---|---|---|
| 10 | **Structs & methods**: definition, fields, initialization, nested & anonymous structs, comparison, methods, value vs pointer receivers, method sets, useful zero values | 🔴 | 13, 14 | ⬜ | | | |
| 11 | **Composition & embedding**: struct embedding, promoted fields & methods, delegation, composition vs inheritance | 🔴 | 20 | ⬜ | | | |
| 12 | **Interfaces & polymorphism**: definition, implicit implementation, interface values (internals), `any`, type assertions, type switches, nil interfaces, interface composition, dynamic dispatch, substitutability | 🔴 | 16, 17, 19 | ⬜ | | | |
| 13 | **The type system**: defined types vs aliases, underlying types, type identity, assignability, conversions, untyped constants revisited | 🔴 | 26, 3 (aliases & defined types) | ⬜ | | | |
| 14 | **Generics**: type parameters, constraints, `comparable`, generic functions & types, type inference, when not to use them | 🔴 | 25 | ⬜ | | | |
| 15 | **Packages, modules & encapsulation**: exported vs unexported, package boundaries & naming, `internal/`, `go.mod`, dependencies, versions, `go get`, `go mod tidy`, module structure | 🔴 | 15, 23, 24 | ⬜ | | | |
| 16 | **Designing abstractions**: interfaces as contracts, small consumer-defined interfaces, implementation hiding, dependency boundaries, "accept interfaces, return structs", package API design | 🔴 | 18, 15 (API design) | ⬜ | | | |

### Phase 4 — Errors, Testing & Debugging
| # | Topic | Pri | Covers | Status | Link | Last session | Notes / gaps |
|---|---|---|---|---|---|---|---|
| 17 | **Errors, panic & recover**: the `error` interface, creating errors, wrapping with `%w`, `errors.Is`/`As`/`Join`, sentinel vs custom errors, `panic`, `recover`, when to use and avoid them | 🔴 | 21, 22 | ⬜ | | | |
| 18 | **Testing**: the `testing` package, table-driven tests, subtests, helpers, examples, benchmarks, fuzzing, `httptest`, coverage | 🔴 | 39 | ⬜ | | | |
| 19 | **Mocking & test design**: interfaces for testing, dependency injection, fakes/stubs/mocks, integration tests | 🟠 | 40 | ⬜ | | | |
| 20 | **Debugging**: reading stack traces & panics, Delve, debugging tests, logging for debugging, first profiling steps | 🟠 | 41 | ⬜ | | | |

### Phase 5 — Concurrency
| # | Topic | Pri | Covers | Status | Link | Last session | Notes / gaps |
|---|---|---|---|---|---|---|---|
| 21 | **Goroutines & the scheduler**: concurrency vs parallelism, goroutines, the G-M-P scheduler, preemption, growable stacks, synchronization concepts | 🔴 | 28 | ⬜ | | | |
| 22 | **Channels & select**: creation, send/receive, blocking, buffered vs unbuffered, closing, `range` over channels, nil channels, `select`, timeouts, `default`, internals | 🔴 | 29, 30 | ⬜ | | | |
| 23 | **sync & atomics**: `Mutex`, `RWMutex`, `WaitGroup`, `Once`, `sync/atomic`, choosing a mutex vs a channel | 🔴 | 31 | ⬜ | | | |
| 24 | **Race conditions & the Go memory model**: data races, the race detector, shared state, happens-before, synchronization strategies | 🔴 | 34 | ⬜ | | | |
| 25 | **Context**: `context.Context`, cancellation, deadlines, timeouts, request-scoped values, propagation | 🔴 | 33 | ⬜ | | | |
| 26 | **Concurrency patterns**: worker pools, pipelines, fan-in/fan-out, producer-consumer, cancellation, errgroup, semaphores, leak prevention, testing concurrent code | 🔴 | 32 | ⬜ | | | |

### Phase 6 — Runtime & Performance
| # | Topic | Pri | Covers | Status | Link | Last session | Notes / gaps |
|---|---|---|---|---|---|---|---|
| 27 | **Memory management**: value vs reference semantics revisited, stack vs heap, allocation, escape analysis, the garbage collector (`GOGC`, `GOMEMLIMIT`), aliasing | 🔴 | 27 | ⬜ | | | Not the "Go memory model" (that's 24) |
| 28 | **Profiling & performance**: CPU & memory profiling with pprof, allocations, copying, slices, strings and interfaces costs, benchmarks with benchstat, optimization workflow | 🟠 | 43, 44 | ⬜ | | | |
| 29 | **Reflection & unsafe**: `reflect` types & values, use cases & trade-offs, `unsafe` pointers, memory layout, the risks | 🟡 | 45, 46 | ⬜ | | | |

### Phase 7 — Standard Library & I/O
| # | Topic | Pri | Covers | Status | Link | Last session | Notes / gaps |
|---|---|---|---|---|---|---|---|
| 30 | **Standard library essentials**: `fmt`, `strings`, `strconv`, `bytes`, `time`, how to read package docs | 🔴 | 35 | ⬜ | | | |
| 31 | **Files & I/O**: `io.Reader`/`io.Writer`, `os`, files, `bufio`, streams, `path/filepath`, `io/fs` | 🟠 | 36, 35 (`os`, `io`, `path/filepath`) | ⬜ | | | |
| 32 | **JSON & serialization**: `encoding/json`, marshal/unmarshal, struct tags, custom serialization, streaming | 🔴 | 37, 36 (serialization) | ⬜ | | | |
| 33 | **HTTP**: `net/http` servers, handlers, routing, middleware, requests/responses, clients, timeouts | 🔴 | 38 | ⬜ | | | |

### Phase 8 — Idioms, Patterns & Production
| # | Topic | Pri | Covers | Status | Link | Last session | Notes / gaps |
|---|---|---|---|---|---|---|---|
| 34 | **Functional patterns**: callbacks, higher-order functions & closures in practice, functional options, middleware chains | 🟠 | 47 | ⬜ | | | |
| 35 | **Design patterns in Go**: Factory, Strategy, Adapter, Decorator, Observer, Repository, dependency injection, and when Go doesn't need them | 🟠 | 48 | ⬜ | | | |
| 36 | **Idiomatic Go**: Effective Go principles, simplicity, composition, small interfaces, error handling, unlearning Java & Python habits | 🔴 | 49 | ⬜ | | | |
| 37 | **Production Go**: service structure, configuration, `log/slog`, graceful shutdown, observability, reliability | 🔴 | 50 | ⬜ | | | |
| 38 | **Capstone**: a production service (HTTP API + worker pool + storage + graceful shutdown + observability + tests) | 🔴 | all | ⬜ | | | |
| 39 | **Go interview practice**: mixed mock interviews. Say `mock interview` | 🔴 | all | ⬜ | | | |

---

## 📋 Syllabus coverage

The original 50-item syllabus, and the topic that covers each item.

| # | Section | Topic | # | Section | Topic |
|---|---|---|---|---|---|
| 1 | Go Fundamentals | 01 | 26 | Type System | 13 |
| 2 | Variables & Constants | 02 | 27 | Memory Model (memory management) | 27 (intro in 09) |
| 3 | Data Types | 03 (aliases/defined types in 13) | 28 | Concurrency Fundamentals | 21 |
| 4 | Operators & Expressions | 03 | 29 | Channels | 22 |
| 5 | Control Flow | 04 | 30 | Select | 22 |
| 6 | Functions | 05 | 31 | Synchronization | 23 |
| 7 | Function Concepts | 05 | 32 | Concurrency Patterns | 26 |
| 8 | Arrays | 06 | 33 | Context | 25 |
| 9 | Slices | 06 | 34 | Race Conditions | 24 |
| 10 | Maps | 07 | 35 | Standard Library | 30 (`os`/`io`/`filepath` in 31) |
| 11 | Strings | 08 | 36 | Files & I/O | 31 (serialization in 32) |
| 12 | Pointers | 09 | 37 | JSON | 32 |
| 13 | Structs | 10 | 38 | HTTP | 33 |
| 14 | Methods | 10 | 39 | Testing | 18 |
| 15 | Encapsulation | 15 (API design in 16) | 40 | Mocking & Test Design | 19 |
| 16 | Interfaces | 12 | 41 | Debugging | 20 |
| 17 | Interface Deep Dive | 12 | 42 | Tooling | 01 |
| 18 | Abstraction | 16 | 43 | Profiling | 28 |
| 19 | Polymorphism | 12 | 44 | Performance | 28 |
| 20 | Composition | 11 | 45 | Reflection | 29 |
| 21 | Error Handling | 17 | 46 | Unsafe | 29 |
| 22 | Panic & Recovery | 17 | 47 | Functional Patterns | 34 |
| 23 | Packages | 15 | 48 | Design Patterns in Go | 35 |
| 24 | Go Modules | 15 | 49 | Idiomatic Go | 36 |
| 25 | Generics | 14 | 50 | Production Go | 37 |
