#!/bin/bash
set -e

# Rampart Auto-Installer v0.1.0
# Usage: curl -sSL rampartfw.com/install | sh

REPO="ersinkoc/Rampart"
VERSION="0.1.0"
BINARY_NAME="rampart"

echo "🛡️  Installing Rampart $VERSION..."

# 1. Detect OS/Arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

if [ "$OS" != "linux" ] && [ "$OS" != "darwin" ]; then
    echo "❌ Unsupported OS: $OS"; exit 1
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/v$VERSION/rampart-$OS-$ARCH"

# 2. Download Binary
echo "📥 Downloading from GitHub..."
curl -L -o /tmp/rampart $DOWNLOAD_URL
chmod +x /tmp/rampart

# 3. Install to /usr/local/bin
echo "🚀 Installing to /usr/local/bin..."
sudo mv /tmp/rampart /usr/local/bin/rampart

# 4. Setup Directories
echo "📁 Creating configuration directories..."
sudo mkdir -p /etc/rampart /var/lib/rampart/snapshots /var/lib/rampart/audit

# 5. Initialize default config if not exists
if [ ! -f /etc/rampart/rampart.yaml ]; then
    echo "⚙️  Generating default configuration..."
    cat <<EOF | sudo tee /etc/rampart/rampart.yaml
server:
  listen: "0.0.0.0:9443"

backend:
  type: "auto"

logging:
  level: "info"
  format: "text"

snapshots:
  directory: "/var/lib/rampart/snapshots"

audit:
  directory: "/var/lib/rampart/audit"
EOF
fi

# 6. Systemd Integration (Linux only)
if [ "$OS" == "linux" ]; then
    echo "🔄 Setting up Systemd service..."
    cat <<EOF | sudo tee /etc/systemd/system/rampart.service
[Unit]
Description=Rampart Unified Firewall Manager
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/rampart serve --config /etc/rampart/rampart.yaml
Restart=on-failure
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_RAW

[Install]
WantedBy=multi-user.target
EOF
    sudo systemctl daemon-reload
    echo "✅ Rampart service created. Start it with: sudo systemctl start rampart"
fi

echo ""
echo "🎉 Rampart $VERSION installed successfully!"
echo "👉 Run 'rampart version' to verify."
