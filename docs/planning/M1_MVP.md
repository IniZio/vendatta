# Milestone: M1_MVP

**Objective**: A single-binary CLI that can initialize a project, spin up a Docker-based isolated codespace with a dedicated worktree, and expose it to an AI agent via MCP.

## 🎯 Success Criteria
- [x] `./vendatta init` scaffolds all necessary files.
- [x] `./vendatta dev <branch>` results in a running container with the worktree mounted.
- [x] Service ports are discovered and injected as environment variables.
- [x] **Multiple AI agents** (Cursor, OpenCode, Claude) can connect via generated configs and MCP server.

## 🛠 Tasks

| ID | Title | Priority | Status |
| :--- | :--- | :--- | :--- |
| **[INF-01](./tasks/INF-01.md)** | Docker Provider implementation (DinD, Ports) | 🔥 High | [Completed] |
| **[COR-01](./tasks/COR-01.md)** | Orchestration Engine (Worktree + Hooks) | 🔥 High | [Completed] |
| **[AGT-01](./tasks/AGT-01.md)** | MCP Server implementation | ⚡ Med | [Completed] |
| **[CLI-01](./tasks/CLI-01.md)** | CLI Scaffolding & Agent Config Generation | ⚡ Med | [Completed] |
| **[CLI-02](./tasks/CLI-02.md)** | Remote Repository Sync | 🟢 Low | [Completed] |
| **[VFY-01](./tasks/V2Y-01.md)** | E2E Verification with Docker-Compose | 🔥 High | [Pending] |
