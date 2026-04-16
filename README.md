# Rampart

**Unified Network Policy Engine for Linux & Cloud**

<p align="center">
  <img src="assets/banner.jpeg" alt="Rampart Banner" width="100%">
</p>

Rampart is a high-performance network policy engine that abstracts the complexity of Linux firewalls and cloud security groups behind a single, human-readable YAML interface. Designed for consistency, performance, and security.

## 🚀 Key Features

- **Unified Policy Engine:** Manage `nftables`, `iptables`, and `eBPF/XDP` with one YAML format.
- **High Performance eBPF/XDP:** Fast-path packet filtering at the network driver level.
- **Raft-based Clustering:** Secure, distributed policy synchronization with mTLS.
- **Time-based Rules:** Automatically activate or expire rules based on flexible schedules.
- **Packet Simulation:** Test and trace rule evaluation without affecting live traffic.
- **AI-Ready (MCP):** Native Model Context Protocol support for agentic management.
- **Observability:** Built-in Prometheus metrics and tamper-evident audit logs.
- **Security Hardened:** Capability dropping and bcrypt-secured API access.

## 🏗️ Architecture

Rampart compiles abstract policies into backend-specific rules. It supports multiple backends (nftables, iptables, eBPF) and uses a custom Raft implementation for cluster-wide consistency.

## 🛠️ Quick Start

### Build from source
```bash
make build
```

### Apply your first policy
```bash
./rampart apply -f test-policy.yaml
```

### Start distributed server
```bash
./rampart serve --config rampart.yaml
```

## 🔒 Security
Rampart drops unnecessary Linux capabilities after startup, uses bcrypt for API key security, and provides an append-only audit log with cryptographic integrity checks.

## 📖 Documentation
Detailed specifications and implementation guides can be found in the [.project/](.project/) directory.

## 📄 License
Rampart is licensed under the Apache License 2.0.
