---
description: "Vendatta Project Overview and Conventions"
globs: ["**/*"]
alwaysApply: true
---

## OVERVIEW
Vendatta provides isolated, reproducible development environments for coding agents.

## STRUCTURE
- `cmd/vendatta/`: Main entry point
- `pkg/config/`: Configuration parsing and validation
- `pkg/ctrl/`: Orchestration logic and lifecycle management
- `pkg/templates/`: Template merging and rendering
- `pkg/provider/`: Environment providers (Docker, LXC)
- `.vendatta/`: Project-specific configuration and templates

## CONVENTIONS
- **Language**: Go 1.24
- **Configuration**: YAML-based, supports template merging from multiple sources
- **Agent Integration**: Automatic configuration generation via Model Context Protocol (MCP)
- **Isolated Environments**: Uses Git worktrees and container sessions

## DEVELOPMENT WORKFLOW
1. `vendatta init`: Initialize project
2. `vendatta workspace create <name>`: Create isolated environment
3. `vendatta workspace up <name>`: Start environment and generate agent configs
