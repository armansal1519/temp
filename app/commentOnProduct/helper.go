package commentOnProduct

import (
	"bamachoub-backend-go-v1/app/graphOrder"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
)

func getApprovedOrder(userKey string, productId string) (bool, error) {
	query := fmt.Sprintf("for u in users filter u._key==\"%v\" \nfor v,e in 3..3 outbound u graph \"orderGraph\" filter v.productId==\"sheet/%v\" return v", userKey, productId)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return false, fmt.Errorf("error while running query:%v", query)
	}
	defer cursor.Close()

	var doc graphOrder.GOrderItemOut
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return false, fmt.Errorf("error while reading data:%v", query)
	}
	if cursor.Count() > 0 {
		return true, nil
	}
	return false, nil
}
