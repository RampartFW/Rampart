# Rampart — Task Breakdown

> **Version:** 1.0.0-draft  
> **Date:** 2026-04-11  
> **Scope:** Full implementation (Phase 1–7)  
> **Companion to:** SPECIFICATION.md + IMPLEMENTATION.md  
> **Execution:** Claude Code sequential execution

---

## Task Conventions

```
T-XXX  → Task number (sequential)
P0     → Must have (blocks everything)
P1     → Should have (needed for MVP)
P2     → Nice to have (post-MVP)
```

**Estimates:** S = 1-2h, M = 2-4h, L = 4-8h, XL = 8-16h

---

# ═══════════════════════════════════════════════
# MILESTONE 1 — Project Scaffold (2 tasks)
# ═══════════════════════════════════════════════

### T-001: Initialize Go Module & Project Structure
- **Priority:** P0
- **Estimate:** M
- **Dependencies:** None
- **Description:** `go mod init github.com/rampartfw/rampart`. Create full directory structure per SPECIFICATION §25. Create `cmd/rampart/main.go` with subcommand router (serve, agent, apply, plan, simulate, rollback, snapshot, rules, audit, cluster, cert, validate, fmt, diff, import, export, version). Empty handler stubs.
- **Acceptance Criteria:**
  - [ ] `go mod init github.com/rampartfw/rampart`
  - [ ] `go.mod` has only `golang.org/x/crypto`, `golang.org/x/sys`, `gopkg.in/yaml.v3`
  - [ ] All directories from §25 created with placeholder `.go` files
  - [ ] `main.go` routes all subcommands to stub functions
  - [ ] `go build ./cmd/rampart` compiles successfully
  - [ ] `./rampart version` prints version info
  - [ ] Makefile with targets: build, test, clean, ui-build, lint

### T-002: Core Type Definitions
- **Priority:** P0
- **Estimate:** L
- **Dependencies:** T-001
- **Description:** Define all core types in `internal/model/`: Rule, CompiledRule, PolicySet, Match, PortRange, Schedule, RecurringSpec, Direction, Action, Protocol, ConnState, IPVersion, Conflict, Snapshot, AuditEvent, ExecutionPlan, RuleModification, BackendCapabilities, RuleStats, SimulatedPacket, SimulationResult. Include JSON/YAML struct tags. Include String() methods for enums.
- **Acceptance Criteria:**
  - [ ] `internal/model/rule.go` — Rule, CompiledRule, CompiledMatch, PortRange, RateLimit
  - [ ] `internal/model/policy.go` — PolicySetYAML, PolicyYAML, RuleYAML, MatchYAML, PolicyDefaults, PolicyMetadata, IncludeRef, ScheduleYAML, VariablesYAML
  - [ ] `internal/model/enums.go` — Direction, Action, Protocol, ConnState, IPVersion with String()/MarshalJSON()
  - [ ] `internal/model/plan.go` — ExecutionPlan, RuleModification
  - [ ] `internal/model/conflict.go` — Conflict, ConflictType, Severity
  - [ ] `internal/model/snapshot.go` — Snapshot
  - [ ] `internal/model/audit.go` — AuditEvent, AuditActor, AuditAction, AuditResource, AuditResult
  - [ ] `internal/model/backend.go` — BackendCapabilities, RuleStats, SimulatedPacket, SimulationResult
  - [ ] All types have JSON tags
  - [ ] All YAML-facing types have yaml tags
  - [ ] `go vet ./...` passes

---

# ═══════════════════════════════════════════════
# MILESTONE 2 — YAML Parser & Validator (4 tasks)
# ═══════════════════════════════════════════════

### T-003: YAML Policy Parser
- **Priority:** P0
- **Estimate:** L
- **Dependencies:** T-002
- **Description:** `internal/engine/parser.go`. Parse PolicySet YAML files. Handle polymorphic fields (ports as int/[]int/string, protocol as string/[]string). Validate apiVersion and kind.
- **Acceptance Criteria:**
  - [ ] `ParsePolicyFile(path string) (*PolicySetYAML, error)` function
  - [ ] Handles port specifications: single int, array of int, "start-end" range string
  - [ ] Handles protocol: single string ("tcp") and array (["tcp", "udp"])
  - [ ] Handles CIDR: IPv4 and IPv6, bare IPs auto-suffixed with /32 or /128
  - [ ] Returns clear error messages with file path context
  - [ ] Unit tests with valid and invalid fixtures

### T-004: Schema Validation
- **Priority:** P0
- **Estimate:** M
- **Dependencies:** T-003
- **Description:** `internal/engine/validator.go`. Programmatic validation of parsed YAML. Check all constraints from SPECIFICATION §5.2.
- **Acceptance Criteria:**
  - [ ] apiVersion must be "rampart.dev/v1"
  - [ ] kind must be "PolicySet"
  - [ ] metadata.name required, non-empty
  - [ ] Policy names unique within PolicySet
  - [ ] Rule names unique within Policy
  - [ ] Priority 0-999
  - [ ] Valid CIDR notation (IPv4 + IPv6)
  - [ ] Port numbers 1-65535, start ≤ end
  - [ ] Schedule: valid RFC 3339 dates, activeFrom < activeUntil
  - [ ] Rate limit: rate > 0, burst ≥ rate
  - [ ] Protocol values: tcp, udp, icmp, icmpv6, any
  - [ ] Direction values: inbound, outbound, forward
  - [ ] Action values: accept, drop, reject, log, rate-limit
  - [ ] Returns all validation errors (not just first)
  - [ ] Tests for each validation rule

### T-005: Variable Substitution
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-003
- **Description:** `internal/engine/variable.go`. Parse `rampart-vars.yaml` files. Substitute `${var_name}` in policy YAML before parsing. Support string, int, and array variables.
- **Acceptance Criteria:**
  - [ ] `ParseVariablesFile(path string) (map[string]interface{}, error)`
  - [ ] `SubstituteVars(data []byte, vars map[string]interface{}) ([]byte, error)`
  - [ ] Regex pattern: `\$\{[a-zA-Z_][a-zA-Z0-9_]*\}`
  - [ ] Unresolved variables produce error
  - [ ] Nested variable references not supported (simple single-pass)
  - [ ] Tests with various variable types

### T-006: Policy Includes
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-003
- **Description:** `internal/engine/includes.go`. Resolve `includes` in policy files. Support relative paths, absolute paths. Max depth 10 (circular reference protection).
- **Acceptance Criteria:**
  - [ ] `ResolveIncludes(ps *PolicySetYAML, basePath string) error`
  - [ ] Relative path resolution (relative to including file)
  - [ ] Absolute path support
  - [ ] Circular include detection (max depth 10)
  - [ ] Included policies merged into parent PolicySet
  - [ ] Clear error on missing include files
  - [ ] Tests with nested includes

---

# ═══════════════════════════════════════════════
# MILESTONE 3 — Rule Compiler (3 tasks)
# ═══════════════════════════════════════════════

### T-007: Rule Compiler Core
- **Priority:** P0
- **Estimate:** L
- **Dependencies:** T-004
- **Description:** `internal/engine/compiler.go`. Transform PolicySetYAML → CompiledRuleSet. Normalize CIDRs to net.IPNet, parse ports to PortRange, resolve protocols. Apply policy defaults. Sort by priority (stable). Generate deterministic rule IDs. Compute ruleset hash (SHA-256).
- **Acceptance Criteria:**
  - [ ] `Compile(ps *PolicySetYAML, vars map[string]interface{}) (*CompiledRuleSet, error)`
  - [ ] CIDR normalization: bare IPs → /32 or /128, validate all CIDRs
  - [ ] Port parsing: int → PortRange, "start-end" → PortRange
  - [ ] Protocol normalization: lowercase, validate
  - [ ] Default application: direction, action, states from PolicyDefaults
  - [ ] Stable sort by priority
  - [ ] Deterministic rule ID generation (content-based hash)
  - [ ] Ruleset hash: SHA-256 of canonical serialization
  - [ ] Source file + line tracking per rule
  - [ ] Comprehensive tests

### T-008: UUID v7 Generator
- **Priority:** P0
- **Estimate:** S
- **Dependencies:** T-001
- **Description:** `internal/model/uuid.go`. Implement UUID v7 (time-sortable, random suffix). No google/uuid dependency.
- **Acceptance Criteria:**
  - [ ] `GenerateUUIDv7() string`
  - [ ] 48-bit timestamp (milliseconds)
  - [ ] Version 7 marker
  - [ ] Variant 2 marker
  - [ ] Random suffix from crypto/rand
  - [ ] Standard UUID string format (8-4-4-4-12)
  - [ ] Monotonically increasing for same millisecond
  - [ ] Tests verify format and ordering

### T-009: Execution Plan Generator
- **Priority:** P0
- **Estimate:** M
- **Dependencies:** T-007
- **Description:** `internal/engine/planner.go`. Compare current CompiledRuleSet with desired CompiledRuleSet. Produce ExecutionPlan (add/remove/modify).
- **Acceptance Criteria:**
  - [ ] `GeneratePlan(current, desired *CompiledRuleSet) *ExecutionPlan`
  - [ ] Detect added rules (in desired, not in current)
  - [ ] Detect removed rules (in current, not in desired)
  - [ ] Detect modified rules (in both, different content)
  - [ ] Track which fields changed in modifications
  - [ ] Count statistics (add/remove/modify counts)
  - [ ] Tests with various scenarios

---

# ═══════════════════════════════════════════════
# MILESTONE 4 — Conflict Detection (2 tasks)
# ═══════════════════════════════════════════════

### T-010: Conflict Detection Engine
- **Priority:** P0
- **Estimate:** XL
- **Dependencies:** T-007
- **Description:** `internal/engine/conflict.go`. Pairwise conflict detection for compiled rules. Detect shadow, contradiction, redundancy, subset, and overlap conflicts. CIDR overlap via interval comparison.
- **Acceptance Criteria:**
  - [ ] `DetectConflicts(rules []CompiledRule) []Conflict`
  - [ ] Shadow detection: higher-priority rule makes lower unreachable
  - [ ] Contradiction detection: same priority, overlap, different action
  - [ ] Redundancy detection: identical match + action
  - [ ] Subset detection: strict subset, same action
  - [ ] Overlap detection: partial overlap, different action
  - [ ] CIDR overlap: `cidrsOverlap(a, b []net.IPNet) bool`
  - [ ] Port overlap: `portsOverlap(a, b []PortRange) bool`
  - [ ] Protocol overlap check
  - [ ] Direction must match for conflict
  - [ ] Severity levels: Error (contradiction), Warning (shadow/overlap), Info (redundancy/subset)
  - [ ] Tests for each conflict type
  - [ ] Performance: < 100ms for 1,000 rules

### T-011: Conflict Report Formatter
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-010
- **Description:** Format conflict detection results for CLI (colorized text) and API (JSON). Include human-readable explanations and suggestions.
- **Acceptance Criteria:**
  - [ ] `FormatConflicts(conflicts []Conflict, format string) string`
  - [ ] Text format: colorized, human-readable, with rule references
  - [ ] JSON format: structured for API responses
  - [ ] Warning/error icons (⚠/✗)
  - [ ] Include rule names, priorities, match conditions in report
  - [ ] Suggest resolution for each conflict type

---

# ═══════════════════════════════════════════════
# MILESTONE 5 — nftables Backend (4 tasks)
# ═══════════════════════════════════════════════

### T-012: Backend Interface & Registry
- **Priority:** P0
- **Estimate:** M
- **Dependencies:** T-002
- **Description:** `internal/backend/backend.go` and `internal/backend/registry.go`. Define Backend interface, BackendFactory type, Register/NewBackend/AutoDetect functions.
- **Acceptance Criteria:**
  - [ ] Backend interface with all methods from SPECIFICATION §6.1
  - [ ] BackendFactory type
  - [ ] `Register(name string, factory BackendFactory)`
  - [ ] `NewBackend(name string, cfg BackendConfig) (Backend, error)`
  - [ ] `AutoDetect() (Backend, error)` — probe nftables > iptables > ebpf
  - [ ] Thread-safe registry

### T-013: nftables Backend — Probe & CurrentState
- **Priority:** P0
- **Estimate:** L
- **Dependencies:** T-012
- **Description:** `internal/backend/nftables/nftables.go`. Implement Probe() (check nft binary exists), CurrentState() (parse `nft -j list table inet rampart`). Extract Rampart-managed rules by comment prefix.
- **Acceptance Criteria:**
  - [ ] `Probe()` checks nft binary exists and executable
  - [ ] `Probe()` checks kernel support (run `nft list tables`)
  - [ ] `CurrentState()` executes `nft -j list table inet rampart`
  - [ ] Parse JSON output to extract rules
  - [ ] Filter by comment prefix "rampart:"
  - [ ] Convert nft JSON rules back to CompiledRule
  - [ ] Handle "table not found" gracefully (return empty set)
  - [ ] Name(), Capabilities(), Close() implemented

### T-014: nftables Backend — Rule Compiler
- **Priority:** P0
- **Estimate:** L
- **Dependencies:** T-013
- **Description:** `internal/backend/nftables/compiler.go`. Translate CompiledRule → nftables script commands. Handle IPv4/IPv6, ports, protocols, CIDRs, states, rate limiting, logging. Generate complete table/chain structure.
- **Acceptance Criteria:**
  - [ ] `generateScript(rs *CompiledRuleSet) string` — full nft script
  - [ ] Table creation: `table inet rampart`
  - [ ] Chain creation: input (filter, hook input, priority 0, policy drop), forward, output
  - [ ] Stateful tracking: `ct state established,related accept` at top
  - [ ] Loopback accept: `iifname "lo" accept`
  - [ ] ICMP allow: `meta l4proto icmp accept` + `meta l4proto icmpv6 accept`
  - [ ] Rule → nft command translation (all match types)
  - [ ] Anonymous sets for multiple CIDRs: `{ 10.0.1.0/24, 10.0.2.0/24 }`
  - [ ] Rate limiting: `limit rate X/second burst Y`
  - [ ] Logging: `log prefix "rampart:rulename: "`
  - [ ] Counter per rule
  - [ ] Comment with rule name
  - [ ] Tests with expected nft output

### T-015: nftables Backend — Apply, DryRun, Rollback, Stats
- **Priority:** P0
- **Estimate:** L
- **Dependencies:** T-014
- **Description:** Implement Apply (atomic via `nft -f`), DryRun (compile without applying), Rollback (apply snapshot state), Flush (delete rampart table), Stats (parse counters from JSON output).
- **Acceptance Criteria:**
  - [ ] `Apply()` writes script to temp file → `nft -f` → remove temp
  - [ ] `Apply()` returns error on failure (does not leave partial state)
  - [ ] `DryRun()` compiles and returns ExecutionPlan without applying
  - [ ] `Rollback()` takes Snapshot, reconstructs script, applies
  - [ ] `Flush()` deletes `table inet rampart` entirely
  - [ ] `Stats()` parses counter values from `nft -j list table`
  - [ ] Stats returns map[ruleID] → {packets, bytes}

---

# ═══════════════════════════════════════════════
# MILESTONE 6 — Snapshot & Audit (3 tasks)
# ═══════════════════════════════════════════════

### T-016: Snapshot Store
- **Priority:** P0
- **Estimate:** L
- **Dependencies:** T-012
- **Description:** `internal/snapshot/`. Create, list, load, delete, diff snapshots. Gob encoding. File-based storage with metadata index.
- **Acceptance Criteria:**
  - [ ] `Create(trigger, description string, state *CompiledRuleSet) (*Snapshot, error)`
  - [ ] `List() ([]Snapshot, error)` — sorted by creation time
  - [ ] `Load(id string) (*Snapshot, *CompiledRuleSet, error)`
  - [ ] `Delete(id string) error`
  - [ ] `Diff(id string, current *CompiledRuleSet) (*ExecutionPlan, error)`
  - [ ] Gob encoding for serialization
  - [ ] File naming: `{uuid}-{trigger}.snap`
  - [ ] Metadata index file: `snapshots.json` (fast listing without reading all files)
  - [ ] Tests for full lifecycle

### T-017: Snapshot Retention & Cleanup
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-016
- **Description:** Automatic cleanup based on retention config (max count + max age). Background goroutine for periodic cleanup.
- **Acceptance Criteria:**
  - [ ] `Cleanup() error` — delete expired snapshots
  - [ ] Retention by count: keep last N snapshots
  - [ ] Retention by age: keep snapshots from last N days
  - [ ] Combined: keep if within count AND within age
  - [ ] Background goroutine with configurable interval
  - [ ] Tests with mock filesystem

### T-018: Audit Store
- **Priority:** P0
- **Estimate:** L
- **Dependencies:** T-002
- **Description:** `internal/audit/`. Append-only JSONL store with hash chain integrity. Record, search, verify. Daily file rotation.
- **Acceptance Criteria:**
  - [ ] `Record(event AuditEvent) error` — append to JSONL file
  - [ ] Hash chain: each entry's hash = SHA-256(prev_hash + entry_json)
  - [ ] `Search(query AuditQuery) ([]AuditEvent, int, error)` — filter by action, actor, time range
  - [ ] `Get(id string) (*AuditEvent, error)`
  - [ ] `VerifyIntegrity() (bool, error)` — verify hash chain
  - [ ] Daily file rotation: `audit-2026-04-11.jsonl`
  - [ ] Buffered channel → single writer goroutine
  - [ ] Gzip compression of files older than 24h
  - [ ] Retention cleanup (configurable max age)
  - [ ] Tests for write, search, integrity verification

---

# ═══════════════════════════════════════════════
# MILESTONE 7 — CLI Commands (5 tasks)
# ═══════════════════════════════════════════════

### T-019: CLI Framework (Flag Parsing + Output)
- **Priority:** P0
- **Estimate:** M
- **Dependencies:** T-001
- **Description:** `internal/cli/root.go`. Custom CLI framework using stdlib `flag.FlagSet`. Subcommand routing, global flags (--config, --output, --verbose), colorized output helpers.
- **Acceptance Criteria:**
  - [ ] Subcommand routing in main.go
  - [ ] Per-subcommand FlagSet
  - [ ] Global flags: `--config`, `-o/--output` (text/json), `--verbose`, `--no-color`
  - [ ] Color output helpers (ANSI codes, respects --no-color and NO_COLOR env)
  - [ ] `exitOnError(err, context)` helper
  - [ ] `confirm(prompt) bool` for interactive confirmation
  - [ ] Version subcommand with build info (ldflags)

### T-020: CLI — apply, plan, validate, fmt
- **Priority:** P0
- **Estimate:** L
- **Dependencies:** T-007, T-009, T-010, T-015, T-016, T-018, T-019
- **Description:** Implement `rampart apply`, `rampart plan`, `rampart validate`, `rampart fmt`.
- **Acceptance Criteria:**
  - [ ] `apply -f policy.yaml` — compile → conflict check → plan → confirm → snapshot → apply → audit
  - [ ] `apply --auto-approve` — skip confirmation
  - [ ] `apply --dry-run` — same as plan
  - [ ] `plan -f policy.yaml` — show execution plan without applying
  - [ ] `plan -o json` — JSON output for CI/CD
  - [ ] `validate -f policy.yaml` — validate schema only, exit 0/1
  - [ ] `fmt -f policy.yaml` — format YAML (consistent indentation, key ordering)
  - [ ] `fmt --check` — check if already formatted (exit 0/1 for CI)
  - [ ] Colorized plan output (green=add, yellow=modify, red=remove)

### T-021: CLI — rules (list, add, remove, stats)
- **Priority:** P0
- **Estimate:** M
- **Dependencies:** T-013, T-019
- **Description:** Implement `rampart rules` subcommands for quick rule management (bypass YAML workflow).
- **Acceptance Criteria:**
  - [ ] `rules list` — table format with name, priority, protocol, ports, source, action
  - [ ] `rules list -o json` — JSON output
  - [ ] `rules add --name X --protocol tcp --dport 22 --source 10.0.0.0/8 --action accept`
  - [ ] `rules add --ttl 2h` — time-to-live for temporary rules
  - [ ] `rules remove --name X` or `rules remove --id UUID`
  - [ ] `rules stats` — table with packet/byte counters per rule
  - [ ] `rules stats --name X` — stats for specific rule

### T-022: CLI — snapshot, rollback
- **Priority:** P0
- **Estimate:** M
- **Dependencies:** T-016, T-019
- **Description:** Implement snapshot and rollback CLI commands.
- **Acceptance Criteria:**
  - [ ] `snapshot list` — table with ID, created, trigger, rules count
  - [ ] `snapshot create --description "before maintenance"` — manual snapshot
  - [ ] `snapshot diff ID` — show diff between snapshot and current state
  - [ ] `snapshot export ID -o policy.yaml` — reverse-compile to YAML
  - [ ] `rollback ID` — confirm → rollback → audit
  - [ ] `rollback --last` — shortcut for most recent snapshot
  - [ ] `rollback --auto-approve` — skip confirmation

### T-023: CLI — audit, import, export, diff, simulate
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-018, T-019
- **Description:** Implement remaining CLI commands.
- **Acceptance Criteria:**
  - [ ] `audit list --last 20` — recent audit events
  - [ ] `audit list --action policy.apply --since 2026-04-01`
  - [ ] `audit list --actor ersin`
  - [ ] `audit show ID` — full event detail with before/after diff
  - [ ] `import --from iptables-save -o rules.yaml` — import iptables-save output to YAML
  - [ ] `import --from nftables -o rules.yaml` — import current nftables rules
  - [ ] `export -o current.yaml` — export current Rampart rules as YAML
  - [ ] `diff file1.yaml file2.yaml` — diff two policy files
  - [ ] `simulate --src IP --dst IP --protocol tcp --dport PORT --direction inbound`
  - [ ] Simulation output: verdict, matched rule, evaluation trace

---

# ═══════════════════════════════════════════════
# MILESTONE 8 — iptables Backend (3 tasks)
# ═══════════════════════════════════════════════

### T-024: iptables Backend — Probe & CurrentState
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-012
- **Description:** `internal/backend/iptables/`. Probe iptables binary. Parse `iptables-save` output to extract Rampart-managed rules.
- **Acceptance Criteria:**
  - [ ] `Probe()` checks iptables binary, ip6tables binary
  - [ ] `CurrentState()` runs `iptables-save` + `ip6tables-save`
  - [ ] Parse output to identify RAMPART-* chains
  - [ ] Extract rules with "rampart:" comment
  - [ ] Convert to CompiledRule format
  - [ ] Capabilities: IPv4=true, IPv6=true (via ip6tables), AtomicReplace=false

### T-025: iptables Backend — Rule Compiler
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-024
- **Description:** Translate CompiledRule → iptables arguments. Handle dual IPv4/IPv6.
- **Acceptance Criteria:**
  - [ ] `ruleToArgs(rule CompiledRule, chain string) []string`
  - [ ] `-p` protocol flag
  - [ ] `--dport` destination port (with `-m multiport` for multiple)
  - [ ] `-s` source CIDR
  - [ ] `-d` destination CIDR
  - [ ] `-m state --state` for connection tracking
  - [ ] `-i`/`-o` interface
  - [ ] `-m limit --limit` for rate limiting
  - [ ] `-m comment --comment "rampart:name"` for identification
  - [ ] `-j ACCEPT/DROP/REJECT/LOG` target
  - [ ] LOG prefix: `--log-prefix "rampart:name: "`
  - [ ] Tests with expected iptables command output

### T-026: iptables Backend — Apply (Chain Swap), Rollback, Stats
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-025
- **Description:** Apply via chain swap strategy (create new → populate → swap → cleanup). Rollback, Flush, Stats.
- **Acceptance Criteria:**
  - [ ] Create RAMPART-*-NEW chains
  - [ ] Populate new chains with all rules
  - [ ] Swap jump targets (add new, remove old)
  - [ ] Rename chains (old → delete, new → active)
  - [ ] Cleanup on failure (delete new chains)
  - [ ] `Rollback()` via chain swap
  - [ ] `Flush()` removes all RAMPART chains + jump rules
  - [ ] `Stats()` parses `iptables -L -v -n` for packet/byte counters

---

# ═══════════════════════════════════════════════
# MILESTONE 9 — Packet Simulator (2 tasks)
# ═══════════════════════════════════════════════

### T-027: Simulation Engine
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-007
- **Description:** `internal/engine/simulator.go`. Evaluate a simulated packet against compiled ruleset. Return verdict, matched rule, evaluation trace.
- **Acceptance Criteria:**
  - [ ] `Simulate(rules []CompiledRule, pkt SimulatedPacket) SimulationResult`
  - [ ] Direction filtering
  - [ ] Schedule evaluation (skip inactive rules)
  - [ ] Protocol matching
  - [ ] Source CIDR matching (net.IPNet.Contains)
  - [ ] Dest CIDR matching
  - [ ] Source port matching
  - [ ] Dest port matching
  - [ ] Interface matching
  - [ ] Connection state matching
  - [ ] Negation (NOT) matching
  - [ ] First-match wins (priority order)
  - [ ] Default deny if no match
  - [ ] Build human-readable match path
  - [ ] Track evaluation count and duration
  - [ ] Tests with comprehensive packet scenarios

### T-028: Import from iptables-save / nft list
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-007
- **Description:** Parse existing firewall rules and convert to Rampart YAML policy format.
- **Acceptance Criteria:**
  - [ ] `ImportIptablesSave(data []byte) (*PolicySetYAML, error)`
  - [ ] Parse -A rules: chain, protocol, ports, CIDRs, target
  - [ ] Group rules by chain → policies
  - [ ] Assign priorities based on rule order
  - [ ] `ImportNftables(data []byte) (*PolicySetYAML, error)`
  - [ ] Parse nft JSON output → policies
  - [ ] Output valid Rampart YAML
  - [ ] Tests with real-world iptables-save samples

---

# ═══════════════════════════════════════════════
# MILESTONE 10 — Configuration & Server (3 tasks)
# ═══════════════════════════════════════════════

### T-029: Configuration System
- **Priority:** P0
- **Estimate:** M
- **Dependencies:** T-002
- **Description:** `internal/config/`. Load config from YAML file (search order: --config flag, ./rampart.yaml, /etc/rampart/rampart.yaml, ~/.config/rampart/rampart.yaml). Environment variable overrides with RAMPART_ prefix. Default values.
- **Acceptance Criteria:**
  - [ ] `LoadConfig(paths ...string) (*Config, error)`
  - [ ] Config struct matching SPECIFICATION §21
  - [ ] Search path order
  - [ ] Env override: `RAMPART_LISTEN`, `RAMPART_BACKEND_TYPE`, etc.
  - [ ] Sensible defaults for all fields
  - [ ] Validation of config values
  - [ ] Tests with various config scenarios

### T-030: REST API Server
- **Priority:** P0
- **Estimate:** XL
- **Dependencies:** T-007, T-012, T-016, T-018, T-029
- **Description:** `internal/api/`. Custom HTTP router (trie-based). All endpoints from SPECIFICATION §16. Middleware: request ID, logging, CORS, auth.
- **Acceptance Criteria:**
  - [ ] Custom Router with method routing and path params (`:id`)
  - [ ] Middleware chain: requestID → logging → CORS → auth → handler
  - [ ] API key auth: `Authorization: Bearer rmp_xxx`
  - [ ] Unix socket support for local access (no auth)
  - [ ] All policy endpoints (plan, apply, simulate, current, conflicts)
  - [ ] All rule endpoints (list, add, get, delete, stats)
  - [ ] All snapshot endpoints (list, create, rollback, diff, export)
  - [ ] All audit endpoints (list, get, search)
  - [ ] System endpoints (info, health, backends, metrics)
  - [ ] Standard response envelope (status, data, error, meta)
  - [ ] Proper HTTP status codes
  - [ ] Request/response JSON serialization
  - [ ] Tests for each endpoint

### T-031: SSE (Server-Sent Events) for Real-Time Updates
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-030
- **Description:** `internal/api/sse.go`. SSE endpoint for real-time events (rule changes, audit events, cluster status). Pub/sub pattern for event distribution.
- **Acceptance Criteria:**
  - [ ] `GET /api/v1/events` — SSE stream
  - [ ] Event types: `rule.applied`, `rule.expired`, `audit.event`, `cluster.change`, `snapshot.created`
  - [ ] Subscribe/Unsubscribe pattern
  - [ ] Broadcast to all connected clients
  - [ ] Automatic reconnection support (Last-Event-ID)
  - [ ] Connection cleanup on client disconnect
  - [ ] Max connections limit (configurable)

---

# ═══════════════════════════════════════════════
# MILESTONE 11 — React WebUI (6 tasks)
# ═══════════════════════════════════════════════

### T-032: React Project Setup
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-030
- **Description:** Initialize React 19 + TypeScript + Vite + Tailwind CSS v4. Configure build output to `ui/dist/`. Setup Go embed for serving.
- **Acceptance Criteria:**
  - [ ] `ui/` directory with React 19 + TypeScript
  - [ ] Vite build config (output to `ui/dist/`)
  - [ ] Tailwind CSS v4 configuration
  - [ ] API client module with auth header injection
  - [ ] React Router for client-side routing
  - [ ] `//go:embed ui/dist/*` in Go server
  - [ ] Serve UI at `/ui/` path
  - [ ] SPA fallback (all non-API routes → index.html)
  - [ ] Dark/light mode support
  - [ ] `make ui-build` target

### T-033: Dashboard Page
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-032
- **Description:** Overview page: stats cards (rules count, cluster nodes, snapshots, recent audit events), rule hit heatmap (optional D3.js), recent activity feed (SSE-powered).
- **Acceptance Criteria:**
  - [ ] Stats cards: active rules, cluster nodes, snapshots count, backend type
  - [ ] Recent audit events feed (auto-updates via SSE)
  - [ ] Quick actions: apply policy, create snapshot, quick add rule
  - [ ] Responsive layout (mobile-friendly)
  - [ ] Loading skeletons during data fetch

### T-034: Policies & Rules Pages
- **Priority:** P1
- **Estimate:** XL
- **Dependencies:** T-032
- **Description:** Policy editor (CodeMirror 6 with YAML highlighting), live validation, conflict panel. Rules table (sortable, filterable, searchable), rule detail view, quick add rule form.
- **Acceptance Criteria:**
  - [ ] Policy editor: CodeMirror 6 with YAML mode
  - [ ] Live validation (debounced, show errors inline)
  - [ ] Plan preview (compile → show diff)
  - [ ] Apply button with confirmation dialog
  - [ ] Conflict panel (sidebar with warnings/errors)
  - [ ] Rules table: columns (name, priority, protocol, ports, source, action, stats)
  - [ ] Sortable by any column
  - [ ] Filterable by protocol, action, priority range
  - [ ] Searchable by name
  - [ ] Rule detail slide-over (full rule info + stats)
  - [ ] Quick add rule form (modal)
  - [ ] Delete rule with confirmation

### T-035: Simulator Page
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-032
- **Description:** Interactive packet simulation form. Visual trace of rule evaluation. Result card with verdict.
- **Acceptance Criteria:**
  - [ ] Form: source IP, dest IP, protocol, source port, dest port, direction, interface
  - [ ] Protocol dropdown: TCP, UDP, ICMP, ICMPv6
  - [ ] Direction dropdown: Inbound, Outbound, Forward
  - [ ] "Simulate" button → POST /api/v1/policies/simulate
  - [ ] Result card: ACCEPT (green) / DROP (red) / REJECT (yellow)
  - [ ] Matched rule display (name, policy, priority)
  - [ ] Evaluation trace: list of rules checked, matched/skipped
  - [ ] Match path visualization (which conditions matched)

### T-036: Snapshots & Audit Pages
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-032
- **Description:** Snapshot timeline, diff viewer, one-click rollback. Audit log table with search, filter, event detail with before/after diff.
- **Acceptance Criteria:**
  - [ ] Snapshot list: timeline or table view
  - [ ] Snapshot detail: metadata, rule count, trigger
  - [ ] Diff viewer: side-by-side comparison (snapshot vs current)
  - [ ] Rollback button with confirmation dialog
  - [ ] Create snapshot button
  - [ ] Audit log table: paginated (20 per page)
  - [ ] Filter by: action, actor, date range
  - [ ] Search by keyword
  - [ ] Audit event detail: full event with before/after JSON diff
  - [ ] Inline diff highlighting (added=green, removed=red)

### T-037: Cluster & Settings Pages
- **Priority:** P2
- **Estimate:** L
- **Dependencies:** T-032
- **Description:** Cluster node list, health indicators, Raft status. Settings page for backend config, API keys, snapshot retention.
- **Acceptance Criteria:**
  - [ ] Node list: ID, state (leader/follower), backend, rules count, last sync, health
  - [ ] Health indicators: green/yellow/red dots
  - [ ] Raft status: current term, commit index, log size
  - [ ] Settings: read-only display of current config
  - [ ] API key management: list, create, revoke
  - [ ] Snapshot retention settings
  - [ ] Backend selection (if multiple available)

---

# ═══════════════════════════════════════════════
# MILESTONE 12 — Time-Based Rules (2 tasks)
# ═══════════════════════════════════════════════

### T-038: Schedule Evaluator
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-007
- **Description:** `internal/engine/scheduler.go`. Evaluate Schedule structs. Support one-time (activeFrom/Until) and recurring (days + time range + timezone). Deterministic evaluation (same result on all cluster nodes).
- **Acceptance Criteria:**
  - [ ] `IsActive(sched *Schedule, now time.Time) bool`
  - [ ] One-time: check activeFrom ≤ now ≤ activeUntil
  - [ ] Recurring: check day of week + time of day in timezone
  - [ ] Timezone support via time.LoadLocation
  - [ ] Edge cases: nil From (immediately active), nil Until (permanent)
  - [ ] Deterministic: no random, no local state
  - [ ] Tests with mock time for all edge cases

### T-039: Background Scheduler Service
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-038, T-015
- **Description:** Background goroutine that periodically evaluates scheduled rules and activates/deactivates them. Generates audit events for transitions.
- **Acceptance Criteria:**
  - [ ] `Scheduler.Run(ctx context.Context)` — background loop
  - [ ] Configurable check interval (default 30s)
  - [ ] Detect rule activation (was inactive, now active → apply)
  - [ ] Detect rule deactivation (was active, now expired → remove)
  - [ ] Reapply full ruleset on any change
  - [ ] Audit events: `rule.activated`, `rule.expired`
  - [ ] Graceful shutdown via context cancellation
  - [ ] Tests with fast-forwarded clock

---

# ═══════════════════════════════════════════════
# MILESTONE 13 — Raft Cluster (5 tasks)
# ═══════════════════════════════════════════════

### T-040: Raft Core — State Machine & Log
- **Priority:** P1
- **Estimate:** XL
- **Dependencies:** T-002
- **Description:** `internal/cluster/raft.go`, `internal/cluster/log.go`. Core Raft consensus: leader election, log replication, commit. In-memory log with WAL persistence.
- **Acceptance Criteria:**
  - [ ] Node states: Follower, Candidate, Leader
  - [ ] Leader election with randomized timeout (150-300ms)
  - [ ] RequestVote RPC
  - [ ] AppendEntries RPC (heartbeat + log replication)
  - [ ] Log entry types: PolicyUpdate, ConfigChange, NodeJoin, NodeLeave
  - [ ] Commit index tracking
  - [ ] WAL persistence: append-only file for durability
  - [ ] Log compaction via snapshots
  - [ ] Tests with 3-node simulated cluster

### T-041: Raft Transport — TCP + TLS
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-040
- **Description:** `internal/cluster/transport.go`. TCP transport with mTLS for Raft RPCs. Binary frame protocol for efficiency.
- **Acceptance Criteria:**
  - [ ] TCP listener with TLS (mutual auth)
  - [ ] Connection pooling per peer
  - [ ] Binary frame: [4-byte length][type byte][gob-encoded payload]
  - [ ] RPC types: RequestVote, RequestVoteResponse, AppendEntries, AppendEntriesResponse
  - [ ] Connection timeout (5s dial, 10s idle)
  - [ ] Reconnection with exponential backoff
  - [ ] Tests with local TCP connections

### T-042: Raft FSM — Policy Application
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-040, T-015
- **Description:** `internal/cluster/fsm.go`. Finite state machine that applies committed Raft log entries. On PolicyUpdate → compile and apply to local backend.
- **Acceptance Criteria:**
  - [ ] `Apply(entry LogEntry) error` — apply committed entry
  - [ ] PolicyUpdate: deserialize → compile → apply to local backend
  - [ ] ConfigChange: update local config
  - [ ] NodeJoin/NodeLeave: update peer list
  - [ ] Snapshot: serialize current state for new nodes
  - [ ] Restore: load state from snapshot
  - [ ] Audit logging for all applied entries

### T-043: Cluster CLI Commands
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-040, T-041
- **Description:** Implement `rampart cluster` subcommands: init, join, leave, status, elect.
- **Acceptance Criteria:**
  - [ ] `cluster init --listen ADDR --advertise ADDR` — bootstrap new cluster
  - [ ] `cluster join --leader ADDR --listen ADDR` — join existing cluster
  - [ ] `cluster leave` — graceful leave (transfer leadership if leader)
  - [ ] `cluster status` — table: node, state, backend, rules, last-sync, health
  - [ ] `cluster elect --force` — force new election (emergency)

### T-044: Certificate Management
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-001
- **Description:** `internal/cert/`. Generate cluster CA and node certificates using crypto/x509 (stdlib). CLI commands: `rampart cert init`, `rampart cert generate`.
- **Acceptance Criteria:**
  - [ ] `cert init --ca-dir DIR` — generate CA key + cert (Ed25519)
  - [ ] `cert generate --node-name NAME --ca-dir DIR` — generate node cert signed by CA
  - [ ] Certificate validity: CA=10 years, node=1 year
  - [ ] SAN: node name + IP addresses + DNS names
  - [ ] Output: PEM format (.crt + .key files)
  - [ ] Verify cert against CA
  - [ ] Tests for full certificate lifecycle

---

# ═══════════════════════════════════════════════
# MILESTONE 14 — MCP Server (2 tasks)
# ═══════════════════════════════════════════════

### T-045: MCP Server Implementation
- **Priority:** P2
- **Estimate:** L
- **Dependencies:** T-030
- **Description:** `internal/mcp/`. JSON-RPC 2.0 based MCP server. Expose tools: list_rules, add_rule, remove_rule, plan_policy, apply_policy, simulate_packet, rollback, list_snapshots, audit_search, cluster_status, get_rule_stats.
- **Acceptance Criteria:**
  - [ ] JSON-RPC 2.0 transport (stdio or TCP)
  - [ ] Tool definitions with JSON Schema parameter descriptions
  - [ ] All 11 tools from SPECIFICATION §19.1
  - [ ] Resource definitions from SPECIFICATION §19.2
  - [ ] apply_policy requires explicit confirmation parameter
  - [ ] Error handling with MCP error codes
  - [ ] Tests for each tool

### T-046: MCP Resources
- **Priority:** P2
- **Estimate:** M
- **Dependencies:** T-045
- **Description:** MCP resource endpoints: rampart://policies/current, rampart://rules, rampart://audit/recent, rampart://cluster/status.
- **Acceptance Criteria:**
  - [ ] `rampart://policies/current` → current YAML policy
  - [ ] `rampart://rules` → active rules JSON
  - [ ] `rampart://audit/recent` → last 50 audit events
  - [ ] `rampart://cluster/status` → cluster health
  - [ ] Resource change notifications via MCP protocol

---

# ═══════════════════════════════════════════════
# MILESTONE 15 — Observability (2 tasks)
# ═══════════════════════════════════════════════

### T-047: Prometheus Metrics Endpoint
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-030
- **Description:** `GET /metrics` endpoint in Prometheus exposition format. No prometheus/client_golang — implement text format from scratch.
- **Acceptance Criteria:**
  - [ ] Custom Prometheus text format writer
  - [ ] Rule metrics: packets_total, bytes_total per rule
  - [ ] Backend metrics: apply_duration_seconds, apply_total
  - [ ] Cluster metrics: raft_term, raft_commit_index, raft_peers, raft_state
  - [ ] Snapshot metrics: snapshots_total, snapshot_size_bytes
  - [ ] Audit metrics: audit_events_total per action
  - [ ] Scheduler metrics: scheduled_rules_active, scheduled_rules_total
  - [ ] System metrics: uptime_seconds, go_goroutines
  - [ ] Tests verifying exposition format

### T-048: Structured Logging
- **Priority:** P0
- **Estimate:** M
- **Dependencies:** T-001
- **Description:** Custom structured logger (no slog external deps — use Go 1.21+ slog if available, else custom). JSON and text output formats. Log levels. Component tagging.
- **Acceptance Criteria:**
  - [ ] Log levels: debug, info, warn, error
  - [ ] JSON format: `{"time": "...", "level": "info", "msg": "...", "component": "engine", ...}`
  - [ ] Text format: `2026-04-11T10:30:00 INFO [engine] policy applied rules=12`
  - [ ] Component field (engine, api, cluster, backend, scheduler, audit)
  - [ ] Structured fields (key=value)
  - [ ] File output option
  - [ ] Log rotation (optional, recommend external tool)
  - [ ] Configurable via config file

---

# ═══════════════════════════════════════════════
# MILESTONE 16 — eBPF Backend (3 tasks)
# ═══════════════════════════════════════════════

### T-049: eBPF Backend — Loader & Maps
- **Priority:** P2
- **Estimate:** XL
- **Dependencies:** T-012
- **Description:** `internal/backend/ebpf/`. eBPF program loader (via bpf() syscall — no cilium/ebpf). BPF map management: LPM trie for CIDRs, hash map for ports, percpu array for counters.
- **Acceptance Criteria:**
  - [ ] BPF syscall wrapper (golang.org/x/sys)
  - [ ] Program load from pre-compiled bytecode (ELF)
  - [ ] Map creation: LPM_TRIE, HASH, PERCPU_ARRAY, LRU_HASH
  - [ ] Map CRUD operations
  - [ ] XDP attach/detach to network interface
  - [ ] Probe: check kernel BPF support
  - [ ] Capabilities: RateLimiting=true, ConnectionTracking=limited

### T-050: eBPF Backend — XDP Program
- **Priority:** P2
- **Estimate:** XL
- **Dependencies:** T-049
- **Description:** Pre-compiled XDP program for packet filtering. Parse Ethernet → IP → TCP/UDP headers. Match against BPF maps. Return XDP_DROP/XDP_PASS.
- **Acceptance Criteria:**
  - [ ] CO-RE compatible BPF C source
  - [ ] Ethernet header parsing
  - [ ] IPv4/IPv6 header parsing
  - [ ] TCP/UDP header parsing
  - [ ] LPM trie lookup for source/dest CIDR matching
  - [ ] Hash map lookup for port matching
  - [ ] Per-rule packet/byte counters
  - [ ] Token bucket rate limiting in BPF
  - [ ] Perf event output for logging
  - [ ] Pre-compiled .o files for common architectures (amd64, arm64)

### T-051: Hybrid Backend (eBPF + nftables)
- **Priority:** P2
- **Estimate:** L
- **Dependencies:** T-015, T-049
- **Description:** Hybrid backend that routes rules to eBPF (fast path) or nftables (slow path) based on rule characteristics.
- **Acceptance Criteria:**
  - [ ] Rule classification: simple (eBPF-capable) vs complex (nftables)
  - [ ] Configuration: which rules go to fast path
  - [ ] Coordinated apply: eBPF first, then nftables
  - [ ] Combined stats from both backends
  - [ ] Fallback: if eBPF unavailable, all rules to nftables

---

# ═══════════════════════════════════════════════
# MILESTONE 17 — Cloud Backends (3 tasks)
# ═══════════════════════════════════════════════

### T-052: AWS Security Groups Backend
- **Priority:** P2
- **Estimate:** XL
- **Dependencies:** T-012
- **Description:** `internal/backend/aws/`. AWS SigV4 signing (from scratch). Security Group rule management via EC2 API.
- **Acceptance Criteria:**
  - [ ] AWS SigV4 request signing (no aws-sdk-go)
  - [ ] HTTP client for EC2 API
  - [ ] DescribeSecurityGroupRules → CurrentState
  - [ ] AuthorizeSecurityGroupIngress/Egress → Apply
  - [ ] RevokeSecurityGroupIngress/Egress → remove rules
  - [ ] Rule translation: CompiledRule → SG rule format
  - [ ] Respect SG limits (60 inbound + 60 outbound)
  - [ ] Exponential backoff for API rate limiting
  - [ ] Eventual consistency: poll until converged

### T-053: GCP Firewall Rules Backend
- **Priority:** P2
- **Estimate:** L
- **Dependencies:** T-012
- **Description:** `internal/backend/gcp/`. GCP OAuth2 (from scratch). Compute Engine Firewall Rules API.
- **Acceptance Criteria:**
  - [ ] OAuth2 token acquisition (service account JSON key)
  - [ ] Compute Engine REST API client
  - [ ] List firewall rules → CurrentState
  - [ ] Create/Delete/Patch firewall rules → Apply
  - [ ] Priority mapping (Rampart 0-999 → GCP 0-65535)
  - [ ] Target tags for rule scoping
  - [ ] Eventual consistency handling

### T-054: Azure NSG Backend
- **Priority:** P2
- **Estimate:** L
- **Dependencies:** T-012
- **Description:** `internal/backend/azure/`. Azure AD OAuth2 (from scratch). NSG Security Rules API.
- **Acceptance Criteria:**
  - [ ] Azure AD OAuth2 token acquisition
  - [ ] NSG REST API client
  - [ ] List security rules → CurrentState
  - [ ] Create/Delete/Update rules → Apply
  - [ ] Priority mapping (Rampart 0-999 → Azure 100-4096)
  - [ ] Eventual consistency handling

---

# ═══════════════════════════════════════════════
# MILESTONE 18 — Production Hardening (5 tasks)
# ═══════════════════════════════════════════════

### T-055: Security Hardening
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** T-030
- **Description:** Capability dropping (CAP_NET_ADMIN + CAP_NET_RAW only), API key bcrypt hashing, rate limiting per API key, hash chain integrity verification.
- **Acceptance Criteria:**
  - [ ] Drop all capabilities except CAP_NET_ADMIN, CAP_NET_RAW after startup
  - [ ] API keys stored as bcrypt hashes
  - [ ] Per-key rate limiting (configurable)
  - [ ] Hash chain integrity check: `rampart audit verify`
  - [ ] TLS minimum version 1.2
  - [ ] HSTS header on WebUI

### T-056: Benchmark Suite
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-007, T-010, T-027
- **Description:** Go benchmarks for critical paths: compilation, conflict detection, simulation, snapshot create/restore.
- **Acceptance Criteria:**
  - [ ] `BenchmarkCompile100Rules`
  - [ ] `BenchmarkCompile10000Rules`
  - [ ] `BenchmarkConflictDetection1000Rules`
  - [ ] `BenchmarkSimulatePacket`
  - [ ] `BenchmarkSnapshotCreate`
  - [ ] `BenchmarkSnapshotRestore`
  - [ ] Results must meet SPECIFICATION §24 targets
  - [ ] CI integration: benchmark regression detection

### T-057: Docker Image & Systemd Service
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-001
- **Description:** Multi-stage Dockerfile (build + scratch runtime). Systemd service file. Docker Compose example.
- **Acceptance Criteria:**
  - [ ] Multi-stage Dockerfile: Go build → scratch/alpine runtime
  - [ ] Image size < 50 MB
  - [ ] `rampart.service` systemd unit file
  - [ ] After=network.target, Restart=on-failure
  - [ ] `docker-compose.yml` example (3-node cluster)
  - [ ] Health check in Docker + systemd

### T-058: CI/CD Pipeline
- **Priority:** P1
- **Estimate:** M
- **Dependencies:** T-001
- **Description:** GitHub Actions: lint, test, build, release. Matrix build (linux/amd64, linux/arm64). Goreleaser for releases.
- **Acceptance Criteria:**
  - [ ] `.github/workflows/ci.yml`: go vet, staticcheck, test, build
  - [ ] Matrix: Go versions, OS/arch
  - [ ] Race detector enabled in tests
  - [ ] Coverage report
  - [ ] Release workflow: goreleaser, Docker push
  - [ ] Artifact: single binary per platform

### T-059: Documentation & README
- **Priority:** P1
- **Estimate:** L
- **Dependencies:** All previous
- **Description:** Final README.md polish, man pages (`rampart.1`), quick start guide, configuration reference, API reference.
- **Acceptance Criteria:**
  - [ ] README.md: ASCII architecture diagram, feature list, quick start, comparison table
  - [ ] man page: `rampart(1)` with all subcommands
  - [ ] Quick start: install → first policy → apply → verify
  - [ ] Configuration reference: all YAML keys with descriptions
  - [ ] API reference: all endpoints with request/response examples
  - [ ] Contributing guide

---

## Summary

| Milestone | Tasks | Priority | Est. Duration |
|-----------|-------|----------|---------------|
| M1 Project Scaffold | T-001, T-002 | P0 | 1 day |
| M2 YAML Parser | T-003 — T-006 | P0/P1 | 2 days |
| M3 Rule Compiler | T-007 — T-009 | P0 | 2 days |
| M4 Conflict Detection | T-010, T-011 | P0/P1 | 2 days |
| M5 nftables Backend | T-012 — T-015 | P0 | 3 days |
| M6 Snapshot & Audit | T-016 — T-018 | P0/P1 | 2 days |
| M7 CLI Commands | T-019 — T-023 | P0/P1 | 3 days |
| M8 iptables Backend | T-024 — T-026 | P1 | 2 days |
| M9 Packet Simulator | T-027, T-028 | P1 | 2 days |
| M10 Config & Server | T-029 — T-031 | P0/P1 | 3 days |
| M11 React WebUI | T-032 — T-037 | P1/P2 | 5 days |
| M12 Time-Based Rules | T-038, T-039 | P1 | 1 day |
| M13 Raft Cluster | T-040 — T-044 | P1 | 5 days |
| M14 MCP Server | T-045, T-046 | P2 | 2 days |
| M15 Observability | T-047, T-048 | P0/P1 | 1 day |
| M16 eBPF Backend | T-049 — T-051 | P2 | 4 days |
| M17 Cloud Backends | T-052 — T-054 | P2 | 4 days |
| M18 Production Hardening | T-055 — T-059 | P1 | 3 days |
| **Total** | **59 tasks** | | **~45 days** |

**MVP (P0 only):** M1 + M2 + M3 + M4 + M5 + M6 + M7 + M10 + M15 = ~19 days
**Full v1.0:** All milestones = ~45 days (~9 weeks)
