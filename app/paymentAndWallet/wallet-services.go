package paymentAndWallet

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/payment"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func createSupplierWalletHistory(amount int64, supplierKey string, in bool, txType string) error {
	wh := SupplierWalletHistory{
		Amount:      amount,
		SupplierKey: supplierKey,
		CreatedAt:   time.Now().Unix(),
		Income:      in,
		TxType:      txType,
	}
	col := database.GetCollection("supplierWalletHistory")
	_, err := col.CreateDocument(context.Background(), wh)
	if err != nil {
		return err
	}
	return nil
}

func GetSupplierWalletHistoryBySupplierKey(supplierKey string, fromTime int64) (*[]supplierWalletOut, error) {
	query := fmt.Sprintf("for i in supplierWalletHistory  filter i.supplierKey==\"%v\" and i.createdAt > %v return i", supplierKey, fromTime)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error while running query:%v", query)
	}
	defer cursor.Close()
	var data []supplierWalletOut
	for {
		var doc supplierWalletOut
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error while getting data from query:%v", err)

		}
		data = append(data, doc)
	}
	return &data, err
}

// getDataForSupplierWalletPage  return wallet amount , income , outcome
// @Summary return wallet amount , income , outcome
// @Description return wallet amount , income , outcome and if withHistory query eq ture wallet history
// @Tags wallet
// @Accept json
// @Produce json
// @Param   withHistory      query    bool     true        "withHistory"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} supplierPageResponse{}
// @Failure 404 {object} string{}
// @Router /wallet/data [get]
func getDataForSupplierWalletPage(c *fiber.Ctx) error {
	supplierId := c.Locals("supplierId").(string)
	supplierKey := supplierId
	withHistory := c.Query("withHistory")
	query := fmt.Sprintf("let data=(for i in supplierWalletHistory filter i.supplierKey==\"%v\" and i.createdAt>%v return i )\nlet income=(for i in data  filter i.income==true return i.amount)\nlet outcome=(for i in data filter i.income==false return i.amount)\nreturn {income:sum(income),outcome:sum(outcome)}", supplierKey, time.Now().Add(-30*24*time.Hour).Unix())
	log.Println(query)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return c.JSON(err)
	}
	defer cursor.Close()

	var doc temp1
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return c.JSON(err)
	}

	w, err := getWalletAmount(supplierKey)
	if err != nil {
		return c.JSON(err)
	}
	if withHistory == "false" {
		resp := supplierPageResponse{
			WalletAmount: w,
			TotalIn:      doc.Income,
			TotalOut:     doc.Outcome,
			History:      nil,
		}
		return c.JSON(resp)
	}
	swh, err := GetSupplierWalletHistoryBySupplierKey(supplierKey, 0)
	if err != nil {
		return c.JSON(err)
	}
	resp := supplierPageResponse{
		WalletAmount: w,
		TotalIn:      doc.Income,
		TotalOut:     doc.Outcome,
		History:      *swh,
	}
	return c.JSON(resp)

}

func getWalletAmount(supplierKey string) (int64, error) {
	query := fmt.Sprintf("for i in suppliers\nfilter i._key==\"%v\"\nreturn i.walletAmount", supplierKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return 0, fmt.Errorf("error while running query:%v", query)
	}
	defer cursor.Close()

	var doc int64
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return 0, fmt.Errorf("error while running query:%v", query)
	}

	return doc, err
}

// getPaymentUrlForSupplierWallet  get payment url for wallet
// @Summary get payment url for wallet
// @Description get payment url for wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param data body addToWallet true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /wallet/url [post]
func getPaymentUrlForSupplierWallet(c *fiber.Ctx) error {
	w := new(addToWallet)
	if err := utils.ParseBodyAndValidate(c, w); err != nil {
		return c.JSON(err)
	}
	supplierId := c.Locals("supplierId").(string)
	//supplierKey := strings.Split(supplierId, "/")[1]
	supplierKey := supplierId

	p := paymentHistory{
		PayerKey:      supplierKey,
		OrderKey:      "-",
		TxType:        "supplier-wallet-add",
		Amount:        w.Amount,
		Status:        "not",
		CardHolder:    "",
		ShaparakRefId: "",
		TransId:       "",
		ImageUrl:      "-",
		CheckNumber:   "-",
	}
	paymentCol := database.GetCollection("paymentHistory")

	meta, err := paymentCol.CreateDocument(context.Background(), p)
	if err != nil {
		return c.JSON(err)
	}
	transId, err := payment.GetPaymentUrl(fmt.Sprintf("%v", w.Amount), meta.Key, fmt.Sprintf("https://choonet.com/wallet-return.html?id=%v", meta.Key))
	if err != nil {
		return c.JSON(err)
	}
	ut := updateTransId{TransId: transId}
	_, err = paymentCol.UpdateDocument(context.Background(), meta.Key, ut)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(fiber.Map{
		"paymentUrl": fmt.Sprintf("https://nextpay.org/nx/gateway/payment/%v", transId),
	})

}

// getPaymentUrlForSupplierWallet  get payment url for wallet
// @Summary get payment url for wallet
// @Description get payment url for wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param data body addToWallet true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /wallet/user/url [post]
func getPaymentUrlForUserWallet(c *fiber.Ctx) error {
	w := new(addToWallet)
	if err := utils.ParseBodyAndValidate(c, w); err != nil {
		return c.JSON(err)
	}
	userKey := c.Locals("userKey").(string)

	p := paymentHistory{
		PayerKey:      userKey,
		OrderKey:      "-",
		TxType:        "user-wallet-add",
		Amount:        w.Amount,
		Status:        "not",
		CardHolder:    "",
		ShaparakRefId: "",
		TransId:       "",
		ImageUrl:      "-",
		CheckNumber:   "-",
	}
	paymentCol := database.GetCollection("paymentHistory")

	meta, err := paymentCol.CreateDocument(context.Background(), p)
	if err != nil {
		return c.JSON(err)
	}
	transId, err := payment.GetPaymentUrl(fmt.Sprintf("%v", w.Amount), meta.Key, fmt.Sprintf("http://localhost:3000/payment/verify/%v", meta.Key))
	if err != nil {
		return c.JSON(err)
	}
	ut := updateTransId{TransId: transId}
	_, err = paymentCol.UpdateDocument(context.Background(), meta.Key, ut)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(fiber.Map{
		"paymentUrl": fmt.Sprintf("https://nextpay.org/nx/gateway/payment/%v", transId),
	})

}

// VerifyPaymentAndAddToWallet  verify tx to wallet
// @Summary verify tx to wallet
// @Description verify tx to wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param   key      path   string     true  "key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /wallet/verify/{key} [post]
func VerifyPaymentAndAddToWallet(c *fiber.Ctx) error {
	PaymentKey := c.Params("key")
	supplierId := c.Locals("supplierId").(string)
	//supplierId := strings.Split(supplierId, "/")[1]
	var p paymentOut
	paymentCol := database.GetCollection("paymentHistory")
	_, err := paymentCol.ReadDocument(context.Background(), PaymentKey, &p)
	if err != nil {
		return c.JSON(err)
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
	meta, err := paymentCol.UpdateDocument(context.Background(), PaymentKey, up)
	if err != nil {
		return c.JSON(err)
	}
	sw := SupplierWalletHistory{
		Amount:      p.Amount,
		SupplierKey: supplierId,
		PaymentKey:  meta.Key,
		CreatedAt:   time.Now().Unix(),
		Income:      true,
		TxType:      p.TxType,
		TxStatus:    "done",
	}
	walletCol := database.GetCollection("supplierWalletHistory")
	_, err = walletCol.CreateDocument(context.Background(), sw)
	if err != nil {
		return c.JSON(err)
	}
	query := fmt.Sprintf("for i in suppliers filter i._key==\"%v\" update i with {walletAmount: i.walletAmount + %v } in suppliers ", supplierId, p.Amount)
	database.ExecuteGetQuery(query)
	return c.JSON(fiber.Map{"status": "ok"})

}

// VerifyUserPaymentAndAddToWallet  verify tx to wallet
// @Summary verify tx to wallet
// @Description verify tx to wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param   key      path   string     true  "key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /wallet/user/verify/{key} [post]
func VerifyUserPaymentAndAddToWallet(c *fiber.Ctx) error {
	PaymentKey := c.Params("key")
	userKey := c.Locals("userKey").(string)
	//supplierId := strings.Split(supplierId, "/")[1]
	var p paymentOut
	paymentCol := database.GetCollection("paymentHistory")
	_, err := paymentCol.ReadDocument(context.Background(), PaymentKey, &p)
	if err != nil {
		return c.JSON(err)
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
	meta, err := paymentCol.UpdateDocument(context.Background(), PaymentKey, up)
	if err != nil {
		return c.JSON(err)
	}
	sw := UserWalletHistory{
		Amount:     p.Amount,
		UserKey:    userKey,
		PaymentKey: meta.Key,
		CreatedAt:  time.Now().Unix(),
		Income:     true,
		TxType:     p.TxType,
		TxStatus:   "done",
	}
	walletCol := database.GetCollection("userWalletHistory")
	_, err = walletCol.CreateDocument(context.Background(), sw)
	if err != nil {
		return c.JSON(err)
	}
	query := fmt.Sprintf("for i in users filter i._key==\"%v\" update i with {walletAmount: i.walletAmount + %v } in users ", userKey, p.Amount)
	database.ExecuteGetQuery(query)
	return c.JSON(fiber.Map{"status": "ok"})

}

// GetUserWalletHistoryByUserKey   get user wallet history
// @Summary get user wallet history
// @Description get user wallet history
// @Tags wallet
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
//@Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []userWalletOut{}
// @Failure 404 {object} string{}
// @Router /wallet/user/history [get]
func GetUserWalletHistoryByUserKey(supplierKey string, offset string, limit string) (*[]userWalletOut, temp1, error) {
	query := fmt.Sprintf("for i in userWalletHistory  filter i.userKey==\"%v\" limit %v,%v return i", supplierKey, offset, limit)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return nil, temp1{}, fmt.Errorf("error while running query:%v", query)
	}
	defer cursor.Close()
	var data []userWalletOut
	for {
		var doc userWalletOut
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, temp1{}, fmt.Errorf("error while getting data from query:%v", err)

		}
		data = append(data, doc)
	}
	t := time.Now().Add(-30 * 24 * time.Hour).Unix()
	q := fmt.Sprintf("let income=( for i in userWalletHistory filter i.userKey==\"%v\" and i.income==true and i.createdAt> %v  return i.amount)\nlet outcome=( for i in userWalletHistory filter i.userKey==\"%v\" and i.income==false and i.createdAt> %v   return i.amount)\nreturn {income:sum(income),outcome:sum(outcome)}", supplierKey, t, supplierKey, t)
	cursor, err = db.Query(ctx, q, nil)
	if err != nil {
		return nil, temp1{}, fmt.Errorf("error while running query:%v", query)
	}
	var d temp1
	_, err = cursor.ReadDocument(ctx, &d)
	if err != nil {
		return nil, temp1{}, fmt.Errorf("error while running query:%v", query)
	}

	return &data, d, err
}
