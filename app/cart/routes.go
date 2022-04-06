package cart

import (
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/middleware"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/Cart")
	r.Post("", middleware.AddToCartAuth, middleware.IsAuthenticated, func(c *fiber.Ctx) error {
		cIn := new(cartIn)
		if err := utils.ParseBodyAndValidate(c, cIn); err != nil {
			return c.JSON(err)
		}
		isAuthenticated := c.Locals("isAuthenticated").(bool)
		isLogin := c.Locals("isLogin").(bool)
		userKey := c.Locals("userKey").(string)
		tempUserKey := c.Get("temp-user-key", "")
		fmt.Println(*cIn)
		resp, err := addToCart(*cIn, isLogin, userKey, tempUserKey, isAuthenticated)
		if err.Status != -1 {
			return c.Status(err.Status).JSON(err)
		}
		return c.JSON(resp)
	})
	r.Get("/", middleware.AddToCartAuth, func(c *fiber.Ctx) error {

		isLogin := c.Locals("isLogin").(bool)
		userKey := c.Locals("userKey").(string)

		tempUserKey := c.Get("temp-user-key", "")

		if tempUserKey == "" && userKey == "" {
			return c.Status(400).SendString("user key missing ")
		}
		if userKey == "" {
			userKey = tempUserKey
		}

		resp, err := getCartByUserKey(isLogin, userKey)
		if err.Status != -1 {
			return c.Status(err.Status).JSON(err)
		}
		return c.JSON(resp)
	})

	r.Patch("/:key", middleware.AddToCartAuth, update)
	r.Delete("/:key", middleware.AddToCartAuth, remove)

}
