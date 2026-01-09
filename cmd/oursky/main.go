package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vibegear/oursky/pkg/agent"
	"github.com/vibegear/oursky/pkg/ctrl"
	"github.com/vibegear/oursky/pkg/provider"
	"github.com/vibegear/oursky/pkg/provider/docker"
)

func main() {
	var providers []provider.Provider
	dProvider, err := docker.NewDockerProvider()
	if err == nil {
		providers = append(providers, dProvider)
	}

	controller := ctrl.NewBaseController(providers)

	rootCmd := &cobra.Command{
		Use:   "oursky",
		Short: "Oursky Dev Environment Manager",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize .oursky in the current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return controller.Init(context.Background())
		},
	}

	devCmd := &cobra.Command{
		Use:   "dev [branch]",
		Short: "Start a development session for a branch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return controller.Dev(context.Background(), args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List active sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			sessions, err := controller.List(context.Background())
			if err != nil {
				return err
			}
			for _, s := range sessions {
				fmt.Printf("%s\t%s\t%s\n", s.Labels["oursky.session.id"], s.Provider, s.Status)
			}
			return nil
		},
	}

	killCmd := &cobra.Command{
		Use:   "kill [session-id]",
		Short: "Stop and destroy a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return controller.Kill(context.Background(), args[0])
		},
	}

	agentCmd := &cobra.Command{
		Use:   "agent [session-id]",
		Short: "Start MCP agent gateway for a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sessions, _ := controller.List(context.Background())
			var targetSession *provider.Session
			for _, s := range sessions {
				if s.ID == args[0] || s.Labels["oursky.session.id"] == args[0] {
					targetSession = &s
					break
				}
			}
			if targetSession == nil {
				return fmt.Errorf("session %s not found", args[0])
			}

			p, ok := controller.Providers[targetSession.Provider]
			if !ok {
				return fmt.Errorf("provider %s not found", targetSession.Provider)
			}

			s := agent.NewAgentServer(targetSession.ID, p)
			return s.Serve()
		},
	}

	rootCmd.AddCommand(initCmd, devCmd, listCmd, killCmd, agentCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
