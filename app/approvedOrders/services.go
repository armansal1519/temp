package approvedOrders

import (
	"bamachoub-backend-go-v1/config/database"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func getApprovedOrders(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")

	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	q := fmt.Sprintf("for i in approvedOrder limit %v,%v return i ", offset, limit)
	//TODO add more filters
	return c.JSON(database.ExecuteGetQuery(q))
}

func getApprovedOrderByUserKey(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	q := fmt.Sprintf("for i in approvedOrder filter i.userKey==\"%v\" return i", userKey)
	return c.JSON(database.ExecuteGetQuery(q))
}
