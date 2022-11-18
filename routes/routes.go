package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/BrosSquad/vokapi/container"
	"github.com/BrosSquad/vokapi/handlers_v1/vokativ"
)

func Register(di *container.Container, app *fiber.App) {
	apiV1 := app.Group("/api/v1")

	vokativ.Register(di, apiV1.Group("/vokativ"))
}
