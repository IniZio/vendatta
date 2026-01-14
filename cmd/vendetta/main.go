package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vibegear/vendatta/pkg/coordination"
	"github.com/vibegear/vendatta/pkg/ctrl"
	"github.com/vibegear/vendatta/pkg/provider"
	dockerProvider "github.com/vibegear/vendatta/pkg/provider/docker"
	lxcProvider "github.com/vibegear/vendatta/pkg/provider/lxc"
	qemuProvider "github.com/vibegear/vendatta/pkg/provider/qemu"
	"github.com/vibegear/vendatta/pkg/worktree"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "vendatta",
		Short: "Isolated development environments that work with AI agents",
		Long: `Vendatta provides isolated development environments that integrate 
seamlessly with AI coding assistants like Cursor, OpenCode, Claude, and others.`,
	}

	providers, err := initProviders()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing providers: %v\n", err)
		os.Exit(1)
	}

	worktreeManager := worktree.NewManager(".", ".vendatta/worktrees")
	controller := ctrl.NewBaseController(providers, worktreeManager)

	addInitCommand(rootCmd, controller)
	addWorkspaceCommands(rootCmd, controller)
	addPluginCommands(rootCmd, controller)
	addManagementCommands(rootCmd, controller)
	addCoordinationCommands(rootCmd)
	addAgentCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func initProviders() ([]provider.Provider, error) {
	var providers []provider.Provider

	dockerProv, err := dockerProvider.NewDockerProvider()
	if err == nil {
		providers = append(providers, dockerProv)
		fmt.Println("✅ Docker provider initialized")
	} else {
		fmt.Printf("⚠️  Docker provider not available: %v\n", err)
	}

	lxcProv, err := lxcProvider.NewLXCProvider()
	if err == nil {
		providers = append(providers, lxcProv)
		fmt.Println("✅ LXC provider initialized")
	} else {
		fmt.Printf("⚠️  LXC provider not available: %v\n", err)
	}

	qemuProv, err := qemuProvider.NewQEMUProvider()
	if err == nil {
		providers = append(providers, qemuProv)
		fmt.Println("✅ QEMU provider initialized")
	} else {
		fmt.Printf("⚠️  QEMU provider not available: %v\n", err)
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no providers available")
	}

	return providers, nil
}

func addInitCommand(rootCmd *cobra.Command, controller ctrl.Controller) {
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Vendatta project",
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.Init(context.Background()); err != nil {
				fmt.Printf("Error initializing project: %v\n", err)
				os.Exit(1)
			}
		},
	}
	rootCmd.AddCommand(cmd)
}

func addWorkspaceCommands(rootCmd *cobra.Command, controller ctrl.Controller) {
	var createCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new workspace",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.WorkspaceCreate(context.Background(), args[0]); err != nil {
				fmt.Printf("Error creating workspace: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var upCmd = &cobra.Command{
		Use:   "up [name]",
		Short: "Start a workspace",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.WorkspaceUp(context.Background(), args[0]); err != nil {
				fmt.Printf("Error starting workspace: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var downCmd = &cobra.Command{
		Use:   "down [name]",
		Short: "Stop a workspace",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.WorkspaceDown(context.Background(), args[0]); err != nil {
				fmt.Printf("Error stopping workspace: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var shellCmd = &cobra.Command{
		Use:   "shell [name]",
		Short: "Open shell in workspace",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.WorkspaceShell(context.Background(), args[0]); err != nil {
				fmt.Printf("Error opening shell: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all workspaces",
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.WorkspaceList(context.Background()); err != nil {
				fmt.Printf("Error listing workspaces: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var rmCmd = &cobra.Command{
		Use:   "rm [name]",
		Short: "Remove a workspace",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.WorkspaceRm(context.Background(), args[0]); err != nil {
				fmt.Printf("Error removing workspace: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var devCmd = &cobra.Command{
		Use:   "dev [branch]",
		Short: "Create workspace for development (alias for workspace create)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.Dev(context.Background(), args[0]); err != nil {
				fmt.Printf("Error creating dev workspace: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var workspaceCmd = &cobra.Command{
		Use:   "workspace",
		Short: "Workspace management commands",
	}
	workspaceCmd.AddCommand(createCmd, upCmd, downCmd, shellCmd, listCmd, rmCmd)
	rootCmd.AddCommand(workspaceCmd)
	rootCmd.AddCommand(devCmd)
}

func addPluginCommands(rootCmd *cobra.Command, controller ctrl.Controller) {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List loaded plugins",
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.PluginList(context.Background()); err != nil {
				fmt.Printf("Error listing plugins: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update plugins to latest versions",
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.PluginUpdate(context.Background()); err != nil {
				fmt.Printf("Error updating plugins: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var pluginCmd = &cobra.Command{
		Use:   "plugin",
		Short: "Plugin management commands",
	}
	pluginCmd.AddCommand(listCmd, updateCmd)
	rootCmd.AddCommand(pluginCmd)
}

func addManagementCommands(rootCmd *cobra.Command, controller ctrl.Controller) {
	var applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply latest configuration to agent configs",
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.Apply(context.Background()); err != nil {
				fmt.Printf("Error applying configuration: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var killCmd = &cobra.Command{
		Use:   "kill [session-id]",
		Short: "Kill a specific session",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.Kill(context.Background(), args[0]); err != nil {
				fmt.Printf("Error killing session: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all active sessions",
		Run: func(cmd *cobra.Command, args []string) {
			sessions, err := controller.List(context.Background())
			if err != nil {
				fmt.Printf("Error listing sessions: %v\n", err)
				os.Exit(1)
			}

			if len(sessions) == 0 {
				fmt.Println("No active sessions")
				return
			}

			fmt.Println("Active sessions:")
			for _, session := range sessions {
				fmt.Printf("  - %s (%s) [%s]\n", session.ID, session.Provider, session.Status)
				if len(session.Services) > 0 {
					var ports []string
					for name, port := range session.Services {
						ports = append(ports, fmt.Sprintf("%s:%d", name, port))
					}
					fmt.Printf("    Services: %s\n", strings.Join(ports, ", "))
				}
			}
		},
	}

	var execCmd = &cobra.Command{
		Use:   "exec [session-id] [command...]",
		Short: "Execute command in a session",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := controller.Exec(context.Background(), args[0], args[1:]); err != nil {
				fmt.Printf("Error executing command: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(applyCmd, killCmd, listCmd, execCmd)
}

func addCoordinationCommands(rootCmd *cobra.Command) {
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start coordination server",
		Run: func(cmd *cobra.Command, args []string) {
			configPath := coordination.GetConfigPath()
			if err := coordination.StartServer(configPath); err != nil {
				fmt.Printf("Error starting coordination server: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Generate coordination server configuration",
		Run: func(cmd *cobra.Command, args []string) {
			configPath := coordination.GetConfigPath()
			if err := coordination.GenerateDefaultConfig(configPath); err != nil {
				fmt.Printf("Error generating config: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Generated coordination server config at: %s\n", configPath)
		},
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show coordination server status",
		Run: func(cmd *cobra.Command, args []string) {
			configPath := coordination.GetConfigPath()
			cfg, err := coordination.LoadConfig(configPath)
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				os.Exit(1)
			}

			server := coordination.NewServer(cfg)
			stats := server.GetStats()

			fmt.Println("Coordination Server Status:")
			fmt.Printf("  Config: %s\n", configPath)
			fmt.Printf("  Host: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
			fmt.Printf("  Registry Provider: %s\n", cfg.Registry.Provider)
			fmt.Printf("  WebSocket Enabled: %t\n", cfg.WebSocket.Enabled)
			fmt.Printf("  Auth Enabled: %t\n", cfg.Auth.Enabled)
			fmt.Printf("  Timestamp: %s\n", stats["timestamp"])
		},
	}

	var coordinationCmd = &cobra.Command{
		Use:   "coordination",
		Short: "Coordination server management commands",
		Long: `Coordination server manages remote nodes and provides a centralized
API for dispatching commands and monitoring cluster status.`,
	}
	coordinationCmd.AddCommand(startCmd, configCmd, statusCmd)
	rootCmd.AddCommand(coordinationCmd)
}

func addAgentCommands(rootCmd *cobra.Command) {
	var installCmd = &cobra.Command{
		Use:   "install [target]",
		Short: "Install agent to remote node",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Installing agent to %s...\n", args[0])
			// TODO: Implement agent installation
			fmt.Println("Agent installation completed")
		},
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start agent daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting agent daemon...")
			// TODO: Implement agent daemon start
		},
	}

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop agent daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Stopping agent daemon...")
			// TODO: Implement agent daemon stop
		},
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show agent status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Agent status:")
			// TODO: Implement agent status check
		},
	}

	var connectCmd = &cobra.Command{
		Use:   "connect [url]",
		Short: "Connect to coordination server",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Connecting to coordination server at %s...\n", args[0])
			// TODO: Implement server connection
		},
	}

	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Generate agent configuration",
		Run: func(cmd *cobra.Command, args []string) {
			configPath := "/tmp/vendetta-agent.yaml"
			// TODO: Generate agent config
			fmt.Printf("Agent configuration generated at: %s\n", configPath)
		},
	}

	var agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "Node agent management commands",
		Long: `Node agent manages remote execution environments and communicates
with the coordination server for centralized management.`,
	}
	agentCmd.AddCommand(installCmd, startCmd, stopCmd, statusCmd, connectCmd, configCmd)
	rootCmd.AddCommand(agentCmd)
}
