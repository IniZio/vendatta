# M4 Execution Plan - Executive Summary & Delegation Strategy

**Version**: 1.0  
**Status**: Ready for Immediate Execution  
**Created**: January 17, 2026  
**Target Completion**: March 3, 2026 (6 weeks)  
**Investment**: ~940 engineering hours (~$35-50k)

---

## CRITICAL BLOCKERS TO RESOLVE FIRST

Before delegating M4 Phase 1 work, **2 blocking issues must be fixed immediately**:

### BLOCKER 1: Failing E2E Tests
**Issue**: 5 E2E test files failing (timeout, environment setup failures)
- `e2e/lifecycle_test.go` - Failing
- `e2e/m3_verification_test.go` - Failing
- `e2e/transport_test.go` - Failing
- `e2e/transport_local_test.go` - Failing
- `e2e/testenv.go` - Has useful utilities but tests won't run

**Root Causes**:
- Docker/LXC environment not properly configured on test machines
- SSH keys not setup for test runner
- LXC memory limits causing container creation failures
- Test environment (testenv.go) not fully initialized

**Resolution** (Priority: IMMEDIATE):
- [ ] Diagnose why Docker tests fail (check Docker daemon)
- [ ] Fix SSH key generation in test environment
- [ ] Resolve LXC memory/resource issues
- [ ] Update CI/CD environment to include proper setup
- [ ] Verify all e2e tests pass before Phase 1 starts

**Owner**: Systems Engineer  
**Effort**: 8-12 hours  
**Deadline**: Jan 20, 2026

---

### BLOCKER 2: Makefile Test Targets
**Issue**: Makefile declares test targets but implementations missing
```makefile
# These exist but have no implementation:
test:
	@echo "Running tests..."
test-unit:
	@echo "Running unit tests..."
test-integration:
	@echo "Running integration tests..."
test-e2e:
	@echo "Running e2e tests..."
```

**Impact**: CI/CD pipeline can't run proper test categories, TDD workflow broken

**Resolution** (Priority: HIGH):
- [ ] Implement `make test` - Run all tests with coverage
- [ ] Implement `make test-unit` - Unit tests only
- [ ] Implement `make test-integration` - Integration tests
- [ ] Implement `make test-e2e` - E2E tests
- [ ] Add `make test-coverage` - Generate coverage report
- [ ] Add targets to CI/CD pipeline

**Owner**: Systems Engineer  
**Effort**: 4-6 hours  
**Deadline**: Jan 22, 2026

---

## CURRENT TEST COVERAGE ASSESSMENT

### Summary
- **Total Test Files**: 27 files across 13 packages
- **Test Framework**: testify v1.11.1 (excellent)
- **Coverage Range**: 9.6% - 88.3%

### By Coverage Level

**Strong (>70%)**:
- ✅ plugins/registry: 88.3%
- ✅ lock/manager: 73.3%
- ✅ worktree: 71.7%

**Medium (40-70%)**:
- ⚠️ auth: 57.1%
- ⚠️ provider/qemu: 56.3%
- ⚠️ metrics: 61.9%
- ⚠️ config: 45.0%
- ⚠️ provider/docker: 34.8%

**Weak (<40%)**:
- ❌ provider/lxc: 6.5% (CRITICAL - needed for M4)
- ❌ templates: 33.6%
- ❌ coordination: 22.8% (CRITICAL - needed for M4)
- ❌ agent: 11.3% (CRITICAL - needed for M4)
- ❌ cmd/nexus: 9.6%

### Action Items for M4
For Phase 1 work, **MUST achieve 90%+ coverage** on new code:
- Extended coordination server endpoints
- New database models & DAOs
- Service discovery logic
- Port allocation & forwarding

**This will be enforced at each Phase gate review.**

---

## M4 EXECUTION TIMELINE

### Pre-Phase 1: Blocker Resolution (Jan 17-22)
```
✓ DONE: Strategic execution planning
⏳ TODO: Fix failing E2E tests (due Jan 20)
⏳ TODO: Implement Makefile targets (due Jan 22)
✓ VERIFICATION: All tests passing, CI/CD ready
```

### PHASE 1: Coordination Server Foundation (Jan 27 - Feb 10)
**Duration**: 2 weeks (2 parallel weeks)  
**Goal**: Extend coordination server for user registration, workspace creation, service discovery  
**Success Criteria**: Server stable 99.9%, 10+ concurrent workspaces, all tests >90% coverage

**4 Parallel Work Streams**:

| Stream | Lead | Focus | Hours | Tests |
|--------|------|-------|-------|-------|
| 1.1: API Endpoints | Backend Engineer | GitHub user registration, workspace CRUD, service discovery | 80 | Integration + Unit |
| 1.2: LXC Integration | Systems Engineer | Container lifecycle, SSH setup, service startup | 60 | E2E container tests |
| 1.3: Database | Backend Engineer | SQLite schema, DAOs, state management | 40 | Database + State |
| 1.4: Testing | QA Specialist | Unit + Integration + E2E + Performance | 60 | Full test suite |

**Verification Gate 1** (Feb 10): Sisyphus reviews code, architecture, tests, performance

---

### PHASE 2: GitHub CLI Integration (Feb 3 - Feb 17)
**Duration**: 1.5 weeks (overlaps Phase 1 end)  
**Goal**: Automate user onboarding (GitHub auth, SSH keys, workspace creation)  
**Success Criteria**: Full workflow <3 minutes, all errors handled, works macOS+Linux

**5 Parallel Work Streams**:

| Stream | Lead | Focus | Hours | Tests |
|--------|------|-------|-------|-------|
| 2.1: GitHub Auth | Backend Engineer | `nexus auth github`, token storage | 80 | Mock GitHub API |
| 2.2: Fork/Clone | Backend Engineer | `nexus workspace create`, repo ownership | 80 | GitHub API mocking |
| 2.3: Editor Launch | CLI Engineer | `nexus workspace connect`, deep links | 40 | Editor mock tests |
| 2.4: Service Discovery | CLI Engineer | `nexus workspace services`, port display | 40 | Service list tests |
| 2.5: E2E Testing | QA Specialist | Full workflow, error scenarios, multi-OS | 60 | E2E workflow tests |

**Verification Gate 2** (Feb 17): Sisyphus validates end-to-end workflow with fresh user test

---

### PHASE 3: One-Line Install Script (Feb 17 - Feb 24)
**Duration**: 1 week  
**Goal**: Zero-friction user bootstrap  
**Success Criteria**: Works on macOS, Linux, RHEL; signed binary; documented

**3 Work Streams**:

| Stream | Lead | Focus | Hours | Tests |
|--------|------|-------|-------|-------|
| 3.1: Install Script | Systems Engineer | Bash script, dependency checking, config init | 40 | Platform testing |
| 3.2: Binary Distribution | Systems Engineer | Multi-platform builds, code signing, GitHub releases | 30 | Build automation |
| 3.3: Testing | QA Specialist | Installation on all platforms, documentation | 50 | Cross-platform tests |

**Verification Gate 3** (Feb 24): Sisyphus validates installation process

---

### PHASE 4: Polish, Testing & Launch (Feb 24 - Mar 3)
**Duration**: 1-2 weeks  
**Goal**: Production-ready system, beta launch  
**Success Criteria**: 95%+ E2E pass rate, performance targets met, security audit passed

**5 Work Streams**:

| Stream | Lead | Focus | Hours | Tests |
|--------|------|-------|-------|-------|
| 4.1: Comprehensive Testing | QA Specialist | Load tests, reliability, performance benchmarks | 100 | Full test matrix |
| 4.2: UX Polish | CLI Engineer | Progress messages, error clarity, help text | 60 | UX validation |
| 4.3: Security Review | Systems Engineer | SSH, GitHub tokens, API auth, credentials | 40 | Security audit |
| 4.4: Documentation | Any Engineer | User guide, troubleshooting, API docs | 60 | Doc review |
| 4.5: Final Integration | Backend Engineer | Build, release, staging setup, monitoring | 80 | Staging validation |

**Verification Gate 4** (Mar 3): **FINAL GO/NO-GO** for public launch

---

## TEAM & EFFORT ALLOCATION

### Engineering Team (Recommended)
```
Backend Engineer (Senior)    - 360 hours (9 weeks full-time)
  ├─ Phase 1: API endpoints, database
  ├─ Phase 2: GitHub integration
  └─ Phase 4: Final integration & launch

Systems Engineer (Mid-level) - 170 hours (4 weeks full-time)
  ├─ Phase 1: LXC integration
  ├─ Phase 3: Install script, binary distribution
  └─ Phase 4: Security review

CLI/UX Engineer (Mid-level)  - 140 hours (3.5 weeks full-time)
  ├─ Phase 2: Editor detection, service discovery
  └─ Phase 4: UX polish

QA Specialist (Mid-level)    - 270 hours (6.5 weeks full-time)
  ├─ All phases: Comprehensive testing
  └─ Phase 4: Load testing & reliability
```

**Total**: ~4.5 full-time engineers for 6 weeks = **~940 engineering hours**

---

## QUALITY GATES & APPROVAL PROCESS

### Gate 1: Phase 1 Completion (Feb 10)
**Sisyphus Review Duration**: 1 hour

**Verification Checklist**:
- [ ] Code review: Architecture patterns, no anti-patterns
- [ ] Test coverage: 90%+ on all new code
- [ ] E2E test: Full container lifecycle works
- [ ] Performance: <30s startup, <100ms SSH latency
- [ ] Security: SSH & auth implementation reviewed
- [ ] Concurrent workspaces: 10+ without interference

**Decision**: PASS / CONDITIONAL / FAIL

---

### Gate 2: Phase 2 Completion (Feb 17)
**Sisyphus Review Duration**: 2 hours

**Verification Checklist**:
- [ ] Fresh user test: Install → Auth → Create → Connect
- [ ] Workflow timing: <3 minutes total
- [ ] All CLI commands working
- [ ] Error handling: Bulletproof GitHub API failures
- [ ] Multi-OS: Verified on macOS & Linux
- [ ] Test coverage: 90%+ on new code

**Decision**: PASS / CONDITIONAL / FAIL

---

### Gate 3: Phase 3 Completion (Feb 24)
**Sisyphus Review Duration**: 1 hour

**Verification Checklist**:
- [ ] Installation script: Works on macOS, Ubuntu, RHEL
- [ ] Binary security: Signed and verifiable
- [ ] Dependency detection: Robust error handling
- [ ] Documentation: Complete & accurate

**Decision**: PASS / CONDITIONAL / FAIL

---

### Gate 4: FINAL LAUNCH APPROVAL (Mar 3)
**Sisyphus Review Duration**: 3 hours

**Verification Checklist**:
- [ ] E2E test suite: 95%+ pass rate, all scenarios
- [ ] Load testing: 10+ concurrent workspaces stable
- [ ] Performance: <3min workflow, <30s startup, <100ms SSH
- [ ] Security audit: All findings addressed
- [ ] Documentation: User guide, troubleshooting, API docs complete
- [ ] Staging environment: Ready for beta users

**Decision**: **GO / NO-GO for public launch**

---

## DELEGATION PROMPT TEMPLATE

Each engineer receives a detailed prompt like this:

```
TASK: Implement POST /api/v1/users/register-github Endpoint

EXPECTED OUTCOME:
  - HTTP endpoint responds to requests
  - Accepts: github_username, github_id, ssh_pubkey, ssh_pubkey_fingerprint
  - Returns: user_id, registered_at, workspaces array
  - Status 201 on success, 400/409 on validation/duplicate errors
  - All responses match M4_API_SPECIFICATION.md exactly

REQUIRED SKILLS:
  - Go HTTP API development
  - SQLite database operations
  - JSON request/response handling

TESTING REQUIREMENTS:
  - Unit tests: Valid/invalid input handling
  - Integration tests: Database persistence
  - Error tests: Duplicate key (409), validation (400)
  - 90%+ code coverage

MUST DO:
  - Follow existing code patterns in pkg/coordination/
  - Use chi/mux routing (existing dependency)
  - Wrap errors with fmt.Errorf for context
  - Database schema per M4_TECHNICAL_SPECIFICATION.md
  - Request validation before database operations
  - Idempotent operations where possible

MUST NOT DO:
  - Use interface{} (use concrete types)
  - Skip error handling
  - Make assumptions about API contract
  - Modify existing endpoints without approval
  - Leave TODOs without issue tracking

CONTEXT:
  - Review: M4_API_SPECIFICATION.md
  - Review: M4_TECHNICAL_SPECIFICATION.md
  - Reference: pkg/coordination/handlers.go (existing patterns)
  - Schema: See STREAM 1.3 (Database design)

DELIVERABLES:
  - Implementation in pkg/coordination/handlers.go
  - Tests in pkg/coordination/handlers_test.go
  - Updated API documentation
  - PR with reference to Phase 1.1 checklist item
```

---

## SUCCESS METRICS & ACCEPTANCE

### Must-Have Metrics (Cannot Ship Without)
```
✓ First workspace creation: <3 minutes (user flow)
✓ Container startup: <30 seconds (technical)
✓ SSH latency: <100ms (technical)
✓ Server uptime: 99.9% over 48h test (reliability)
✓ Code coverage: 90%+ on all new code (quality)
✓ All gates passed: Phase 1-4 approvals required (governance)
✓ Zero critical bugs in beta: First week (quality)
```

### Nice-to-Have Metrics
```
◇ E2E test pass rate: 95%+
◇ Support tickets/day: <5
◇ Beta user satisfaction: 4.5+/5.0
◇ Install script success rate: 99%+
```

---

## RISK MITIGATION

### Top 3 Risks

| Risk | Impact | Mitigation |
|------|--------|-----------|
| GitHub API integration complexity | Delays Phase 2 | Early prototyping, comprehensive error handling |
| SSH connection stability | Test failures, poor UX | Connection pooling, retry logic, extensive testing |
| LXC performance on staging | Container startup >30s | Profiling, image optimization, load testing |

### Contingency Plans
- **Phase 1 delays**: Reduce scope to Docker only, extend to 3 weeks
- **GitHub API issues**: Implement manual workspace creation first, fork/clone as enhancement
- **Performance issues**: Extra optimization week in Phase 4, may push launch by 1 week

---

## SUMMARY TABLE

| Aspect | Status | Details |
|--------|--------|---------|
| **Current Codebase** | Strong (⭐⭐⭐⭐) | 64 Go files, 13 packages, solid foundations |
| **Blockers** | 2 CRITICAL | E2E tests failing, Makefile targets missing |
| **Timeline** | 6 weeks | Jan 20 - Mar 3, 2026 |
| **Team Size** | ~4.5 FTE | Backend lead, systems engineer, CLI engineer, QA |
| **Total Effort** | ~940 hours | ~$35-50k engineering cost |
| **Quality Target** | 90%+ coverage | Enforced at every gate |
| **Launch Readiness** | ON TRACK | Assuming blockers fixed by Jan 22 |

---

## NEXT STEPS (IN ORDER)

### Week of Jan 17-22: Blocker Resolution
1. **TODAY (Jan 17)**: Approve this plan
2. **Jan 18-20**: Fix failing E2E tests (Systems Engineer)
3. **Jan 20-22**: Implement Makefile targets (Systems Engineer)
4. **Jan 22**: Verify all tests passing in CI/CD

### Week of Jan 27-Feb 3: Phase 1 Begins
1. **Jan 27**: Kick off Phase 1 parallel work streams
2. **Daily**: 15-min standup, progress tracking
3. **Feb 3-4**: Mid-phase checkpoint (are we on track?)
4. **Feb 10**: Gate 1 review & approval

### Week of Feb 10-17: Phase 1 Polish + Phase 2 Ramp
1. **Feb 10**: Phase 1 gate approval
2. **Feb 10-17**: Phase 2 work begins (overlaps Phase 1 finish)
3. **Feb 17**: Gate 2 review & approval

### Week of Feb 17-24: Phase 3 + Phase 4 Prep
1. **Feb 17**: Phase 3 work begins
2. **Feb 24**: Gate 3 review & approval
3. **Feb 24**: Phase 4 work begins

### Week of Feb 24-Mar 3: Final Push
1. **Feb 24-Mar 3**: Phase 4 execution
2. **Mar 1**: Final testing & validation
3. **Mar 3**: Gate 4 review & **GO/NO-GO decision**

---

## Document References

**Strategic Planning**:
- M4_STRATEGIC_PLAN.md - Detailed 4-phase plan with per-stream breakdown
- M4_EXECUTION_PLAN.md - Comprehensive execution plan (in /docs/planning/M4/)

**M4 Specification** (reference):
- M4_OVERVIEW.md - Executive summary
- M4_USER_FLOW_SPECIFICATION.md - 7-step user journey
- M4_TECHNICAL_SPECIFICATION.md - Architecture & APIs
- M4_API_SPECIFICATION.md - REST API reference
- M4_IMPLEMENTATION_CHECKLIST.md - Task checklist

---

## Approval

**Prepared by**: Sisyphus  
**Status**: Ready for Execution  
**Requires**: Your approval to proceed with delegation  
**Next Action**: Fix blockers, then start Phase 1 Jan 27

---

**Created**: January 17, 2026  
**Last Updated**: January 17, 2026  
**Version**: 1.0 (Ready for Execution)
