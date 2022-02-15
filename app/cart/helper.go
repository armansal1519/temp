package cart

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
)

//func reserveProduct(priceKey string, priceColName string,number int,userKey string,cartKey string) error {
//	rCol :=database.GetCollection("reserveProducts")
//	r:=rp{
//		PriceCol:  priceColName,
//		PriceKey:  priceKey,
//		UserKey:   userKey,
//		Number:    number,
//		CreatedAt: time.Now().Unix(),
//		Status:    "reserve",
//	}
//	meta,err:= rCol.CreateDocument(context.Background(),r)
//	if err != nil {
//		return err
//	}
//	query := fmt.Sprintf("for i in %v filter i._key==\"%v\" update i with {totalNumber:i.totalNumber - %v} in %v ",priceColName,priceKey,number,priceColName)
//	database.ExecuteGetQuery(query)
//
//	errCode:=0
//	time.AfterFunc(2*time.Hour , func() {
//		var rBack rp
//		_,err=rCol.ReadDocument(context.Background(),meta.Key,&rBack)
//		if err != nil	{
//			errCode=-1
//		}
//		if rBack.Status=="buy" {
//			_,err:=rCol.ReadDocument()
//		}
//
//
//	})
//
//
//}

type rp struct {
	PriceKey  string `json:"priceKey"`
	PriceCol  string `json:"priceCol"`
	UserKey   string `json:"userKey"`
	Number    int    `json:"number"`
	CreatedAt int64  `json:"createdAt"`
	Status    string `json:"status"`
}

func GetCartGroupByBuyingMethod(userKey string) (*[]GroupedCart, error) {
	query := fmt.Sprintf("for i in cart\nfilter i.userKey==\"%v\" COLLECT PricingTypes = i.PricingType INTO carts  RETURN {\"type\" : PricingTypes,\"cart\" : carts[*].i \n  }", userKey)
	db := database.GetDB()
	cursor, err := db.Query(context.Background(), query, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var data []GroupedCart
	for {
		var doc GroupedCart
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		data = append(data, doc)
	}
	return &data, nil

}
