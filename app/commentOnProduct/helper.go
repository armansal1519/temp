package commentOnProduct

import (
	"bamachoub-backend-go-v1/app/paymentAndWallet"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
)

func getApprovedOrder(userKey string, productId string) (*paymentAndWallet.ApprovedOrderOut, error) {
	query := fmt.Sprintf("for i in approvedOrder filter i.userKey==\"%v\" && i.productId==\"%v\" limit 1 return i", userKey, productId)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error while running query:%v", query)
	}
	defer cursor.Close()

	var doc paymentAndWallet.ApprovedOrderOut
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return nil, fmt.Errorf("error while reading data:%v", query)
	}
	return &doc, nil
}
