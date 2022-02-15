package paymentAndWallet

import (
	"bamachoub-backend-go-v1/app/orders"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
)

func updateOrderWithPaymentKey(order orders.Order, paymentKey string, txType string, includeTransportation bool) error {

	for i, item := range order.OrderItems {
		if item.Type == txType {
			order.OrderItems[i].PaymentKey = paymentKey
		}
	}

	if includeTransportation {
		order.TransportationPaymentId = fmt.Sprintf("payment/%v", paymentKey)
	}
	orderCol := database.GetCollection("order")
	_, err := orderCol.UpdateDocument(context.Background(), order.Key, order)
	if err != nil {
		return err
	}
	return nil
}
