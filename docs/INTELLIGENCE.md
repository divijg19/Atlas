# Atlas Intelligence Model

This document defines the **canonical intelligence model** for Atlas. It is the
conceptual companion to [`ARCHITECTURE.md`](./ARCHITECTURE.md), which defines
layer ownership. Where `ARCHITECTURE.md` answers *"which package owns this?"*,
this document answers *"what is the shape of intelligence as data flows through
Atlas?"*

The model is a strictly downward pipeline. Each stage transforms the previous
stage's output into a richer, still-deterministic representation. Nothing is
ever inferred probabilistically: given the same inputs, Atlas always produces
the same intelligence.

---

## The Eight Stages

```
Vestiges      raw normalized observations        (signals.Repo, profile.UserMetadata, contributions.Summary)
   ↓
Facts         deterministic aggregates            (signals.Facts)
   ↓
Indicators    measured component signals          (signals.Signals → signals.RawScore)
   ↓
Profile       canonical candidate aggregate       (index.Profile)
   ↓
Evaluation    score calibration & confidence      (internal/evaluation)
   ↓
Intelligence  the evaluated candidate             (Profile + Evaluation output)
   ↓
Projections   consumer-shaped read models         (internal/projection)
   ↓
Consumers     CLI, HTTP API, web UI               (cmd/atlas, cmd/server, web)
```

### 1. Vestiges

The raw, normalized observations pulled from GitHub and mapped into the
domain. These are facts about the world that Atlas does not compute — it only
records them.

- `signals.Repo` — a normalized repository observation (name, fork, size,
  timestamps, visibility, archived, template, license, topics, stars, forks,
  watchers, open issues, created/pushed dates, default branch).
- `profile.UserMetadata` — normalized account metadata.
- `contributions.Summary` — normalized contribution totals.

Vestiges are owned by **Acquisition** and **Normalization**
(`internal/acquisition`).

### 2. Facts

Deterministic, evidence-backed aggregates computed from vestiges. Facts answer
*"what do we know about this user's repository portfolio?"* They are the
inspectable layer between raw observations and derived indicators.

Defined by `signals.Facts` and produced by `signals.FromRepos`. Examples:
`TotalRepos`, `OriginalRepos`, `ForkRepos`, `RecentRepos`, `DeepRepos`,
`ValidRepos`, `LargestRepoSize`, `LatestActivity`, plus the Phase 9 metadata
facts (`ArchivedRepos`, `PublicRepos`, `PrivateRepos`, `LicensedRepos`,
`TotalStars`, `TotalForks`, `TotalWatchers`, `TotalOpenIssues`, `TotalTopics`,
`OldestCreated`, `NewestCreated`).

Facts introduce **no new indicators**. They only enrich observations with
deterministic counts and sums.

### 3. Indicators

Measured component signals derived from facts. Indicators quantify observations:
ownership, consistency, depth, activity. Produced by `signals.ExtractSignals`
and `signals.ExtractSignalsFromFacts`, yielding `signals.Signals` and a
`signals.RawScore` of the three component scores (ownership, consistency,
depth).

Indicators are measurements, not judgments. The overall score is **not** an
indicator — it is an evaluation.

### 4. Profile

The canonical candidate aggregate. `index.Profile` stores vestiges-level
observations (facts, signals, metadata, contributions) and is the single source
of truth for what Atlas knows about a candidate. It stores observations, not
evaluations.

### 5. Evaluation

How scores are interpreted. `internal/evaluation` owns:

- Overall score assembly — `OverallScore(RawScore)` combines the three
  component scores with the canonical ranking weights (ownership `0.3`, consistency
  `0.4`, depth `0.3`).
- Small-sample penalty — `ApplySmallSamplePenalty(rawScore, repoCount)`
  down-weights evaluations built on too few repositories (threshold `3`,
  multiplier `7`).
- Confidence classification — `ClassifyConfidence` maps a repo count to
  `high` / `moderate` / `low`.
- Ranking policy — `RankingPolicy{ Score }` is the single type consumed by the
  engine and projection layers for ordering.

Evaluation is the single source of truth for score interpretation.

### 6. Intelligence

The assembled result of applying Evaluation to a Profile. Intelligence is the
evaluated candidate: a Profile enriched with an overall score, confidence, and
ranking. It is what the projection layer renders into consumer shapes.

### 7. Projections

Consumer-shaped, read-only, deterministic views (`internal/projection`):
`AnalyzeProjection`, `InspectProjection`, `SearchProjection`. Projections
**do not** compute overall scores or penalties — those are supplied by
Evaluation. Projections only re-shape and order.

### 8. Consumers

Presentation surfaces that render intelligence: `cmd/atlas` (CLI),
`cmd/server` (HTTP API), and `web` (UI). Consumers format; they never derive
facts, signals, scores, or rankings.

---

## Design Principles

Every stage of the intelligence model is governed by the same principles:

- **Deterministic** — identical inputs always yield identical output. No
  randomness, no clock-dependent results in scored paths (the activity indicator
  is deterministic *given a reference time*; see caveat below).
- **Explainable** — every score can be traced to the facts and weights that
  produced it.
- **Composable** — stages stack cleanly; each consumes only the output of the
  stage above it.
- **Observable** — intermediate stages are inspectable (`InspectProjection`
  exposes raw vestiges, facts, and indicators).
- **Inspectable** — there is no "black box"; the full pipeline can be dumped.
- **Evidence-backed** — conclusions reference the observations that support
  them.
- **Layer-owned** — each transformation has exactly one owning package.
- **Transport-independent** — intelligence does not depend on HTTP, the GitHub
  API, or any transport.
- **Presentation-independent** — intelligence does not depend on CLI, JSON, or
  web rendering.

---

## Non-Determinism Caveat

The `activity` indicator uses `time.Now()` during signal extraction to bucket
repositories into recent windows. This makes activity deterministic **given a
reference time**, but not purely a function of stored vestiges. This is a known
limitation and is explicitly **out of scope for v0.8.13**. All other stages are
fully deterministic from stored observations.

---

## Historical Notes

- **v0.8.13 (Phase 1)** — rebranded GH-Analyzer → Atlas; module path
  `github.com/divijg19/Atlas`; removed the legacy `cmd/gha` and
  `internal/ghanalyzer` packages.
- **v0.8.13 (Phase 6)** — consolidated Evaluation: ranking weights, overall
  score, and small-sample penalty moved into `internal/evaluation`
  (`scoring.go`). `internal/engine/ranking.go` now defines only the
  `RankingStrategy` interface; projection and engine consume
  `evaluation.RankingPolicy`.
- **v0.8.13 (Phase 9)** — Repository Intelligence Foundation. Enriched
  `RepoDTO`, `signals.Repo`, `NormalizeRepos`, and `signals.Facts` with
  repository metadata (visibility, archived, template, license, topics, stars,
  forks, watchers, open issues, created/pushed dates, default branch). **No new
  indicators were introduced** — only observations and deterministic facts.
