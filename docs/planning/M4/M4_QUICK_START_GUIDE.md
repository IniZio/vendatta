# M3 Post-Completion: Staging Environment & User Flow Specification

**Executive Summary for Stakeholders**

---

## The Vision

After completing M3 (Provider-Agnostic Remote Nodes), we hyper-focus on one thing:

> **A stable staging environment that makes it dead simple for designers to start designing.**

One install command → GitHub login → Project ready → Start designing → Create PR

---

## What We're Building

### Staging Infrastructure

```
Staging Host (Ubuntu 22.04)
├── Coordination Server (HTTP API, port 3001)
│   ├── Manages workspace lifecycle
│   ├── Assigns SSH ports (2222-2299)
│   ├── Discovers services & ports
│   └── Forwards SSH connections
│
└── LXC Driver Node
    ├── Runs containers
    ├── Sets up SSH in containers
    ├── Starts services
    └── Reports status back to coordination server
```

### User Journey

```
1. curl install.sh
   ↓
2. GitHub login (opens browser)
   ↓
3. SSH key setup (auto-generated, uploaded to GitHub)
   ↓
4. Fork repo (if needed)
   ↓
5. Create workspace (coordination server creates container)
   ↓
6. Services start (npm/python/etc)
   ↓
7. Editor opens (Cursor with SSH remote connection)
   ↓
8. Designer makes changes → git commit → gh pr create
```

**Total Time**: ~3 minutes from start to editor opening

---

## Key Deliverables

### 1. Coordination Server
- **Purpose**: Central workspace management
- **Functionality**: 
  - User registration (GitHub handle + SSH key)
  - Workspace CRUD (create, start, stop, delete)
  - SSH port forwarding (2222 → container:22, etc.)
  - Service discovery (list running services + their ports)
  - Status tracking (pending, creating, running, stopped)
- **API**: REST HTTP on port 3001
- **Deployment**: Single binary, runs on staging host

### 2. Node Agent
- **Purpose**: Remote command executor
- **Functionality**:
  - Runs on LXC driver node
  - Receives commands from coordination server
  - Executes provider operations (lxc launch, lxc exec, etc.)
  - Reports container status back to coordination server
- **Deployment**: Single binary per driver node

### 3. CLI Extensions
- **`nexus auth github`**: Orchestrate GitHub authentication
- **`nexus ssh setup`**: Generate SSH key, upload to GitHub
- **`nexus workspace create <repo>`**: Fork/clone repo, create workspace
- **`nexus workspace connect`**: Launch editor with SSH remote connection
- **`nexus workspace services`**: List running services & their URLs
- **`nexus workspace status`**: Poll creation progress

### 4. One-Line Install Script
```bash
curl https://nexus.example.com/install.sh | bash -s -- \
  --repo my-org/my-project --server staging.example.com
```

**What it does**:
- Downloads & installs nexus CLI
- Checks dependencies (git, ssh, gh)
- Initializes config
- Kicks off GitHub authentication
- Starts workspace creation flow

---

## Technical Specification (Quick Reference)

### Coordination Server API Endpoints

```
POST   /api/v1/users/register
       Body: {github_username, github_id, ssh_pubkey, ssh_pubkey_fingerprint}

POST   /api/v1/workspaces/create
       Body: {github_username, workspace_name, repo, provider, image, services}
       Response: {workspace_id, status, ssh_port, polling_url}

GET    /api/v1/workspaces/{id}/status
       Response: {status, ssh, services[], repo}

POST   /api/v1/workspaces/{id}/stop
       Response: {status: "stopped"}

DELETE /api/v1/workspaces/{id}

GET    /health
       Response: {status: "ok"}
```

### Data Model

**Workspace Metadata**:
```
workspace:
  id: "ws-abc123"
  owner: "github-username"
  name: "my-project-feature"
  provider: "lxc"
  image: "ubuntu:22.04"
  status: "running"  # pending, creating, running, stopped
  ssh:
    port: 2222
    user: "dev"
  services:
    web: {port: 3000, mapped_port: 23000, status: "running"}
    api: {port: 4000, mapped_port: 23001, status: "running"}
  repo:
    owner: "my-org"
    name: "my-project"
    branch: "main"
    commit: "abc123..."
```

### Configuration Files

**`.nexus/config.yaml`** (in repo):
```yaml
name: my-project
provider: lxc

services:
  web:
    command: "npm run dev"
    port: 3000
    health_check: {type: "http", path: "/"}
  
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

**`~/.nexus/config.yaml`** (user config, auto-generated):
```yaml
server: staging.example.com
server_port: 3001
github:
  username: alice
  auth_token: <token>
ssh:
  key_path: ~/.ssh/id_ed25519
  key_type: ed25519
editor: cursor
```

---

## Implementation Timeline

### Phase 1: Coordination Server (Week 1-2)
- [ ] HTTP API server
- [ ] Workspace CRUD logic
- [ ] Driver node communication
- [ ] SSH port allocation & forwarding
- [ ] Metadata storage (SQLite)

### Phase 2: GitHub CLI Integration (Week 2-3)
- [ ] `nexus auth github` command
- [ ] `nexus ssh setup` command
- [ ] `nexus workspace create` (fork/clone + create)
- [ ] `nexus workspace connect` (editor launch)
- [ ] Full integration testing

### Phase 3: One-Line Install (Week 3)
- [ ] Install script (download binary, check deps)
- [ ] Platform support (macOS, Linux)
- [ ] Dependency installation (gh, git, ssh)

### Phase 4: Polish & Testing (Week 4)
- [ ] E2E testing (full flow start-to-finish)
- [ ] Error scenarios (network, permissions, etc.)
- [ ] Multi-user testing (concurrent workspaces)
- [ ] UX refinement based on feedback

---

## Success Criteria

### User Experience
-  Complete setup in < 3 minutes
-  All steps automated (no manual SSH key copying)
-  Clear error messages with next steps
-  Editor opens automatically with working connection
-  Service URLs discoverable & clickable

### Technical
-  Coordination server: 99.9% uptime
-  Container startup: < 30 seconds
-  Service health checks: 100% reliable
-  SSH latency: < 100ms
-  No cross-workspace interference

### Coverage
-  Node.js projects (npm)
-  Python projects (pip/poetry)
-  Static sites
-  Projects with .nexus/config.yaml
-  Projects without config (sensible defaults)

---

## Why This Matters

### Current State
- Users can use Nexus locally or on remote QEMU VMs
- Manual setup required for SSH, keys, workspace creation
- Designer experience is "technical"

### After This Work
- One-command setup from install to editing
- Fully automated GitHub integration
- Services discovered & accessible
- Editor opens automatically
- Designer experience is "effective"

### Impact
- Removes all friction for onboarding designers/developers
- Makes staging environment the "go-to" for team development
- Enables rapid iteration on design work
- Creates a clear, reproducible workflow

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Staging Host                              │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ Coordination Server (Port 3001)                        │ │
│  │ • Workspace API                                        │ │
│  │ • User registration (GitHub + SSH key)                │ │
│  │ • Port allocation & forwarding                        │ │
│  │ • Service discovery                                   │ │
│  └────────────────────────────────────────────────────────┘ │
│                         ↓ commands                           │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ LXC Driver Node                                        │ │
│  │ ┌──────────────────────────────────────────────────┐  │ │
│  │ │ Node Agent                                       │  │ │
│  │ │ • Receives commands from coordination server     │  │ │
│  │ │ • Executes provider ops (lxc launch, etc.)      │  │ │
│  │ │ • Reports status back                           │  │ │
│  │ └──────────────────────────────────────────────────┘  │ │
│  │                   ↓ executes                           │ │
│  │ ┌──────────────────────────────────────────────────┐  │ │
│  │ │ LXC Daemon                                       │  │ │
│  │ │ ┌──────────────┐ ┌──────────────┐              │  │ │
│  │ │ │ Container 1  │ │ Container 2  │              │  │ │
│  │ │ │ (my-project) │ │ (other-proj) │              │  │ │
│  │ │ │ SSH server   │ │ SSH server   │              │  │ │
│  │ │ │ Services     │ │ Services     │              │  │ │
│  │ │ └──────────────┘ └──────────────┘              │  │ │
│  │ └──────────────────────────────────────────────────┘  │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                              │
│  SSH Port Forwarding: 2222 → container1:22                 │
│                       2223 → container2:22                 │
└─────────────────────────────────────────────────────────────┘
         ↑ SSH connections
         │
    Designer's Local Machine
    ├── nexus CLI (installed via install.sh)
    ├── GitHub CLI (auto-installed)
    ├── Editor: Cursor/VS Code (SSH remote)
    └── Terminal: SSH direct access
```

---

## Rollout Strategy

### Week 1: Internal
- Deploy coordination server + LXC node
- Manual E2E testing
- Fix critical bugs

### Week 2-3: Beta
- 5-10 designer beta testers
- Gather feedback
- Iterate on UX

### Week 4: GA
- Public announcement
- Documentation
- Support runbooks

---

## Key Metrics to Track

| What | Target | How to Measure |
|------|--------|---|
| First workspace time | < 3 min | Time from install.sh to editor opening |
| Repeat workspace time | < 1 min | Cached deps + existing config |
| Service startup | < 45 sec | From container creation to healthy |
| SSH latency | < 100ms | `time ssh ... 'echo hi'` |
| Error messages clarity | 100% | All errors have next steps |
| Editor launch success | 95% | Deep links work for Cursor/Code |
| Server uptime | 99.9% | Monitoring & alerting |
| Concurrent workspaces | 10+ | Load test with multiple users |

---

## Risks & Contingencies

| Risk | Likelihood | Impact | Mitigation |
|------|---|---|---|
| SSH key upload fails | Medium | Blocker | Retry logic + clear error + manual fallback |
| Container doesn't start | Low | Blocker | Health checks + timeout + helpful error |
| Service depends_on broken | Medium | Partial fail | Parse config, start in order, partial status |
| Editor deep link fails | Medium | Annoying | Print SSH command as fallback |
| LXC daemon crashes | Low | Major outage | Auto-restart + monitoring |

---

## Related Documents

1. **`M3_STAGING_ENV_USER_FLOW_SPEC.md`**
   - Detailed step-by-step user flow
   - Decision trees for edge cases
   - API contract specifications
   - Configuration file examples

2. **`M3_IMPLEMENTATION_PLAN.md`**
   - Phase-by-phase implementation guide
   - Detailed checklist
   - Testing strategy
   - Success criteria

3. **`M3_NEXT_STEPS.md`** (existing)
   - Coordination server priority
   - Provider remote support
   - SSH automation

---

## Questions & Clarifications Needed

Before starting implementation, confirm:

1. **GitHub App vs. CLI Auth**: Use GitHub CLI (`gh auth login`) or custom OAuth app?
   - **Recommendation**: GitHub CLI (simpler, users already have it)

2. **SSH Key Storage**: Store user SSH private keys locally only?
   - **Recommendation**: Yes, never on coordination server (security)

3. **Container Images**: Pre-build images with dev tools or install on-demand?
   - **Recommendation**: Start with generic ubuntu:22.04, add dev tools on first run

4. **Service Configuration**: Require .nexus/config.yaml or auto-detect?
   - **Recommendation**: Both: auto-detect for simple projects, config.yaml for complex

5. **Multi-Provider**: Support Docker/QEMU in staging or LXC only?
   - **Recommendation**: LXC only for staging (lighter, faster), multi-provider in future

---

## Success Looks Like

A designer opens their terminal and runs:

```bash
$ curl https://nexus.example.com/install.sh | bash -s -- \
    --repo my-org/my-project --server staging.example.com

📦 Downloading nexus CLI...
 Downloaded and installed

⚙  Checking dependencies...
 git, ssh, curl found
📦 Installing GitHub CLI...
 GitHub CLI installed

🔐 GitHub Authentication
[Opens browser]
 Authenticated as: alice

🔑 SSH Key Setup
 SSH key generated
 Uploaded to GitHub

📦 Repository Setup
 Repository forked to alice/my-project
 Cloned to ~/my-project

🔨 Workspace Initialization
⏳ Creating workspace (estimated 45s)...
 Container ready
 SSH configured
 Services starting...
 Services ready

 Workspace Ready
Server: staging.example.com
SSH: ssh -p 2222 dev@staging.example.com
Services:
  - web: http://localhost:23000 (npm dev)
  - api: http://localhost:23001 (npm server)

 Opening Cursor...
[Cursor opens with SSH remote connected]

Ready to design! Start by editing files in ~/my-project/src
```

**Total time**: 3 minutes  
**Designer friction**: Zero  
**Ready to contribute**: Yes 

---

**Document Status**: Specification Ready for Implementation  
**Approval**: Pending review  
**Next Action**: Begin Phase 1 implementation (Coordination Server)  
**Target Completion**: 4 weeks (late February 2026)
