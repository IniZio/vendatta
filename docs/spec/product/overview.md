# Product Overview: Project Oursky

> 📖 **Configuration Guide**: For detailed configuration options, see [Configuration Reference](./configuration.md)

## 1. Vision
Oursky is a developer-centric, local-first development environment manager. It aims to eliminate the "it works on my machine" problem by providing high-isolation, reproducible codespaces that are natively compatible with modern AI agents (Cursor, OpenCode).

## 2. Problem Statement
- **Environment Drift**: Developers struggle with conflicting node/python versions across projects.
- **Microservice Complexity**: Local development often requires complex `docker-compose` setups that are hard to isolate.
- **Agent Friction**: AI agents lack a standard, secure way to execute tools and understand environment-specific rules.
- **Bloat**: Existing solutions are often "heavy" (cloud-only) or "fragmented" (manual worktrees + manual docker).

## 3. Core Value Propositions

### **High Isolation (LXC/Docker + Worktrees)**
Unlike standard `docker-compose` which shares the host filesystem, Oursky uses **Git Worktrees** to provide a unique, branch-specific filesystem for every session. This prevents file-locking and allows parallel work on multiple branches without pollution.

### **BYOA (Bring Your Own Agent)**
Oursky provides a comprehensive **agent configuration generation system** with standardized **Model Context Protocol (MCP)** gateway. Any agent (Cursor, OpenCode, Claude Desktop/Code) automatically gets configured to connect to isolated environments with:
- **Environment-aware Tools**: MCP `exec` tool for secure command execution
- **Shared Standard Capabilities**: Reusable skills, commands, and rules following open standards
- **Dynamic Configuration**: Templates generate agent-specific configs with project context
- **Multi-Agent Support**: Simultaneous configuration for different AI tools

### **Single-Binary Portability**
Zero-setup installation. A single Go binary manages everything from worktree creation to container orchestration and port forwarding.

## 4. Target Personas

| Persona | Needs | Oursky Solution |
| :--- | :--- | :--- |
| **Senior Dev** | Complex orchestration, fast branch switching. | Automated Worktree + DinD orchestration. |
| **Agent User** | Secure tool execution for AI. | Built-in MCP Server with session boundaries. |
| **Team Lead** | Onboarding consistency. | Standardized `.vendatta` config & lifecycle hooks. |

## 5. Configuration Guide

### **Understanding the `.vendatta/` Structure**
Oursky uses a simple, intuitive configuration system centered around the `.vendatta/` directory:

```
.vendatta/
├── config.yaml          # Main project configuration
├── templates/           # Shared capabilities (skills, commands, rules)
│   ├── skills/          # Reusable AI skills
│   ├── commands/        # Standardized command definitions
│   └── rules/           # Development guidelines & rules
├── agents/              # Agent-specific configuration templates
│   ├── cursor/          # Cursor IDE settings
│   ├── opencode/        # OpenCode AI settings
│   ├── claude-desktop/  # Claude Desktop settings
│   └── claude-code/     # Claude Code CLI settings
└── worktrees/           # Generated isolated environments (auto-managed)
```

### **Main Configuration (`config.yaml`)**
The heart of your Oursky setup. Here's what you can configure:

```yaml
# Project identity
name: "my-awesome-project"
description: "Full-stack web application"

# Container & services
services:
  db:
    command: "docker-compose up -d postgres"
    healthcheck:
      url: "http://localhost:5432/health"
      interval: 5s
  api:
    command: "cd server && npm run dev"
    healthcheck:
      url: "http://localhost:5000/health"
    depends_on: ["db"]
  web:
    command: "cd client && npm run dev"
    healthcheck:
      url: "http://localhost:3000"
    depends_on: ["api"]

# AI agents to enable
agents:
  - name: "cursor"
    enabled: true
  - name: "opencode"
    enabled: true

# Remote sync targets
sync_targets:
  - name: "upstream"
    url: "https://github.com/example/upstream.git"

# MCP server configuration
mcp:
  enabled: true
  port: 3001
  host: "localhost"

# Container settings
docker:
  image: "ubuntu:22.04"
  dind: true  # Docker-in-Docker support

# Lifecycle hooks
hooks:
  setup: ".vendatta/hooks/setup.sh"
  dev: ".vendatta/hooks/dev.sh"
```

#### **Services Configuration**
Define your development stack:
- **command**: How to start each service
- **healthcheck**: Health monitoring with automatic retries
- **depends_on**: Service startup ordering

#### **Agent Configuration**
Choose which AI agents to support:
- **cursor**: VS Code extension with MCP support
- **opencode**: Standalone AI coding assistant
- **claude-desktop**: Anthropic's desktop application
- **claude-code**: Anthropic's CLI tool

#### **MCP (Model Context Protocol)**
Secure communication bridge between agents and your environment:
- **port/host**: Where the MCP server listens
- **enabled**: Toggle MCP functionality

### **Shared Templates**
Reusable capabilities that work across all agents:

#### **Skills** (`.vendatta/templates/skills/`)
AI capabilities following the [agentskills.io](https://agentskills.io) standard:
```yaml
name: "web-search"
description: "Search the web for information"
parameters:
  query: { type: "string", description: "Search query" }
execute:
  type: "http"
  url: "https://api.searchengine.com/search"
```

#### **Commands** (`.vendatta/templates/commands/`)
Standardized development workflows:
```yaml
name: "build"
description: "Build the project"
steps:
  - name: "Install dependencies"
    command: "npm install"
  - name: "Run build"
    command: "npm run build"
```

#### **Rules** (`.vendatta/templates/rules/`)
Development guidelines following [agents.md](https://github.com/agentsmd/agents.md):
```markdown
---
title: "Code Quality Standards"
applies_to: ["**/*.js", "**/*.ts"]
---

# Code Quality Standards

## General Principles
- Write clean, readable code
- Use meaningful variable names
- Add comments for complex logic
```

### **Agent-Specific Templates**
Each agent gets customized configuration:

#### **Cursor** (`.cursor/mcp.json`)
```json
{
  "mcpServers": {
    "our-project": {
      "type": "http",
      "url": "http://localhost:3001",
      "headers": {
        "Authorization": "Bearer YOUR_TOKEN"
      }
    }
  }
}
```

#### **OpenCode** (`opencode.json`)
```json
{
  "mcp": {
    "our-project": {
      "type": "remote",
      "url": "http://localhost:3001"
    }
  },
  "rules": ["code-quality", "collaboration"],
  "skills": ["web-search", "file-ops"],
  "commands": ["build", "deploy"]
}
```

### **Configuration Workflow**
1. **Initialize**: `vendatta init` creates the `.vendatta/` structure
2. **Customize**: Edit `config.yaml` and templates for your project
3. **Generate**: `vendatta dev <branch>` auto-generates agent configs
4. **Use**: Open worktree in your preferred AI agent
5. **Iterate**: Modify templates and regenerate as needed

## 6. The Workspace-Centric Workflow
1.  **Onboard**: Run `vendatta init`. Define services and capabilities.
2.  **Create**: Run `vendatta workspace create feature-x`. Creates branch, worktree, and generates agent configs.
3.  **Develop**: Run `vendatta workspace up feature-x`. Starts isolated container with port forwarding and lifecycle hooks.
4.  **Code**: Open worktree in any AI agent. Agents connect via their generated configs to MCP server.
5.  **Collaborate**: Multiple workspaces can run simultaneously with complete isolation.
6.  **Sync**: Use `vendatta config sync` for sharing configurations. Standard `git` for code.
7.  **Clean**: Use `workspace down` to stop, `workspace rm` to delete workspace entirely.
