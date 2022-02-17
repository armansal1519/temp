package orders

import (
	"bamachoub-backend-go-v1/app/cart"
	"bamachoub-backend-go-v1/app/transportation"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

// InitializeOrder Initialize order for first time
// @Summary Initialize order for first time
// @Description Initialize order for first time , use when user enter the "Checkout - Shipping info"
// @Tags order
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} Order{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /order/init [post]
func InitializeOrder(userKey string) (*Order, error) {
	groupedCart, err := cart.GetCartGroupByBuyingMethod(userKey)
	if err != nil {
		return nil, err
	}
	orderItems := make([]OrderItem, 0)
	for _, gc := range *groupedCart {
		var totalPrice int64
		for _, cartItem := range gc.Cart {
			totalPrice += cartItem.PricePerNumber * int64(cartItem.Number)
		}

		statusArr := make([]string, 0)
		for j := 0; j < len(gc.Cart); j++ {
			statusArr = append(statusArr, "wait-payment")
		}
		temp := OrderItem{
			Type:              gc.Type,
			Cart:              gc.Cart,
			TotalPrice:        totalPrice,
			RemainingPrice:    totalPrice,
			FromWallet:        0,
			PaymentKey:        "",
			StatusForEachItem: statusArr,
			Status:            "wait-payment",
		}
		orderItems = append(orderItems, temp)
	}
	//tp := transportation.GetTransportationPrice()
	o := Order{
		UserKey:                      userKey,
		OrderItems:                   orderItems,
		Status:                       "wait-payment",
		TransportationPrice:          0,
		IsTransportationPriceIsPayed: false,
		TransportationPriceWithPrice: false,
	}
	orderCol := database.GetCollection("order")
	meta, err := orderCol.CreateDocument(context.Background(), o)
	o.Key = meta.Key
	o.ID = meta.ID.String()
	return &o, nil
}

// InitializeOrder get order by jwt
// @Summary get order by jwt
// @Description get order by jwt
// @Tags order
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} Order{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /order/by-user [get]
func getOrderByUserKey(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	q := fmt.Sprintf("for o in order filter i.userKey==\"%v\" limit 1 return o", userKey)
	res := database.ExecuteGetQuery(q)
	if res == nil {
		return c.Status(404).JSON("order not found")
	}
	return c.JSON(res[0])
}

// InitializeOrder update order by sending info
// @Summary update order by sending info
// @Description update order by sending info , sending info must be created first ,returns order with transportation price
// @Tags order
// @Accept json
// @Produce json
// @Param data body sendingInfo true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} Order{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /order/sending-info [patch]
func addSendingInfoToOrder(c *fiber.Ctx) error {

	s := new(sendingInfo)
	if err := utils.ParseBodyAndValidate(c, s); err != nil {
		return c.JSON(err)
	}
	userKey := c.Locals("userKey").(string)

	orderCol := database.GetCollection("order")
	sendingInfoCol := database.GetCollection("sendingInfo")
	flag, err := orderCol.DocumentExists(context.Background(), s.OrderKey)
	if err != nil {
		return c.Status(404).JSON("error asserting order is exist")
	}
	if !flag {
		return c.Status(404).JSON("order not found")
	}
	flag, err = sendingInfoCol.DocumentExists(context.Background(), s.OrderKey)
	if err != nil {
		return c.Status(404).JSON("error asserting sending info is exist")
	}
	if !flag {
		return c.Status(404).JSON("sending info not found")
	}
	tp := transportation.GetTransportationPrice()
	uos := updateOrderBySendingInfo{
		TransportationPrice: tp,
		SendingInfoKey:      s.SendingInfoKey,
	}
	var o Order
	_, err = orderCol.ReadDocument(context.Background(), s.OrderKey, &o)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	if o.UserKey != userKey {
		return c.Status(403).JSON("Unauthorized")
	}

	ctx := driver.WithReturnNew(context.Background(), &o)
	_, err = orderCol.UpdateDocument(ctx, s.OrderKey, uos)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	return c.JSON(o)

}
