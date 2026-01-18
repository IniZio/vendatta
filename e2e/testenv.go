package e2e

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestEnvironment provides a testing harness for E2E tests
type TestEnvironment struct {
	t          *testing.T
	tempDir    string
	binaryPath string
}

// NewTestEnvironment creates a new test environment
func NewTestEnvironment(t *testing.T) *TestEnvironment {
	tempDir, err := os.MkdirTemp("", "nexus-e2e-*")
	require.NoError(t, err)

	return &TestEnvironment{
		t:       t,
		tempDir: tempDir,
	}
}

// Cleanup removes the test environment
func (env *TestEnvironment) Cleanup() {
	env.cleanupDockerContainers()
	if env.tempDir != "" {
		os.RemoveAll(env.tempDir)
	}
}

func (env *TestEnvironment) cleanupDockerContainers() {
	dockerContainerIDsByLabel := func(label string) []string {
		cmd := exec.Command("docker", "ps", "-q", "--filter", "label="+label)
		if output, err := cmd.Output(); err == nil {
			return strings.Fields(string(output))
		}
		return nil
	}

	dockerContainerIDsByName := func(pattern string) []string {
		cmd := exec.Command("docker", "ps", "-q", "--filter", "name="+pattern)
		if output, err := cmd.Output(); err == nil {
			return strings.Fields(string(output))
		}
		return nil
	}

	for _, id := range dockerContainerIDsByLabel("nexus.session.id") {
		exec.Command("docker", "rm", "-f", id).Run()
	}

	for _, id := range dockerContainerIDsByName("my-project-") {
		exec.Command("docker", "rm", "-f", id).Run()
	}
}

// CreateTestProject creates a test project with the given files
func (env *TestEnvironment) CreateTestProject(t *testing.T, files map[string]string) string {
	projectDir := filepath.Join(env.tempDir, "project")
	require.NoError(t, os.MkdirAll(projectDir, 0755))

	for path, content := range files {
		fullPath := filepath.Join(projectDir, path)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0644))
	}

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = projectDir
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = projectDir
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = projectDir
	require.NoError(t, cmd.Run())

	// Create initial commit
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = projectDir
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = projectDir
	require.NoError(t, cmd.Run())

	return projectDir
}

// BuildnexusBinary builds the nexus binary for testing
func (env *TestEnvironment) BuildnexusBinary(t *testing.T) string {
	if env.binaryPath != "" {
		return env.binaryPath
	}

	// Find the nexus source directory by walking up from test executable
	// until we find a directory containing go.mod
	repoRoot := findRepoRoot(t)

	binaryPath := filepath.Join(env.tempDir, "nexus")
	cmd := exec.Command("go", "build", "-o", binaryPath, "cmd/nexus/main.go")
	cmd.Dir = repoRoot
	require.NoError(t, cmd.Run())

	env.binaryPath = binaryPath
	return binaryPath
}

// findRepoRoot finds repository root by looking for go.mod
func findRepoRoot(t *testing.T) string {
	dir, err := os.Getwd()
	require.NoError(t, err)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("Could not find repository root (go.mod not found)")
		}
		dir = parent
	}
}

// RunnexusCommand runs a nexus command and returns the output
func (env *TestEnvironment) RunnexusCommand(t *testing.T, binaryPath, projectDir string, args ...string) string {
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = projectDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err != nil {
		t.Logf("Command failed: %s %v in %s", binaryPath, args, projectDir)
		t.Logf("Output: %s", output)
		require.NoError(t, err)
	}
	return output
}

// RunnexusCommandWithError runs a nexus command and returns output and error
func (env *TestEnvironment) RunnexusCommandWithError(binaryPath, projectDir string, args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = projectDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	return output, err
}

// VerifyServiceHealth checks if a service is healthy
func (env *TestEnvironment) VerifyServiceHealth(t *testing.T, url string, timeout time.Duration) {
	client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{DisableKeepAlives: true}}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return
		}
		if err != nil {
			t.Logf("Error connecting to %s: %v", url, err)
		} else {
			t.Logf("Service at %s returned status %d", url, resp.StatusCode)
			resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf("Service at %s did not become healthy within %v", url, timeout)
}

// VerifyServiceDown checks if a service is down
func (env *TestEnvironment) VerifyServiceDown(t *testing.T, url string, timeout time.Duration) {
	client := &http.Client{Timeout: 5 * time.Second}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err != nil {
			return
		}
		resp.Body.Close()
		if resp.StatusCode >= 500 {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf("Service at %s did not go down within %v", url, timeout)
}

// VerifyFileExists checks if a file exists
func (env *TestEnvironment) VerifyFileExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	require.NoError(t, err, "File should exist: %s", path)
}

// GenerateTestSSHKey generates a test RSA private key and returns the path
func (env *TestEnvironment) GenerateTestSSHKey(t *testing.T) string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	keyPath := filepath.Join(env.tempDir, "test_rsa")
	err = os.WriteFile(keyPath, privateKeyPem, 0600)
	require.NoError(t, err)

	return keyPath
}
