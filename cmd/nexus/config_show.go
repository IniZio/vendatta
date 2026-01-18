package main

import (
	"fmt"
	"os"

	"github.com/nexus/nexus/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show user configuration",
	Long:  `Display the current user configuration including GitHub credentials and SSH keys.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runConfigShow()
	},
}

func init() {
	if configCmd != nil {
		configCmd.AddCommand(configShowCmd)
	}
}

func runConfigShow() error {
	fmt.Println("📋 User Configuration")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	configPath := config.GetUserConfigPath()

	if _, err := os.Stat(configPath); err != nil {
		fmt.Println("No configuration found at", configPath)
		fmt.Println("Run 'nexus auth github' to get started")
		return nil
	}

	userCfg, err := config.LoadUserConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	data, err := yaml.Marshal(userCfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	fmt.Println(string(data))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Location: %s\n", configPath)

	return nil
}
