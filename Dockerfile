# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o rampart ./cmd/rampart

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies (nftables, iptables)
RUN apk add --no-color nftables iptables ip6tables ca-certificates

# Copy the binary from builder
COPY --from=builder /app/rampart /usr/local/bin/rampart

# Create config and data directories
RUN mkdir -p /etc/rampart /var/lib/rampart/snapshots /var/lib/rampart/audit

# Expose API and Raft ports
EXPOSE 9443 7946

# Set entrypoint
ENTRYPOINT ["rampart"]
CMD ["serve", "--config", "/etc/rampart/rampart.yaml"]

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD rampart version || exit 1
