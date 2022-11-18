package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

func registerCommands(root *cobra.Command) {
	root.AddCommand(InsertCommand())
	root.AddCommand(ServerCommand())
}

func Execute() {
	rootCmd := &cobra.Command{
		Use:     "vokapi",
		Short:   "vokapi",
		Long:    ``,
		Version: "0.0.1",
	}

	flags := rootCmd.PersistentFlags()

	flags.StringP("db-path", "d", "./db/database.badgerdb", "path to the database file")
	registerCommands(rootCmd)

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		log.Fatalf("error executing command: %v", err)
	}
}
