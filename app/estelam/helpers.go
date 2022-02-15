package estelam

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
)

func isSupplierValid(supplierKey string, estelamCartKey string) (bool, error) {
	query := fmt.Sprintf("let a=(for i in supplierEstelam filter i.estelamCartKey==\"%v\" && i.supplierKey==\"%v\" return i)\nreturn LENGTH(a)", estelamCartKey, supplierKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return false, fmt.Errorf("error while running query:%v", query)
	}
	defer cursor.Close()
	var resp []int
	_, err = cursor.ReadDocument(ctx, &resp)
	if err != nil {
		return false, fmt.Errorf("error while running query:%v", query)
	}
	if resp[0] == 1 {
		return true, nil
	}
	return false, nil
}
