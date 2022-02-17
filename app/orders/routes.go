package orders

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/order")

	r.Post("/init", middleware.Auth, func(c *fiber.Ctx) error {
		userKey := c.Locals("userKey").(string)
		resp, err := InitializeOrder(userKey)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(resp)
	})

	r.Get("/by-user", middleware.Auth, getOrderByUserKey)

	r.Patch("/sending-info", middleware.Auth, addSendingInfoToOrder)

}
