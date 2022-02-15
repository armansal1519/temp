package estelam

import (
	"bamachoub-backend-go-v1/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/estelam")

	r.Post("/create", middleware.Auth, createEstelamRequest)
	r.Post("/resp", middleware.GetSupplierByEmployee, createEstelamRequest)

	r.Get("/user", middleware.Auth, func(c *fiber.Ctx) error {
		userKey := c.Locals("key").(string)
		resp, err := getEstelamCart(userKey)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(resp)

	})

	r.Get("/supplier", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		supplierId := c.Locals("supplierId").(string)

		resp, err := getEstelamForSupplier(supplierId)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(resp)

	})

	r.Get("/supplier/response", middleware.GetSupplierByEmployee, responseToEstelam)

}
