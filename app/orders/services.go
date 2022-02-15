package orders

import (
	"bamachoub-backend-go-v1/app/cart"
	"bamachoub-backend-go-v1/app/transportation"
	"bamachoub-backend-go-v1/config/database"
	"context"
)

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
	tp := transportation.GetTransportationPrice()
	o := Order{
		UserKey:                      userKey,
		OrderItems:                   orderItems,
		Status:                       "wait-payment",
		TransportationPrice:          tp,
		IsTransportationPriceIsPayed: false,
		TransportationPriceWithPrice: false,
	}
	orderCol := database.GetCollection("order")
	meta, err := orderCol.CreateDocument(context.Background(), o)
	o.Key = meta.Key
	o.ID = meta.ID.String()
	return &o, nil
}
