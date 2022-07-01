package graphOrder

import (
	"bamachoub-backend-go-v1/app/cart"
	"bamachoub-backend-go-v1/app/graphPayment"
	"bamachoub-backend-go-v1/app/transportation"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"time"
)

// InitOrder Initialize order for first time
// @Summary Initialize order for first time
// @Description Initialize order for first time , use when user enter the "Checkout - Shipping info"
// @Tags order
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} GOrder{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /order/init [post]
func InitOrder(userKey string) ([]edgeData, error) {
	groupedCart, err := cart.GetCartGroupByBuyingMethod(userKey)
	if err != nil {
		return nil, err
	}
	edArr := make([]edgeData, 0)

	var totalOrderAmount int64
	for _, gc := range *groupedCart {
		var totalAmount int64
		oiArr := make([]GOrderItem, 0)
		for _, out := range gc.Cart {
			oi := GOrderItem{
				PriceId:                out.PriceId,
				SupplierKey:            out.SupplierKey,
				ProductId:              out.ProductId,
				PricePerNumber:         out.PricePerNumber,
				Number:                 out.Number,
				Variant:                out.Variant,
				ProductTitle:           out.ProductTitle,
				ProductImageUrl:        out.ProductImageUrl,
				UserKey:                out.UserKey,
				UserAuthType:           out.UserAuthType,
				CommissionPercent:      out.CommissionPercent,
				CheckCommissionPercent: out.CheckCommissionPercent,
				IsWaitingForPayment:    true,
				IsApprovedBySupplier:   false,
				SupplierEmployeeId:     "",
				IsProcessing:           false,
				IsArrived:              false,
				IsCancelled:            false,
				CancelledById:          "",
				IsReferred:             false,
				ReferredReason:         "",
				CreatedAt:              time.Now().Unix(),
			}
			oiArr = append(oiArr, oi)
			totalAmount += int64(out.Number) * out.PricePerNumber
		}

		p := graphPayment.GPayment{
			Type:           gc.Type,
			TotalPrice:     totalAmount,
			RemainingPrice: totalAmount,
			FromWallet:     0,
			Status:         "wait",
		}

		edArr = append(edArr, edgeData{
			Payment:    p,
			OrderItems: oiArr,
		})

		totalOrderAmount += totalAmount
	}

	gOrderCol := database.GetCollection("gOrder")
	gOrderItemCol := database.GetCollection("gOrderItem")
	gPaymentCol := database.GetCollection("gPayment")

	edgeIds := make([]edgeIdData, 0)
	for _, data := range edArr {
		var ei edgeIdData
		pMeta, err := gPaymentCol.CreateDocument(context.Background(), data.Payment)
		if err != nil {
			return []edgeData{}, err
		}
		ei.PaymentIds = pMeta.ID.String()

		metaArr, errArr, err := gOrderItemCol.CreateDocuments(context.Background(), data.OrderItems)
		if err != nil {
			return []edgeData{}, errors.New(fmt.Sprintf("%v", errArr))
		}
		idArr := make([]string, 0)
		for _, meta := range metaArr {
			idArr = append(idArr, meta.ID.String())
		}
		ei.OrderItemIdes = idArr
		edgeIds = append(edgeIds, ei)
	}

	o := GOrder{
		SendingInfoKey:               "",
		UserKey:                      userKey,
		TransportationPrice:          0,
		IsTransportationPriceIsPayed: false,
		TransportationPriceWithPrice: false,
		UseWalletForTransportation:   false,
		TransportationPaymentId:      "",
		TotalAmount:                  -1,
		CreateAt:                     time.Now().Unix(),
	}
	oMeta, err := gOrderCol.CreateDocument(context.Background(), o)
	if err != nil {
		return []edgeData{}, err
	}

	edgeDocumentArr := make([]orderEdge, 0)

	for _, id := range edgeIds {
		for _, ide := range id.OrderItemIdes {
			edgeDocumentArr = append(edgeDocumentArr, orderEdge{
				From: id.PaymentIds,
				To:   ide,
			})
		}
		edgeDocumentArr = append(edgeDocumentArr, orderEdge{
			From: oMeta.ID.String(),
			To:   id.PaymentIds,
		})

	}
	edgeDocumentArr = append(edgeDocumentArr, orderEdge{
		From: fmt.Sprintf("users/%v", userKey),
		To:   oMeta.ID.String(),
	})
	fmt.Println(edgeDocumentArr)

	orderEdgeCol := database.GetCollection("gOrderEdge")
	_, errArr, err := orderEdgeCol.CreateDocuments(context.Background(), edgeDocumentArr)
	if err != nil {
		return []edgeData{}, errors.New(fmt.Sprintf("%v", errArr))
	}

	q := fmt.Sprintf("for i in cart filter  i.userKey==\"%v\" remove i in cart", userKey)
	database.ExecuteGetQuery(q)

	return edArr, nil
}

// getOrderByUserKey get order by jwt
// @Summary get order by jwt
// @Description get order by jwt
// @Tags order
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Param   tab      query    string     false        "all / wait-payment / processing / arrived / cancelled / referred"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} GOrder{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /order/user [get]
func getOrdersByUserKey(userKey string, filter string, offset string, limit string) (*[]GOrderResponseOut, error) {
	var q string
	if filter == "all" {
		q = fmt.Sprintf("for u in users filter u._key==\"%v\" \nlet order =(for v,e in 1..1 outbound u graph \"orderGraph\" return v)\nfor o in order \nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\" return v)\nlet result=(for p in payment let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" return v) return {payment:p,orderItems:oi})\n  sort o.createdAt desc limit %v,%v return {order:o,items:result}", userKey, offset, limit)
	} else if filter == "wait-payment" {
		q = fmt.Sprintf("for u in users filter u._key==\"%v\" \nlet order =(for v,e in 1..1 outbound u graph \"orderGraph\" return v)\nfor o in order \nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\" filter v.remainingPrice!=0 return v)\nlet result=(for p in payment let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" return v) return {payment:p,orderItems:oi})\n  sort o.createdAt desc limit %v,%v return {order:o,items:result}", userKey, offset, limit)
	} else if filter == "processing" {
		q = fmt.Sprintf("for u in users filter u._key==\"%v\" \nlet order =(for v,e in 1..1 outbound u graph \"orderGraph\" return v)\nfor o in order \nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\"  return v)\nlet result=(for p in payment let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" filter v.isProcessing==true return v) return {payment:p,orderItems:oi})\n  sort o.createdAt desc limit %v,%v return {order:o,items:result}", userKey, offset, limit)
	} else if filter == "arrived" {
		q = fmt.Sprintf("for u in users filter u._key==\"%v\" \nlet order =(for v,e in 1..1 outbound u graph \"orderGraph\" return v)\nfor o in order \nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\"  return v)\nlet result=(for p in payment let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" filter v.isArrived==true return v) return {payment:p,orderItems:oi})\n  sort o.createdAt desc limit %v,%v return {order:o,items:result}", userKey, offset, limit)
	} else if filter == "cancelled" {
		q = fmt.Sprintf("for u in users filter u._key==\"%v\" \nlet order =(for v,e in 1..1 outbound u graph \"orderGraph\" return v)\nfor o in order \nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\"  return v)\nlet result=(for p in payment let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" filter v.isCancelled==true return v) return {payment:p,orderItems:oi})\n  sort o.createdAt desc limit %v,%v return {order:o,items:result}", userKey, offset, limit)
	} else if filter == "referred" {
		q = fmt.Sprintf("for u in users filter u._key==\"%v\" \nlet order =(for v,e in 1..1 outbound u graph \"orderGraph\" return v)\nfor o in order \nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\"  return v)\nlet result=(for p in payment let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" filter v.isReferred==true return v) return {payment:p,orderItems:oi})\n  sort o.createdAt desc limit %v,%v return {order:o,items:result}", userKey, offset, limit)
	}
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

	//calc price
	for i, datum := range data {
		var totalPrice int64
		for _, i := range datum.Items {
			if filter == "wait-payment" {
				totalPrice += i.Payment.RemainingPrice
			} else {
				totalPrice += i.Payment.TotalPrice
			}

		}

		data[i].Order.TotalAmount = totalPrice

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
	statusMap[1] = "WaitingForPayment"
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
				} else if item.IsWaitingForPayment {
					statusArr = append(statusArr, 1)
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

	return &final, nil
}

// getAllOrdersForAdmin get order for admin
// @Summary get order for admin
// @Description get order for admin
// @Tags order
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Success 200 {object} GOrder{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /order/admin [get]
func getAllOrdersForAdmin(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	q := fmt.Sprintf("for u in users \nlet orderl=(for i in gOrder return i)\nlet order =(for v,e in 1..1 outbound u graph \"orderGraph\" return v)\nfor o in order \nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\" return v)\nlet result=(for p in payment let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" return v) return {payment:p,orderItems:oi})\n  sort o.createdAt desc limit %v,%v  return {user:u,order:o,items:result,length:length(orderl)}", offset, limit)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		//fmt.Println(q)
		return c.Status(500).JSON(err)
	}
	defer cursor.Close()
	var data []GOrderAdminRespOut
	for {
		var doc GOrderAdminRespOut
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}

	for i, datum := range data {
		var totalPrice int64
		for _, i := range datum.Items {

			totalPrice += i.Payment.TotalPrice

		}

		data[i].Order.TotalAmount = totalPrice

	}

	//remove empty orders
	final := make([]GOrderAdminRespOut, 0)
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
	statusMap[1] = "WaitingForPayment"
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
				} else if item.IsWaitingForPayment {
					statusArr = append(statusArr, 1)
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
	return c.JSON(final)
}

// getOrderByKey get order by order key
// @Summary get order by order key
// @Description get order by order key
// @Tags order
// @Accept json
// @Produce json
// @Param   key     path    int     true        "order key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} GOrder{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /order/{key} [get]
func getOrderByKey(c *fiber.Ctx) error {
	orderKey := c.Params("key")
	userKey := c.Locals("userKey").(string)

	q := fmt.Sprintf("for u in users filter u._key==\"%v\" \nlet order =(for v,e in 1..1 outbound u graph \"orderGraph\" return v)\nfor o in order filter o._key==\"%v\" \nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\" return v)\nlet result=(for p in payment let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" return v) return {payment:p,orderItems:oi})\n  return {order:o,items:result}", userKey, orderKey)

	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		//fmt.Println(q)
		return c.Status(500).JSON(err)
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
	statusMap[1] = "WaitingForPayment"
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
				} else if item.IsWaitingForPayment {
					statusArr = append(statusArr, 1)
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

	return c.JSON(final[0])

}

// addSendingInfoToOrder update order by sending info
// @Summary update order by sending info
// @Description update order by sending info , sending info must be created first ,returns order with transportation price
// @Tags order
// @Accept json
// @Produce json
// @Param data body sendingInfo true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} GOrder{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /order/sending-info [patch]
func addSendingInfoToOrder(c *fiber.Ctx) error {

	s := new(sendingInfo)
	if err := utils.ParseBodyAndValidate(c, s); err != nil {
		return c.JSON(err)
	}
	userKey := c.Locals("userKey").(string)

	orderCol := database.GetCollection("gOrder")
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
	var o GOrder
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
