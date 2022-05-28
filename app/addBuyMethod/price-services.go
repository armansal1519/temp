package addBuyMethod

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
	"time"

	"github.com/arangodb/go-driver"
)

// getPriceWithProductBySupplierKey  return products with price
// @Summary return products with price
// @Description return products with price
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
// @Router /add-buy-method/price/{categoryurl} [get]
func getPriceWithProductBySupplierKey(categoryUrl string, supplierKey string, offset string, limit string) (*[]priceAndProduct, error) {
	query := fmt.Sprintf("for j in supplier_%v_price filter j._from==\"supplier/%v\" for s in %v filter s._id==j._to  limit %v,%v return {product:s,price:j}", categoryUrl, supplierKey, categoryUrl, offset, limit)
	fmt.Println(query)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var data []priceAndProduct
	for {
		var doc priceAndProduct
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err

		}
		data = append(data, doc)
	}
	return &data, nil
}

func getAllPricesWithProductsBySupplierKey(supplierKey string, brand string, search string, offset string, limit string) (*[]priceAndProduct, error) {
	catNameArr, err := getCategoriesName(true)
	if err != nil {
		return nil, err
	}
	var s string
	if search != "" {
		s = fmt.Sprintf("and like(s.title,\"%v\")", "%"+search+"%")
	}
	var b string
	if brand != "" {
		b = fmt.Sprintf("and s.brand==\"%v\"", brand)
	}

	q := fmt.Sprintf("for cat in [%v]\nfor j in cat filter j._from==\"supplier/%v\" for s in sheet filter s._id==j._to %v  %v limit %v,%v return {product:s,price:j}", strings.Join(catNameArr, ","), supplierKey, b, s, offset, limit)
	fmt.Println(q)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var data []priceAndProduct
	for {
		var doc priceAndProduct
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err

		}
		data = append(data, doc)
	}
	return &data, nil
}

func getPriceBrandsBySupplierKey(c *fiber.Ctx) error {

	catNameArr, err := getCategoriesName(true)
	supplierId := c.Locals("supplierId").(string)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	q := fmt.Sprintf("for cat in [%v]\nfor j in cat filter j._from==\"supplier/%v\" for s in sheet filter s._id==j._to collect b=s.brand return b", strings.Join(catNameArr, ","), supplierId)
	fmt.Println(q)
	return c.JSON(database.ExecuteGetQuery(q))
}

// AddPriceToProduct    add price from supplier to product
// @Summary add price from supplier to product
// @Description add price from supplier to product
// @Tags buy method
// @Accept json
// @Produce json
// @Param data body PriceIn true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} PriceOut{}
// @Failure 404 {object} string{}
// @Router /add-buy-method/price [post]
func AddPriceToProduct(pg PriceIn, supplierId string) (*PriceOut, error) {
	q := fmt.Sprintf("for i in supplier_sheet_price filter i._to==\"%v\" && i._from==\"%v\" return i", pg.ProductId, supplierId)
	res := database.ExecuteGetQuery(q)
	//supplierKey := strings.Split(supplierId, "/")[1]
	productKey := strings.Split(pg.ProductId, "/")[1]
	productCol := strings.Split(pg.ProductId, "/")[0]

	variantArr := getVariationsName(productKey, productCol)
	flag := false
	for _, v := range variantArr {
		if v == pg.Variant {
			flag = true
		}
	}
	if !flag {
		return nil, fmt.Errorf("variant does not match product variant: %v", pg.Variant)
	}

	log.Println(variantArr)
	log.Println("edge data:", pg.ProductId, supplierId)
	pgc := PriceGroupCreate{
		To:                pg.ProductId,
		From:              fmt.Sprintf("supplier/%v", supplierId),
		Price:             pg.Price,
		OneMonthPrice:     pg.OneMonthPrice,
		TwoMonthPrice:     pg.TwoMonthPrice,
		ThreeMonthPrice:   pg.ThreeMonthPrice,
		Variant:           pg.Variant,
		TotalNumber:       pg.TotalNumber,
		TotalNumberInCart: pg.TotalNumberInCart,
		CodeForSupplier:   pg.CodeForSupplier,
		PriceRepetition:   len(res),
		CreatedAt:         time.Now().Unix(),
	}

	colName := strings.Split(pg.ProductId, "/")[0]

	edgeCol := database.GetCollection("supplier_" + colName + "_price")

	lp, err := getLowestPrice(pg.ProductId, productCol)
	if err != nil {
		return nil, err
	}
	lowestPrice := compareTwoPrice(lp, pg.Price)
	log.Println(lp, pg.Price)
	ulp := updateLowestPriceInProduct{LowestPrice: lowestPrice}
	pc := database.GetCollection(productCol)
	log.Println(ulp)
	_, err = pc.UpdateDocument(context.Background(), productKey, ulp)
	if err != nil {
		return nil, err
	}

	var po PriceOut
	ctx := driver.WithReturnNew(context.Background(), &po)

	_, err = edgeCol.CreateDocument(ctx, pgc)
	if err != nil {
		log.Println("error edge", pgc)
		return nil, err
	}
	return &po, nil

}

// updatePriceOfProduct    update price from supplier to product
// @Summary update price from supplier to product
// @Description update price from supplier to product
// @Tags buy method
// @Accept json
// @Produce json
// @Param   priceColName     path   string     true  "price col name"
// @Param   PriceKey     path   string     true  "price key"
// @Param data body updatePrice true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 204 {object} PriceOut{}
// @Failure 404 {object} string{}
// @Router /add-buy-method/price/{priceColName}/{PriceKey} [put]
func updatePriceOfProduct(data updatePrice, priceCol string, priceKey string) (*PriceOut, error) {
	var po PriceOut
	col := database.GetCollection(priceCol)

	ctx := driver.WithReturnNew(context.Background(), &po)
	_, err := col.UpdateDocument(ctx, priceKey, data)
	if err != nil {
		return nil, err
	}
	return &po, nil
}

// groupUpdatePrice    update group of prices from supplier to product
// @Summary update group of prices from supplier to product
// @Description update group of prices from supplier to product
// @Tags buy method
// @Accept json
// @Produce json
// @Param   estelamColName     path   string     true  "price col name"
// @Param data body groupUpdatePriceIn true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 204 {object} []PriceOut{}
// @Failure 404 {object} string{}
// @Router /add-buy-method/price/group_update/{estelamColName} [put]
func groupUpdatePrice(data groupUpdatePriceIn, priceCol string) (*[]PriceOut, error) {
	col := database.GetCollection(priceCol)
	keyArrLen := len(data.PriceKeys)
	fmt.Println(keyArrLen)
	fmt.Println(data)
	resp := make([]PriceOut, keyArrLen)

	keyStr := "["
	for i, key := range data.PriceKeys {
		keyStr += fmt.Sprintf("\"%v\"", key)
		if i < keyArrLen-1 {
			keyStr += " , "
		}
	}
	keyStr += "] "

	ctx := driver.WithReturnNew(context.Background(), resp)
	if data.ChangeStatus {
		u1 := groupUpdatePrice1{Show: data.Show}
		u1Arr := make([]groupUpdatePrice1, keyArrLen)
		for i := 0; i < keyArrLen; i++ {
			u1Arr = append(u1Arr, u1)
		}
		_, errArr, err := col.UpdateDocuments(ctx, data.PriceKeys, u1Arr)
		if err != nil {
			return nil, fmt.Errorf("%v \n %v", err, errArr)
		}
	}
	if data.ChangeNumber {
		changeQuary := ""
		fmt.Println("in change number")
		if data.ChangeNumberMethod == "add" {
			changeQuary += fmt.Sprintf("totalNumber : i.totalNumber + %v ", data.ChangeNumberValue)
		} else if data.ChangeNumberMethod == "mines" {
			changeQuary += fmt.Sprintf("totalNumber : i.totalNumber - %v ", data.ChangeNumberValue)
		} else if data.ChangeNumberMethod == "replace" {
			changeQuary += fmt.Sprintf("totalNumber :  %v ", data.ChangeNumberValue)
		} else if data.ChangeNumberMethod == "zero" {
			changeQuary += fmt.Sprintf("totalNumber : 0")
		} else {
			return nil, fmt.Errorf("only acceptable ChangeNumberMethod: mines add replace zero")
		}

		query := fmt.Sprintf("for i in %v filter i._key in %v update i with { %v } in %v return NEW", priceCol, keyStr, changeQuary, priceCol)
		fmt.Println("change price ", query)
		docs, err := runUpdatePriceQuery(query)
		if err != nil {
			return nil, err
		}
		resp = *docs
	}
	if data.ChangePrice {

		priceTypeArr := make([]string, 0)
		if data.Price {
			priceTypeArr = append(priceTypeArr, "price")
		}
		if data.OneMonthPrice {
			priceTypeArr = append(priceTypeArr, "oneMonthPrice")
		}
		if data.TwoMonthPrice {
			priceTypeArr = append(priceTypeArr, "twoMonthPrice")
		}
		if data.ThreeMonthPrice {
			priceTypeArr = append(priceTypeArr, "threeMonthPrice")
		}
		changeQuary := ""

		if data.ChangePriceMethod == "replace" {
			for i, s := range priceTypeArr {
				changeQuary += fmt.Sprintf(" %v: %v ", s, data.ChangePriceValue)
				if i < len(priceTypeArr)-1 {
					changeQuary += " , "
				}
			}
		} else if data.ChangePriceMethod == "add" {
			for i, s := range priceTypeArr {
				changeQuary += fmt.Sprintf(" %v: i.%v + %v ", s, s, data.ChangePriceValue)
				if i < len(priceTypeArr)-1 {
					changeQuary += " , "
				}
			}
		} else if data.ChangePriceMethod == "mines" {
			for i, s := range priceTypeArr {
				changeQuary += fmt.Sprintf(" %v: i.%v - %v ", s, s, data.ChangePriceValue)
				if i < len(priceTypeArr)-1 {
					changeQuary += " , "
				}
			}
		} else if data.ChangePriceMethod == "add-percent" {
			for i, s := range priceTypeArr {
				changeQuary += fmt.Sprintf(" %v: i.%v + floor(i.%v * 0.%v )", s, s, s, data.ChangePriceValue)
				if i < len(priceTypeArr)-1 {
					changeQuary += " , "
				}
			}
		} else if data.ChangePriceMethod == "mines-percent" {
			for i, s := range priceTypeArr {
				changeQuary += fmt.Sprintf(" %v: i.%v - floor(i.%v * 0.%v )", s, s, s, data.ChangePriceValue)
				if i < len(priceTypeArr)-1 {
					changeQuary += " , "
				}
			}
		} else {
			return nil, fmt.Errorf("only acceptable ChangeNumberMethod: mines add replace mines-percent add-percent")
		}
		query := fmt.Sprintf("for i in %v filter i._key in %v  update i with { %v } in %v return NEW", priceCol, keyStr, changeQuary, priceCol)
		docs, err := runUpdatePriceQuery(query)
		if err != nil {
			return nil, err
		}
		resp = *docs

	}
	return &resp, nil

}

// deletePrice    delete price from supplier to product
// @Summary delete price from supplier to product
// @Description delete price from supplier to product
// @Tags buy method
// @Accept json
// @Produce json
// @Param   categoryUrl     path   string     true  "category url"
// @Param   productKey     path   string     true  "product key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 204
// @Failure 404 {object} string{}
// @Router /add-buy-method/price/{categoryUrl}/{productKey} [delete]
func deletePrice(productKey string, productCol string) error {
	//productKey := strings.Split(productId, "/")[1]
	//productCol :=strings.Split(productId, "/")[0]
	edgeCol := database.GetCollection("supplier_" + productCol + "_price")
	_, err := edgeCol.RemoveDocument(context.Background(), productKey)
	if err != nil {
		return err
	}
	return nil

}

func getVariationsName(productKey string, productCol string) []string {
	query := fmt.Sprintf("for i in %v filter i._key==\"%v\" let x= i.variationsObj.variations for j in x return j.name", productCol, productKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []string
	for {
		var doc string
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return data

}

func getLowestPrice(productId string, productCol string) (int64, error) {
	query := fmt.Sprintf("for i in supplier_%v_price filter i._to==\"%v\" && i.price!=0 sort i.price return i.price", productCol, productId)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []int64
	for {
		var doc int64
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	if len(data) > 0 {
		return data[0], nil
	}
	return -1, nil
}

func compareTwoPrice(a int64, b int64) int64 {
	if a <= 0 && b <= 0 {
		return -1000
	}
	if a <= 0 && b > 0 {
		return b
	}
	if a > 0 && b <= 0 {
		return a
	}
	if a > b {
		return b
	} else {
		return a
	}
}

func runUpdatePriceQuery(query string) (*[]PriceOut, error) {
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []PriceOut
	for {
		var doc PriceOut
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Println(err)
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return &data, nil
}
