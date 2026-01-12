# Task: INF-01 Docker Provider implementation

**Priority**: 🔥 High
**Status**: [Completed]

## 🎯 Objective
Implement the Docker-based environment provider with support for Docker-in-Docker (DinD) and dynamic port mapping.

## 🛠 Implementation Details
1.  **Mount Management**: 
    - Absolute path bind-mount for git worktree.
    - `/var/run/docker.sock` mount for DinD.
2.  **Privileged Mode**: Enable when DinD is requested in config.
3.  **Port Mapping**: 
    - Dynamic mapping of private service ports to available host ports.
    - Expose SSH (port 22).
4.  **Labeling**: Use `vendatta.session.id` for container discovery and tracking.

## 🧪 Proof of Work
- Container is created with `vendatta.session.id` label.
- Running `docker exec [id] docker version` works inside the container.
- Ports are mapped and visible in `vendatta list`.
