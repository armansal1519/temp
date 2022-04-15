package categories

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

// CreateBaseCategory create base category
// @Summary create base category
// @Description create base category in categories database
// @Tags categories
// @Accept json
// @Produce json
// @Param category body BaseCategoryDto true "category"
// @Success 200 {object} string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /categories/base [post]
func CreateBaseCategory(c *fiber.Ctx) error {
	bc := new(BaseCategoryDto)
	fmt.Println(bc)

	if err := c.BodyParser(bc); err != nil {
		return err
	}
	//utils.Validate(bc)

	errors := utils.Validate(bc)
	if errors != nil {
		c.JSON(errors)
		return nil
	}
	bcSave := BaseCategorySave{
		Name:      bc.Name,
		Url:       bc.Url,
		Status:    "start",
		GraphPath: bc.Name,
		ImageUrl:  bc.ImageUrl,
		Text:      bc.Text,
	}

	newProduct := database.CreateDocument(bcSave, "categories")

	db := database.GetDB()
	ctx := context.Background()
	options := &driver.CreateCollectionOptions{ /* ... */ }
	_, err := db.CreateCollection(ctx, bc.Url, options)
	if err != nil {
		return utils.CustomErrorResponse(409, 4091, err, "", c)
	}

	err = createGraphAndEdge(bc.Url)

	if err != nil {
		panic(err)
	}
	err = createSupplierToProductEdgeAndAddToGraph(bc.Url)

	if err != nil {
		panic(err)
	}
	return c.JSON(newProduct)
}

// CreateCategory  create category
// @Summary create category
// @Description create category in categories database
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CreateCategoryType true "category"
// @Success 200 {object} string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /categories [post]
func CreateCategory(c *fiber.Ctx) error {
	cat := new(CreateCategoryType)
	cp := 5.0

	if err := c.BodyParser(cat); err != nil {
		return err
	}
	errors := utils.Validate(cat)
	if errors != nil {
		c.JSON(errors)
		return nil
	}
	var meta driver.DocumentMeta
	ctx := context.Background()
	categoryCollection := database.GetCollection("categories")

	splitedId := strings.Split(cat.From, "/")
	key := splitedId[1]
	var bc BaseCategorySave
	meta, err := categoryCollection.ReadDocument(ctx, key, &bc)
	if err != nil {
		panic(fmt.Sprintf("175 %v", err))
	}

	if bc.Status == "start" || bc.Status == "mid" {
		newCategory := Category{
			Name:                cat.Name,
			GraphPath:           fmt.Sprintf("%v-%v", bc.GraphPath, cat.Name),
			Status:              "end",
			CommissionPercent:   cp,
			Text:                cat.Text,
			ImageUrl:            cat.ImageUrl,
			CustomerReviewItems: cat.CustomerReviewItems,
			Url:                 cat.Url,
		}
		meta, _ = categoryCollection.CreateDocument(ctx, newCategory)

	} else if bc.Status == "end" {
		newCategory := Category{
			Name:              cat.Name,
			GraphPath:         fmt.Sprintf("%v-%v", bc.GraphPath, cat.Name),
			Status:            "end",
			CommissionPercent: cp,
			Url:               cat.Url,
		}
		meta, _ = categoryCollection.CreateDocument(ctx, newCategory)

		fromStatus := updateStatus{
			Status: "mid",
		}
		_ = database.UpdateDocument(key, fromStatus, "categories")

	} else {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("from: %v is not a valid", cat.From))
	}
	fmt.Println(meta)

	meta = createEdge(cat.From, meta.ID.String())

	return c.JSON(meta)

}

// getCategoryByKey  return at least one level to three level of categories
// @Summary return at least one level to three level of categories
// @Description return at least one level to three level of categories ,first level is main products
// @Tags categories
// @Accept json
// @Produce json
// @Param   key      path   string     true  "category key"
// @Param   level     query    int     true        "level"
// @Success 200 {object} BaseCategoryOut{}
// @Failure 401 {object} ResponseHTTP{}
// @Router /categories/{key} [get]
func getCategoryByKey(c *fiber.Ctx) error {
	key := c.Params("key")
	level := c.Query("level")

	var query string

	if level == "2" {
		query = fmt.Sprintf("for i in categories"+
			" filter i._key==\"%v\" "+
			"let b=(FOR v IN 1..1 OUTBOUND i._id GRAPH 'categoryGraph'"+
			"let a= (FOR vv IN 1..1 OUTBOUND v._id GRAPH 'categoryGraph' return vv)"+
			"return{\"level1\":v,\"level2\":a})\nreturn {\"category\":i,\"subCategories\":b}", key)
	} else if level == "1" {
		query = fmt.Sprintf("for i in categories "+
			" filter i._key==\"%v\" "+
			"let a=( FOR v IN 1..1 OUTBOUND i._id GRAPH 'categoryGraph'   return v)"+
			"return {\"category\":i,\"subCategories\":a}", key)
	} else {
		query = fmt.Sprintf("for i in categories  filter i._key==\"%v\" return i", key)
	}

	fmt.Println(query)
	data := database.ExecuteGetQuery(query)
	fmt.Println(data)
	return c.JSON(data[0])
}

//func getCategoriesFor()  {
//
//}

// getBaseCategories  return at least one level to three level of categories
// @Summary return at least one level to three level of categories
// @Description return at least one level to three level of categories ,first level is main products
// @Tags categories
// @Accept json
// @Produce json
// @Param   level     query    int     true        "level"
// @Success 200 {object} BaseCategoryOut{}
// @Failure 401 {object} ResponseHTTP{}
// @Router /categories [get]
func getBaseCategories(c *fiber.Ctx) error {
	level := c.Query("level")
	var query string

	if level == "2" {
		query = fmt.Sprintf("for i in categories" +
			" filter i.status==\"start\" " +
			"let b=(FOR v IN 1..1 OUTBOUND i._id GRAPH 'categoryGraph'" +
			"let a= (FOR vv IN 1..1 OUTBOUND v._id GRAPH 'categoryGraph' return vv)" +
			"return{\"level1\":v,\"level2\":a})\nreturn {\"category\":i,\"subCategories\":b}")
	} else if level == "1" {
		query = "for i in categories " +
			" filter i.status==\"start\"  " +
			"let a=( FOR v IN 1..1 OUTBOUND i._id GRAPH 'categoryGraph'   return v)" +
			"return {\"category\":i,\"subCategories\":a}"
	} else {
		query = "for i in categories filter i.status==\"start\"  return i"
	}
	data := database.ExecuteGetQuery(query)
	return c.JSON(data)
}

// getPriceRangeUnderOnCategory  return price range in products attached to that category
// @Summary return price range in products attached to that category
// @Description return price range in products attached to that category
// @Tags categories
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param   categoryKey      path   string     true  "category key"
// @Success 200 {object} []int{}
// @Failure 404 {object} string{}
// @Router /categories/price-range/{categoryurl}/{categoryKey} [get]
func getPriceRangeUnderOnCategory(c *fiber.Ctx) error {
	key := c.Params("key")
	dbName := c.Params("dbName")

	db := database.GetDB()
	flag, err := db.CollectionExists(context.Background(), dbName)
	if err != nil {
		return c.JSON(err)
	}
	if !flag {
		return c.Status(404).JSON(fmt.Errorf("collection with name : %v dose not exit", dbName))
	}
	graphPathString := "\"%\",i.graphPath,\"%\""
	query := fmt.Sprintf("for i in categories  filter i._key==\"%v\" for s in %v  filter like(s.categoryName,concat(%v))   collect b =s.lowestPrice RETURN b", key, dbName, graphPathString)
	log.Println(query)
	cursor, err := db.Query(context.Background(), query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []int
	for {
		var doc int
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return c.JSON(data)
}

// update create category
// @Summary update category
// @Description update category
// @Tags categories
// @Accept json
// @Produce json
// @Param CreateCategoryType body CreateCategoryType true "category"
// @Param key path string true "key"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /categories/{key}  [put]
func update(c *fiber.Ctx) error {
	ci := new(CreateCategoryType)
	if err := utils.ParseBodyAndValidate(c, ci); err != nil {
		return c.JSON(err)
	}
	key := c.Params("key")
	contactCol := database.GetCollection("categories")
	meta, err := contactCol.UpdateDocument(context.Background(), key, ci)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

func GetBaseCategory(categoryKey string) (BaseCategoryOut, error) {
	q := fmt.Sprintf("for c in categories\nfilter c._key==\"%v\"\nfor v in 1..10 inbound c graph categoryGraph\nfilter v.status==\"start\"\nreturn v", categoryKey)
	db := database.GetDB()
	cursor, err := db.Query(context.Background(), q, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", q))
	}
	defer cursor.Close()

	var doc BaseCategoryOut
	_, err = cursor.ReadDocument(context.Background(), &doc)
	if err != nil {
		return doc, err
	}
	return doc, nil

}
