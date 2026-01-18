# Nexus Installation

Three ways to install nexus.

## 1. One-Liner (Recommended)

```bash
curl -fsSL https://nexus.example.com/install.sh | bash
```

Done in ~10 seconds. Install script:
- Detects OS (macOS/Linux) and architecture
- Downloads nexus binary
- Installs to ~/.local/bin/
- Sets up ~/.nexus/ and ~/.ssh/
- Prints next steps

## 2. Manual Download

Download binary from [releases](https://github.com/nexus/nexus/releases):

```bash
curl -L https://github.com/nexus/nexus/releases/download/v1.0.0/nexus-linux-amd64 \
  -o ~/.local/bin/nexus
chmod +x ~/.local/bin/nexus
```

Then setup:
```bash
nexus auth github
nexus ssh setup
mkdir -p ~/.nexus
```

## 3. From Source

```bash
git clone https://github.com/nexus/nexus
cd nexus
make build
mv bin/nexus ~/.local/bin/
```

---

## Prerequisites (Auto-Installed)

- **git** - Repository cloning
- **ssh** - Remote connections
- **curl** - Binary download
- **GitHub CLI** - GitHub authentication (auto-installed if missing)

---

## Verify Installation

```bash
nexus version
nexus auth status
```

---

## First Time Setup

```bash
nexus auth github        # Authenticate with GitHub
nexus ssh setup          # Generate/upload SSH keys
nexus workspace create owner/repo  # Create workspace
nexus workspace connect my-workspace  # Open in editor
```

---

## Environment Variables

```bash
# Override installation directory (default: ~/.local/bin)
export NEXUS_BIN_DIR=/usr/local/bin

# Override config directory (default: ~/.nexus)
export NEXUS_CONFIG_DIR=~/.config/nexus

# Override release version (default: latest)
export NEXUS_VERSION=v1.0.0
```

---

## Troubleshooting

**Binary not in PATH**
```bash
export PATH="$HOME/.local/bin:$PATH"
# Add this line to ~/.bashrc or ~/.zshrc
```

**GitHub CLI not installed**
```bash
# macOS
brew install gh

# Ubuntu/Debian
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | \
  sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | \
  sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update && sudo apt install gh
```

**Old version installed**
```bash
nexus version
# If version < latest, reinstall:
curl -fsSL https://nexus.example.com/install.sh | bash
```

---

## Security Notes

- Scripts downloaded from HTTPS only
- Binary checksums verified
- Install to user home directory (no sudo)
- SSH keys never transmitted
- GitHub tokens managed by `gh` CLI (system credential store)

---

## Support

Issues? Check [docs](../docs/) or open [GitHub issue](https://github.com/nexus/nexus/issues)
