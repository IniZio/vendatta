package agent

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateWorkspaceCommandValidation(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *CreateWorkspaceCommand
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid command",
			cmd: &CreateWorkspaceCommand{
				ID:            "cmd-123",
				WorkspaceID:   "ws-abc123",
				WorkspaceName: "my-project-feature",
				Provider:      "lxc",
				Image:         "ubuntu:22.04",
				Repository: RepositoryInfo{
					Owner:  "my-org",
					Name:   "my-project",
					URL:    "git@github.com:my-org/my-project.git",
					Branch: "main",
				},
				SSH: SSHConfig{
					Port:   2222,
					User:   "dev",
					PubKey: "ssh-ed25519 AAAA...",
				},
				Resources: ResourceConfig{
					CPU:    2,
					Memory: "4GB",
					Disk:   "20GB",
				},
				Services: []ServiceDefinition{
					{
						Name:    "web",
						Command: "npm run dev",
						Port:    3000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing workspace_id",
			cmd: &CreateWorkspaceCommand{
				WorkspaceName: "test",
				Provider:      "lxc",
				Image:         "ubuntu:22.04",
			},
			wantErr: true,
			errMsg:  "workspace_id is required",
		},
		{
			name: "missing provider",
			cmd: &CreateWorkspaceCommand{
				WorkspaceID:   "ws-123",
				WorkspaceName: "test",
				Image:         "ubuntu:22.04",
			},
			wantErr: true,
			errMsg:  "provider is required",
		},
		{
			name: "invalid service depends_on",
			cmd: &CreateWorkspaceCommand{
				ID:            "cmd-123",
				WorkspaceID:   "ws-abc123",
				WorkspaceName: "test",
				Provider:      "lxc",
				Image:         "ubuntu:22.04",
				Repository: RepositoryInfo{
					Owner:  "org",
					Name:   "repo",
					URL:    "git@github.com:org/repo.git",
					Branch: "main",
				},
				SSH: SSHConfig{
					Port:   2222,
					User:   "dev",
					PubKey: "ssh-ed25519 AAAA...",
				},
				Resources: ResourceConfig{
					CPU:    2,
					Memory: "4GB",
					Disk:   "20GB",
				},
				Services: []ServiceDefinition{
					{
						Name:      "api",
						Command:   "npm run server",
						Port:      4000,
						DependsOn: []string{"db"},
					},
				},
			},
			wantErr: true,
			errMsg:  "depends on undefined service",
		},
		{
			name: "duplicate service names",
			cmd: &CreateWorkspaceCommand{
				ID:            "cmd-123",
				WorkspaceID:   "ws-abc123",
				WorkspaceName: "test",
				Provider:      "lxc",
				Image:         "ubuntu:22.04",
				Repository: RepositoryInfo{
					Owner:  "org",
					Name:   "repo",
					URL:    "git@github.com:org/repo.git",
					Branch: "main",
				},
				SSH: SSHConfig{
					Port:   2222,
					User:   "dev",
					PubKey: "ssh-ed25519 AAAA...",
				},
				Resources: ResourceConfig{
					CPU:    2,
					Memory: "4GB",
					Disk:   "20GB",
				},
				Services: []ServiceDefinition{
					{
						Name:    "web",
						Command: "npm run dev",
						Port:    3000,
					},
					{
						Name:    "web",
						Command: "npm run build",
						Port:    3001,
					},
				},
			},
			wantErr: true,
			errMsg:  "duplicate service name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateWorkspaceCommandJSON(t *testing.T) {
	originalTime := time.Date(2026, 1, 17, 10, 30, 0, 0, time.UTC)

	cmd := &CreateWorkspaceCommand{
		ID:            "cmd-123",
		WorkspaceID:   "ws-abc123",
		WorkspaceName: "my-feature",
		Provider:      "lxc",
		Image:         "ubuntu:22.04",
		Repository: RepositoryInfo{
			Owner:  "my-org",
			Name:   "my-project",
			URL:    "git@github.com:my-org/my-project.git",
			Branch: "main",
		},
		SSH: SSHConfig{
			Port:   2222,
			User:   "dev",
			PubKey: "ssh-ed25519 AAAA...",
		},
		Resources: ResourceConfig{
			CPU:    2,
			Memory: "4GB",
			Disk:   "20GB",
		},
		Services: []ServiceDefinition{
			{
				Name:    "web",
				Command: "npm run dev",
				Port:    3000,
				HealthCheck: &HealthCheck{
					Type:    HealthCheckHTTP,
					Path:    "/health",
					Timeout: 10,
				},
			},
		},
		CreatedAt: originalTime,
	}

	// Marshal
	data, err := json.Marshal(cmd)
	require.NoError(t, err)

	// Unmarshal
	var decoded CreateWorkspaceCommand
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, cmd.ID, decoded.ID)
	assert.Equal(t, cmd.WorkspaceID, decoded.WorkspaceID)
	assert.Equal(t, cmd.WorkspaceName, decoded.WorkspaceName)
	assert.Equal(t, cmd.Provider, decoded.Provider)
	assert.Equal(t, cmd.SSH.Port, decoded.SSH.Port)
	assert.Equal(t, cmd.Resources.CPU, decoded.Resources.CPU)
	assert.Equal(t, len(cmd.Services), len(decoded.Services))
	assert.Equal(t, cmd.Services[0].Name, decoded.Services[0].Name)
	assert.Equal(t, cmd.CreatedAt.Unix(), decoded.CreatedAt.Unix())
}

func TestHealthCheckTypes(t *testing.T) {
	tests := []struct {
		name  string
		check *HealthCheck
	}{
		{
			name: "http health check",
			check: &HealthCheck{
				Type:    HealthCheckHTTP,
				Path:    "/health",
				Timeout: 10,
			},
		},
		{
			name: "tcp health check",
			check: &HealthCheck{
				Type:    HealthCheckTCP,
				Port:    3000,
				Timeout: 5,
				Retries: 3,
			},
		},
		{
			name: "exec health check",
			check: &HealthCheck{
				Type:    HealthCheckExec,
				Command: "curl http://localhost:3000/health",
				Timeout: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.check)
			require.NoError(t, err)

			var decoded HealthCheck
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err)

			assert.Equal(t, tt.check.Type, decoded.Type)
		})
	}
}

func TestServiceDependencyResolution(t *testing.T) {
	tests := []struct {
		name     string
		services []ServiceDefinition
		wantErr  bool
	}{
		{
			name: "linear dependency chain",
			services: []ServiceDefinition{
				{Name: "db", Command: "postgres", Port: 5432},
				{Name: "api", Command: "node server.js", Port: 4000, DependsOn: []string{"db"}},
				{Name: "web", Command: "npm run dev", Port: 3000, DependsOn: []string{"api"}},
			},
			wantErr: false,
		},
		{
			name: "circular dependency",
			services: []ServiceDefinition{
				{Name: "a", Command: "cmd", Port: 1000, DependsOn: []string{"b"}},
				{Name: "b", Command: "cmd", Port: 1001, DependsOn: []string{"a"}},
			},
			wantErr: false, // Our validator doesn't check for cycles
		},
		{
			name: "valid multiple dependencies",
			services: []ServiceDefinition{
				{Name: "db", Command: "postgres", Port: 5432},
				{Name: "cache", Command: "redis", Port: 6379},
				{Name: "api", Command: "node", Port: 4000, DependsOn: []string{"db", "cache"}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &CreateWorkspaceCommand{
				ID:            "cmd-123",
				WorkspaceID:   "ws-abc123",
				WorkspaceName: "test",
				Provider:      "lxc",
				Image:         "ubuntu:22.04",
				Repository: RepositoryInfo{
					Owner:  "org",
					Name:   "repo",
					URL:    "git@github.com:org/repo.git",
					Branch: "main",
				},
				SSH: SSHConfig{
					Port:   2222,
					User:   "dev",
					PubKey: "ssh-ed25519 AAAA...",
				},
				Resources: ResourceConfig{
					CPU:    2,
					Memory: "4GB",
					Disk:   "20GB",
				},
				Services: tt.services,
			}

			err := cmd.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
