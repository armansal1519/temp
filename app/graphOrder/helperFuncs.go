package graphOrder

import (
	"bamachoub-backend-go-v1/app/graphPayment"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
)

func GetPaymentForOrderWithOrderKey(orderKey string) (*[]graphPayment.GPaymentOut, error) {
	q := fmt.Sprintf(" for o in gOrder  filter o._key==\"%v\" for v,e in 1..1 outbound o graph \"orderGraph\"  return v ", orderKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var data []graphPayment.GPaymentOut
	for {
		var doc graphPayment.GPaymentOut
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return &data, nil
}

func GetOrderPaymentAndOrderItem(orderKey string, orderItemKey string) (*OrderPaymentAndOrderItem, error) {
	q := fmt.Sprintf("for o in gOrder  filter o._key==\"%v\"\nfor v,e,p in 2..2 outbound o graph \"orderGraph\"  filter v._key==\"%v\" return {order:p.vertices[0],payment:p.vertices[1],orderItem:p.vertices[2]}", orderKey, orderItemKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var doc OrderPaymentAndOrderItem
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil

}

func GetOrderItemsAndPaymentByPaymentKey(paymentKey string, filter string) (OrderItemsAndPayment, error) {
	q := fmt.Sprintf("\nfor p in gPayment  filter p._key==\"%v\" let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" %v  return v)\nreturn {payment:p,orderItems:oi}", paymentKey, filter)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		return OrderItemsAndPayment{}, err
	}
	defer cursor.Close()

	var doc OrderItemsAndPayment
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return OrderItemsAndPayment{}, err
	}

	return doc, nil
}

func GetOrderAndPayments(orderKey string) (getOrderAndPayments, error) {
	q := fmt.Sprintf("\nfor o in gOrder \nfilter o._key==\"%v\"\nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\"  return v)\nreturn {order:o,payment:payment} ", orderKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		return getOrderAndPayments{}, err
	}
	defer cursor.Close()

	var doc getOrderAndPayments
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return getOrderAndPayments{}, err
	}

	return doc, nil

}

type getOrderAndPayments struct {
	Order   GOrderOut
	Payment []graphPayment.GPaymentOut
}