package ssh

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func UploadPublicKeyToGitHub(ghCLIPath, pubKeyPath string) error {
	_, err := ReadPublicKey(pubKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	cmd := exec.Command(ghCLIPath, "ssh-key", "add", pubKeyPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		errOutput := stderr.String()
		code := ParseGitHubSSHKeyError(errOutput)

		if code == 409 {
			return fmt.Errorf("ssh key already exists on github: %w", err)
		}

		return fmt.Errorf("failed to upload ssh key to github: %w (stderr: %s)", err, errOutput)
	}

	return nil
}

func ParseGitHubSSHKeyError(errOutput string) int {
	if strings.Contains(errOutput, "409") {
		return 409
	}
	if strings.Contains(errOutput, "401") {
		return 401
	}
	if strings.Contains(errOutput, "403") {
		return 403
	}

	for _, part := range strings.Split(errOutput, " ") {
		if code, err := strconv.Atoi(strings.TrimPrefix(part, "HTTP")); err == nil {
			return code
		}
	}

	return 0
}
