package paymentAndWallet

import (
	"bamachoub-backend-go-v1/app/orders"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/payment"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"time"
)

func createPayment(userKey string, orderKey string, amount int64, txType string, ImageUrl string, IncludeTransportation bool) (*paymentOut, error) {
	p := paymentHistory{
		PayerKey:              userKey,
		OrderKey:              orderKey,
		TxType:                txType,
		Amount:                amount,
		Status:                "not",
		CardHolder:            "",
		ShaparakRefId:         "",
		TransId:               "",
		ImageUrl:              ImageUrl,
		CheckNumber:           "",
		IncludeTransportation: IncludeTransportation,
		CreatedAt:             time.Now().Unix(),
	}
	var pOut paymentOut
	ctx := driver.WithReturnNew(context.Background(), &pOut)
	paymentCol := database.GetCollection("payment")
	meta, err := paymentCol.CreateDocument(ctx, p)
	if err != nil {
		return nil, err
	}

	u := UpdateOrderWithPaymentKey{PaymentKey: meta.Key}
	orderCol := database.GetCollection("order")
	_, err = orderCol.UpdateDocument(context.Background(), orderKey, u)
	if err != nil {
		return nil, err
	}
	return &pOut, nil
}

// getPaymentUrl get payment by url
// @Summary get payment by url
// @Description get payment by url
// @Tags payment
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /payment/by-url [post]
func getPaymentUrl(c *fiber.Ctx) error {
	cpp := new(createPaymentByPortal)
	if err := utils.ParseBodyAndValidate(c, cpp); err != nil {
		return c.JSON(err)
	}
	var o orders.Order
	orderCol := database.GetCollection("order")
	_, err := orderCol.ReadDocument(context.Background(), cpp.OrderKey, &o)
	if err != nil {
		return c.JSON(err)
	}

	var amount int64
	for _, item := range o.OrderItems {
		if item.Type == "price" {
			amount = item.TotalPrice
		}
	}
	if cpp.IncludeTransportation {
		amount += o.TransportationPrice
	}
	userKey := c.Locals("userKey").(string)
	p, err := createPayment(userKey, cpp.OrderKey, amount, "price-portal", "", cpp.IncludeTransportation)
	if err != nil {
		return c.JSON(err)
	}
	transId, err := payment.GetPaymentUrl(fmt.Sprintf("%v", cpp.Amount), p.Key, fmt.Sprintf("https://localhost:3000/payment-varification/%v", p.Key))
	if err != nil {
		return c.JSON(err)
	}
	ut := updateTransId{TransId: transId}
	paymentCol := database.GetCollection("payment")
	_, err = paymentCol.UpdateDocument(context.Background(), p.Key, ut)
	if err != nil {
		return c.JSON(err)
	}
	err = updateOrderWithPaymentKey(o, p.Key, "price", cpp.IncludeTransportation)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(fiber.Map{
		"paymentUrl": fmt.Sprintf("https://nextpay.org/nx/gateway/payment/%v", transId),
	})

}

// verifyPaymentUrl verity payment by url
// @Summary verity payment by url
// @Description verity payment by url
// @Tags payment
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /payment/verify-url/{key} [post]
func verifyPaymentUrl(c *fiber.Ctx) error {
	PaymentKey := c.Params("key")
	userKey := c.Locals("userKey").(string)
	var p paymentOut
	paymentCol := database.GetCollection("payment")
	_, err := paymentCol.ReadDocument(context.Background(), PaymentKey, &p)
	if err != nil {
		return c.JSON(err)
	}
	if userKey != p.PayerKey {
		return c.Status(403).JSON("this is not your order")
	}
	v, err := payment.Verify(p.Amount, p.TransId)
	if err != nil {
		return c.JSON(err)
	}
	up := updatePaymentHistory{
		CardHolder:    v.CardHolder,
		ShaparakRefId: v.ShaparakRefId,
		Status:        "valid",
	}

	_, err = paymentCol.UpdateDocument(context.Background(), PaymentKey, up)

	orderCol := database.GetCollection("order")
	var o orders.Order
	_, err = orderCol.ReadDocument(context.Background(), p.OrderKey, &o)
	var oi orders.OrderItem
	for i, item := range o.OrderItems {
		if item.Type == "price" {
			o.OrderItems[i].Status = "valid"
			o.OrderItems[i].RemainingPrice = 0
			oi = item
			break
		}
	}
	if p.IncludeTransportation {
		o.IsTransportationPriceIsPayed = true
		o.TransportationPriceWithPrice = true
	}

	meta, err := orderCol.UpdateDocument(context.Background(), o.Key, o)
	if err != nil {
		return c.JSON(err)
	}
	err = supplierConfirmation(oi, o.Key)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// createPaymentByImage get payment by image
// @Summary get payment by image
// @Description get payment by image
// @Tags payment
// @Accept json
// @Produce json
// @Param PaymentByImage body PaymentByImage true "data"
// @Success 200 {object} paymentOut{}
// @Failure 404 {object} string{}
// @Router /payment/by-image [post]
func createPaymentByImage(c *fiber.Ctx) error {
	pbi := new(PaymentByImage)
	if err := utils.ParseBodyAndValidate(c, pbi); err != nil {
		return c.JSON(err)
	}
	var o orders.Order
	orderCol := database.GetCollection("order")
	_, err := orderCol.ReadDocument(context.Background(), pbi.OrderKey, &o)
	if err != nil {
		return c.JSON(err)
	}

	var amount int64
	itemIndex := 0
	for i, item := range o.OrderItems {
		if item.Type == "price" {
			amount = item.TotalPrice
			itemIndex = i
		}
	}
	//if there is already a payment Key and overwritePaymentKey was true
	if !pbi.OverwritePaymentKey {
		if o.OrderItems[itemIndex].PaymentKey != "" {
			return c.Status(409).SendString("this item already have paymentKey")
		}
	}

	if o.OrderItems[itemIndex].Status == "done" {
		return c.Status(409).SendString("this already payed")
	}

	if pbi.IncludeTransportation {
		amount += o.TransportationPrice
	}
	if pbi.Type != "s-p" && pbi.Type != "ctoc" && pbi.Type != "place" {
		return c.Status(409).SendString("only s-p or ctoc or place is allowed as type")
	}
	userKey := c.Locals("userKey").(string)
	p, err := createPayment(userKey, pbi.OrderKey, amount, fmt.Sprintf("price-%v", pbi.Type), pbi.ImageUrl, pbi.IncludeTransportation)
	if err != nil {
		return c.JSON(err)
	}
	err = updateOrderWithPaymentKey(o, p.Key, "price", pbi.IncludeTransportation)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(p)

}

// verifyPaymentImage verity payment
// @Summary verity payment
// @Description verity payment
// @Tags payment
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /payment/verify-image/{key} [post]
func verifyPaymentImage(c *fiber.Ctx) error {
	PaymentKey := c.Params("key")
	var p paymentOut
	paymentCol := database.GetCollection("payment")
	_, err := paymentCol.ReadDocument(context.Background(), PaymentKey, &p)
	if err != nil {
		return c.JSON(err)
	}

	orderCol := database.GetCollection("order")
	var o orders.Order
	_, err = orderCol.ReadDocument(context.Background(), p.OrderKey, &o)

	var oi orders.OrderItem

	for i, item := range o.OrderItems {
		if item.Type == "price" {
			o.OrderItems[i].Status = "valid"
			o.OrderItems[i].RemainingPrice = 0

			oi = item
			break
		}
	}
	if p.IncludeTransportation {
		o.IsTransportationPriceIsPayed = true
		o.TransportationPriceWithPrice = true
	}

	meta, err := orderCol.UpdateDocument(context.Background(), o.Key, o)
	if err != nil {
		return c.JSON(err)
	}
	err = supplierConfirmation(oi, o.Key)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// createPaymentByImage create check by image
// @Summary create check by image
// @Description create check by image
// @Tags payment
// @Accept json
// @Produce json
// @Param checkByImage body checkByImage true "data"
// @Success 200 {object} paymentOut{}
// @Failure 404 {object} string{}
// @Router /payment/by-check [post]
func createCheckPayment(c *fiber.Ctx) error {
	pbi := new(checkByImage)
	if err := utils.ParseBodyAndValidate(c, pbi); err != nil {
		return c.JSON(err)
	}
	var o orders.Order
	orderCol := database.GetCollection("order")
	_, err := orderCol.ReadDocument(context.Background(), pbi.OrderKey, &o)
	if err != nil {
		return c.JSON(err)
	}
	if pbi.Type != "one" && pbi.Type != "two" && pbi.Type != "three" {
		return c.Status(409).SendString("type can be only one or two or three")
	}

	var amount int64
	for _, item := range o.OrderItems {
		if item.Type == pbi.Type {
			amount = item.TotalPrice
		}
	}
	userKey := c.Locals("userKey").(string)
	p, err := createPayment(userKey, pbi.OrderKey, amount, fmt.Sprintf("check-%v", pbi.Type), pbi.ImageUrl, false)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(p)
}

// verifyCheckImage verity check
// @Summary verity check
// @Description verity check
// @Tags payment
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /payment/verify-check/{key} [post]
func verifyCheckImage(c *fiber.Ctx) error {
	PaymentKey := c.Params("key")
	var p paymentOut
	paymentCol := database.GetCollection("payment")
	_, err := paymentCol.ReadDocument(context.Background(), PaymentKey, &p)
	if err != nil {
		return c.JSON(err)
	}

	orderCol := database.GetCollection("order")
	var o orders.Order
	_, err = orderCol.ReadDocument(context.Background(), p.OrderKey, &o)

	var oi orders.OrderItem

	for i, item := range o.OrderItems {
		if item.Type == "price" {
			o.OrderItems[i].Status = "valid"
			o.OrderItems[i].RemainingPrice = 0

			oi = item
			break
		}
	}

	meta, err := orderCol.UpdateDocument(context.Background(), o.Key, o)
	if err != nil {
		return c.JSON(err)
	}
	err = supplierConfirmation(oi, o.Key)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

func getPaymentByKey(paymentKey string) (*paymentOut, error) {
	var p paymentOut
	paymentCol := database.GetCollection("payment")
	_, err := paymentCol.ReadDocument(context.Background(), paymentKey, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil

}

// filterPayment fet filtered payment
// @Summary fet filtered payment
// @Description fet filtered payment if filter is empty return all
// @Tags payment
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Param filter body filter true "data"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /payment/filter [post]
func filterPayment(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	f := new(filter)
	if err := utils.ParseBodyAndValidate(c, f); err != nil {
		return c.JSON(err)
	}
	filterString := ""
	if f.TxType != "" {
		filterString += fmt.Sprintf(" filter  i.txType== %v \n", f.TxType)
	}
	if f.Status != "" {
		filterString += fmt.Sprintf(" filter  i.status== %v \n", f.Status)
	}
	if f.ShaparakRefId != "" {
		filterString += fmt.Sprintf(" filter  i.shaparakRefId== %v \n", f.ShaparakRefId)
	}
	if f.CheckNumber != "" {
		filterString += fmt.Sprintf(" filter  i.checkNumber== %v \n", f.CheckNumber)
	}
	if f.OrderKey != "" {
		filterString += fmt.Sprintf(" filter  i.orderKey== %v \n", f.OrderKey)
	}
	if f.PayerKey != "" {
		filterString += fmt.Sprintf("  filter i.payerKey== %v \n", f.PayerKey)
	}

	query := fmt.Sprintf("for i in payment  %v  sort i.createdAt limit %v,%v return i", filterString, offset, limit)
	fmt.Println(query)

	return c.JSON(database.ExecuteGetQuery(query))
}
