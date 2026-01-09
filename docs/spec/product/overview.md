# Product Overview: Project Oursky

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
| **Team Lead** | Onboarding consistency. | Standardized `.oursky` config & lifecycle hooks. |

## 5. The "Oursky" Workflow
1.  **Onboard**: Run `oursky init`. Define services, agents, and MCP settings.
2.  **Configure**: CLI generates agent configs (Cursor `.cursor/mcp.json`, OpenCode `opencode.json`, etc.) from templates.
3.  **Dev**: Run `oursky dev feature-x`. Clean worktree + container + MCP server start.
4.  **Code**: Open worktree in any AI agent. Automatic MCP connection with full capabilities.
5.  **Collaborate**: Multiple agents can work simultaneously with isolated environments.
6.  **Clean**: Run `oursky kill`. All resources and generated configs are wiped.
