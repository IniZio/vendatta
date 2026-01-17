# M4 Strategic Execution Plan with Subagent Delegation
**Version**: 1.0  
**Status**: Ready for Approval & Execution  
**Created**: January 17, 2026  
**Duration**: 6 weeks (Jan 20 - Mar 3, 2026)  
**Complexity**: 4 phases, 8 parallel work streams, $~1200-1500 in engineering effort

---

## CURRENT STATE ASSESSMENT

### Codebase Maturity: ⭐⭐⭐⭐ (MATURE)

The project has **strong foundational infrastructure**:

| Component | Status | Readiness |
|-----------|--------|-----------|
| Coordination Server (`pkg/coordination/`) | ✅ EXISTS | 60% complete - needs extensions |
| Node Agent (`pkg/agent/`) | ✅ EXISTS | 70% complete - needs M4 integration |
| SSH Transport (`pkg/transport/`) | ✅ EXISTS | 80% complete - needs hardening |
| Config System | ✅ EXISTS | 100% complete |
| Providers (Docker/LXC/QEMU) | ✅ EXISTS | 85% complete locally |
| Testing Infrastructure | ✅ EXISTS | 55 tests, testify framework, Makefile, CI/CD |

**Key Advantage**: We're not building from scratch. M3 completed major infrastructure. M4 is **augmentation + integration + user experience**.

### What Works Today
- ✅ Local branch creation (Docker, LXC, QEMU)
- ✅ Git worktree management
- ✅ Service orchestration (basic)
- ✅ SSH key generation & management
- ✅ Configuration parsing
- ✅ Agent communication (SSH, heartbeat)
- ✅ Coordination server HTTP API (basic)

### What M4 Adds
- ❌ GitHub OAuth integration
- ❌ Automatic fork/clone workflow
- ❌ SSH key upload to GitHub
- ❌ One-line install script
- ❌ Editor detection & deep links
- ❌ Service discovery & port forwarding
- ❌ Advanced UX (progress, errors, formatting)
- ❌ Multi-workspace isolation verification

---

## PHASE TIMELINE & DEPENDENCIES

```
┌─────────────────────────────────────────────────────────────┐
│                 M4 EXECUTION TIMELINE                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ Week 1 (Jan 20-26): Phase 1 - Coordination Server Setup    │
│   - STREAM 1.1: Extended API endpoints                      │
│   - STREAM 1.2: LXC integration improvements                │
│   - STREAM 1.3: Service discovery & port forwarding         │
│   - STREAM 1.4: Comprehensive testing                       │
│   GATE: All endpoints working, >90% test coverage           │
│                                                              │
│ Week 2 (Jan 27-Feb 2): Phase 1 Continues + Phase 2 Starts  │
│   - Phase 1: Final verification & bug fixes                 │
│   - Phase 2 STREAM 2.1: GitHub auth orchestration           │
│   - Phase 2 STREAM 2.2: Workspace fork/clone logic          │
│   GATE PHASE 1: Server 99.9% uptime, 10+ concurrent spaces │
│                                                              │
│ Week 3 (Feb 3-9): Phase 2 - GitHub Integration             │
│   - STREAM 2.3: Editor detection & deep links              │
│   - STREAM 2.4: Service discovery UI                        │
│   - STREAM 2.5: Full E2E testing (<3 min workflow)         │
│   GATE: Full user flow working, auth → create → connect    │
│                                                              │
│ Week 4 (Feb 10-16): Phase 2 Polish + Phase 3 Starts        │
│   - Phase 2: UX refinement, error handling                  │
│   - Phase 3 STREAM 3.1: Install script                      │
│   - Phase 3 STREAM 3.2: Binary distribution                 │
│   GATE: Install script works on macOS, Linux, RHEL          │
│                                                              │
│ Week 5 (Feb 17-23): Phase 3 + Phase 4 Polish Begins        │
│   - Phase 3: Final testing & documentation                  │
│   - Phase 4 STREAM 4.1: Load testing                        │
│   - Phase 4 STREAM 4.2: UX polish                           │
│                                                              │
│ Week 6 (Feb 24-Mar 3): Phase 4 - Final Launch              │
│   - Security audit & fixes                                   │
│   - Complete documentation                                  │
│   - Staging environment launch                              │
│   - Beta user recruitment & testing                         │
│   FINAL GATE: Go/no-go for public launch                    │
│                                                              │
└─────────────────────────────────────────────────────────────┘

CRITICAL PATHS:
  - Phase 1 → 2: Coordination Server must be stable
  - Phase 2 → 3: GitHub integration must be bulletproof
  - Phase 3 → 4: Install script must work on all platforms
  - Phase 4: All gates must pass before launch
```

---

## DELEGATION MATRIX

### Subagent Team Structure

```
YOU (Sisyphus)
├─ Strategic Direction & Architecture Review
│
├─ Backend Engineer (sisyphus-junior-high)
│  ├─ Phase 1.1: Extended API endpoints (80h)
│  ├─ Phase 1.3: Database & metadata (40h)
│  ├─ Phase 2.1: GitHub auth (80h)
│  ├─ Phase 2.2: Fork/clone logic (80h)
│  └─ Phase 4.5: Final integration (80h)
│  Total: 360 hours
│
├─ Systems Engineer (sisyphus-junior)
│  ├─ Phase 1.2: LXC integration (60h)
│  ├─ Phase 3.1: Install script (40h)
│  ├─ Phase 3.2: Binary signing (30h)
│  ├─ Phase 4.3: Security review (40h)
│  └─ Contingency planning
│  Total: 170 hours
│
├─ CLI/UX Engineer (sisyphus-junior)
│  ├─ Phase 2.3: Editor detection (40h)
│  ├─ Phase 2.4: Service discovery UI (40h)
│  ├─ Phase 4.2: UX polish (60h)
│  └─ Help text & completions
│  Total: 140 hours
│
└─ QA Specialist (qa-tester)
   ├─ Phase 1.4: Testing (60h)
   ├─ Phase 2.5: E2E testing (60h)
   ├─ Phase 3.3: Installation tests (50h)
   ├─ Phase 4.1: Load & reliability (100h)
   └─ Documentation review
   Total: 270 hours
```

**Total Engineering**: ~940 hours (~5 full-time engineer-weeks)  
**Cost Estimate**: $35-50k (assuming $75/hr blended rate)

---

## PHASE 1: COORDINATION SERVER FOUNDATION (2 WEEKS)

### Objective
Extend existing coordination server to support:
- GitHub repository metadata
- SSH port allocation & forwarding
- Service discovery & port mapping
- Workspace status lifecycle
- Multi-user isolation

### Deliverables by Stream

#### STREAM 1.1: Extended API Endpoints (80 hours)
**Assigned**: Backend Engineer  
**Existing Code**: `pkg/coordination/handlers.go`, `server.go`

**Specific Tasks**:
1. **POST /api/v1/users/register-github**
   - Accept GitHub username, ID, public key
   - Store in SQLite users table
   - Return user_id for future requests

2. **POST /api/v1/workspaces/create-from-repo**
   - Input: GitHub username, repo (owner/name), branch
   - Output: workspace_id, status, ssh_port
   - Dispatch to LXC node agent for container creation
   - Return 202 Accepted with polling URL

3. **GET /api/v1/workspaces/{id}/status**
   - Status: pending, creating, running, stopped, error
   - SSH details: host, port, username
   - Services: name, port, local_port, health_status
   - Repository info: owner, name, branch, commit

4. **GET /api/v1/workspaces/{id}/services**
   - List all services with URLs
   - Health check status
   - Port mapping details

5. **POST /api/v1/workspaces/{id}/stop**
6. **DELETE /api/v1/workspaces/{id}**
7. **GET /api/v1/workspaces** (list all)

**Testing Requirements**:
- Unit tests for all status transitions
- Integration tests for API endpoints
- Database persistence tests
- Error handling (validation, permissions)

**Success Criteria**:
- [ ] All endpoints respond per M4_API_SPECIFICATION.md
- [ ] Status transitions work correctly
- [ ] 90%+ test coverage
- [ ] Error messages clear & actionable

**Code Pattern to Follow**:
- Existing handlers use `chi.Router`
- Responses use JSON with `StatusCode` first
- Errors wrapped with `fmt.Errorf`
- Database layer in separate package

---

#### STREAM 1.2: LXC Driver Integration (60 hours)
**Assigned**: Systems Engineer  
**Existing Code**: `pkg/provider/lxc/`, `pkg/agent/node.go`

**Specific Tasks**:
1. **Node Agent Enhancement**
   - Receive `CreateWorkspaceCommand` from coordination server
   - Launch LXC container with Ubuntu 22.04 image
   - Configure SSH inside container
   - Install Git, SSH client, Node.js (default)
   - Setup service startup from `.nexus/config.yaml`
   - Report status back to server

2. **Container Lifecycle**
   - `lxc launch ubuntu:22.04 workspace-{id}`
   - Setup authorized_keys with user's SSH public key
   - Create `/home/dev` directory structure
   - Mount workspace repository (via git clone)
   - Port forwarding: 2222→22, 3000→23000, 4000→23001, etc.

3. **Service Orchestration**
   - Parse `.nexus/config.yaml` from repo
   - Start services in dependency order
   - Health checks (HTTP, TCP, timeout)
   - Report service status

4. **SSH Connection**
   - Ensure SSH server runs inside container
   - Configure for remote key authentication
   - Test connection from coordination server

**Testing Requirements**:
- Container creation/deletion tests
- SSH access validation
- Service startup tests
- Port forwarding tests

**Success Criteria**:
- [ ] Container creates in <30s
- [ ] SSH connection established <5s
- [ ] Services start correctly
- [ ] Multiple containers isolated

---

#### STREAM 1.3: Database & Metadata (40 hours)
**Assigned**: Backend Engineer  
**Existing Code**: `pkg/coordination/` (if any)

**Specific Tasks**:
1. **SQLite Schema**
   ```sql
   -- Users table
   CREATE TABLE users (
     user_id TEXT PRIMARY KEY,
     github_username TEXT UNIQUE NOT NULL,
     github_id INTEGER NOT NULL,
     ssh_pubkey TEXT NOT NULL,
     ssh_pubkey_fingerprint TEXT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   -- Workspaces table
   CREATE TABLE workspaces (
     workspace_id TEXT PRIMARY KEY,
     user_id TEXT NOT NULL REFERENCES users(user_id),
     workspace_name TEXT NOT NULL,
     status TEXT NOT NULL, -- pending, creating, running, stopped, error
     ssh_port INTEGER NOT NULL,
     repo_owner TEXT NOT NULL,
     repo_name TEXT NOT NULL,
     repo_url TEXT NOT NULL,
     repo_branch TEXT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   -- Services table
   CREATE TABLE services (
     service_id TEXT PRIMARY KEY,
     workspace_id TEXT NOT NULL REFERENCES workspaces(workspace_id),
     service_name TEXT NOT NULL,
     command TEXT NOT NULL,
     port INTEGER NOT NULL,
     local_port INTEGER,
     status TEXT, -- running, stopped, error
     health_status TEXT, -- healthy, unhealthy
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

2. **DAOs (Data Access Objects)**
   - User: Create, Read, Update, Delete
   - Workspace: CRUD, Status updates, List by user
   - Service: CRUD, Status updates

3. **Database Initialization**
   - Schema creation on startup
   - Migrations support
   - Connection pooling

**Testing Requirements**:
- Schema validation tests
- DAO CRUD tests
- State transition tests
- Transaction tests

**Success Criteria**:
- [ ] Schema matches spec
- [ ] All CRUD operations work
- [ ] 90%+ test coverage
- [ ] No data corruption

---

#### STREAM 1.4: Testing & Documentation (60 hours)
**Assigned**: QA Specialist  
**Tools**: testify, Go test, manual E2E

**Specific Tasks**:
1. **Unit Tests**
   - Workspace state transitions (20 tests)
   - API request/response parsing (30 tests)
   - Database CRUD operations (25 tests)
   - Port allocation logic (10 tests)

2. **Integration Tests**
   - API endpoint integration (30 tests)
   - Database persistence (15 tests)
   - Service startup (20 tests)

3. **E2E Tests**
   - Full workspace creation workflow
   - Container lifecycle (create → run → stop → delete)
   - SSH access validation
   - Service health checks

4. **Performance Tests**
   - Container startup timing (<30s)
   - API response time (<100ms)
   - SSH latency (<100ms)
   - Concurrent workspace handling (5, 10, 20 workspaces)

5. **Documentation**
   - API OpenAPI spec (Swagger)
   - Database schema docs
   - Architecture diagram updates
   - Setup & deployment guide

**Testing Commands**:
```bash
make test                    # Run all tests
make test-coverage          # Coverage report
make test-integration       # Integration tests
make test-e2e              # E2E tests
make test-performance      # Performance benchmarks
```

**Success Criteria**:
- [ ] 90%+ code coverage
- [ ] All tests passing
- [ ] Performance targets met
- [ ] OpenAPI spec complete

---

### Phase 1 Gate Review (CRITICAL)

**Before Moving to Phase 2, verify**:
1. Code review: Architecture correctness, pattern compliance
2. Test execution: All tests passing, coverage >90%
3. Performance: Container startup <30s, SSH latency <100ms
4. Security: SSH setup reviewed, no credential leaks
5. E2E validation: Create workspace, SSH in, start services

**Approval**: Sisyphus sign-off required

---

## PHASE 2: GITHUB CLI INTEGRATION (1.5 WEEKS)

### Objective
Automate user onboarding:
- GitHub OAuth authentication
- SSH key generation & upload to GitHub
- Automatic repository fork/clone
- Workspace creation from repository
- Editor launch with SSH connection

### Deliverables by Stream

#### STREAM 2.1: GitHub Auth & SSH Setup (80 hours)
**Assigned**: Backend Engineer  
**Existing Code**: `pkg/auth/`, `cmd/nexus/main.go`

**Tasks**:
1. **`nexus auth github` Command**
   - Detect GitHub CLI installation (`gh --version`)
   - Check authentication status (`gh auth status`)
   - If not authenticated: Orchestrate `gh auth login`
   - Extract username: `gh api user --jq '.login'`
   - Extract user ID: `gh api user --jq '.id'`
   - Extract auth token from GitHub CLI config
   - Save to `~/.nexus/credentials.json` (encrypted)
   - Verify with `gh auth status`

2. **`nexus ssh setup` Command**
   - Detect existing keys: `ls ~/.ssh/id_ed25519*`
   - Validate key format (ed25519)
   - If missing: Generate with `ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -N ""`
   - Upload to GitHub: `gh ssh-key add ~/.ssh/id_ed25519.pub --title "nexus@staging"`
   - Handle errors: Duplicate keys (409), network issues (retry)
   - Extract fingerprint from public key
   - Register fingerprint with coordination server
   - Verify with `gh ssh-key list`

3. **Error Handling**
   - GitHub CLI not installed → Clear instructions to install
   - Not authenticated → Guide to `gh auth login`
   - SSH key already on GitHub → Confirm reuse
   - Network issues → Retry with exponential backoff

**Testing**:
- Mock GitHub CLI responses
- Test SSH key detection
- Test key generation & validation
- Test error scenarios

**Success Criteria**:
- [ ] Full auth flow works <30s
- [ ] SSH key uploaded to GitHub
- [ ] Clear error messages for failures
- [ ] Works on macOS & Linux

---

#### STREAM 2.2: Workspace Creation & Fork Logic (80 hours)
**Assigned**: Backend Engineer  
**Existing Code**: CLI structure in `cmd/nexus/`

**Tasks**:
1. **`nexus workspace create <repo>` Command**
   - Input: repo in format `owner/name` (e.g., `my-org/my-project`)
   - Check ownership: `gh repo view {owner}/{name} --json owner`
   - If owned by user: Clone directly
   - If not owned: Fork with `gh repo fork {owner}/{name} --clone`
   - Clone to `~/{name}` (or user-specified path)
   - Change directory to cloned repo
   - Load `.nexus/config.yaml` (if exists)
   - Generate default config (if missing)
   - Call coordination server:
     ```json
     POST /api/v1/workspaces/create-from-repo
     {
       "github_username": "alice",
       "workspace_name": "my-project-main",
       "repo": {
         "owner": "my-org",
         "name": "my-project",
         "branch": "main"
       },
       "provider": "lxc"
     }
     ```
   - Poll `/api/v1/workspaces/{id}/status` until running
   - Display progress with spinner

2. **Repository Fork Logic**
   - Check if repo owned: `gh repo view {owner}/{name} --json owner.login`
   - If owned by user: Skip fork, clone directly
   - If not owned: Fork with `gh repo fork {owner}/{name}`
   - Handle fork conflicts (already forked)
   - Clone to local machine

3. **Configuration**
   - Load `.nexus/config.yaml` from repo root
   - If missing: Generate default based on:
     - `package.json` (Node.js)
     - `pyproject.toml` / `requirements.txt` (Python)
     - `Gemfile` (Ruby)
     - Static site (default web server)
   - Example config:
     ```yaml
     name: my-project
     provider: lxc
     services:
       web:
         command: "npm run dev"
         port: 3000
       api:
         command: "npm run server"
         port: 4000
     ```

**Testing**:
- Mock GitHub API responses
- Test fork detection & creation
- Test config loading
- Test default config generation

**Success Criteria**:
- [ ] Works with owned repos
- [ ] Works with forks
- [ ] Config loading works
- [ ] Error handling robust

---

#### STREAM 2.3: Editor Detection & Deep Links (40 hours)
**Assigned**: CLI/UX Engineer

**Tasks**:
1. **`nexus workspace connect <name>` Command**
   - Retrieve workspace from coordination server
   - Get SSH details: host, port, username
   - Detect editor:
     - `which cursor` → Cursor
     - `which code` → VS Code
     - `which vim` → Vim (fallback)
   - Generate deep links:
     - Cursor: `cursor://ssh/remote?host=staging&port=2222&user=dev&path=/home/dev/my-project`
     - VS Code: `vscode://vscode-remote/ssh-remote+dev@staging:2222/home/dev/my-project`
   - Launch editor:
     - Cursor: `cursor --remote ssh-remote+dev@staging:2222 /home/dev/my-project`
     - VS Code: `code --remote ssh-remote+dev@staging:2222 /home/dev/my-project`
   - Fallback: Print SSH command for manual use

2. **SSH Validation**
   - Test SSH connection before marking ready
   - `ssh -p 2222 dev@staging.example.com 'echo OK'`
   - Timeout: 5 seconds
   - If fails: Display troubleshooting steps

3. **Summary Display**
   - Workspace name & status
   - Services with URLs
   - SSH command (fallback)
   - Editor launch confirmation

**Testing**:
- Mock editor executables
- Test deep link generation
- Test SSH validation
- Test fallback scenarios

**Success Criteria**:
- [ ] Editor detection works
- [ ] Deep links correct
- [ ] Editor launches successfully
- [ ] Fallback to SSH works

---

#### STREAM 2.4: Service Discovery UI (40 hours)
**Assigned**: CLI/UX Engineer

**Tasks**:
1. **`nexus workspace services <name>` Command**
   - Query coordination server: `GET /api/v1/workspaces/{id}/services`
   - Display services table:
     ```
     Service    Status     Port    Local URL
     ─────────────────────────────────────
     web        running    3000    http://localhost:23000
     api        running    4000    http://localhost:23001
     db         running    5432    localhost:23002
     ```
   - Include health status (healthy/unhealthy)
   - Clickable URLs (in terminal that supports it)

2. **Port Forwarding Info**
   - Show local port mappings
   - Explain forwarding setup
   - Provide `ssh -L` command if manual forwarding needed

3. **Error Handling**
   - Service not found → List available services
   - Service unhealthy → Show last error
   - Connection issues → Retry logic

**Testing**:
- Service list API mocking
- Table formatting tests
- Error display tests

**Success Criteria**:
- [ ] Service list accurate
- [ ] URLs correct
- [ ] Pretty formatting works
- [ ] Error messages clear

---

#### STREAM 2.5: E2E Testing (<3 min workflow) (60 hours)
**Assigned**: QA Specialist

**Tests**:
1. **Full User Flow**
   ```bash
   # Start with fresh shell, no GitHub auth
   1. nexus auth github           # ~30s (browser redirect)
   2. nexus ssh setup             # ~15s (generate or reuse key)
   3. nexus workspace create my-org/my-project  # ~20s (fork/clone)
   4. # Workspace creation async: ~60s (polling)
   5. nexus workspace connect my-project        # ~5s (open editor)
   # Total: ~3 minutes
   ```

2. **Error Scenarios**
   - GitHub CLI not installed
   - Not authenticated
   - SSH key already exists (confirm reuse)
   - Fork already exists
   - Workspace creation timeout
   - Editor not installed (fallback to SSH)

3. **Multi-OS Testing**
   - macOS (Intel + Apple Silicon)
   - Ubuntu 20.04, 22.04
   - RHEL 8, 9

4. **Documentation**
   - Step-by-step guide
   - Screenshots/GIFs
   - Troubleshooting
   - FAQs

**Success Criteria**:
- [ ] Full workflow <3 minutes
- [ ] All error scenarios handled
- [ ] Works on macOS & Linux
- [ ] Documentation complete

---

### Phase 2 Gate Review

**Before Phase 3, verify**:
1. Full end-to-end workflow <3 minutes
2. GitHub auth & SSH key upload automated
3. All error messages actionable
4. Multi-OS compatibility confirmed
5. 90%+ test coverage

**Approval**: Sisyphus sign-off required

---

## PHASE 3: ONE-LINE INSTALL SCRIPT (1 WEEK)

### Objective
Single command to bootstrap everything:
```bash
curl https://nexus.example.com/install.sh | bash -s -- \
  --repo my-org/my-project --server staging.example.com
```

### Deliverables by Stream

#### STREAM 3.1: Install Script (40 hours)
**Assigned**: Systems Engineer

**Script Behavior**:
```bash
#!/bin/bash
set -e

# Parse arguments
REPO=${REPO:-""}
SERVER=${SERVER:-""}
DEST=${DEST:-"$HOME/.local/bin"}

# Step 1: Check dependencies
check_dependency() {
  which $1 >/dev/null || {
    echo "❌ $1 not found. Install with: [instructions]"
    exit 1
  }
}

check_dependency git
check_dependency ssh
check_dependency curl

# Step 2: Install GitHub CLI if missing
if ! which gh >/dev/null; then
  echo "📦 Installing GitHub CLI..."
  # Platform-specific installation
  # macOS: brew install gh
  # Linux: apt/yum/dnf install gh
fi

# Step 3: Download nexus binary
echo "📦 Installing nexus..."
NEXUS_VERSION="v0.1.0"
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

curl -L \
  -o "$DEST/nexus" \
  "https://github.com/nexus/releases/download/$NEXUS_VERSION/nexus-$PLATFORM-$ARCH"

chmod +x "$DEST/nexus"

# Step 4: Initialize config
echo "⚙️  Initializing configuration..."
mkdir -p "$HOME/.nexus"
cat > "$HOME/.nexus/config.yaml" <<EOF
server: $SERVER
server_port: 3001
repo: $REPO
EOF

# Step 5: Kick off setup
echo "🚀 Starting setup workflow..."
export PATH="$DEST:$PATH"
nexus workspace setup-from-repo "$REPO" --server "$SERVER"

echo "✅ Installation complete!"
```

#### STREAM 3.2: Binary Distribution (30 hours)
**Assigned**: Systems Engineer

**Tasks**:
1. **Build Artifacts**
   - macOS Intel (darwin-amd64)
   - macOS ARM (darwin-arm64)
   - Linux Intel (linux-amd64)
   - Linux ARM (linux-arm64)

2. **Code Signing**
   - macOS: Sign binary with developer certificate
   - Checksum: SHA256 for all binaries

3. **Release Automation**
   - GitHub release creation
   - Binary uploads
   - Checksum file generation

---

#### STREAM 3.3: Installation Testing (50 hours)
**Assigned**: QA Specialist

**Tests**:
- macOS 12+ (Intel + ARM)
- Ubuntu 20.04, 22.04
- RHEL 8, 9
- Dependency detection
- Clean uninstall

---

### Phase 3 Gate Review

**Before Phase 4, verify**:
1. Installation script works on all platforms
2. Binary signed and verifiable
3. Dependency detection working
4. Documentation complete

---

## PHASE 4: POLISH, TESTING & LAUNCH (2 WEEKS)

### Objective
Production-ready system, beta launch

### Deliverables

#### STREAM 4.1: Comprehensive Testing (100 hours)
- E2E test suite execution
- Load testing (10+ concurrent)
- Performance benchmarking
- Reliability testing (24h+)
- Multi-user scenarios
- Documentation

#### STREAM 4.2: UX Polish (60 hours)
- Progress messaging refinement
- Error message clarity
- Command help text
- Bash completion
- Color/emoji consistency

#### STREAM 4.3: Security Review (40 hours)
- SSH key handling audit
- GitHub token storage
- API authentication
- Credentials protection
- Access control

#### STREAM 4.4: Documentation (60 hours)
- User guide
- Installation guide (all platforms)
- Troubleshooting
- Architecture docs
- API documentation
- Example repos

#### STREAM 4.5: Final Integration & Launch (80 hours)
- Build & release preparation
- Staging environment setup
- Monitoring & alerting
- Rollout plan
- Beta user recruitment

---

## QUALITY GATES & VERIFICATION STRATEGY

### Gate 1: Phase 1 Completion
```
CRITERIA:
  ✓ All coordination server endpoints respond per spec
  ✓ Port allocation: 2222-2299 working
  ✓ Container startup: <30 seconds
  ✓ SSH latency: <100ms
  ✓ Concurrent workspaces: 10+ without interference
  ✓ Code coverage: 90%+ on new code
  ✓ All tests passing
  
REVIEW BY: Sisyphus (1 hour)
  - Code architecture
  - Test coverage analysis
  - E2E workflow validation
  - Performance metrics
  - Security review

SIGN-OFF: Required before Phase 2 starts
```

### Gate 2: Phase 2 Completion
```
CRITERIA:
  ✓ Full workflow: install → auth → create → connect <3 min
  ✓ All CLI commands working (auth, ssh setup, workspace create/connect)
  ✓ GitHub integration bulletproof (error handling)
  ✓ Multi-OS testing passed (macOS, Linux)
  ✓ Code coverage: 90%+ on new code
  ✓ All tests passing
  
REVIEW BY: Sisyphus (2 hours)
  - User flow validation (fresh test)
  - Error handling completeness
  - GitHub API robustness
  - Multi-OS compatibility
  - Code quality

SIGN-OFF: Required before Phase 3 starts
```

### Gate 3: Phase 3 Completion
```
CRITERIA:
  ✓ Install script works on all platforms
  ✓ Binary signed and verifiable
  ✓ Dependency detection robust
  ✓ Installation documentation complete
  
REVIEW BY: Sisyphus (1 hour)
  - Installation process validation
  - Binary security
  - Documentation quality

SIGN-OFF: Required before Phase 4 starts
```

### Gate 4: Phase 4 / Launch
```
CRITERIA:
  ✓ E2E test suite: 95%+ pass rate
  ✓ Load test: 10+ workspaces stable
  ✓ Performance: <3min workflow, <30s startup, <100ms SSH
  ✓ Security audit: Passed
  ✓ Documentation: Complete & accurate
  ✓ Staging environment: Ready
  
REVIEW BY: Sisyphus (3 hours)
  - Test results validation
  - Performance analysis
  - Security findings review
  - Documentation review
  - Go/no-go decision

SIGN-OFF: Go/no-go for public launch
```

---

## Risk Assessment & Mitigation

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|-----------|
| GitHub API rate limiting | Phase 2 delays | LOW | Implement caching, token refresh |
| SSH connection instability | Test failures | MEDIUM | Connection pooling, retry logic |
| Editor deep links fail | User frustration | LOW | Fallback to SSH command |
| LXC performance issues | Container startup >30s | LOW | Profiling, image optimization |
| Multi-workspace interference | Data corruption | LOW | Comprehensive isolation tests |
| Install script edge cases | Setup failures | MEDIUM | Platform-specific testing |

---

## Success Metrics

### User Experience
- [ ] First workspace: <3 minutes
- [ ] Container startup: <30 seconds
- [ ] SSH latency: <100ms
- [ ] Error clarity: 100% of errors actionable

### Technical
- [ ] Server uptime: 99.9%+
- [ ] Concurrent workspaces: 10+
- [ ] Code coverage: 90%+
- [ ] Test pass rate: 95%+

### Adoption
- [ ] Beta users: 20+
- [ ] Critical bugs in week 1: 0
- [ ] Support tickets/day: <5

---

## Summary

This plan converts M4 from a wishlist into executable work streams with clear delegation, concrete deliverables, and verification gates. Each phase builds on the previous, with strategic checkpoints to ensure quality.

**Key Success Factors**:
1. **Clear Delegation**: Each engineer knows exactly what they're building
2. **Strategic Oversight**: Sisyphus reviews each phase before approval
3. **Comprehensive Testing**: 90%+ coverage on all new code
4. **Quality Gates**: No moving forward until criteria met
5. **Risk Mitigation**: Contingency plans for high-risk areas

**Next Steps**:
1. **Today (Jan 17)**: Approve this plan
2. **Jan 20**: Begin Phase 1 work (parallel streams)
3. **Feb 10**: Phase 1 gate review & approval
4. **Feb 24**: Phase 4 final integration
5. **Mar 3**: Go/no-go for launch

---

**Prepared by**: Sisyphus  
**Status**: Ready for Approval & Delegation  
**Last Updated**: January 17, 2026
