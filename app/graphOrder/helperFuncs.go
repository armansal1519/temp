package graphOrder

import (
	"bamachoub-backend-go-v1/app/graphPayment"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"time"
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

func GetOrderAndPayments(orderKey string) (GetOrderAndPaymentsDto, error) {
	q := fmt.Sprintf("\nfor o in gOrder \nfilter o._key==\"%v\"\nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\"  return v)\nreturn {order:o,payment:payment} ", orderKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		return GetOrderAndPaymentsDto{}, err
	}
	defer cursor.Close()

	var doc GetOrderAndPaymentsDto
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return GetOrderAndPaymentsDto{}, err
	}

	return doc, nil

}

func GetOrderAndPaymentsByUserKey(orderKey string) (GetOrderAndPaymentsDto, error) {
	q := fmt.Sprintf("\nfor o in gOrder \nfilter o.userKey==\"%v\"\nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\"  return v)\nreturn {order:o,payment:payment} ", orderKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		return GetOrderAndPaymentsDto{}, err
	}
	defer cursor.Close()

	var doc GetOrderAndPaymentsDto
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return GetOrderAndPaymentsDto{}, err
	}

	return doc, nil

}

func getOrderByKeyHelper(orderKey, userKey string) (*GOrderResponseOut, error) {
	q := fmt.Sprintf("for u in users filter u._key==\"%v\" \nlet order =(for v,e in 1..1 outbound u graph \"orderGraph\" return v)\nfor o in order filter o._key==\"%v\" \nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\" return v)\nlet result=(for p in payment let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" return v) return {payment:p,orderItems:oi})\n  return {order:o,items:result}", userKey, orderKey)

	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		//fmt.Println(q)
		return nil, err
	}
	defer cursor.Close()
	var data []GOrderResponseOut
	for {
		var doc GOrderResponseOut
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}

	//remove empty orders
	final := make([]GOrderResponseOut, 0)
	for _, datum := range data {
		flag := false
		for _, item := range datum.Items {
			if len(item.OrderItems) > 0 {
				flag = true
			}
		}
		if flag {
			final = append(final, datum)
		}
	}

	//calc status
	statusMap := make(map[int]string)
	statusMap[0] = "WaitingForPayment"
	statusMap[1] = "WaitingForSupplierToApprove"
	statusMap[2] = "ApprovedBySupplier"
	statusMap[3] = "Processing"
	statusMap[4] = "Arrived"
	statusMap[5] = "Cancelled"
	statusMap[6] = "Referred"
	for ii, f := range final {
		statusArr := make([]int, 0)
		for _, i := range f.Items {
			for _, item := range i.OrderItems {
				if item.IsReferred {
					statusArr = append(statusArr, 6)
				} else if item.IsCancelled {
					statusArr = append(statusArr, 5)
				} else if item.IsArrived {
					statusArr = append(statusArr, 4)
				} else if item.IsProcessing {
					statusArr = append(statusArr, 3)
				} else if item.IsApprovedBySupplier {
					statusArr = append(statusArr, 2)
				} else if !item.IsWaitingForPayment && !item.IsApprovedBySupplier {
					statusArr = append(statusArr, 1)
				} else if item.IsWaitingForPayment {
					statusArr = append(statusArr, 0)
				}
			}
		}

		statusScore := 6
		for _, i := range statusArr {
			if i < statusScore {
				statusScore = i
			}
		}
		final[ii].Order.Status = statusMap[statusScore]

	}

	//calc price
	for i, datum := range final {
		var totalPrice int64
		for _, j := range datum.Items {
			if datum.Order.Status == "wait-payment" {
				totalPrice += j.Payment.RemainingPrice
			} else {
				totalPrice += j.Payment.TotalPrice
			}

		}

		data[i].Order.TotalAmount = totalPrice

	}

	//add reserved
	for i, out := range final {
		for _, out2 := range out.Items {
			if out2.Payment.IsRejected {
				if time.Now().Unix() < out2.Payment.RejectionTime {
					final[i].Reserved = reservedInfo{
						IsReserved: true,
						TimeToEnd:  out2.Payment.RejectionTime - time.Now().Unix(),
					}
				}
			}
		}
	}

	return &final[0], nil
}

type GetOrderAndPaymentsDto struct {
	Order   GOrderOut
	Payment []graphPayment.GPaymentOut
}
