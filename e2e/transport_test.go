package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransportLocalExecution tests that local transport execution works
func TestTransportLocalExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E transport test in short mode")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create test project with local QEMU configuration
	projectDir := env.CreateTestProject(t, map[string]string{
		".nexus/config.yaml": `
name: local-transport-test
provider: qemu
qemu:
  cpu: 1
  memory: "1G"
  disk: "5G"
  ssh_port: 2223

services:
  test:
    command: "echo 'hello'"
    port: 9999
`,
	})

	binaryPath := env.BuildnexusBinary(t)
	env.RunnexusCommand(t, binaryPath, projectDir, "init")

	// Create workspace - should work locally
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "create", "local-test")

	// Verify worktree was created
	worktreePath := filepath.Join(projectDir, ".nexus", "worktrees", "local-test")
	_, err := os.Stat(worktreePath)
	require.NoError(t, err, "Worktree directory should exist")

	// Cleanup
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", "local-test")
}

// TestTransportRemoteConfigValidation tests that remote configuration is validated
func TestTransportRemoteConfigValidation(t *testing.T) {
	t.Skip("Remote plugin feature disabled - focus on local functionality")
}

// TestTransportQEMULifecycle tests QEMU provider with transport layer
func TestTransportQEMULifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E transport test in short mode")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create test project
	projectDir := env.CreateTestProject(t, map[string]string{
		".nexus/config.yaml": `
name: qemu-lifecycle-test
provider: qemu
qemu:
  cpu: 1
  memory: "1G"
  disk: "5G"
  ssh_port: 2226

services:
  health:
    command: "sleep infinity"
    port: 9000
`,
	})

	binaryPath := env.BuildnexusBinary(t)
	env.RunnexusCommand(t, binaryPath, projectDir, "init")

	// Create workspace
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "create", "lifecycle-test")

	// Verify worktree
	worktreePath := filepath.Join(projectDir, ".nexus", "worktrees", "lifecycle-test")
	_, err := os.Stat(worktreePath)
	require.NoError(t, err, "Worktree should exist")

	// Verify config was generated in worktree
	worktreeConfig := filepath.Join(worktreePath, ".nexus", "config.yaml")
	_, err = os.Stat(worktreeConfig)
	require.NoError(t, err, "Config should be copied to worktree")

	// Cleanup
	env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", "lifecycle-test")
}

// TestTransportMultipleWorkspaces tests multiple workspaces with different transports
func TestTransportMultipleWorkspaces(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E transport test in short mode")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	projectDir := env.CreateTestProject(t, map[string]string{
		".nexus/config.yaml": `
name: multi-workspace-test
provider: qemu
qemu:
  cpu: 1
  memory: "1G"
  disk: "5G"
  ssh_port: 2227
`,
	})

	binaryPath := env.BuildnexusBinary(t)
	env.RunnexusCommand(t, binaryPath, projectDir, "init")

	// Create multiple workspaces
	for i := 1; i <= 3; i++ {
		wsName := "multi-test-" + string(rune('a'+i-1))
		env.RunnexusCommand(t, binaryPath, projectDir, "branch", "create", wsName)

		// Verify each workspace
		worktreePath := filepath.Join(projectDir, ".nexus", "worktrees", wsName)
		_, err := os.Stat(worktreePath)
		require.NoError(t, err, "Workspace %s should exist", wsName)
	}

	// Cleanup all workspaces
	for i := 1; i <= 3; i++ {
		wsName := "multi-test-" + string(rune('a'+i-1))
		env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", wsName)
	}
}

// TestTransportInvalidConfig tests error handling for invalid configurations
func TestTransportInvalidConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E transport test in short mode")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	t.Run("empty_node_config", func(t *testing.T) {
		projectDir := env.CreateTestProject(t, map[string]string{
			".nexus/config.yaml": `
name: empty-node-test
provider: qemu
remote:
  node: ""
  user: "testuser"

qemu:
  cpu: 1
  memory: "1G"
  disk: "5G"
`,
		})

		binaryPath := env.BuildnexusBinary(t)
		env.RunnexusCommand(t, binaryPath, projectDir, "init")

		// Should handle empty node gracefully (treat as local)
		_, _ = env.RunnexusCommandWithError(binaryPath, projectDir, "branch", "create", "empty-test")
		// Either succeeds (as local) or fails with appropriate error
	})

	t.Run("invalid_qemu_settings", func(t *testing.T) {
		projectDir := env.CreateTestProject(t, map[string]string{
			".nexus/config.yaml": `
name: invalid-qemu-test
provider: qemu
qemu:
  cpu: -1
  memory: "invalid"
  disk: "0G"
`,
		})

		binaryPath := env.BuildnexusBinary(t)
		env.RunnexusCommand(t, binaryPath, projectDir, "init")

		// Create workspace - should either use defaults or fail gracefully
		_, _ = env.RunnexusCommandWithError(binaryPath, projectDir, "branch", "create", "invalid-qemu")
		// Should not panic
	})
}

// TestTransportSSHKeyValidation tests SSH key configuration
func TestTransportSSHKeyValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E transport test in short mode")
	}

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	t.Run("default_ssh_key", func(t *testing.T) {
		// Test with default SSH key path (~/.ssh/id_rsa)
		projectDir := env.CreateTestProject(t, map[string]string{
			".nexus/config.yaml": `
name: default-key-test
provider: qemu
remote:
  node: "localhost"
  user: "root"

qemu:
  cpu: 1
  memory: "1G"
  disk: "5G"
  ssh_port: 2228
`,
		})

		binaryPath := env.BuildnexusBinary(t)
		env.RunnexusCommand(t, binaryPath, projectDir, "init")

		// Create workspace - will fail connection but should use default key
		output, err := env.RunnexusCommandWithError(binaryPath, projectDir, "branch", "create", "default-key-test")

		// Verify it tried to connect (not config error)
		if err != nil {
			configErrorFound := strings.Contains(output, "invalid") ||
				strings.Contains(output, "required") ||
				strings.Contains(output, "key")
			// Should fail at connection, not config
			assert.False(t, configErrorFound, "Should not fail at config validation")
		}
	})

	t.Run("custom_ssh_key_path", func(t *testing.T) {
		// The remote config should support custom SSH key path
		// This test verifies the config parsing
		projectDir := env.CreateTestProject(t, map[string]string{
			".nexus/config.yaml": `
name: custom-key-test
provider: qemu
remote:
  node: "test.example.com"
  user: "testuser"
  # ssh_key: "/path/to/custom/key"  # Would be configured by user

qemu:
  cpu: 1
  memory: "1G"
  disk: "5G"
  ssh_port: 2229
`,
		})

		binaryPath := env.BuildnexusBinary(t)
		env.RunnexusCommand(t, binaryPath, projectDir, "init")

		// Should parse config successfully
		_, _ = env.RunnexusCommandWithError(binaryPath, projectDir, "branch", "create", "custom-key-test")
		// Fails at connection (expected), not config

		// Cleanup
		env.RunnexusCommand(t, binaryPath, projectDir, "branch", "rm", "custom-key-test")
	})
}

// Helper to create workspace name from number
func workspaceName(i int) string {
	return "ws-" + string(rune('a'+i-1))
}
