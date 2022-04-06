package addBuyMethod

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"log"
	"strings"
	"time"
)

// getEstelamWithProductBySupplierKey  return products with estelam
// @Summary return products with estelam
// @Description return products with estelam
// @Tags buy method
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []priceAndProduct{}
// @Failure 404 {object} string{}
// @Router /add-buy-method/estelam/{categoryurl} [get]
func getEstelamWithProductBySupplierKey(categoryUrl string, supplierKey string, offset string, limit string) (*[]estelamAndProduct, error) {
	query := fmt.Sprintf("for j in supplier_%v_estelam filter j._from==\"suppliers/%v\" for s in %v filter s._id==j._to sort j.createAt  limit %v,%v return {product:s,estelam:j}", categoryUrl, supplierKey, categoryUrl, offset, limit)
	log.Println(query)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var data []estelamAndProduct
	for {
		var doc estelamAndProduct
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return &data, nil
}

// AddEstelamToProduct   add estelam from supplier to product
// @Summary add estelam from supplier to product
// @Description add estelam from supplier to product
// @Tags buy method
// @Accept json
// @Produce json
// @Param data body estelamIn true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} EstelamOut{}
// @Failure 404 {object} string{}
// @Router /add-buy-method/estelam [post]
func AddEstelamToProduct(est estelamIn, supplierId string) (*EstelamOut, error) {
	supplierKey := supplierId
	productKey := strings.Split(est.ProductId, "/")[1]
	productCol := strings.Split(est.ProductId, "/")[0]

	variantArr := getVariationsName(productKey, productCol)
	flag := true
	index := -1
	for i, v := range variantArr {
		if v == est.Variant {
			flag = false
			index = i
		}

	}
	if flag {
		return nil, fmt.Errorf("variant does not match product variant: %v", est.Variant)
	}

	if index == -1 {
		return nil, fmt.Errorf("index of varent is not acceptable")
	}

	//ce := CreateEstelam{
	//	From:            supplierId,
	//	To:              est.ProductId,
	//	Key:             supplierKey + productKey,
	//	VariantArr:      est.VariantArr,
	//	CodeForSupplier: est.CodeForSupplier,
	//
	//	CreatedAt: time.Now().Unix(),
	//}

	ce := CreateEstelam{
		From:            fmt.Sprintf("supplier/%v", supplierId),
		To:              est.ProductId,
		Key:             fmt.Sprintf("%v%v%v", supplierKey, productKey, index),
		CodeForSupplier: est.CodeForSupplier,
		Variant:         est.Variant,
		Price:           est.Price,
		OneMonthPrice:   est.OneMonthPrice,
		TwoMonthPrice:   est.TwoMonthPrice,
		ThreeMonthPrice: est.ThreeMonthPrice,
		Show:            est.Show,
		CreatedAt:       time.Now().Unix(),
	}

	colName := strings.Split(est.ProductId, "/")[0]
	edgeCol := database.GetCollection("supplier_" + colName + "_estelam")

	var estOut EstelamOut
	ctx := driver.WithReturnNew(context.Background(), &estOut)

	_, err := edgeCol.CreateDocument(ctx, ce)
	if err != nil {
		log.Println(ce)
		return nil, err
	}
	return &estOut, nil

}

// updateEstelamOfProduct    update estelam from supplier to product
// @Summary update estelam from supplier to product
// @Description update estelam from supplier to product
// @Tags buy method
// @Accept json
// @Produce json
// @Param   estelamColName     path   string     true  "estelam col name"
// @Param   estelamKey     path   string     true  "estelam key"
// @Param data body updateEstelam true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 204 {object} EstelamOut{}
// @Failure 404 {object} string{}
// @Router /add-buy-method/estelam/{estelamColName}/{estelamKey} [put]
func updateEstelamOfProduct(data updateEstelam, priceCol string, priceKey string) (*EstelamOut, error) {
	var po EstelamOut
	col := database.GetCollection(priceCol)

	ctx := driver.WithReturnNew(context.Background(), &po)
	_, err := col.UpdateDocument(ctx, priceKey, data)
	if err != nil {
		return nil, err
	}
	return &po, nil
}

// updateEstelamOfProduct    update group os estelams from supplier to product
// @Summary update group os estelams from supplier to product
// @Description update group os estelams from supplier to product
// @Tags buy method
// @Accept json
// @Produce json
// @Param   estelamColName     path   string     true  "estelam col name"
// @Param data body groupUpdateEstelamIn true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 204 {object} []EstelamOut{}
// @Failure 404 {object} string{}
// @Router /add-buy-method/estelam/group_update/{estelamColName} [put]
func groupUpdateEstelam(data groupUpdateEstelamIn, priceCol string) (*[]EstelamOut, error) {
	col := database.GetCollection(priceCol)
	resp := make([]EstelamOut, 0)
	var doc EstelamOut
	ctx := driver.WithReturnNew(context.Background(), &doc)
	if data.ChangeBuyMode {
		u1 := groupUpdateEstelam1{
			Price:           data.Price,
			OneMonthPrice:   data.OneMonthPrice,
			TwoMonthPrice:   data.TwoMonthPrice,
			ThreeMonthPrice: data.ThreeMonthPrice,
		}
		errArr := make([]error, 0)
		for _, key := range data.PriceKeys {
			log.Println(key)
			_, err := col.UpdateDocument(ctx, key, u1)
			if err != nil {
				errArr = append(errArr, err)
			} else {
				resp = append(resp, doc)
			}
		}
		if len(errArr) > 0 {
			return nil, fmt.Errorf("%v ", errArr)
		}
	}
	if data.ChangeStatus {
		u2 := groupUpdateEstelam2{Show: data.Show}
		errArr := make([]error, 0)
		for _, key := range data.PriceKeys {
			_, err := col.UpdateDocument(context.Background(), key, u2)
			if err != nil {
				errArr = append(errArr, err)
			} else {
				resp = append(resp, doc)
			}
		}
		if len(errArr) > 0 {
			return nil, fmt.Errorf("%v ", errArr)
		}
	}
	return &resp, nil

}

// deleteEstelam    delete estelam from supplier to product
// @Summary delete estelam from supplier to product
// @Description delete estelam from supplier to product
// @Tags buy method
// @Accept json
// @Produce json
// @Param   categoryUrl     path   string     true  "category url"
// @Param   productKey     path   string     true  "product key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 204
// @Failure 404 {object} string{}
// @Router /add-buy-method/estelam/{categoryUrl}/{productKey} [delete]
func deleteEstelam(productKey string, productCol string) error {
	//productKey := strings.Split(productId, "/")[1]
	//productCol :=strings.Split(productId, "/")[0]
	edgeCol := database.GetCollection("supplier_" + productCol + "_estelam")
	_, err := edgeCol.RemoveDocument(context.Background(), productKey)
	if err != nil {
		return err
	}
	return nil

}
