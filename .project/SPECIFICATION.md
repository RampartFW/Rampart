# Rampart — Technical Specification

> **Version:** 1.0.0-draft  
> **Status:** Design Phase  
> **Language:** Go 1.23+  
> **Dependencies:** Zero (stdlib + golang.org/x/crypto, golang.org/x/sys, gopkg.in/yaml.v3 only)  
> **Binary:** Single static binary (server + agent + CLI unified)  
> **License:** Apache 2.0  
> **Domain:** rampartfw.com  
> **Repository:** github.com/rampartfw/rampart

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Problem Statement](#2-problem-statement)
3. [Architecture Overview](#3-architecture-overview)
4. [Policy Model](#4-policy-model)
5. [Policy-as-Code (YAML Schema)](#5-policy-as-code-yaml-schema)
6. [Backend Abstraction Layer](#6-backend-abstraction-layer)
7. [nftables Backend](#7-nftables-backend)
8. [iptables Backend](#8-iptables-backend)
9. [eBPF/XDP Backend](#9-ebpfxdp-backend)
10. [Cloud Security Group Backends](#10-cloud-security-group-backends)
11. [Rule Compiler & Conflict Detection](#11-rule-compiler--conflict-detection)
12. [Dry-Run & Policy Simulation](#12-dry-run--policy-simulation)
13. [Snapshot & Rollback Engine](#13-snapshot--rollback-engine)
14. [Audit System](#14-audit-system)
15. [Raft Cluster & Multi-Host Sync](#15-raft-cluster--multi-host-sync)
16. [REST API](#16-rest-api)
17. [CLI Design](#17-cli-design)
18. [React WebUI](#18-react-webui)
19. [MCP Server](#19-mcp-server)
20. [Time-Based Rules & Scheduling](#20-time-based-rules--scheduling)
21. [Configuration](#21-configuration)
22. [Security Model](#22-security-model)
23. [Observability](#23-observability)
24. [Performance Targets](#24-performance-targets)
25. [Project Structure](#25-project-structure)
26. [Version Roadmap](#26-version-roadmap)

---

## 1. Executive Summary

Rampart is a **network policy engine** — a unified firewall rule manager that abstracts away the complexity of iptables, nftables, eBPF, and cloud security groups behind a single policy-as-code interface. Written entirely in Go with zero external dependencies, it deploys as a single static binary.

### Core Value Propositions

- **Pluggable backends:** nftables, iptables, eBPF/XDP, AWS Security Groups, GCP Firewall Rules, Azure NSGs — all managed through the same YAML policy files.
- **Policy-as-code:** Firewall rules defined in version-controllable YAML. No more manual `iptables -A` commands.
- **Dry-run mode:** See exactly what will change before applying (Terraform plan for firewalls).
- **Instant rollback:** Snapshot-based state management. One command to revert to any previous state.
- **Audit trail:** Every change recorded — who, when, what, before/after diff.
- **Multi-host sync:** Raft consensus ensures all nodes in a cluster converge to the same policy state.
- **Policy simulation:** Test if a packet would be allowed/denied without applying any rules.
- **Time-based rules:** Temporary rules with auto-expiry (maintenance windows, incident response).
- **Single binary:** Server, agent, CLI — all in one Go binary. `rampart serve`, `rampart agent`, `rampart apply`.

### What Rampart Replaces

| Tool | Problem | Rampart Solution |
|------|---------|------------------|
| Raw iptables/nftables | Manual, error-prone, no audit | Policy-as-code + audit + rollback |
| UFW / firewalld | Single-host only, no sync | Raft-based multi-host sync |
| Terraform (firewall) | Cloud-only, no Linux host support | Unified local + cloud backends |
| Ansible firewall roles | Slow convergence, push-based | Real-time Raft sync, pull-based |
| Custom scripts | No conflict detection, no dry-run | Rule compiler + simulator |

---

## 2. Problem Statement

### Current Pain Points

1. **Fragmentation:** Linux firewalls (iptables vs nftables), cloud firewalls (AWS SG, GCP FW, Azure NSG), container firewalls (Calico, Cilium) — all have different syntaxes and semantics.

2. **No version control:** `iptables -A INPUT -p tcp --dport 22 -j ACCEPT` is executed and forgotten. No git history, no blame, no rollback.

3. **No audit trail:** Who opened port 3306 to the world? When? Why? Nobody knows until the breach.

4. **Multi-host drift:** 50 servers should have identical firewall rules. After 6 months, they don't. Configuration drift is invisible.

5. **No dry-run:** One wrong iptables rule locks you out of SSH. There's no `--dry-run` flag.

6. **No conflict detection:** Rule A allows port 80 from 0.0.0.0/0. Rule B blocks port 80 from 10.0.0.0/8. Which wins? Depends on insertion order.

7. **No time-based rules:** "Open port 8080 for 2 hours during maintenance" requires manual rule + manual removal + hoping you remember.

### Target Users

- **DevOps engineers** managing fleet firewall rules
- **SREs** doing incident response (quick block/unblock with audit)
- **Security teams** enforcing network policies across hybrid infrastructure
- **Platform teams** providing firewall-as-a-service to developers

---

## 3. Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Rampart Architecture                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ CLI      │  │ REST API │  │ React UI │  │ MCP Srvr │    │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘    │
│       │              │              │              │          │
│       └──────────────┴──────┬───────┴──────────────┘          │
│                             │                                 │
│                    ┌────────▼────────┐                        │
│                    │  Policy Engine  │                        │
│                    │  ┌────────────┐ │                        │
│                    │  │ YAML Parser│ │                        │
│                    │  ├────────────┤ │                        │
│                    │  │ Rule       │ │                        │
│                    │  │ Compiler   │ │                        │
│                    │  ├────────────┤ │                        │
│                    │  │ Conflict   │ │                        │
│                    │  │ Detector   │ │                        │
│                    │  ├────────────┤ │                        │
│                    │  │ Simulator  │ │                        │
│                    │  └────────────┘ │                        │
│                    └────────┬────────┘                        │
│                             │                                 │
│              ┌──────────────┼──────────────┐                  │
│              │              │              │                   │
│     ┌────────▼──────┐ ┌────▼─────┐ ┌──────▼───────┐         │
│     │ Snapshot &    │ │ Audit    │ │ Raft Cluster │         │
│     │ Rollback Eng. │ │ System   │ │ (Multi-Host) │         │
│     └────────┬──────┘ └──────────┘ └──────┬───────┘         │
│              │                             │                  │
│              └──────────────┬──────────────┘                  │
│                             │                                 │
│                    ┌────────▼────────┐                        │
│                    │ Backend         │                        │
│                    │ Abstraction     │                        │
│                    │ Layer (BAL)     │                        │
│                    └────────┬────────┘                        │
│                             │                                 │
│       ┌──────────┬──────────┼──────────┬──────────┐          │
│       │          │          │          │          │           │
│  ┌────▼───┐ ┌───▼────┐ ┌───▼───┐ ┌───▼───┐ ┌───▼───┐      │
│  │nftables│ │iptables│ │eBPF/  │ │AWS SG │ │GCP/   │      │
│  │Backend │ │Backend │ │XDP    │ │Backend│ │Azure  │      │
│  └────────┘ └────────┘ └───────┘ └───────┘ └───────┘      │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility |
|-----------|---------------|
| **Policy Engine** | Parse YAML → compile rules → detect conflicts → build execution plan |
| **Rule Compiler** | Transform abstract policy into backend-specific rule sets |
| **Conflict Detector** | Identify overlapping, shadowed, or contradictory rules |
| **Simulator** | Test packet flows against compiled ruleset without applying |
| **Snapshot Engine** | Capture current firewall state, store, compare, restore |
| **Audit System** | Append-only log of every change with full context |
| **Raft Cluster** | Strong consistency multi-host policy distribution |
| **Backend Abstraction Layer** | Pluggable interface for different firewall implementations |
| **REST API** | HTTP/JSON API for programmatic access |
| **React WebUI** | Dashboard for rule management, audit, cluster status |
| **MCP Server** | AI agent integration for firewall management |
| **CLI** | Command-line interface for all operations |

---

## 4. Policy Model

### 4.1 Core Concepts

```
Policy Set (cluster-wide)
  └── Policy (named group of rules)
       └── Rule (single firewall rule)
            ├── Match (conditions: src, dst, port, protocol, ...)
            ├── Action (accept, drop, reject, log, rate-limit)
            ├── Priority (numeric, lower = higher priority)
            ├── Schedule (optional: time-based activation)
            └── Metadata (description, tags, owner, created_at)
```

### 4.2 Policy Inheritance

Policies can be layered with explicit priority ordering:

1. **System policies** (priority 0-99): SSH access, management ports — cannot be overridden
2. **Organization policies** (priority 100-499): Company-wide security baselines
3. **Team policies** (priority 500-799): Team-specific rules
4. **Service policies** (priority 800-999): Application-specific rules

Higher-priority (lower number) rules always win. Conflicts within the same priority level are flagged.

### 4.3 Internal Rule Representation

```go
type Rule struct {
    ID          string            // UUID v7 (time-sortable)
    Name        string            // Human-readable name
    PolicyID    string            // Parent policy reference
    Priority    int               // 0-999
    Direction   Direction         // Inbound | Outbound | Forward
    Action      Action            // Accept | Drop | Reject | Log | RateLimit
    Match       Match             // Matching conditions
    Schedule    *Schedule         // Optional time-based activation
    Tags        map[string]string // Arbitrary key-value metadata
    Description string            // Human-readable description
    CreatedAt   time.Time
    CreatedBy   string
    Version     uint64            // Monotonically increasing
}

type Match struct {
    SourceCIDRs      []string // Source IP/CIDR list
    DestCIDRs        []string // Destination IP/CIDR list
    SourcePorts      []PortRange
    DestPorts        []PortRange
    Protocols        []Protocol // TCP, UDP, ICMP, ICMPv6, Any
    Interfaces       []string   // Network interfaces (eth0, wg0, etc.)
    States           []ConnState // New, Established, Related, Invalid
    Not              *Match     // Negation (match everything except this)
    ICMPTypes        []uint8    // ICMP type filter
    IPVersion        IPVersion  // IPv4, IPv6, Both
}

type PortRange struct {
    Start uint16
    End   uint16 // Same as Start for single port
}

type Schedule struct {
    ActiveFrom  *time.Time     // Start time (nil = immediately)
    ActiveUntil *time.Time     // End time (nil = permanent)
    Recurring   *RecurringSpec // Cron-like recurring schedule
}

type RecurringSpec struct {
    Days      []time.Weekday // Which days
    StartTime string         // "09:00" (local time)
    EndTime   string         // "17:00" (local time)
    Timezone  string         // IANA timezone
}
```

---

## 5. Policy-as-Code (YAML Schema)

### 5.1 Policy File Structure

```yaml
# rampart-policy.yaml
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: production-web-tier
  description: "Firewall rules for production web servers"
  owner: platform-team
  tags:
    environment: production
    tier: web

# Global defaults for all rules in this policy set
defaults:
  direction: inbound
  action: drop  # Default deny
  ipVersion: both
  states: [established, related]  # Stateful by default

policies:
  - name: ssh-access
    priority: 10  # System-level
    description: "SSH access from bastion hosts only"
    rules:
      - name: allow-ssh-bastion
        match:
          protocol: tcp
          destPorts: [22]
          sourceCIDRs:
            - 10.0.1.0/24     # Bastion subnet
            - 172.16.0.5/32   # Jump host
        action: accept

      - name: deny-ssh-all
        match:
          protocol: tcp
          destPorts: [22]
        action: drop
        log: true  # Log denied SSH attempts

  - name: web-traffic
    priority: 500
    description: "Public HTTP/HTTPS access"
    rules:
      - name: allow-http
        match:
          protocol: tcp
          destPorts: [80, 443]
          sourceCIDRs: ["0.0.0.0/0", "::/0"]
        action: accept

      - name: allow-http3
        match:
          protocol: udp
          destPorts: [443]
          sourceCIDRs: ["0.0.0.0/0", "::/0"]
        action: accept

  - name: monitoring
    priority: 600
    description: "Monitoring and health check access"
    rules:
      - name: allow-prometheus
        match:
          protocol: tcp
          destPorts: [9090, 9100]
          sourceCIDRs: [10.0.10.0/24]
        action: accept

      - name: allow-healthcheck
        match:
          protocol: tcp
          destPorts: [8080]
          sourceCIDRs: [10.0.0.0/8]
        action: accept

  - name: maintenance-window
    priority: 800
    description: "Temporary debug access"
    rules:
      - name: temp-debug-port
        match:
          protocol: tcp
          destPorts: [9999]
          sourceCIDRs: [10.0.1.100/32]
        action: accept
        schedule:
          activeFrom: "2026-04-15T10:00:00Z"
          activeUntil: "2026-04-15T14:00:00Z"

  - name: rate-limiting
    priority: 200
    description: "Rate limiting for public services"
    rules:
      - name: syn-flood-protection
        match:
          protocol: tcp
          destPorts: [80, 443]
          states: [new]
        action: rate-limit
        rateLimit:
          rate: 100
          per: second
          burst: 200
          action: drop  # Action when limit exceeded

  - name: outbound
    priority: 300
    direction: outbound
    description: "Outbound access control"
    rules:
      - name: allow-dns
        match:
          protocol: [tcp, udp]
          destPorts: [53]
        action: accept

      - name: allow-http-out
        match:
          protocol: tcp
          destPorts: [80, 443]
        action: accept

      - name: allow-ntp
        match:
          protocol: udp
          destPorts: [123]
        action: accept

      - name: deny-all-outbound
        match: {}  # Match everything
        action: drop
        log: true
```

### 5.2 YAML Validation Rules

- `apiVersion` must be `rampartfw.com/v1`
- `kind` must be `PolicySet`
- Policy names must be unique within a PolicySet
- Rule names must be unique within a Policy
- Priority must be 0-999
- CIDR notation must be valid IPv4 or IPv6
- Port numbers must be 1-65535
- Port ranges: `start` must be ≤ `end`
- Schedule times must be valid RFC 3339
- `activeFrom` must be before `activeUntil` if both specified
- Rate limit `rate` must be > 0, `burst` must be ≥ `rate`

### 5.3 Variable Substitution

```yaml
# rampart-vars.yaml
apiVersion: rampartfw.com/v1
kind: Variables
metadata:
  name: production-vars

variables:
  bastion_subnet: "10.0.1.0/24"
  monitoring_subnet: "10.0.10.0/24"
  internal_network: "10.0.0.0/8"
  web_ports: [80, 443]
  ssh_port: 22
```

Referenced in policies:

```yaml
rules:
  - name: allow-ssh-bastion
    match:
      protocol: tcp
      destPorts: ["${ssh_port}"]
      sourceCIDRs: ["${bastion_subnet}"]
    action: accept
```

### 5.4 Policy Includes

```yaml
# rampart-policy.yaml
includes:
  - path: ./base-policies.yaml      # Relative path
  - path: /etc/rampart/org-base.yaml # Absolute path
  - url: https://policies.internal.company.com/security-baseline.yaml
```

---

## 6. Backend Abstraction Layer

### 6.1 Backend Interface

```go
// Backend is the core interface that all firewall backends must implement.
type Backend interface {
    // Name returns the backend identifier (e.g., "nftables", "iptables", "ebpf")
    Name() string

    // Capabilities reports what this backend supports
    Capabilities() BackendCapabilities

    // Probe checks if this backend is available on the current system
    Probe() error

    // CurrentState returns the active firewall rules in normalized form
    CurrentState() (*RuleSet, error)

    // Apply atomically applies a complete RuleSet, replacing all managed rules
    Apply(rs *RuleSet) error

    // DryRun returns what Apply would do without actually doing it
    DryRun(rs *RuleSet) (*ExecutionPlan, error)

    // Rollback restores a previously captured snapshot
    Rollback(snapshot *Snapshot) error

    // Flush removes all Rampart-managed rules (leaves system rules intact)
    Flush() error

    // Stats returns per-rule packet/byte counters
    Stats() (map[string]RuleStats, error)

    // Close releases any resources held by the backend
    Close() error
}

type BackendCapabilities struct {
    IPv4              bool
    IPv6              bool
    RateLimiting      bool
    ConnectionTracking bool
    Logging           bool
    NAT               bool
    PerRuleCounters   bool
    AtomicReplace     bool // Can replace entire ruleset atomically
    InterfaceFiltering bool
    MarkPackets       bool
    GeoIP             bool
}

type ExecutionPlan struct {
    Add    []CompiledRule // Rules to add
    Remove []CompiledRule // Rules to remove
    Modify []RuleChange   // Rules to modify (before/after)
    Order  []string       // Execution order (rule IDs)
}

type RuleChange struct {
    Before CompiledRule
    After  CompiledRule
    Diff   string // Human-readable diff
}

type RuleStats struct {
    RuleID  string
    Packets uint64
    Bytes   uint64
    LastHit time.Time
}
```

### 6.2 Backend Registry

```go
var backends = map[string]BackendFactory{}

type BackendFactory func(cfg BackendConfig) (Backend, error)

func Register(name string, factory BackendFactory) {
    backends[name] = factory
}

func NewBackend(name string, cfg BackendConfig) (Backend, error) {
    factory, ok := backends[name]
    if !ok {
        return nil, fmt.Errorf("unknown backend: %s", name)
    }
    return factory(cfg)
}

// Auto-detect best available backend
func AutoDetect() (Backend, error) {
    // Priority: nftables > iptables > eBPF
    for _, name := range []string{"nftables", "iptables", "ebpf"} {
        if factory, ok := backends[name]; ok {
            b, err := factory(DefaultConfig())
            if err == nil && b.Probe() == nil {
                return b, nil
            }
        }
    }
    return nil, fmt.Errorf("no supported firewall backend found")
}
```

---

## 7. nftables Backend

### 7.1 Overview

Primary backend for modern Linux systems (kernel ≥ 3.13, recommended ≥ 4.10).

### 7.2 Table & Chain Structure

Rampart manages its own nftables table to avoid conflicts with system rules:

```
table inet rampart {
    # Base chains (hooks into netfilter)
    chain input {
        type filter hook input priority 0; policy drop;
        # Rampart-managed rules here
    }
    chain forward {
        type filter hook forward priority 0; policy drop;
    }
    chain output {
        type filter hook output priority 0; policy accept;
    }
    
    # Named sets for IP lists
    set blocked_ips_v4 { type ipv4_addr; flags interval; }
    set blocked_ips_v6 { type ipv6_addr; flags interval; }
    set allowed_ips_v4 { type ipv4_addr; flags interval; }
    
    # Per-rule counters
    chain rule_ssh_bastion { counter accept; }
    chain rule_web_http { counter accept; }
}
```

### 7.3 Rule Translation

Abstract rule → nftables syntax:

```
Rule: allow-ssh-bastion (TCP, dst:22, src:10.0.1.0/24, ACCEPT)
  ↓
nft add rule inet rampart input ip saddr 10.0.1.0/24 tcp dport 22 counter accept comment "rampart:allow-ssh-bastion"
```

### 7.4 Atomic Replacement

nftables supports atomic ruleset replacement via `nft -f`:

```bash
# Generate complete ruleset file
# Apply atomically (all-or-nothing)
nft -f /tmp/rampart-ruleset.nft
```

This ensures no moment where the firewall is in a partial state.

### 7.5 Implementation Details

- **Execution:** Via `nft` CLI (parsed JSON output with `nft -j list ruleset`)
- **Monitoring:** nfnetlink socket for rule change notifications
- **Counters:** `nft list ruleset -j` includes counter values
- **Sets:** nftables named sets for efficient large IP list matching
- **Maps:** Verdict maps for port-to-action mappings

---

## 8. iptables Backend

### 8.1 Overview

Fallback backend for older Linux systems or environments where nftables is not available.

### 8.2 Chain Structure

```
# Rampart chains (do NOT modify system chains directly)
iptables -N RAMPART-INPUT
iptables -N RAMPART-FORWARD
iptables -N RAMPART-OUTPUT

# Jump from system chains to Rampart chains
iptables -A INPUT -j RAMPART-INPUT
iptables -A FORWARD -j RAMPART-FORWARD
iptables -A OUTPUT -j RAMPART-OUTPUT
```

### 8.3 Rule Translation

```
Rule: allow-ssh-bastion (TCP, dst:22, src:10.0.1.0/24, ACCEPT)
  ↓
iptables -A RAMPART-INPUT -p tcp --dport 22 -s 10.0.1.0/24 -m comment --comment "rampart:allow-ssh-bastion" -j ACCEPT
```

### 8.4 Atomic Replacement

iptables does NOT support atomic replacement. Rampart uses:

1. Create new chains: `RAMPART-INPUT-NEW`, `RAMPART-FORWARD-NEW`, `RAMPART-OUTPUT-NEW`
2. Populate new chains with all rules
3. Swap jump targets atomically (rename chains)
4. Delete old chains

### 8.5 Limitations vs nftables

| Feature | nftables | iptables |
|---------|----------|----------|
| Atomic replace | Native | Emulated (chain swap) |
| Named sets | Yes | ipset (separate) |
| Per-rule counters | Built-in | Built-in |
| Rate limiting | Yes | `-m limit` / `-m hashlimit` |
| IPv4+IPv6 unified | `inet` family | Separate `iptables` + `ip6tables` |
| JSON output | `nft -j` | `iptables-save` (text parsing) |

---

## 9. eBPF/XDP Backend

### 9.1 Overview

High-performance backend for XDP (eXpress Data Path) fast-path packet filtering. Operates at the network driver level — before the kernel networking stack.

### 9.2 Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│ Network     │ ──► │ XDP Program  │ ──► │ Kernel      │
│ Driver      │     │ (eBPF)       │     │ Network     │
│ (NIC)       │     │              │     │ Stack       │
└─────────────┘     │ XDP_DROP     │     └─────────────┘
                    │ XDP_PASS     │
                    │ XDP_TX       │
                    └──────────────┘
```

### 9.3 BPF Map Design

```
// Per-rule counters
BPF_MAP_TYPE_PERCPU_ARRAY: rule_stats[rule_index] → {packets, bytes}

// IP blocklist (LPM trie for CIDR matching)
BPF_MAP_TYPE_LPM_TRIE: blocked_cidrs[prefix] → action

// Port allowlist
BPF_MAP_TYPE_HASH: allowed_ports[port] → action

// Connection tracking (simplified)
BPF_MAP_TYPE_LRU_HASH: conntrack[5-tuple] → state
```

### 9.4 Compilation Strategy

Rampart ships pre-compiled eBPF bytecode (CO-RE — Compile Once, Run Everywhere) for common rule patterns. For custom rules, Rampart generates eBPF C source and compiles via the system's `clang` (if available).

### 9.5 Capabilities & Limitations

| Feature | Support |
|---------|---------|
| Basic L3/L4 filtering | Full |
| CIDR matching | Full (LPM trie) |
| Rate limiting | Full (token bucket in BPF) |
| Connection tracking | Limited (simplified 5-tuple) |
| Logging | Via perf events / ring buffer |
| NAT | Not supported |
| Application-layer | Not supported |

### 9.6 Hybrid Mode

eBPF can be used alongside nftables:
- **XDP fast path:** High-volume rules (DDoS mitigation, IP blocklists, rate limiting)
- **nftables slow path:** Complex rules (stateful, NAT, logging, application-aware)

```yaml
# rampart.yaml
backend:
  type: hybrid
  fastPath: ebpf   # XDP for high-volume, simple rules
  slowPath: nftables # nftables for complex rules
  fastPathRules:
    - "priority < 100"  # System-level rules on XDP
    - "action == rate-limit"
    - "tags.fastpath == true"
```

---

## 10. Cloud Security Group Backends

### 10.1 AWS Security Groups

```go
type AWSBackend struct {
    region     string
    sgID       string  // Security Group ID
    // Uses AWS SDK Go v2 (stdlib HTTP client, custom signer)
    // Zero dependency: implements SigV4 signing from scratch
}
```

**Mapping:**
- Rampart inbound rules → SG Ingress rules
- Rampart outbound rules → SG Egress rules
- CIDR matching → IpRanges / Ipv6Ranges
- Port ranges → FromPort/ToPort
- Protocol → IpProtocol

**Limitations:**
- Max 60 inbound + 60 outbound rules per SG (soft limit)
- No rate limiting
- No logging (use VPC Flow Logs separately)
- No connection state filtering (always stateful)

### 10.2 GCP Firewall Rules

```go
type GCPBackend struct {
    project  string
    network  string
    // Implements GCP REST API with OAuth2 (from scratch)
}
```

**Mapping:**
- Rampart rules → GCP Firewall Rules
- Priority → GCP priority (0-65535, reversed: lower = higher)
- Direction → INGRESS / EGRESS
- Tags for target filtering

### 10.3 Azure NSGs

```go
type AzureBackend struct {
    subscriptionID string
    resourceGroup  string
    nsgName        string
    // Implements Azure REST API with Azure AD OAuth2
}
```

**Mapping:**
- Rampart rules → NSG Security Rules
- Priority → Azure priority (100-4096)
- Direction → Inbound / Outbound

### 10.4 Cloud Backend Common Constraints

- API rate limiting → exponential backoff with jitter
- Eventual consistency → verify after apply (poll until converged)
- Cost → minimize API calls (batch where possible)
- No atomic replace → ordered add/remove with safety checks

---

## 11. Rule Compiler & Conflict Detection

### 11.1 Compilation Pipeline

```
YAML Policy Files
       │
       ▼
┌──────────────┐
│ YAML Parser  │ → Validate schema, resolve variables, process includes
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Rule Normalizer│ → Expand port ranges, normalize CIDRs, resolve DNS
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Priority     │ → Sort by priority, assign evaluation order
│ Sorter       │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Conflict     │ → Detect overlaps, shadows, contradictions
│ Detector     │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Schedule     │ → Evaluate time-based rules, filter active rules
│ Evaluator    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Backend      │ → Translate to backend-specific format
│ Compiler     │
└──────┬───────┘
       │
       ▼
CompiledRuleSet (ready for Apply or DryRun)
```

### 11.2 Conflict Detection Types

| Conflict Type | Description | Severity |
|--------------|-------------|----------|
| **Shadow** | Higher-priority rule makes lower-priority rule unreachable | Warning |
| **Contradiction** | Same priority, overlapping match, different actions | Error |
| **Redundancy** | Two rules with identical match and action | Warning |
| **Subset** | Rule A is a strict subset of Rule B (same action) | Info |
| **Overlap** | Partially overlapping CIDRs with different actions | Warning |

### 11.3 Conflict Detection Algorithm

For each pair of rules (R_i, R_j) where i < j (lower index = higher priority):

1. **Compute match intersection:** Do the source CIDRs overlap? Do the dest ports overlap? Same protocol?
2. **If intersection is non-empty:**
   - Same action → Redundancy (info/warning)
   - Different action + R_i fully covers R_j → Shadow (warning)
   - Different action + partial overlap → Overlap conflict (warning)
   - Same priority + different action → Contradiction (error, blocks apply)

CIDR overlap detection uses interval tree for efficient O(n log n) checking.

### 11.4 Conflict Report Format

```
$ rampart plan -f policy.yaml

Rampart Policy Plan
====================

⚠ 2 warnings, 0 errors

WARNING [shadow]: Rule "deny-ssh-all" (priority 10) is completely shadowed by 
  "allow-ssh-bastion" (priority 10) for source 10.0.1.0/24, port 22/tcp.
  → The deny rule will never match traffic from the bastion subnet.

WARNING [redundancy]: Rules "allow-http" and "allow-web-traffic" have identical 
  match conditions and actions. Consider removing one.

Plan: 12 rules to add, 3 to remove, 2 to modify.

  + [P10]  allow-ssh-bastion      TCP :22 ← 10.0.1.0/24         ACCEPT
  + [P10]  deny-ssh-all           TCP :22 ← 0.0.0.0/0           DROP+LOG
  + [P500] allow-http             TCP :80,:443 ← 0.0.0.0/0      ACCEPT
  ~ [P600] allow-prometheus       TCP :9090 ← 10.0.10.0/24      ACCEPT
    (changed: added port 9100)
  - [P800] old-temp-debug         TCP :9999 ← 10.0.1.100/32     ACCEPT
    (reason: schedule expired)

Apply? [y/N]:
```

---

## 12. Dry-Run & Policy Simulation

### 12.1 Dry-Run Mode

`rampart plan` compiles policies and shows the execution plan without applying:

```bash
# Show what would change
rampart plan -f policy.yaml

# Output as JSON (for CI/CD pipelines)
rampart plan -f policy.yaml -o json

# Compare with current state
rampart plan -f policy.yaml --diff
```

### 12.2 Policy Simulation

Test if a specific packet would be allowed or denied:

```bash
# Simulate a packet
rampart simulate \
  --src 10.0.1.50 \
  --dst 192.168.1.10 \
  --protocol tcp \
  --dport 22 \
  --direction inbound

# Output:
# ACCEPT by rule "allow-ssh-bastion" (policy: ssh-access, priority: 10)
# Match path: src 10.0.1.50 ∈ 10.0.1.0/24, protocol tcp, dport 22
```

```bash
# Simulate against a policy file (not yet applied)
rampart simulate \
  --src 203.0.113.50 \
  --dst 192.168.1.10 \
  --protocol tcp \
  --dport 3306 \
  --direction inbound \
  -f new-policy.yaml

# Output:
# DROP by default policy (no matching rule)
# Evaluated 12 rules, none matched.
```

### 12.3 Simulation Engine

```go
type Simulator struct {
    compiledRules []CompiledRule // Priority-sorted
}

type SimulatedPacket struct {
    SourceIP   net.IP
    DestIP     net.IP
    Protocol   Protocol
    SourcePort uint16
    DestPort   uint16
    Direction  Direction
    Interface  string
    State      ConnState
}

type SimulationResult struct {
    Verdict     Action        // Accept, Drop, Reject
    MatchedRule *CompiledRule // nil if default policy
    MatchPath   string        // Human-readable match explanation
    Evaluated   int           // Number of rules evaluated
    Duration    time.Duration // Simulation time
}

func (s *Simulator) Simulate(pkt SimulatedPacket) SimulationResult
```

---

## 13. Snapshot & Rollback Engine

### 13.1 Snapshot Model

```go
type Snapshot struct {
    ID          string    // UUID v7
    CreatedAt   time.Time
    CreatedBy   string    // User/system that triggered
    Trigger     string    // "manual", "pre-apply", "scheduled"
    Description string
    PolicyHash  string    // SHA-256 of compiled policy
    RuleCount   int
    Backend     string    // Which backend was active
    State       []byte    // Serialized backend state (gob-encoded)
    Metadata    map[string]string
}
```

### 13.2 Automatic Snapshots

Rampart automatically creates snapshots:

- **Pre-apply:** Before every `rampart apply` operation
- **Post-rollback:** After every rollback (for rollback-of-rollback)
- **Scheduled:** Configurable periodic snapshots (default: every 6 hours)

### 13.3 Snapshot Storage

```
/var/lib/rampart/snapshots/
  ├── 01JQXYZ123-pre-apply.snap
  ├── 01JQXYZ456-scheduled.snap
  └── 01JQXYZ789-manual.snap
```

Storage format: gob-encoded Go structs with zstd compression.

Retention policy:
- Keep last 100 snapshots (configurable)
- Keep all snapshots from last 30 days
- Purge older snapshots automatically

### 13.4 Rollback Operations

```bash
# List available snapshots
rampart snapshot list

# ID                    CREATED              TRIGGER    RULES  DESCRIPTION
# 01JQXYZ789           2026-04-11 10:30:00  manual     12     Before maintenance
# 01JQXYZ456           2026-04-11 04:00:00  scheduled  12     Scheduled snapshot
# 01JQXYZ123           2026-04-10 15:22:00  pre-apply  10     Pre-apply: policy update

# Rollback to a specific snapshot
rampart rollback 01JQXYZ789

# Rollback to the previous state (shortcut)
rampart rollback --last

# Show diff between current state and a snapshot
rampart snapshot diff 01JQXYZ789

# Export snapshot as YAML policy (reverse-compile)
rampart snapshot export 01JQXYZ789 -o policy-backup.yaml
```

---

## 14. Audit System

### 14.1 Audit Event Model

```go
type AuditEvent struct {
    ID        string          // UUID v7
    Timestamp time.Time       // Nanosecond precision
    NodeID    string          // Which cluster node
    Actor     AuditActor      // Who did it
    Action    AuditAction     // What was done
    Resource  AuditResource   // What was affected
    Before    json.RawMessage // State before change (nil for creates)
    After     json.RawMessage // State after change (nil for deletes)
    Result    AuditResult     // Success, Failure, DryRun
    Metadata  map[string]string
}

type AuditActor struct {
    Type     string // "user", "api", "system", "mcp", "raft-sync"
    Identity string // Username, API key ID, "system:scheduler"
    SourceIP string // Client IP address
}

type AuditAction string

const (
    AuditApply    AuditAction = "policy.apply"
    AuditRollback AuditAction = "policy.rollback"
    AuditPlan     AuditAction = "policy.plan"
    AuditSimulate AuditAction = "policy.simulate"
    AuditSnapshot AuditAction = "snapshot.create"
    AuditFlush    AuditAction = "policy.flush"
    AuditSync     AuditAction = "cluster.sync"
    AuditJoin     AuditAction = "cluster.join"
    AuditLeave    AuditAction = "cluster.leave"
)
```

### 14.2 Audit Storage

Append-only file with JSON Lines format:

```
/var/lib/rampart/audit/
  ├── audit-2026-04-11.jsonl
  ├── audit-2026-04-10.jsonl
  └── audit-2026-04-09.jsonl
```

Rotation: daily files, gzip compressed after 24 hours, configurable retention (default: 90 days).

### 14.3 Audit Query API

```bash
# View recent audit events
rampart audit list --last 20

# Filter by action
rampart audit list --action policy.apply --since 2026-04-01

# Filter by actor
rampart audit list --actor ersin --since 2026-04-01

# Show diff for a specific event
rampart audit show 01JQXYZ789
```

---

## 15. Raft Cluster & Multi-Host Sync

### 15.1 Cluster Architecture

```
         ┌─────────────────────┐
         │     Raft Leader     │
         │  (Policy Authority) │
         │                     │
         │  ┌───────────────┐  │
         │  │ Policy Store  │  │
         │  │ Audit Log     │  │
         │  │ Snapshot Store│  │
         │  └───────────────┘  │
         └──────────┬──────────┘
                    │ Raft replication
         ┌──────────┼──────────┐
         │          │          │
    ┌────▼────┐ ┌───▼────┐ ┌──▼─────┐
    │Follower │ │Follower│ │Follower│
    │ Node 1  │ │ Node 2 │ │ Node 3 │
    │         │ │        │ │        │
    │ nftables│ │iptables│ │  eBPF  │
    │ backend │ │backend │ │backend │
    └─────────┘ └────────┘ └────────┘
```

### 15.2 What Gets Replicated

- **Compiled policy state** (the desired ruleset)
- **Audit events** (cluster-wide audit trail)
- **Snapshot metadata** (actual snapshots stored locally)

### 15.3 What Stays Local

- **Backend implementation** (each node can run a different backend)
- **Snapshot data** (raw firewall state is node-specific)
- **Rule counters** (per-node statistics)

### 15.4 Raft Implementation

Custom Raft implementation (no hashicorp/raft dependency):

```go
type RaftNode struct {
    id          string
    state       NodeState    // Leader, Follower, Candidate
    currentTerm uint64
    votedFor    string
    log         []LogEntry
    commitIndex uint64
    lastApplied uint64
    
    // Leader-only
    nextIndex   map[string]uint64
    matchIndex  map[string]uint64
    
    // Transport
    transport   Transport // TCP + TLS
    peers       []Peer
}

type LogEntry struct {
    Term    uint64
    Index   uint64
    Type    EntryType   // PolicyUpdate, ConfigChange, NodeJoin, NodeLeave
    Data    []byte      // gob-encoded
}
```

### 15.5 Sync Flow

1. User runs `rampart apply -f policy.yaml` on any node
2. If not leader → forward to leader via Raft
3. Leader compiles policy, validates, creates audit event
4. Leader proposes log entry via Raft
5. Once committed (majority ack), leader responds to client
6. Each follower applies the committed policy to its local backend
7. Followers report apply result back (for monitoring)

### 15.6 Node Operations

```bash
# Initialize a new cluster (first node becomes leader)
rampart cluster init --listen 0.0.0.0:7946 --advertise 10.0.1.1:7946

# Join an existing cluster
rampart cluster join --leader 10.0.1.1:7946 --listen 0.0.0.0:7946

# Leave the cluster gracefully
rampart cluster leave

# Show cluster status
rampart cluster status

# NODE          STATE     BACKEND    RULES  LAST-SYNC           HEALTHY
# 10.0.1.1:7946 leader    nftables   12     2026-04-11 10:30:00 ✓
# 10.0.1.2:7946 follower  nftables   12     2026-04-11 10:30:01 ✓
# 10.0.1.3:7946 follower  iptables   12     2026-04-11 10:30:01 ✓

# Force leader election (emergency)
rampart cluster elect --force
```

### 15.7 TLS Mutual Authentication

All Raft communication is encrypted with mTLS:

```bash
# Generate cluster CA and node certificates
rampart cert init --ca-dir /etc/rampart/certs
rampart cert generate --node-name node-1 --ca-dir /etc/rampart/certs
```

---

## 16. REST API

### 16.1 Authentication

- **API Key:** `Authorization: Bearer rmp_xxxx`
- **mTLS:** Client certificate authentication
- **Local socket:** `/var/run/rampart.sock` (no auth needed for local access)

### 16.2 Endpoints

#### Policy Management

```
POST   /api/v1/policies/plan        # Dry-run: compile and show plan
POST   /api/v1/policies/apply       # Apply policy (body: YAML or JSON)
POST   /api/v1/policies/simulate    # Simulate a packet
GET    /api/v1/policies/current     # Get current active policy
GET    /api/v1/policies/conflicts   # Get conflict report for current policy
DELETE /api/v1/policies             # Flush all Rampart rules
```

#### Snapshots

```
GET    /api/v1/snapshots            # List snapshots
POST   /api/v1/snapshots            # Create manual snapshot
GET    /api/v1/snapshots/:id        # Get snapshot details
POST   /api/v1/snapshots/:id/rollback  # Rollback to snapshot
GET    /api/v1/snapshots/:id/diff   # Diff snapshot vs current
GET    /api/v1/snapshots/:id/export # Export as YAML policy
DELETE /api/v1/snapshots/:id        # Delete snapshot
```

#### Audit

```
GET    /api/v1/audit                # List audit events (paginated)
GET    /api/v1/audit/:id            # Get audit event details
GET    /api/v1/audit/search         # Search audit events
```

#### Cluster

```
GET    /api/v1/cluster/status       # Cluster status
POST   /api/v1/cluster/join         # Join cluster
POST   /api/v1/cluster/leave        # Leave cluster
GET    /api/v1/cluster/nodes        # List nodes
GET    /api/v1/cluster/health       # Health check
```

#### Rules (Quick CRUD — bypass YAML workflow)

```
GET    /api/v1/rules                # List active rules
POST   /api/v1/rules                # Add single rule
GET    /api/v1/rules/:id            # Get rule details
PUT    /api/v1/rules/:id            # Update rule
DELETE /api/v1/rules/:id            # Delete rule
GET    /api/v1/rules/:id/stats      # Get rule counters
```

#### System

```
GET    /api/v1/system/info          # Version, backend, uptime
GET    /api/v1/system/backends      # Available backends
GET    /api/v1/system/health        # Health check
GET    /api/v1/system/metrics       # Prometheus metrics
```

### 16.3 Response Format

```json
{
  "status": "success",
  "data": { ... },
  "meta": {
    "requestId": "01JQXYZ789",
    "timestamp": "2026-04-11T10:30:00Z",
    "node": "node-1"
  }
}
```

Error response:

```json
{
  "status": "error",
  "error": {
    "code": "CONFLICT_DETECTED",
    "message": "2 rule conflicts detected",
    "details": [ ... ]
  },
  "meta": { ... }
}
```

---

## 17. CLI Design

### 17.1 Command Structure

```
rampart
├── serve                    # Start server (API + WebUI + Raft)
├── agent                    # Start agent mode (follower-only)
├── apply                    # Apply policy from YAML file
├── plan                     # Dry-run: show execution plan
├── simulate                 # Simulate a packet
├── rollback                 # Rollback to a snapshot
├── snapshot
│   ├── list                 # List snapshots
│   ├── create               # Create manual snapshot
│   ├── diff                 # Diff snapshot vs current
│   └── export               # Export snapshot as YAML
├── rules
│   ├── list                 # List active rules
│   ├── add                  # Add single rule (quick mode)
│   ├── remove               # Remove single rule
│   └── stats                # Show rule counters
├── audit
│   ├── list                 # List audit events
│   └── show                 # Show audit event detail
├── cluster
│   ├── init                 # Initialize new cluster
│   ├── join                 # Join existing cluster
│   ├── leave                # Leave cluster
│   ├── status               # Show cluster status
│   └── elect                # Force leader election
├── cert
│   ├── init                 # Generate cluster CA
│   └── generate             # Generate node certificate
├── validate                 # Validate YAML policy file
├── fmt                      # Format YAML policy file
├── diff                     # Diff two policy files
├── import                   # Import from iptables-save / nft list
├── export                   # Export current rules as YAML policy
└── version                  # Show version info
```

### 17.2 Common Usage Examples

```bash
# Start Rampart server with WebUI
rampart serve --config /etc/rampart/rampart.yaml

# Apply a policy
rampart apply -f production-web.yaml

# Plan first, then apply
rampart plan -f production-web.yaml
rampart apply -f production-web.yaml --auto-approve

# Quick add a rule (no YAML needed)
rampart rules add \
  --name temp-debug \
  --protocol tcp \
  --dport 9999 \
  --source 10.0.1.100/32 \
  --action accept \
  --ttl 2h  # Auto-expire in 2 hours

# Import existing iptables rules
rampart import --from iptables-save --output current-rules.yaml

# Check what rules are active
rampart rules list

# Rollback to previous state
rampart rollback --last

# Validate a policy file
rampart validate -f policy.yaml

# Format a policy file (consistent style)
rampart fmt -f policy.yaml
```

---

## 18. React WebUI

### 18.1 Dashboard

Embedded React application served by the Go binary at `https://<host>:9443/ui/`.

### 18.2 Pages

| Page | Description |
|------|-------------|
| **Dashboard** | Overview: active rules count, cluster status, recent audit events, rule hit heatmap |
| **Policies** | YAML editor with syntax highlighting, live validation, conflict warnings |
| **Rules** | Table view of active rules with search, filter, sort. Per-rule stats (packets/bytes) |
| **Simulator** | Interactive packet simulation form. Visual trace of rule evaluation |
| **Snapshots** | Timeline view of snapshots. Diff viewer (side-by-side). One-click rollback |
| **Audit Log** | Searchable, filterable audit event timeline. Before/after diff for each event |
| **Cluster** | Node list with health status. Leader indicator. Raft log status |
| **Settings** | Backend config, API keys, notification webhooks, snapshot retention |

### 18.3 Technical Stack

- **React 19** with TypeScript
- **Tailwind CSS v4** for styling
- **CodeMirror 6** for YAML editor with Rampart schema validation
- **D3.js** for rule hit heatmap and network topology
- **Built at compile time** → embedded in Go binary via `embed.FS`
- **SPA** with client-side routing
- **SSE** for real-time updates (audit events, rule counter refresh, cluster status)

### 18.4 Build & Embedding

```go
//go:embed ui/dist/*
var uiFS embed.FS

func (s *Server) serveUI() http.Handler {
    return http.FileServer(http.FS(uiFS))
}
```

---

## 19. MCP Server

### 19.1 Tools

| Tool | Description |
|------|-------------|
| `list_rules` | List active firewall rules with optional filters |
| `add_rule` | Add a single rule (quick mode) |
| `remove_rule` | Remove a rule by name or ID |
| `plan_policy` | Dry-run a YAML policy and show execution plan |
| `apply_policy` | Apply a YAML policy (requires confirmation) |
| `simulate_packet` | Test if a packet would be allowed/denied |
| `rollback` | Rollback to a specific snapshot |
| `list_snapshots` | List available snapshots |
| `audit_search` | Search audit events |
| `cluster_status` | Show cluster node status |
| `get_rule_stats` | Get packet/byte counters for a rule |

### 19.2 MCP Resources

| Resource | Description |
|----------|-------------|
| `rampart://policies/current` | Current active policy as YAML |
| `rampart://rules` | Active rules in JSON format |
| `rampart://audit/recent` | Recent audit events |
| `rampart://cluster/status` | Cluster health and node list |

### 19.3 Example Interaction

```
User: "Block all traffic from 203.0.113.0/24 for the next hour"
AI → rampart.add_rule({
  name: "block-suspicious-range",
  source: "203.0.113.0/24",
  action: "drop",
  log: true,
  ttl: "1h"
})
→ "Rule 'block-suspicious-range' applied. Will auto-expire at 11:30 UTC."
```

---

## 20. Time-Based Rules & Scheduling

### 20.1 One-Time Rules (TTL)

```bash
# CLI
rampart rules add --name temp-access --dport 8080 --action accept --ttl 2h

# YAML
schedule:
  activeFrom: "2026-04-15T10:00:00Z"
  activeUntil: "2026-04-15T14:00:00Z"
```

### 20.2 Recurring Rules

```yaml
schedule:
  recurring:
    days: [monday, tuesday, wednesday, thursday, friday]
    startTime: "09:00"
    endTime: "17:00"
    timezone: "Europe/Istanbul"
```

### 20.3 Scheduler Implementation

- Background goroutine checks rule schedules every 30 seconds
- On activation: rule is compiled and applied to backend
- On deactivation: rule is removed from backend
- Audit event logged for both activation and deactivation
- Schedule evaluation is deterministic (same input → same result on all cluster nodes)

---

## 21. Configuration

### 21.1 Server Configuration

```yaml
# /etc/rampart/rampart.yaml
server:
  listen: "0.0.0.0:9443"
  unixSocket: "/var/run/rampart.sock"
  tls:
    cert: /etc/rampart/certs/server.crt
    key: /etc/rampart/certs/server.key
    ca: /etc/rampart/certs/ca.crt  # For mTLS

backend:
  type: auto  # auto | nftables | iptables | ebpf | hybrid | aws | gcp | azure
  # nftables-specific
  nftables:
    tableName: rampart
    binary: /usr/sbin/nft
  # iptables-specific
  iptables:
    chainPrefix: RAMPART
    binary: /usr/sbin/iptables
  # eBPF-specific
  ebpf:
    xdpMode: native  # native | skb | offload
    interface: eth0
  # AWS-specific
  aws:
    region: eu-west-1
    securityGroupId: sg-0123456789
  # Hybrid mode
  hybrid:
    fastPath: ebpf
    slowPath: nftables

cluster:
  enabled: true
  nodeId: node-1
  listen: "0.0.0.0:7946"
  advertise: "10.0.1.1:7946"
  peers:
    - "10.0.1.2:7946"
    - "10.0.1.3:7946"
  tls:
    cert: /etc/rampart/certs/node.crt
    key: /etc/rampart/certs/node.key
    ca: /etc/rampart/certs/ca.crt

snapshots:
  directory: /var/lib/rampart/snapshots
  retention:
    maxCount: 100
    maxAge: 720h  # 30 days
  autoSnapshot:
    interval: 6h
    preApply: true

audit:
  directory: /var/lib/rampart/audit
  retention: 2160h  # 90 days
  compress: true     # gzip after 24h

scheduler:
  checkInterval: 30s

api:
  keys:
    - name: admin
      key: rmp_xxxxx
      permissions: ["*"]
    - name: readonly
      key: rmp_yyyyy
      permissions: ["read"]

webui:
  enabled: true
  path: /ui

mcp:
  enabled: true
  listen: "127.0.0.1:9444"

logging:
  level: info  # debug, info, warn, error
  format: json # json, text
  output: stderr
  file: /var/log/rampart/rampart.log

metrics:
  enabled: true
  path: /metrics
```

---

## 22. Security Model

### 22.1 Principle of Least Privilege

- Rampart server runs as root (required for nftables/iptables manipulation)
- Drops all capabilities except `CAP_NET_ADMIN`, `CAP_NET_RAW`
- After binding ports, drops `CAP_NET_BIND_SERVICE`
- WebUI and API accessible without root

### 22.2 API Security

- All API communication over TLS (HTTP/2)
- API keys stored as bcrypt hashes in config
- Rate limiting per API key (configurable)
- RBAC: read, write, admin permissions
- CORS configuration for WebUI

### 22.3 Cluster Security

- mTLS for all Raft communication
- Certificate pinning (CA must match)
- Raft messages signed with shared secret (HMAC-SHA256)
- Node authentication before cluster join

### 22.4 Audit Security

- Audit log is append-only (no modification API)
- Audit entries include cryptographic hash chain
- Each entry's hash = SHA-256(previous_hash + entry_data)
- Tamper detection: verify hash chain integrity

### 22.5 Policy File Security

- YAML includes: only local files and HTTPS URLs
- URL includes: TLS verification, certificate pinning optional
- Variable substitution: no shell expansion (prevent injection)
- Max policy file size: 10 MB (configurable)

---

## 23. Observability

### 23.1 Prometheus Metrics

```
# Rule statistics
rampart_rule_packets_total{rule="allow-ssh", policy="ssh-access"} 1234
rampart_rule_bytes_total{rule="allow-ssh", policy="ssh-access"} 567890

# Backend operations
rampart_apply_duration_seconds{backend="nftables"} 0.015
rampart_apply_total{backend="nftables", status="success"} 42

# Cluster health
rampart_raft_term 5
rampart_raft_commit_index 128
rampart_raft_peers 3
rampart_raft_state{state="leader"} 1

# Snapshot
rampart_snapshots_total 25
rampart_snapshot_size_bytes 4096

# Audit
rampart_audit_events_total{action="policy.apply"} 42

# Scheduler
rampart_scheduled_rules_active 3
rampart_scheduled_rules_total 7
```

### 23.2 Structured Logging

All logs are structured JSON:

```json
{
  "time": "2026-04-11T10:30:00.123Z",
  "level": "info",
  "msg": "policy applied",
  "component": "engine",
  "rules_added": 3,
  "rules_removed": 1,
  "backend": "nftables",
  "duration_ms": 15,
  "audit_id": "01JQXYZ789"
}
```

### 23.3 Health Check

```
GET /api/v1/system/health

{
  "status": "healthy",
  "checks": {
    "backend": {"status": "ok", "backend": "nftables"},
    "cluster": {"status": "ok", "role": "leader", "peers": 2},
    "storage": {"status": "ok", "snapshots": 25},
    "audit": {"status": "ok", "lastEvent": "2026-04-11T10:30:00Z"}
  }
}
```

---

## 24. Performance Targets

| Metric | Target |
|--------|--------|
| Policy compilation (100 rules) | < 10ms |
| Policy compilation (10,000 rules) | < 500ms |
| Rule apply (nftables, atomic) | < 50ms |
| Rule apply (iptables, chain swap) | < 200ms |
| Conflict detection (1,000 rules) | < 100ms |
| Packet simulation | < 1ms |
| Snapshot create | < 100ms |
| Rollback (nftables) | < 50ms |
| Raft consensus (3 nodes, LAN) | < 10ms |
| API response (rule list) | < 5ms |
| Memory usage (1,000 rules) | < 50 MB |
| Binary size | < 30 MB |
| Cold start to serving | < 2s |

---

## 25. Project Structure

```
rampart/
├── cmd/
│   └── rampart/
│       └── main.go              # Unified binary entry point
├── internal/
│   ├── engine/                  # Policy engine core
│   │   ├── compiler.go          # YAML → compiled rules
│   │   ├── conflict.go          # Conflict detection
│   │   ├── simulator.go         # Packet simulation
│   │   ├── scheduler.go         # Time-based rule scheduler
│   │   ├── variable.go          # Variable substitution
│   │   └── validator.go         # Policy validation
│   ├── backend/                 # Backend abstraction layer
│   │   ├── backend.go           # Interface definition
│   │   ├── registry.go          # Backend registry
│   │   ├── nftables/            # nftables backend
│   │   │   ├── nftables.go
│   │   │   ├── compiler.go      # Rule → nft commands
│   │   │   ├── parser.go        # nft JSON output parser
│   │   │   └── sets.go          # Named set management
│   │   ├── iptables/            # iptables backend
│   │   │   ├── iptables.go
│   │   │   ├── compiler.go
│   │   │   └── parser.go        # iptables-save parser
│   │   ├── ebpf/                # eBPF/XDP backend
│   │   │   ├── ebpf.go
│   │   │   ├── loader.go        # BPF program loader
│   │   │   ├── maps.go          # BPF map management
│   │   │   └── programs/        # Pre-compiled BPF bytecode
│   │   ├── aws/                 # AWS Security Groups
│   │   │   ├── aws.go
│   │   │   └── sigv4.go         # AWS SigV4 signing (from scratch)
│   │   ├── gcp/                 # GCP Firewall Rules
│   │   │   └── gcp.go
│   │   └── azure/               # Azure NSGs
│   │       └── azure.go
│   ├── snapshot/                # Snapshot & rollback engine
│   │   ├── snapshot.go
│   │   ├── store.go             # File-based snapshot storage
│   │   └── retention.go         # Cleanup policy
│   ├── audit/                   # Audit system
│   │   ├── audit.go
│   │   ├── store.go             # JSONL append-only store
│   │   ├── query.go             # Search/filter
│   │   └── integrity.go         # Hash chain verification
│   ├── cluster/                 # Raft cluster
│   │   ├── raft.go              # Raft consensus
│   │   ├── log.go               # Raft log
│   │   ├── transport.go         # TCP + TLS transport
│   │   ├── fsm.go               # Finite state machine
│   │   └── discovery.go         # Peer discovery
│   ├── api/                     # REST API
│   │   ├── server.go            # HTTP server
│   │   ├── middleware.go        # Auth, logging, CORS
│   │   ├── handlers/            # Route handlers
│   │   │   ├── policy.go
│   │   │   ├── rules.go
│   │   │   ├── snapshot.go
│   │   │   ├── audit.go
│   │   │   ├── cluster.go
│   │   │   └── system.go
│   │   └── sse.go               # Server-Sent Events
│   ├── mcp/                     # MCP server
│   │   ├── server.go
│   │   ├── tools.go
│   │   └── resources.go
│   ├── cli/                     # CLI commands
│   │   ├── root.go
│   │   ├── serve.go
│   │   ├── apply.go
│   │   ├── plan.go
│   │   ├── simulate.go
│   │   ├── rollback.go
│   │   ├── snapshot.go
│   │   ├── rules.go
│   │   ├── audit.go
│   │   ├── cluster.go
│   │   ├── cert.go
│   │   ├── validate.go
│   │   ├── fmt.go
│   │   ├── diff.go
│   │   ├── importcmd.go
│   │   └── export.go
│   ├── model/                   # Core data models
│   │   ├── rule.go
│   │   ├── policy.go
│   │   ├── snapshot.go
│   │   ├── audit.go
│   │   └── cluster.go
│   ├── config/                  # Configuration
│   │   ├── config.go
│   │   └── defaults.go
│   ├── cert/                    # Certificate management
│   │   ├── ca.go
│   │   └── generate.go
│   └── version/                 # Build info
│       └── version.go
├── ui/                          # React WebUI
│   ├── src/
│   │   ├── App.tsx
│   │   ├── pages/
│   │   │   ├── Dashboard.tsx
│   │   │   ├── Policies.tsx
│   │   │   ├── Rules.tsx
│   │   │   ├── Simulator.tsx
│   │   │   ├── Snapshots.tsx
│   │   │   ├── AuditLog.tsx
│   │   │   ├── Cluster.tsx
│   │   │   └── Settings.tsx
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── api/
│   │   └── styles/
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── tailwind.config.ts
├── SPECIFICATION.md
├── IMPLEMENTATION.md
├── TASKS.md
├── BRANDING.md
├── README.md
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── .github/
    └── workflows/
        └── ci.yml
```

---

## 26. Version Roadmap

### Phase 1 — Core Engine (Weeks 1-3)
- Policy YAML parser & validator
- Rule compiler & normalizer
- Conflict detection engine
- nftables backend (primary)
- Snapshot & rollback engine
- Audit system (append-only JSONL)
- CLI: apply, plan, validate, fmt, rules, snapshot, rollback

### Phase 2 — iptables Backend & Simulation (Weeks 4-5)
- iptables backend implementation
- Packet simulation engine
- Import from iptables-save / nft list
- Export current state as YAML
- Policy diff tool
- Variable substitution
- Policy includes

### Phase 3 — REST API & WebUI (Weeks 6-8)
- REST API server (full CRUD)
- Authentication (API keys, mTLS)
- React WebUI (all 8 pages)
- SSE for real-time updates
- YAML editor with live validation
- Rule hit heatmap visualization
- Build & embed UI in Go binary

### Phase 4 — Raft Cluster (Weeks 9-11)
- Raft consensus implementation
- TLS transport
- Policy replication
- Multi-node apply flow
- Certificate management (CA + node certs)
- Cluster CLI commands
- Cluster UI page

### Phase 5 — Scheduling & Advanced Features (Weeks 12-13)
- Time-based rules (TTL + recurring)
- Scheduler background service
- Rate limiting rules
- Rule statistics & counters
- Prometheus metrics endpoint
- MCP server

### Phase 6 — Cloud & eBPF Backends (Weeks 14-16)
- eBPF/XDP backend
- Hybrid mode (eBPF + nftables)
- AWS Security Groups backend
- GCP Firewall Rules backend
- Azure NSG backend
- Cloud API authentication (SigV4, OAuth2)

### Phase 7 — Production Hardening (Weeks 17-18)
- Audit hash chain integrity verification
- Capability dropping (security hardening)
- Benchmark suite
- Load testing (10,000 rule sets)
- Documentation (man pages, website)
- Docker image
- Systemd service file
- CI/CD pipeline

---

## DEPENDENCY POLICY

### Allowed External Dependencies (Extended Stdlib Only)

| Package | Purpose | Justification |
|---------|---------|---------------|
| golang.org/x/crypto | TLS helpers, certificate management | Go extended stdlib |
| golang.org/x/sys | mmap, epoll, capabilities, syscalls | Go extended stdlib |
| gopkg.in/yaml.v3 | YAML config & policy parsing | Standard YAML parser |

**Everything else is built from scratch:**
- Raft consensus ✋ (no hashicorp/raft)
- HTTP router ✋ (no chi, gin, echo)
- CLI parser ✋ (no cobra, urfave/cli)
- eBPF loader ✋ (no cilium/ebpf)
- AWS SigV4 ✋ (no aws-sdk-go)
- GCP OAuth2 ✋ (no google/cloud)
- Azure auth ✋ (no azure-sdk-for-go)
- UUID generation ✋ (no google/uuid)
- JSON Lines ✋ (encoding/json is stdlib)
- Prometheus exposition ✋ (no prometheus/client_golang)
- Certificate generation ✋ (crypto/x509 is stdlib)

---

## TAGLINE & POSITIONING

**Primary tagline:** "Policy-as-Code Firewall. One Binary. Every Backend."

**Secondary taglines:**
- "Stop managing iptables rules by hand."
- "Terraform plan, but for firewalls."
- "The wall that remembers everything."

**Positioning statement:**
Rampart is a network policy engine that replaces manual firewall rule management with policy-as-code, audit trails, instant rollback, and multi-host synchronization — all in a single Go binary with zero dependencies. It speaks nftables, iptables, eBPF, and cloud security groups through a unified YAML interface.

