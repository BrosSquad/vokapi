package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/cobra"

	"github.com/BrosSquad/vokapi/container"
	"github.com/BrosSquad/vokapi/routes"
	"github.com/goccy/go-json"
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

	app := fiber.New(fiber.Config{
		JSONEncoder:          json.Marshal,
		JSONDecoder:          json.Unmarshal,
		Prefork:              false,
		ServerHeader:         "",
		Immutable:            false,
		AppName:              "vokapi",
		StrictRouting:        true,
		EnablePrintRoutes:    false,
		EnableIPValidation:   true,
		ProxyHeader:          "",
		CompressedFileSuffix: ".gz",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Accept-Encoding, Authorization, X-Request-ID",
		AllowMethods:     "GET, OPTIONS",
		AllowCredentials: true,
		MaxAge:           3600,
	}))
	app.Use(etag.New(etag.Config{
		Weak: false,
	}))
	app.Use(requestid.New())

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	go func(di *container.Container, app *fiber.App) {
		routes.Register(di, app)
		if err := app.Listen(fmt.Sprintf("%s:%d", host, port)); err != nil {
			log.Fatal("Cannot run server: " + err.Error())
		}
	}(di, app)

	<-sig

	if err := app.Shutdown(); err != nil {
		log.Printf("Cannot shutdown server: %s", err.Error())
	}

	if err := di.Close(); err != nil {
		log.Printf("Cannot close container: %s", err.Error())
	}

	return nil
}
