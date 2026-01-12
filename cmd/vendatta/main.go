package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vibegear/vendatta/pkg/ctrl"
	"github.com/vibegear/vendatta/pkg/provider"
	dockerProvider "github.com/vibegear/vendatta/pkg/provider/docker"
	lxcProvider "github.com/vibegear/vendatta/pkg/provider/lxc"
	"github.com/vibegear/vendatta/pkg/worktree"
)

var rootCmd = &cobra.Command{
	Use:   "vendatta",
	Short: "Vendatta - Isolated Development Environments",
	Long: `Vendatta eliminates the "it works on my machine" problem by providing
isolated, reproducible development environments that work seamlessly with
Coding Agents like Cursor, OpenCode, Claude, etc.`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vendatta project",
	Long:  `Initialize a new vendatta project by creating the .vendatta directory and default configuration files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		return controller.Init(ctx)
	},
}

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspaces",
	Long:  `Create, start, stop, and manage isolated development workspaces.`,
}

var workspaceCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new workspace",
	Long:  `Create a new workspace with the specified name. This will set up a Git worktree and generate AI agent configurations.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		return controller.WorkspaceCreate(ctx, args[0])
	},
}

var workspaceUpCmd = &cobra.Command{
	Use:   "up [name]",
	Short: "Start a workspace",
	Long:  `Start the specified workspace or auto-detect if no name is provided. This will create and start the isolated environment.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		return controller.WorkspaceUp(ctx, name)
	},
}

var workspaceDownCmd = &cobra.Command{
	Use:   "down [name]",
	Short: "Stop a workspace",
	Long:  `Stop the specified workspace or auto-detect if no name is provided.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		return controller.WorkspaceDown(ctx, name)
	},
}

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workspaces",
	Long:  `List all workspaces, showing their status and provider information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		return controller.WorkspaceList(ctx)
	},
}

var workspaceRmCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "Remove a workspace",
	Long:  `Remove the specified workspace, stopping it first if it's running.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		return controller.WorkspaceRm(ctx, args[0])
	},
}

var workspaceShellCmd = &cobra.Command{
	Use:   "shell [name]",
	Short: "Open shell in workspace",
	Long:  `Open an interactive shell in the specified workspace or auto-detect if no name is provided.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		return controller.WorkspaceShell(ctx, name)
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply latest configuration to agent configs",
	Long:  `Apply the latest vendatta configuration to all enabled AI agent configuration files (Cursor, OpenCode, Claude, etc.).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		return controller.Apply(ctx)
	},
}

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage plugins",
	Long:  `Manage plugins: add, remove, update, and list available plugins.`,
}

var pluginUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all plugins to latest versions",
	Long:  `Update all loaded plugins to their latest versions and refresh the lockfile.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		return controller.PluginUpdate(ctx)
	},
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all loaded plugins",
	Long:  `List all currently loaded plugins with their versions and status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		return controller.PluginList(ctx)
	},
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(pluginCmd)
	rootCmd.AddCommand(workspaceCmd)

	// Plugin subcommands
	pluginCmd.AddCommand(pluginUpdateCmd)
	pluginCmd.AddCommand(pluginListCmd)

	// Workspace subcommands
	workspaceCmd.AddCommand(workspaceCreateCmd)
	workspaceCmd.AddCommand(workspaceUpCmd)
	workspaceCmd.AddCommand(workspaceDownCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceRmCmd)
	workspaceCmd.AddCommand(workspaceShellCmd)
}

func createController() ctrl.Controller {
	// Create providers
	dockerProv, err := dockerProvider.NewDockerProvider()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Docker provider: %v\n", err)
		os.Exit(1)
	}

	lxcProv, err := lxcProvider.NewLXCProvider()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create LXC provider: %v\n", err)
		os.Exit(1)
	}

	providers := []provider.Provider{dockerProv, lxcProv}

	// Create worktree manager
	wtManager := worktree.NewManager(".", ".vendatta/worktrees")

	// Create controller
	return ctrl.NewBaseController(providers, wtManager)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
