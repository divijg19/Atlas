# GH-Analyzer Architecture

This document is the canonical reference for layer ownership in GH-Analyzer.
It is the source of truth for where new code belongs. When proposing a change,
first answer: **Which layer owns this?** If the answer is not obvious, do not
implement it yet — revisit the architecture.

The pipeline is strictly downward:

```
GitHub
   ↓
Transport          internal/github
   ↓
Acquisition        internal/acquisition
   ↓
Normalization      internal/acquisition/normalize.go
   ↓
Domain             signals · profile · contributions · index · projection
   ↓
Projection         internal/projection
   ↓
Presentation       cmd/gh-analyzer · cmd/server · web
```

---

## Layer Ownership

### Transport — `internal/github`

**Owns**

- HTTP client construction and lifecycle
- Authentication (`GITHUB_TOKEN` → `Authorization: Bearer`)
- Request headers (User-Agent, auth)
- Request execution (`Do`)

**Never owns**

- Facts, Signals, Profiles, Contributions
- Presentation
- GitHub endpoint semantics or DTOs

Transport is the only package permitted to perform raw `http` I/O on behalf of
authentication and execution. It has no knowledge of the domain.

### Acquisition — `internal/acquisition`

**Owns**

- GitHub REST endpoints and their URLs
- GitHub DTOs (mirroring GitHub's JSON schema exactly)
- GitHub endpoint semantics, response handling, pagination
- GitHub API errors (`APIError`)
- The acquisition client (`Client`) used by consumers

**Never owns**

- Signals, Facts, Profiles, business logic
- Ranking, scoring, evidence
- Presentation

Acquisition owns **GitHub's schema, not ours**. If GitHub changes its API,
only this layer changes; everything above remains untouched. DTOs mirror
GitHub field names and types verbatim — timestamps remain raw strings inside
DTOs and are parsed only during normalization.

### Normalization — `internal/acquisition/normalize.go`

**Owns**

- Mapping GitHub DTOs to domain models:
  - `RepoDTO` → `signals.Repo`
  - `UserDTO` → `profile.UserMetadata`
  - `ContributionsDTO` → `contributions.Summary`
- Timestamp parsing (`created_at` / `updated_at` → `time.Time`, zero-fallback)

**Never performs**

- Network requests
- Ranking, signal computation, or any business logic

Normalization is co-located in `internal/acquisition` (not a separate
package). It is the single boundary between GitHub's representation and the
domain's.

### Domain

Packages: `signals`, `profile`, `contributions`, `index`, `projection`.

**Owns**

- Facts (`signals.Facts`)
- Signals (`signals.Signals`, scoring, evidence)
- Metadata (`profile.UserMetadata`)
- Contributions (`contributions.Summary`)
- The candidate aggregate (`index.Profile`)
- Consumer-facing read models (`projection.CandidateProjection`)

**Rules**

- Domain packages are **pure**: they import no `net/http` and no
  `internal/github`. They contain only models, facts, signals, evidence, and
  pure derivation.
- Acquisition is always the consumer of the domain, never the reverse.

### Projection — `internal/projection`

**Owns**

- Consumer-facing read models (`CandidateProjection`)

**Never performs**

- Networking, normalization, ranking, or storage

Projection is a read-only, deterministic view of the domain for presentation
layers. It owns no business logic and no persistence.

### Presentation — `cmd/gh-analyzer`, `cmd/server`, `web`

**Owns**

- CLI command handling and formatting
- HTTP API surface and JSON shaping
- Web UI

**Never computes**

- Intelligence. Presentation formats; it does not derive facts, signals, or
  rankings. All intelligence is produced by the domain and reached through
  projection.

---

## Decision Rule

Before implementing any change, ask:

> **Which layer owns this?**

- Network call to GitHub? → **Acquisition**
- GitHub JSON → internal model? → **Normalization**
- A fact, signal, or aggregate? → **Domain**
- A view for a consumer? → **Projection**
- Output text or HTTP response? → **Presentation**

If a change spans multiple layers, it is likely too large — split it. If no
single layer clearly owns it, the architecture must be revisited before code
is written.

---

## Historical Note

Prior to v0.8.11, acquisition (HTTP fetch, decode, partial normalization) was
spread across `signals`, `profile`, `contributions`, and `live`. v0.8.11
consolidated all GitHub REST access into `internal/acquisition`, made the
domain packages pure, and established the normalization boundary. The
`Profile` remains the canonical candidate aggregate; `Report` (analyze path) is
a presentation-oriented derivation and is not the persisted source of truth.
