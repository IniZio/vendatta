package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vibegear/vendetta/pkg/ctrl"
	"github.com/vibegear/vendetta/pkg/provider"
)

var rootCmd = &cobra.Command{
	Use:   "vendatta",
	Short: "Vendatta - Isolated development environments",
	Long:  `Vendatta creates isolated development environments with AI agent integration.`,
}

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspaces",
	Long:  `Create, start, stop, and manage isolated workspaces.`,
}

var workspaceCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation would go here
		fmt.Printf("Creating workspace: %s\n", args[0])
	},
}

var workspaceUpCmd = &cobra.Command{
	Use:   "up <name>",
	Short: "Start a workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation would go here
		fmt.Printf("Starting workspace: %s\n", args[0])
	},
}

var workspaceDownCmd = &cobra.Command{
	Use:   "down <name>",
	Short: "Stop a workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation would go here
		fmt.Printf("Stopping workspace: %s\n", args[0])
	},
}

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation would go here
		fmt.Println("Listing workspaces...")
	},
}

var workspaceRmCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "Remove a workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation would go here
		fmt.Printf("Removing workspace: %s\n", args[0])
	},
}

var workspacePortsCmd = &cobra.Command{
	Use:   "ports <name>",
	Short: "Show workspace service ports",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// Create controller
		providers := []provider.Provider{}
		controller := ctrl.NewBaseController(providers, nil)

		// Get workspace services
		mappings, err := controller.WorkspaceServices(context.Background(), name)
		if err != nil {
			fmt.Printf("Error getting workspace services: %v\n", err)
			return
		}

		if len(mappings) == 0 {
			fmt.Printf("No services found for workspace '%s'\n", name)
			return
		}

		fmt.Printf("Services for workspace '%s':\n", name)
		for _, mapping := range mappings {
			fmt.Printf("  %s: %s (port %d -> %d)\n",
				mapping.ServiceName,
				mapping.URL,
				mapping.RemotePort,
				mapping.LocalPort)
		}
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceCreateCmd)
	workspaceCmd.AddCommand(workspaceUpCmd)
	workspaceCmd.AddCommand(workspaceDownCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceRmCmd)
	workspaceCmd.AddCommand(workspacePortsCmd)

	rootCmd.AddCommand(workspaceCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
