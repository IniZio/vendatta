# M4 API Reference Specification

**API Version**: v1  
**Base URL**: `http://coordination-server:3001/api/v1`  
**Content-Type**: `application/json`  
**Authentication**: GitHub OAuth Token (via Authorization header)

---

## Quick Reference

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/users/register` | POST | Register user with SSH key |
| `/workspaces/create` | POST | Create new workspace |
| `/workspaces/{id}/status` | GET | Get workspace status |
| `/workspaces/{id}/stop` | POST | Stop workspace |
| `/workspaces/{id}` | DELETE | Delete workspace |
| `/workspaces` | GET | List workspaces |
| `/health` | GET | Server health check |

---

## Authentication

All requests require GitHub token authentication:

```
Authorization: Bearer ghp_xxxxxxxxxxxx
```

Get token from:
```bash
gh auth token
```

---

## Error Handling

All error responses follow this format:

```json
{
  "error": "error_code",
  "message": "Human-readable description",
  "details": {
    "field": "specific_information"
  },
  "request_id": "req-abc123def456"
}
```

Common HTTP Status Codes:
- `200` - OK
- `201` - Created
- `202` - Accepted (async operation)
- `204` - No Content
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

---

## Endpoints

### 1. Register User

**Endpoint**
```
POST /users/register
```

**Purpose**: Register user with their SSH public key with coordination server

**Request**
```json
{
  "github_username": "alice",
  "github_id": 123456789,
  "ssh_pubkey": "ssh-ed25519 AAAA... user@example.com",
  "ssh_pubkey_fingerprint": "SHA256:abcd1234..."
}
```

**Response (201 Created)**
```json
{
  "user_id": "user-abc123",
  "github_username": "alice",
  "ssh_pubkey_fingerprint": "SHA256:abcd1234...",
  "registered_at": "2026-01-17T10:30:00Z",
  "workspaces": []
}
```

**Errors**
- `400` - Invalid SSH public key format
- `409` - User already registered

**Example**
```bash
curl -X POST http://localhost:3001/api/v1/users/register \
  -H "Authorization: Bearer $(gh auth token)" \
  -H "Content-Type: application/json" \
  -d '{
    "github_username": "alice",
    "github_id": 123456789,
    "ssh_pubkey": "ssh-ed25519 AAAA... alice@example.com",
    "ssh_pubkey_fingerprint": "SHA256:abcd1234..."
  }'
```

---

### 2. Create Workspace

**Endpoint**
```
POST /workspaces/create
```

**Purpose**: Create new workspace and request container creation

**Request**
```json
{
  "github_username": "alice",
  "workspace_name": "my-project-feature",
  "repo": {
    "owner": "my-org",
    "name": "my-project",
    "url": "git@github.com:my-org/my-project.git",
    "branch": "main",
    "is_fork": false
  },
  "provider": "lxc",
  "image": "ubuntu:22.04",
  "services": [
    {
      "name": "web",
      "command": "npm run dev",
      "port": 3000,
      "depends_on": [],
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
      "depends_on": ["web"],
      "health_check": {
        "type": "tcp",
        "timeout": 5
      }
    }
  ]
}
```

**Response (202 Accepted)**
```json
{
  "workspace_id": "ws-abc123",
  "status": "creating",
  "ssh_port": 2222,
  "polling_url": "/api/v1/workspaces/ws-abc123/status",
  "estimated_time_seconds": 60,
  "created_at": "2026-01-17T10:30:00Z"
}
```

**Status Values**
- `pending` - Request received, awaiting assignment
- `creating` - Container & services being set up
- `running` - Ready for use
- `error` - Error during creation
- `stopped` - Stopped (can be restarted)

**Errors**
- `400` - Invalid request (missing required fields)
- `404` - Repository not found on GitHub
- `409` - Workspace name already exists for user
- `503` - No available nodes

**Example**
```bash
curl -X POST http://localhost:3001/api/v1/workspaces/create \
  -H "Authorization: Bearer $(gh auth token)" \
  -H "Content-Type: application/json" \
  -d @workspace-request.json

# Then poll status
curl http://localhost:3001/api/v1/workspaces/ws-abc123/status \
  -H "Authorization: Bearer $(gh auth token)"
```

---

### 3. Get Workspace Status

**Endpoint**
```
GET /workspaces/{workspace_id}/status
```

**Purpose**: Get current status of workspace, services, and connection info

**Parameters**
- `workspace_id` (path) - Workspace ID from creation response

**Response (200 OK)**
```json
{
  "workspace_id": "ws-abc123",
  "owner": "alice",
  "name": "my-project-feature",
  "status": "running",
  "provider": "lxc",
  
  "ssh": {
    "host": "staging.example.com",
    "port": 2222,
    "user": "dev",
    "key_required": "~/.ssh/id_ed25519"
  },
  
  "services": {
    "web": {
      "name": "web",
      "status": "running",
      "port": 3000,
      "mapped_port": 23000,
      "health": "healthy",
      "url": "http://localhost:23000",
      "last_check": "2026-01-17T10:31:30Z"
    },
    "api": {
      "name": "api",
      "status": "running",
      "port": 4000,
      "mapped_port": 23001,
      "health": "healthy",
      "url": "http://localhost:23001",
      "last_check": "2026-01-17T10:31:30Z"
    }
  },
  
  "repository": {
    "owner": "my-org",
    "name": "my-project",
    "branch": "main",
    "commit": "abc123def456...",
    "url": "git@github.com:my-org/my-project.git"
  },
  
  "node": "lxc-node-1",
  "created_at": "2026-01-17T10:30:00Z",
  "updated_at": "2026-01-17T10:31:30Z"
}
```

**Service Status Values**
- `pending` - Configured, awaiting startup
- `starting` - Startup in progress
- `running` - Running and healthy
- `unhealthy` - Running but health check failed
- `stopped` - Stopped
- `error` - Error during startup

**Service Health Values**
- `healthy` - Last health check passed
- `unhealthy` - Last health check failed
- `unknown` - No health check performed yet
- `timeout` - Health check timed out

**Errors**
- `404` - Workspace not found

**Example**
```bash
# Poll until status is "running"
for i in {1..60}; do
  curl http://localhost:3001/api/v1/workspaces/ws-abc123/status \
    -H "Authorization: Bearer $(gh auth token)" | jq '.status'
  
  if [ "$(curl -s http://localhost:3001/api/v1/workspaces/ws-abc123/status \
    -H 'Authorization: Bearer '$TOKEN'' | jq -r '.status')" = "running" ]; then
    echo "Workspace ready!"
    break
  fi
  
  sleep 1
done
```

---

### 4. List Workspaces

**Endpoint**
```
GET /workspaces?owner=alice&status=running
```

**Purpose**: List all workspaces for authenticated user

**Query Parameters**
- `owner` (optional) - Filter by owner (GitHub username)
- `status` (optional) - Filter by status (running, stopped, etc.)
- `limit` (optional) - Max results (default 50)
- `offset` (optional) - Pagination offset

**Response (200 OK)**
```json
{
  "workspaces": [
    {
      "workspace_id": "ws-abc123",
      "name": "my-project-feature",
      "owner": "alice",
      "status": "running",
      "provider": "lxc",
      "ssh_port": 2222,
      "created_at": "2026-01-17T10:30:00Z",
      "services_count": 2
    }
  ],
  "total": 1,
  "limit": 50,
  "offset": 0
}
```

**Example**
```bash
# List all running workspaces
curl "http://localhost:3001/api/v1/workspaces?status=running" \
  -H "Authorization: Bearer $(gh auth token)"

# List and show in JSON
curl "http://localhost:3001/api/v1/workspaces?owner=alice" \
  -H "Authorization: Bearer $(gh auth token)" | jq '.workspaces[]'
```

---

### 5. Stop Workspace

**Endpoint**
```
POST /workspaces/{workspace_id}/stop
```

**Purpose**: Stop a running workspace (can be restarted)

**Parameters**
- `workspace_id` (path) - Workspace ID

**Request Body**
```json
{
  "force": false
}
```

**Response (200 OK)**
```json
{
  "workspace_id": "ws-abc123",
  "status": "stopped",
  "stopped_at": "2026-01-17T10:35:00Z"
}
```

**Query Parameters**
- `force` (optional, boolean) - Force stop without graceful shutdown

**Errors**
- `404` - Workspace not found
- `409` - Workspace already stopped

**Example**
```bash
curl -X POST http://localhost:3001/api/v1/workspaces/ws-abc123/stop \
  -H "Authorization: Bearer $(gh auth token)" \
  -H "Content-Type: application/json"
```

---

### 6. Delete Workspace

**Endpoint**
```
DELETE /workspaces/{workspace_id}
```

**Purpose**: Permanently delete a workspace (must be stopped first)

**Parameters**
- `workspace_id` (path) - Workspace ID

**Response (204 No Content)**

**Errors**
- `404` - Workspace not found
- `409` - Workspace still running (stop first)

**Example**
```bash
# Stop workspace first
curl -X POST http://localhost:3001/api/v1/workspaces/ws-abc123/stop \
  -H "Authorization: Bearer $(gh auth token)"

# Then delete
curl -X DELETE http://localhost:3001/api/v1/workspaces/ws-abc123 \
  -H "Authorization: Bearer $(gh auth token)"
```

---

### 7. Health Check

**Endpoint**
```
GET /health
```

**Purpose**: Check server health and connected nodes

**Response (200 OK)**
```json
{
  "status": "ok",
  "timestamp": "2026-01-17T10:35:00Z",
  "uptime_seconds": 5400,
  "nodes": {
    "lxc-node-1": {
      "status": "healthy",
      "last_heartbeat": "2026-01-17T10:34:50Z",
      "workspaces": 5
    },
    "docker-node-1": {
      "status": "healthy",
      "last_heartbeat": "2026-01-17T10:34:55Z",
      "workspaces": 3
    }
  },
  "version": "1.0.0"
}
```

**No Authentication Required**

**Example**
```bash
curl http://localhost:3001/api/v1/health | jq
```

---

## Response Envelopes

### Success Response
```json
{
  "data": { /* actual response data */ },
  "status": "success"
}
```

### Error Response
```json
{
  "error": "error_code",
  "message": "Human-readable message",
  "details": { /* additional context */ },
  "request_id": "req-abc123def456"
}
```

---

## Rate Limiting

Current implementation: **No rate limiting**

Future implementation:
- 100 requests per minute per user
- Returned in headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`

---

## Webhooks (Future)

Planned for post-M4:
- Workspace status changes
- Service health changes
- Node availability changes

---

## Testing the API

### Using curl

```bash
# Set token
TOKEN=$(gh auth token)

# Create workspace
curl -X POST http://localhost:3001/api/v1/workspaces/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "github_username": "test-user",
    "workspace_name": "test-ws",
    "repo": {"owner": "test", "name": "repo", "branch": "main"},
    "provider": "lxc",
    "image": "ubuntu:22.04",
    "services": []
  }' | jq

# Get status
WS_ID=$(curl -s http://localhost:3001/api/v1/workspaces \
  -H "Authorization: Bearer $TOKEN" | jq -r '.workspaces[0].workspace_id')

curl http://localhost:3001/api/v1/workspaces/$WS_ID/status \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Using httpie

```bash
# Much cleaner output
http POST localhost:3001/api/v1/users/register \
  Authorization:"Bearer $TOKEN" \
  github_username=alice \
  github_id:=123456789 \
  ssh_pubkey="ssh-ed25519 AAAA..." \
  ssh_pubkey_fingerprint="SHA256:..."
```

### Using Python

```python
import requests
import os

token = os.getenv('GH_TOKEN')
headers = {'Authorization': f'Bearer {token}'}

# Create workspace
response = requests.post(
    'http://localhost:3001/api/v1/workspaces/create',
    headers=headers,
    json={
        'github_username': 'alice',
        'workspace_name': 'test',
        'repo': {...},
        'provider': 'lxc',
        'image': 'ubuntu:22.04',
        'services': []
    }
)
print(response.json())
```

---

**Document**: M4 API Reference Specification  
**Version**: 1.0  
**Status**: Final  
**Date**: January 17, 2026
