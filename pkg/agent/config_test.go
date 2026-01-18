package agent

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigWithFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "agent-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
coordination_url: http://localhost:3001
provider: docker
heartbeat:
  interval: 30s
  timeout: 10s
  retries: 3
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:3001", cfg.CoordinationURL)
	assert.Equal(t, "docker", cfg.Provider)
	assert.Equal(t, 30*time.Second, cfg.Heartbeat.Interval)
	assert.Equal(t, 10*time.Second, cfg.Heartbeat.Timeout)
	assert.Equal(t, 3, cfg.Heartbeat.Retries)
}

func TestLoadConfigDefaults(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "agent-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
coordination_url: http://test:3001
provider: docker
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	require.NoError(t, err)

	assert.NotNil(t, cfg)
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "nexus")
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := tmpDir + "/config.yaml"

	cfg := NodeConfig{
		CoordinationURL: "http://localhost:3001",
		Provider:        "docker",
		Heartbeat: HeartbeatConfig{
			Interval: 30 * time.Second,
			Timeout:  10 * time.Second,
			Retries:  3,
		},
	}

	err := SaveConfig(cfg, configPath)
	assert.NoError(t, err)

	loaded, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, cfg.CoordinationURL, loaded.CoordinationURL)
	assert.Equal(t, cfg.Provider, loaded.Provider)
}
