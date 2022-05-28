package supplierRequests

import (
	"bamachoub-backend-go-v1/app/suppliers"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// removeFromWallet create supplier request
// @Summary create supplier request
// @Description create supplier request
// @Tags  supplierRequest
// @Accept json
// @Produce json
// @Param   amount      path   string     true  "amount"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []string{}
// @Failure 400 {object} string
// @Router /supplier-request/remove-from-wallet/{amount} [post]
func removeFromWallet(c *fiber.Ctx) error {
	supplierKey := c.Locals("supplierId").(string)
	amount := c.Params("amount")
	s, err := suppliers.GetSupplierByKey(supplierKey)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	intAmount, err := strconv.Atoi(amount)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	if int64(intAmount) > s.WalletAmount {
		return c.Status(400).JSON("amount bigger than wallet")
	}
	r := sRequest{
		SupplierKey: supplierKey,
		Amount:      int64(intAmount),
		Type:        "remove from wallet",
	}

	srCol := database.GetCollection("supplierRequests")
	meta, err := srCol.CreateDocument(context.Background(), r)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.JSON(meta)
}
