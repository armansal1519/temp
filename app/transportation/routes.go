package transportation

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/transportation")

	r.Post("/", middleware.Auth, getTransportationPrice)

}
