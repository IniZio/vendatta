# vendetta

**Isolated development environments with remote SSH access**

Vendatta creates isolated development environments (Docker/LXC/QEMU) on a coordination server. Each workspace gets its own SSH access with your public key seeded, so you can SSH directly into your workspace from anywhere.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                  Coordination Server (Your Dev Host)             │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  vendetta coordination start                              │  │
│  │  - Runs on port 3001                                      │  │
│  │  - Manages workspace lifecycle                            │  │
│  │  - Maps SSH ports to workspaces                           │  │
│  └───────────────────────────────────────────────────────────┘  │
│                            │                                    │
│                            ▼                                    │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Providers (Docker/LXC/QEMU)                              │  │
│  │  - Spawn workspaces locally on this host                  │  │
│  │  - Configure SSH with user's public key                   │  │
│  └───────────────────────────────────────────────────────────┘  │
│                            │                                    │
│            ┌───────────────┼───────────────┐                    │
│            ▼               ▼               ▼                    │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐           │
│  │ Workspace 1 │   │ Workspace 2 │   │ Workspace 3 │           │
│  │ SSH: :2222  │   │ SSH: :2223  │   │ SSH: :2224  │           │
│  │ User: alice │   │ User: bob   │   │ User: carol │           │
│  └─────────────┘   └─────────────┘   └─────────────┘           │
└─────────────────────────────────────────────────────────────────┘
                            │
                            │ SSH (your key grants access)
                            ▼
                  ┌─────────────────┐
                  │  You (Anywhere) │
                  │  SSH Client     │
                  └─────────────────┘

# Example: SSH into your workspace
ssh -p 2222 dev@your-server.com
```

## 🚀 Complete Onboarding Guide

### Step 1: Generate Your SSH Key

```bash
# Generate an SSH key for vendetta (if you don't have one)
vendetta ssh generate

# This creates ~/.ssh/id_ed25519_vendetta
# Your public key will be displayed - share it with your admin
```

### Step 2: Get Access to a Workspace

**Option A: Register with the coordination server (if enabled)**
```bash
vendetta ssh register your-server.com:3001
```

**Option B: Share your public key with your administrator**
```
Your public key:
ssh-ed25519 AAAA... your-email@example.com
```

Your admin will add this key to your workspace and tell you:
- Server address (e.g., `dev.company.com`)
- SSH port (e.g., `2222`)
- Username (e.g., `alice`)

### Step 3: Connect to Your Workspace

```bash
# Get connection info including deep links
vendetta workspace connect my-feature

# Example output:
# 🔗 Workspace Connection Info
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Workspace: my-feature
# Path:      /home/user/my-app/.vendetta/worktrees/my-feature
#
# 🐚 SSH Access:
#   ssh -p 2222 alice@dev.company.com
#
# 💻 Deep Links:
#   VSCode:  vscode://vscode-remote/ssh-remote+dev.company.com:2222/home/alice
#   Cursor:  cursor://ssh/remote?host=dev.company.com&port=2222&user=alice
#
# 🌐 Services:
#   - http://localhost:23000 (web)
#   - http://localhost:23001 (api)
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Step 4: Open in Your Editor

**VSCode:**
```bash
# Click the deep link or use:
code --remote ssh-remote+dev.company.com:2222 /home/alice/my-feature
```

**Cursor:**
```bash
# Click the deep link or use cursor:// scheme
```

**Terminal:**
```bash
ssh -p 2222 alice@dev.company.com
```

### Step 5: Check Running Services

```bash
# List all services and their URLs
vendetta workspace services my-feature

# Output:
# 📦 Services for workspace 'my-feature'
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Service     Port     Local Port  URL
# web         3000     23000       http://localhost:23000
# api         4000     23001       http://localhost:23001
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## 📋 Command Reference

### SSH Management

| Command | Description |
|---------|-------------|
| `vendetta ssh generate` | Generate SSH key for vendetta access |
| `vendetta ssh register <server>` | Register your public key with coordination server |
| `vendetta ssh info <workspace>` | Show connection info for a workspace |

### Workspace Management

| Command | Description |
|---------|-------------|
| `vendetta workspace create <name>` | Create a new workspace |
| `vendetta workspace up <name>` | Start a workspace |
| `vendetta workspace down <name>` | Stop a workspace |
| `vendetta workspace list` | List all workspaces |
| `vendetta workspace rm <name>` | Remove a workspace |
| `vendetta workspace connect <name>` | Show connection info and deep links |
| `vendetta workspace services <name>` | List services and their URLs |

### Server Administration

| Command | Description |
|---------|-------------|
| `vendetta coordination start` | Start the coordination server |
| `vendetta coordination stop` | Stop the coordination server |
| `vendetta coordination status` | Show server status |

---

## 🔧 Configuration

### Server Configuration (`.vendetta/coordination.yaml`)

```yaml
server:
  host: "0.0.0.0"
  port: 3001
  auth_token: "change-this-in-production"

registry:
  provider: "memory"
  health_check_interval: "10s"
  node_timeout: "60s"

auth:
  enabled: false  # Enable for production with JWT
  jwt_secret: "your-secure-secret"
```

### Workspace Configuration (`.vendetta/config.yaml`)

```yaml
name: my-workspace
provider: docker

services:
  web:
    command: "npm run dev"
    port: 3000
  api:
    command: "npm run server"
    port: 4000

docker:
  image: node:20
```

---

## 👥 User Management

### Adding a User (Admin)

Administrators can add users and their SSH public keys via API:

```bash
curl -X POST http://localhost:3001/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "public_key": "ssh-ed25519 AAAA... alice@example.com",
    "workspace": "my-feature"
  }'
```

### List Users (Admin)

```bash
curl http://localhost:3001/api/v1/users
# Returns: {"users":[...],"count":N}
```

---

## ✨ Features

- **🔒 Isolated Environments**: Docker/LXC/QEMU containers and VMs
- **🔑 SSH Key Authentication**: Users access workspaces via SSH with their keys
- **🌐 Remote Access**: SSH to workspaces from anywhere
- **💻 Editor Integration**: Deep links for VSCode and Cursor
- **📦 Single Binary**: Zero dependencies on the coordination server
- **🔌 Plugin System**: Extensible rules, skills, and commands
- **🚀 Service Discovery**: Automatic port mapping and service URLs
- **🤖 AI Agent Ready**: Auto-configures Cursor, OpenCode, Claude, and more

---

## 📚 Documentation

- [Configuration Guide](docs/spec/product/configuration.md)
- [Plugin System](docs/spec/technical/plugins.md)
- [Architecture](docs/spec/technical/architecture.md)
- [Coordination API](docs/coordination-api.md)

---

## License

MIT
