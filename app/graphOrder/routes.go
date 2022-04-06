package graphOrder

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Routes(app fiber.Router) {
	r := app.Group("/order")

	r.Post("/init", middleware.Auth, func(c *fiber.Ctx) error {
		userKey := c.Locals("userKey").(string)
		resp, err := InitOrder(userKey)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint violated ") {
				return c.Status(409).SendString("order must be unique")
			}
			return c.JSON(err)
		}
		return c.JSON(resp)
	})

	r.Get("/user", middleware.Auth, func(c *fiber.Ctx) error {
		userKey := c.Locals("userKey").(string)
		tab := c.Query("tab")
		offset := c.Query("offset")
		limit := c.Query("limit")
		if offset == "" || limit == "" {
			return c.Status(400).SendString("Offset and Limit must have a value")
		}
		resp, err := getOrdersByUserKey(userKey, tab, offset, limit)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(resp)
	})

	r.Get("/:key", middleware.Auth, getOrderByKey)

	r.Patch("/sending-info", middleware.Auth, addSendingInfoToOrder)

}
