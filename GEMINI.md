# Rampart - Unified Network Policy Engine

## Project Overview
Rampart is a high-performance, unified network policy engine that abstracts the complexity of Linux firewalls (`nftables`, `iptables`, `eBPF/XDP`) and cloud security groups (AWS, GCP, Azure) behind a single, human-readable YAML policy-as-code interface. It is designed for consistency, auditability, and safety in both single-host and multi-host (cluster) environments.

### Core Technologies
- **Language:** Go 1.26+ (strictly following a zero-external-dependency policy, using only stdlib + `golang.org/x/crypto`, `golang.org/x/sys`, and `gopkg.in/yaml.v3`).
- **Frontend:** React 19 (TypeScript, Tailwind CSS v4, Vite, D3.js) embedded into the Go binary.
- **Backends:** `nftables` (primary), `iptables` (fallback), `eBPF/XDP` (fast-path), and Cloud APIs (AWS, GCP, Azure).
- **Consensus:** Custom Raft implementation for multi-host policy synchronization.
- **State Management:** Local file-based storage for snapshots (zstd compressed) and audit logs (JSON Lines with cryptographic hash chains).

### Architecture Highlights
- **Policy Engine:** Compiles abstract YAML into backend-specific rules, performs conflict detection (shadowing, contradictions), and provides a packet simulation engine.
- **Backend Abstraction Layer (BAL):** A stateless adapter interface for various firewall implementations.
- **Clustering:** Raft-based consensus ensures all nodes converge to the same desired state.
- **Security:** Drops Linux capabilities after startup (`CAP_NET_ADMIN`, `CAP_NET_RAW` retained), uses mTLS for cluster communication, and provides a tamper-evident audit trail.

---

## Building and Running

### Prerequisites
- Go 1.26+
- Node.js & npm (for WebUI build)
- Linux (for native firewall backends)

### Key Commands
- **Build Everything:** `make build` (builds UI then compiles the `rampart` binary).
- **Run Server:** `./rampart serve --config rampart.yaml` (starts API, WebUI, and engine).
- **Plan Changes:** `./rampart plan -f policy.yaml` (dry-run to see execution plan and conflicts).
- **Apply Policy:** `./rampart apply -f policy.yaml` (compiles and applies policy to the local backend).
- **Run Tests:** `make test` or `go test -v ./...`.
- **Linting:** `make lint` or `go vet ./...`.
- **Clean Build:** `make clean`.

---

## Development Conventions

### Go Backend
- **Surgical Logic:** The Policy Engine is implemented using pure functions for determinism and testability (same input $\rightarrow$ same output).
- **Dependency Policy:** Do NOT add external dependencies. Implement required functionality (like UUIDs, routers, or CLI parsers) from scratch or use `golang.org/x` packages.
- **Stateless Backends:** Backends should not store internal state; they are adapters for the current policy.
- **Error Handling:** Use wrapped errors with component context (e.g., `fmt.Errorf("component.Operation: %w", err)`).
- **Security First:** Always verify capability dropping and ensure sensitive data (like API keys) is never logged or exposed.

### React Frontend (`ui/` directory)
- **Tech Stack:** React 19, TypeScript, Tailwind CSS v4.
- **Embedding:** The UI is built into `ui/dist` and embedded into the Go binary using `//go:embed`.
- **Real-time Updates:** Uses Server-Sent Events (SSE) for live audit logs and rule hit statistics.
- **Styling:** Follow the established Tailwind CSS v4 patterns for a modern, high-performance UI.

### Testing & Validation
- **Unit Tests:** Mandatory for engine logic, compiler, and conflict detection.
- **Integration Tests:** Used for verifying backend-specific rule application (requires root/privileged environment).
- **Verification:** Always run `make test` and `make lint` before concluding a task. Every policy change must be validated using the built-in conflict detector.

---

## Key Files & Directories
- `cmd/rampart/main.go`: Entry point for the unified binary.
- `internal/engine/`: Core policy compilation and conflict detection logic.
- `internal/backend/`: Implementations for `nftables`, `iptables`, etc.
- `internal/cluster/`: Raft consensus implementation.
- `internal/api/`: REST API handlers and custom router.
- `ui/`: React source code and Vite configuration.
- `.project/`: Detailed specifications and implementation guides.
- `rampart.yaml`: Main configuration file.
