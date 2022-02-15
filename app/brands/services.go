package brands

import (
	"bamachoub-backend-go-v1/app/categories"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
)


// GetBrandByKey  get brand by key
// @Summary get brand by key
// @Description get brand by key
// @Tags brands
// @Accept json
// @Produce json
// @Param   Key      path   string     true  " key"
// @Success 200 {object} Brand{}
// @Failure 404 {object} string{}
// @Router /brands/{Key} [get]
func GetBrandByKey(c *fiber.Ctx)error {
	key := c.Params("key")
	q:=fmt.Sprintf("for i in brands filter i._key==\"%v\" update i with {seen:i.seen + 1} in brands return NEW",key)
	res:=database.ExecuteGetQuery(q)
	return c.JSON(res[0])

}


// getAllBrands  get all brand
// @Summary get all brand
// @Description get all brand
// @Tags brands
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Param   sort      query    string     true        "only seen or name"
//@Success 200 {object} []Brand{}
// @Failure 404 {object} string{}
// @Router /brands [get]
func getAllBrands(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	sort := c.Query("sort")

	limitStr:=""
	if limit !="" && offset!="" {
		limitStr=fmt.Sprintf(" limit %v,%v ",offset,limit)
	}
	if sort != "seen" && sort!="name" {
		return c.Status(400).JSON(fiber.Map{
			"error":"seen query only can be seen or name",
		})
	}
	if sort=="name" {
		query := fmt.Sprintf("for i in brands sort i.name %v return i",limitStr)
		return c.JSON(database.ExecuteGetQuery(query))
	}
	query := fmt.Sprintf("for i in brands sort i.seen desc limit 24 return i")
	return c.JSON(database.ExecuteGetQuery(query))


}


// getAllBrandsByCategoryUrl  get all brand by category url
// @Summary get all brand by category url
// @Description get all brand by category url
// @Tags brands
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Param   sort      query    string     true        "only seen or name"
// @Param   categoryurl      path   string     true  "categoryurl"
//@Success 200 {object} []Brand{}
// @Failure 404 {object} string{}
// @Router /brands/url/{categoryurl} [get]
func getAllBrandsByCategoryUrl(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	sort := c.Query("sort")
	categoryUrl:=c.Params("categoryurl")

	limitStr:=""
	if limit !="" && offset!="" {
		limitStr=fmt.Sprintf(" limit %v,%v ",offset,limit)
	}
	if sort != "seen" && sort!="name" {
		return c.Status(400).JSON(fiber.Map{
			"error":"seen query only can be seen or name",
		})
	}
	if sort=="name" {
		query := fmt.Sprintf("for i in categories filter i.url==\"%v\"\nfor v,e,p in 1..1 any i graph brandCategory sort v.name %v return v",categoryUrl,limitStr)
		return c.JSON(database.ExecuteGetQuery(query))
	}
	query := fmt.Sprintf("for i in categories filter i.url==\"%v\"\nfor v,e,p in 1..1 any i graph brandCategory sort v.seen desc %v return v" ,categoryUrl,limitStr)
	log.Println(query)
	return c.JSON(database.ExecuteGetQuery(query))


}



func getBrandsByCategoryName(c *fiber.Ctx) error{
	s := new(cat)

	if err := utils.ParseBodyAndValidate(c, s); err != nil {
		return c.JSON(err)
	}
	q := fmt.Sprintf("for i in categories filter i.name==\"%v\"\nfor v,e,p in 1..1 any i graph brandCategory return v", s.Cat)
	log.Println(q)
	return c.JSON(database.ExecuteGetQuery(q))
}

// getBrandsUsedUnderCategoryByCategoryKey  return brands that used in products attached to that category
// @Summary return brands that used in products attached to that category
// @Description return brands that used in products attached to that category
// @Tags brands
// @Accept json
// @Produce json
// @Param   categoryurl      path   string     true  "category url"
// @Param   categoryKey      path   string     true  "category key"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /brands/used/{categoryurl}/{categoryKey} [get]
func getBrandsUsedUnderCategoryByCategoryKey(categoryKey string, dbName string) ([]string, error) {
	db := database.GetDB()
	flag, err := db.CollectionExists(context.Background(), dbName)
	if err != nil {
		return []string{}, err
	}
	if !flag {
		return []string{}, fmt.Errorf("collection with name : %v dose not exit", dbName)
	}
	graphPathString := "\"%\",i.graphPath,\"%\""
	query := fmt.Sprintf("for i in categories  filter i._key==\"%v\" for s in %v  filter like(s.categoryName,concat(%v))   collect b =s.brand return b", categoryKey, dbName, graphPathString)
	cursor, err := db.Query(context.Background(), query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []string
	for {
		var doc string
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

// CreateBrand  create a brand
// @Summary create a brand
// @Description create a brand
// @Tags brands
// @Accept json
// @Produce json
// @Param data body brandDto true "data"
// @Success 200 {object} Brand{}
// @Failure 404 {object} string{}
// @Router /brands [post]
func CreateBrand(brand brandDto) (Brand, error) {
	brandCol := database.GetCollection("brands")
	catCol := database.GetCollection("categories")
	var b Brand
	var cat categories.BaseCategorySave
	catMeta, err := catCol.ReadDocument(context.Background(), brand.CategoryKey, &cat)
	if cat.Status != "start" {
		return b, fiber.NewError(fiber.StatusBadRequest, "category is not a base category")
	}
	ctx := driver.WithReturnNew(context.Background(), &b)
	brandMeta, err := brandCol.CreateDocument(ctx, brand)
	if err != nil {
		return b, err
	}
	e := database.MyEdgeObject{
		From: brandMeta.ID.String(),
		To:   catMeta.ID.String(),
	}
	edgeCol := database.GetEdgeCollection("brandCategory", "brandsToCategories")
	_, err = edgeCol.CreateDocument(context.Background(), e)

	return b, nil
}

// updateBrand  update a brand
// @Summary update a brand
// @Description update a brand
// @Tags brands
// @Accept json
// @Produce json
// @Param data body editBrand true "data"
// @Param   key      path   string     true  "key"
// @Success 200 {object} Brand{}
// @Failure 404 {object} string{}
// @Router /brands/{key} [put]
func updateBrand(c *fiber.Ctx) error {
	key := c.Params("key")
	eb := new(editBrand)
	if err := utils.ParseBodyAndValidate(c, eb); err != nil {
		return c.JSON(err)
	}
	brandCol := database.GetCollection("brands")

	var b Brand
	ctx := driver.WithReturnNew(context.Background(), &b)
	_, err := brandCol.UpdateDocument(ctx, key, eb)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(b)
}

// removeBrand get brand
// @Summary delete brand
// @Description delete brand
// @Tags brands
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /brands/{key} [delete]
func removeBrand(c *fiber.Ctx) error {
	key := c.Params("key")
	brandCol := database.GetCollection("brands")

	_, err := brandCol.RemoveDocument(context.Background(), key)
	if err != nil {
		return c.JSON(err)
	}
	return c.Status(204).SendString("deleted")
}
