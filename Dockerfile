# Build Web UI
FROM node:20-alpine AS ui-builder
WORKDIR /ui
COPY ui/package*.json ./
RUN npm ci
COPY ui/ ./
RUN npm run build

# Build Go Backend
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy pre-built UI from ui-builder stage
COPY --from=ui-builder /ui/dist ./ui/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o rampart ./cmd/rampart

# Final stage
FROM alpine:3.19
WORKDIR /app

# Install runtime dependencies (firewall tools)
RUN apk add --no-cache nftables iptables ip6tables ca-certificates

# Copy the binary from builder
COPY --from=builder /app/rampart /usr/local/bin/rampart

# Create necessary directories
RUN mkdir -p /etc/rampart /var/lib/rampart/snapshots /var/lib/rampart/audit

# Expose API and Raft ports
EXPOSE 9443 7946

# Set entrypoint
ENTRYPOINT ["rampart"]
CMD ["serve", "--config", "/etc/rampart/rampart.yaml"]

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD rampart version || exit 1
