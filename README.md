# Oursky (Vibegear Rework)

Oursky is a developer-centric, single-binary dev environment manager. It abstracts complex infrastructure into a simple CLI, providing isolated, reproducible, and AI-agent-friendly codespaces using local providers (Docker, LXC).

## Key Features

- **Single Binary**: No complex dependencies on the host machine.
- **Git Worktree Managed**: Automatic isolation of code for every branch/session.
- **BYOA (Bring Your Own Agent)**: Built-in Model Context Protocol (MCP) server for Cursor, OpenCode, and Claude.
- **Service Discovery**: Automatic port mapping and environment variable injection (`OURSKY_SERVICE_WEB_URL`, etc.).
- **Docker-in-Docker**: Seamless support for `docker-compose` projects inside isolated environments.

## Quick Start

### 1. Installation
Build the binary using Go (requires Go 1.24+):
```bash
go build -o oursky cmd/oursky/main.go
```

### 2. Onboard a Project
Initialize the `.oursky` configuration in your repository:
```bash
./oursky init
```

### 3. Configure Your Stack
Edit `.oursky/config.yaml` to define your development environment:

```yaml
# Project settings
name: "my-web-app"

# Define your services (database, API, frontend, etc.)
services:
  db:
    command: "docker-compose up -d postgres"
    healthcheck:
      url: "http://localhost:5432/health"
  api:
    command: "cd server && npm run dev"
    depends_on: ["db"]
  web:
    command: "cd client && npm run dev"
    depends_on: ["api"]

# Enable AI agents
agents:
  - name: "cursor"
    enabled: true
  - name: "opencode"
    enabled: true

# MCP server configuration
mcp:
  enabled: true
  port: 3001
```

### 4. Start Developing
Spin up an isolated environment for a feature branch:
```bash
./oursky dev feature-login
```

The CLI automatically generates AI agent configurations and starts your services in an isolated environment.

## AI Agent Configuration

Oursky automatically configures your favorite AI coding assistants to work securely with your isolated development environments.

### Supported Agents

| Agent | Description | Generated Config |
|-------|-------------|------------------|
| **Cursor** | VS Code extension with AI | `.cursor/mcp.json` |
| **OpenCode** | Standalone AI assistant | `opencode.json` + `.opencode/` |
| **Claude Desktop** | Anthropic's desktop app | `claude_desktop_config.json` |
| **Claude Code** | Anthropic's CLI tool | `claude_code_config.json` |

### How It Works

1. **Configure agents** in `.oursky/config.yaml`:
   ```yaml
   agents:
     - name: "cursor"
       enabled: true
     - name: "opencode"
       enabled: true
   ```

2. **Run development session**:
   ```bash
   ./oursky dev my-feature
   ```

3. **Open in your AI agent**:
   - Cursor: Open `.oursky/worktrees/my-feature/`
   - OpenCode: The generated `opencode.json` connects automatically
   - Claude: Uses generated config files

### Shared Capabilities

Oursky includes standard AI capabilities that work across all agents:

- **Skills**: Web search, file operations, data analysis
- **Commands**: Build, deploy, git operations
- **Rules**: Code quality standards, collaboration guidelines

Customize these in `.oursky/templates/` and they'll be available to all your enabled agents.

## Configuration System

### File Structure
```
.oursky/
├── config.yaml          # Main project configuration
├── templates/           # Shared AI capabilities
│   ├── skills/          # Reusable AI skills
│   ├── commands/        # Development commands
│   └── rules/           # Coding guidelines
├── agents/              # Agent-specific templates
│   ├── cursor/          # Cursor configuration
│   ├── opencode/        # OpenCode configuration
│   ├── claude-desktop/  # Claude Desktop config
│   └── claude-code/     # Claude Code config
└── worktrees/           # Auto-generated environments
```

### Customizing Templates

**Add a new skill** in `.oursky/templates/skills/my-skill.yaml`:
```yaml
name: "my-custom-skill"
description: "Does something useful"
parameters:
  type: object
  properties:
    input: { type: "string" }
execute:
  command: "node"
  args: ["scripts/my-skill.js"]
```

**Add coding rules** in `.oursky/templates/rules/my-rules.md`:
```markdown
---
title: "My Team Standards"
applies_to: ["**/*.ts", "**/*.js"]
---

# Team Coding Standards
- Use TypeScript for new code
- Maximum function length: 30 lines
- Always add return types
```

**Enable for agents** by updating `.oursky/config.yaml`:
```yaml
# In your agent config sections
rules: ["my-rules"]
skills: ["my-custom-skill"]
```

### Environment Variables

Use environment variables for secrets and dynamic configuration:

```yaml
# In config.yaml
mcp:
  port: "{{.Env.MCP_PORT}}"
```

```bash
export MCP_PORT=3001
./oursky dev my-branch
```

## Example Usage

### Full-Stack Web Development

1. **Initialize project**:
   ```bash
   ./oursky init
   ```

2. **Configure for web development** (edit `.oursky/config.yaml`):
   ```yaml
   services:
     db:
       command: "docker-compose up -d postgres"
     api:
       command: "cd server && npm run dev"
       depends_on: ["db"]
     web:
       command: "cd client && npm run dev"
       depends_on: ["api"]

   agents:
     - name: "cursor"
       enabled: true
     - name: "opencode"
       enabled: true
   ```

3. **Start development**:
   ```bash
   ./oursky dev new-feature
   ```

4. **Code with AI assistance**:
   - Open `.oursky/worktrees/new-feature/` in Cursor
   - Use OpenCode with the generated config
   - All agents have access to your full development environment

## Dogfooding (Developing Oursky with Oursky)

Oursky is designed to be self-hosting. To develop the `oursky` project itself using an isolated environment:

1.  **Build the CLI**:
    ```bash
    go build -o oursky cmd/oursky/main.go
    ```
2.  **Initialize Oursky on Oursky**:
    ```bash
    ./oursky init
    ```
3.  **Configure agents** (optional, edit `.oursky/config.yaml`):
    ```yaml
    agents:
      - name: "cursor"
        enabled: true
    ```
4.  **Spin up a dev session**:
    ```bash
    ./oursky dev main
    ```
    *This creates a worktree at `.oursky/worktrees/main`, starts a container, generates AI agent configs, and mounts it.*
5.  **Connect your Agent**:
    Open `.oursky/worktrees/main/` in Cursor - the MCP connection is automatically configured.
6.  **Verify Changes**:
    Since the Docker socket is bind-mounted (DinD), you can run `./oursky` commands *inside* the Oursky container to spawn further sub-environments or run tests.

## Project Structure

### Code
- `pkg/ctrl`: Control Plane - Orchestration and lifecycle logic.
- `pkg/provider`: Execution Plane - Environment abstraction (Docker, LXC).
- `pkg/worktree`: Filesystem isolation using Git Worktrees.
- `pkg/agent`: Agent Plane - MCP server and configuration generation.

### Configuration (`.oursky/`)
- `config.yaml`: Main project configuration
- `templates/`: Shared AI capabilities (skills, commands, rules)
- `agents/`: Agent-specific configuration templates
- `worktrees/`: Auto-generated isolated development environments

### Documentation
- `docs/`: Technical specifications and planning
- `example/`: Complete working example project

## Roadmap

See [docs/planning/README.md](./docs/planning/README.md) for milestones and tasks.

- **M1: CLI MVP** (Current) - Docker + Worktree + MCP.
- **M2: Alpha** - LXC Provider, Advanced Port Forwarding.
- **M3: Beta** - QEMU/Virtualization for macOS/Windows.

---
*Powered by OhMyOpenCode.*
