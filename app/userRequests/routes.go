package userRequests

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/user-request")

	r.Post("/remove-from-wallet/:amount", middleware.Auth, removeFromWallet)
	r.Post("/admin/:key", middleware.Auth, approveUserRequest)

}
