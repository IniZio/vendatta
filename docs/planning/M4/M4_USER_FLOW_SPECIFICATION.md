# M3 Staging Environment & User Flow Specification

**Status**: Planning / Specification  
**Date**: 2026-01-17  
**Priority**: Critical (Unblock end-to-end user flow)  
**Target Completion**: Before general availability

## Executive Summary

This spec defines the complete staging environment setup and the end-to-end user flow for designers/developers to:

1. Install CLI with one-line script
2. Auto-authenticate with GitHub
3. Generate/upload SSH keys
4. Create workspace (fork/clone repo)
5. Start development (services running, editor open)
6. Commit and create PRs

**Key Principle**: fully automated after initial install script. Everything auto-detected, auto-configured.

---

## Part 1: Staging Environment Architecture

### 1.1 Infrastructure Components

```
┌─────────────────────────────────────────────────────────────────┐
│                   STAGING ENVIRONMENT (Host Machine)             │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Coordination Server (nexus coordination start)           │   │
│  │  • Port: 3001 (HTTP API)                                │   │
│  │  • Port: 2222-2299 (SSH forwarding to workspaces)       │   │
│  │  Responsibilities:                                        │   │
│  │  - Track registered users (GitHub handle → SSH pubkey)  │   │
│  │  - Manage workspace lifecycle (create, up, down)        │   │
│  │  - Forward SSH connections to LXC containers            │   │
│  │  - Expose workspace metadata (ports, services, status)  │   │
│  │  - Health checking for driver nodes                     │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                 │                                │
│              ┌──────────────────┼──────────────────┐             │
│              │                  │                  │             │
│              ▼                  ▼                  ▼             │
│  ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐ │
│  │ LXC Driver Node  │ │   Docker Node    │ │   QEMU Node      │ │
│  │ (Staging)        │ │   (Optional)     │ │ (Optional)       │ │
│  │                  │ │                  │ │                  │ │
│  │ Runs:            │ │ Runs:            │ │ Runs:            │ │
│  │ • LXC daemon     │ │ • Docker daemon  │ │ • QEMU/KVM       │ │
│  │ • Agent process  │ │ • Agent process  │ │ • Agent process  │ │
│  │ • Workspaces     │ │ • Containers     │ │ • VMs            │ │
│  │   (LXC instances)│ │                  │ │                  │ │
│  │ • SSH server     │ │ • SSH server     │ │ • SSH server     │ │
│  │   (per workspace)│ │   (per container)│ │   (per VM)       │ │
│  └──────────────────┘ └──────────────────┘ └──────────────────┘ │
│                                                                   │
│  Agent Responsibilities:                                         │
│  - Communicate with coordination server (heartbeat, ready status)│
│  - Execute provider commands (create/start/stop container)      │
│  - Report workspace status (running, services, ports)           │
│  - Handle SSH key setup in containers                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ User SSH Access (port forwarding)
                              │
                   ┌──────────▼──────────┐
                   │   Designers/Devs    │
                   │   Local Machine     │
                   │                     │
                   │ • nexus CLI         │
                   │ • Editor (Cursor)   │
                   │ • SSH (from editor) │
                   └─────────────────────┘
```

### 1.2 Networking

**Ports**:
- `3001` - Coordination server HTTP API
- `2222-2299` - SSH forwarding (forwarded from containers to host via coordination server)
- Container internal: `22` (SSH server inside container)

**Flow**:
```
User SSH (e.g., ssh -p 2222 user@staging-server)
  ↓
Coordination Server (port 2222 → workspace container:22)
  ↓
LXC Container (SSH server + dev environment)
```

### 1.3 Key Data Structures

#### Workspace Metadata (stored by coordination server)

```yaml
workspace:
  id: "workspace-abc123"
  name: "my-project-feature"
  owner: "github-username"  # GitHub username
  
  config:
    provider: "lxc"         # docker, lxc, qemu
    image: "ubuntu:22.04"
    services:
      web: {port: 3000}
      api: {port: 4000}
  
  source:
    type: "github"          # github, gitlab, local
    repo: "owner/repo"
    branch: "main"
    is_fork: false          # true if user forked it
  
  status: "running"         # pending, creating, running, stopped
  
  ssh:
    port: 2222              # forwarded port on staging server
    user: "dev"             # user inside container
    pubkey_fingerprint: "..."
  
  services:
    web:
      status: "running"
      port: 3000
      mapped_port: 23000    # localhost:23000 on staging server
```

#### User Registration (in coordination server)

```yaml
users:
  github-username:
    github_username: "github-username"
    github_id: 123456789
    ssh_pubkey: "ssh-ed25519 AAAA..."
    ssh_pubkey_fingerprint: "..."
    registered_at: "2026-01-17T10:30:00Z"
    workspaces: ["workspace-abc123", ...]
```

---

## Part 2: End-to-End User Flow

### 2.0 Prerequisites (System Admin Setup - Once)

1. **Provision staging host** with:
   - OS: Linux (Ubuntu 22.04+)
   - LXC installed and configured (`lxd init`)
   - Docker optional (for multi-provider testing)
   - Public IP or domain name (for SSH access)

2. **Configure firewall**:
   - SSH port 22 open (for dev access to coordination server if needed)
   - HTTP port 3001 open (for CLI API calls)
   - SSH forwarding ports 2222-2299 (for workspace access)

3. **Start coordination server**:
   ```bash
   ssh admin@staging-server
   nexus coordination start
   # Runs on port 3001
   ```

4. **Initialize LXC driver node**:
   ```bash
   # On staging server
   nexus agent start --coordination-url http://localhost:3001
   # Or for remote driver node:
   nexus agent start --coordination-url http://staging-server:3001 --drivers lxc
   ```

---

### 2.1 User Flow: Designer Installing CLI & Starting First Workspace

```
┌─────────────────────────────────────────────────────────────────┐
│ Step 0: Designer Intent                                         │
├─────────────────────────────────────────────────────────────────┤
│ "I want to design for github.com/my-org/my-project"            │
│ Action: Runs install script with repo URL                      │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 1: One-Line Install & Initial Setup                        │
├─────────────────────────────────────────────────────────────────┤
│ Command:                                                         │
│ $ curl https://nexus.example.com/install.sh | bash \           │
│   -s -- --repo my-org/my-project --server staging.example.com  │
│                                                                  │
│ Script Does:                                                     │
│ 1. Download nexus CLI binary (latest stable)                    │
│ 2. Install to ~/.local/bin/nexus                               │
│ 3. Initialize ~/.nexus/config.yaml (with server address)       │
│ 4. Check system: git, ssh-keygen, editor (cursor/code/vim)     │
│ 5. Print: "GitHub authentication needed next"                  │
│ 6. Call: nexus workspace setup-from-repo my-org/my-project \   │
│    --server staging.example.com                                │
│                                                                  │
│ Exit Code: 0 on success, user continues to Step 2              │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 2: GitHub Authentication                                   │
├─────────────────────────────────────────────────────────────────┤
│ Check GitHub Auth Status:                                        │
│ - If `gh auth status` succeeds → GO TO Step 3a                 │
│ - If `gh` not installed → INSTALL: brew install gh (macOS)    │
│ - If not authenticated → START: gh auth login                  │
│                                                                  │
│ UX Output:                                                       │
│ ┌──────────────────────────────────────────────────────────┐   │
│ │ 🔐 GitHub Authentication                                │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ │ Checking GitHub CLI...                                │   │
│ │ Checking authentication...                        ❌    │   │
│ │                                                          │   │
│ │ Running: gh auth login                                 │   │
│ │ Follow the prompts to authenticate with GitHub        │   │
│ │ (This opens browser for OAuth)                        │   │
│ │                                                          │   │
│ │ ⏳ Waiting for authentication...                       │   │
│ │  Authentication complete! (github-username)          │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ │                                                          │   │
│ │ Next: Checking SSH keys...                             │   │
│ └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│ Exit Condition: User authenticated with GitHub, username known │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 3: SSH Key Setup (Auto-Detection → Generation → Upload)    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│ 3a: Detect Existing SSH Keys                                   │
│ ────────────────────────────────────────────────────────────   │
│ Check ~/.ssh/ for valid public keys:                            │
│ - id_ed25519.pub (preferred)                                   │
│ - id_rsa.pub (acceptable)                                      │
│ - id_ecdsa.pub (acceptable)                                    │
│                                                                  │
│ Decision Tree:                                                   │
│ ┌─────────────────────────┐                                    │
│ │ Valid key exists?       │                                    │
│ └────────┬────────────────┘                                    │
│          │                                                      │
│    YES  │  NO                                                  │
│         │                                                       │
│         ▼  ▼                                                    │
│     3b   3b'                                                   │
│                                                                  │
│ 3b: Use Existing Key (if valid public/private pair exists)     │
│ ───────────────────────────────────────────────────────────   │
│ - Read public key content (e.g., ~/.ssh/id_ed25519.pub)       │
│ - UX: "Using existing SSH key: ssh-ed25519 AAA...xyz"         │
│ - Go to 3c (Upload)                                            │
│                                                                  │
│ 3b': Generate New Key (if none exist or user prefers)          │
│ ────────────────────────────────────────────────────────────  │
│ - Run: ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -N "" \    │
│       -C "nexus-$(whoami)@$(hostname)-$(date +%s)"            │
│ - UX: "Generating SSH key (ed25519)..."                        │
│ - UX: " SSH key generated: ~/.ssh/id_ed25519"                │
│ - Continue to 3c                                               │
│                                                                  │
│ 3c: Upload Key to GitHub                                       │
│ ────────────────────────────────────────────────────────────  │
│ - Use GitHub CLI to upload:                                    │
│   $ gh ssh-key add ~/.ssh/id_ed25519.pub \                    │
│     --title "nexus@staging"                                    │
│ - Or use GitHub API if key already exists                      │
│ - UX: "Uploading SSH key to GitHub..."                         │
│ - UX: " SSH key registered with GitHub (fingerprint: ...)"   │
│                                                                  │
│ 3d: Register with Coordination Server                          │
│ ────────────────────────────────────────────────────────────  │
│ HTTP POST to staging-server:3001/api/v1/users/register        │
│                                                                  │
│ Request Body:                                                   │
│ {                                                               │
│   "github_username": "github-username",                        │
│   "github_id": 123456789,                                      │
│   "ssh_pubkey": "ssh-ed25519 AAAA...",                         │
│   "ssh_pubkey_fingerprint": "SHA256:..."                       │
│ }                                                               │
│                                                                  │
│ UX Output:                                                       │
│ ┌──────────────────────────────────────────────────────────┐   │
│ │ 🔑 SSH Key Setup                                         │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ │ Detecting SSH keys...                             (1) │   │
│ │ SSH key: ~/.ssh/id_ed25519 (ed25519)                    │   │
│ │                                                          │   │
│ │ Uploading to GitHub...                                │   │
│ │ Key fingerprint: SHA256:abcd1234...                     │   │
│ │                                                          │   │
│ │ Registering with coordination server...               │   │
│ │ Server: staging.example.com:3001                        │   │
│ │ User: github-username                                   │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│ Exit Condition: SSH key uploaded to GitHub & registered        │
│                with coordination server                         │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 4: Repository Setup (Fork or Clone)                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│ 4a: Determine Repository Ownership                              │
│ ────────────────────────────────────────────────────────────   │
│ Check: Is repo under current user's GitHub account?            │
│                                                                  │
│ $ gh repo view my-org/my-project --json owner                 │
│                                                                  │
│ Decision:                                                        │
│ ┌──────────────────────┐                                       │
│ │ Owner == current user?                                       │
│ └───────┬──────────────┘                                       │
│         │                                                       │
│    YES  │  NO                                                  │
│         │                                                       │
│         ▼  ▼                                                    │
│     4b   4b'                                                   │
│                                                                  │
│ 4b: User Owns Repo → Clone                                     │
│ ────────────────────────────────────────────────────────────  │
│ - Use git clone (SSH) to user's fork/repo                     │
│ - UX: "Cloning github-username/my-project..."                │
│ - Continue to 4c                                               │
│                                                                  │
│ 4b': User Doesn't Own Repo → Fork Then Clone                   │
│ ──────────────────────────────────────────────────────────────│
│ - Check: Is this a GitHub org repo without a fork?            │
│ - If yes, fork using GitHub CLI:                              │
│   $ gh repo fork my-org/my-project --clone                    │
│ - Clones to ~/my-project and sets up remotes (origin, upstream)│
│ - UX: "Forking my-org/my-project..."                          │
│ - UX: " Forked to github-username/my-project"               │
│ - UX: "Cloning to local disk..."                              │
│ - Continue to 4c                                               │
│                                                                  │
│ 4c: Create Workspace on Staging Server                         │
│ ────────────────────────────────────────────────────────────  │
│ HTTP POST to staging-server:3001/api/v1/workspaces/create     │
│                                                                  │
│ Request Body:                                                   │
│ {                                                               │
│   "github_username": "github-username",                        │
│   "workspace_name": "my-project-main",                         │
│   "repo": {                                                     │
│     "owner": "my-org" OR "github-username" (if forked),       │
│     "name": "my-project",                                      │
│     "url": "git@github.com:...",                               │
│     "branch": "main"                                           │
│   },                                                            │
│   "provider": "lxc",  # auto-selected for staging              │
│   "image": "ubuntu:22.04"  # or user-configured               │
│ }                                                               │
│                                                                  │
│ Response:                                                        │
│ {                                                               │
│   "workspace_id": "ws-abc123",                                 │
│   "ssh_port": 2222,                                            │
│   "status": "creating",                                        │
│   "estimated_time": "30s"                                      │
│ }                                                               │
│                                                                  │
│ UX Output:                                                       │
│ ┌──────────────────────────────────────────────────────────┐   │
│ │ 📦 Repository Setup                                      │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ │ Checking repository ownership...                       │   │
│ │ Repo: my-org/my-project                                  │   │
│ │ Owner: my-org (not your account)                         │   │
│ │                                                          │   │
│ │ Forking repository...                              ⏳    │   │
│ │  Forked to github-username/my-project                 │   │
│ │                                                          │   │
│ │ Cloning to local disk...                           ⏳    │   │
│ │  Cloned to ~/my-project                                │   │
│ │                                                          │   │
│ │ Creating workspace on staging server...            ⏳    │   │
│ │ Server: staging.example.com                             │   │
│ │ Provider: lxc                                           │   │
│ │ Image: ubuntu:22.04                                     │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│ Exit Condition: Workspace created on server (status: creating) │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 5: Workspace Initialization (Container Creation)           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│ Server-Side (Coordination Server + LXC Driver):                 │
│ ──────────────────────────────────────────────────────────────│
│                                                                  │
│ 5a: LXC Container Creation                                     │
│ - Launch LXC container from image (ubuntu:22.04)               │
│ - Container name: workspace-<id> (e.g., workspace-abc123)      │
│ - Network: default LXC bridge                                  │
│ - UX (user-facing): "Initializing container... ⏳"              │
│ - Expected time: 20-30 seconds                                 │
│                                                                  │
│ 5b: SSH Server Setup Inside Container                          │
│ - Install openssh-server                                        │
│ - Generate host keys                                            │
│ - Create 'dev' user (if not exists)                            │
│ - Add user's SSH pubkey to /home/dev/.ssh/authorized_keys      │
│ - Start SSH server (systemctl start ssh)                        │
│                                                                  │
│ 5c: Repository Cloning Inside Container                        │
│ - Clone repo from GitHub into /home/dev/workspace              │
│ - Use SSH clone (git@github.com:...)                           │
│ - Check: Does .nexus/config.yaml exist in repo?               │
│   - If yes: Load it                                            │
│   - If no: Use defaults (node:20, basic npm dev script)        │
│                                                                  │
│ 5d: Service Definition & Environment Setup                     │
│ - Parse services from config (or defaults)                     │
│ - Install dependencies (npm install, pip install, etc.)        │
│ - Prepare service startup scripts                              │
│ - UX (user-facing): "Installing dependencies... ⏳"             │
│                                                                  │
│ 5e: SSH Port Mapping                                           │
│ - Get container IP address (e.g., 10.0.0.42)                  │
│ - Set up SSH forwarding: staging-server:2222 → container:22    │
│ - Verify connectivity test                                     │
│ - Update workspace metadata with SSH port (2222)               │
│ - UX (user-facing): "SSH access ready"                         │
│                                                                  │
│ Exit Condition: Container running, SSH ready, services staged  │
│                                                                  │
│ User-Facing UX:                                                 │
│ ┌──────────────────────────────────────────────────────────┐   │
│ │ 🔨 Workspace Initialization                             │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ │ Creating LXC container...                              ⏳    │   │
│ │  Container created: workspace-abc123                   │   │
│ │                                                          │   │
│ │ Setting up SSH...                                   ⏳    │   │
│ │  SSH server ready (dev@staging:2222)                  │   │
│ │                                                          │   │
│ │ Cloning repository...                               ⏳    │   │
│ │  Cloned to /home/dev/workspace                         │   │
│ │                                                          │   │
│ │ Installing dependencies...                          ⏳    │   │
│ │  npm packages installed                               │   │
│ │                                                          │   │
│ │ Initializing services...                            ⏳    │   │
│ │ Services found:                                        │   │
│ │   - web (npm run dev, port 3000)                       │   │
│ │   - api (npm run server, port 4000)                    │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ └──────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 6: Start Services & Port Mapping                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│ 6a: Parse Service Configuration                                │
│ - From .nexus/config.yaml or auto-detect from repo type        │
│ - Build dependency graph                                        │
│ - Example:                                                      │
│   ```yaml                                                       │
│   services:                                                     │
│     db:                                                         │
│       command: "postgres -D /data/postgres"                     │
│       port: 5432                                                │
│       health_check: "psql://localhost/postgres"                 │
│     api:                                                         │
│       command: "npm run server"                                 │
│       port: 4000                                                │
│       depends_on: [db]                                          │
│     web:                                                         │
│       command: "npm run dev"                                    │
│       port: 3000                                                │
│       depends_on: [api]                                         │
│   ```                                                            │
│                                                                  │
│ 6b: Start Services in Dependency Order                          │
│ - Start DB first → wait for health check (30s timeout)         │
│ - Start API → wait for readiness                               │
│ - Start web → wait for readiness                               │
│ - Error handling: If a service fails, halt with clear message  │
│                                                                  │
│ 6c: Port Mapping (Container to Staging Host)                   │
│ - For each service port, allocate host port (3000+)            │
│ - Example:                                                      │
│   - Container service web:3000 → Host 23000                    │
│   - Container service api:4000 → Host 23001                    │
│   - Container service db:5432 → Host 23002 (internal only)     │
│ - Use SSH tunneling or direct iptables for mapping             │
│ - Store in workspace metadata for discovery                    │
│                                                                  │
│ 6d: Health Verification                                        │
│ - Ping each service port (HTTP GET or TCP connect)             │
│ - Timeout: 30s per service                                     │
│ - If all pass → workspace status = "ready"                     │
│ - If any fail → status = "partial" with warnings               │
│                                                                  │
│ UX Output:                                                       │
│ ┌──────────────────────────────────────────────────────────┐   │
│ │  Starting Services                                     │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ │ Starting db...                                      ⏳    │   │
│ │  Database ready (postgres)                            │   │
│ │                                                          │   │
│ │ Starting api (depends on db)...                   ⏳     │   │
│ │  API ready (http://localhost:23001)                  │   │
│ │                                                          │   │
│ │ Starting web (depends on api)...                  ⏳     │   │
│ │  Web ready (http://localhost:23000)                  │   │
│ │                                                          │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ │  All services running and healthy                     │   │
│ └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│ Exit Condition: All services healthy, ports mapped             │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 7: Open Editor and Display Summary                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│ 7a: Detect Editor Preference                                   │
│ - Check (in order):                                             │
│   1. Environment variable: $NEXUS_EDITOR                       │
│   2. Command: cursor --version (preferred)                     │
│   3. Command: code --version (VS Code)                         │
│   4. Command: vim (fallback)                                   │
│ - Store choice in ~/.nexus/preferences.yaml                    │
│                                                                  │
│ 7b: Generate Remote Editor Deep Link                           │
│ Example for Cursor:                                             │
│ cursor://ssh/remote?host=staging.example.com&port=2222&user=dev&path=/home/dev/workspace
│                                                                  │
│ Example for VS Code:                                            │
│ vscode://vscode-remote/ssh-remote+dev@staging.example.com:2222/home/dev/workspace
│                                                                  │
│ 7c: Open Editor                                                │
│ - Execute: cursor --remote ssh-remote+dev@staging:2222 /home/dev/workspace
│ - Or: code --remote ssh-remote+dev@staging:2222 /home/dev/workspace
│ - Or: ssh dev@staging.example.com -p 2222 (if vim)            │
│ - Non-blocking: Shows deep link in case editor doesn't open   │
│                                                                  │
│ 7d: Display Summary & Connection Info                          │
│ ┌──────────────────────────────────────────────────────────┐   │
│ │                                                          │   │
│ │  WORKSPACE READY                                       │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ │                                                          │   │
│ │ 📍 Project: my-project (github-username/my-project)    │   │
│ │ 🖥  Environment: LXC Container (workspace-abc123)       │   │
│ │  Status: Running                                       │   │
│ │                                                          │   │
│ │ 🌐 Services:                                            │   │
│ │    • web     → http://localhost:23000 (npm dev)        │   │
│ │    • api     → http://localhost:23001 (npm server)     │   │
│ │    • db      → localhost:23002 (postgresql, internal)  │   │
│ │                                                          │   │
│ │ 🐚 SSH Access:                                          │   │
│ │    Command: ssh -p 2222 dev@staging.example.com        │   │
│ │    Or use editor deep link (opening Cursor now...)     │   │
│ │                                                          │   │
│ │ 💻 Editor Deep Link:                                   │   │
│ │    cursor://ssh/remote?host=staging.example.com&port=2222... │
│ │                                                          │   │
│ │  Git Status:                                          │   │
│ │    Branch: main                                         │   │
│ │    Upstream: my-org/my-project (original)              │   │
│ │    Origin: github-username/my-project (your fork)      │   │
│ │                                                          │   │
│ │ 🔧 Next Steps:                                         │   │
│ │    1. Editor opens automatically → start designing!    │   │
│ │    2. Make changes, test with running services         │   │
│ │    3. Commit: git commit -m "..."                      │   │
│ │    4. Create PR: gh pr create --base my-org/main \     │   │
│ │       --head github-username/my-project                │   │
│ │                                                          │   │
│ │  More commands:                                      │   │
│ │    • nexus workspace services        # List ports      │   │
│ │    • nexus workspace logs <service> # View output     │   │
│ │    • nexus workspace exec <cmd>     # Run in container│   │
│ │    • nexus workspace down            # Stop workspace  │   │
│ │                                                          │   │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │   │
│ │ ready to work!                                      │   │
│ └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│ Exit Condition: Editor opened, summary displayed               │
└─────────────────────────────────────────────────────────────────┘
```

---

### 2.2 Decision Trees & Edge Cases

#### Edge Case 2.2a: User has SSH key but not on GitHub

**Scenario**: User has ~/.ssh/id_ed25519.pub locally, but hasn't uploaded to GitHub yet.

**Flow**:
```
Step 3c: Upload to GitHub
  ↓
API Response 409 Conflict (key already registered)
  ↓
UX: "SSH key already exists locally. Registering with GitHub..."
  ↓
gh ssh-key add ~/.ssh/id_ed25519.pub --title "nexus@staging"
  ↓
Continue to 3d (register with coordination server)
```

#### Edge Case 2.2b: User's SSH key exists on GitHub but under different email

**Scenario**: User has ~/.ssh/id_rsa, uploaded long ago with different GitHub account.

**Flow**:
```
Step 3: SSH Key Setup
  ↓
Detect key exists (id_rsa)
  ↓
UX: "Found existing key. Verify it's registered with current GitHub account?"
  ↓
User confirms or generates new key
  ↓
Continue with either key
```

#### Edge Case 2.2c: Repository is in private org

**Scenario**: User forking my-org/private-repo where they're already a member.

**Flow**:
```
Step 4: Repository Setup
  ↓
Check GitHub permissions for my-org/private-repo
  ↓
If user is member: Allow fork/clone
  ↓
If user is NOT member: Error "You don't have access to this repository"
  ↓
UX: "Request access at https://github.com/my-org/private-repo/settings/access"
```

#### Edge Case 2.2d: LXC/Docker not available on staging server

**Scenario**: LXC daemon not running or container image missing.

**Flow**:
```
Step 5a: LXC Container Creation
  ↓
lxc launch ubuntu:22.04 workspace-abc123 (fails)
  ↓
Error: "LXC daemon not responding. Check: sudo systemctl status lxd"
  ↓
UX: Clear message with troubleshooting steps
```

---

## Part 3: Technical Implementation Details

### 3.1 API Contracts

#### Register User with Coordination Server

```http
POST /api/v1/users/register
Content-Type: application/json

{
  "github_username": "github-username",
  "github_id": 123456789,
  "ssh_pubkey": "ssh-ed25519 AAAA...",
  "ssh_pubkey_fingerprint": "SHA256:..."
}

Response: 200 OK
{
  "user_id": "user-abc123",
  "github_username": "github-username",
  "registered_at": "2026-01-17T10:30:00Z",
  "workspaces": []
}
```

#### Create Workspace

```http
POST /api/v1/workspaces/create
Content-Type: application/json
Authorization: Bearer <coordination-server-token> (optional)

{
  "github_username": "github-username",
  "workspace_name": "my-project-main",
  "repo": {
    "owner": "my-org",
    "name": "my-project",
    "url": "git@github.com:my-org/my-project.git",
    "branch": "main"
  },
  "provider": "lxc",
  "image": "ubuntu:22.04"
}

Response: 202 Accepted
{
  "workspace_id": "ws-abc123",
  "status": "creating",
  "ssh_port": 2222,
  "estimated_time": "30s",
  "polling_url": "/api/v1/workspaces/ws-abc123/status"
}
```

#### Get Workspace Status

```http
GET /api/v1/workspaces/ws-abc123/status

Response: 200 OK
{
  "workspace_id": "ws-abc123",
  "status": "running",  # pending, creating, running, stopped, error
  "ssh": {
    "host": "staging.example.com",
    "port": 2222,
    "user": "dev"
  },
  "services": {
    "web": {
      "status": "running",
      "port": 3000,
      "mapped_port": 23000,
      "url": "http://localhost:23000"
    },
    "api": {
      "status": "running",
      "port": 4000,
      "mapped_port": 23001,
      "url": "http://localhost:23001"
    }
  },
  "repo": {
    "url": "git@github.com:github-username/my-project.git",
    "branch": "main",
    "commit": "abc123def456"
  }
}
```

### 3.2 Configuration File Structure

#### .nexus/config.yaml (In Repository)

```yaml
version: "1.0"
name: my-project

# Service definitions
services:
  web:
    command: "npm run dev"
    port: 3000
    health_check:
      type: "http"
      path: "/"
      timeout: 10s

  api:
    command: "npm run server"
    port: 4000
    depends_on: ["db"]
    health_check:
      type: "tcp"
      timeout: 5s

  db:
    command: "postgres -D /data/postgres"
    port: 5432
    environment:
      POSTGRES_DB: "dev"
      POSTGRES_USER: "dev"
      POSTGRES_PASSWORD: "password"
    health_check:
      type: "custom"
      command: "psql -U dev -d dev -c 'SELECT 1'"

# Provider configuration
providers:
  default: "lxc"
  
  lxc:
    image: "ubuntu:22.04"
    cpu: 2
    memory: "4GB"
    disk: "20GB"
    
  docker:
    image: "node:20-alpine"
    
  qemu:
    image: "ubuntu:22.04"
    cpu: 4
    memory: "8GB"
    disk: "50GB"

# Optional: lifecycle hooks
lifecycle:
  pre_start: |
    #!/bin/bash
    npm install
  post_start: |
    #!/bin/bash
    npm run setup
  pre_stop: |
    #!/bin/bash
    npm run cleanup
```

#### ~/.nexus/config.yaml (User Config - Generated)

```yaml
version: "1.0"
server: "staging.example.com"
server_port: 3001

github:
  username: "github-username"
  auth_token: "<github-cli-token>"

ssh:
  key_path: "~/.ssh/id_ed25519"
  key_type: "ed25519"

editor:
  preferred: "cursor"  # cursor, code, vim, neovim

workspaces:
  my-project-main:
    workspace_id: "ws-abc123"
    repo: "github-username/my-project"
    provider: "lxc"
    status: "running"
    ssh_port: 2222
    created_at: "2026-01-17T10:30:00Z"
```

### 3.3 SSH Key Management Flow

**Key Storage Locations**:
- Local: `~/.ssh/id_ed25519` (private key)
- Local: `~/.ssh/id_ed25519.pub` (public key)
- GitHub: Uploaded via `gh ssh-key add`
- Coordination Server: Fingerprint stored for user registration
- Container: Copied to `/home/dev/.ssh/authorized_keys`

**Key Upload Chain**:
```
1. Generate or detect local key (~/.ssh/id_ed25519.pub)
2. Upload to GitHub via `gh ssh-key add`
3. Register with coordination server (store fingerprint + public key)
4. Container init: Adds key to authorized_keys
5. User SSH: Uses local private key to connect
```

---

## Part 4: One-Line Install Script Specification

### 4.1 Script Location & Download

```bash
# Primary
curl https://nexus.example.com/install.sh | bash -s -- --repo OWNER/REPO --server staging.example.com

# Or with arguments
curl https://staging.example.com/api/v1/install-script | bash
```

### 4.2 Script Behavior

```bash
#!/bin/bash
set -e

# Parse arguments
REPO=""
SERVER=""
EDITOR=""

while [[ $# -gt 0 ]]; do
  case $1 in
    --repo)
      REPO="$2"
      shift 2
      ;;
    --server)
      SERVER="$2"
      shift 2
      ;;
    --editor)
      EDITOR="$2"
      shift 2
      ;;
    *)
      shift
      ;;
  esac
done

# Validate
if [[ -z "$SERVER" ]]; then
  echo "Error: --server required"
  exit 1
fi

# Download binary
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
BINARY_URL="https://releases.nexus.example.com/nexus-latest-${OS}-${ARCH}.tar.gz"

echo "📦 Downloading nexus CLI..."
curl -fsSL "$BINARY_URL" | tar xz -C ~/.local/bin/

# Initialize config
echo "⚙  Initializing configuration..."
mkdir -p ~/.nexus
cat > ~/.nexus/config.yaml << EOF
server: "$SERVER"
server_port: 3001
EOF

# Check dependencies
echo " Checking system requirements..."
command -v git >/dev/null || (echo "Error: git not found" && exit 1)
command -v ssh >/dev/null || (echo "Error: ssh not found" && exit 1)

# Check GitHub CLI
if ! command -v gh &> /dev/null; then
  echo "📦 Installing GitHub CLI..."
  if [[ "$OS" == "darwin" ]]; then
    brew install gh
  else
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
    # ... additional setup
  fi
fi

# Start setup flow
echo " Starting nexus workspace setup..."
nexus workspace setup-from-repo "$REPO" --server "$SERVER" ${EDITOR:+--editor "$EDITOR"}
```

### 4.3 Success Criteria

After script completes:
1.  nexus CLI binary in PATH
2.  ~/.nexus/config.yaml created with server address
3.  GitHub CLI installed and checked
4.  SSH keys checked/generated
5.  User directed to GitHub authentication
6.  First workspace creation started

---

## Part 5: Implementation Checklist

### Phase 1: Coordination Server & Infrastructure

- [ ] **Coordination Server Core**
  - [ ] HTTP API server (port 3001)
  - [ ] User registration endpoint (`/api/v1/users/register`)
  - [ ] Workspace CRUD endpoints
  - [ ] Node agent health check (`/health`, `/metrics`)
  - [ ] SSH port forwarding setup & management

- [ ] **Node Agent**
  - [ ] Heartbeat to coordination server
  - [ ] Command reception interface
  - [ ] Provider dispatch logic
  - [ ] Status reporting back to server

- [ ] **LXC Driver Integration**
  - [ ] Remote container creation via agent
  - [ ] SSH server setup in container
  - [ ] Port mapping management
  - [ ] Service lifecycle management

- [ ] **Testing Infrastructure**
  - [ ] Staging host setup (Ubuntu 22.04 + LXC)
  - [ ] Network connectivity verification
  - [ ] Basic sanity tests (start/stop container)

### Phase 2: CLI User Flow

- [ ] **One-Line Install Script**
  - [ ] Binary download & install
  - [ ] Dependency checking (git, ssh, gh)
  - [ ] Configuration initialization
  - [ ] Entry point to setup flow

- [ ] **GitHub Authentication**
  - [ ] Detect/install GitHub CLI
  - [ ] Orchestrate `gh auth login`
  - [ ] Store username & tokens securely

- [ ] **SSH Key Management**
  - [ ] Detect existing keys
  - [ ] Generate new keys if needed
  - [ ] Upload to GitHub via `gh ssh-key add`
  - [ ] Register fingerprint with coordination server

- [ ] **Repository Handling**
  - [ ] Check repo ownership (owned vs. org)
  - [ ] Fork if necessary via `gh repo fork`
  - [ ] Clone to local disk
  - [ ] Create workspace on coordination server

- [ ] **Workspace Initialization**
  - [ ] Container creation with proper image
  - [ ] SSH configuration in container
  - [ ] Repository cloning in container
  - [ ] Dependency installation (npm/pip/etc)
  - [ ] Service startup & health checking

- [ ] **Editor Integration**
  - [ ] Detect editor preference
  - [ ] Generate deep links (Cursor, VS Code)
  - [ ] Launch editor with SSH remote connection
  - [ ] Display summary & next steps

### Phase 3: Polish & Testing

- [ ] **Error Handling**
  - [ ] Clear error messages for all failure paths
  - [ ] Actionable guidance for troubleshooting
  - [ ] Graceful degradation (fallbacks)

- [ ] **UX Refinement**
  - [ ] Progress indicators (spinners, %) for long operations
  - [ ] Colored output for clarity
  - [ ] Summary display with quick copy commands
  - [ ] Contextual help & links

- [ ] **E2E Testing**
  - [ ] Full flow: install → auth → workspace ready
  - [ ] Multi-user scenario (concurrent workspaces)
  - [ ] Different repo types (Node, Python, monorepo)
  - [ ] Error scenarios (network failure, etc.)

- [ ] **Documentation**
  - [ ] User guide for first-time setup
  - [ ] Troubleshooting guide
  - [ ] Architecture diagrams
  - [ ] API documentation

---

## Part 6: Success Metrics

### User Experience

- **First Workspace Time**: < 3 minutes from install script start
  - Download & install: 20s
  - GitHub auth: 30s
  - SSH setup: 15s
  - Repo fork/clone: 20s
  - Container creation: 50s
  - Service startup: 45s
  - Total: ~3 minutes

- **Error Recovery**: All errors with clear next steps
  - No cryptic messages
  - All URLs/commands clickable or copy-able
  - Actionable guidance for every failure

- **One-Click Editor Launch**: Deep links work for Cursor, VS Code
  - Editor opens automatically
  - SSH connection established
  - User can start editing immediately

### Technical

- **Coordination Server Stability**: 99.9% uptime
- **SSH Port Mapping**: < 1s latency
- **Service Health Checks**: 100% detection of failed services
- **Container Isolation**: No cross-workspace interference

---

## Part 7: Rollout Plan

### Staging Phase 1: Internal Testing (Week 1)
- Deploy coordination server + LXC node
- Manual end-to-end testing
- Fix critical bugs

### Staging Phase 2: Soft Launch (Week 2-3)
- 5-10 beta users
- Collect feedback
- Iterate on UX

### Staging Phase 3: General Availability (Week 4+)
- Public documentation
- Marketing launch
- Scale infrastructure as needed

---

## Appendix: CLI Command Examples

### User Commands

```bash
# Install
curl https://nexus.example.com/install.sh | bash -s -- \
  --repo my-org/my-project --server staging.example.com

# After setup, these commands are available:
nexus workspace list                # List all workspaces
nexus workspace create <name>       # Create new workspace
nexus workspace up <name>           # Start workspace
nexus workspace down <name>         # Stop workspace
nexus workspace services <name>     # List running services
nexus workspace logs <name> <svc>   # View service logs
nexus workspace exec <name> <cmd>   # Run command in workspace
nexus workspace connect <name>      # Show SSH/editor links
nexus workspace rm <name>           # Delete workspace
```

---

## Appendix: Environment Variables

```bash
# Configuration
NEXUS_SERVER=staging.example.com
NEXUS_SERVER_PORT=3001

# GitHub
GITHUB_TOKEN=ghp_...  # Optional, for faster auth

# SSH
NEXUS_SSH_KEY_PATH=~/.ssh/id_ed25519

# Editor
NEXUS_EDITOR=cursor  # Default: cursor, fallback: code → vim

# Debugging
NEXUS_DEBUG=1         # Enable debug output
NEXUS_TRACE=1         # Enable trace output
```

---

**Document Status**: Initial Specification (Ready for Review)  
**Next Steps**: 
1. Review & refine with team
2. Create detailed API specs (OpenAPI/Swagger)
3. Begin Phase 1 implementation (coordination server)
4. Set up staging infrastructure
