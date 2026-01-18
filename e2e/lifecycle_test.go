package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWorkspaceLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	projectDir := env.CreateTestProject(t, map[string]string{
		".nexus/config.yaml": `
name: lifecycle-test
provider: docker
services:
  web:
    command: "python3 -m http.server 28080"
    port: 28080
  api:
    command: "python3 -m http.server 23000"
    port: 23000
    depends_on: ["web"]
`,
		".nexus/hooks/up.sh": `#!/bin/bash
echo "Starting test services..."
timeout 60 /usr/bin/python3 -m http.server -b 0.0.0.0 18080 &
timeout 60 /usr/bin/python3 -m http.server -b 0.0.0.0 13000 &
echo "Services started"
`,
	})

	require.NoError(t, os.Chmod(filepath.Join(projectDir, ".nexus/hooks/up.sh"), 0755))

	binaryPath := env.BuildnexusBinary(t)
	env.RunnexusCommand(t, binaryPath, projectDir, "init")

	require.NoError(t, os.WriteFile(filepath.Join(projectDir, ".nexus/config.yaml"), []byte(`
name: lifecycle-test
provider: docker
services:
  db:
    command: "docker-compose up postgres"
  api:
    command: "npm start"
    depends_on: ["db"]
  web:
    command: "npm start"
    depends_on: ["api"]
`), 0644))

	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "create", "lifecycle-test")

	worktreePath := filepath.Join(projectDir, ".nexus", "worktrees", "lifecycle-test")
	_, err := os.Stat(worktreePath)
	require.NoError(t, err)

	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "up", "lifecycle-test")
	time.Sleep(3 * time.Second)

	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "down", "lifecycle-test")
	// down stops the container, rm removes the worktree
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", "lifecycle-test")

	_, err = os.Stat(worktreePath)
	require.Error(t, err)
}

func TestWorkspaceList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	projectDir := env.CreateTestProject(t, map[string]string{
		".nexus/config.yaml": `
name: list-test
provider: docker
services:
  app:
    command: "sleep infinity"
    port: 3000
`,
	})

	binaryPath := env.BuildnexusBinary(t)
	env.RunnexusCommand(t, binaryPath, projectDir, "init")

	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "create", "ws1")
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "create", "ws2")
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "create", "ws3")

	output := env.RunnexusCommand(t, binaryPath, projectDir, "branch", "list")

	for _, ws := range []string{"ws1", "ws2", "ws3"} {
		if !strings.Contains(output, ws) {
			t.Errorf("Expected workspace %s in list output", ws)
		}
	}

	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", "ws1")
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", "ws2")
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", "ws3")
}

func TestPluginSystem(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	projectDir := env.CreateTestProject(t, map[string]string{
		".nexus/hooks/up.sh": `#!/bin/bash
echo "Starting test environment..."
wait
`,
	})

	require.NoError(t, os.Chmod(filepath.Join(projectDir, ".nexus/hooks/up.sh"), 0755))

	binaryPath := env.BuildnexusBinary(t)
	env.RunnexusCommand(t, binaryPath, projectDir, "init")

	require.NoError(t, os.WriteFile(filepath.Join(projectDir, ".nexus/config.yaml"), []byte(`
name: plugin-test
provider: docker
agents:
  - name: "cursor"
    enabled: true
`), 0644))

	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "create", "plugin-test")

	worktreePath := filepath.Join(projectDir, ".nexus", "worktrees", "plugin-test")
	agentConfigPath := filepath.Join(worktreePath, "AGENTS.md")
	require.FileExists(t, agentConfigPath, "AGENTS.md should exist")

	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", "plugin-test")
}

func TestLXCProvider(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping LXC test in short mode")
	}

	if os.Getenv("LXC_TEST") == "" {
		t.Skip("Skipping LXC test - set LXC_TEST=1 to run")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	projectDir := env.CreateTestProject(t, map[string]string{
		".nexus/config.yaml": `
name: lxc-test
provider: lxc
services:
  app:
    command: "sleep infinity"
    port: 3000
lxc:
  image: ubuntu:22.04
`,
		".nexus/hooks/up.sh": `#!/bin/bash
echo "Starting LXC test environment..."
wait
`,
	})

	require.NoError(t, os.Chmod(filepath.Join(projectDir, ".nexus/hooks/up.sh"), 0755))

	binaryPath := env.BuildnexusBinary(t)
	env.RunnexusCommand(t, binaryPath, projectDir, "init")

	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "create", "lxc-ws")
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "up", "lxc-ws")
	time.Sleep(2 * time.Second)
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "down", "lxc-ws")
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", "lxc-ws")
}

func TestDockerProvider(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	projectDir := env.CreateTestProject(t, map[string]string{
		".nexus/config.yaml": `
name: docker-test
provider: docker
services:
  app:
    command: "sleep infinity"
    port: 3000
`,
		".nexus/hooks/up.sh": `#!/bin/bash
echo "Starting Docker test environment..."
wait
`,
	})

	require.NoError(t, os.Chmod(filepath.Join(projectDir, ".nexus/hooks/up.sh"), 0755))

	binaryPath := env.BuildnexusBinary(t)
	env.RunnexusCommand(t, binaryPath, projectDir, "init")

	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "create", "docker-ws")
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "up", "docker-ws")
	time.Sleep(2 * time.Second)
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "down", "docker-ws")
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", "docker-ws")
}

func TestErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	projectDir := env.CreateTestProject(t, map[string]string{
		".nexus/config.yaml": `
name: error-test
provider: docker
services:
  failing:
    command: "exit 1"
    port: 5000
`,
		".nexus/hooks/up.sh": `#!/bin/bash
echo "Starting test environment..."
wait
`,
	})

	require.NoError(t, os.Chmod(filepath.Join(projectDir, ".nexus/hooks/up.sh"), 0755))

	binaryPath := env.BuildnexusBinary(t)

	t.Log("Testing invalid workspace name...")
	output, err := env.RunnexusCommandWithError(binaryPath, projectDir, "branch", "create", "invalid/name")
	require.Error(t, err)
	require.Contains(t, output, "invalid")

	t.Log("Testing stop of non-existent workspace...")
	output, err = env.RunnexusCommandWithError(binaryPath, projectDir, "branch", "down", "nonexistent")
	require.Error(t, err)
	require.Contains(t, output, "not found")

	t.Log("Testing duplicate workspace creation...")
	_, err = env.RunnexusCommandWithError(binaryPath, projectDir, "branch", "create", "test-ws")
	require.Error(t, err)
}
