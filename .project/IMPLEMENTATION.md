# Rampart — Implementation Guide

> **Version:** 0.1.0  
> **Companion to:** Rampart SPECIFICATION.md v0.1.0  
> **Scope:** Phase 1–3 MVP (Core Engine + iptables + API + WebUI)  
> **Estimated LOC:** ~12,000 Go + ~8,000 TypeScript/React  
> **Timeline:** 8 weeks

---

## Table of Contents

1. [Architecture Decisions](#1-architecture-decisions)
2. [Core Data Structures](#2-core-data-structures)
3. [YAML Policy Parser](#3-yaml-policy-parser)
4. [Rule Compiler](#4-rule-compiler)
5. [Conflict Detection Engine](#5-conflict-detection-engine)
6. [Backend Abstraction Layer](#6-backend-abstraction-layer)
7. [nftables Backend Implementation](#7-nftables-backend-implementation)
8. [iptables Backend Implementation](#8-iptables-backend-implementation)
9. [Snapshot Engine](#9-snapshot-engine)
10. [Audit System](#10-audit-system)
11. [Packet Simulator](#11-packet-simulator)
12. [Time-Based Rule Scheduler](#12-time-based-rule-scheduler)
13. [REST API Server](#13-rest-api-server)
14. [CLI Implementation](#14-cli-implementation)
15. [React WebUI Architecture](#15-react-webui-architecture)
16. [Configuration Loading](#16-configuration-loading)
17. [Error Handling Strategy](#17-error-handling-strategy)
18. [Testing Strategy](#18-testing-strategy)

---

## 1. Architecture Decisions

### AD-1: Unified Binary Architecture

Rampart is a single binary that operates in three modes based on the subcommand:

```
rampart serve   → Server mode (API + WebUI + engine + optional Raft)
rampart agent   → Agent mode (follower-only, receives policies from leader)
rampart apply   → One-shot CLI (compile YAML → apply to local backend → exit)
```

All modes share the same compiled binary. Mode is determined at runtime by the subcommand.

### AD-2: Policy Engine as Pure Functions

The policy engine (YAML → compiled rules) is implemented as pure functions with no side effects. This enables:
- Deterministic compilation (same input → same output on all nodes)
- Easy testing (no mocking required)
- Dry-run mode (compile without applying)

```go
// Pure function: no I/O, no side effects
func Compile(policySet PolicySet, vars Variables) (CompiledRuleSet, []Conflict, error)
```

### AD-3: Backend as Stateless Adapter

Backends do NOT store state. They translate compiled rules to backend-specific commands and execute them. The source of truth is always the compiled policy in the engine.

### AD-4: CLI Parsing Strategy

Custom CLI parser (no Cobra/urfave). Simple recursive descent:

```go
func main() {
    args := os.Args[1:]
    if len(args) == 0 {
        printUsage()
        os.Exit(1)
    }
    switch args[0] {
    case "serve":
        runServe(args[1:])
    case "apply":
        runApply(args[1:])
    // ...
    }
}
```

Flags parsed with stdlib `flag.FlagSet` per subcommand.

### AD-5: HTTP Router

Custom trie-based HTTP router. No external dependencies.

```go
type Router struct {
    root *node
}

func (r *Router) Handle(method, pattern string, handler http.Handler)
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request)
```

Supports path parameters (`:id`), method routing, middleware chain.

### AD-6: JSON Handling

All JSON serialization uses stdlib `encoding/json`. Model structs have JSON tags. API responses use a standard envelope:

```go
type Response struct {
    Status string      `json:"status"`          // "success" or "error"
    Data   interface{} `json:"data,omitempty"`
    Error  *APIError   `json:"error,omitempty"`
    Meta   ResponseMeta `json:"meta"`
}
```

### AD-7: Concurrency Model

- Policy compilation: single-threaded (fast enough, deterministic)
- API server: `net/http` goroutine-per-request
- Backend operations: serialized through a single goroutine (avoid concurrent nft/iptables calls)
- Audit writes: buffered channel → single writer goroutine
- Scheduler: single goroutine with time.Ticker
- Raft: dedicated goroutine pool (Phase 4)

```go
type Engine struct {
    mu       sync.RWMutex     // Protects current policy state
    backendC chan backendOp     // Serialized backend operations
    auditC   chan AuditEvent    // Buffered audit writes
}
```

### AD-8: Embedded React UI

React app is built at compile time and embedded via `//go:embed`:

```go
//go:embed ui/dist/*
var uiAssets embed.FS
```

Build process: `cd ui && npm run build` → output to `ui/dist/` → `go build` embeds it.

Makefile target:

```makefile
build: ui-build
	go build -o rampart ./cmd/rampart

ui-build:
	cd ui && npm ci && npm run build
```

---

## 2. Core Data Structures

### 2.1 Policy Set (YAML representation)

```go
// PolicySetYAML is the top-level YAML structure
type PolicySetYAML struct {
    APIVersion string            `yaml:"apiVersion"`
    Kind       string            `yaml:"kind"`
    Metadata   PolicyMetadata    `yaml:"metadata"`
    Defaults   *PolicyDefaults   `yaml:"defaults,omitempty"`
    Includes   []IncludeRef      `yaml:"includes,omitempty"`
    Policies   []PolicyYAML      `yaml:"policies"`
}

type PolicyMetadata struct {
    Name        string            `yaml:"name"`
    Description string            `yaml:"description,omitempty"`
    Owner       string            `yaml:"owner,omitempty"`
    Tags        map[string]string `yaml:"tags,omitempty"`
}

type PolicyDefaults struct {
    Direction Direction `yaml:"direction,omitempty"`
    Action    Action    `yaml:"action,omitempty"`
    IPVersion IPVersion `yaml:"ipVersion,omitempty"`
    States    []string  `yaml:"states,omitempty"`
}

type PolicyYAML struct {
    Name        string      `yaml:"name"`
    Priority    int         `yaml:"priority"`
    Direction   Direction   `yaml:"direction,omitempty"`
    Description string      `yaml:"description,omitempty"`
    Rules       []RuleYAML  `yaml:"rules"`
}

type RuleYAML struct {
    Name        string       `yaml:"name"`
    Match       MatchYAML    `yaml:"match"`
    Action      Action       `yaml:"action"`
    Log         bool         `yaml:"log,omitempty"`
    RateLimit   *RateLimitYAML `yaml:"rateLimit,omitempty"`
    Schedule    *ScheduleYAML  `yaml:"schedule,omitempty"`
    Description string       `yaml:"description,omitempty"`
    Tags        map[string]string `yaml:"tags,omitempty"`
}

type MatchYAML struct {
    Protocol     interface{} `yaml:"protocol,omitempty"`     // string or []string
    SourceCIDRs  []string    `yaml:"sourceCIDRs,omitempty"`
    DestCIDRs    []string    `yaml:"destCIDRs,omitempty"`
    SourcePorts  interface{} `yaml:"sourcePorts,omitempty"`  // int, []int, or "start-end"
    DestPorts    interface{} `yaml:"destPorts,omitempty"`    // int, []int, or "start-end"
    Interfaces   []string    `yaml:"interfaces,omitempty"`
    States       []string    `yaml:"states,omitempty"`
    ICMPTypes    []int       `yaml:"icmpTypes,omitempty"`
    Not          *MatchYAML  `yaml:"not,omitempty"`
}
```

### 2.2 Compiled Rule (Internal representation)

```go
// CompiledRule is the normalized, validated, backend-agnostic rule
type CompiledRule struct {
    ID          string             // UUID v7 (generated)
    Name        string             // From YAML
    PolicyName  string             // Parent policy name
    Priority    int                // 0-999
    Direction   Direction          // Inbound, Outbound, Forward
    Action      Action             // Accept, Drop, Reject, Log, RateLimit
    Match       CompiledMatch
    Log         bool
    RateLimit   *RateLimit
    Schedule    *Schedule
    Tags        map[string]string
    Description string
    SourceFile  string             // Which YAML file this came from
    SourceLine  int                // Line number in YAML
}

type CompiledMatch struct {
    SourceNets  []net.IPNet
    DestNets    []net.IPNet
    SourcePorts []PortRange
    DestPorts   []PortRange
    Protocols   []Protocol
    Interfaces  []string
    States      []ConnState
    ICMPTypes   []uint8
    IPVersion   IPVersion
    Negated     *CompiledMatch     // NOT condition
}

type PortRange struct {
    Start uint16
    End   uint16
}

type RateLimit struct {
    Rate   int    // Packets per interval
    Per    string // "second", "minute"
    Burst  int
    Action Action // Action when limit exceeded
}
```

### 2.3 Compiled Rule Set

```go
type CompiledRuleSet struct {
    Rules       []CompiledRule     // Priority-sorted
    Hash        string             // SHA-256 of deterministic serialization
    CompiledAt  time.Time
    SourceFiles []string           // All YAML files involved
    Metadata    PolicyMetadata
}
```

### 2.4 Execution Plan

```go
type ExecutionPlan struct {
    ToAdd    []CompiledRule
    ToRemove []CompiledRule
    ToModify []RuleModification
    Warnings []Conflict
    Errors   []Conflict

    // Statistics
    CurrentRuleCount int
    PlannedRuleCount int
    AddCount         int
    RemoveCount      int
    ModifyCount      int
}

type RuleModification struct {
    Before CompiledRule
    After  CompiledRule
    Fields []string // Which fields changed
}
```

---

## 3. YAML Policy Parser

### 3.1 Parsing Pipeline

```go
func ParsePolicyFile(path string) (*PolicySetYAML, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("parse: read %s: %w", path, err)
    }
    
    var ps PolicySetYAML
    if err := yaml.Unmarshal(data, &ps); err != nil {
        return nil, fmt.Errorf("parse: unmarshal %s: %w", path, err)
    }
    
    if err := validateSchema(&ps); err != nil {
        return nil, fmt.Errorf("parse: validate %s: %w", path, err)
    }
    
    return &ps, nil
}
```

### 3.2 Schema Validation

Validation is performed programmatically (no JSON Schema dependency):

```go
func validateSchema(ps *PolicySetYAML) error {
    var errs []error
    
    if ps.APIVersion != "rampartfw.com/v1" {
        errs = append(errs, fmt.Errorf("unsupported apiVersion: %s", ps.APIVersion))
    }
    if ps.Kind != "PolicySet" {
        errs = append(errs, fmt.Errorf("unsupported kind: %s", ps.Kind))
    }
    if ps.Metadata.Name == "" {
        errs = append(errs, fmt.Errorf("metadata.name is required"))
    }
    
    policyNames := map[string]bool{}
    for i, p := range ps.Policies {
        if policyNames[p.Name] {
            errs = append(errs, fmt.Errorf("policies[%d]: duplicate name %q", i, p.Name))
        }
        policyNames[p.Name] = true
        
        if p.Priority < 0 || p.Priority > 999 {
            errs = append(errs, fmt.Errorf("policies[%d]: priority must be 0-999", i))
        }
        
        errs = append(errs, validateRules(i, p)...)
    }
    
    return errors.Join(errs...)
}
```

### 3.3 Variable Substitution

```go
var varPattern = regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)

func substituteVars(data []byte, vars map[string]interface{}) ([]byte, error) {
    result := varPattern.ReplaceAllFunc(data, func(match []byte) []byte {
        name := string(match[2 : len(match)-1]) // Extract variable name
        val, ok := vars[name]
        if !ok {
            return match // Leave unresolved (will be caught by validation)
        }
        b, _ := yaml.Marshal(val)
        return bytes.TrimSpace(b)
    })
    return result, nil
}
```

### 3.4 Include Resolution

```go
func resolveIncludes(ps *PolicySetYAML, basePath string, depth int) error {
    if depth > 10 {
        return fmt.Errorf("include depth exceeded (max 10), possible circular include")
    }
    
    for _, inc := range ps.Includes {
        var data []byte
        var err error
        
        if inc.URL != "" {
            data, err = fetchURL(inc.URL) // HTTPS only
        } else {
            path := inc.Path
            if !filepath.IsAbs(path) {
                path = filepath.Join(filepath.Dir(basePath), path)
            }
            data, err = os.ReadFile(path)
        }
        if err != nil {
            return fmt.Errorf("include %s: %w", inc.Path+inc.URL, err)
        }
        
        var included PolicySetYAML
        if err := yaml.Unmarshal(data, &included); err != nil {
            return fmt.Errorf("include parse: %w", err)
        }
        
        // Recursive include resolution
        if err := resolveIncludes(&included, inc.Path, depth+1); err != nil {
            return err
        }
        
        ps.Policies = append(ps.Policies, included.Policies...)
    }
    
    return nil
}
```

### 3.5 Port Parsing

YAML ports can be specified as int, []int, or "start-end" ranges:

```go
func parsePorts(raw interface{}) ([]PortRange, error) {
    switch v := raw.(type) {
    case nil:
        return nil, nil
    case int:
        return []PortRange{{Start: uint16(v), End: uint16(v)}}, nil
    case []interface{}:
        var ranges []PortRange
        for _, item := range v {
            switch p := item.(type) {
            case int:
                ranges = append(ranges, PortRange{Start: uint16(p), End: uint16(p)})
            case string:
                pr, err := parsePortRange(p) // "8000-9000" → PortRange{8000, 9000}
                if err != nil {
                    return nil, err
                }
                ranges = append(ranges, pr)
            }
        }
        return ranges, nil
    case string:
        pr, err := parsePortRange(v)
        return []PortRange{pr}, err
    default:
        return nil, fmt.Errorf("invalid port specification: %T", raw)
    }
}
```

---

## 4. Rule Compiler

### 4.1 Compilation Flow

```go
func Compile(ps *PolicySetYAML, vars map[string]interface{}) (*CompiledRuleSet, error) {
    // 1. Resolve includes
    if err := resolveIncludes(ps, "", 0); err != nil {
        return nil, err
    }
    
    // 2. Apply defaults
    applyDefaults(ps)
    
    // 3. Compile each rule
    var rules []CompiledRule
    for _, policy := range ps.Policies {
        for _, ruleYAML := range policy.Rules {
            compiled, err := compileRule(ruleYAML, policy)
            if err != nil {
                return nil, fmt.Errorf("policy %s, rule %s: %w",
                    policy.Name, ruleYAML.Name, err)
            }
            rules = append(rules, compiled)
        }
    }
    
    // 4. Sort by priority (stable sort: preserve order within same priority)
    sort.SliceStable(rules, func(i, j int) bool {
        return rules[i].Priority < rules[j].Priority
    })
    
    // 5. Generate IDs (deterministic: based on content hash)
    for i := range rules {
        rules[i].ID = generateRuleID(rules[i])
    }
    
    // 6. Compute ruleset hash
    hash := computeRuleSetHash(rules)
    
    return &CompiledRuleSet{
        Rules:      rules,
        Hash:       hash,
        CompiledAt: time.Now(),
        Metadata:   ps.Metadata,
    }, nil
}
```

### 4.2 CIDR Normalization

```go
func compileCIDRs(cidrs []string) ([]net.IPNet, error) {
    var nets []net.IPNet
    for _, cidr := range cidrs {
        // Handle bare IPs (add /32 or /128)
        if !strings.Contains(cidr, "/") {
            ip := net.ParseIP(cidr)
            if ip == nil {
                return nil, fmt.Errorf("invalid IP/CIDR: %s", cidr)
            }
            if ip.To4() != nil {
                cidr = cidr + "/32"
            } else {
                cidr = cidr + "/128"
            }
        }
        _, ipNet, err := net.ParseCIDR(cidr)
        if err != nil {
            return nil, fmt.Errorf("invalid CIDR: %s: %w", cidr, err)
        }
        nets = append(nets, *ipNet)
    }
    return nets, nil
}
```

### 4.3 Rule ID Generation (UUID v7)

```go
// UUID v7: time-sortable, random suffix
func generateUUIDv7() string {
    now := time.Now()
    ms := uint64(now.UnixMilli())
    
    var uuid [16]byte
    // Timestamp (48 bits)
    uuid[0] = byte(ms >> 40)
    uuid[1] = byte(ms >> 32)
    uuid[2] = byte(ms >> 24)
    uuid[3] = byte(ms >> 16)
    uuid[4] = byte(ms >> 8)
    uuid[5] = byte(ms)
    // Version 7
    uuid[6] = 0x70 | (uuid[6] & 0x0f)
    // Random (62 bits)
    rand.Read(uuid[8:])
    uuid[8] = 0x80 | (uuid[8] & 0x3f) // Variant 2
    
    return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
        uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}
```

---

## 5. Conflict Detection Engine

### 5.1 Interval Tree for CIDR Overlap

CIDR ranges are converted to integer intervals for overlap detection:

```go
type IPInterval struct {
    Start *big.Int // IP as big integer
    End   *big.Int // End of CIDR range
    RuleID string
}

// Convert CIDR to interval
func cidrToInterval(ipNet net.IPNet) IPInterval {
    start := ipToInt(ipNet.IP)
    ones, bits := ipNet.Mask.Size()
    size := new(big.Int).Lsh(big.NewInt(1), uint(bits-ones))
    end := new(big.Int).Add(start, size)
    end.Sub(end, big.NewInt(1))
    return IPInterval{Start: start, End: end}
}
```

### 5.2 Pairwise Conflict Check

```go
func DetectConflicts(rules []CompiledRule) []Conflict {
    var conflicts []Conflict
    
    for i := 0; i < len(rules); i++ {
        for j := i + 1; j < len(rules); j++ {
            if c := checkPair(rules[i], rules[j]); c != nil {
                conflicts = append(conflicts, *c)
            }
        }
    }
    
    return conflicts
}

func checkPair(a, b CompiledRule) *Conflict {
    // Same direction required for conflict
    if a.Direction != b.Direction {
        return nil
    }
    
    // Check protocol overlap
    if !protocolsOverlap(a.Match.Protocols, b.Match.Protocols) {
        return nil
    }
    
    // Check port overlap
    if !portsOverlap(a.Match.DestPorts, b.Match.DestPorts) {
        return nil
    }
    
    // Check CIDR overlap
    srcOverlap := cidrsOverlap(a.Match.SourceNets, b.Match.SourceNets)
    if !srcOverlap {
        return nil
    }
    
    // Overlap found — classify the conflict
    if a.Action == b.Action {
        if matchEquals(a.Match, b.Match) {
            return &Conflict{Type: ConflictRedundancy, RuleA: a, RuleB: b}
        }
        return &Conflict{Type: ConflictSubset, RuleA: a, RuleB: b}
    }
    
    if a.Priority == b.Priority {
        return &Conflict{Type: ConflictContradiction, RuleA: a, RuleB: b,
            Severity: SeverityError}
    }
    
    // Different priority, different action
    if matchSubset(b.Match, a.Match) {
        return &Conflict{Type: ConflictShadow, RuleA: a, RuleB: b,
            Severity: SeverityWarning}
    }
    
    return &Conflict{Type: ConflictOverlap, RuleA: a, RuleB: b,
        Severity: SeverityWarning}
}
```

### 5.3 CIDR Overlap Detection

```go
func cidrsOverlap(a, b []net.IPNet) bool {
    // Empty means "any" → always overlaps
    if len(a) == 0 || len(b) == 0 {
        return true
    }
    
    for _, netA := range a {
        for _, netB := range b {
            if netA.Contains(netB.IP) || netB.Contains(netA.IP) {
                return true
            }
        }
    }
    return false
}
```

---

## 6. Backend Abstraction Layer

### 6.1 Interface Implementation Pattern

```go
// Each backend registers itself via init()
func init() {
    backend.Register("nftables", NewNftablesBackend)
    backend.Register("iptables", NewIptablesBackend)
}

// Backend selection at runtime
func selectBackend(cfg config.BackendConfig) (Backend, error) {
    if cfg.Type == "auto" {
        return AutoDetect()
    }
    return NewBackend(cfg.Type, cfg)
}
```

### 6.2 Execution Plan Generation

```go
func GeneratePlan(current, desired *CompiledRuleSet) *ExecutionPlan {
    plan := &ExecutionPlan{
        CurrentRuleCount: len(current.Rules),
        PlannedRuleCount: len(desired.Rules),
    }
    
    currentMap := indexByID(current.Rules)
    desiredMap := indexByID(desired.Rules)
    
    // Rules to add (in desired but not in current)
    for id, rule := range desiredMap {
        if _, exists := currentMap[id]; !exists {
            plan.ToAdd = append(plan.ToAdd, rule)
        }
    }
    
    // Rules to remove (in current but not in desired)
    for id, rule := range currentMap {
        if _, exists := desiredMap[id]; !exists {
            plan.ToRemove = append(plan.ToRemove, rule)
        }
    }
    
    // Rules to modify (in both, but different)
    for id, desired := range desiredMap {
        if current, exists := currentMap[id]; exists {
            if !rulesEqual(current, desired) {
                plan.ToModify = append(plan.ToModify, RuleModification{
                    Before: current,
                    After:  desired,
                    Fields: diffFields(current, desired),
                })
            }
        }
    }
    
    plan.AddCount = len(plan.ToAdd)
    plan.RemoveCount = len(plan.ToRemove)
    plan.ModifyCount = len(plan.ToModify)
    
    return plan
}
```

---

## 7. nftables Backend Implementation

### 7.1 Command Execution

```go
type NftablesBackend struct {
    binary    string // Path to nft binary (default: /usr/sbin/nft)
    tableName string // Rampart table name (default: "rampart")
}

func (b *NftablesBackend) exec(args ...string) ([]byte, error) {
    cmd := exec.Command(b.binary, args...)
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("nft %s: %s: %w",
            strings.Join(args, " "), stderr.String(), err)
    }
    return stdout.Bytes(), nil
}

// JSON output for parsing
func (b *NftablesBackend) listJSON() ([]byte, error) {
    return b.exec("-j", "list", "table", "inet", b.tableName)
}
```

### 7.2 Atomic Apply via Ruleset File

```go
func (b *NftablesBackend) Apply(rs *CompiledRuleSet) error {
    // Generate complete nftables ruleset file
    nftScript := b.generateScript(rs)
    
    // Write to temp file
    tmpFile, err := os.CreateTemp("", "rampart-*.nft")
    if err != nil {
        return fmt.Errorf("nftables.Apply: create temp: %w", err)
    }
    defer os.Remove(tmpFile.Name())
    
    if _, err := tmpFile.WriteString(nftScript); err != nil {
        return fmt.Errorf("nftables.Apply: write: %w", err)
    }
    tmpFile.Close()
    
    // Atomic apply
    if _, err := b.exec("-f", tmpFile.Name()); err != nil {
        return fmt.Errorf("nftables.Apply: apply: %w", err)
    }
    
    return nil
}
```

### 7.3 Rule Translation

```go
func (b *NftablesBackend) ruleToNft(rule CompiledRule) string {
    var parts []string
    
    // Direction → chain
    chain := b.directionToChain(rule.Direction)
    parts = append(parts, fmt.Sprintf("add rule inet %s %s", b.tableName, chain))
    
    // Source CIDR
    if len(rule.Match.SourceNets) > 0 {
        if len(rule.Match.SourceNets) == 1 {
            parts = append(parts, fmt.Sprintf("ip saddr %s", rule.Match.SourceNets[0].String()))
        } else {
            // Use anonymous set
            cidrs := make([]string, len(rule.Match.SourceNets))
            for i, n := range rule.Match.SourceNets {
                cidrs[i] = n.String()
            }
            parts = append(parts, fmt.Sprintf("ip saddr { %s }", strings.Join(cidrs, ", ")))
        }
    }
    
    // Protocol
    if len(rule.Match.Protocols) == 1 {
        parts = append(parts, rule.Match.Protocols[0].String())
    }
    
    // Destination ports
    if len(rule.Match.DestPorts) > 0 {
        ports := formatPortsNft(rule.Match.DestPorts)
        parts = append(parts, fmt.Sprintf("dport %s", ports))
    }
    
    // Connection state
    if len(rule.Match.States) > 0 {
        states := formatStatesNft(rule.Match.States)
        parts = append(parts, fmt.Sprintf("ct state %s", states))
    }
    
    // Counter + action
    parts = append(parts, "counter")
    parts = append(parts, rule.Action.ToNft())
    
    // Comment (for identification)
    parts = append(parts, fmt.Sprintf(`comment "rampart:%s"`, rule.Name))
    
    return strings.Join(parts, " ")
}
```

### 7.4 Current State Parsing

```go
func (b *NftablesBackend) CurrentState() (*CompiledRuleSet, error) {
    output, err := b.listJSON()
    if err != nil {
        return nil, err
    }
    
    var nftJSON struct {
        Nftables []map[string]interface{} `json:"nftables"`
    }
    if err := json.Unmarshal(output, &nftJSON); err != nil {
        return nil, fmt.Errorf("nftables: parse JSON: %w", err)
    }
    
    // Extract rules with "rampart:" comment prefix
    var rules []CompiledRule
    for _, item := range nftJSON.Nftables {
        if ruleData, ok := item["rule"]; ok {
            rule, err := b.parseNftRule(ruleData)
            if err != nil {
                continue // Skip non-Rampart rules
            }
            if rule != nil {
                rules = append(rules, *rule)
            }
        }
    }
    
    return &CompiledRuleSet{Rules: rules}, nil
}
```

---

## 8. iptables Backend Implementation

### 8.1 Chain Swap Strategy

```go
func (b *IptablesBackend) Apply(rs *CompiledRuleSet) error {
    // 1. Create new chains
    for _, chain := range []string{"INPUT", "FORWARD", "OUTPUT"} {
        newChain := fmt.Sprintf("%s-%s-NEW", b.chainPrefix, chain)
        b.exec("-N", newChain)
    }
    
    // 2. Populate new chains
    for _, rule := range rs.Rules {
        chain := b.directionToChain(rule.Direction) + "-NEW"
        args := b.ruleToArgs(rule, chain)
        if err := b.exec(args...); err != nil {
            // Cleanup: delete new chains
            b.cleanupNewChains()
            return fmt.Errorf("iptables.Apply: add rule %s: %w", rule.Name, err)
        }
    }
    
    // 3. Swap: update jump targets
    for _, chain := range []string{"INPUT", "FORWARD", "OUTPUT"} {
        oldChain := fmt.Sprintf("%s-%s", b.chainPrefix, chain)
        newChain := fmt.Sprintf("%s-%s-NEW", b.chainPrefix, chain)
        
        // Add jump to new chain
        b.exec("-I", chain, "1", "-j", newChain)
        // Remove jump to old chain
        b.exec("-D", chain, "-j", oldChain)
        // Flush and delete old chain
        b.exec("-F", oldChain)
        b.exec("-X", oldChain)
        // Rename new chain
        b.exec("-E", newChain, oldChain)
    }
    
    return nil
}
```

### 8.2 iptables-save Parser

```go
func parseIptablesSave(data []byte) ([]CompiledRule, error) {
    var rules []CompiledRule
    scanner := bufio.NewScanner(bytes.NewReader(data))
    
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        
        // Skip comments and table headers
        if line == "" || line[0] == '#' || line[0] == '*' || line == "COMMIT" {
            continue
        }
        
        // Skip chain definitions
        if line[0] == ':' {
            continue
        }
        
        // Parse -A CHAIN rules
        if strings.HasPrefix(line, "-A ") {
            rule, err := parseIptablesRule(line)
            if err != nil {
                continue
            }
            // Only include Rampart-managed rules (by comment)
            if strings.HasPrefix(rule.Comment, "rampart:") {
                rules = append(rules, *rule)
            }
        }
    }
    
    return rules, nil
}
```

---

## 9. Snapshot Engine

### 9.1 Snapshot Creation

```go
type SnapshotStore struct {
    dir       string
    retention RetentionConfig
}

func (s *SnapshotStore) Create(trigger string, description string, state *CompiledRuleSet, backendState []byte) (*Snapshot, error) {
    snap := &Snapshot{
        ID:          generateUUIDv7(),
        CreatedAt:   time.Now(),
        CreatedBy:   currentUser(),
        Trigger:     trigger,
        Description: description,
        PolicyHash:  state.Hash,
        RuleCount:   len(state.Rules),
        Backend:     state.Backend,
    }
    
    // Serialize: gob encode + zstd compress
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    if err := enc.Encode(snapData{Rules: state.Rules, BackendState: backendState}); err != nil {
        return nil, fmt.Errorf("snapshot: encode: %w", err)
    }
    
    compressed := zstdCompress(buf.Bytes())
    
    // Write to file
    filename := fmt.Sprintf("%s-%s.snap", snap.ID, trigger)
    path := filepath.Join(s.dir, filename)
    if err := os.WriteFile(path, compressed, 0600); err != nil {
        return nil, fmt.Errorf("snapshot: write: %w", err)
    }
    
    snap.Size = int64(len(compressed))
    
    // Write metadata index
    s.appendIndex(snap)
    
    return snap, nil
}
```

### 9.2 zstd Compression (from scratch)

Simplified zstd implementation for snapshot compression:

```go
// Minimal zstd compressor using entropy coding + LZ77
// For snapshots (small data, ~KB), simple implementation is sufficient
func zstdCompress(data []byte) []byte {
    // Frame header
    var out bytes.Buffer
    out.Write([]byte{0x28, 0xb5, 0x2f, 0xfd}) // Magic number
    
    // For small data (<128KB): use raw block (no compression)
    // For larger data: implement basic LZ77 + FSE
    if len(data) < 128*1024 {
        writeRawBlock(&out, data)
    } else {
        writeCompressedBlock(&out, data)
    }
    
    return out.Bytes()
}

func zstdDecompress(data []byte) ([]byte, error) {
    // Verify magic number
    if !bytes.HasPrefix(data, []byte{0x28, 0xb5, 0x2f, 0xfd}) {
        return nil, fmt.Errorf("invalid zstd magic")
    }
    // Decode blocks
    return decodeBlocks(data[4:])
}
```

### 9.3 Retention Cleanup

```go
func (s *SnapshotStore) Cleanup() error {
    snaps, err := s.List()
    if err != nil {
        return err
    }
    
    // Sort by creation time (newest first)
    sort.Slice(snaps, func(i, j int) bool {
        return snaps[i].CreatedAt.After(snaps[j].CreatedAt)
    })
    
    cutoff := time.Now().Add(-s.retention.MaxAge)
    
    for i, snap := range snaps {
        // Keep if within max count AND within max age
        if i < s.retention.MaxCount && snap.CreatedAt.After(cutoff) {
            continue
        }
        // Delete snapshot file
        os.Remove(filepath.Join(s.dir, snap.Filename))
    }
    
    return nil
}
```

---

## 10. Audit System

### 10.1 Append-Only JSONL Store

```go
type AuditStore struct {
    dir       string
    retention time.Duration
    mu        sync.Mutex
    file      *os.File
    lastHash  string // Previous entry hash (for chain)
    eventC    chan AuditEvent
}

func NewAuditStore(dir string) *AuditStore {
    s := &AuditStore{
        dir:    dir,
        eventC: make(chan AuditEvent, 1000),
    }
    go s.writer() // Single writer goroutine
    return s
}

func (s *AuditStore) writer() {
    for event := range s.eventC {
        s.mu.Lock()
        s.writeEvent(event)
        s.mu.Unlock()
    }
}

func (s *AuditStore) writeEvent(event AuditEvent) error {
    // Ensure file for today
    today := time.Now().Format("2006-01-02")
    filename := fmt.Sprintf("audit-%s.jsonl", today)
    path := filepath.Join(s.dir, filename)
    
    if s.file == nil || !s.isCurrentFile(path) {
        s.rotateFile(path)
    }
    
    // Compute hash chain
    data, _ := json.Marshal(event)
    chainInput := s.lastHash + string(data)
    hash := sha256.Sum256([]byte(chainInput))
    event.ChainHash = hex.EncodeToString(hash[:])
    s.lastHash = event.ChainHash
    
    // Write JSON line
    line, _ := json.Marshal(event)
    s.file.Write(line)
    s.file.Write([]byte("\n"))
    
    return nil
}
```

### 10.2 Audit Query

```go
type AuditQuery struct {
    Action  string
    Actor   string
    Since   time.Time
    Until   time.Time
    Limit   int
    Offset  int
}

func (s *AuditStore) Search(q AuditQuery) ([]AuditEvent, int, error) {
    var results []AuditEvent
    
    // Determine which files to scan
    files := s.filesInRange(q.Since, q.Until)
    
    for _, f := range files {
        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            var event AuditEvent
            if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
                continue
            }
            
            if matchesQuery(event, q) {
                results = append(results, event)
            }
        }
    }
    
    // Sort by timestamp (newest first)
    sort.Slice(results, func(i, j int) bool {
        return results[i].Timestamp.After(results[j].Timestamp)
    })
    
    total := len(results)
    
    // Apply offset/limit
    if q.Offset < len(results) {
        results = results[q.Offset:]
    }
    if q.Limit > 0 && len(results) > q.Limit {
        results = results[:q.Limit]
    }
    
    return results, total, nil
}
```

---

## 11. Packet Simulator

### 11.1 Simulation Engine

```go
func (s *Simulator) Simulate(pkt SimulatedPacket) SimulationResult {
    start := time.Now()
    
    for i, rule := range s.rules {
        // Skip rules for wrong direction
        if rule.Direction != pkt.Direction {
            continue
        }
        
        // Check schedule (is rule active now?)
        if rule.Schedule != nil && !rule.Schedule.IsActive(time.Now()) {
            continue
        }
        
        if matchesPacket(rule.Match, pkt) {
            return SimulationResult{
                Verdict:     rule.Action,
                MatchedRule: &s.rules[i],
                MatchPath:   buildMatchPath(rule, pkt),
                Evaluated:   i + 1,
                Duration:    time.Since(start),
            }
        }
    }
    
    // No rule matched → default policy
    return SimulationResult{
        Verdict:   ActionDrop, // Default deny
        Evaluated: len(s.rules),
        Duration:  time.Since(start),
        MatchPath: "no matching rule; default policy: drop",
    }
}

func matchesPacket(match CompiledMatch, pkt SimulatedPacket) bool {
    // Protocol check
    if len(match.Protocols) > 0 && !containsProtocol(match.Protocols, pkt.Protocol) {
        return false
    }
    
    // Source CIDR check
    if len(match.SourceNets) > 0 {
        if !anyContains(match.SourceNets, pkt.SourceIP) {
            return false
        }
    }
    
    // Dest CIDR check
    if len(match.DestNets) > 0 {
        if !anyContains(match.DestNets, pkt.DestIP) {
            return false
        }
    }
    
    // Dest port check
    if len(match.DestPorts) > 0 {
        if !portInRanges(pkt.DestPort, match.DestPorts) {
            return false
        }
    }
    
    // Source port check
    if len(match.SourcePorts) > 0 {
        if !portInRanges(pkt.SourcePort, match.SourcePorts) {
            return false
        }
    }
    
    // Interface check
    if len(match.Interfaces) > 0 && pkt.Interface != "" {
        if !contains(match.Interfaces, pkt.Interface) {
            return false
        }
    }
    
    // Connection state check
    if len(match.States) > 0 && pkt.State != "" {
        if !containsState(match.States, pkt.State) {
            return false
        }
    }
    
    // Negation check
    if match.Negated != nil {
        if matchesPacketMatch(*match.Negated, pkt) {
            return false // Negated match hit → overall miss
        }
    }
    
    return true
}
```

---

## 12. Time-Based Rule Scheduler

### 12.1 Scheduler Loop

```go
type Scheduler struct {
    engine   *Engine
    interval time.Duration
    stopC    chan struct{}
}

func (s *Scheduler) Run() {
    ticker := time.NewTicker(s.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            s.evaluate()
        case <-s.stopC:
            return
        }
    }
}

func (s *Scheduler) evaluate() {
    now := time.Now()
    
    changed := false
    rules := s.engine.CurrentRules()
    
    for _, rule := range rules {
        if rule.Schedule == nil {
            continue
        }
        
        wasActive := rule.Schedule.wasActive
        isActive := rule.Schedule.IsActive(now)
        
        if wasActive && !isActive {
            // Rule expired → remove from backend
            s.engine.DeactivateRule(rule.ID)
            s.engine.Audit(AuditEvent{
                Action:   "rule.expired",
                Resource: AuditResource{Type: "rule", ID: rule.ID, Name: rule.Name},
            })
            changed = true
        } else if !wasActive && isActive {
            // Rule became active → add to backend
            s.engine.ActivateRule(rule.ID)
            s.engine.Audit(AuditEvent{
                Action:   "rule.activated",
                Resource: AuditResource{Type: "rule", ID: rule.ID, Name: rule.Name},
            })
            changed = true
        }
        
        rule.Schedule.wasActive = isActive
    }
    
    if changed {
        s.engine.ReapplyRules()
    }
}
```

### 12.2 Schedule Evaluation

```go
func (sched *Schedule) IsActive(now time.Time) bool {
    // One-time schedule
    if sched.ActiveFrom != nil && now.Before(*sched.ActiveFrom) {
        return false
    }
    if sched.ActiveUntil != nil && now.After(*sched.ActiveUntil) {
        return false
    }
    
    // Recurring schedule
    if sched.Recurring != nil {
        loc, _ := time.LoadLocation(sched.Recurring.Timezone)
        localNow := now.In(loc)
        
        // Check day of week
        if len(sched.Recurring.Days) > 0 {
            dayMatch := false
            for _, d := range sched.Recurring.Days {
                if localNow.Weekday() == d {
                    dayMatch = true
                    break
                }
            }
            if !dayMatch {
                return false
            }
        }
        
        // Check time of day
        currentTime := localNow.Format("15:04")
        if currentTime < sched.Recurring.StartTime || currentTime >= sched.Recurring.EndTime {
            return false
        }
    }
    
    return true
}
```

---

## 13. REST API Server

### 13.1 Custom Router

```go
type Router struct {
    trees       map[string]*node // method → trie root
    middlewares []Middleware
}

type node struct {
    path     string
    handler  http.Handler
    children []*node
    param    string // ":id" → "id"
    isParam  bool
}

func (r *Router) Handle(method, pattern string, handler http.HandlerFunc) {
    // Build with middleware chain
    h := http.Handler(handler)
    for i := len(r.middlewares) - 1; i >= 0; i-- {
        h = r.middlewares[i](h)
    }
    r.addRoute(method, pattern, h)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    tree, ok := r.trees[req.Method]
    if !ok {
        http.Error(w, "Method Not Allowed", 405)
        return
    }
    
    handler, params := tree.find(req.URL.Path)
    if handler == nil {
        http.Error(w, "Not Found", 404)
        return
    }
    
    // Store params in context
    ctx := context.WithValue(req.Context(), paramsKey, params)
    handler.ServeHTTP(w, req.WithContext(ctx))
}
```

### 13.2 Server Setup

```go
func NewAPIServer(cfg config.ServerConfig, engine *Engine) *Server {
    router := NewRouter()
    
    // Middleware
    router.Use(requestIDMiddleware)
    router.Use(loggingMiddleware)
    router.Use(corsMiddleware(cfg.CORS))
    router.Use(authMiddleware(cfg.API.Keys))
    
    s := &Server{engine: engine, router: router}
    
    // Policy endpoints
    router.Handle("POST", "/api/v1/policies/plan", s.handlePlan)
    router.Handle("POST", "/api/v1/policies/apply", s.handleApply)
    router.Handle("POST", "/api/v1/policies/simulate", s.handleSimulate)
    router.Handle("GET", "/api/v1/policies/current", s.handleGetCurrent)
    
    // Rules endpoints
    router.Handle("GET", "/api/v1/rules", s.handleListRules)
    router.Handle("POST", "/api/v1/rules", s.handleAddRule)
    router.Handle("GET", "/api/v1/rules/:id", s.handleGetRule)
    router.Handle("DELETE", "/api/v1/rules/:id", s.handleDeleteRule)
    router.Handle("GET", "/api/v1/rules/:id/stats", s.handleRuleStats)
    
    // Snapshot endpoints
    router.Handle("GET", "/api/v1/snapshots", s.handleListSnapshots)
    router.Handle("POST", "/api/v1/snapshots", s.handleCreateSnapshot)
    router.Handle("POST", "/api/v1/snapshots/:id/rollback", s.handleRollback)
    router.Handle("GET", "/api/v1/snapshots/:id/diff", s.handleSnapshotDiff)
    
    // Audit endpoints
    router.Handle("GET", "/api/v1/audit", s.handleListAudit)
    router.Handle("GET", "/api/v1/audit/:id", s.handleGetAudit)
    
    // Cluster endpoints
    router.Handle("GET", "/api/v1/cluster/status", s.handleClusterStatus)
    router.Handle("GET", "/api/v1/cluster/health", s.handleHealth)
    
    // System
    router.Handle("GET", "/api/v1/system/info", s.handleSystemInfo)
    router.Handle("GET", "/api/v1/system/health", s.handleHealth)
    router.Handle("GET", "/metrics", s.handleMetrics)
    
    // SSE
    router.Handle("GET", "/api/v1/events", s.handleSSE)
    
    // WebUI (embedded React)
    router.Handle("GET", "/ui/*", s.serveUI())
    
    return s
}
```

### 13.3 SSE (Server-Sent Events)

```go
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "SSE not supported", 500)
        return
    }
    
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    
    // Subscribe to events
    ch := s.engine.Subscribe()
    defer s.engine.Unsubscribe(ch)
    
    for {
        select {
        case event := <-ch:
            data, _ := json.Marshal(event)
            fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, data)
            flusher.Flush()
        case <-r.Context().Done():
            return
        }
    }
}
```

---

## 14. CLI Implementation

### 14.1 Flag Parsing Pattern

```go
func runApply(args []string) {
    fs := flag.NewFlagSet("apply", flag.ExitOnError)
    
    file := fs.String("f", "", "Policy YAML file path")
    server := fs.String("server", "", "Rampart server URL (for remote apply)")
    autoApprove := fs.Bool("auto-approve", false, "Skip confirmation prompt")
    dryRun := fs.Bool("dry-run", false, "Show plan without applying")
    output := fs.String("o", "text", "Output format: text, json")
    
    fs.Parse(args)
    
    if *file == "" {
        fmt.Fprintln(os.Stderr, "error: -f flag is required")
        os.Exit(1)
    }
    
    // Parse and compile
    ps, err := ParsePolicyFile(*file)
    exitOnError(err, "parse")
    
    compiled, err := Compile(ps, loadVars())
    exitOnError(err, "compile")
    
    // Detect conflicts
    conflicts := DetectConflicts(compiled.Rules)
    if hasErrors(conflicts) {
        printConflicts(conflicts)
        os.Exit(1)
    }
    
    // Get current state
    backend, err := selectBackend(loadConfig().Backend)
    exitOnError(err, "backend")
    
    current, err := backend.CurrentState()
    exitOnError(err, "current state")
    
    // Generate plan
    plan := GeneratePlan(current, compiled)
    
    // Show plan
    printPlan(plan, *output)
    
    if *dryRun {
        return
    }
    
    // Confirm
    if !*autoApprove {
        if !confirm("Apply these changes?") {
            fmt.Println("Cancelled.")
            return
        }
    }
    
    // Create pre-apply snapshot
    snapStore := NewSnapshotStore(loadConfig().Snapshots.Directory)
    snapStore.Create("pre-apply", "Pre-apply snapshot", current, nil)
    
    // Apply
    err = backend.Apply(compiled)
    exitOnError(err, "apply")
    
    // Audit
    auditStore := NewAuditStore(loadConfig().Audit.Directory)
    auditStore.Record(AuditEvent{
        Action: AuditApply,
        Actor:  AuditActor{Type: "user", Identity: currentUser()},
    })
    
    fmt.Printf("✓ Applied %d rules (%d added, %d removed, %d modified)\n",
        len(compiled.Rules), plan.AddCount, plan.RemoveCount, plan.ModifyCount)
}
```

### 14.2 Colorized Output

```go
// ANSI color codes (no external dependency)
const (
    colorReset  = "\033[0m"
    colorRed    = "\033[31m"
    colorGreen  = "\033[32m"
    colorYellow = "\033[33m"
    colorBlue   = "\033[34m"
    colorBold   = "\033[1m"
    colorDim    = "\033[2m"
)

func printPlan(plan *ExecutionPlan, format string) {
    if format == "json" {
        data, _ := json.MarshalIndent(plan, "", "  ")
        fmt.Println(string(data))
        return
    }
    
    fmt.Printf("\n%sRampart Policy Plan%s\n", colorBold, colorReset)
    fmt.Println(strings.Repeat("=", 40))
    
    for _, rule := range plan.ToAdd {
        fmt.Printf("  %s+ [P%d] %-25s %s%s\n",
            colorGreen, rule.Priority, rule.Name,
            formatRuleSummary(rule), colorReset)
    }
    
    for _, mod := range plan.ToModify {
        fmt.Printf("  %s~ [P%d] %-25s %s%s\n",
            colorYellow, mod.After.Priority, mod.After.Name,
            formatRuleSummary(mod.After), colorReset)
    }
    
    for _, rule := range plan.ToRemove {
        fmt.Printf("  %s- [P%d] %-25s %s%s\n",
            colorRed, rule.Priority, rule.Name,
            formatRuleSummary(rule), colorReset)
    }
    
    fmt.Printf("\nPlan: %d to add, %d to remove, %d to modify.\n",
        plan.AddCount, plan.RemoveCount, plan.ModifyCount)
}
```

---

## 15. React WebUI Architecture

### 15.1 Component Hierarchy

```
App
├── Layout (sidebar + topbar)
│   ├── Sidebar (navigation)
│   └── TopBar (cluster status badge, user menu)
├── Pages
│   ├── Dashboard
│   │   ├── StatsCards (rules count, cluster nodes, recent events)
│   │   ├── RuleHitHeatmap (D3.js)
│   │   └── RecentAuditFeed (SSE-powered)
│   ├── Policies
│   │   ├── PolicyEditor (CodeMirror 6, YAML mode)
│   │   ├── ConflictPanel (warnings/errors sidebar)
│   │   └── PlanPreview (diff view)
│   ├── Rules
│   │   ├── RuleTable (sortable, filterable)
│   │   ├── RuleDetail (modal/slide-over)
│   │   └── QuickAddRule (form)
│   ├── Simulator
│   │   ├── PacketForm (src, dst, port, protocol)
│   │   ├── TraceVisualization (rule evaluation path)
│   │   └── ResultCard (verdict + matched rule)
│   ├── Snapshots
│   │   ├── SnapshotTimeline (visual timeline)
│   │   ├── DiffViewer (side-by-side)
│   │   └── RollbackDialog (confirmation)
│   ├── AuditLog
│   │   ├── AuditTable (paginated, filterable)
│   │   ├── AuditDetail (before/after diff)
│   │   └── AuditSearch (filters bar)
│   ├── Cluster
│   │   ├── NodeList (status, role, backend)
│   │   ├── RaftStatus (term, commit index)
│   │   └── NodeDetail (per-node stats)
│   └── Settings
│       ├── BackendConfig
│       ├── APIKeys
│       ├── SnapshotRetention
│       └── WebhookConfig
└── Hooks
    ├── useSSE (real-time event stream)
    ├── useAPI (fetch wrapper with auth)
    └── useTheme (dark/light mode)
```

### 15.2 SSE Hook for Real-Time Updates

```typescript
function useSSE(url: string) {
  const [events, setEvents] = useState<SSEEvent[]>([]);
  
  useEffect(() => {
    const source = new EventSource(url);
    
    source.addEventListener('rule.applied', (e) => {
      setEvents(prev => [JSON.parse(e.data), ...prev].slice(0, 100));
    });
    
    source.addEventListener('audit.event', (e) => {
      // Trigger re-fetch of audit list
    });
    
    source.addEventListener('cluster.change', (e) => {
      // Update cluster status
    });
    
    return () => source.close();
  }, [url]);
  
  return events;
}
```

---

## 16. Configuration Loading

```go
func LoadConfig(paths ...string) (*Config, error) {
    cfg := defaultConfig()
    
    // Search order:
    // 1. Explicit path (--config flag)
    // 2. ./rampart.yaml
    // 3. /etc/rampart/rampart.yaml
    // 4. ~/.config/rampart/rampart.yaml
    
    searchPaths := paths
    if len(searchPaths) == 0 {
        searchPaths = []string{
            "rampart.yaml",
            "/etc/rampart/rampart.yaml",
            filepath.Join(homeDir(), ".config", "rampart", "rampart.yaml"),
        }
    }
    
    for _, path := range searchPaths {
        data, err := os.ReadFile(path)
        if err != nil {
            continue
        }
        if err := yaml.Unmarshal(data, cfg); err != nil {
            return nil, fmt.Errorf("config: parse %s: %w", path, err)
        }
        cfg.loadedFrom = path
        break
    }
    
    // Environment variable overrides
    applyEnvOverrides(cfg) // RAMPART_LISTEN, RAMPART_BACKEND, etc.
    
    return cfg, nil
}
```

---

## 17. Error Handling Strategy

### 17.1 Error Wrapping Convention

All errors include component and operation context:

```go
// Pattern: component.Operation: detail: underlying error
return fmt.Errorf("nftables.Apply: atomic replace: %w", err)
return fmt.Errorf("engine.Compile: policy %q, rule %q: %w", policyName, ruleName, err)
return fmt.Errorf("snapshot.Create: write %s: %w", path, err)
```

### 17.2 User-Facing Errors

```go
type UserError struct {
    Code    string // Machine-readable code
    Message string // Human-readable message
    Details interface{}
}

func (e *UserError) Error() string { return e.Message }

// Usage in API handlers
func (s *Server) handleApply(w http.ResponseWriter, r *http.Request) {
    // ...
    if conflicts := DetectConflicts(compiled.Rules); hasErrors(conflicts) {
        respondError(w, &UserError{
            Code:    "CONFLICT_DETECTED",
            Message: fmt.Sprintf("%d rule conflicts detected", len(conflicts)),
            Details: conflicts,
        }, http.StatusConflict)
        return
    }
}
```

---

## 18. Testing Strategy

### 18.1 Unit Tests

- Policy parser: test valid/invalid YAML, edge cases
- Rule compiler: test CIDR normalization, port parsing, variable substitution
- Conflict detector: test all conflict types with synthetic rules
- Simulator: test packet matching against known rulesets
- Scheduler: test time-based rule activation/deactivation (mock clock)
- Backend compilers: test rule → nft/iptables command translation

### 18.2 Integration Tests

- nftables backend: requires root + nftables installed (CI with privileged container)
- iptables backend: requires root + iptables installed
- API server: HTTP client tests against running server
- Full pipeline: YAML → compile → plan → apply → verify → rollback

### 18.3 Test Fixtures

```
testdata/
├── policies/
│   ├── valid-basic.yaml
│   ├── valid-complex.yaml
│   ├── valid-with-vars.yaml
│   ├── valid-with-includes.yaml
│   ├── valid-with-schedule.yaml
│   ├── invalid-schema.yaml
│   ├── invalid-cidr.yaml
│   ├── invalid-port.yaml
│   └── conflicts/
│       ├── shadow.yaml
│       ├── contradiction.yaml
│       └── redundancy.yaml
├── snapshots/
│   └── sample.snap
└── nftables/
    ├── sample-output.json
    └── expected-script.nft
```

### 18.4 Mock Backend for Testing

```go
type MockBackend struct {
    rules   []CompiledRule
    applied int
    probeOK bool
}

func (m *MockBackend) Name() string           { return "mock" }
func (m *MockBackend) Probe() error           { if !m.probeOK { return errors.New("mock unavailable") }; return nil }
func (m *MockBackend) CurrentState() (*CompiledRuleSet, error) { return &CompiledRuleSet{Rules: m.rules}, nil }
func (m *MockBackend) Apply(rs *CompiledRuleSet) error { m.rules = rs.Rules; m.applied++; return nil }
func (m *MockBackend) DryRun(rs *CompiledRuleSet) (*ExecutionPlan, error) { return GeneratePlan(&CompiledRuleSet{Rules: m.rules}, rs), nil }
// ...
```

