# Technical Specification: Lifecycle & Config

## 1. Project Configuration (`.vendatta/config.yaml`)

```yaml
name: project-name
provider: docker

# Service port definitions for discovery
services:
  web: 3000
  api: 8080

docker:
  image: node:20-alpine
  dind: true  # Enables Docker-in-Docker

agent:
  enabled: true

hooks:
  setup: .vendatta/hooks/setup.sh
  dev: .vendatta/hooks/dev.sh
```

## 2. Lifecycle States

### **`init`**
Scaffolds the `.vendatta` directory. Creates the base configuration and templates.

### **`workspace create <name>`**
1.  **Branch**: Creates or switches to the specified git branch.
2.  **Worktree**: Creates a git worktree in `.vendatta/worktrees/<name>/` (if `-w` flag used).
3.  **Agent Configs**: Generates AI agent configurations (Cursor, OpenCode, etc.) from templates.
4.  **Hooks**: Runs `.vendatta/hooks/create.sh` if it exists.

### **`workspace up [name]`**
1.  **Container**: Starts the Docker container with worktree bind-mounted.
2.  **Port Forwarding**: Maps service ports and injects `OURSKY_SERVICE_*` environment variables.
3.  **Hooks**: Executes `.vendatta/hooks/up.sh` if it exists.
4.  **Blocking**: Streams logs and maintains session until Ctrl+C (or detached with `-d`).

### **`workspace stop [name]`**
Stops the container but preserves state and resources.

### **`workspace down [name]`**
Stops and removes the container, networks, and temporary resources.

### **`workspace rm <name>`**
Deletes the worktree directory and all associated workspace resources.
