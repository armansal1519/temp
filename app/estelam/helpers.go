package estelam

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"log"
)

func isSupplierValid(supplierKey string, estelamCartKey string) (bool, error) {
	query := fmt.Sprintf("let a=(for i in supplierEstelam filter i.estelamCartKey==\"%v\" && i.supplierKey==\"%v\" return i) return LENGTH(a)", estelamCartKey, supplierKey)
	log.Println(query)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return false, fmt.Errorf("error while running query:%v", query)
	}
	defer cursor.Close()
	var resp int
	_, err = cursor.ReadDocument(ctx, &resp)
	if err != nil {
		return false, fmt.Errorf("error while reading data query:%v", query)
	}
	log.Println("resp==", resp)
	if resp == 1 {
		return true, nil
	}
	return false, nil
}
