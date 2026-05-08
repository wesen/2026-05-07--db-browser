package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "db-browser",
		Short: "Run Goja-backed SQLite browser scripts and repository JS verbs",
		Long: `db-browser is a Goja-backed playground for SQLite-focused web apps.

It scans configured JavaScript verb repositories, exposes explicit __verb__
functions as CLI commands, and will host Express-style database browser apps.`,
	}

	root.AddCommand(newServeCommand())
	root.AddCommand(newInspectCommand())
	root.AddCommand(newLazyVerbsCommand())
	return root
}

func newServeCommand() *cobra.Command {
	var addr string
	var dbPath string
	var scriptsDir string
	var dev bool

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve a Goja-backed SQLite browser web app",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("serve is not implemented yet (addr=%s db=%s scripts=%s dev=%v)", addr, dbPath, scriptsDir, dev)
		},
	}
	cmd.Flags().StringVar(&addr, "addr", ":8080", "HTTP listen address")
	cmd.Flags().StringVar(&dbPath, "db", "", "SQLite database path")
	cmd.Flags().StringVar(&scriptsDir, "scripts-dir", "./scripts", "Directory containing app JavaScript files")
	cmd.Flags().BoolVar(&dev, "dev", false, "Show detailed JavaScript errors in HTTP responses")
	return cmd
}

func newInspectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect db-browser runtime configuration",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "modules",
		Short: "List planned JavaScript modules",
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "database")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "db")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "fs")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "yaml")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "express")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "ui.dsl")
		},
	})
	return cmd
}

func newLazyVerbsCommand() *cobra.Command {
	return &cobra.Command{
		Use:                "verbs",
		Short:              "Run repository-scanned JavaScript verbs",
		DisableFlagParsing: true,
		Args:               cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("verbs are not implemented yet; received args: %v", args)
		},
	}
}
