package github

import (
	"bytes"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

func ParseRepoURL(repoString string) (owner string, repo string, err error) {
	repoString = strings.TrimSpace(repoString)

	if strings.Contains(repoString, "://") {
		u, err := url.Parse(repoString)
		if err != nil {
			return "", "", fmt.Errorf("invalid repo URL: %w", err)
		}

		path := strings.TrimPrefix(u.Path, "/")
		path = strings.TrimSuffix(path, ".git")

		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			return "", "", fmt.Errorf("invalid repo URL format: expected owner/repo")
		}

		return parts[0], parts[1], nil
	}

	parts := strings.Split(repoString, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repo format: expected owner/repo")
	}

	if parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo format: owner and repo must not be empty")
	}

	return parts[0], parts[1], nil
}

func VerifyRepoOwnership(ghCLIPath, owner, repo string) (bool, error) {
	cmd := exec.Command(ghCLIPath, "api", fmt.Sprintf("repos/%s/%s", owner, repo), "--jq", ".id")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("failed to verify repo ownership: %w", err)
	}

	output := strings.TrimSpace(stdout.String())
	return output != "" && output != "null", nil
}

func ForkRepository(ghCLIPath, owner, repo string) error {
	cmd := exec.Command(ghCLIPath, "repo", "fork", fmt.Sprintf("%s/%s", owner, repo), "--clone=false")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "already forked") {
			return nil
		}
		return fmt.Errorf("failed to fork repository: %w", err)
	}

	return nil
}

func BuildCloneURL(owner, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
}

func CloneRepository(url, destDir string) error {
	args := []string{"clone", url}
	if destDir != "" {
		args = append(args, destDir)
	}

	cmd := exec.Command("git", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w (stderr: %s)", err, stderr.String())
	}

	return nil
}

func GetRepositoryCommit(ghCLIPath, owner, repo string) (string, error) {
	cmd := exec.Command(ghCLIPath, "api", fmt.Sprintf("repos/%s/%s/commits/HEAD", owner, repo), "--jq", ".sha")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get repository commit: %w", err)
	}

	return strings.TrimSpace(stdout.String()), nil
}
