package discountCode

import "github.com/gofiber/fiber/v2"

func Routes(app fiber.Router) {
	r := app.Group("/discount")

	r.Post("/", createDiscountForPhoneNumbers)
	r.Get("/:key", getDiscountByKey)

}
