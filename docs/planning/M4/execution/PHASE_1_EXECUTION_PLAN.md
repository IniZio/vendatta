# M4 Phase 1: Coordination Server Foundation - Detailed Execution Plan

**Duration**: 2 weeks (Jan 27 - Feb 10, 2026)  
**Objective**: Extend coordination server for user registration, workspace creation, service discovery  
**Success Criteria**: Server stable 99.9%, 10+ concurrent workspaces, all tests >90% coverage

---

## Overview

Phase 1 builds the core coordination server that manages:
- User registration with GitHub profiles & SSH keys
- Workspace lifecycle (create, list, delete, get)
- Service discovery and port allocation
- SSH port forwarding management
- Metadata persistence (SQLite)

This is the **critical path blocker** for M4. Once complete, M4 Phase 2-4 can run in parallel.

---

## 4 Parallel Work Streams

### Stream 1.1: API Endpoints (Backend Engineer Lead)
**Duration**: 2 weeks  
**Effort**: 80 hours  
**Deliverable**: REST API endpoints working with real data

#### Endpoints to Implement

**Users**
- `POST /api/v1/users/register-github` - Register GitHub user with SSH key
- `GET /api/v1/users/:username` - Get user details
- `PUT /api/v1/users/:username/ssh-key` - Update SSH key

**Workspaces**
- `POST /api/v1/workspaces` - Create new workspace
- `GET /api/v1/workspaces` - List user's workspaces
- `GET /api/v1/workspaces/:id` - Get workspace details
- `DELETE /api/v1/workspaces/:id` - Delete workspace
- `PATCH /api/v1/workspaces/:id/status` - Update workspace status

**Services**
- `GET /api/v1/workspaces/:id/services` - List running services
- `POST /api/v1/workspaces/:id/services/:name/port-forward` - Request port forward

**Health**
- `GET /api/v1/health` - Server health check
- `GET /api/v1/status` - Detailed server status

#### Success Criteria
- All endpoints return JSON matching M4_API_SPECIFICATION.md exactly
- Proper HTTP status codes (201 create, 400 validation, 409 conflict, etc.)
- Request validation on all endpoints
- Comprehensive error messages
- Integration tests for all endpoints

#### References
- M4_API_SPECIFICATION.md - Exact API contract
- M4_TECHNICAL_SPECIFICATION.md - Architecture & design
- pkg/coordination/handlers.go - Where to implement

---

### Stream 1.2: LXC Integration (Systems Engineer Lead)
**Duration**: 2 weeks  
**Effort**: 60 hours  
**Deliverable**: Containers can be created/started/stopped via coordination server

#### Tasks

1. **Container Lifecycle**
   - Create LXC container from image (ubuntu:22.04)
   - Start/stop container operations
   - Delete container cleanup
   - Handle resource constraints (CPU, memory, disk)

2. **SSH Configuration**
   - Generate container SSH keys
   - Configure SSH server in container
   - Set user public keys in authorized_keys
   - Enable remote SSH access

3. **Port Forwarding**
   - Allocate SSH port from range (2222-2299)
   - Forward port to container SSH (port 22)
   - Manage port allocation table
   - Cleanup on container deletion

4. **Service Startup**
   - Execute setup scripts in container
   - Parse services from config
   - Start services with dependency ordering
   - Implement health checks

#### Success Criteria
- Container creates in <30 seconds
- SSH port allocated correctly
- Can SSH into container immediately
- Port forwarding is stable
- Services start in dependency order
- Health checks pass before marking ready
- Clean up on container deletion (no orphans)

#### References
- M4_TECHNICAL_SPECIFICATION.md - Container design
- pkg/provider/lxc/ - LXC provider implementation
- pkg/transport/ - SSH transport layer

---

### Stream 1.3: Database & State Management (Backend Engineer Lead)
**Duration**: 2 weeks  
**Effort**: 40 hours  
**Deliverable**: Persistent SQLite database with proper schema

#### Database Design

**Schema**
```sql
-- Users
CREATE TABLE users (
    id TEXT PRIMARY KEY,                  -- github_username
    github_id INTEGER UNIQUE,
    email TEXT,
    public_key_fingerprint TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Workspaces
CREATE TABLE workspaces (
    id TEXT PRIMARY KEY,                  -- uuid or hash
    user_id TEXT FOREIGN KEY,
    name TEXT,
    status TEXT,                          -- creating, running, stopped, deleted
    provider TEXT,                        -- docker, lxc, qemu
    repository TEXT,                      -- for git clone
    ssh_port INTEGER UNIQUE,              -- allocated port
    container_id TEXT,                    -- provider-specific ID
    created_at TIMESTAMP,
    started_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Services
CREATE TABLE services (
    id TEXT PRIMARY KEY,
    workspace_id TEXT FOREIGN KEY,
    name TEXT,
    port INTEGER,                         -- internal port
    forwarded_port INTEGER,               -- external port
    status TEXT,                          -- running, stopped, failed
    last_health_check TIMESTAMP,
    created_at TIMESTAMP
);

-- Port Allocations
CREATE TABLE port_allocations (
    port INTEGER PRIMARY KEY,             -- 2222-2299
    workspace_id TEXT,
    allocated_at TIMESTAMP,
    released_at TIMESTAMP
);
```

#### Data Access Objects (DAOs)
- UserDAO - CRUD for users, SSH key management
- WorkspaceDAO - CRUD for workspaces, status transitions
- ServiceDAO - CRUD for services, health check tracking
- PortAllocationDAO - Allocate/release ports atomically

#### Tasks
1. Create database schema migrations
2. Implement DAOs with proper error handling
3. Add transaction support for atomic operations
4. Implement indexes on frequently queried fields
5. Add database connection pooling

#### Success Criteria
- Database initializes automatically
- All DAOs tested with >90% coverage
- Transactions prevent race conditions
- Database properly cleans up deleted records
- Schema versioning for migrations

#### References
- M4_TECHNICAL_SPECIFICATION.md - Data model
- Database schema: See M4_TECHNICAL_SPECIFICATION.md appendix
- pkg/coordination/ - Where to place DAO code

---

### Stream 1.4: Testing & Quality (QA Specialist Lead)
**Duration**: 2 weeks  
**Effort**: 60 hours  
**Deliverable**: Full test coverage (unit, integration, E2E)

#### Test Categories

**Unit Tests** (40 hours)
- API endpoint handlers
  - Valid/invalid requests
  - Request validation
  - Error responses
  - Status codes
- DAO operations
  - CRUD operations
  - Transaction handling
  - Error conditions
- Business logic
  - Port allocation
  - Status transitions
  - Dependency ordering

**Integration Tests** (15 hours)
- Database with DAOs
- API → Database → DAO flow
- Transaction rollback scenarios
- Concurrent operations

**E2E Tests** (15 hours)
- Full workspace lifecycle (create → start → stop → delete)
- Port forwarding verification
- Service startup and health checks
- Concurrent workspace operations
- Container cleanup verification

#### Coverage Requirements
- Minimum 90% coverage on new code
- All error paths tested
- All status transitions tested
- All API endpoints tested (success & failure)

#### Success Criteria
- All unit tests pass
- All integration tests pass
- All E2E tests pass
- 90%+ coverage on all new code
- No flaky tests
- Performance benchmarks meet targets:
  - Container creation: <30s
  - Service startup: <60s
  - API endpoint response: <200ms

#### Testing Tools
- testify/assert and testify/require
- Table-driven tests
- Mock providers for isolation
- Load testing (10+ concurrent workspaces)

---

## Verification Gate 1 (Feb 10)

**Sisyphus Review Duration**: 1-2 hours

### Verification Checklist

- [ ] Code review
  - [ ] No anti-patterns (interface{}, @ts-ignore equivalents)
  - [ ] Proper error wrapping with context
  - [ ] Database transactions for consistency
  - [ ] SSH key isolation (no ~/.ssh usage)
  - [ ] Resource cleanup on error paths

- [ ] Test coverage
  - [ ] 90%+ coverage on all new code
  - [ ] All endpoints tested (success & error)
  - [ ] Database operations tested
  - [ ] Port allocation tested
  - [ ] Service startup tested

- [ ] E2E validation
  - [ ] Full workspace lifecycle works
  - [ ] 10+ concurrent workspaces without interference
  - [ ] Services start in dependency order
  - [ ] Health checks pass
  - [ ] Container cleanup works

- [ ] Performance
  - [ ] Container creation <30s
  - [ ] SSH port allocation <100ms
  - [ ] API response time <200ms
  - [ ] No memory leaks

- [ ] Security
  - [ ] SSH keys properly handled
  - [ ] No credentials in logs
  - [ ] Input validation on all endpoints
  - [ ] Database transactions prevent race conditions
  - [ ] Port allocation prevents conflicts

---

## Deliverables by Stream

### Stream 1.1 Deliverables
```
pkg/coordination/
├── handlers.go (updated)              # All API endpoints
├── handlers_test.go (updated)         # Endpoint tests
├── router.go (new)                    # HTTP routing
├── router_test.go (new)               # Router tests
└── types.go (updated)                 # Request/response types
```

### Stream 1.2 Deliverables
```
pkg/coordination/
├── container.go (new)                 # Container management
├── container_test.go (new)            # Container tests
├── ssh.go (new)                       # SSH key management
├── ssh_test.go (new)                  # SSH tests
└── port_allocator.go (new)            # Port management
```

### Stream 1.3 Deliverables
```
pkg/coordination/
├── db/ (new directory)
│   ├── schema.sql (new)               # Database schema
│   ├── migrations.go (new)            # Schema migrations
│   ├── user_dao.go (new)              # User data access
│   ├── user_dao_test.go (new)         # User tests
│   ├── workspace_dao.go (new)         # Workspace data access
│   ├── workspace_dao_test.go (new)    # Workspace tests
│   ├── service_dao.go (new)           # Service data access
│   ├── service_dao_test.go (new)      # Service tests
│   └── types.go (new)                 # Database models
└── coordination.go (updated)          # Inject database
```

### Stream 1.4 Deliverables
```
docs/planning/M4/
├── test-results.md                    # Test run results
├── coverage-report.html               # Coverage visualization
└── performance-benchmarks.md          # Performance metrics
```

---

## Dependencies & Coordination

**Critical Path**:
1. Stream 1.1 & 1.3 start in parallel (API routes first, DB schema ready)
2. Stream 1.2 integrates with 1.1 for lifecycle operations
3. Stream 1.4 integrates with all three for comprehensive testing

**Daily Sync**:
- 15-minute standup on progress
- Identify blockers immediately
- Coordinate on schema/API changes

**Integration Points**:
- All endpoints must persist to database
- Database writes trigger container lifecycle
- Container operations must report status to API
- Tests must exercise full workflow

---

## Effort Breakdown

| Stream | Role | Hours | FTE Weeks |
|--------|------|-------|-----------|
| 1.1 | Backend | 80 | 2.0 |
| 1.2 | Systems | 60 | 1.5 |
| 1.3 | Backend | 40 | 1.0 |
| 1.4 | QA | 60 | 1.5 |
| **TOTAL** | **4 engineers** | **240** | **6.0 FTE-weeks** |

**Timeline**: 2 calendar weeks with 4 people working in parallel = **1 full team sprint**

---

## Success Metrics

### Must-Have (Gates Cannot Pass Without)
- ✓ 90%+ test coverage on all new code
- ✓ All API endpoints match specification
- ✓ Container creation <30 seconds
- ✓ 10+ concurrent workspaces work
- ✓ SSH keys isolated (not in ~/.ssh)
- ✓ Zero critical bugs found in testing

### Nice-to-Have
- ◇ E2E test pass rate 95%+
- ◇ Performance benchmarks publish-ready
- ◇ Load testing with 20+ concurrent workspaces
- ◇ Database schema documentation

---

## References

**API Specification**:
- `/docs/planning/M4/api/M4_API_SPECIFICATION.md` - Exact contracts

**Technical Specification**:
- `/docs/planning/M4/M4_TECHNICAL_SPECIFICATION.md` - Architecture & design

**Code References**:
- `pkg/coordination/` - Main coordination server package
- `pkg/provider/lxc/` - LXC provider for reference
- `pkg/transport/` - SSH transport layer
- `e2e/` - E2E testing infrastructure

---

## Next Phase

Once Gate 1 passes, proceed immediately to:

**Phase 2**: GitHub CLI Integration (Feb 3 - Feb 17)
- `nexus auth github` command
- `nexus workspace create <repo>` automation
- `nexus workspace connect` with editor launch
- `nexus workspace services` discovery

---

**Document**: M4 Phase 1 Execution Plan  
**Status**: Ready for Execution (Jan 27)  
**Created**: January 17, 2026
