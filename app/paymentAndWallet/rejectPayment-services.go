package paymentAndWallet

import (
	"bamachoub-backend-go-v1/app/graphOrder"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
	"time"
)

//
//func rejectPaymentImage(c *fiber.Ctx) error {
//	PaymentKey := c.Params("key")
//	rpo, err := getReservedProductByPaymentKey(PaymentKey)
//	noDoc := errors.Is(err, driver.NoMoreDocumentsError{})
//	if err != nil {
//		if !noDoc {
//			return c.JSON(err)
//		}
//	}
//	if noDoc {
//		var p paymentOut
//		paymentCol := database.GetCollection("payment")
//		_, err := paymentCol.ReadDocument(context.Background(), PaymentKey, &p)
//		if err != nil {
//			return c.JSON(err)
//		}
//		orderCol := database.GetCollection("order")
//		var o orders.Order
//		_, err = orderCol.ReadDocument(context.Background(), p.OrderKey, &o)
//
//		var oi orders.OrderItem
//
//		for i, item := range o.OrderItems {
//			if item.Type == "price" {
//				o.OrderItems[i].Status = "rejected"
//
//				oi = item
//				break
//			}
//		}
//		meta, err := orderCol.UpdateDocument(context.Background(), o.Key, o)
//		if err != nil {
//			return c.JSON(err)
//		}
//		err = supplierConfirmation(oi, o.Key)
//		if err != nil {
//			return c.JSON(err)
//		}
//		err = updateNumberInPriceAndMoveToReserved(&oi.Cart, PaymentKey)
//
//		var rt time.Duration
//		if rpo.TxType == "price" {
//			rt = 2 * time.Hour
//		} else {
//			rt = 24 * time.Hour
//		}
//		time.AfterFunc(rt, func() {
//			deleteReserveProductAddNUmberToPrice(&oi.Cart, rpo.Key, 0)
//		})
//		return c.JSON(meta)
//
//	} else {
//
//		if rpo.FailedCount > 2 {
//			return c.Status(409).JSON("call support more than three try")
//		}
//
//		var rt time.Duration
//		if rpo.TxType == "price" {
//			rt = 2 * time.Hour
//		} else {
//			rt = 24 * time.Hour
//		}
//		u := updateReservedProduct{
//			FailedCount: rpo.FailedCount + 1,
//			EndTime:     time.Now().Add(rt).Unix(),
//		}
//		rCol := database.GetCollection("reserveProducts")
//		_, err = rCol.UpdateDocument(context.Background(), rpo.Key, u)
//	}
//	return c.JSON(rpo)
//}
//

// rejectPaymentImage reject an image payment
// @Summary rejects an image payment
// @Description rejects an image payment
// @Tags payment rejection
// @Accept json
// @Produce json
// @Param rejectRequest body rejectRequest true "data"
// @Param key path string true "payment key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
//@Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /rejection/by-image/{key} [post]
func rejectPaymentImage(c *fiber.Ctx) error {
	PaymentKey := c.Params("key")
	ms := new(rejectRequest)
	if err := utils.ParseBodyAndValidate(c, ms); err != nil {
		return c.JSON(err)
	}
	data, err := graphOrder.GetOrderItemsAndPaymentByPaymentKey(PaymentKey, " filter v.isWaitingForPayment==true ")
	if err != nil {
		return c.JSON(err)

	}
	if len(data.OrderItems) == 0 {
		return c.Status(409).JSON("payment is no connected to any order")
	}
	if data.Payment.IsRejected {
		return c.Status(400).SendString("payment already been rejected")
	}
	var rt time.Duration
	if data.Payment.Type == "price" {
		rt = 2 * time.Hour
		//rt = 2 * time.Minute
	} else {
		rt = 24 * time.Hour
	}
	uj := updatePaymentForRejection{
		IsRejected:      true,
		RejectionTime:   time.Now().Add(rt).Unix(),
		RejectionReason: ms.RejectionReason,
	}
	paymentCol := database.GetCollection("gPayment")
	_, err = paymentCol.UpdateDocument(context.Background(), PaymentKey, uj)

	reservedKeyArr, err := updateNumberInPriceAndMoveToReserved(data)

	time.AfterFunc(rt, func() {
		err := deleteReserveProductAddNUmberToPrice(reservedKeyArr)
		if err != nil {
			return
		}
	})

	return c.JSON(reservedKeyArr)
}

func updateNumberInPriceAndMoveToReserved(cArr graphOrder.OrderItemsAndPayment) ([]string, error) {
	rArr := make([]reservedProduct, 0)
	for _, out := range cArr.OrderItems {
		temp := reservedProduct{
			OrderItemKey: out.Key,
			PaymentKey:   cArr.Payment.Key,
			PriceId:      out.PriceId,
			FailedCount:  0,
			Number:       out.Number,
			EndTime:      0,
		}
		rArr = append(rArr, temp)
		priceCol := strings.Split(out.PriceId, "/")[0]
		query := fmt.Sprintf("for i in %v filter i._id==\"%v\" update i with {totalNumber: i.totalNumber - %v} in %v", priceCol, out.PriceId, out.Number, priceCol)
		log.Println(query)

		database.ExecuteGetQuery(query)

	}

	rCol := database.GetCollection("reserveProducts")
	metaArr, errArr, err := rCol.CreateDocuments(context.Background(), rArr)
	if err != nil {
		return nil, fmt.Errorf("%v", errArr)

	}

	return metaArr.Keys(), nil

}

func deleteReserveProductAddNUmberToPrice(reservedKeyArr []string) error {

	rpArr := make([]reservedProduct, 0)
	rCol := database.GetCollection("reserveProducts")

	for _, rk := range reservedKeyArr {
		var rp reservedProduct
		_, err := rCol.ReadDocument(context.Background(), rk, &rp)
		if err != nil {
			return err
		}
		rpArr = append(rpArr, rp)

	}

	for _, out := range rpArr {
		priceCol := strings.Split(out.PriceId, "/")[0]

		query := fmt.Sprintf("for i in %v filter i._id==\"%v\" update i with {totalNumber: i.totalNumber + %v} in %v", priceCol, out.PriceId, out.Number, priceCol)
		log.Println(query)
		database.ExecuteGetQuery(query)
	}

	_, errArr, err := rCol.RemoveDocuments(context.Background(), reservedKeyArr)
	if err != nil {
		return fmt.Errorf("%v", errArr)
	}
	return nil

}

//func deleteReserveProductAddNUmberToPrice(cArr *[]cart.CartOut, reservedKey string, carry int) error {
//	if carry > 10 {
//		log.Fatalf("fucking recersive func \n %v \n\n %v \n\n", *cArr, reservedKey)
//	}
//	var rp reservedProduct
//	rCol := database.GetCollection("reserveProducts")
//	_, err := rCol.ReadDocument(context.Background(), reservedKey, &rp)
//	if err != nil {
//		return err
//	}
//	if rp.EndTime > time.Now().Unix() {
//		remainingTime := time.Duration(rp.EndTime - time.Now().Unix())
//
//		time.AfterFunc(remainingTime*time.Second, func() {
//			deleteReserveProductAddNUmberToPrice(cArr, reservedKey, carry+1)
//		})
//		return nil
//	} else {
//		for _, out := range *cArr {
//			priceCol := strings.Split(out.PriceId, "/")[0]
//
//			query := fmt.Sprintf("for i in %v filter i._id==\"%v\" update i with {totalNumber: i.totalNumber + %v} in %v", priceCol, out.PriceId, out.Number, priceCol)
//			go database.ExecuteGetQuery(query)
//		}
//		return nil
//	}
//}

func getReservedProductByPaymentKey(paymentKey string) (*reservedProductOut, error) {
	query := fmt.Sprintf("for i in reserveProducts filter i.paymentKey == \"%v\" return i ", paymentKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()

	var doc reservedProductOut
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		if !cursor.HasMore() {
			return nil, driver.NoMoreDocumentsError{}
		}

		return nil, err
	}
	return &doc, nil
}
