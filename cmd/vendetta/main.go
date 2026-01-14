package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vibegear/vendetta/cmd/internal"
	"github.com/vibegear/vendetta/pkg/agent"
	"github.com/vibegear/vendetta/pkg/config"
	"github.com/vibegear/vendetta/pkg/coordination"
	"github.com/vibegear/vendetta/pkg/ctrl"
	"github.com/vibegear/vendetta/pkg/metrics"
	"github.com/vibegear/vendetta/pkg/provider"
	dockerProvider "github.com/vibegear/vendetta/pkg/provider/docker"
	lxcProvider "github.com/vibegear/vendetta/pkg/provider/lxc"
	"github.com/vibegear/vendetta/pkg/templates"
	"github.com/vibegear/vendetta/pkg/worktree"
	goyaml "gopkg.in/yaml.v3"
)

var (
	version   = "dev"
	goVersion = runtime.Version()
	buildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "vendetta",
	Short: "Isolated development environments that work with AI agents",
	Long: `Vendetta provides isolated development environments that integrate
seamlessly with AI coding assistants like Cursor, OpenCode, Claude, and others.`,
	Version: version,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of vendetta",
	Long:  `Print the version number, Go version, and build information of vendetta`,
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Printf("Vendetta Version: %s\n", version)
		fmt.Printf("Go Version: %s\n", goVersion)
		fmt.Printf("Build Date: %s\n", buildDate)
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vendetta project",
	Long:  `Initialize a new vendetta project by creating the .vendetta directory and default configuration files.`,
	RunE: func(_ *cobra.Command, _ []string) error {
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

var (
	workspaceNode string
)

var workspaceCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new workspace",
	Long:  `Create a new workspace with the specified name. This will set up a Git worktree and generate AI agent configurations.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
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
	RunE: func(_ *cobra.Command, args []string) error {
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
	RunE: func(_ *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		return controller.WorkspaceDown(ctx, name)
	},
}

var workspaceShellCmd = &cobra.Command{
	Use:   "shell [name]",
	Short: "Open shell in workspace",
	Long:  `Open an interactive shell in the specified workspace or auto-detect if no name is provided.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		return controller.WorkspaceShell(ctx, name)
	},
}

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workspaces",
	Long:  `List all workspaces, showing their status and provider information.`,
	RunE: func(_ *cobra.Command, _ []string) error {
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
	RunE: func(_ *cobra.Command, args []string) error {
		ctx := context.Background()
		controller := createController()
		return controller.WorkspaceRm(ctx, args[0])
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply latest configuration to agent configs",
	Long:  `Apply the latest vendetta configuration to all enabled AI agent configuration files (Cursor, OpenCode, Claude, etc.).`,
	RunE: func(_ *cobra.Command, _ []string) error {
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
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx := context.Background()
		controller := createController()
		return controller.PluginUpdate(ctx)
	},
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all loaded plugins",
	Long:  `List all currently loaded plugins with their versions and status.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx := context.Background()
		controller := createController()
		return controller.PluginList(ctx)
	},
}

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Analyze usage metrics and productivity",
	Long:  `Generate reports and insights from usage logs.`,
}

var usageSummaryCmd = &cobra.Command{
	Use:   "summary [date]",
	Short: "Generate daily summary report",
	Long:  `Generate a daily summary of usage metrics and insights. Date format: YYYY-MM-DD (defaults to today).`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return runUsageSummary(args)
	},
}

var usageMetricsCmd = &cobra.Command{
	Use:   "metrics [days]",
	Short: "Calculate productivity metrics",
	Long:  `Calculate detailed productivity metrics for the specified number of days (defaults to 7).`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return runUsageMetrics(args)
	},
}

var usagePatternsCmd = &cobra.Command{
	Use:   "patterns [days]",
	Short: "Analyze usage patterns",
	Long:  `Analyze usage patterns and trends for the specified number of days (defaults to 7).`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return runUsagePatterns(args)
	},
}

var usageBenchmarkCmd = &cobra.Command{
	Use:   "benchmark <baseline-days> <current-days>",
	Short: "Compare baseline and current metrics",
	Long:  `Compare productivity metrics between baseline period and current period.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		return runUsageBenchmark(args)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Manage vendetta configuration files.`,
}

var configExtractCmd = &cobra.Command{
	Use:   "extract <plugin-name>",
	Short: "Extract configuration to plugin",
	Long: `Extract local configuration (rules, skills, commands) into a reusable plugin.
This allows teams to share their coding standards and configurations as distributable plugins.`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		pluginName := args[0]

		// Default to extracting all types
		return internal.ExtractConfigToPlugin(pluginName, true, true, true)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all extends to latest versions",
	Long: `Fetch the latest versions of all configured extends and update the lockfile.
This ensures you have the most recent templates from remote repositories.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return runUpdate()
	},
}

var coordinationCmd = &cobra.Command{
	Use:   "coordination",
	Short: "Manage coordination server",
	Long:  `Start and manage the coordination server for remote node communication.`,
}

var coordinationStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the coordination server",
	Long: `Start the coordination server for remote node communication.
The coordination server manages remote nodes and enables distributed workspace execution.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Println("Starting coordination server...")
		configPath := coordination.GetConfigPath()
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configPath = ".vendetta/coordination.yaml"
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				fmt.Println("Generating default configuration...")
				if err := coordination.GenerateDefaultConfig(configPath); err != nil {
					return fmt.Errorf("failed to generate config: %w", err)
				}
				fmt.Printf("Configuration written to %s\n", configPath)
			}
		}
		return coordination.StartServer(configPath)
	},
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage node agent",
	Long:  `Start and manage the node agent for remote workspace execution.`,
}

var agentStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the node agent",
	Long: `Start the node agent to connect to a coordination server.
The node agent provides workspace execution capabilities on this machine.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Println("Starting node agent...")
		// Get coordination URL from environment or config
		coordURL := os.Getenv("VENDETTA_COORD_URL")
		if coordURL == "" {
			coordURL = "http://localhost:3001"
		}
		cfg := agent.NodeConfig{
			CoordinationURL: coordURL,
			Heartbeat: agent.HeartbeatConfig{
				Interval: 30 * time.Second,
				Timeout:  10 * time.Second,
				Retries:  3,
			},
		}
		agnt, err := agent.NewAgent(cfg)
		if err != nil {
			return fmt.Errorf("failed to create agent: %w", err)
		}
		return agnt.Start(context.Background())
	},
}

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage remote nodes",
	Long:  `Add, list, and manage remote nodes for distributed workspace execution.`,
}

var nodeAddCmd = &cobra.Command{
	Use:   "add <name> <host>",
	Short: "Add a remote node",
	Long: `Add a remote node to the configuration.
This stores the node configuration for later use with workspace creation.`,
	Args: cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		name := args[0]
		host := args[1]
		fmt.Printf("Adding node '%s' at %s...\n", name, host)

		// Auto-generate SSH key if needed
		pubKey, keyPath, err := ensureSSHKey()
		if err != nil {
			return err
		}

		// Store node configuration
		configDir := filepath.Join(os.Getenv("HOME"), ".config", "vendetta")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config dir: %w", err)
		}
		configPath := filepath.Join(configDir, "nodes.yaml")
		var nodes map[string]map[string]string
		data, err := os.ReadFile(configPath)
		if err == nil {
			goyaml.Unmarshal(data, &nodes)
		}
		if nodes == nil {
			nodes = make(map[string]map[string]string)
		}
		nodes[name] = map[string]string{
			"host":     host,
			"key_path": keyPath,
		}
		output, _ := goyaml.Marshal(nodes)
		if err := os.WriteFile(configPath, output, 0644); err != nil {
			return fmt.Errorf("failed to save node config: %w", err)
		}

		fmt.Printf("\n✅ Node '%s' added successfully!\n\n", name)
		fmt.Println("🔐 SSH Key Setup Required:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("To enable passwordless SSH access, copy your public key to the remote server:")
		fmt.Println("")
		fmt.Printf("  ssh-copy-id -i %s %s\n", keyPath, host)
		fmt.Println("")
		fmt.Println("Or manually add this key to ~/.ssh/authorized_keys on the remote:")
		fmt.Println("")
		fmt.Println(pubKey)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("\nAfter setting up SSH, you can create workspaces on this node:")
		fmt.Printf("  vendetta workspace create my-feature --node %s\n", name)
		return nil
	},
}

var nodeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured nodes",
	Long:  `List all configured remote nodes.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Println("Configured nodes:")
		configPath := filepath.Join(os.Getenv("HOME"), ".config", "vendetta", "nodes.yaml")
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Println("  No nodes configured")
			return nil
		}
		var nodes map[string]map[string]string
		goyaml.Unmarshal(data, &nodes)
		for name, cfg := range nodes {
			fmt.Printf("  - %s: %s\n", name, cfg["host"])
		}
		return nil
	},
}

// ensureSSHKey generates an SSH key for vendetta if one doesn't exist
func ensureSSHKey() (string, string, error) {
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	keyPath := filepath.Join(sshDir, "id_ed25519_vendetta")
	pubPath := keyPath + ".pub"

	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return "", "", fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		fmt.Println("🔑 Generating SSH key for vendetta remote access...")
		cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", keyPath, "-N", "", "-C", "vendetta@localhost")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", "", fmt.Errorf("failed to generate SSH key: %w, output: %s", err, output)
		}
		fmt.Println("✅ SSH key generated")
	}

	pubKey, err := os.ReadFile(pubPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read public key: %w", err)
	}

	return string(pubKey), keyPath, nil
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(pluginCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(usageCmd)
	rootCmd.AddCommand(workspaceCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(coordinationCmd)
	rootCmd.AddCommand(agentCmd)
	rootCmd.AddCommand(nodeCmd)

	// Plugin subcommands
	pluginCmd.AddCommand(pluginUpdateCmd)
	pluginCmd.AddCommand(pluginListCmd)

	// Config subcommands
	configCmd.AddCommand(configExtractCmd)

	// Usage subcommands
	usageCmd.AddCommand(usageSummaryCmd)
	usageCmd.AddCommand(usageMetricsCmd)
	usageCmd.AddCommand(usagePatternsCmd)
	usageCmd.AddCommand(usageBenchmarkCmd)

	// Workspace subcommands
	workspaceCmd.AddCommand(workspaceCreateCmd)
	workspaceCmd.AddCommand(workspaceUpCmd)
	workspaceCmd.AddCommand(workspaceDownCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceRmCmd)
	workspaceCmd.AddCommand(workspaceShellCmd)

	// Add --node flag to workspace commands
	workspaceCreateCmd.Flags().StringVarP(&workspaceNode, "node", "n", "", "Remote node to create workspace on")
	workspaceUpCmd.Flags().StringVarP(&workspaceNode, "node", "n", "", "Remote node to start workspace on")
	workspaceDownCmd.Flags().StringVarP(&workspaceNode, "node", "n", "", "Remote node to stop workspace on")
	workspaceShellCmd.Flags().StringVarP(&workspaceNode, "node", "n", "", "Remote node to shell into")

	// Coordination subcommands
	coordinationCmd.AddCommand(coordinationStartCmd)

	// Agent subcommands
	agentCmd.AddCommand(agentStartCmd)

	// Node subcommands
	nodeCmd.AddCommand(nodeAddCmd)
	nodeCmd.AddCommand(nodeListCmd)
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
	wtManager := worktree.NewManager(".", ".vendetta/worktrees")

	// Create controller
	return ctrl.NewBaseController(providers, wtManager)
}

// runUpdate updates all extends to their latest versions
func runUpdate() error {
	fmt.Println("📦 Updating extends to latest versions...")

	// Load config
	cfg, err := config.LoadConfig(".vendetta/config.yaml")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create templates manager
	templateManager := templates.NewManager(".")

	// Update each extend
	updated := 0
	for _, ext := range cfg.Extends {
		extStr, ok := ext.(string)
		if !ok {
			continue
		}
		parts := strings.Split(extStr, "/")
		if len(parts) != 2 {
			continue
		}

		owner, repo := parts[0], parts[1]

		// Parse optional branch
		branch := ""
		if strings.Contains(repo, "@") {
			repoParts := strings.SplitN(repo, "@", 2)
			repo = repoParts[0]
			branch = repoParts[1]
		}

		repoURL := fmt.Sprintf("https://github.com/%s/%s", owner, repo)
		repoTemplate := templates.TemplateRepo{
			URL:    repoURL,
			Branch: branch,
		}

		fmt.Printf("  Updating %s...\n", ext)
		if err := templateManager.PullWithUpdate(repoTemplate); err != nil {
			fmt.Printf("  ⚠️  Failed to update %s: %v\n", ext, err)
			continue
		}

		// Get new SHA
		sha, err := templateManager.GetRepoSHA(repoTemplate)
		if err != nil {
			fmt.Printf("  ⚠️  Failed to get SHA for %s: %v\n", ext, err)
			continue
		}

		fmt.Printf("  ✅ %s (SHA: %s)\n", ext, sha[:7])
		updated++
	}

	if updated == 0 {
		fmt.Println("No extends to update")
	} else {
		fmt.Printf("✅ Updated %d extends\n", updated)
	}

	return nil
}

func runUsageSummary(args []string) error {
	logger := metrics.NewLogger(".")
	reporter := metrics.NewReporter()

	var date time.Time
	if len(args) > 0 {
		var err error
		date, err = time.Parse("2006-01-02", args[0])
		if err != nil {
			return fmt.Errorf("invalid date format: %w (use YYYY-MM-DD)", err)
		}
	} else {
		date = time.Now()
	}

	summary, err := reporter.GenerateDailySummary(logger, date)
	if err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal summary: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func runUsageMetrics(args []string) error {
	days := 7
	if len(args) > 0 {
		parsed, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid days value: %w", err)
		}
		days = parsed
	}

	logger := metrics.NewLogger(".")
	reporter := metrics.NewReporter()

	m, summary, patterns, err := reporter.GenerateReport(logger, days)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	result := map[string]interface{}{
		"summary":  summary,
		"metrics":  m,
		"patterns": patterns,
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func runUsagePatterns(args []string) error {
	days := 7
	if len(args) > 0 {
		parsed, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid days value: %w", err)
		}
		days = parsed
	}

	logger := metrics.NewLogger(".")
	analyzer := metrics.NewAnalyzer()

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	logs, err := logger.Query(metrics.Filter{
		StartTime: startDate,
		EndTime:   endDate,
	})

	if err != nil {
		return fmt.Errorf("failed to query logs: %w", err)
	}

	patterns, err := analyzer.AnalyzePatterns(logs)
	if err != nil {
		return fmt.Errorf("failed to analyze patterns: %w", err)
	}

	data, err := json.MarshalIndent(patterns, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal patterns: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func runUsageBenchmark(args []string) error {
	baselineDays, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid baseline days: %w", err)
	}

	currentDays, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid current days: %w", err)
	}

	logger := metrics.NewLogger(".")
	reporter := metrics.NewReporter()

	comparison, err := reporter.GenerateBenchmark(logger, baselineDays, currentDays)
	if err != nil {
		return fmt.Errorf("failed to generate benchmark: %w", err)
	}

	data, err := json.MarshalIndent(comparison, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal benchmark: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
