# M4 Phase 1: Delegation Task Templates

**Purpose**: Copy-paste ready delegation prompts for each work stream  
**Format**: One task per engineer, includes all 7 required sections  
**Duration**: 2 weeks (Jan 27 - Feb 10)

---

## TASK: Stream 1.1 - API Endpoints (Backend Engineer)

```
TASK: Implement M4 Coordination Server REST API Endpoints

EXPECTED OUTCOME:
  ✓ All user registration endpoints working (POST /api/v1/users/register-github)
  ✓ All workspace CRUD endpoints working
  ✓ All service discovery endpoints working
  ✓ All endpoints return JSON exactly matching M4_API_SPECIFICATION.md
  ✓ Proper HTTP status codes (201, 400, 409, etc.)
  ✓ Request validation before database operations
  ✓ Integration tests for all endpoints (>90% coverage)
  ✓ All tests passing

REQUIRED SKILLS:
  - Go HTTP API development (chi/mux routing)
  - SQLite integration (DAO layer)
  - Request/response JSON handling
  - HTTP status code semantics
  - Integration testing with databases

TESTING REQUIREMENTS:
  - Unit tests: Valid/invalid request handling
  - Unit tests: Error response formatting
  - Unit tests: Status code accuracy
  - Integration tests: API → Database → DAO flow
  - Error scenario tests: 400, 409, 500 responses
  - Minimum 90% code coverage

MUST DO:
  ✓ Follow pkg/coordination/ code patterns exactly
  ✓ Use chi/mux for routing (existing dependency)
  ✓ Wrap errors with fmt.Errorf("context: %w", err) for stack traces
  ✓ Validate ALL inputs before database operations
  ✓ Use database transactions for consistency
  ✓ Return JSON responses matching M4_API_SPECIFICATION.md
  ✓ Implement all 9 endpoints (3 user, 4 workspace, 2 service)
  ✓ HTTP status codes: 200 OK, 201 Created, 400 Bad Request, 409 Conflict, 500 Error
  ✓ Write tests FIRST (TDD), then implement
  ✓ Run tests frequently (not end of day)
  ✓ Use testify/assert and testify/require for assertions

MUST NOT DO:
  ✗ Use interface{} (use concrete types instead)
  ✗ Skip error handling with _ = err
  ✗ Return unvalidated user input in errors
  ✗ Write tests after implementation
  ✗ Use fmt.Println for logging (use structured logging)
  ✗ Make assumptions about API contract - check spec constantly
  ✗ Modify existing endpoints without approval
  ✗ Leave TODO comments without GitHub issue links
  ✗ Commit without running full test suite
  ✗ Mix concerns (API + database in handlers)

CONTEXT & REFERENCES:
  - Read: M4_API_SPECIFICATION.md (exact endpoint contracts)
  - Read: M4_TECHNICAL_SPECIFICATION.md (architecture & design)
  - Read: pkg/coordination/handlers.go (existing patterns)
  - Read: pkg/coordination/types.go (request/response structures)
  - Reference: Stream 1.3 (database DAOs) - coordinate schema
  - Database: See M4_TECHNICAL_SPECIFICATION.md appendix

DELIVERABLES:
  1. pkg/coordination/handlers.go - All 9 endpoints implemented
  2. pkg/coordination/handlers_test.go - Comprehensive tests
  3. pkg/coordination/router.go - HTTP routing setup
  4. pkg/coordination/router_test.go - Router tests
  5. Updated pkg/coordination/types.go - Request/response types
  6. Integration test results (all passing)
  7. Coverage report (90%+ coverage)

VERIFICATION:
  - Run: go test ./pkg/coordination/... -v -race
  - Check: go tool cover -html=coverage.out
  - Verify: Each test name matches "Test<Endpoint><Scenario>"
  - Verify: All error responses have meaningful messages
  - Verify: All status codes match specification

BLOCKERS & ESCALATIONS:
  - If database schema changes needed: Coordinate with Stream 1.3 immediately
  - If API spec conflicts with implementation needs: Document and ask Sisyphus
  - If test infrastructure missing: Use pkg/coordination/*_test.go patterns

START DATE: Jan 27, 2026
DEADLINE: Feb 10, 2026 (Gate 1 review)
ESTIMATED EFFORT: 80 hours (2 weeks full-time)
```

---

## TASK: Stream 1.2 - LXC Container Integration (Systems Engineer)

```
TASK: Implement LXC Container Lifecycle & SSH Configuration for M4

EXPECTED OUTCOME:
  ✓ LXC containers create successfully from ubuntu:22.04 image
  ✓ Containers start and stop on demand
  ✓ SSH server runs inside container
  ✓ SSH port allocated (2222-2299) and forwarded
  ✓ Can SSH into container immediately after creation
  ✓ SSH keys properly configured (isolated from ~/.ssh)
  ✓ Services start in dependency order
  ✓ Health checks pass before marking workspace ready
  ✓ Container cleanup works (no orphaned containers)
  ✓ Performance targets met (create <30s, start <60s)
  ✓ All E2E container tests passing (>90% coverage)

REQUIRED SKILLS:
  - LXC container management
  - SSH server configuration
  - Linux system administration
  - Go concurrent programming
  - E2E testing with containers

TESTING REQUIREMENTS:
  - E2E tests: Full container lifecycle (create → start → stop → delete)
  - E2E tests: SSH connectivity after container creation
  - E2E tests: Port forwarding verification
  - E2E tests: Service startup with dependencies
  - E2E tests: Health check execution
  - Performance tests: Measure startup time
  - Stress tests: 10+ concurrent containers
  - Cleanup tests: Verify no orphaned containers

MUST DO:
  ✓ Create containers from ubuntu:22.04 base image
  ✓ Configure SSH server inside container
  ✓ Allocate SSH port from isolated range (2222-2299)
  ✓ Forward host port → container port 22
  ✓ Generate SSH keys for container (use isolated .nexus/test-ssh/ paths)
  ✓ Parse service configs and execute in order
  ✓ Implement health checks (HTTP, TCP, process)
  ✓ Clean up all container resources on deletion
  ✓ Handle errors gracefully (report status to API)
  ✓ Use pkg/transport/ for SSH operations
  ✓ Coordinate with Stream 1.1 for status updates
  ✓ Write tests FIRST (TDD), then implement
  ✓ Use testify/require and testify/assert for assertions

MUST NOT DO:
  ✗ Use user's ~/.ssh directory (keep test keys isolated)
  ✗ Hard-code SSH keys or ports
  ✗ Leave containers running if cleanup fails
  ✗ Skip health checks
  ✗ Run services without dependency ordering
  ✗ Ignore container startup errors
  ✗ Write tests after implementation
  ✗ Create containers without cleanup handlers
  ✗ Use direct SSH commands (use pkg/transport/ instead)
  ✗ Leave TODO comments without issue tracking

CONTEXT & REFERENCES:
  - Read: M4_TECHNICAL_SPECIFICATION.md (container design)
  - Reference: pkg/provider/lxc/lxc.go (existing LXC implementation)
  - Reference: pkg/provider/docker/docker.go (Docker patterns)
  - Reference: pkg/transport/transport_e2e_test.go (SSH isolation pattern)
  - Reference: e2e/lifecycle_test.go (E2E testing patterns)
  - Port Range: 2222-2299 (allocated by Stream 1.3 PortAllocationDAO)

DELIVERABLES:
  1. pkg/coordination/container.go - Container lifecycle management
  2. pkg/coordination/container_test.go - Comprehensive E2E tests
  3. pkg/coordination/ssh.go - SSH key and server setup
  4. pkg/coordination/ssh_test.go - SSH integration tests
  5. pkg/coordination/port_allocator.go - Port management (work with Stream 1.3)
  6. Updated pkg/coordination/coordination.go - Container integration
  7. Performance benchmark results (container creation timing)

VERIFICATION:
  - Run: go test ./pkg/coordination/... -v -race -timeout 10m
  - Verify: Container creation completes in <30 seconds
  - Verify: SSH port allocated in range 2222-2299
  - Verify: Can SSH into container after creation
  - Verify: Health checks pass before reporting ready
  - Verify: No orphaned containers after deletion
  - Verify: Services start in dependency order

BLOCKERS & ESCALATIONS:
  - If LXC not available in test environment: Use Docker as fallback temporarily
  - If port allocation conflicts: Coordinate with Stream 1.3 on DAO implementation
  - If container creation slow: Profile and optimize image/startup
  - If SSH not working: Check key isolation (.nexus/test-ssh/ paths)

START DATE: Jan 27, 2026
DEADLINE: Feb 10, 2026 (Gate 1 review)
ESTIMATED EFFORT: 60 hours (1.5 weeks full-time)
```

---

## TASK: Stream 1.3 - Database & State Management (Backend Engineer)

```
TASK: Implement M4 SQLite Database Schema & Data Access Objects

EXPECTED OUTCOME:
  ✓ SQLite database with complete schema (5 tables)
  ✓ Database initializes automatically on startup
  ✓ All DAOs (User, Workspace, Service, PortAllocation) working
  ✓ CRUD operations for all entities
  ✓ Transactions prevent race conditions
  ✓ Port allocation atomic (no conflicts)
  ✓ Database migrations support schema evolution
  ✓ Indexes on frequently queried fields
  ✓ Connection pooling for performance
  ✓ All DAO tests passing (>90% coverage)
  ✓ No memory leaks or resource leaks

REQUIRED SKILLS:
  - SQLite database design
  - Go database/sql package
  - Transaction management
  - Data access patterns (DAO)
  - SQL query optimization

TESTING REQUIREMENTS:
  - Unit tests: All CRUD operations
  - Unit tests: Transaction rollback scenarios
  - Unit tests: Constraint violations
  - Integration tests: Multi-operation transactions
  - Concurrency tests: 10+ concurrent writes
  - Cleanup tests: Resource cleanup on errors
  - Performance tests: Query speed benchmarks
  - Migration tests: Schema version management

MUST DO:
  ✓ Create SQLite schema with 5 tables (users, workspaces, services, port_allocations, schema_version)
  ✓ Implement migrations for schema version management
  ✓ Create UserDAO with methods: CreateUser, GetUser, UpdateSSHKey, ListUsers
  ✓ Create WorkspaceDAO with methods: Create, Get, List, Update, Delete, UpdateStatus
  ✓ Create ServiceDAO with methods: Create, Get, List, Update, UpdateStatus
  ✓ Create PortAllocationDAO with methods: Allocate, Release, GetAllocated, List
  ✓ Use database transactions for consistency
  ✓ Implement proper error handling with wrapped errors
  ✓ Add connection pooling (sqlite driver parameters)
  ✓ Create indexes on foreign keys and frequently queried columns
  ✓ Write tests FIRST (TDD), then implement DAOs
  ✓ Use testify/assert and testify/require

MUST NOT DO:
  ✗ Hard-code port ranges (use constants)
  ✗ Skip constraint checks in code
  ✗ Use raw SQL strings (prepare statements)
  ✗ Leave database connections open on error
  ✗ Write tests after implementation
  ✗ Mix DAOs with business logic
  ✗ Skip transaction rollback tests
  ✗ Use interface{} for query results
  ✗ Leave TODO comments without issue tracking
  ✗ Create schema without migration support

CONTEXT & REFERENCES:
  - Read: M4_TECHNICAL_SPECIFICATION.md (database design section)
  - Schema: See database schema in M4_TECHNICAL_SPECIFICATION.md appendix
  - Reference: pkg/lock/manager.go (SQLite patterns in codebase)
  - Port Range: 2222-2299 (coordinate with Stream 1.2)
  - Coordinate: All DAOs used by Stream 1.1 and Stream 1.2

DELIVERABLES:
  1. pkg/coordination/db/schema.sql - Database schema (5 tables)
  2. pkg/coordination/db/migrations.go - Schema version management
  3. pkg/coordination/db/user_dao.go - User data access
  4. pkg/coordination/db/user_dao_test.go - User DAO tests
  5. pkg/coordination/db/workspace_dao.go - Workspace data access
  6. pkg/coordination/db/workspace_dao_test.go - Workspace DAO tests
  7. pkg/coordination/db/service_dao.go - Service data access
  8. pkg/coordination/db/service_dao_test.go - Service DAO tests
  9. pkg/coordination/db/port_allocation_dao.go - Port management
  10. pkg/coordination/db/port_allocation_dao_test.go - Port DAO tests
  11. pkg/coordination/db/types.go - Database model types
  12. pkg/coordination/db/db.go - Database initialization and pooling

VERIFICATION:
  - Run: go test ./pkg/coordination/db/... -v -race
  - Check: go tool cover -html=coverage.out (target 90%+)
  - Verify: Each DAO method has corresponding test
  - Verify: All constraints enforced (unique, foreign key)
  - Verify: Transactions rollback on error
  - Verify: Port allocation never conflicts (atomic)
  - Verify: Database cleanup works (no leaks)

BLOCKERS & ESCALATIONS:
  - If SQLite concurrency issues: Use WAL mode and timeouts
  - If port allocation conflicts: Add UNIQUE constraint and retry logic
  - If schema changes needed from Stream 1.1: Update both files atomically
  - If migration issues: Add explicit down() migration support

START DATE: Jan 27, 2026
DEADLINE: Feb 10, 2026 (Gate 1 review)
ESTIMATED EFFORT: 40 hours (1 week full-time)
```

---

## TASK: Stream 1.4 - Testing & Quality (QA Specialist)

```
TASK: Comprehensive Testing for M4 Phase 1 (Unit, Integration, E2E)

EXPECTED OUTCOME:
  ✓ All unit tests passing (API handlers, DAOs, business logic)
  ✓ All integration tests passing (API → DB flow)
  ✓ All E2E tests passing (full workspace lifecycle)
  ✓ 90%+ code coverage on all new code
  ✓ No flaky tests
  ✓ Performance benchmarks meet targets:
    - Container creation <30 seconds
    - Service startup <60 seconds
    - API response <200ms
  ✓ Load testing with 10+ concurrent workspaces
  ✓ Comprehensive test report for Gate 1 review

REQUIRED SKILLS:
  - Go unit testing (testify framework)
  - Test design and coverage
  - Integration testing with databases
  - Container E2E testing
  - Performance benchmarking
  - Load testing methodology

TESTING REQUIREMENTS (BY CATEGORY):

  UNIT TESTS:
    - API endpoint handlers (success & error paths)
    - All DAO CRUD operations
    - Port allocation logic
    - SSH key generation
    - Status transition validation
    - Request validation
    - Error message formatting

  INTEGRATION TESTS:
    - API endpoint → Database DAO flow
    - Transaction rollback scenarios
    - Concurrent database writes
    - Port allocation conflicts
    - SSH key storage and retrieval
    - Workspace status transitions

  E2E TESTS:
    - Full workspace lifecycle (create → start → stop → delete)
    - SSH connectivity after creation
    - Port forwarding verification
    - Service startup with dependencies
    - Service health checks
    - Concurrent workspace operations
    - Container cleanup verification
    - Error recovery scenarios

  PERFORMANCE TESTS:
    - Container creation timing (<30s target)
    - Service startup timing (<60s target)
    - API response time (<200ms target)
    - Database query performance
    - 10+ concurrent workspace creation

MUST DO:
  ✓ Write tests FIRST (TDD approach)
  ✓ Use testify/assert and testify/require consistently
  ✓ Table-driven tests for multiple scenarios
  ✓ Mock external dependencies (LXC, SSH)
  ✓ Test both success and failure paths
  ✓ Measure code coverage (target 90%+)
  ✓ Create coverage report (HTML visualization)
  ✓ Run tests frequently (not batch at end)
  ✓ Identify and fix flaky tests immediately
  ✓ Document performance benchmarks
  ✓ Test concurrent operations
  ✓ Verify resource cleanup on errors
  ✓ Create comprehensive test summary for Gate 1

MUST NOT DO:
  ✗ Write tests after implementation
  ✗ Skip error path testing
  ✗ Use t.Skip() without reason
  ✗ Create flaky timing-dependent tests
  ✗ Mock too aggressively (miss real issues)
  ✗ Skip performance verification
  ✗ Leave untested edge cases
  ✗ Write vague test names
  ✗ Mix unit and integration tests
  ✗ Assume other streams' code quality

CONTEXT & REFERENCES:
  - Reference: e2e/lifecycle_test.go (E2E test patterns)
  - Reference: e2e/m3_verification_test.go (Verification patterns)
  - Reference: pkg/transport/transport_e2e_test.go (SSH isolation testing)
  - Coordinate: All three streams (1.1, 1.2, 1.3) for integration points
  - Target coverage: 90%+ on all new code in pkg/coordination/

DELIVERABLES:
  1. Test suite results (all passing)
  2. Code coverage report (HTML: coverage.html)
  3. Performance benchmark results
  4. Load test results (10+ concurrent workspaces)
  5. Test documentation (what each test verifies)
  6. Flaky test analysis (if any)
  7. Performance optimization recommendations
  8. Gate 1 review test summary

VERIFICATION:
  - Run: make test (all categories)
  - Check: go tool cover -html=coverage.out (90%+ coverage)
  - Benchmark: go test -bench=. -benchmem ./pkg/coordination/...
  - Load test: Create 10+ workspaces concurrently
  - Verify: No orphaned resources
  - Verify: All error scenarios covered

BLOCKERS & ESCALATIONS:
  - If other streams slow on implementation: Create mock implementations for testing
  - If coverage below 90%: Work with engineering teams to add tests
  - If flaky tests found: Debug and fix root cause (timing issues, race conditions)
  - If performance below targets: Profile and optimize with respective teams

START DATE: Jan 27, 2026
DEADLINE: Feb 10, 2026 (Gate 1 review)
ESTIMATED EFFORT: 60 hours (1.5 weeks full-time)
```

---

## Task Coordination Rules

### Daily Standup (15 minutes)
- Report progress on assigned stream
- Identify blockers immediately
- Coordinate on integration points
- Update burndown chart

### Integration Points
| Stream | Depends On | Feeds |
|--------|-----------|-------|
| 1.1 (API) | 1.3 (DB schema) | All other streams |
| 1.2 (LXC) | 1.1 (API contract), 1.3 (port allocation) | 1.1 (container lifecycle) |
| 1.3 (DB) | None | 1.1, 1.2 |
| 1.4 (Testing) | All three | Gate 1 review |

### Code Review Process
1. Push to feature branch (one per stream)
2. Run full test suite locally
3. Open PR with test results
4. Stream lead approves
5. Merge when complete

### Escalation Path
- **Blocker**: Notify stream lead immediately
- **Design question**: Ask Sisyphus in Slack/PR
- **External blocker**: Escalate to project lead
- **Out of scope**: Document and defer to Phase 2

---

## Success Criteria for Phase 1

**MUST HAVE (Gates Required)**
- ✓ 90%+ code coverage on all new code
- ✓ All unit tests passing
- ✓ All integration tests passing
- ✓ All E2E tests passing
- ✓ Container creation <30 seconds
- ✓ SSH key isolation maintained
- ✓ Zero critical bugs found
- ✓ All API endpoints match specification

**NICE TO HAVE**
- ◇ 95%+ test pass rate
- ◇ Performance benchmarks in documentation
- ◇ Load testing with 20+ concurrent workspaces
- ◇ Database schema documentation
- ◇ Performance optimization recommendations

---

**Document**: M4 Phase 1 Delegation Task Templates  
**Created**: January 17, 2026  
**Ready for**: Immediate Delegation (Jan 27, 2026 kickoff)
