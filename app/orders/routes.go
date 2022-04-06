package orders

//import (
//	"bamachoub-backend-go-v1/utils/middleware"
//	"github.com/gofiber/fiber/v2"
//	"strings"
//)
//
//func Routes(app fiber.Router) {
//	r := app.Group("/order")
//
//	r.Post("/init", middleware.Auth, func(c *fiber.Ctx) error {
//		userKey := c.Locals("userKey").(string)
//		resp, err := InitializeOrder(userKey)
//		if err != nil {
//			if strings.Contains(err.Error(), "unique constraint violated ") {
//				return c.Status(409).SendString("order must be unique")
//			}
//			return c.JSON(err)
//		}
//		return c.JSON(resp)
//	})
//
//	r.Get("/by-user", middleware.Auth, getOrderByUserKey)
//
//	r.Patch("/sending-info", middleware.Auth, addSendingInfoToOrder)
//
//}
