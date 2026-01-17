# M4 Execution Plan: Strategic Coordination & Delegation
**Version**: 1.0  
**Status**: IN PROGRESS (Planning Phase)  
**Created**: January 17, 2026  
**Target**: 4-6 weeks to full M4 completion (Feb 20 - Mar 1, 2026)

---

## Executive Summary

This document outlines the **strategic execution plan** for M4 with complete delegation to specialized subagents. The focus is on:

1. **Parallel M3 completion** (20% → 100%) while starting M4 Phase 1
2. **Clear delegation** - Backend engineers, systems engineers, QA specialists assigned to parallel work streams
3. **Direction oversight** - Strategic checkpoints to ensure architecture correctness
4. **Comprehensive verification** - All phases must pass strict testing before proceeding

### Critical Decision: M3 Blockers for M4

**Status**: M3 is 20% complete. Several critical components block M4:
- ❌ Coordination Server (0%) - **REQUIRED FOR M4 Phase 1**
- ❌ Node Agents (0%) - **REQUIRED FOR M4 Phase 1**  
- ❌ Transport Layer (0%) - **REQUIRED FOR M4 Phase 1**
- ❌ SSH Automation (60%) - Partially complete, needs finishing

**Decision**: 
- **M3 Phases 1-2 MUST complete before M4 Phase 1 can finish**
- M4 Phase 1 can START in parallel (leverage existing M3 infrastructure)
- Plan assumes **2-week M3 acceleration sprint** → Full M3 ready by Feb 3, 2026

---

## Overall Timeline

```
PHASE                  DURATION    DATES            M3 DEPENDENCY    STATUS
─────────────────────────────────────────────────────────────────────────
M3 Acceleration        2 weeks     Jan 20 - Feb 3   Critical         ← BLOCKING M4
M4 Phase 1 (Coord)     2 weeks     Jan 27 - Feb 10  Requires M3 Pfx  IN PLAN
M4 Phase 2 (GitHub)    1.5 weeks   Feb 3 - Feb 17   Low              IN PLAN
M4 Phase 3 (Install)   1 week      Feb 17 - Feb 24  Low              IN PLAN
M4 Phase 4 (Polish)    1 week      Feb 24 - Mar 3   None             IN PLAN
─────────────────────────────────────────────────────────────────────────
TOTAL                  ~6 weeks    Jan 20 - Mar 3
```

---

## M3 Acceleration Plan (BLOCKING CRITICAL PATH)

### Why M3 Must Finish First

M4 Phase 1 requires:
- ✅ Configuration System (EXISTS - 100%)
- ✅ Local Providers (EXISTS - 50%, needs fixes)
- ❌ Transport Layer (MUST BUILD - 0%)
- ❌ Coordination Server (MUST BUILD - 0%)
- ❌ Node Agents (MUST BUILD - 0%)
- ⚠️ SSH Automation (MUST FINISH - 60%)

### M3 Work Streams (Parallel Delegation)

**STREAM 1: Coordination Server (Backend Engineer)**
- Duration: 2 weeks
- Primary deliverable: HTTP API + workspace management
- Subagent: `sisyphus-junior-high` (complex backend work)
- Tasks:
  - HTTP server on port 3001
  - Workspace CRUD endpoints
  - User registration & SSH key management
  - Workspace status tracking
  - Health check endpoint

**STREAM 2: Node Agents & Transport (Systems Engineer)**
- Duration: 2 weeks  
- Primary deliverable: SSH-based agent + command dispatch
- Subagent: `sisyphus-junior` or `explorer`
- Tasks:
  - Node agent binary (runs on LXC/Docker nodes)
  - SSH command protocol design
  - Command execution & result reporting
  - Container lifecycle management

**STREAM 3: SSH Automation Completion (Backend Engineer)**
- Duration: 1 week
- Primary deliverable: Automated SSH key management
- Subagent: `sisyphus-junior`
- Tasks:
  - Detect/generate SSH keys
  - Upload to GitHub via `gh ssh-key add`
  - Register with coordination server
  - SSH config management

**STREAM 4: Testing & Verification (QA Specialist)**
- Duration: 2 weeks (parallel)
- Verification: All M3 components work together
- Subagent: `qa-tester` 
- Tests:
  - Unit tests: 90%+ coverage
  - Integration tests: Component interactions
  - E2E tests: Full remote workflow

### M3 Success Criteria (Before M4 Phase 1 Finish)

- [ ] Coordination server running and responsive (99.9% uptime)
- [ ] Node agents can be registered and controlled
- [ ] Container creation/deletion works remotely
- [ ] SSH authentication functional
- [ ] All tests passing with 90%+ coverage
- [ ] No architectural violations or anti-patterns

**Verification**: Sisyphus (you) will review code, architecture, test coverage

---

## M4 Phase Breakdown & Delegation

### PHASE 1: Coordination Server Foundation (Weeks 1-2 of M4)
**Dates**: Jan 27 - Feb 10  
**Dependencies**: M3 Transport Layer + Coordination Server (from M3 work)  
**Primary Goal**: Build staging environment infrastructure

#### Phase 1 Work Streams (Parallel Delegation)

**STREAM 1.1: API Specification & Workspace Management**
- **Assigned to**: Backend Engineer (sisyphus-junior-high)
- **Effort**: 80 hours
- **Key Deliverables**:
  - Extended `/api/v1/users/register` with GitHub OAuth
  - Extended `/api/v1/workspaces/create` with GitHub repo fork
  - `/api/v1/workspaces/{id}/services` for service discovery
  - `/api/v1/workspaces/{id}/ssh` for SSH connection info
  - Port allocation & forwarding logic
  - Workspace status polling endpoints
- **Testing**: Unit + Integration tests (90%+ coverage)
- **Success Criteria**:
  - All endpoints respond correctly
  - Port allocation works (2222-2299 range)
  - SSH connections established successfully

**STREAM 1.2: LXC Driver Integration**
- **Assigned to**: Systems Engineer
- **Effort**: 60 hours
- **Key Deliverables**:
  - LXC node agent receiving workspace creation commands
  - Container launch with Ubuntu 22.04 image
  - SSH server setup inside container
  - Service startup from `.nexus/config.yaml`
  - Port forwarding to host
- **Testing**: E2E container lifecycle tests
- **Success Criteria**:
  - Container creates and starts <30s
  - SSH accessible from coordination server
  - Services run inside container

**STREAM 1.3: Database & Metadata Storage**
- **Assigned to**: Backend Engineer  
- **Effort**: 40 hours
- **Key Deliverables**:
  - SQLite schema for workspaces, users, services
  - Workspace state transitions (pending → creating → running → stopped)
  - User metadata (GitHub handle, SSH keys)
  - Service health tracking
  - Data migration support
- **Testing**: Database schema tests, state transition tests
- **Success Criteria**:
  - Correct data persistence
  - State transitions work as designed
  - No data corruption

**STREAM 1.4: Testing & Documentation**
- **Assigned to**: QA Specialist
- **Effort**: 60 hours
- **Key Deliverables**:
  - API integration tests
  - Container lifecycle E2E tests
  - Multi-user concurrent workspace tests
  - SSH latency benchmarks
  - Coordination server stability tests (99.9% uptime)
  - API documentation (OpenAPI spec)
- **Testing**: Full test suite execution
- **Success Criteria**:
  - All tests passing
  - 90%+ code coverage
  - API docs complete

#### Phase 1 Success Criteria

- [ ] Coordination server stable (99.9% uptime over 24h test)
- [ ] 10+ concurrent workspaces without interference
- [ ] SSH latency <100ms
- [ ] Container startup <30 seconds
- [ ] All tests passing, 90%+ coverage
- [ ] API spec complete & documented

#### Phase 1 Verification (CRITICAL CHECKPOINT)

**Sisyphus Review**:
1. Code review for architecture correctness
2. Test coverage analysis
3. E2E workflow validation
4. Performance benchmarking
5. Security audit for SSH & authentication

**Approval Required**: Phase 1 must pass all criteria before Phase 2 starts

---

### PHASE 2: GitHub CLI Integration (Weeks 2-3 of M4)
**Dates**: Feb 3 - Feb 17  
**Dependencies**: Phase 1 coordination server  
**Primary Goal**: Automate user onboarding (auth, SSH keys, workspace creation)

#### Phase 2 Work Streams

**STREAM 2.1: CLI Auth & SSH Commands**
- **Assigned to**: Backend Engineer (sisyphus-junior-high)
- **Effort**: 80 hours
- **Key Deliverables**:
  - `nexus auth github` - Orchestrate `gh auth login`
  - Extract & cache GitHub username & token
  - `nexus ssh setup` - Detect/generate SSH keys
  - Upload public key to GitHub via `gh ssh-key add`
  - Handle duplicate key detection & errors
  - Store SSH key fingerprint with coordination server
- **Testing**: Integration tests with GitHub API mocking
- **Success Criteria**:
  - GitHub auth flow works end-to-end
  - SSH keys upload to GitHub successfully
  - Error handling for network issues, duplicate keys

**STREAM 2.2: Workspace Creation & Fork Logic**
- **Assigned to**: Backend Engineer
- **Effort**: 80 hours
- **Key Deliverables**:
  - `nexus workspace create <repo>` command
  - Repository ownership check via GitHub API
  - Fork logic via `gh repo fork` (if not owned)
  - Clone logic via `git clone` (if already owned)
  - Load `.nexus/config.yaml` from repository
  - Generate sensible defaults (if config missing)
  - Call coordination server to create workspace
  - Poll status until ready
- **Testing**: Integration with GitHub API, coordination server
- **Success Criteria**:
  - Works with owned repos (direct clone)
  - Works with unowned repos (automatic fork)
  - Config loading & defaults work
  - Status polling & completion detection

**STREAM 2.3: Editor Detection & Deep Links**
- **Assigned to**: CLI/UX Engineer  
- **Effort**: 40 hours
- **Key Deliverables**:
  - `nexus workspace connect <name>` command
  - Detect installed editors (Cursor > VS Code > Vim)
  - Generate SSH deep links for Cursor & VS Code
  - Launch editor with remote SSH connection
  - Display summary with service URLs
  - Handle launch failures gracefully
- **Testing**: Editor launch tests (Docker mocked editors)
- **Success Criteria**:
  - Editor detection works
  - Deep links generated correctly
  - Editor launches successfully
  - Fallback to SSH command if needed

**STREAM 2.4: Service Discovery & Display**
- **Assigned to**: CLI/UX Engineer
- **Effort**: 40 hours
- **Key Deliverables**:
  - `nexus workspace services <name>` command
  - Query coordination server for service list
  - Display services with ports & URLs
  - Pretty formatting with colors & emojis
  - Handle service failures gracefully
- **Testing**: Service discovery tests
- **Success Criteria**:
  - Service list accurate
  - URLs correct for local forwarding
  - Display clear and readable

**STREAM 2.5: Testing & Documentation**
- **Assigned to**: QA Specialist
- **Effort**: 60 hours
- **Key Deliverables**:
  - Full workflow tests (auth → fork → create → connect)
  - Error scenario tests (GitHub auth failure, network issues)
  - Multi-OS testing (macOS, Linux)
  - Documentation of all new commands
  - User guides & troubleshooting
- **Testing**: E2E workflow tests
- **Success Criteria**:
  - Full workflow works <3 minutes
  - Clear error messages for all failures
  - Works on macOS & Linux

#### Phase 2 Success Criteria

- [ ] `nexus auth github` - Complete auth flow works
- [ ] `nexus ssh setup` - SSH keys to GitHub automated
- [ ] `nexus workspace create` - Workspace creation fully automated
- [ ] `nexus workspace connect` - Editor opens with working SSH
- [ ] `nexus workspace services` - Service discovery working
- [ ] Full workflow: install → auth → create → connect <3 minutes
- [ ] All tests passing, 90%+ coverage

#### Phase 2 Verification (CRITICAL CHECKPOINT)

**Sisyphus Review**:
1. End-to-end user flow validation (fresh user test)
2. Error handling for all edge cases
3. GitHub API integration robustness
4. Multi-OS compatibility verification
5. Performance under typical load

**Approval Required**: Phase 2 must pass all criteria before Phase 3 starts

---

### PHASE 3: One-Line Install Script (Week 3 of M4)
**Dates**: Feb 17 - Feb 24  
**Dependencies**: Phase 2 CLI commands  
**Primary Goal**: Zero-friction user onboarding

#### Phase 3 Work Streams

**STREAM 3.1: Install Script Development**
- **Assigned to**: Systems Engineer
- **Effort**: 40 hours
- **Key Deliverables**:
  - Bash install script (install.sh)
  - Binary download from GitHub releases (signed)
  - Dependency checking (git, ssh, gh)
  - Platform detection (macOS, Linux variants)
  - Installation to ~/.local/bin
  - Config initialization
  - Kick off setup workflow
- **Testing**: Script testing on multiple platforms
- **Success Criteria**:
  - Script runs on macOS 11+, Ubuntu 18.04+, RHEL 8+
  - Graceful failure with helpful messages
  - All dependencies checked

**STREAM 3.2: Binary Signing & Distribution**
- **Assigned to**: DevOps/Systems Engineer
- **Effort**: 30 hours
- **Key Deliverables**:
  - Build artifacts for macOS (Intel + ARM) & Linux
  - Code signing for macOS
  - Checksum verification
  - GitHub release automation
  - CDN distribution (if applicable)
- **Testing**: Binary verification tests
- **Success Criteria**:
  - Binaries buildable & runnable
  - Checksums verifiable
  - GitHub releases automated

**STREAM 3.3: Testing & Documentation**
- **Assigned to**: QA Specialist
- **Effort**: 50 hours
- **Key Deliverables**:
  - Platform-specific installation tests
  - Dependency detection tests
  - Failure scenario tests
  - Installation documentation
  - Troubleshooting guide
- **Testing**: E2E installation tests on CI
- **Success Criteria**:
  - Installation works on all supported platforms
  - Clear error messages for failures

#### Phase 3 Success Criteria

- [ ] Install script works on macOS, Ubuntu, RHEL
- [ ] Binary signed and verifiable
- [ ] Installation <1 minute
- [ ] All error cases handled gracefully
- [ ] Documentation complete

#### Phase 3 Verification

**Sisyphus Review**:
1. Installation process validation (clean machines)
2. Security review of binary signing
3. Documentation completeness
4. Multi-platform testing results

---

### PHASE 4: Polish, Testing & Launch (Weeks 4-5 of M4)
**Dates**: Feb 24 - Mar 3  
**Dependencies**: All previous phases  
**Primary Goal**: Production-ready system

#### Phase 4 Work Streams

**STREAM 4.1: Comprehensive Testing**
- **Assigned to**: QA Specialist  
- **Effort**: 100 hours
- **Key Deliverables**:
  - Full E2E test suite (install → workspace → services → edit)
  - Error scenario comprehensive testing
  - Load testing (10+ concurrent workspaces)
  - Performance benchmarking
  - Reliability testing (24h+ uptime)
  - Multi-user scenarios
  - Service failure recovery
  - Documentation of test results
- **Testing**: Execute full test matrix
- **Success Criteria**:
  - All tests passing consistently
  - 95%+ pass rate on E2E tests
  - Performance targets met

**STREAM 4.2: UX Polish & Error Handling**
- **Assigned to**: CLI/UX Engineer
- **Effort**: 60 hours
- **Key Deliverables**:
  - Refined progress messaging (spinners, percentages)
  - Comprehensive error messages with next steps
  - Command help text & examples
  - Bash completion scripts
  - Log output formatting
  - Color & emoji consistency
- **Testing**: UX feedback from test users
- **Success Criteria**:
  - Error messages are clear & actionable
  - Progress is visible throughout workflow
  - All commands have help text

**STREAM 4.3: Security Review**
- **Assigned to**: Systems Engineer
- **Effort**: 40 hours
- **Key Deliverables**:
  - SSH key handling security audit
  - GitHub token storage security
  - API authentication review
  - Network communication encryption
  - Access control verification
  - Credentials protection review
- **Testing**: Security audit execution
- **Success Criteria**:
  - No security vulnerabilities identified
  - All credentials properly protected

**STREAM 4.4: Documentation**
- **Assigned to**: Documentation Writer (or Backend Engineer)
- **Effort**: 60 hours
- **Key Deliverables**:
  - Complete user guide
  - Installation guide (all platforms)
  - Troubleshooting guide
  - Architecture documentation
  - API documentation (OpenAPI)
  - Developer guide for contributors
  - Migration guide from M3
  - Example repositories with .nexus/config.yaml
- **Testing**: Documentation review & testing
- **Success Criteria**:
  - Docs cover all commands
  - Examples are runnable
  - Troubleshooting covers common issues

**STREAM 4.5: Final Integration & Launch**
- **Assigned to**: Backend Engineer (Lead)
- **Effort**: 80 hours
- **Key Deliverables**:
  - Final integration testing
  - Performance optimization
  - Build & release preparation
  - Staging environment setup
  - Monitoring & alerting setup
  - Rollout plan execution
  - Beta user outreach
- **Testing**: Staging environment validation
- **Success Criteria**:
  - Staging environment ready for users
  - Monitoring configured
  - Rollout plan ready

#### Phase 4 Success Criteria

- [ ] All E2E tests passing consistently
- [ ] 95%+ test pass rate
- [ ] Performance: <3 min full workflow, <30s container startup, <100ms SSH
- [ ] Load test: 10+ concurrent workspaces stable
- [ ] Security audit passed
- [ ] Documentation complete & accurate
- [ ] Staging environment running
- [ ] Beta testing plan ready

#### Phase 4 Verification (FINAL CHECKPOINT)

**Sisyphus Review**:
1. Comprehensive test results validation
2. Performance benchmark analysis
3. Security audit findings review
4. Documentation completeness check
5. Staging environment readiness
6. Go/no-go decision for public launch

---

## Delegation Strategy

### Subagent Assignments

| Role | Subagent Type | Key Responsibilities | Phases |
|------|---------------|---------------------|--------|
| **Backend Lead** | `sisyphus-junior-high` | Coordination Server, API, GitHub integration | 1.1, 1.3, 2.1, 2.2, 4.5 |
| **Systems Engineer** | `sisyphus-junior` | Node agents, SSH, LXC integration, install script | M3 stream 2, 1.2, 3.1, 4.3 |
| **CLI/UX Engineer** | `sisyphus-junior` | CLI commands, editor integration, UX polish | 2.3, 2.4, 4.2 |
| **QA Specialist** | `qa-tester` | Testing, verification, documentation | 1.4, 2.5, 3.3, 4.1 |
| **You (Sisyphus)** | Strategic Lead | Direction oversight, architecture review, final approval | All phases |

### How Delegation Works

**For Each Task**:
1. **Clear scope** - Specific deliverable, success criteria, testing requirements
2. **Context provided** - Code examples, API specs, configuration details
3. **Async execution** - Subagent works independently
4. **Verification gate** - You review code & tests before approval
5. **Iterative refinement** - Feedback loop until criteria met

**Example Delegation Prompt**:
```
TASK: Implement Coordination Server HTTP API
EXPECTED OUTCOME: 
  - HTTP server on port 3001
  - POST /api/v1/workspaces/create endpoint
  - GET /api/v1/workspaces/{id}/status endpoint
  - DELETE /api/v1/workspaces/{id} endpoint
  - All endpoints return JSON as specified in M4_API_SPECIFICATION.md
REQUIRED SKILLS: Backend Go engineering, HTTP APIs, database design
TESTING:
  - Unit tests for workspace state transitions
  - Integration tests for all endpoints
  - 90%+ code coverage
  - Error handling for all edge cases
MUST DO:
  - Follow existing code patterns in pkg/coordination
  - Use chi/mux for HTTP routing (existing dependency)
  - SQLite for persistence (embedded)
  - All errors wrapped with context
MUST NOT DO:
  - Use interface{} (use concrete types)
  - Modify existing API contracts without approval
  - Skip tests or documentation
CONTEXT:
  - Review M4_API_SPECIFICATION.md for exact contract
  - Review M4_TECHNICAL_SPECIFICATION.md for architecture
  - Look at existing pkg/provider/ for pattern examples
  - Database schema: docs/planning/M4/specs/schema.sql
```

---

## Verification & Quality Gates

### Per-Phase Gate Criteria

Each phase must pass **all** criteria before moving to next:

**Phase 1 Gate** (Coordination Server):
- ✅ Code review passed (architecture, patterns, quality)
- ✅ Unit tests: 90%+ coverage
- ✅ Integration tests: All passing
- ✅ E2E tests: Container lifecycle works
- ✅ Performance: <30s container startup, <100ms SSH latency
- ✅ Security: SSH & auth implementation reviewed
- ✅ API contract matches spec exactly

**Phase 2 Gate** (GitHub Integration):
- ✅ End-to-end user flow works (<3 minutes)
- ✅ All CLI commands implemented & tested
- ✅ Error handling for network issues, GitHub API failures
- ✅ Multi-OS compatibility verified (macOS, Linux)
- ✅ GitHub CLI orchestration robust (handles missing tools, auth issues)

**Phase 3 Gate** (Install Script):
- ✅ Installation works on all supported platforms
- ✅ Binary signed and verifiable
- ✅ Dependency detection working
- ✅ Clean uninstall possible

**Phase 4 Gate** (Polish & Launch):
- ✅ E2E test suite: 95%+ pass rate
- ✅ Load test: 10+ concurrent workspaces stable
- ✅ Performance benchmarks meet targets
- ✅ Security audit passed
- ✅ Documentation complete & accurate
- ✅ Staging environment operational

### Code Review Checklist

Every PR requires review against:
- ✅ Follows existing code patterns (no cowboys)
- ✅ No `interface{}` without justification
- ✅ All errors wrapped with `fmt.Errorf`
- ✅ All tests passing, coverage ≥90%
- ✅ Follows conventional commits
- ✅ PR description references checklist item
- ✅ No security issues (credentials, auth, SSH)

---

## Risk Assessment & Mitigation

### High-Risk Areas

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|-----------|
| M3 coordination server not ready on time | Blocks M4 Phase 1 | MEDIUM | Parallel M3 acceleration sprint (2 weeks) |
| GitHub API rate limiting or changes | Phase 2 delays | LOW | Implement retry logic, cache tokens |
| SSH connection stability issues | SSH latency >100ms | MEDIUM | Connection pooling, timeout tuning |
| LXC/Docker performance on staging | Container startup >30s | LOW | Profiling, image optimization |
| Editor deep links don't work | User frustration | LOW | Fallback to SSH command |
| Multi-workspace interference | Data corruption | LOW | Comprehensive isolation tests |

### Contingency Plans

**If Phase 1 Delays**:
- Reduce Phase 1 scope: Docker only, defer LXC
- Extend Phase 1 to 3 weeks
- Parallelize Phase 2 work on CLI while Phase 1 finishes

**If GitHub Integration Proves Complex**:
- Implement manual workspace creation first (without fork/clone)
- Defer automatic fork/clone to Phase 2.5
- Focus on SSH setup automation first

**If Performance Issues Found**:
- Profiling & optimization focused work
- May require architecture adjustments
- Budget extra week in Phase 4

---

## Success Metrics

### User Experience Metrics
- [ ] First workspace creation: <3 minutes (target)
- [ ] Container startup: <30 seconds
- [ ] SSH latency: <100ms
- [ ] Editor launch: <10 seconds
- [ ] Error clarity: 100% of errors have actionable messages

### Technical Metrics
- [ ] Server uptime: 99.9% (or better)
- [ ] Concurrent workspaces: 10+ without interference
- [ ] Code coverage: 90%+ on new code
- [ ] Test pass rate: 95%+ (E2E tests)

### Adoption Metrics
- [ ] Successful beta testing with 20+ users
- [ ] 0 critical bugs in first week of launch
- [ ] <5 support tickets per day (target)

---

## Document Index

### M4 Specification Documents
- `M4_OVERVIEW.md` - Executive summary
- `M4_USER_FLOW_SPECIFICATION.md` - 7-step user journey
- `M4_TECHNICAL_SPECIFICATION.md` - Architecture & APIs
- `api/M4_API_SPECIFICATION.md` - REST API reference
- `checklists/M4_IMPLEMENTATION_CHECKLIST.md` - Task list

### This Execution Plan
- `M4_EXECUTION_PLAN.md` - This document
- Work item assignments (created as PR descriptions)
- Phase gate reviews (documented in commit messages)

---

## Next Steps

1. **Today (Jan 17)**: Finalize M3 acceleration plan, assign subagents
2. **Jan 20**: M3 acceleration sprint starts
3. **Jan 27**: M4 Phase 1 work begins (parallel with M3 finish)
4. **Feb 3**: M3 complete, M4 Phase 2 starts
5. **Feb 24**: Phase 4 launches, staging environment opens to beta users
6. **Mar 3**: Public launch

---

**Created by**: Sisyphus  
**Status**: Ready for Subagent Delegation  
**Last Updated**: January 17, 2026
