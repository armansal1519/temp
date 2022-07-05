package paymentAndWallet

import (
	"bamachoub-backend-go-v1/app/graphOrder"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"time"
)

func supplierConfirmation(oi []graphOrder.GOrderItemOut, orderKey string) error {

	infoArr := make([]supplierInfoForConfirmation, 0)
	for _, cart := range oi {
		temp := supplierInfoForConfirmation{
			SupplierKey:  cart.SupplierKey,
			OrderKey:     orderKey,
			OrderItemKey: cart.Key,
		}
		infoArr = append(infoArr, temp)
	}
	scCol := database.GetCollection("supplierConfirmation")
	metaArr, errArr, err := scCol.CreateDocuments(context.Background(), infoArr)
	//fmt.Println(111111111,meta)
	if err != nil {
		return fmt.Errorf("%v", errArr)
	}

	for _, meta := range metaArr {
		time.AfterFunc(24*time.Hour, func() {
			err := rejectOrder(meta.Key, "", "system")
			if err != nil {
				return
			}
		})
	}
	return nil

}

// GetOrderConfirmationBySupplierKey get orders for supplier confirmation
// @Summary get orders for supplier confirmation
// @Description get orders for supplier confirmation
// @Tags supplier confirmation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []getSupplierConfirmationResponse{}
// @Failure 404 {object} string{}
// @Router /suppliers-confirmation [get]
func GetOrderConfirmationBySupplierKey(c *fiber.Ctx) error {
	supplierKey := c.Locals("supplierId").(string)

	query := fmt.Sprintf("for i in gOrderItem filter i.supplierKey==\"%v\"  sort i.isApprovedBySupplier  \nlet data=(for v,e,p  in 0..3 inbound i graph \"orderGraph\"  return v)\nreturn data ", supplierKey)
	return c.JSON(database.ExecuteGetQuery(query))

}

// approveOrder approve order
// @Summary approve order
// @Description approve order
// @Tags supplier confirmation
// @Accept json
// @Produce json
// @Param infoKey path string true "infoKey"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /suppliers-confirmation/approve/{infoKey} [post]
func approveOrder(orderItemKey string, supplierKey string, supplierEmployeeId string) error {
	//var sc supplierInfoForConfirmationOut
	//col := database.GetCollection("supplierConfirmation")
	//_, err := col.ReadDocument(context.Background(), orderItem, &sc)
	//if err != nil {
	//	return err
	//}

	var o graphOrder.GOrderItem
	gOrderItemCol := database.GetCollection("gOrderItem")
	_, err := gOrderItemCol.ReadDocument(context.Background(), orderItemKey, &o)

	if o.SupplierKey != supplierKey {
		return errors.New("different suppliers")
	}

	u := updateOrder{
		IsApprovedBySupplier: true,
		SupplierEmployeeId:   supplierEmployeeId,
	}
	_, err = gOrderItemCol.UpdateDocument(context.Background(), orderItemKey, u)
	if err != nil {
		return err
	}
	return nil

}

//func approveOrder(infoKey string, supplierKey string) error {
//	var sc supplierInfoForConfirmationOut
//	col := database.GetCollection("supplierConfirmation")
//	_, err := col.ReadDocument(context.Background(), infoKey, &sc)
//	if err != nil {
//		return err
//	}
//
//	var c cart.CartOut
//
//	cartCol := database.GetCollection("cart")
//	_, err = cartCol.ReadDocument(context.Background(), sc.OrderItemKey, &c)
//	if err != nil {
//		return err
//	}
//
//	if sc.SupplierKey != supplierKey {
//		return fmt.Errorf("not your data")
//	}
//
//	var o graphOrder.GOrder
//	orderCol := database.GetCollection("gOrder")
//	_, err = orderCol.ReadDocument(context.Background(), sc.OrderKey, &o)
//
//	var oi orders.OrderItem
//	for _, item := range o.OrderItems {
//		for _, out := range item.Cart {
//			if out.Key == sc.CartKey {
//				oi = item
//			}
//		}
//	}
//	ao := ApprovedOrder{
//		UserKey:           o.UserKey,
//		ProductId:         c.ProductId,
//		ProductTitle:      c.ProductTitle,
//		ProductImageUrl:   c.ProductImageUrl,
//		SupplierKey:       c.SupplierKey,
//		PaymentKey:        oi.PaymentKey,
//		CommissionPercent: c.CommissionPercent,
//		TxType:            oi.Type,
//		Price:             c.PricePerNumber * int64(c.Number),
//		Number:            c.Number,
//		CreatedAt:         time.Now().Unix(),
//		SendInfoKey:       o.SendingInfoKey,
//		Status:            "wait-send",
//	}
//
//	aoCol := database.GetCollection("approvedOrder")
//	_, err = aoCol.CreateDocument(context.Background(), ao)
//	if err != nil {
//		return err
//	}
//	_, err = col.RemoveDocument(context.Background(), sc.Key)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// rejectOrder reject order
// @Summary reject order
// @Description reject order
// @Tags supplier confirmation
// @Accept json
// @Produce json
// @Param infoKey path string true "infoKey"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /suppliers-confirmation/reject/{infoKey} [post]
func rejectOrder(infoKey string, supplierKey string, callBy string) error {
	col := database.GetCollection("supplierConfirmation")
	flag, err := col.DocumentExists(context.Background(), infoKey)
	if err != nil {
		return err
	}
	if !flag {
		if callBy == "system" {
			return nil
		} else {
			return fmt.Errorf("info not Found")
		}
	}
	var sc supplierInfoForConfirmationOut
	_, err = col.ReadDocument(context.Background(), infoKey, &sc)
	if err != nil {
		return err
	}

	//var c cart.CartOut
	//
	//cartCol := database.GetCollection("cart")
	//_, err = cartCol.ReadDocument(context.Background(), sc.OrderItemKey, &c)
	//if err != nil {
	//	return err
	//}

	if callBy != "system" {
		if sc.SupplierKey != supplierKey {
			return fmt.Errorf("not your data")
		}
	}

	//var o orders.Order
	//orderCol := database.GetCollection("order")
	//_, err = orderCol.ReadDocument(context.Background(), sc.OrderKey, &o)
	//
	//var oi orders.OrderItem
	//for _, item := range o.OrderItems {
	//	for _, out := range item.Cart {
	//		if out.Key == sc.CartKey {
	//			oi = item
	//		}
	//	}
	//}

	data, err := graphOrder.GetOrderPaymentAndOrderItem(sc.OrderKey, sc.OrderItemKey)

	ro := rejectionPoolItem{
		UserKey:         data.Order.UserKey,
		ProductId:       data.OrderItem.ProductId,
		ProductTitle:    data.OrderItem.ProductTitle,
		ProductImageUrl: data.OrderItem.ProductImageUrl,
		RejectBy:        data.OrderItem.SupplierKey,
		PaymentKey:      data.Payment.Key,
		TxType:          data.Payment.Type,
		Price:           data.OrderItem.PricePerNumber * int64(data.OrderItem.Number),
		Number:          data.OrderItem.Number,
		CreatedAt:       time.Now().Unix(),
		SendInfoKey:     data.Order.SendingInfoKey,
		Status:          "wait-accept",
	}

	rpCol := database.GetCollection("rejectionPool")
	_, err = rpCol.CreateDocument(context.Background(), ro)
	if err != nil {
		return err
	}
	_, err = col.RemoveDocument(context.Background(), sc.Key)
	if err != nil {
		return err
	}
	return nil

}

// rejectOrder new reject order
// @Summary new reject order
// @Description new reject order
// @Tags supplier confirmation
// @Accept json
// @Produce json
// @Param orderItemKey path string true "orderItemKey"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} graphOrder.GOrderItemOut
// @Failure 404 {object} string{}
// @Router /suppliers-confirmation/g-reject/{orderItemKey} [post]
func graphRejectOrder(orderItemKey string, supplierKey string, callBySystem bool) (*graphOrder.GOrderItemOut, error) {
	goiCol := database.GetCollection("gOrderItem")
	var orderItemOut graphOrder.GOrderItemOut
	ctx := driver.WithReturnNew(context.Background(), &orderItemOut)
	if callBySystem {
		u := updateOrderItemFromRejection{
			IsRejected:         true,
			IsRejectedBySystem: true,
			RejectedById:       "",
			RejectedAt:         time.Now().Unix(),
		}
		_, err := goiCol.UpdateDocument(ctx, orderItemKey, u)
		if err != nil {
			return nil, err
		}
	} else {
		u := updateOrderItemFromRejection{
			IsRejected:         true,
			IsRejectedBySystem: false,
			RejectedById:       "suppliers/" + supplierKey,
			RejectedAt:         time.Now().Unix(),
		}
		_, err := goiCol.UpdateDocument(ctx, orderItemKey, u)
		if err != nil {
			return nil, err
		}
	}

	return &orderItemOut, nil
}

// rejectOrder get reject orders
// @Summary get reject orders
// @Description get reject orders
// @Tags supplier confirmation
// @Accept json
// @Produce json
// @Param offset query string true "offset"
// @Param limit query string true "limit"
// @Param rejected-by query string false "key of supplier you want to filter "
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []graphOrder.GOrderItemOut
// @Failure 404 {object} string{}
// @Router /suppliers-confirmation/g-reject [get]
func getGraphRejectedOrder(c *fiber.Ctx) error {
	isAdmin := c.Locals("isAdmin").(bool)

	offset := c.Query("offset")
	limit := c.Query("limit")
	rejectedBy := c.Query("rejected-by")

	if offset == "" || limit == "" {
		return c.Status(400).SendString("offset or limit is empty")
	}
	if isAdmin {

		f := ""
		if rejectedBy != "" {
			f = fmt.Sprintf(" filter i.supplierKey == \"%v\" ", rejectedBy)
		}
		q := fmt.Sprintf("for i in gOrderItem %v sort i.rejectedAt limit %v,%v return i", f, offset, limit)
		return c.JSON(database.ExecuteGetQuery(q))
	}
	supplierKey := c.Locals("supplierId").(string)
	q := fmt.Sprintf("for i in gOrderItem filter i.supplierKey != \"%v\" sort i.rejectedAt limit %v,%v return i", supplierKey, offset, limit)
	return c.JSON(database.ExecuteGetQuery(q))

}

// rejectOrder accept a  rejected order
// @Summary accept a  rejected order
// @Description accept a  rejected order
// @Tags supplier confirmation
// @Accept json
// @Produce json
// @Param orderItemKey path string true "orderItemKey"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} graphOrder.GOrderItemOut
// @Failure 404 {object} string{}
// @Failure 409 {object} string{}
// @Failure 500 {object} string{}
// @Router /suppliers-confirmation/g-accept/{orderItemKey} [post]
func acceptARejectedOrder(c *fiber.Ctx) error {
	orderItemKey := c.Params("orderItemKey")
	supplierKey := c.Locals("supplierId").(string)
	orderItemCol := database.GetCollection("gOrderItem")
	var oi graphOrder.GOrderItemOut

	_, err := orderItemCol.ReadDocument(context.Background(), orderItemKey, &oi)
	if err != nil {
		if driver.IsNotFound(err) {
			return c.Status(404).JSON(err)
		}
		return c.Status(500).JSON(err)
	}

	if !oi.IsRejected {
		return c.Status(409).SendString("you can not accept a nor rejected order item")
	}
	u := updateOrderItemFromAccept{
		IsAcceptedAfterRejection: true,
		AcceptedById:             "supplier/" + supplierKey,
		AcceptedAt:               time.Now().Unix(),
	}
	var newOI graphOrder.GOrderItemOut
	ctx := driver.WithReturnNew(context.Background(), &newOI)
	_, err = orderItemCol.UpdateDocument(ctx, orderItemKey, u)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.JSON(newOI)

}

func getApprovedOrderForUser(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	q := fmt.Sprintf("for i in approvedOrder filter i.userKey==\"%v\" sort i.createdAt return i", userKey)
	return c.JSON(database.ExecuteGetQuery(q))
}
