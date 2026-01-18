# nexus

Isolated dev environments. SSH + Services.

```bash
curl -fsSL https://raw.githubusercontent.com/IniZio/nexus/main/scripts/install.sh | bash
```

## Quick Start

```bash
nexus auth github              # GitHub auth
nexus ssh setup                # SSH key setup
nexus workspace create owner/repo  # Create workspace
nexus workspace connect name   # Connect editor
```

## For Development

```bash
make build          # Build binary
make test           # Run tests
make test-coverage  # Coverage report
```

## Deployment

Local staging: `cd deploy/envs/staging && ./ops/start.sh`
