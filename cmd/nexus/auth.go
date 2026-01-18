package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nexus/nexus/pkg/config"
	"github.com/nexus/nexus/pkg/github"
	"github.com/nexus/nexus/pkg/ssh"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication with GitHub",
	Long:  `Authenticate with GitHub and manage SSH keys for secure access.`,
}

var authGitHubCmd = &cobra.Command{
	Use:   "github",
	Short: "Authenticate with GitHub",
	Long: `Authenticate with GitHub using the gh CLI.
This sets up your GitHub credentials for workspace creation and management.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runAuthGitHub()
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Display current authentication status for GitHub and SSH.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runAuthStatus()
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authGitHubCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(sshSetupCmd)
}

func runAuthGitHub() error {
	fmt.Println("🔐 GitHub Authentication")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	ghPath, err := github.DetectGHCLI()
	if err != nil {
		return fmt.Errorf("❌ gh CLI not found: %w\nPlease install it: https://cli.github.com", err)
	}

	authenticated, err := github.CheckAuthStatus(ghPath)
	if err != nil {
		return fmt.Errorf("failed to check auth status: %w", err)
	}

	if !authenticated {
		fmt.Println("ℹ️  Not currently authenticated with GitHub")
		fmt.Println("\nStarting authentication flow...")
		if err := github.AuthenticateWithGH(ghPath); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	fmt.Println("✅ Extracting GitHub user info...")
	userInfo, err := github.ExtractUserInfo(ghPath)
	if err != nil {
		return fmt.Errorf("failed to extract user info: %w", err)
	}

	fmt.Printf("   Username: %s\n", userInfo.Username)
	fmt.Printf("   User ID: %d\n", userInfo.UserID)

	configPath := config.GetUserConfigPath()
	configDir := filepath.Dir(configPath)

	if err := config.EnsureConfigDirectory(configDir); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	var userCfg *config.UserConfig
	if _, err := os.Stat(configPath); err == nil {
		userCfg, _ = config.LoadUserConfig(configPath)
		if userCfg == nil {
			userCfg = &config.UserConfig{}
		}
	} else {
		userCfg = &config.UserConfig{}
	}

	userCfg.GitHub.Username = userInfo.Username
	userCfg.GitHub.UserID = userInfo.UserID
	userCfg.GitHub.AvatarURL = userInfo.AvatarURL

	if err := config.SaveUserConfig(configPath, userCfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("")
	fmt.Println("✅ GitHub authentication complete!")
	fmt.Println("")
	fmt.Println("📝 Next steps:")
	fmt.Println("  1. Set up SSH keys: nexus ssh setup")
	fmt.Println("  2. Create a workspace: nexus workspace create <repo>")

	return nil
}

func runAuthStatus() error {
	fmt.Println("🔍 Authentication Status")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	ghPath, err := github.DetectGHCLI()
	if err != nil {
		fmt.Println("❌ gh CLI: Not found")
	} else {
		authenticated, _ := github.CheckAuthStatus(ghPath)
		if authenticated {
			userInfo, _ := github.ExtractUserInfo(ghPath)
			if userInfo != nil {
				fmt.Printf("✅ GitHub: Authenticated as @%s\n", userInfo.Username)
			} else {
				fmt.Println("⚠️  GitHub: CLI found but couldn't read user info")
			}
		} else {
			fmt.Println("⚠️  GitHub: Not authenticated")
		}
	}

	configPath := config.GetUserConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		userCfg, err := config.LoadUserConfig(configPath)
		if err == nil && userCfg != nil {
			fmt.Printf("✅ Saved config: %s\n", configPath)
			if userCfg.SSH.KeyPath != "" {
				fmt.Printf("✅ SSH key: %s\n", userCfg.SSH.KeyPath)
			}
		}
	} else {
		fmt.Println("⚠️  No saved configuration found")
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

var sshSetupCmd = &cobra.Command{
	Use:   "ssh-setup",
	Short: "Set up SSH keys and upload to GitHub",
	Long: `Generate or use existing SSH keys and upload them to GitHub.
This enables secure authentication for workspace access.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runSSHSetup()
	},
}

func runSSHSetup() error {
	fmt.Println("🔑 SSH Key Setup")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	home := os.Getenv("HOME")
	if home == "" {
		return fmt.Errorf("HOME environment variable not set")
	}

	sshDir := filepath.Join(home, ".ssh")

	hasEd25519, hasRSA, err := ssh.DetectExistingKeys(sshDir)
	if err != nil {
		return fmt.Errorf("failed to detect existing keys: %w", err)
	}

	var keyPath string
	if !hasEd25519 && !hasRSA {
		fmt.Println("Generating new ED25519 SSH key...")
		if err := ssh.EnsureSSHKey(sshDir, "ed25519"); err != nil {
			return fmt.Errorf("failed to generate SSH key: %w", err)
		}
		keyPath = filepath.Join(sshDir, "id_ed25519")
		fmt.Println("✅ SSH key generated")
	} else if hasEd25519 {
		keyPath = filepath.Join(sshDir, "id_ed25519")
		fmt.Println("✅ Using existing ED25519 key")
	} else {
		keyPath = filepath.Join(sshDir, "id_rsa")
		fmt.Println("✅ Using existing RSA key")
	}

	pubKey, err := ssh.ReadPublicKey(keyPath + ".pub")
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	fmt.Println("")
	fmt.Println("📋 Public Key:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println(pubKey)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fmt.Println("")
	fmt.Println("Uploading to GitHub...")

	ghPath, err := github.DetectGHCLI()
	if err != nil {
		fmt.Println("⚠️  gh CLI not found, skipping upload")
		fmt.Println("")
		fmt.Println("To upload manually:")
		fmt.Println("  1. Go to https://github.com/settings/keys")
		fmt.Println("  2. Click 'New SSH key'")
		fmt.Println("  3. Paste the public key above")
	} else {
		if err := ssh.UploadPublicKeyToGitHub(ghPath, keyPath+".pub"); err != nil {
			if err.Error() == "ssh key already exists on github: exit status 1" {
				fmt.Println("⚠️  SSH key already exists on GitHub")
			} else {
				return fmt.Errorf("failed to upload SSH key: %w", err)
			}
		} else {
			fmt.Println("✅ SSH key uploaded to GitHub")
		}
	}

	configPath := config.GetUserConfigPath()
	configDir := filepath.Dir(configPath)

	if err := config.EnsureConfigDirectory(configDir); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	var userCfg *config.UserConfig
	if _, err := os.Stat(configPath); err == nil {
		userCfg, _ = config.LoadUserConfig(configPath)
		if userCfg == nil {
			userCfg = &config.UserConfig{}
		}
	} else {
		userCfg = &config.UserConfig{}
	}

	userCfg.SSH.KeyPath = keyPath
	userCfg.SSH.PublicKey = pubKey

	if err := config.SaveUserConfig(configPath, userCfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("")
	fmt.Println("✅ SSH setup complete!")

	return nil
}
