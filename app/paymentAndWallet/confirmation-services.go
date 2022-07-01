package paymentAndWallet

import (
	"bamachoub-backend-go-v1/app/graphOrder"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"errors"
	"fmt"
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

func getApprovedOrderForUser(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	q := fmt.Sprintf("for i in approvedOrder filter i.userKey==\"%v\" sort i.createdAt return i", userKey)
	return c.JSON(database.ExecuteGetQuery(q))
}
