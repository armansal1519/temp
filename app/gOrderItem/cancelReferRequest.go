package gOrderItem

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"github.com/gofiber/fiber/v2"
	"time"
)

// cancelOrderItem cancel order item
// @Summary cancel order item
// @Description cancel order item
// @Tags order item
// @Accept json
// @Produce json
// @Param data body orderItemCancelRequest true "data"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /order-req/cancel [post]
func cancelOrderItem(c *fiber.Ctx) error {
	m := new(orderItemCancelRequest)
	userKey := c.Locals("userKey").(string)
	if err := utils.ParseBodyAndValidate(c, m); err != nil {
		return c.JSON(err)
	}
	r := orderItemCancel{
		ProductId:    m.ProductId,
		Number:       m.Number,
		CancelReason: m.CancelReason,
		UserKey:      userKey,
		CreatedAt:    time.Now().Unix(),
		CancelAll:    m.CancelAll,
	}
	rCol := database.GetCollection("gOrderItemRequests")
	meta, err := rCol.CreateDocument(context.Background(), r)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.JSON(meta)
}

// referOrderItem refer order item
// @Summary refer order item
// @Description refer order item
// @Tags order item
// @Accept json
// @Produce json
// @Param data body OrderItemReferRequest true "data"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /order-req/refer [post]
func referOrderItem(c *fiber.Ctx) error {
	m := new(OrderItemReferRequest)
	userKey := c.Locals("userKey").(string)
	if err := utils.ParseBodyAndValidate(c, m); err != nil {
		return c.JSON(err)
	}
	r := OrderItemRefer{
		ProductId:          m.ProductId,
		Number:             m.Number,
		ReferReason:        m.ReferReason,
		ReferReasonDetails: m.ReferReasonDetails,
		UserKey:            userKey,
		CreatedAt:          time.Now().Unix(),
		ImageArr:           m.ImageArr,
		CancelAll:          m.CancelAll,
	}
	rCol := database.GetCollection("gOrderItemRequests")
	meta, err := rCol.CreateDocument(context.Background(), r)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.JSON(meta)
}
