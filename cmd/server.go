package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"

	"github.com/BrosSquad/vokapi/container"
	"github.com/BrosSquad/vokapi/routes"
)

var (
	port int
	host string
)

func ServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "server",
		RunE: func(command *cobra.Command, args []string) error {
			return executeServer(command)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&host, "host", "H", "0.0.0.0", "HTTP Servers Listening Address")
	flags.IntVarP(&port, "port", "p", 1389, "HTTP Server Port")
	return cmd
}

func executeServer(command *cobra.Command) error {
	ctx, cancel := context.WithCancel(command.Context())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer cancel()

	di := container.New(ctx, &container.Config{
		BadgedDBPath: command.Root().PersistentFlags().Lookup("db-path").Value.String(),
	})
	app := fiber.New()

	go func(di *container.Container, app *fiber.App) {
		routes.Register(di, app)

		err := app.Listen(fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			log.Fatal("Cannot run server: " + err.Error())
		}
	}(di, app)

	<-sig
	cancel()
	app.Shutdown()
	di.Close()

	return nil
}
