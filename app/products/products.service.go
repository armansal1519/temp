package products

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
)

func GetProductByUrl(c *fiber.Ctx) error {
	productName := c.Params("productName")
	query := fmt.Sprintf("for i in products\n filter i.url==\"%v\" let b=(FOR v IN 1..1 OUTBOUND i._id GRAPH 'pro-cat-pro'\nlet a= (FOR vv IN 1..1 OUTBOUND v._id GRAPH 'pro-cat-pro' return vv)\nreturn{\"mainCategory\":v,\"subCategories\":a})\nreturn {\"Product\":i,\"categories\":b}", productName)
	data := database.ExecuteGetQuery(query)
	return c.JSON(data[0])
}

func GetProductBodyByCategoryKey(c *fiber.Ctx) error {
	key := c.Params("key")
	q := fmt.Sprintf("for p in productBody filter p.categoryKey==\"%v\" return p", key)
	data := database.ExecuteGetQuery(q)
	return c.JSON(data[0])

}

func createProductBody(c *fiber.Ctx) error {
	p := new(productBodyDto)
	if err := c.BodyParser(p); err != nil {
		return err
	}
	err := validateProductBody(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println(p)
	var pIn productBodyIn
	pIn.Status = "ok"
	pIn.ImageArr = []string{}
	pIn.Tags = []string{}
	pIn.Description = ""
	pIn.CategoryKey = p.CategoryKey
	pIn.TitleMaker = p.TitleMaker
	pIn.VariationObj = p.VariationObj
	//create map from array for speacs
	//pIn.MainSpecs=make(map[string]string)
	for _, spec := range p.MainSpecs {
		pIn.MainSpecs = append(pIn.MainSpecs, f{
			Name:  spec,
			Value: "",
		})
	}
	for _, cs := range p.CompleteSpec {
		tempMap := make([]f, 0)
		for _, item := range cs.Items {
			tempMap = append(tempMap, f{
				Name:  item,
				Value: "",
			})
		}
		tempCs := csType{
			Name:  cs.Name,
			Items: tempMap,
		}
		pIn.CompleteSpec = append(pIn.CompleteSpec, tempCs)
	}

	productBodyCol := database.GetCollection("productBody")
	meta, err := productBodyCol.CreateDocument(context.Background(), pIn)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

func getProducts(dbName []string) ([]productOut, error) {
	db := database.GetDB()
	//flag, err := db.CollectionExists(context.Background(), dbName)
	//if err != nil {
	//	return []productOut{}, err
	//}
	//if !flag {
	//	return []productOut{}, fmt.Errorf("collection with name : %v dose not exit", dbName)
	//}
	s := ""
	for i, s2 := range dbName {
		s += " " + s2 + " "
		if i < len(dbName)-1 {
			s += " , "
		}
	}
	query := fmt.Sprintf("for i in [%v] for j in i return j", s)
	cursor, err := db.Query(context.Background(), query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []productOut
	for {
		var doc productOut
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return data, nil
}

// getProductWithColorCode  return products with same color code
// @Summary return products with same color code
// @Description return products with same color code
// @Tags products
// @Accept json
// @Produce json
// @Param   spId      path   string     true  "spId"
// @Param   productKey  path   string     true  "productKey"
// @Success 200 {object} colorOut{}
// @Failure 404 {object} string{}
// @Router /products/color/{spId}/{productKey} [post]
func getProductWithColorCode(c *fiber.Ctx) error {
	spId := c.Params("spId")
	productKey := c.Params("productKey")
	query := fmt.Sprintf("let productList=(for i in sheet filter i.spId==\"%v\" return i)\nlet s=(for j in productList filter j._key==\"%v\" return j)\nreturn {main:s[0],sub:REMOVE_VALUE(productList,s[0])}\n", spId, productKey)
	//log.Println(query)
	return c.JSON(database.ExecuteGetQuery(query))
}

// getProductFromCategory  return products attached to that category
// @Summary return products attached to that category
// @Description return products attached to that category
// @Tags products
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param   categorykey      path   string     true  "category key"
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Success 200 {object} []productOut{}
// @Failure 404 {object} string{}
// @Router /products/cat/{categoryurl}/{categorykey} [get]
func getProductFromCategory(dbName string, key string, offset string, limit string) ([]productOut, error) {
	db := database.GetDB()
	flag, err := db.CollectionExists(context.Background(), dbName)
	if err != nil {
		return []productOut{}, err
	}
	if !flag {
		return []productOut{}, fmt.Errorf("collection with name : %v dose not exit", dbName)
	}
	graphPathString := "\"%\",i.graphPath,\"%\""

	query := fmt.Sprintf("for i in categories filter i._key==\"%v\"\nfor s in %v filter like(s.categoryName,concat(%v)) limit %v,%v return s\n", key, dbName, graphPathString, offset, limit)
	log.Println(query)
	cursor, err := db.Query(context.Background(), query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []productOut
	for {
		var doc productOut
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return data, nil
}

// getProductFromCategory  return length of products attached to that category
// @Summary return length of products attached to that category
// @Description return length of products attached to that category
// @Tags products
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param   categorykey      path   string     true  "category key"
// @Success 200 {object} int{}
// @Failure 404 {object} string{}
// @Router /products/length/{categoryurl}/{categorykey} [get]
func getLengthByCategory(c *fiber.Ctx) error {
	dbName := c.Params("dbName")
	key := c.Params("key")
	graphPathString := "\"%\",i.graphPath,\"%\""
	q := fmt.Sprintf("for i in categories filter i._key==\"%v\"\nlet p=(for s in %v filter like(s.categoryName,concat(%v)) return  s)\nreturn length(p)\n", key, dbName, graphPathString)
	resp := database.ExecuteGetQuery(q)
	return c.JSON(fiber.Map{
		"length": resp[0],
	})
}

// getProductByKey  return Product by key
// @Summary return Product by key
// @Description return Product by key with estelamArr and priceArr
// @Tags products
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param   key      path   string     true  "Product key"
// @Success 200 {object} getProductByKeyResponse{}
// @Failure 404 {object} string{}
// @Router /products/one/{categoryurl}/{key} [get]
func getProductByKey(categoryUrl string, key string) (*getProductByKeyResponse, error) {
	query := fmt.Sprintf("for s in %v filter s._key==\"%v\" \nlet est=(for i in supplier_%v_estelam filter i._to==s._id return i)\nlet prc=(for i in supplier_%v_price filter i._to==s._id return i)\nreturn {\"Product\":s,\"estelamArr\":est,\"priceArr\":prc}", categoryUrl, key, categoryUrl, categoryUrl)
	//log.Println(query)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error while running query: %v \n err: %v", query, err)
	}
	defer cursor.Close()
	var data getProductByKeyResponse

	_, err = cursor.ReadDocument(ctx, &data)
	//log.Println(data)

	return &data, nil

}

// updateProduct  updates product
// @Summary updates product
// @Description update product by key
// @Tags products
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param   key      path   string     true  "Product key"
// @Param data body Product true "data"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /products/{categoryurl}/{key} [put]
func updateProduct(c *fiber.Ctx) error {
	categoryUrl := c.Params("categoryUrl")
	key := c.Params("key")
	p := new(Product)
	fmt.Println(p)
	if err := c.BodyParser(p); err != nil {
		return err
	}
	col := database.GetCollection(categoryUrl)
	meta, err := col.UpdateDocument(context.Background(), key, p)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// updateProduct  deletes product
// @Summary deletes product
// @Description delete product by key
// @Tags products
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param   key      path   string     true  "Product key"
// @Success 204 {object} string{}
// @Failure 404 {object} string{}
// @Router /products/{categoryurl}/{key} [delete]
func deleteProduct(c *fiber.Ctx) error {
	categoryUrl := c.Params("categoryUrl")
	key := c.Params("key")
	col := database.GetCollection(categoryUrl)
	meta, err := col.RemoveDocument(context.Background(), key)
	if err != nil {
		return c.JSON(err)
	}
	return c.Status(204).JSON(meta)
}

// basicSearchProducts  basic search in products
// @Summary basic search in products
// @Description basic search in products
// @Tags products
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param data body search true "data"
// @Success 200 {object} []productOut{}
// @Failure 404 {object} string{}
// @Router /products/basic-search/{categoryurl} [post]
func basicSearchProducts(c *fiber.Ctx) error {
	dbName := c.Params("dbName")
	offset := c.Query("offset")
	limit := c.Query("limit")
	ss := new(search)
	if err := utils.ParseBodyAndValidate(c, ss); err != nil {
		return c.JSON(err)
	}
	searchString := "%" + ss.SearchString + "%"
	query := fmt.Sprintf("for i in %v filter like(i.title,\"%v\") limit %v,%v return i", dbName, searchString, offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}

type search struct {
	SearchString string `json:"searchString"`
}

// basicFilter  basic filter in products
// @Summary basic filter in products
// @Description basic filter in products
// @Tags products
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Param data body filter true "data"
// @Success 200 {object} []productOut{}
// @Failure 404 {object} string{}
// @Router /products/basic-filter/{categoryurl} [post]
func basicFilter(c *fiber.Ctx) error {
	dbName := c.Params("dbName")
	offset := c.Query("offset")
	limit := c.Query("limit")
	f := new(filter)
	if err := utils.ParseBodyAndValidate(c, f); err != nil {
		return c.JSON(err)
	}
	query := fmt.Sprintf("for i in %v filter i.brand==\"%v\" limit %v,%v sort i.createdAt return i", dbName, f.Brand, offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}

type filter struct {
	Brand string `json:"brand"`
}

// AdvanceFilter  advance filter in products
// @Summary advance filter in products
// @Description advance filter in products
// @Tags products
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param   categorykey      path   string     true  "category key"
// @Param data body advanceFilter true "data"
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Param   sample      query    bool     true        "sample"
// @Success 200 {object} []productOut{}
// @Failure 404 {object} string{}
// @Router /products/advance-filter/{categoryurl}/{categorykey} [post]
func AdvanceFilter(c *fiber.Ctx) error {
	dbName := c.Params("dbName")
	categoryKey := c.Params("categoryKey")
	offset := c.Query("offset")
	limit := c.Query("limit")
	sample := c.Query("sample")
	f := new(advanceFilter)
	if err := utils.ParseBodyAndValidate(c, f); err != nil {
		return c.JSON(err)
	}
	filterQuery := ""
	if len(f.FilterStringArr) > 0 {
		filterQuery += " filter "
	}
	log.Println(f.FilterStringArr)

	for i, s := range f.FilterStringArr {

		if i < len(f.FilterStringArr) && i > 0 {
			filterQuery += " and "

		}

		if len(s.Values) > 1 {
			for i2, s2 := range s.Values {
				if i2 == len(s.Values)-1 {
					filterQuery = filterQuery + fmt.Sprintf("\"%v=%v\"  in i.filterArr ", s.Name, s2)
				} else {
					filterQuery = filterQuery + fmt.Sprintf("\"%v=%v\"  in i.filterArr  or ", s.Name, s2)
				}
			}
		}

		if len(s.Values) == 1 {
			filterQuery = filterQuery + fmt.Sprintf("\"%v=%v\" in i.filterArr", s.Name, s.Values[0])
		}

		if len(s.Values) == 0 {
			continue
		}
	}

	if len(f.Brand) > 0 {
		brandString := "["
		for i, b := range f.Brand {
			brandString += fmt.Sprintf("\"%v\"", b)
			if i < len(f.Brand) {
				brandString += " , "
			}
		}
		brandString += "] "

		filterQuery = filterQuery + fmt.Sprintf(" filter i.brand in %v  ", brandString)
	}
	if f.PriceFrom > 0 {
		filterQuery = filterQuery + fmt.Sprintf(" filter i.lowestPrice>=%v  ", f.PriceFrom)
	}
	if f.PriceTo > 0 {
		filterQuery = filterQuery + fmt.Sprintf(" filter i.lowestPrice<=%v  ", f.PriceTo)

	}
	if f.InStock {
		filterQuery = filterQuery + fmt.Sprintf(" filter i.lowestPrice != -1  ")

	}

	s := "seenNumber desc"
	if f.Sort != "" {
		if f.Sort == "buy" {
			s = "buyNumber"
		} else if f.Sort == "high" {
			s = "lowestPrice desc"
		} else if f.Sort == "new" {
			s = "createdAt desc"
		} else if f.Sort == "less" {
			s = "lowestPrice "
		} else if f.Sort == "discount" {
			s = "discount desc "
		} else if f.Sort == "spId-desc" {
			s = "spId desc "
		} else if f.Sort == "spId" {
			s = "spId "
		} else {
			return c.Status(400).SendString("sort in not acceptable only : discount / less / new / buy / spId desc / spId  and default : seen number")
		}
	}

	r := " return i"
	r2 := ""
	if sample == "true" {
		r = "  return  groups[0] "
		r2 = " COLLECT sp = i.spId INTO groups "
	}

	query := fmt.Sprintf(" let ck=(for c in categories filter c._key==\"%v\" for v in 0..10 outbound c graph \"categoryGraph\" filter v.status==\"end\" return v._key) for i in %v filter i.categoryKey in ck  sort i.%v  %v %v limit %v,%v %v", categoryKey, dbName, s, filterQuery, r2, offset, limit, r)
	log.Println(11111, query)
	return c.JSON(database.ExecuteGetQuery(query))

}

type advanceFilter struct {
	FilterStringArr []fItem  `json:"FilterStringArr"`
	Brand           []string `json:"brand"`
	PriceFrom       int      `json:"priceFrom"`
	PriceTo         int      `json:"priceTo"`
	InStock         bool     `json:"inStock"`
	Sort            string   `json:"sort"`
}

type fItem struct {
	Name   string   `json:"name"`
	Values []string `json:"value"`
}
