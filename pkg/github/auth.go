package github

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func DetectGHCLI() (string, error) {
	path, err := exec.LookPath("gh")
	if err != nil {
		return "", fmt.Errorf("gh CLI not found in PATH: %w", err)
	}
	return path, nil
}

func CheckAuthStatus(ghCLIPath string) (bool, error) {
	cmd := exec.Command(ghCLIPath, "auth", "status")
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()
	return err == nil, nil
}

func ExecuteGHCommand(ghCLIPath string, args ...string) (string, error) {
	cmd := exec.Command(ghCLIPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("gh command failed: %w (stderr: %s)", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func AuthenticateWithGH(ghCLIPath string) error {
	cmd := exec.Command(ghCLIPath, "auth", "login")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
