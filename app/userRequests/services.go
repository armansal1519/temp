package userRequests

import (
	"bamachoub-backend-go-v1/app/paymentAndWallet"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

// removeFromWallet create user request
// @Summary create user request
// @Description create user request
// @Tags  userRequest
// @Accept json
// @Produce json
// @Param   amount      path   string     true  "amount"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []string{}
// @Failure 400 {object} string
// @Router /user-request/remove-from-wallet/{amount} [post]
func removeFromWallet(c *fiber.Ctx) error {
	amount := c.Params("amount")
	userKey := c.Locals("userKey").(string)
	intAmount, err := strconv.Atoi(amount)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	w := paymentAndWallet.UserWalletHistory{
		Amount:     int64(intAmount),
		UserKey:    userKey,
		PaymentKey: "-",
		CreatedAt:  time.Now().Unix(),
		Income:     false,
		TxType:     "remove-from-wallet",
		TxStatus:   "wait",
	}

	whCol := database.GetCollection("userWalletHistory")
	meta, err := whCol.CreateDocument(context.Background(), w)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	ur := userWalletRequest{
		UserKey:          userKey,
		Amount:           int64(intAmount),
		Type:             "remove from user wallet",
		WalletHistoryKey: meta.Key,
	}
	urCol := database.GetCollection("userRequests")
	meta, err = urCol.CreateDocument(context.Background(), ur)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.JSON(meta)
}

// removeFromWallet create user request
// @Summary create user request
// @Description create user request
// @Tags  userRequest
// @Accept json
// @Produce json
// @Param   key      path   string     true  "key"
// @Param   op      query   string     true  "done or reject"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []string{}
// @Failure 400 {object} string
// @Router /user-request/admin/{key} [post]
func approveUserRequest(c *fiber.Ctx) error {
	reqKey := c.Params("key")
	op := c.Query("op")
	if op != "reject" && op != "done" {
		return c.Status(400).JSON("op query must be reject or done")
	}
	col := database.GetCollection("userRequests")
	var ur userWalletRequest
	_, err := col.ReadDocument(context.Background(), reqKey, &ur)

	if err != nil {
		return c.Status(500).JSON(err)
	}

	u := updateWalletStatus{TxStatus: op}
	whCol := database.GetCollection("userWalletHistory")
	meta, err := whCol.UpdateDocument(context.Background(), ur.WalletHistoryKey, u)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.JSON(meta)

}
