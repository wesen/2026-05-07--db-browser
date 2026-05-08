package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-go-golems/db-browser/internal/app"
	"github.com/go-go-golems/db-browser/internal/verbcli"
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
	root.AddCommand(verbcli.NewLazyCommand())
	return root
}

func newServeCommand() *cobra.Command {
	var addr string
	var dbPath string
	var scriptsDir string
	var dev bool
	var readonly bool
	var allowWrites bool

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve a Goja-backed SQLite browser web app",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			server, err := app.NewServer(ctx, app.Config{Addr: addr, DBPath: dbPath, ScriptsDir: scriptsDir, Dev: dev, ReadOnly: readonly, AllowWrites: allowWrites})
			if err != nil {
				return err
			}
			defer server.Close(context.Background())
			fmt.Fprintf(cmd.ErrOrStderr(), "serving db-browser on %s\n", addr)
			return server.Run(ctx)
		},
	}
	cmd.Flags().StringVar(&addr, "addr", ":8080", "HTTP listen address")
	cmd.Flags().StringVar(&dbPath, "db", "", "SQLite database path")
	cmd.Flags().StringVar(&scriptsDir, "scripts-dir", "./scripts", "Directory containing app JavaScript files")
	cmd.Flags().BoolVar(&dev, "dev", false, "Show detailed JavaScript errors in HTTP responses")
	cmd.Flags().BoolVar(&readonly, "readonly", true, "Disable db.exec writes in served scripts")
	cmd.Flags().BoolVar(&allowWrites, "allow-writes", false, "Allow db.exec writes when --readonly=false")
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
