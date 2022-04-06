package graphPayment

import (
	"bamachoub-backend-go-v1/app/discountCode"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func getPaymentByKey(paymentKey string) (*GPaymentOut, error) {
	paymentCol := database.GetCollection("gPayment")
	var p GPaymentOut
	_, err := paymentCol.ReadDocument(context.Background(), paymentKey, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// addDiscountToPayment add discount to payment
// @Summary add discount to payment
// @Description add discount to payment
// @Tags graph payment
// @Accept json
// @Produce json
// @Param   key      path   string     true  "discount key"
// @Param   paymentkey      path   string     true  " paymentkey"
// @Param   use-less-discount      query   bool     true  " if amount of discount is  more than payment amount this overwrites error"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} GPaymentOut{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /gpayment/add-discount/{key}/{paymentkey} [post]
func addDiscountToPayment(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	discountKey := c.Params("key")
	paymentKey := c.Params("paymentkey")
	useLessDiscount := c.Query("use-less-discount")
	q := fmt.Sprintf("for d in discountEdge filter d._from==\"users/%v\" and d._to==\"discount/%v\" and d.isUsed==false\nfor i in discount filter i._id==d._to\nreturn i", userKey, discountKey)
	db := database.GetDB()
	cursor, err := db.Query(context.Background(), q, nil)
	log.Println(q)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("error while running query:%v", q))
	}
	defer cursor.Close()

	var doc discountCode.DiscountOut
	_, err = cursor.ReadDocument(context.Background(), &doc)
	if errors.Is(err, driver.NoMoreDocumentsError{}) {
		return c.JSON(nil)
	}

	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(err)
	}

	payment, err := getPaymentByKey(paymentKey)
	if payment.Type != doc.Type {
		return c.Status(409).JSON("payment type and discount type dont match")

	}
	if useLessDiscount != "true" {
		if payment.RemainingPrice < doc.Amount {
			return c.Status(409).JSON("discount is less than payment")
		}
	}

	if payment.DiscountKey != "" {
		return c.Status(400).JSON("this payment already have discount")
	}

	if doc.EndAt < time.Now().Unix() {
		return c.Status(409).JSON("discount is expired")
	}

	var newPayment GPaymentOut
	ctx := driver.WithReturnNew(context.Background(), &newPayment)
	paymentCol := database.GetCollection("gPayment")
	u := updatePaymentWithDiscount{
		DiscountKey:    discountKey,
		DiscountAmount: doc.Amount,
		RemainingPrice: payment.RemainingPrice - doc.Amount,
	}
	_, err = paymentCol.UpdateDocument(ctx, paymentKey, u)

	return c.JSON(newPayment)

}
