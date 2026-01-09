# Configuration Reference: Project Oursky

## 1. Overview

Oursky uses a declarative configuration system based on YAML and JSON templates. This reference covers all available configuration options and how to use them effectively.

## 2. Main Configuration (`config.yaml`)

### **Root Structure**
```yaml
name: "project-name"           # Required: Project identifier
description: "Optional description"

services: {}                   # Container services definition
agents: []                     # Enabled AI agents
mcp: {}                        # MCP server configuration
docker: {}                     # Container runtime settings
hooks: {}                      # Lifecycle scripts
sync_targets: []              # Remote targets for .vendatta sync
```

### **Services Configuration**

#### **Basic Service Definition**
```yaml
services:
  web:
    command: "cd client && npm run dev"
    healthcheck:
      url: "http://localhost:3000"
      interval: 5s
      timeout: 3s
      retries: 5
```

#### **Service Options**
| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `command` | string | Yes | Shell command to start the service |
| `healthcheck.url` | string | No | Health check endpoint URL |
| `healthcheck.interval` | duration | No | Check frequency (default: 5s) |
| `healthcheck.timeout` | duration | No | Check timeout (default: 3s) |
| `healthcheck.retries` | int | No | Max retry attempts (default: 5) |
| `depends_on` | array | No | Services that must start first |

#### **Example: Full-Stack Setup**
```yaml
services:
  db:
    command: "docker-compose up -d postgres"
    healthcheck:
      url: "http://localhost:5432/health"
      interval: 10s
      retries: 10

  api:
    command: "cd server && npm run dev"
    healthcheck:
      url: "http://localhost:5000/api/health"
    depends_on: ["db"]

  web:
    command: "cd client && npm run dev"
    healthcheck:
      url: "http://localhost:3000"
    depends_on: ["api"]
```

### **Agents Configuration**

#### **Agent Definition**
```yaml
agents:
  - name: "cursor"              # Agent identifier
    enabled: true               # Enable/disable this agent
  - name: "opencode"
    enabled: true
  - name: "claude-desktop"
    enabled: false
```

#### **Supported Agents**
| Agent | Description | Generated Config |
|-------|-------------|------------------|
| `cursor` | VS Code extension with MCP | `.cursor/mcp.json` |
| `opencode` | Standalone AI assistant | `opencode.json` + `.opencode/` |
| `claude-desktop` | Anthropic desktop app | `claude_desktop_config.json` |
| `claude-code` | Anthropic CLI tool | `claude_code_config.json` |

### **MCP Configuration**

#### **Basic MCP Setup**
```yaml
mcp:
  enabled: true                 # Enable MCP server
  port: 3001                    # Server port
  host: "localhost"             # Server host
```

#### **MCP Options**
| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | boolean | `true` | Toggle MCP functionality |
| `port` | int | `3001` | TCP port for MCP server |
| `host` | string | `"localhost"` | Host address to bind to |

### **Docker Configuration**

#### **Container Runtime**
```yaml
docker:
  image: "ubuntu:22.04"         # Base container image
  dind: true                    # Enable Docker-in-Docker
  privileged: false             # Run in privileged mode
  memory: "2g"                  # Memory limit
  cpu: "1.0"                    # CPU limit (cores)
```

### **Hooks Configuration**

#### **Lifecycle Scripts**
```yaml
hooks:
  setup: ".vendatta/hooks/setup.sh"  # Run after container creation
  dev: ".vendatta/hooks/dev.sh"      # Run before dev session starts
  teardown: ".vendatta/hooks/teardown.sh"  # Run before cleanup
```

### **Sync Targets Configuration**

#### **Remote Sync Targets**
```yaml
sync_targets:
  - name: "upstream"              # Remote name
    url: "https://github.com/example/upstream.git"  # Repository URL
  - name: "configs"
    url: "https://github.com/example/configs.git"
```

#### **Remote Options**
| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `name` | string | Yes | Git remote name |
| `url` | string | Yes | Repository URL |

#### **Usage**
Configured sync targets push only the `.vendatta` directory to the specified remote. Sync individually or all at once:
```bash
# Sync specific target
vendatta remote sync <target-name>

# Sync all configured targets
vendatta remote sync-all
```

## 3. Template System

### **Template Variables**
All `.tpl` files support variable substitution using Go template syntax:

```yaml
# In any .tpl file
mcp:
  server: "http://{{.Host}}:{{.Port}}"
  token: "{{.AuthToken}}"
```

#### **Available Variables**
| Variable | Description | Example |
|----------|-------------|---------|
| `{{.Host}}` | MCP server host | `localhost` |
| `{{.Port}}` | MCP server port | `3001` |
| `{{.AuthToken}}` | Authentication token | `abc123...` |
| `{{.ProjectName}}` | Project name | `my-project` |
| `{{.DatabaseURL}}` | Database connection | `postgresql://...` |

### **Skills Templates** (`templates/skills/`)

Following [agentskills.io](https://agentskills.io) specification:

```yaml
name: "web-search"
description: "Search the web for information"
version: "1.0.0"
author: "Your Team"

parameters:
  type: object
  properties:
    query:
      type: string
      description: "The search query"
    limit:
      type: integer
      default: 10
  required: ["query"]

execute:
  type: "http"
  url: "https://api.searchengine.com/search"
  method: "GET"

permissions:
  - "web:read"
```

### **Commands Templates** (`templates/commands/`)

Standardized command definitions:

```yaml
name: "build"
description: "Build the project"
aliases: ["compile", "make"]

steps:
  - name: "Install dependencies"
    command: "npm install"
  - name: "Lint code"
    command: "npm run lint"
  - name: "Run tests"
    command: "npm test"
  - name: "Build artifacts"
    command: "npm run build"

env:
  NODE_ENV: "production"
```

### **Rules Templates** (`templates/rules/`)

Following [agents.md](https://github.com/agentsmd/agents.md) format:

```markdown
---
title: "Code Quality Standards"
version: "1.0.0"
applies_to: ["**/*.js", "**/*.ts", "**/*.py"]
priority: "high"
---

# Code Quality Standards

## Naming Conventions
- Use camelCase for variables and functions
- Use PascalCase for classes and components
- Use UPPER_CASE for constants

## Code Structure
- Keep functions under 50 lines
- Use early returns to reduce nesting
- Group related functionality together

## Documentation
- Add JSDoc comments for public APIs
- Document complex business logic
- Keep README files current
```

## 4. Agent-Specific Configuration

### **Cursor Configuration**
Generated: `.cursor/mcp.json`

```json
{
  "mcpServers": {
    "project-name": {
      "type": "http",
      "url": "http://localhost:3001",
      "headers": {
        "Authorization": "Bearer YOUR_TOKEN"
      }
    }
  }
}
```

### **OpenCode Configuration**
Generated: `opencode.json` + `.opencode/` directory

```json
{
  "$schema": "https://opencode.ai/config.json",
  "model": "anthropic/claude-sonnet-4-5",
  "mcp": {
    "project-name": {
      "type": "remote",
      "url": "http://localhost:3001",
      "enabled": true
    }
  },
  "rules": ["code-quality", "collaboration"],
  "skills": ["web-search", "file-operations"],
  "commands": ["build", "deploy"]
}
```

### **Claude Desktop/Code Configuration**
Generated: `claude_desktop_config.json` / `claude_code_config.json`

```json
{
  "mcpServers": {
    "project-name": {
      "command": "npx",
      "args": ["-y", "mcp-remote", "http://localhost:3001"],
      "env": {
        "MCP_AUTH_TOKEN": "YOUR_TOKEN"
      }
    }
  }
}
```

## 5. Environment Variables

Oursky supports environment variable substitution in configuration:

### **In config.yaml**
```yaml
mcp:
  port: "{{.Env.MCP_PORT}}"
  host: "{{.Env.MCP_HOST}}"
```

### **In Templates**
```yaml
# In any .tpl file
database_url: "{{.Env.DATABASE_URL}}"
api_key: "{{.Env.OPENAI_API_KEY}}"
```

### **Common Environment Variables**
```bash
# MCP Configuration
MCP_PORT=3001
MCP_HOST=localhost
MCP_AUTH_TOKEN=your-secret-token

# Project Configuration
PROJECT_NAME=my-awesome-project
DATABASE_URL=postgresql://user:pass@localhost:5432/db

# API Keys
OPENAI_API_KEY=sk-...
GITHUB_TOKEN=ghp_...
```

## 6. Best Practices

### **Configuration Organization**
- Keep `config.yaml` focused on services and agents
- Use templates for reusable capabilities
- Document custom configurations with comments

### **Security**
- Never commit sensitive data to version control
- Use environment variables for secrets
- Regularly rotate authentication tokens

### **Performance**
- Use appropriate health check intervals
- Configure resource limits for containers
- Enable D in D only when needed

### **Maintenance**
- Version control your `.vendatta/` directory
- Test configuration changes in isolated branches
- Document custom templates and their purpose

## 7. Troubleshooting

### **Common Issues**

#### **MCP Connection Failed**
```yaml
# Check mcp configuration
mcp:
  enabled: true
  port: 3001
  host: "localhost"
```

#### **Agent Config Not Generated**
```yaml
# Ensure agent is enabled
agents:
  - name: "cursor"
    enabled: true
```

#### **Container Won't Start**
```yaml
# Check service dependencies
services:
  api:
    depends_on: ["db"]  # Wait for database first
```

#### **Template Variables Not Substituted**
- Ensure variables are defined in the correct context
- Check for typos in variable names
- Verify environment variables are set

### **Debugging Commands**
```bash
# Check MCP server status
curl http://localhost:3001/health

# View generated configurations
cat .cursor/mcp.json
cat opencode.json

# Check container logs
docker logs vendatta-session-123
```

## 8. Migration Guide

### **Upgrading from Manual Config**
1. Run `vendatta init` to generate new structure
2. Move existing configs to appropriate template directories
3. Update `config.yaml` with your settings
4. Test with `vendatta dev test-branch`

### **From Other Tools**
- **docker-compose**: Move service definitions to `services:` section
- **Manual scripts**: Convert to templates with `{{.Variable}}` syntax
- **Environment files**: Use `{{.Env.VAR_NAME}}` substitution