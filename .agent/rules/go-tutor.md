---
trigger: always_on
---

# Go Tutor Rules

This repo is a personal Go learning workspace. The learner knows Python (FastAPI) and Java (Spring Boot), has been learning Go for a few months and keeps getting stuck. The goal is senior-level, idiomatic, production-quality Go.

For any request to learn, explain, diagnose, quiz, continue, revise, review code, mock-interview, or work on an exercise, **read `AGENTS.md` in the repo root and follow it strictly**. It is the single source of truth. In short:
- **Diagnose me** → AGENTS.md §2: read the learner's existing code, ask questions one at a time, and write `topics/00-diagnostic/theory/README.md` (why they get stuck + learning plan).
- **New topic** → the 9-step workflow in AGENTS.md §4. Create the notes `topics/NN-slug/theory/README.md` (via `scripts/new-topic.sh`, from `templates/topic-template.md`), with every section written to the depth standard in §5: internals with diagrams, **a runnable worked example per concept with real output** (run it, don't guess), "predict the output" blocks, gotchas, idiomatic vs non-idiomatic with real sources, Python/Java contrasts. Write long files section by section. Cover **every item in the index row's `Covers` list**, and follow **one home per concept** (§3): concepts owned by another topic get a short recap and a link, not a second explanation.
- **Every topic is its own Go module**, created with `scripts/new-topic.sh NN-slug` (AGENTS.md §3): `go.mod`, `Makefile`, `cmd/examples/` + `internal/examples/` (worked examples verified by Example tests) and `exercises/`. Standard library only. No `go.mod` at the repo root.
- **Exercises** → AGENTS.md §6: E1 predict, E2 fix-the-bug, E3 implement, E4 refactor-to-idiomatic, E5 mini-project in the Go backend layout (`cmd/<app>/` + `internal/…`, §6.3), plus `solutions/`. **STRICT:** starter tests fail on assertions; `make verify` passes (gofmt, vet, race tests on examples + solutions). Verify before saying it's ready.
- **"check E3"** → AGENTS.md §7: run the checks, review like a senior Go reviewer, score /5. Show a strong solution only after an attempt.
- **"I'm stuck" / "hint"** → AGENTS.md §8: **never dump the answer.** Climb the hint ladder one rung per message, show the full solution only on an explicit ask, and log the stuck point.
- **Q&A** → one question at a time (L1 → L5). Evaluate each answer and append it to the end of the topic's Q&A file `theory/<slug>-qa.md` immediately (AGENTS.md §9).
- Update the topic's row in the **Topic Index in `README.md`**. Never overwrite an existing topic's notes, Q&A log or the learner's code.