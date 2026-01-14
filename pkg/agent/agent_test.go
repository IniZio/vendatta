package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("/tmp/test-config.yaml")
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:3001", config.CoordinationURL)
	assert.Equal(t, "docker", config.Provider)
	assert.Equal(t, 30*time.Second, config.Heartbeat.Interval)
}

func TestNewAgent(t *testing.T) {
	config := NodeConfig{
		CoordinationURL: "http://test:3001",
		Provider:        "test",
		Heartbeat: HeartbeatConfig{
			Interval: 10 * time.Second,
		},
	}

	agent, err := NewAgent(config)
	require.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, "test", agent.config.Provider)
	assert.NotEmpty(t, agent.node.ID)
}

func TestNewExecutor(t *testing.T) {
	config := NodeConfig{
		CoordinationURL: "http://test:3001",
		Provider:        "test",
	}

	agent, err := NewAgent(config)
	require.NoError(t, err)

	executor := NewExecutor(agent)
	assert.NotNil(t, executor)
	assert.Equal(t, agent, executor.agent)
}

func TestExecutorSessionCommand(t *testing.T) {
	config := NodeConfig{
		CoordinationURL: "http://test:3001",
		Provider:        "test",
	}

	agent, err := NewAgent(config)
	require.NoError(t, err)

	executor := NewExecutor(agent)

	cmd := Command{
		ID:     "test-cmd-1",
		Type:   "session",
		Action: "list",
	}

	result := executor.ExecuteSessionCommand(cmd)
	assert.Equal(t, "success", result.Status)
	assert.NotEmpty(t, result.Output)
}

func TestExecutorServiceCommand(t *testing.T) {
	config := NodeConfig{
		CoordinationURL: "http://test:3001",
		Provider:        "test",
	}

	agent, err := NewAgent(config)
	require.NoError(t, err)

	executor := NewExecutor(agent)

	cmd := Command{
		ID:     "test-cmd-2",
		Type:   "service",
		Action: "list",
	}

	result := executor.ExecuteServiceCommand(cmd)
	assert.Equal(t, "success", result.Status)
	assert.NotEmpty(t, result.Output)
}

func TestExecutorSystemCommand(t *testing.T) {
	config := NodeConfig{
		CoordinationURL: "http://test:3001",
		Provider:        "test",
	}

	agent, err := NewAgent(config)
	require.NoError(t, err)

	executor := NewExecutor(agent)

	cmd := Command{
		ID:     "test-cmd-3",
		Type:   "system",
		Action: "status",
	}

	result := executor.ExecuteSystemCommand(cmd)
	assert.Equal(t, "success", result.Status)
	assert.NotEmpty(t, result.Output)
}
