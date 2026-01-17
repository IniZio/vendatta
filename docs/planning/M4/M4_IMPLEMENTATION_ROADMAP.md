# M3 Implementation Plan: Staging Environment & User Flow

**Status**: Ready for Implementation  
**Priority**: Critical Path (Blocks general availability)  
**Timeline**: 4-6 weeks to completion  
**Target**: Hyper-focused staging environment with end-to-end user flow  

---

## Executive Summary

After M3, the immediate focus is establishing a **stable staging environment** that enables the complete user flow:

> **Designer's Journey**: Install CLI (1-line) → Authenticate GitHub → Generate SSH keys → Create workspace (fork/clone) → Services start → Editor opens → Start designing

**Key Constraint**: Everything must be automated after the initial install script. fully automated for SSH, keys, workspace setup, or service discovery.

---

## Current State (As of Jan 17, 2026)

###  What's Working
- **Docker, LXC, QEMU Providers**: Local implementation complete
- **Configuration System**: YAML-based, flexible
- **Agent Integration**: Auto-config generation for Cursor, OpenCode, Claude
- **Coordination Server Infrastructure**: Basic structure in place
- **CLI Commands**: Branch management commands exist

### ❌ Critical Gaps
1. **Coordination Server**: No workspace lifecycle management
2. **Node Agents**: No remote command execution
3. **SSH Automation**: Key generation exists, but not upload/GitHub integration
4. **GitHub CLI Integration**: No orchestration of `gh auth`, `gh repo fork`, `gh ssh-key add`
5. **One-Line Install Script**: Not implemented
6. **Editor Launch**: No deep link generation or editor detection

---

## Architecture (Staging Environment)

### Infrastructure Components

```
┌─ Staging Host ─────────────────────────────────────┐
│                                                    │
│  Coordination Server (Port 3001)                   │
│  • Workspace CRUD                                  │
│  • User registration                              │
│  • SSH port forwarding                            │
│  • Service discovery                              │
│                                                    │
│  ┌─ LXC Driver Node ──────────────────────────┐  │
│  │ Agent + LXC daemon                         │  │
│  │ → Receives commands from coordination      │  │
│  │ → Creates containers                       │  │
│  │ → Reports status back                      │  │
│  └────────────────────────────────────────────┘  │
│                                                    │
└────────────────────────────────────────────────────┘
              ↑ SSH: ports 2222-2299
              │
      Designer's Local Machine
      nexus CLI + Editor (Cursor/Code)
```

### Key Data Flows

1. **User Registration**:
   ```
   CLI → GitHub Auth → SSH Key Gen/Upload → Coordination Server
   ```

2. **Workspace Creation**:
   ```
   CLI (fork/clone repo) → Coordination Server 
   → LXC Driver Node (create container) → Setup SSH & services
   ```

3. **Editor Connection**:
   ```
   CLI (detect editor) → Generate SSH deep link 
   → Launch editor → Editor connects via SSH (port 2222)
   → User can edit files in container
   ```

---

## User Flow (Detailed Steps)

### Step 1: One-Line Install
```bash
curl https://nexus.example.com/install.sh | bash -s -- \
  --repo my-org/my-project --server staging.example.com
```

**Script Does**:
- Download nexus CLI binary
- Check dependencies (git, ssh, gh)
- Initialize ~/.nexus/config.yaml
- Call next step (GitHub auth)

**Success Indicator**: "GitHub authentication needed next"

### Step 2: GitHub Authentication
```bash
nexus workspace setup-from-repo my-org/my-project
→ Check: gh auth status
→ If not authed: `gh auth login` (opens browser)
→ Get username from `gh api user --jq '.login'`
```

**Success Indicator**: GitHub username confirmed

### Step 3: SSH Key Setup
```bash
nexus ssh setup
→ Check ~/.ssh/id_ed25519.pub (exists?)
→ If yes: Ask user to confirm
→ If no: Generate with `ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -N ""`
→ Upload to GitHub: `gh ssh-key add ~/.ssh/id_ed25519.pub --title "nexus@staging"`
→ Register fingerprint with coordination server
```

**Success Indicator**: SSH key registered with GitHub & coordination server

### Step 4: Repository Setup
```bash
nexus workspace create my-project-main
→ Check repo ownership: `gh repo view my-org/my-project --json owner`
→ If owned by user: Clone to ~/my-project
→ If not owned: Fork with `gh repo fork my-org/my-project --clone`
→ Send workspace creation request to coordination server
```

**Request**:
```json
POST /api/v1/workspaces/create
{
  "github_username": "user",
  "workspace_name": "my-project-main",
  "repo": {"owner": "my-org", "name": "my-project", "branch": "main"},
  "provider": "lxc",
  "image": "ubuntu:22.04"
}
```

**Response**: `{"workspace_id": "ws-123", "status": "creating", "ssh_port": 2222}`

### Step 5: Container Initialization
**Coordination Server → LXC Driver Node**:
- Create LXC container from image
- Add SSH key to authorized_keys
- Clone repo inside container
- Install dependencies
- Start services with dependency ordering

**Polling**: Client polls `/api/v1/workspaces/ws-123/status` until `status: running`

### Step 6: Editor Launch
```bash
nexus workspace connect my-project-main
→ Detect editor: cursor > code > vim
→ Generate deep link: 
   cursor://ssh/remote?host=staging.example.com&port=2222&user=dev&path=/home/dev/workspace
→ Execute: cursor --remote ssh-remote+dev@staging:2222 /home/dev/workspace
→ Display summary with service URLs
```

---

## Implementation Phases

### Phase 1: Coordination Server Foundation (Week 1-2)

**Goal**: Workspace lifecycle management + LXC driver integration

**Key Components**:
1. HTTP API server (port 3001)
   - `/api/v1/users/register` - Register GitHub user + SSH key
   - `/api/v1/workspaces/create` - Create workspace
   - `/api/v1/workspaces/{id}/status` - Get status
   - `/api/v1/workspaces/{id}/stop` - Stop workspace
   - `/health` - Health check for driver nodes

2. Workspace CRUD Logic
   - Metadata storage (in-memory or SQLite for staging)
   - SSH port allocation (2222-2299 pool)
   - Status tracking (pending, creating, running, stopped)

3. Driver Node Communication
   - Heartbeat/health check
   - Command dispatch (create container, start services)
   - Status reporting back

4. Node Agent (runs on LXC node)
   - Listens for commands from coordination server
   - Executes provider operations (lxc launch, lxc exec, etc.)
   - Reports container status

**Testing**: Manual end-to-end with staging LXC node

### Phase 2: GitHub CLI Integration (Week 2-3)

**Goal**: Automate GitHub auth, fork, SSH key upload

**Key Commands to Implement**:
1. `nexus auth github`
   - Detect GitHub CLI status
   - Orchestrate `gh auth login` if needed
   - Extract username + store token securely

2. `nexus ssh setup`
   - Detect/generate SSH key
   - Upload to GitHub via `gh ssh-key add`
   - Register fingerprint with coordination server

3. `nexus workspace create <repo>`
   - Parse repo (owner/name/branch)
   - Fork if necessary: `gh repo fork <owner>/<repo> --clone`
   - Clone if already owned
   - Create workspace on coordination server

4. `nexus workspace connect`
   - Detect editor preference
   - Generate SSH deep links
   - Launch editor with remote connection
   - Display summary

**Testing**: Full flow with test GitHub repo

### Phase 3: One-Line Install Script (Week 3)

**Goal**: Single command to bootstrap everything

**Script**:
```bash
curl https://nexus.example.com/install.sh | bash -s -- \
  --repo my-org/my-project --server staging.example.com
```

**Behavior**:
1. Download & install nexus binary
2. Check dependencies (git, ssh, gh)
3. Initialize ~/.nexus/config.yaml
4. Call `nexus workspace setup-from-repo my-org/my-project`
5. Hand off to interactive setup flow

**Platforms**: macOS, Linux (Ubuntu/Debian/RHEL)

### Phase 4: Polish & E2E Testing (Week 4)

**Goal**: Stable, reliable staging environment

**Testing**:
- Full flow: install → auth → workspace → services → editor
- Error scenarios: network issues, missing tools, permission errors
- Multi-user: concurrent workspaces don't interfere
- Service health: all services start reliably
- SSH access: port forwarding works, latency acceptable

**UX Improvements**:
- Clear progress indicators
- Actionable error messages
- Command copy/paste for troubleshooting
- Summary display with quick links

---

## API Specification

### User Registration
```http
POST /api/v1/users/register
Content-Type: application/json

{
  "github_username": "alice",
  "github_id": 123456789,
  "ssh_pubkey": "ssh-ed25519 AAAA...",
  "ssh_pubkey_fingerprint": "SHA256:..."
}

Response: 200 OK
{
  "user_id": "user-abc123",
  "registered_at": "2026-01-17T10:30:00Z",
  "workspaces": []
}
```

### Workspace Creation
```http
POST /api/v1/workspaces/create
Content-Type: application/json

{
  "github_username": "alice",
  "workspace_name": "my-project-feature",
  "repo": {
    "owner": "my-org",
    "name": "my-project",
    "url": "git@github.com:my-org/my-project.git",
    "branch": "main"
  },
  "provider": "lxc",
  "image": "ubuntu:22.04",
  "services": [
    {"name": "web", "command": "npm run dev", "port": 3000},
    {"name": "api", "command": "npm run server", "port": 4000}
  ]
}

Response: 202 Accepted
{
  "workspace_id": "ws-abc123",
  "status": "creating",
  "ssh_port": 2222,
  "polling_url": "/api/v1/workspaces/ws-abc123/status"
}
```

### Workspace Status
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
    "owner": "my-org",
    "name": "my-project",
    "branch": "main",
    "commit": "abc123def456"
  }
}
```

---

## Configuration Files

### .nexus/config.yaml (in repository)
```yaml
name: my-project
provider: lxc

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
    depends_on: ["web"]

providers:
  lxc:
    image: "ubuntu:22.04"
    cpu: 2
    memory: "4GB"
```

### ~/.nexus/config.yaml (user config, generated)
```yaml
server: staging.example.com
server_port: 3001

github:
  username: alice
  auth_token: <gh-token>

ssh:
  key_path: ~/.ssh/id_ed25519
  key_type: ed25519

editor: cursor
```

---

## Success Criteria

### User Experience
-  First workspace in < 3 minutes (from install script)
-  All steps automated after initial install
-  Clear error messages for all failure paths
-  Editor opens automatically with working SSH connection
-  Services discoverable and accessible (URLs printed)

### Technical
-  Coordination server stable (99.9% uptime)
-  SSH port forwarding latency < 100ms
-  Container creation < 30 seconds
-  Service health checks 100% reliable
-  No cross-workspace interference

### Coverage
-  Works with Node.js projects (npm)
-  Works with Python projects (pip/poetry)
-  Works with static sites (no build needed)
-  Handles repos with .nexus/config.yaml
-  Handles repos without config (sensible defaults)

---

## Risks & Mitigations

### Risk: SSH Key Upload Failures (GitHub API)
**Mitigation**: 
- Retry logic with exponential backoff
- Detect duplicate keys (409 conflict)
- Clear error message if key already exists

### Risk: LXC Container Doesn't Start
**Mitigation**:
- Health check before marking ready
- Timeout after 2 minutes with helpful error
- Include LXC daemon status check in setup

### Risk: Service Dependencies Don't Work
**Mitigation**:
- Parse .nexus/config.yaml for depends_on
- Start services in order
- Health checks with timeouts per service
- Mark service as "partial" if some fail

### Risk: SSH Deep Links Don't Work in Editor
**Mitigation**:
- Fallback to printed SSH command: `ssh -p 2222 dev@staging`
- Validate SSH connection before printing ready message
- Provide troubleshooting guide in summary

---

## Implementation Checklist

### Coordination Server
- [ ] HTTP API server on port 3001
- [ ] `/api/v1/users/register` endpoint
- [ ] `/api/v1/workspaces/create` endpoint
- [ ] `/api/v1/workspaces/{id}/status` endpoint
- [ ] SSH port allocation & forwarding
- [ ] Workspace metadata storage (SQLite)
- [ ] Health check endpoint

### Node Agent
- [ ] Agent binary that runs on LXC nodes
- [ ] Heartbeat to coordination server
- [ ] Command reception (create container, start service, etc.)
- [ ] Status reporting back to server
- [ ] Container lifecycle management (create, start, stop)
- [ ] SSH setup inside container
- [ ] Service startup & health checks

### CLI Commands
- [ ] `nexus auth github` - orchestrate gh auth
- [ ] `nexus ssh setup` - key generation & GitHub upload
- [ ] `nexus workspace create <repo>` - fork/clone + create
- [ ] `nexus workspace connect` - editor launch & summary
- [ ] `nexus workspace status` - poll creation progress
- [ ] `nexus workspace services` - list running services

### One-Line Install
- [ ] Download script (curl https://...install.sh)
- [ ] Binary download & install
- [ ] Dependency checking
- [ ] Config initialization
- [ ] Platform support (macOS, Linux)

### Testing
- [ ] E2E test: install → workspace → services running
- [ ] E2E test: error scenarios (network, permissions)
- [ ] E2E test: multi-user concurrent workspaces
- [ ] Manual testing with designers/developers
- [ ] Load testing with 10+ concurrent workspaces

---

## Dependencies & Prerequisites

### For Staging Server
- Ubuntu 22.04+
- LXC installed (`lxd init`)
- Go 1.24+ for building coordination server & agent
- SQLite for metadata storage (embedded)

### For Developers
- Git
- SSH client
- GitHub CLI (`gh`)
- Cursor or VS Code (recommended)

### External Services
- GitHub.com (API, OAuth)
- DNS (for staging server address)

---

## Success Metrics

| Metric | Target | How to Measure |
|--------|--------|----------------|
| First workspace time | < 3 min | Time from install script to editor opening |
| Service startup | < 45 sec | Time from container creation to services healthy |
| SSH latency | < 100ms | `time ssh ... 'echo hi'` |
| Error clarity | 100% | All errors have actionable next steps |
| Editor launch success | 95% | Deep links work for Cursor/VS Code |
| Uptime | 99.9% | Coordination server availability |

---

## Rollout Plan

### Week 1: Internal Testing
- Deploy coordination server + LXC node on staging host
- Manual E2E testing with team
- Fix critical bugs

### Week 2: Soft Launch
- 5-10 beta testers (designers)
- Gather feedback on UX
- Iterate quickly on issues

### Week 3: Polish
- Refine error messages based on feedback
- Performance optimization
- Documentation

### Week 4: General Availability
- Public announcement
- Marketing content
- Support runbooks

---

## Document References

- **Full User Flow Spec**: `M3_STAGING_ENV_USER_FLOW_SPEC.md` (in this directory)
- **Coordination Server API**: See API Specification section above
- **Installation Script**: To be created in `scripts/install.sh`
- **Node Agent**: To be created in `cmd/nexus-agent/main.go`

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-17  
**Status**: Ready for Implementation  
**Next Action**: Begin Phase 1 (Coordination Server) implementation
