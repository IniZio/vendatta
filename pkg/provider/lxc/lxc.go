package lxc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/vibegear/oursky/pkg/provider"
)

type LXCProvider struct {
	name string
}

func NewLXCProvider() (provider.Provider, error) {
	// Check if lxc command is available
	if _, err := exec.LookPath("lxc"); err != nil {
		return nil, fmt.Errorf("lxc command not found: %w", err)
	}
	return &LXCProvider{name: "lxc"}, nil
}

func (p *LXCProvider) Name() string {
	return p.name
}

func (p *LXCProvider) Create(ctx context.Context, sessionID string, workspacePath string, config interface{}) (*provider.Session, error) {
	containerName := fmt.Sprintf("oursky-%s", sessionID)

	// Launch LXC container
	cmd := exec.CommandContext(ctx, "lxc", "launch", "ubuntu:22.04", containerName)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to launch LXC container: %w", err)
	}

	// Wait for container to start
	cmd = exec.CommandContext(ctx, "lxc", "exec", containerName, "--", "sleep", "2")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to wait for container startup: %w", err)
	}

	session := &provider.Session{
		ID:       sessionID,
		Provider: p.name,
		Status:   "running",
		Labels: map[string]string{
			"oursky.session.id": sessionID,
		},
	}

	return session, nil
}

func (p *LXCProvider) Start(ctx context.Context, sessionID string) error {
	containerName := fmt.Sprintf("oursky-%s", sessionID)
	cmd := exec.CommandContext(ctx, "lxc", "start", containerName)
	return cmd.Run()
}

func (p *LXCProvider) Stop(ctx context.Context, sessionID string) error {
	containerName := fmt.Sprintf("oursky-%s", sessionID)
	cmd := exec.CommandContext(ctx, "lxc", "stop", containerName)
	return cmd.Run()
}

func (p *LXCProvider) Destroy(ctx context.Context, sessionID string) error {
	containerName := fmt.Sprintf("oursky-%s", sessionID)

	// Stop first if running
	p.Stop(ctx, sessionID)

	// Delete container
	cmd := exec.CommandContext(ctx, "lxc", "delete", containerName)
	return cmd.Run()
}

func (p *LXCProvider) List(ctx context.Context) ([]provider.Session, error) {
	cmd := exec.CommandContext(ctx, "lxc", "list", "--format", "csv", "-c", "n,s")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list LXC containers: %w", err)
	}

	var sessions []provider.Session
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines[1:] { // Skip header
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) >= 2 && strings.HasPrefix(parts[0], "oursky-") {
			sessionID := strings.TrimPrefix(parts[0], "oursky-")
			status := strings.ToLower(parts[1])

			sessions = append(sessions, provider.Session{
				ID:       sessionID,
				Provider: p.name,
				Status:   status,
				Labels: map[string]string{
					"oursky.session.id": sessionID,
				},
			})
		}
	}

	return sessions, nil
}

func (p *LXCProvider) Exec(ctx context.Context, sessionID string, opts provider.ExecOptions) error {
	containerName := fmt.Sprintf("oursky-%s", sessionID)

	args := []string{"exec", containerName, "--"}
	args = append(args, opts.Cmd...)

	cmd := exec.CommandContext(ctx, "lxc", args...)
	cmd.Env = append(os.Environ(), opts.Env...)

	if opts.Stdout {
		if opts.StdoutWriter != nil {
			cmd.Stdout = opts.StdoutWriter
		}
	}
	if opts.Stderr {
		if opts.StderrWriter != nil {
			cmd.Stderr = opts.StderrWriter
		}
	}

	return cmd.Run()
}
