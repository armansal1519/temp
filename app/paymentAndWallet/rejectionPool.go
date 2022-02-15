package paymentAndWallet

import (
	"bamachoub-backend-go-v1/config/database"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// filterPayment get all rejection pool items
// @Summary get all rejection pool items
// @Description get all rejection pool items
// @Tags rejection pool
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []rejectionPoolItemOut{}
// @Failure 404 {object} string{}
// @Router /rejection/pool [get]
func getRejectionPool(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	query := fmt.Sprintf("for i in rejectionPool  sort i.createdAt limit %v,%v return i", offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}
