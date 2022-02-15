package cart

import (
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/cart")
	r.Post("", middleware.AddToCartAuth, middleware.IsAuthenticated, func(c *fiber.Ctx) error {
		cIn := new(cartIn)
		if err := utils.ParseBodyAndValidate(c, cIn); err != nil {
			return c.JSON(err)
		}
		isAuthenticated := c.Locals("isAuthenticated").(bool)
		log.Println(1)
		isLogin := c.Locals("isLogin").(bool)
		log.Println(isLogin)

		userKey := c.Locals("userKey").(string)
		log.Println(1)
		  

		resp, err := addToCart(*cIn, isLogin, userKey, isAuthenticated)
		if err.Status != -1 {
			return c.Status(err.Status).JSON(err)
		}
		return c.JSON(resp)
	})

}
