package addBuyMethod

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
)

func getCategoriesName(isPrice bool) ([]string, error) {
	var buyMethod string
	if isPrice {
		buyMethod = "price"
	} else {
		buyMethod = "estelam"
	}
	q1 := fmt.Sprintf("for i in categories filter i.status==\"start\" return concat(\"supplier_\",i.url,\"_%v\")", buyMethod)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q1, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var data []string
	for {
		var doc string
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err

		}
		data = append(data, doc)
	}
	return data, nil
}
