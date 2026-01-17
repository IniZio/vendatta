# M4 Subagent Delegation Quick Reference

**For use by Sisyphus when delegating work to subagents**

---

## PHASE 1 DELEGATION TEMPLATES

### Stream 1.1: Extended API Endpoints (80 hours)
**Assign to**: sisyphus-junior-high (Backend Engineer)

```
TASK: Phase 1.1 - Extend Coordination Server API Endpoints

DELIVERABLE: 
  - 7 new API endpoints for workspace & service management
  - All endpoints in pkg/coordination/handlers.go
  - Full test coverage (90%+) in *_test.go
  - OpenAPI spec updated

MUST IMPLEMENT:
  1. POST /api/v1/users/register-github
     - Accept: github_username, github_id, ssh_pubkey
     - Store in SQLite users table
     - Return: user_id, registered_at
     
  2. POST /api/v1/workspaces/create-from-repo
     - Accept: github_username, workspace_name, repo info
     - Call LXC node agent to create container
     - Return: workspace_id, status, ssh_port
     
  3. GET /api/v1/workspaces/{id}/status
     - Return: status, SSH details, services list
     - Services include: name, port, local_port, health_status
     
  4. GET /api/v1/workspaces/{id}/services
     - List services with port mappings
     - Include health status & URLs
     
  5. POST /api/v1/workspaces/{id}/stop
  6. DELETE /api/v1/workspaces/{id}
  7. GET /api/v1/workspaces (list all, paginated)

REFERENCE:
  - M4_API_SPECIFICATION.md (exact API contracts)
  - M4_TECHNICAL_SPECIFICATION.md (architecture)
  - pkg/coordination/handlers.go (existing patterns)
  - STREAM 1.3 output (database schema)

TESTING:
  - Unit tests for each endpoint
  - Error cases: 400, 409, 500
  - Database persistence tests
  - 90%+ coverage required

GO/NO-GO CRITERIA:
  [ ] All endpoints respond per spec
  [ ] All tests passing
  [ ] 90%+ coverage
  [ ] Error messages clear
  [ ] No interface{} without justification
```

**Duration**: 80 hours (~2 weeks)  
**Deadline**: Feb 10, 2026

---

### Stream 1.2: LXC Integration (60 hours)
**Assign to**: sisyphus-junior (Systems Engineer)

```
TASK: Phase 1.2 - LXC Driver Node Integration

DELIVERABLE:
  - Node agent receives CreateWorkspaceCommand
  - Container launches with proper configuration
  - SSH & services operational
  - Full E2E test coverage

MUST IMPLEMENT:
  1. Command Reception (in pkg/agent/node.go)
     - Receive CreateWorkspaceCommand from coordination server
     - Parse workspace config, repo, services
     
  2. Container Lifecycle
     - lxc launch ubuntu:22.04 workspace-{id}
     - Configure LXC with CPU/memory limits
     - Mount storage for workspace data
     
  3. SSH Setup
     - Add authorized_keys with user's public key
     - Ensure SSH server runs on :22 inside container
     - Configure for remote key authentication
     
  4. Service Startup
     - Parse .nexus/config.yaml from repo
     - Install dependencies (npm, python, etc.)
     - Start services in dependency order
     - Health checks for each service
     
  5. Port Forwarding
     - Map SSH: 2222 → internal 22
     - Map services: 23000+ → internal ports
     - Report port mappings to coordination server

REFERENCE:
  - M4_TECHNICAL_SPECIFICATION.md (container design)
  - pkg/provider/lxc/lxc.go (existing LXC code)
  - pkg/agent/node.go (agent structure)
  - e2e/testenv.go (test utilities)

TESTING:
  - Container creation tests (<30s startup)
  - SSH connection validation
  - Service startup tests
  - Port forwarding tests
  - E2E full lifecycle

GO/NO-GO CRITERIA:
  [ ] Containers create consistently
  [ ] SSH accessible <5s after startup
  [ ] Services start correctly
  [ ] Multiple workspaces isolated
  [ ] Startup time <30s
```

**Duration**: 60 hours (~1.5 weeks)  
**Deadline**: Feb 10, 2026

---

### Stream 1.3: Database & Metadata (40 hours)
**Assign to**: sisyphus-junior-high (Backend Engineer)

```
TASK: Phase 1.3 - SQLite Database & Data Models

DELIVERABLE:
  - Complete SQLite schema with migrations
  - DAOs for users, workspaces, services
  - All CRUD operations working
  - Full test coverage (90%+)

MUST IMPLEMENT SCHEMA:
  CREATE TABLE users (
    user_id TEXT PRIMARY KEY,
    github_username TEXT UNIQUE,
    github_id INTEGER,
    ssh_pubkey TEXT,
    ssh_pubkey_fingerprint TEXT,
    created_at TIMESTAMP
  );

  CREATE TABLE workspaces (
    workspace_id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users,
    workspace_name TEXT,
    status TEXT, -- pending, creating, running, stopped
    ssh_port INTEGER,
    repo_owner TEXT,
    repo_name TEXT,
    repo_branch TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
  );

  CREATE TABLE services (
    service_id TEXT PRIMARY KEY,
    workspace_id TEXT REFERENCES workspaces,
    service_name TEXT,
    command TEXT,
    port INTEGER,
    local_port INTEGER,
    status TEXT,
    health_status TEXT,
    created_at TIMESTAMP
  );

MUST IMPLEMENT DAOs:
  - UserDAO: Create, Read, Update, Delete
  - WorkspaceDAO: CRUD + status updates
  - ServiceDAO: CRUD + health tracking
  - Workspace state transitions (validation)

REFERENCE:
  - M4_TECHNICAL_SPECIFICATION.md (data models)
  - pkg/coordination/ (where this lives)
  - Existing database code patterns

TESTING:
  - Schema validation
  - CRUD operations
  - State transitions (pending→creating→running)
  - Foreign key constraints
  - Transaction tests
  - 90%+ coverage

GO/NO-GO CRITERIA:
  [ ] Schema complete
  [ ] All DAOs implemented
  [ ] Tests passing
  [ ] 90%+ coverage
  [ ] No data corruption
```

**Duration**: 40 hours (~1 week)  
**Deadline**: Feb 10, 2026

---

### Stream 1.4: Comprehensive Testing (60 hours)
**Assign to**: qa-tester (QA Specialist)

```
TASK: Phase 1.4 - Testing & Documentation

DELIVERABLE:
  - 90%+ code coverage on new code
  - All unit/integration/E2E tests passing
  - Performance benchmarks meeting targets
  - API documentation (OpenAPI/Swagger)
  - Deployment guide

MUST TEST:
  1. Unit Tests (70 hours)
     - API endpoint request/response parsing
     - Workspace state transitions
     - Database CRUD operations
     - Port allocation logic
     - Error handling (400/409/500)
     
  2. Integration Tests (50 hours)
     - Full API endpoint flows
     - Database persistence
     - Service startup
     - Port forwarding
     
  3. E2E Tests (40 hours)
     - Full workspace creation → running → delete
     - SSH access validation
     - Service health checks
     - Concurrent workspace isolation
     
  4. Performance Tests (20 hours)
     - Container startup timing (<30s)
     - API response time (<100ms)
     - SSH latency (<100ms)
     - Concurrent load (5, 10, 20 workspaces)

REFERENCE:
  - M4_IMPLEMENTATION_CHECKLIST.md (what to test)
  - e2e/testenv.go (test utilities)
  - pkg/coordination/coordination_test.go (existing patterns)

GO/NO-GO CRITERIA:
  [ ] 90%+ code coverage (all new code)
  [ ] All tests passing
  [ ] Performance targets met
  [ ] OpenAPI spec complete
  [ ] Deployment guide written
  [ ] Error scenarios documented
```

**Duration**: 60 hours (~1.5 weeks)  
**Deadline**: Feb 10, 2026

---

## PHASE 2 DELEGATION TEMPLATES

### Stream 2.1: GitHub Auth & SSH Setup (80 hours)
**Assign to**: sisyphus-junior-high (Backend Engineer)

```
TASK: Phase 2.1 - GitHub CLI Orchestration & SSH Setup

DELIVERABLE:
  - nexus auth github command (fully working)
  - nexus ssh setup command (fully working)
  - GitHub auth token secure storage
  - SSH key generation & validation
  - Full test coverage (90%+)

MUST IMPLEMENT:
  1. nexus auth github
     - Detect GitHub CLI: which gh
     - Check auth: gh auth status
     - If not authed: orchestrate gh auth login
     - Extract username: gh api user --jq '.login'
     - Extract user ID: gh api user --jq '.id'
     - Store token securely in ~/.nexus/credentials.json
     - Verify with gh auth status
     
  2. nexus ssh setup
     - Detect existing keys: ls ~/.ssh/id_ed25519*
     - Validate key format
     - Generate if missing: ssh-keygen -t ed25519
     - Upload to GitHub: gh ssh-key add ~/.ssh/id_ed25519.pub
     - Handle duplicates gracefully (409)
     - Extract fingerprint from public key
     - Register with coordination server
     - Verify with gh ssh-key list

  3. Error Handling
     - GitHub CLI not installed → Install instructions
     - Not authenticated → Guide to gh auth login
     - SSH key already exists → Confirm reuse
     - Network errors → Retry logic
     - All errors must be actionable

REFERENCE:
  - M4_TECHNICAL_SPECIFICATION.md (auth design)
  - M4_USER_FLOW_SPECIFICATION.md (step 2-3)
  - cmd/nexus/main.go (CLI structure)
  - Go GitHub CLI API documentation

TESTING:
  - Mock GitHub CLI responses
  - Test SSH key detection/generation
  - Test error scenarios (network, duplicate, missing)
  - Integration with coordination server
  - 90%+ coverage

GO/NO-GO CRITERIA:
  [ ] Full auth flow works <30s
  [ ] SSH key uploaded to GitHub
  [ ] Error messages clear & actionable
  [ ] Works on macOS & Linux
  [ ] All tests passing
  [ ] 90%+ coverage
```

**Duration**: 80 hours (~2 weeks)  
**Deadline**: Feb 17, 2026

---

### Stream 2.2: Workspace Creation & Fork (80 hours)
**Assign to**: sisyphus-junior-high (Backend Engineer)

```
TASK: Phase 2.2 - Repository Fork/Clone & Workspace Creation

DELIVERABLE:
  - nexus workspace create <repo> command (fully working)
  - Automatic fork detection & creation
  - Repository cloning to local disk
  - .nexus/config.yaml loading
  - Default config generation
  - Full test coverage (90%+)

MUST IMPLEMENT:
  1. nexus workspace create <repo>
     - Input format: owner/name (e.g., my-org/my-project)
     - Check ownership: gh repo view {owner}/{name} --json owner
     - If owned: Clone directly with git clone
     - If not owned: Fork with gh repo fork {owner}/{name}
     - Clone to ~/{name}
     - Load .nexus/config.yaml (if exists)
     - Generate default config (if missing)
     
  2. Fork Logic
     - Detect ownership via GitHub API
     - Handle already-forked repos
     - Create fork if needed
     - Wait for fork completion before cloning
     
  3. Configuration
     - Load .nexus/config.yaml from repo root
     - Auto-generate for: Node.js, Python, Ruby, static sites
     - Example config structure (see M4 spec)
     - Validate config against schema
     
  4. Workspace Creation Request
     - Call POST /api/v1/workspaces/create-from-repo
     - Wait for workspace_id
     - Poll /api/v1/workspaces/{id}/status
     - Display progress indicator
     - Handle async creation errors

REFERENCE:
  - M4_USER_FLOW_SPECIFICATION.md (steps 4-5)
  - M4_TECHNICAL_SPECIFICATION.md (config design)
  - cmd/nexus/main.go (CLI structure)
  - M4_API_SPECIFICATION.md (coordination endpoints)

TESTING:
  - Mock GitHub API responses
  - Test with owned repos (no fork needed)
  - Test with unowned repos (fork needed)
  - Test already-forked repos
  - Test config loading & generation
  - Test error scenarios
  - 90%+ coverage

GO/NO-GO CRITERIA:
  [ ] Works with owned repos
  [ ] Works with unowned repos (auto-fork)
  [ ] Config loading works
  [ ] Default config generated correctly
  [ ] Workspace creation async works
  [ ] Error handling robust
  [ ] All tests passing
  [ ] 90%+ coverage
```

**Duration**: 80 hours (~2 weeks)  
**Deadline**: Feb 17, 2026

---

### Stream 2.3: Editor Detection & Deep Links (40 hours)
**Assign to**: sisyphus-junior (CLI/UX Engineer)

```
TASK: Phase 2.3 - Editor Detection & SSH Deep Links

DELIVERABLE:
  - nexus workspace connect <name> command
  - Automatic editor detection
  - SSH deep link generation
  - Editor launching with SSH
  - Summary display with service URLs
  - Full test coverage (90%+)

MUST IMPLEMENT:
  1. nexus workspace connect <name>
     - Retrieve workspace from coordination server
     - Get SSH details: host, port, user
     
  2. Editor Detection
     - Priority order: Cursor > VS Code > Vim
     - which cursor → Cursor
     - which code → VS Code
     - Fallback: Print SSH command
     
  3. Deep Link Generation
     - Cursor: cursor://ssh/remote?host=...&port=...&user=...
     - VS Code: vscode://vscode-remote/ssh-remote+user@host:port/path
     
  4. Editor Launch
     - Cursor: cursor --remote ssh-remote+user@host:port /path
     - VS Code: code --remote ssh-remote+user@host:port /path
     - Handle launch failures (missing editor)
     
  5. SSH Validation
     - Test SSH connection: ssh -p {port} {user}@{host} 'echo OK'
     - Timeout: 5 seconds
     - If fails: Display troubleshooting
     
  6. Summary Display
     - Workspace name & status
     - Services table (from workspace services call)
     - SSH command (fallback)
     - Editor launch confirmation

REFERENCE:
  - M4_USER_FLOW_SPECIFICATION.md (steps 6-7)
  - M4_TECHNICAL_SPECIFICATION.md (editor integration)
  - cmd/nexus/main.go (CLI structure)
  - M4_API_SPECIFICATION.md (services endpoint)

TESTING:
  - Mock editor executables
  - Test deep link generation (exact format)
  - Test SSH validation
  - Test fallback scenarios
  - Test summary formatting
  - 90%+ coverage

GO/NO-GO CRITERIA:
  [ ] Editor detection works (all three)
  [ ] Deep links formatted correctly
  [ ] Editor launches successfully
  [ ] SSH validation robust
  [ ] Summary display clear
  [ ] All tests passing
  [ ] 90%+ coverage
```

**Duration**: 40 hours (~1 week)  
**Deadline**: Feb 17, 2026

---

## QUICK REFERENCE: SUBAGENT TYPES

| Type | Best For | Examples |
|------|----------|----------|
| sisyphus-junior-high | Complex backend, architecture | Coordination server, API endpoints, data models |
| sisyphus-junior | Systems, infrastructure, scripts | LXC integration, install script, SSH setup |
| qa-tester | Testing, validation, verification | E2E tests, load tests, documentation |

---

## DELEGATION CHECKLIST

Before delegating:
- [ ] Clear task description (see templates above)
- [ ] Success criteria explicit
- [ ] Reference documents provided
- [ ] Test requirements stated (90%+ coverage)
- [ ] Deadline agreed upon
- [ ] No assumptions left unstated

After delegation:
- [ ] Daily 15-min standup
- [ ] Weekly checkpoint (progress review)
- [ ] Gate review before moving to next phase
- [ ] Code review by Sisyphus before merge

---

**Created**: January 17, 2026  
**For Use By**: Sisyphus when delegating M4 phases
