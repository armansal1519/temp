package gOrderItem

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/order-req")

	r.Post("/cancel", middleware.Auth, cancelOrderItem)
	r.Post("/refer", middleware.Auth, referOrderItem)

}
