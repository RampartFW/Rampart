# 🏰 Rampart

**Policy-as-Code Firewall. One Binary. Every Backend.**

Rampart is a network policy engine that replaces manual firewall rule management with version-controlled YAML policies, full audit trails, instant rollback, and multi-host synchronization — all in a single Go binary with zero external dependencies.

---

## The Problem

```bash
# Day 1: Quick fix
iptables -A INPUT -p tcp --dport 3306 -j ACCEPT

# Day 90: Who opened MySQL to the world?
# Day 180: 50 servers, 50 different rulesets
# Day 270: Wrong rule → locked out of SSH → 3am panic
```

No version control. No audit trail. No rollback. No dry-run. No multi-host sync.

## The Solution

```yaml
# rampart-policy.yaml
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: production-web

policies:
  - name: ssh-access
    priority: 10
    rules:
      - name: allow-ssh-bastion
        match:
          protocol: tcp
          destPorts: [22]
          sourceCIDRs: [10.0.1.0/24]
        action: accept

      - name: deny-ssh-all
        match:
          protocol: tcp
          destPorts: [22]
        action: drop
        log: true

  - name: web-traffic
    priority: 500
    rules:
      - name: allow-http
        match:
          protocol: tcp
          destPorts: [80, 443]
        action: accept
```

```bash
$ rampart plan -f production-web.yaml

Rampart Policy Plan
====================
Plan: 4 rules to add, 0 to remove, 0 to modify.

  + [P10]  allow-ssh-bastion    TCP :22 ← 10.0.1.0/24     ACCEPT
  + [P10]  deny-ssh-all         TCP :22 ← 0.0.0.0/0       DROP+LOG
  + [P500] allow-http           TCP :80,:443 ← 0.0.0.0/0  ACCEPT

Apply? [y/N]: y
✓ Applied 4 rules in 12ms
```

---

## Architecture

```
           ┌────────────────────────────────────────────┐
           │            Policy-as-Code (YAML)           │
           └───────────────────┬────────────────────────┘
                               │
                    ┌──────────▼──────────┐
                    │   Rampart Engine    │
                    │                     │
                    │  Parser → Compiler  │
                    │  Conflict Detector  │
                    │  Simulator          │
                    │  Snapshot + Audit   │
                    │  Raft Cluster       │
                    └──────────┬──────────┘
                               │
          ┌──────────┬─────────┼─────────┬──────────┐
          │          │         │         │          │
     ┌────▼───┐ ┌───▼────┐ ┌──▼──┐ ┌───▼───┐ ┌───▼───┐
     │nftables│ │iptables│ │eBPF │ │AWS SG │ │GCP/Az│
     └────────┘ └────────┘ └─────┘ └───────┘ └───────┘
```

---

## Features

### Policy-as-Code
Define firewall rules in version-controlled YAML. Variables, includes, policy inheritance with priority levels (system → org → team → service).

### Dry-Run Mode
See exactly what will change before applying. Like `terraform plan` for firewalls.

### Conflict Detection
Detect shadowed, contradictory, and redundant rules at compile time — before they hit your firewall.

### Instant Rollback
Every `apply` creates an automatic snapshot. One command to revert: `rampart rollback --last`.

### Full Audit Trail
Every change recorded: who, when, what, before/after diff. Cryptographic hash chain for tamper detection.

### Multi-Host Sync
Raft consensus ensures all nodes converge to the same policy. No more configuration drift across 50 servers.

### Packet Simulation
Test if a packet would be allowed or denied without applying any rules:

```bash
rampart simulate --src 10.0.1.50 --dst 192.168.1.10 --protocol tcp --dport 22
# → ACCEPT by rule "allow-ssh-bastion" (priority: 10)
```

### Time-Based Rules
Temporary rules with auto-expiry — perfect for maintenance windows:

```bash
rampart rules add --name temp-debug --dport 9999 --action accept --ttl 2h
```

### Pluggable Backends
One YAML policy → any backend: nftables, iptables, eBPF/XDP, AWS Security Groups, GCP Firewall Rules, Azure NSGs.

### React Dashboard
Built-in WebUI with YAML editor, rule table, packet simulator, snapshot timeline, audit log, and cluster status.

### MCP Server
AI agents can manage your firewall through the Model Context Protocol.

---

## Quick Start

### Install

```bash
# Binary
curl -fsSL https://rampartfw.com/install.sh | sh

# Go install
go install github.com/rampartfw/rampart/cmd/rampart@latest

# Docker
docker run -d --cap-add NET_ADMIN --name rampart rampartfw/rampart serve
```

### First Policy

```bash
# Create a policy file
cat > policy.yaml << 'EOF'
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: my-server

policies:
  - name: basics
    priority: 10
    rules:
      - name: allow-ssh
        match:
          protocol: tcp
          destPorts: [22]
          sourceCIDRs: [10.0.0.0/8]
        action: accept
      - name: allow-http
        match:
          protocol: tcp
          destPorts: [80, 443]
        action: accept
EOF

# Validate
rampart validate -f policy.yaml

# Preview changes
rampart plan -f policy.yaml

# Apply
rampart apply -f policy.yaml

# Check active rules
rampart rules list
```

### Start Server (API + WebUI)

```bash
rampart serve --listen 0.0.0.0:9443
# WebUI: https://localhost:9443/ui/
# API:   https://localhost:9443/api/v1/
```

### Cluster Setup

```bash
# Node 1 (bootstrap)
rampart cluster init --listen 0.0.0.0:7946 --advertise 10.0.1.1:7946

# Node 2
rampart cluster join --leader 10.0.1.1:7946 --listen 0.0.0.0:7946

# Node 3
rampart cluster join --leader 10.0.1.1:7946 --listen 0.0.0.0:7946

# Status
rampart cluster status
```

---

## CLI Reference

```
rampart serve          Start server (API + WebUI + Raft)
rampart apply          Apply policy from YAML file
rampart plan           Dry-run: show execution plan
rampart simulate       Simulate a packet
rampart rollback       Rollback to a snapshot
rampart snapshot       Manage snapshots (list, create, diff, export)
rampart rules          Manage rules (list, add, remove, stats)
rampart audit          View audit trail (list, show)
rampart cluster        Cluster operations (init, join, leave, status)
rampart cert           Certificate management (init, generate)
rampart validate       Validate YAML policy file
rampart fmt            Format YAML policy file
rampart diff           Diff two policy files
rampart import         Import from iptables-save / nft list
rampart export         Export current rules as YAML policy
rampart version        Show version info
```

---

## Comparison

| Feature | Rampart | UFW | firewalld | Terraform | Ansible |
|---------|:-------:|:---:|:---------:|:---------:|:-------:|
| Policy-as-Code | ✅ | ❌ | ❌ | ✅ | ✅ |
| Dry-run | ✅ | ❌ | ❌ | ✅ | ⚠️ |
| Rollback | ✅ | ❌ | ❌ | ⚠️ | ❌ |
| Audit Trail | ✅ | ❌ | ❌ | ❌ | ❌ |
| Multi-host Sync | ✅ | ❌ | ❌ | ❌ | Push |
| Conflict Detection | ✅ | ❌ | ❌ | ❌ | ❌ |
| Packet Simulation | ✅ | ❌ | ❌ | ❌ | ❌ |
| Time-based Rules | ✅ | ❌ | ⚠️ | ❌ | ❌ |
| Multiple Backends | ✅ | ❌ | ❌ | Cloud | ✅ |
| eBPF/XDP | ✅ | ❌ | ❌ | ❌ | ❌ |
| WebUI | ✅ | ❌ | Cockpit | ❌ | AWX |
| Single Binary | ✅ | ❌ | ❌ | ✅ | ❌ |
| Zero Dependencies | ✅ | ❌ | ❌ | ✅ | ❌ |

---

## Technology

- **Language:** Go 1.23+ (zero external dependencies)
- **Allowed deps:** `golang.org/x/crypto`, `golang.org/x/sys`, `gopkg.in/yaml.v3`
- **Binary:** Single static binary (~30 MB) — server + agent + CLI unified
- **WebUI:** React 19 + TypeScript + Tailwind CSS v4 (embedded in binary)
- **Protocol:** REST API (JSON) + SSE (real-time) + MCP (AI agents)
- **Clustering:** Custom Raft consensus with mTLS transport
- **Storage:** File-based (snapshots: gob + zstd, audit: JSONL, config: YAML)

---

## License

Apache License 2.0

---

## Links

- **Website:** [rampartfw.com](https://rampartfw.com)
- **Documentation:** [docs.rampartfw.com](https://docs.rampartfw.com)
- **GitHub:** [github.com/rampartfw/rampart](https://github.com/rampartfw/rampart)
- **Docker Hub:** [rampartfw/rampart](https://hub.docker.com/r/rampartfw/rampart)

---

*Built with 🏰 by [ECOSTACK TECHNOLOGY OÜ](https://ecostack.dev)*

