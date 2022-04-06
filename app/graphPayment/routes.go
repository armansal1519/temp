package graphPayment

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/gpayment")

	r.Post("/add-discount/:key/:paymentkey", middleware.Auth, addDiscountToPayment)

}
