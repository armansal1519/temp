package approvedOrders

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/approved")

	r.Get("/", middleware.CheckAdmin, getApprovedOrders)
	r.Get("/user", middleware.Auth, getApprovedOrderByUserKey)

}
