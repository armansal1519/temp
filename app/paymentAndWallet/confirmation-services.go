package paymentAndWallet

import (
	"bamachoub-backend-go-v1/app/cart"
	"bamachoub-backend-go-v1/app/orders"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"time"
)

func supplierConfirmation(oi orders.OrderItem, orderKey string) error {

	infoArr := make([]supplierInfoForConfirmation, 0)
	for _, cart := range oi.Cart {
		temp := supplierInfoForConfirmation{
			SupplierKey: cart.SupplierKey,
			OrderKey:    orderKey,
			CartKey:     cart.Key,
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
			rejectOrder(meta.Key, "", "system")
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
func GetOrderConfirmationBySupplierKey(supplierKey string) (*[]getSupplierConfirmationResponse, error) {
	query := fmt.Sprintf("for i in supplierConfirmation filter i.supplierKey==\"%v\" for j in cart filter j._key==i.cartKey return {cart:j,info:i}", supplierKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error while running query:%v", query)
	}
	defer cursor.Close()
	var data []getSupplierConfirmationResponse
	for {
		var doc getSupplierConfirmationResponse
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return &data, nil
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
func approveOrder(infoKey string, supplierKey string) error {
	var sc supplierInfoForConfirmationOut
	col := database.GetCollection("supplierConfirmation")
	_, err := col.ReadDocument(context.Background(), infoKey, &sc)
	if err != nil {
		return err
	}

	var c cart.CartOut

	cartCol := database.GetCollection("cart")
	_, err = cartCol.ReadDocument(context.Background(), sc.CartKey, &c)
	if err != nil {
		return err
	}

	if sc.SupplierKey != supplierKey {
		return fmt.Errorf("not your data")
	}

	var o orders.Order
	orderCol := database.GetCollection("order")
	_, err = orderCol.ReadDocument(context.Background(), sc.OrderKey, &o)

	var oi orders.OrderItem
	for _, item := range o.OrderItems {
		for _, out := range item.Cart {
			if out.Key == sc.CartKey {
				oi = item
			}
		}
	}
	ao := ApprovedOrder{
		UserKey:           o.UserKey,
		ProductId:         c.ProductId,
		ProductTitle:      c.ProductTitle,
		ProductImageUrl:   c.ProductImageUrl,
		SupplierKey:       c.SupplierKey,
		PaymentKey:        oi.PaymentKey,
		CommissionPercent: c.CommissionPercent,
		TxType:            oi.Type,
		Price:             c.PricePerNumber * int64(c.Number),
		Number:            c.Number,
		CreatedAt:         time.Now().Unix(),
		SendInfoKey:       o.SendingInfoKey,
		Status:            "wait-send",
	}

	aoCol := database.GetCollection("ApprovedOrder")
	_, err = aoCol.CreateDocument(context.Background(), ao)
	if err != nil {
		return err
	}
	_, err = col.RemoveDocument(context.Background(), sc.Key)
	if err != nil {
		return err
	}
	return nil
}

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

	var c cart.CartOut

	cartCol := database.GetCollection("cart")
	_, err = cartCol.ReadDocument(context.Background(), sc.CartKey, &c)
	if err != nil {
		return err
	}

	if callBy != "system" {
		if sc.SupplierKey != supplierKey {
			return fmt.Errorf("not your data")
		}
	}

	var o orders.Order
	orderCol := database.GetCollection("order")
	_, err = orderCol.ReadDocument(context.Background(), sc.OrderKey, &o)

	var oi orders.OrderItem
	for _, item := range o.OrderItems {
		for _, out := range item.Cart {
			if out.Key == sc.CartKey {
				oi = item
			}
		}
	}

	ro := rejectionPoolItem{
		UserKey:         o.UserKey,
		ProductId:       c.ProductId,
		ProductTitle:    c.ProductTitle,
		ProductImageUrl: c.ProductImageUrl,
		RejectBy:        c.SupplierKey,
		PaymentKey:      oi.PaymentKey,
		TxType:          oi.Type,
		Price:           c.PricePerNumber * int64(c.Number),
		Number:          c.Number,
		CreatedAt:       time.Now().Unix(),
		SendInfoKey:     o.SendingInfoKey,
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
