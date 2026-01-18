# nexus

Isolated dev environments. SSH + Services. Works with Cursor, VSCode, AI agents.

```bash
curl -fsSL https://nexus.example.com/install.sh | bash
```

3 minutes → Workspace ready. Editor connected. Services running.

---

## Install

[Full Installation Guide](scripts/INSTALL.md)

**One-liner:**
```bash
curl -fsSL https://nexus.example.com/install.sh | bash
```

**Or:**
```bash
git clone https://github.com/nexus/nexus && cd nexus && make build
```

---

## Quick Start (After Install)

```bash
nexus auth github         # Login with GitHub
nexus ssh setup           # Generate/upload SSH keys
nexus workspace create owner/repo   # Create workspace
nexus workspace connect my-workspace  # Open in editor
```

Done. Your code is ready.

---

## For Staging Environment

[Staging Deployment](deploy/envs/staging/README.md) - Local testing/development

```bash
cd deploy/envs/staging
./ops/start.sh
```

---

## For Developers (This Project)

Build:
```bash
make build
```

Test:
```bash
make test
make test-coverage
```

Dev setup:
```bash
make dev-setup
```

---

## Documentation

- [Installation](scripts/INSTALL.md) - Binary download & setup
- [Staging](deploy/envs/staging/README.md) - Local server ops
- [M4 Implementation](docs/planning/M4/) - Current phase progress
- [Architecture](docs/specs/m3.md) - System design
