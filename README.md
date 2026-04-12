# Rampart

**Unified Network Policy Engine for Linux & Cloud**

Rampart is a high-performance network policy engine that abstracts the complexity of `nftables`, `iptables`, `eBPF/XDP`, and cloud security groups (AWS, GCP, Azure) behind a single, human-readable YAML policy-as-code interface.

Built entirely in Go with zero external dependencies (beyond extended stdlib), Rampart provides a unified way to manage host-level and cloud-level firewalls with strong consistency, auditability, and safety.

## 🚀 Key Features

- **Policy-as-Code:** Define your firewall rules in version-controllable YAML.
- **Unified Backends:** Manage `nftables`, `iptables`, and `eBPF/XDP` with the same policy.
- **Cloud Integration:** Sync policies to AWS Security Groups, GCP Firewall Rules, and Azure NSGs.
- **Raft Consensus:** Strong consistency across multi-host clusters.
- **Dry-Run Mode:** `rampart plan` shows exactly what will change before applying.
- **Safety First:** Conflict detection (shadowing, contradictions) and policy simulation.
- **Audit Trail:** Append-only cryptographic hash chain of every modification.
- **Rollback:** Snapshot-based instant recovery to any previous state.
- **WebUI & API:** Modern React dashboard and REST API for programmatic access.
- **MCP Server:** AI agent integration for intelligent firewall management.

## 🛠 Installation

```bash
# Clone the repository
git clone https://github.com/rampartfw/rampart
cd rampart

# Build the binary
make build

# Install (optional)
sudo cp rampart /usr/local/bin/
```

## 📖 Quick Start

1. **Initialize the config:**
   ```bash
   rampart serve --config rampart.yaml
   ```

2. **Define a policy (`web-tier.yaml`):**
   ```yaml
   apiVersion: rampart.dev/v1
   kind: PolicySet
   metadata:
     name: web-servers
   policies:
     - name: public-http
       priority: 500
       rules:
         - name: allow-http
           match:
             protocol: tcp
             destPorts: [80, 443]
             sourceCIDRs: ["0.0.0.0/0"]
           action: accept
   ```

3. **Plan and apply:**
   ```bash
   rampart plan -f web-tier.yaml
   rampart apply -f web-tier.yaml
   ```

## 🏗 Architecture

Rampart uses a modular architecture with a central **Policy Engine** that compiles abstract YAML policies into backend-specific rules. The **Backend Abstraction Layer (BAL)** ensures that different firewall implementations can be swapped without changing the policy definition.

For multi-host setups, a **Raft-based Cluster** ensures that all nodes converge to the same desired state, while maintaining local snapshots for autonomous operation.

## 🛡 Security

Rampart is built for production environments. It drops unnecessary Linux capabilities after startup, uses bcrypt for API key security, and provides an append-only audit log with cryptographic integrity checks.

## 📄 License

Rampart is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
