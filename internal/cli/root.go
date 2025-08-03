package cli

import (
	"os"

	"github.com/bxrne/goforge/internal/cli/commands"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "goforge",
		Short: "GoForge: A CLI for scaffolding Go applications",
		Long:  `GoForge is a CLI tool to streamline Go app development with pre-configured templates and integrations.`,
	}

	// Add subcommands
	rootCmd.AddCommand(commands.NewInitCommand())
	rootCmd.AddCommand(commands.NewAddCommand())
	rootCmd.AddCommand(commands.NewGenerateCommand())

	return rootCmd
}

func Execute() {
	rootCmd := NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
