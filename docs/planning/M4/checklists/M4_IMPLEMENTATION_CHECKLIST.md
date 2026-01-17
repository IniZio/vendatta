# M4 Implementation Checklist

**Purpose**: Detailed implementation tasks for M4 development  
**Duration**: 4-6 weeks across 4 phases  
**Team**: Backend engineers, systems engineers, QA

---

## Phase 1: Coordination Server Foundation (Weeks 1-2)

### 1.1 Project Setup
- [ ] Create `cmd/coordination-server/main.go` entry point
- [ ] Create `pkg/coordination/` package structure
- [ ] Set up Go modules & dependencies (chi/mux, sqlc, etc.)
- [ ] Create SQLite schema for workspaces & users
- [ ] Set up logging (structured JSON logs)

### 1.2 HTTP Server & Routing
- [ ] Implement HTTP server on port 3001
- [ ] Create request/response middleware (auth, logging, CORS)
- [ ] Implement request ID generation & tracking
- [ ] Set up error handling middleware
- [ ] Create health check endpoint (`GET /health`)

### 1.3 API Endpoints (Basic)
- [ ] `POST /api/v1/users/register` - User registration
- [ ] `POST /api/v1/workspaces/create` - Create workspace
- [ ] `GET /api/v1/workspaces/{id}/status` - Get status
- [ ] `GET /api/v1/workspaces` - List workspaces
- [ ] `POST /api/v1/workspaces/{id}/stop` - Stop workspace
- [ ] `DELETE /api/v1/workspaces/{id}` - Delete workspace

### 1.4 Data Models & Storage
- [ ] Define Workspace struct with all fields
- [ ] Define User struct with SSH key storage
- [ ] Define Service struct with health check config
- [ ] Implement database layer (save, retrieve, update, delete)
- [ ] Create database migrations
- [ ] Implement workspace status transitions

### 1.5 Node Agent Communication
- [ ] Design agent communication protocol (SSH commands)
- [ ] Create command serialization (JSON)
- [ ] Implement command dispatch to agents
- [ ] Implement result/status collection from agents
- [ ] Create agent registry (track connected nodes)
- [ ] Implement heartbeat monitoring

### 1.6 SSH Port Management
- [ ] Implement SSH port allocation (pool: 2222-2299)
- [ ] Create port forwarding setup
- [ ] Implement port release on workspace deletion
- [ ] Track port assignments in database

### 1.7 Testing (Phase 1)
- [ ] Unit tests for data models
- [ ] Unit tests for workspace status transitions
- [ ] Unit tests for port allocation logic
- [ ] Integration tests for API endpoints
- [ ] Integration tests for database operations

**Phase 1 Success Criteria**:
-  Coordination server starts & runs
-  All API endpoints respond correctly
-  Workspace CRUD operations work
-  Node agents can be registered
-  SSH port allocation works

---

## Phase 2: GitHub CLI Integration (Weeks 2-3)

### 2.1 CLI Commands: GitHub Auth
- [ ] Create `nexus auth github` command
- [ ] Detect GitHub CLI installation
- [ ] Orchestrate `gh auth login` if needed
- [ ] Extract username via `gh api user --jq '.login'`
- [ ] Extract user ID via GitHub API
- [ ] Store credentials securely (~/.netrc or credential manager)
- [ ] Verify authentication with `gh auth status`

### 2.2 CLI Commands: SSH Setup
- [ ] Create `nexus ssh setup` command
- [ ] Detect existing SSH keys in ~/.ssh/
- [ ] Implement key validation (format, permissions)
- [ ] Implement SSH key generation (ed25519)
- [ ] Implement GitHub upload via `gh ssh-key add`
- [ ] Handle duplicate key detection (409 conflict)
- [ ] Store key info in ~/.nexus/config.yaml

### 2.3 CLI Commands: Workspace Creation
- [ ] Create `nexus workspace create <repo>` command
- [ ] Implement repo ownership check via GitHub API
- [ ] Implement fork logic via `gh repo fork`
- [ ] Implement clone logic (`git clone` or `gh repo clone`)
- [ ] Load .nexus/config.yaml from repo (if exists)
- [ ] Generate default config (if not exists)
- [ ] Call coordination server POST /api/v1/workspaces/create
- [ ] Poll status until workspace ready

### 2.4 CLI Commands: Workspace Connect
- [ ] Create `nexus workspace connect` command
- [ ] Detect editor preference (Cursor > Code > Vim)
- [ ] Generate SSH deep links (Cursor, VS Code)
- [ ] Validate SSH connection before launch
- [ ] Launch editor with remote connection
- [ ] Display summary with service URLs
- [ ] Print fallback SSH command

### 2.5 CLI Commands: Utilities
- [ ] Create `nexus workspace services` command
- [ ] Create `nexus workspace status <name>` command
- [ ] Create `nexus workspace logs <name>` command
- [ ] Create `nexus workspace exec <name> <cmd>` command

### 2.6 Configuration Management
- [ ] Create ~/.nexus/config.yaml parser
- [ ] Implement config validation
- [ ] Store GitHub username & token
- [ ] Store SSH key path & type
- [ ] Store editor preference
- [ ] Implement config migration for updates

### 2.7 Testing (Phase 2)
- [ ] Unit tests for GitHub API interactions
- [ ] Unit tests for SSH key operations
- [ ] Integration tests for GitHub auth flow
- [ ] Integration tests for repo fork/clone
- [ ] E2E tests with test GitHub repo
- [ ] Tests for config file handling

**Phase 2 Success Criteria**:
-  Full GitHub authentication flow works
-  SSH key generation & upload works
-  Repository fork/clone works
-  Workspace creation end-to-end works
-  Configuration properly managed

---

## Phase 3: One-Line Install Script (Week 3)

### 3.1 Install Script
- [ ] Create `scripts/install.sh`
- [ ] Parse arguments (--repo, --server)
- [ ] Detect OS (macOS, Linux)
- [ ] Detect architecture (x86_64, arm64)
- [ ] Download nexus binary from releases
- [ ] Verify checksum/signature
- [ ] Install to ~/.local/bin/ or /usr/local/bin
- [ ] Make executable

### 3.2 Dependency Detection
- [ ] Check for git
- [ ] Check for ssh
- [ ] Check for curl
- [ ] Detect and install GitHub CLI (if missing)
- [ ] Provide helpful instructions if deps missing

### 3.3 Configuration Setup
- [ ] Create ~/.nexus/ directory
- [ ] Create ~/.nexus/config.yaml with server address
- [ ] Create ~/.ssh/ if not exists
- [ ] Set proper permissions (0700 for .ssh)

### 3.4 Setup Workflow Kickoff
- [ ] Call `nexus workspace setup-from-repo <repo> --server <server>`
- [ ] Hand off to interactive setup

### 3.5 Platform Support
- [ ] macOS (Intel)
- [ ] macOS (Apple Silicon)
- [ ] Ubuntu/Debian
- [ ] RHEL/CentOS
- [ ] Alpine (optional for Phase 3)

### 3.6 Testing (Phase 3)
- [ ] Test install script on macOS
- [ ] Test install script on Ubuntu
- [ ] Test with fresh user account
- [ ] Test with pre-existing SSH keys
- [ ] Test error scenarios (no internet, etc.)

**Phase 3 Success Criteria**:
-  Install script works on macOS & Linux
-  Dependencies auto-installed
-  Configuration properly initialized
-  Setup workflow begins automatically
-  Users report good experience

---

## Phase 4: Polish & Launch (Weeks 4-5)

### 4.1 End-to-End Testing
- [ ] Complete flow: install → auth → workspace → ready
- [ ] Fresh user account testing
- [ ] User with existing GitHub auth
- [ ] User with existing SSH keys
- [ ] Different repository types (Node, Python, static)
- [ ] Projects with .nexus/config.yaml
- [ ] Projects without config (use defaults)

### 4.2 Error Scenarios
- [ ] Network connectivity failures
- [ ] GitHub API rate limiting
- [ ] SSH key generation failures
- [ ] Repository not found
- [ ] Insufficient disk space
- [ ] LXC daemon not running
- [ ] Container startup timeout
- [ ] Service startup failure
- [ ] SSH connection failure

### 4.3 Multi-User & Concurrency
- [ ] Multiple concurrent workspace creation
- [ ] Multiple users on same staging server
- [ ] SSH port conflicts (should not happen)
- [ ] Service port conflicts
- [ ] Concurrent service health checks

### 4.4 Performance Optimization
- [ ] Container creation time < 30s
- [ ] Service startup time < 45s
- [ ] API response time < 500ms
- [ ] SSH latency < 100ms
- [ ] Health check time < 5s per service

### 4.5 UX Refinement
- [ ] Progress indicators (spinners, percentage)
- [ ] Colored output (success, warning, error)
- [ ] Command copy/paste for troubleshooting
- [ ] Clear next steps in error messages
- [ ] Helpful documentation links
- [ ] Summary with quick reference commands

### 4.6 Monitoring & Observability
- [ ] Set up logging aggregation
- [ ] Create monitoring dashboard
- [ ] Set up alerts for errors
- [ ] Implement metrics collection
- [ ] Create health check endpoint with detailed status

### 4.7 Documentation
- [ ] User onboarding guide
- [ ] Troubleshooting guide
- [ ] API documentation (OpenAPI/Swagger)
- [ ] Architecture diagrams
- [ ] Contributing guide
- [ ] FAQ document

### 4.8 Deployment Preparation
- [ ] Create deployment guide
- [ ] Create systemd service file (if applicable)
- [ ] Create monitoring/alerting setup
- [ ] Create backup strategy
- [ ] Create rollback procedure
- [ ] Create runbook for common issues

### 4.9 Beta Testing
- [ ] Recruit 5-10 beta testers (designers/developers)
- [ ] Collect feedback
- [ ] Document issues & suggestions
- [ ] Iterate on feedback
- [ ] Measure success metrics

### 4.10 GA Launch
- [ ] Final testing & verification
- [ ] Create launch announcement
- [ ] Deploy to staging environment
- [ ] Monitor closely for issues
- [ ] Support beta testers
- [ ] Collect usage metrics

### 4.11 Testing (Phase 4)
- [ ] Full E2E test suite
- [ ] Load testing (10+ concurrent users)
- [ ] Stress testing (resources)
- [ ] Chaos testing (failure scenarios)
- [ ] Security testing (access control)
- [ ] Performance testing

**Phase 4 Success Criteria**:
-  All functionality tested & working
-  Error scenarios handled gracefully
-  Multi-user scenarios stable
-  Performance targets met
-  UX polished & intuitive
-  Documentation complete
-  Ready for public launch

---

## Cross-Phase Requirements

### Code Quality
- [ ] All code follows Go conventions
- [ ] `go fmt` passes
- [ ] `go vet` passes
- [ ] Linter passes (golangci-lint)
- [ ] Test coverage > 80% on new code
- [ ] Code review & approval

### Documentation
- [ ] Code comments for complex logic
- [ ] API documentation (OpenAPI)
- [ ] Architecture documentation
- [ ] Deployment documentation
- [ ] User guides

### Testing
- [ ] Unit tests (80%+ coverage)
- [ ] Integration tests
- [ ] E2E tests
- [ ] Load tests
- [ ] Security tests

### Deployment
- [ ] Staging environment ready
- [ ] Infrastructure documented
- [ ] Deployment process documented
- [ ] Monitoring & alerting set up
- [ ] Backup & recovery plan

---

## Acceptance Criteria

### For "Done"
-  All checklist items completed
-  All tests passing
-  Code reviewed & approved
-  Documentation complete
-  Ready for user testing

### For "Ship"
-  Beta testing positive
-  All known issues resolved
-  Performance targets met
-  Security reviewed
-  Go-live plan ready

---

## Tracking & Status

| Phase | Target | Status | Completion |
|-------|--------|--------|------------|
| Phase 1 | Week 2 | TBD | __ % |
| Phase 2 | Week 3 | TBD | __ % |
| Phase 3 | Week 3 | TBD | __ % |
| Phase 4 | Week 5 | TBD | __ % |

---

**Document**: M4 Implementation Checklist  
**Version**: 1.0  
**Status**: Final  
**Date**: January 17, 2026  
**Used By**: Implementation team for task tracking
