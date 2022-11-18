package cmd

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"os/signal"

	"github.com/BrosSquad/vokapi/container"
	"github.com/spf13/cobra"
)

var csvPath string

func InsertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "insert",
		RunE: func(command *cobra.Command, args []string) error {
			return executeInsert(command)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&csvPath, "csv-path", "c", "./imena.csv", "Path to CSV file")

	return cmd
}

func executeInsert(command *cobra.Command) error {
	ctx, cancel := context.WithCancel(command.Context())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer cancel()

	di := container.New(ctx, &container.Config{
		BadgedDBPath: command.Root().PersistentFlags().Lookup("db-path").Value.String(),
	})
	db := di.GetBadgerDB()
	writer := db.NewWriteBatch()

	done := make(chan struct{}, 1)

	go func() {
		file, err := os.OpenFile(csvPath, os.O_RDONLY, 0o644)
		if err != nil {
			log.Fatalf("Error opening file %s: %v", csvPath, err)
		}

		reader := csv.NewReader(file)
		reader.TrimLeadingSpace = true
		reader.FieldsPerRecord = 2

		data, err := reader.ReadAll()
		if err != nil {
			log.Fatalf("Error reading file %s: %v", csvPath, err)
		}

		items := data[1:]

		for _, item := range items {
			if err := writer.Set([]byte(item[0]), []byte(item[1])); err != nil {
				log.Printf("Error inserting item %s: %v", item[0], err)
			}
		}

		if err := writer.Flush(); err != nil {
			log.Fatalf("Error flushing batch: %v", err)
		}

		done <- struct{}{}
	}()

	select {
	case <-sig:
		cancel()
		log.Println("Interrupted")
		writer.Cancel()
		return nil
	case <-done:
		log.Println("Done")
		return nil
	}
}
