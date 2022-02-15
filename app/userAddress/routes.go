package userAddress

import (
	"bamachoub-backend-go-v1/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/user-address")

	r.Get("/user", middleware.Auth, getAddressByUserKey)
	r.Get("/:key", middleware.Auth, func(c *fiber.Ctx) error {
		key := c.Params("key")
		resp, err := getAddressByKey(key)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(resp)
	})
	r.Post("", middleware.Auth, addAddress)
	r.Put("/:key", middleware.Auth, editAddress)
	r.Delete("/:key", middleware.Auth, removeAddress)

}
