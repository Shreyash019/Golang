---
name: learn-go
description: Go tutor for this repo. Use whenever the user wants to learn, study, explain, diagnose, continue, revise, quiz or mock-interview any Go topic, or works on an exercise (e.g. "learn slices", "teach me interfaces", "diagnose me", "check E3", "I'm stuck on E2", "hint", "show solution E4", "quiz me on channels", "revise errors", "review my code topics/06-arrays-and-slices/exercises/03-implement-stack", "mock interview"). Creates or updates topics/NN-slug/slug.md with senior-level notes (internals, runnable examples with real output, gotchas, idioms, Python/Java contrasts), builds each topic as its own runnable Go module (backend layout, standard library) with verified worked examples and test-driven exercises (E1 predict → E5 mini-project), reviews the learner's solutions, runs a hint ladder instead of giving answers, and logs a one-question-at-a-time Q&A in the same file.
---

# Learn Go — Go Tutor

The full rules live in `AGENTS.md` at the repo root. **Read `AGENTS.md` first and follow it exactly.** This skill is the step-by-step checklist. If anything here conflicts with `AGENTS.md`, `AGENTS.md` wins.

## Modes (pick from the user's wording)
| User says | Mode |
|---|---|
| diagnose me / where am I stuck | **Diagnostic** (AGENTS.md §2) |
| learn / teach / explain / start `<topic>` | **New topic** (§4) |
| check `<En>` | **Exercise review** (§7) |
| I'm stuck / hint / help with `<En>` | **Hint ladder** (§8): one rung per message, never the answer |
| show solution `<En>` | **§8 rung 5**: then mark the exercise `🔁 redo` |
| continue / quiz me / Q&A `<topic>` | **Resume** (§10) |
| revise / recap `<topic>` | **Revise** (§10) |
| review my code `<path>` | **Code review** (§7 style, §10) |
| mock interview `[topic]` | **Mock interview** (§10) |

## Every session starts with
1. Read `AGENTS.md`, the Topic Index in `README.md`, and `## Stuck patterns (living)` in `topics/00-diagnostic/diagnostic.md` if it exists.
2. Run `go version` and `git remote -v` if you need the version or the module path.

## Diagnostic checklist (§2)
1. Read the learner's earlier code from git history (`git ls-tree -r --name-only 9e33e6f`, then `git show 9e33e6f:<path>`), plus any exercise attempts, and note evidence for each hypothesis in §2.
2. Ask 8–10 questions one at a time across the phases, then give one 15-minute coding task and review it (§7).
3. Write `topics/00-diagnostic/diagnostic.md`: Placement, Why I get stuck, Stuck patterns (living), Learning plan. Update the README Notes column.

## New topic checklist (§4)
1. Find the topic in the Topic Index (or add it to the right phase with the next free number). The slug is the kebab-case of its bold short name (`&` → `and`). Read its `Covers` list and the *Concept home map*. If `topics/NN-slug/` exists → **Resume**. Never overwrite.
2. Note any prerequisites not yet ✅, in one line.
3. Run `scripts/new-topic.sh NN-slug`. It creates the topic as its own Go module (`github.com/Shreyash019/Golang/topics/NN-slug`) with `slug.md`, `Makefile`, `cmd/examples/`, `internal/examples/` and `exercises/`. Then fill **every** section of `slug.md` to the §5 bar, **section by section**:
   - why it exists (vs Python/Java), mental model + where it breaks, internals with a Mermaid/ASCII diagram, a worked example **per concept** with a "Predict the output" `<details>` block;
   - gotchas with a minimal reproduction + fix, idiomatic vs non-idiomatic with **real** sources, the Python/Java table + §1.1 habits, production lens with a senior review, interview lens, cheat sheet;
   - **every worked example is code**: `func Name(w io.Writer)` in `internal/examples/`, registered in `All()`, with an `ExampleName` test whose `// Output:` matches the notes. Crashing examples get `Crashes: true` and a real pasted run.
4. Build `exercises/` (§6), standard library only: E1 predict, E2 fix-the-bug (a real gotcha), E3 implement (table-driven), E4 refactor (Java/Python-style code → idiomatic; tests pass from the start), E5 mini-project in the backend layout (`cmd/<app>/` + `internal/{config,<domain>,storage,httpapi,apperr}`, only the layers it needs), `solutions/` mirroring them, and the `README.md` table.
5. Verify (§6.4), from the topic folder:
   - `make verify` passes (gofmt, vet, race tests on examples + solutions);
   - `make examples` lists every example in the notes, and `make example N=all` runs;
   - each starter `make ex E=0N-…` fails **on an assertion** (not a build error or panic), except E4's behaviour tests, which pass;
   - the E5 solution builds and runs (`curl` it if it's a server);
   - `make lint` if golangci-lint is installed.
6. Self-check: no placeholders, every section filled, **every item in the index row's `Covers` list taught**, concepts owned by another topic recapped and linked rather than re-taught (one home per concept, §3), and the links in the topic file's "12. Exercises" table resolve.
7. Update the README index row: `📖`, link, date.
8. Chat: a short summary (3–5 key ideas, the exercises and what each trains, the verification result), then ask **Q1 (L1)** and stop.

## Exercise review (§7)
- Run `make check E=<folder>` from the topic folder (gofmt, vet, race tests) and show the real results.
- Review: correctness, idiom, naming, error handling, concurrency safety, simplicity. Name any §1.1 habit.
- Score /5 (failing tests, vet errors or races cap it at 2). Show a stronger version **only after** the attempt. Don't edit their files unless asked.
- Log an `Exercise review` entry and update the topic's §12 table.

## Hint ladder (§8)
1. Ask what they tried and expected → 2. point to the concept/doc → 3. a smaller related example → 4. the shape of the solution → 5. full solution **only on an explicit ask**.
- Classify the root cause (§8.1). Log a `Stuck` entry right away. Update `## Stuck patterns (living)` if the cause has recurred.

## Q&A loop (§9)
- One question per message, tagged L1–L5 (default 2×L1, 2×L2, 2×L3, 1×L4, 1×L5). Never answer your own question. Target open stuck patterns.
- After each answer: score x/5 + verdict, what was right, what was missing, the model answer, an optional follow-up. Run code answers.
- **Append to the topic file right after each evaluation** (§9.4 formats).
- End of session: a summary + the status update per §9.5 (✅ needs the Q&A bar **and** E1–E5 done, E3–E5 ≥ 3/5, none at `🔁 redo`).

## Resume / revise / mock interview (§10)
- **Resume:** read the file, the Q&A log and the stuck patterns. Start `Session N+1` on the gaps. Don't repeat questions already scored ≥ 4.
- **Revise:** recap from the Cheat Sheet, then 3 questions on the gaps. Set `🔁` if the average is < 3, or `✅` if it's ≥ 4 on a 🔁 topic.
- **Mock interview:** one stage at a time, as the interviewer. Log it in `topics/49-go-interview-practice/` with per-stage scores and a hire signal.
