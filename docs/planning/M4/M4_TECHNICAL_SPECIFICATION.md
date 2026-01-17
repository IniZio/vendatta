# M4: Technical Specification

**Milestone**: M4 - Staging Environment & Production User Flow  
**Focus**: Architecture, APIs, Data Models  
**Audience**: Backend engineers, system architects, implementation team  
**Status**: Final Specification  

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Component Specifications](#component-specifications)
3. [API Contracts](#api-contracts)
4. [Data Models](#data-models)
5. [Configuration Formats](#configuration-formats)
6. [Communication Protocols](#communication-protocols)
7. [Security Model](#security-model)

---

## Architecture Overview

### High-Level Components

```
┌─────────────────────────────────────────────────────────┐
│                 Staging Host (Central)                   │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Coordination Server                             │   │
│  │  • HTTP REST API (port 3001)                     │   │
│  │  • Workspace management                          │   │
│  │  • User registration & auth                      │   │
│  │  • SSH forwarding & port allocation              │   │
│  │  • Service discovery                             │   │
│  │  • Metadata storage (SQLite)                     │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                        ↓ SSH/HTTP
       ┌────────────────────────────────┐
       │                                │
   ┌───▼────────────────┐      ┌───────▼───────────┐
   │  Driver Node 1     │      │  Driver Node 2    │
   │  (LXC or Docker)   │      │  (LXC or Docker)  │
   │                    │      │                   │
   │ ┌────────────────┐ │      │ ┌───────────────┐ │
   │ │  Node Agent    │ │      │ │  Node Agent   │ │
   │ │ • Cmd receiver │ │      │ │ • Cmd receiver│ │
   │ │ • Provider ops │ │      │ │ • Provider ops│ │
   │ │ • Status rept  │ │      │ │ • Status rept │ │
   │ └────────────────┘ │      │ └───────────────┘ │
   │        ↓           │      │        ↓          │
   │ ┌────────────────┐ │      │ ┌───────────────┐ │
   │ │ LXC Daemon     │ │      │ │ Docker Daemon │ │
   │ │ • Containers   │ │      │ │ • Containers  │ │
   │ │ • Networking   │ │      │ │ • Networking  │ │
   │ │ • Storage      │ │      │ │ • Storage     │ │
   │ └────────────────┘ │      │ └───────────────┘ │
   └────────────────────┘      └───────────────────┘
```

### Data Flow

**Workspace Creation**:
```
CLI → Coordination Server (HTTP POST /api/v1/workspaces/create)
    → Agent Selection (round-robin or load-based)
    → Node Agent (SSH) → Provider (LXC/Docker)
    → Container Launch
    → SSH Setup in Container
    → Service Startup
    → Status Report (back to server)
    → CLI Polling (/api/v1/workspaces/{id}/status)
    → Ready
```

**SSH Connection**:
```
User's Editor (Cursor/VS Code)
    ↓ SSH to staging-server:2222
    ↓ Coordination Server (SSH forwarding)
    ↓ LXC Container (internal:22)
    ↓ SSH Server
    ↓ Mounted Workspace
```

---

## Component Specifications

### 1. Coordination Server

**Language**: Go 1.24+  
**Framework**: net/http (standard library) or chi/mux for routing  
**Port**: 3001 (HTTP)  
**Database**: SQLite (embedded, no separate service)  
**Scaling**: Single binary, stateless (except DB)

#### Responsibilities

1. **HTTP API** - RESTful endpoints for all operations
2. **Workspace Management** - CRUD operations on workspaces
3. **User Management** - Registration, SSH key storage
4. **Port Allocation** - SSH port pool (2222-2299)
5. **SSH Forwarding** - Forward SSH connections to containers
6. **Health Monitoring** - Node agent heartbeats & status
7. **Service Discovery** - Track running services & ports

#### Key Data Structures

```go
// Workspace represents an isolated development environment
type Workspace struct {
    ID              string                 // Unique identifier (ws-abc123)
    Owner           string                 // GitHub username
    Name            string                 // Workspace name
    Status          WorkspaceStatus        // pending, creating, running, stopped
    Provider        string                 // "lxc", "docker", "qemu"
    Image           string                 // Base image name
    
    SSH struct {
        Port        int                    // Forwarded port (2222+)
        User        string                 // Username in container (dev)
        PubKey      string                 // Public key authorized in container
    }
    
    Repository struct {
        Owner       string                 // GitHub owner
        Name        string                 // Repository name
        Branch      string                 // Active branch
        URL         string                 // Git URL (git@github.com:...)
        IsFork      bool                   // Is this a user fork?
    }
    
    Services        map[string]Service     // Running services
    Node            string                 // Agent node handling this workspace
    CreatedAt       time.Time
    UpdatedAt       time.Time
    LastActivity    time.Time
}

type Service struct {
    Name           string                 // Service identifier
    Port           int                    // Port inside container
    MappedPort     int                    // Port on staging host
    Status         ServiceStatus          // running, stopped, error
    Command        string                 // Start command
    HealthCheck    HealthCheck            // Health check config
}

type User struct {
    ID             string                 // user-abc123
    GitHubUsername string                 // github handle
    GitHubID       int64                  // GitHub user ID
    SSHPublicKey   string                 // Full public key
    SSHFingerprint string                 // SHA256 fingerprint
    Workspaces     []string               // Workspace IDs
    RegisteredAt   time.Time
}
```

### 2. Node Agent

**Language**: Go 1.24+  
**Transport**: SSH (command execution) or HTTP (future)  
**Deployment**: Single binary per driver node  
**Privileges**: Runs with LXC/Docker daemon access (usually root or docker group)

#### Responsibilities

1. **Command Reception** - Receive commands from coordination server
2. **Provider Operations** - Execute container operations
3. **SSH Setup** - Configure SSH in containers
4. **Service Management** - Start/stop services
5. **Health Reporting** - Report container & service status
6. **Cleanup** - Handle workspace deletion

#### Command Interface

```go
type Command struct {
    ID           string                 // Command ID for tracking
    Type         string                 // "create", "start", "stop", "delete"
    WorkspaceID  string                 // Associated workspace
    Params       map[string]interface{} // Command-specific parameters
}

type CommandResult struct {
    ID           string                 // Command ID
    Status       string                 // "success", "failure"
    Error        string                 // Error message if failed
    Data         map[string]interface{} // Result data
    Timestamp    time.Time
}
```

#### Command Types

**create**: Launch new container
```
params:
  - container_id: unique container name
  - image: base image (ubuntu:22.04)
  - cpu: CPU limit
  - memory: memory limit
  - disk: disk size
```

**start**: Start services in container
```
params:
  - container_id: target container
  - services: list of services to start
  - depends_on: dependency graph
```

**stop**: Stop services or container
```
params:
  - container_id: target container
  - force: force termination (bool)
```

**exec**: Run command inside container
```
params:
  - container_id: target container
  - command: command to execute
  - stdin: input (optional)
```

### 3. CLI Commands

#### New Commands

```bash
# GitHub authentication
nexus auth github [--force]

# SSH key setup
nexus ssh setup [--key-path ~/.ssh/id_ed25519]

# Workspace management
nexus workspace create <repo> [--name <name>] [--provider lxc]
nexus workspace up <name>
nexus workspace down <name>
nexus workspace list [--json]
nexus workspace status <name> [--watch]
nexus workspace delete <name>
nexus workspace connect <name>
nexus workspace services <name>
nexus workspace logs <name> [--service <service>] [--follow]
nexus workspace exec <name> <command> [--stdin]

# Helper commands
nexus editor detect
nexus editor launch <name>
```

---

## API Contracts

### Base URL
```
http://coordination-server:3001/api/v1
```

### Authentication
```
Authorization: Bearer <github-token>
Content-Type: application/json
```

### Endpoints

#### 1. User Registration
```http
POST /users/register
Content-Type: application/json

Request:
{
  "github_username": "alice",
  "github_id": 123456789,
  "ssh_pubkey": "ssh-ed25519 AAAA... user@host",
  "ssh_pubkey_fingerprint": "SHA256:abcd1234..."
}

Response (201 Created):
{
  "user_id": "user-abc123",
  "github_username": "alice",
  "registered_at": "2026-01-17T10:30:00Z",
  "workspaces": []
}

Errors:
- 400: Invalid input
- 409: User already registered
```

#### 2. Create Workspace
```http
POST /workspaces/create
Content-Type: application/json

Request:
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
    {
      "name": "web",
      "command": "npm run dev",
      "port": 3000,
      "health_check": {
        "type": "http",
        "path": "/",
        "timeout": 10
      }
    },
    {
      "name": "api",
      "command": "npm run server",
      "port": 4000,
      "depends_on": ["web"]
    }
  ]
}

Response (202 Accepted):
{
  "workspace_id": "ws-abc123",
  "status": "creating",
  "ssh_port": 2222,
  "polling_url": "/api/v1/workspaces/ws-abc123/status",
  "estimated_time": "60s"
}

Errors:
- 400: Invalid request
- 409: Workspace already exists
```

#### 3. Get Workspace Status
```http
GET /workspaces/{workspace_id}/status

Response (200 OK):
{
  "workspace_id": "ws-abc123",
  "status": "running",
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
      "url": "http://localhost:23000",
      "last_check": "2026-01-17T10:31:00Z"
    },
    "api": {
      "status": "running",
      "port": 4000,
      "mapped_port": 23001,
      "url": "http://localhost:23001",
      "last_check": "2026-01-17T10:31:00Z"
    }
  },
  "repository": {
    "owner": "my-org",
    "name": "my-project",
    "branch": "main",
    "commit": "abc123def456..."
  },
  "created_at": "2026-01-17T10:30:00Z",
  "updated_at": "2026-01-17T10:31:00Z"
}

Errors:
- 404: Workspace not found
```

#### 4. Stop Workspace
```http
POST /workspaces/{workspace_id}/stop

Response (200 OK):
{
  "workspace_id": "ws-abc123",
  "status": "stopped",
  "stopped_at": "2026-01-17T10:35:00Z"
}

Errors:
- 404: Workspace not found
- 409: Workspace already stopped
```

#### 5. Delete Workspace
```http
DELETE /workspaces/{workspace_id}

Response (204 No Content):

Errors:
- 404: Workspace not found
- 409: Workspace still running (stop first)
```

#### 6. List Workspaces
```http
GET /workspaces?owner=alice&status=running

Response (200 OK):
{
  "workspaces": [
    {
      "workspace_id": "ws-abc123",
      "name": "my-project-feature",
      "status": "running",
      "owner": "alice",
      "created_at": "2026-01-17T10:30:00Z"
    }
  ],
  "total": 1
}
```

#### 7. Health Check
```http
GET /health

Response (200 OK):
{
  "status": "ok",
  "timestamp": "2026-01-17T10:35:00Z",
  "uptime": "1h30m",
  "nodes": {
    "lxc-node-1": {
      "status": "healthy",
      "last_heartbeat": "2026-01-17T10:34:50Z"
    }
  }
}
```

---

## Data Models

### Workspace Status Enum
```
pending   - Request received, awaiting agent selection
creating  - Agent creating container & setting up environment
running   - Container running, services started
stopping  - Stopping services and container
stopped   - Container stopped, can be restarted
error     - Error occurred during operation
deleting  - Awaiting final cleanup
```

### Service Status Enum
```
pending   - Service configured, awaiting startup
starting  - Service startup in progress
running   - Service running and healthy
unhealthy - Service running but health check failed
stopped   - Service stopped
error     - Error during startup
```

### Provider Enum
```
lxc    - LXC containers (preferred for staging)
docker - Docker containers
qemu   - QEMU VMs (future)
```

### HealthCheck Types
```
http   - HTTP GET to specified path
tcp    - TCP connection check
exec   - Run command inside container
custom - Custom health check script
```

---

## Configuration Formats

### .nexus/config.yaml (Repository)

```yaml
version: "1.0"
name: my-project
description: "My awesome project"

provider: lxc

services:
  web:
    command: "npm run dev"
    port: 3000
    env:
      NODE_ENV: development
    health_check:
      type: http
      path: /
      timeout: 10
    
  api:
    command: "npm run server"
    port: 4000
    depends_on: [web]
    health_check:
      type: tcp
      timeout: 5
  
  db:
    command: "postgres -D /data/postgres"
    port: 5432
    env:
      POSTGRES_DB: dev
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: "devpass"
    health_check:
      type: exec
      command: "psql -U dev -d dev -c 'SELECT 1'"
      timeout: 10

providers:
  lxc:
    image: ubuntu:22.04
    cpu: 2
    memory: "4GB"
    disk: "20GB"
  
  docker:
    image: node:20-alpine
    memory: "2GB"
  
  qemu:
    image: ubuntu:22.04
    cpu: 4
    memory: "8GB"
    disk: "50GB"

lifecycle:
  pre_start: |
    #!/bin/bash
    npm install
  post_start: |
    #!/bin/bash
    npm run setup
```

### ~/.nexus/config.yaml (User Config)

```yaml
version: "1.0"

server: staging.example.com
server_port: 3001

github:
  username: alice
  auth_token: "ghp_..."

ssh:
  key_path: "~/.ssh/id_ed25519"
  key_type: "ed25519"
  auto_upload_to_github: true

editor:
  preferred: cursor  # cursor, code, vim, nvim
  auto_launch: true

# Cache of recent workspaces
workspaces:
  my-project-feature:
    workspace_id: ws-abc123
    repo: my-org/my-project
    provider: lxc
    status: running
    ssh_port: 2222
    created_at: "2026-01-17T10:30:00Z"
```

---

## Communication Protocols

### Node Agent ↔ Coordination Server

**Transport**: SSH (command execution)

**Command Execution**:
```bash
ssh -i ~/.ssh/id_ed25519_agent node@lxc-node-1 \
  "nexus-agent cmd -j '$(cat command.json)'"
```

**Command JSON**:
```json
{
  "id": "cmd-abc123",
  "type": "create",
  "workspace_id": "ws-abc123",
  "params": {
    "image": "ubuntu:22.04",
    "cpu": 2,
    "memory": "4GB"
  }
}
```

**Response**:
```json
{
  "id": "cmd-abc123",
  "status": "success",
  "data": {
    "container_id": "workspace-abc123",
    "ip_address": "10.0.0.42"
  }
}
```

---

## Security Model

### Authentication

**Client → Coordination Server**:
- GitHub OAuth token (via `gh auth status`)
- Token stored securely in ~/.netrc or credential manager

**Coordination Server → Node Agent**:
- SSH key-based authentication
- Agent identity verified via key fingerprint

**User → Container**:
- SSH key-based authentication
- Public key seeded during container init
- User's local private key (never sent to server)

### SSH Key Management

1. **Generation**: `ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -N ""`
2. **Upload to GitHub**: `gh ssh-key add ~/.ssh/id_ed25519.pub`
3. **Registration**: Send fingerprint to coordination server
4. **Container Setup**: Public key → `/home/dev/.ssh/authorized_keys`
5. **Access**: SSH with local private key

### Network Isolation

- **SSH Forwarding**: Coordination server forwards port 2222 to container:22
- **Service Ports**: Individual mappings for each service
- **Internal Only**: Database ports not forwarded (internal only)

### Data Privacy

- No personal data beyond GitHub username stored
- SSH private keys never leave user's machine
- Public key fingerprints stored (not full keys)
- All communication via HTTPS/SSH (encrypted)

---

## Error Handling

### HTTP Status Codes

```
200 OK            - Request succeeded
201 Created       - Resource created
202 Accepted      - Request accepted, processing async
204 No Content    - Success, no response body
400 Bad Request   - Invalid input
401 Unauthorized  - Missing/invalid auth
403 Forbidden     - Authenticated but not allowed
404 Not Found     - Resource doesn't exist
409 Conflict      - Resource state conflict
500 Server Error  - Internal error
503 Unavailable   - Service temporarily down
```

### Error Response Format

```json
{
  "error": "error_code",
  "message": "Human-readable error message",
  "details": {
    "field": "optional detailed info"
  },
  "request_id": "req-xyz789"
}
```

### Common Errors

```
create_workspace:
  - workspace_already_exists
  - github_repo_not_found
  - insufficient_agent_resources
  - network_connectivity_error

ssh_connection:
  - port_allocation_failed
  - container_not_ready
  - ssh_key_mismatch

service_startup:
  - dependency_failed
  - port_already_in_use
  - health_check_timeout
```

---

## Performance Requirements

| Operation | Target | Method |
|-----------|--------|--------|
| Container creation | < 30s | Direct measurement |
| Service startup (all) | < 45s | Total time from creation |
| SSH connection | < 100ms latency | `time ssh ... 'echo hi'` |
| API response time | < 500ms | avg, 95th percentile < 2s |
| Health check | < 5s | Per service |
| Workspace status poll | < 1s | API response time |

---

## Monitoring & Observability

### Metrics to Track

- Workspace creation time (histogram)
- Service startup time per service (histogram)
- SSH connection latency (histogram)
- Node agent heartbeat success rate (%)
- Container resource utilization (CPU, memory, disk)
- Workspace uptime (%)
- Concurrent workspace count

### Logging

**Coordination Server**:
- Request/response logging (structured JSON)
- Command execution logs
- Error logs with context
- Agent heartbeat logs (summarized)

**Node Agent**:
- Command execution logs
- Container operation logs
- Service status changes
- Resource utilization logs

---

## Testing Strategy

### Unit Tests
- Command parsing & validation
- Configuration parsing
- Service dependency resolution
- Port allocation logic

### Integration Tests
- Workspace creation flow
- SSH setup in container
- Service startup sequence
- Health check execution
- Error recovery

### E2E Tests
- Complete workflow: create → start → connect → stop
- Multi-workspace scenarios
- Concurrent operations
- Failure scenarios

### Load Tests
- 10+ concurrent workspace creation
- Sustained operation (24+ hours)
- Node agent reliability

---

**Document**: M4 Technical Specification  
**Version**: 1.0  
**Status**: Final  
**Date**: January 17, 2026
