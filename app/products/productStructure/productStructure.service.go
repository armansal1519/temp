package productStructure

import (
	"bamachoub-backend-go-v1/app/products/menu"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"errors"
	"fmt"
	driver "github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

func getProductStructureByCategoryKey(c *fiber.Ctx) error {
	key := c.Params("key")
	query := fmt.Sprintf("for i in categories\nfilter i._key==\"%v\"\nfor j in productStructures\nfilter i.productStructureKey ==j._key\nreturn j\n", key)
	data := database.ExecuteGetQuery(query)
	return c.JSON(data)
}

func createProductStructure(c *fiber.Ctx) error {
	b := new(CreateProductStructureDto)

	if err := utils.ParseBodyAndValidate(c, b); err != nil {
		panic(fmt.Sprintf("error parsing dto:%v", err))
	}
	if getCategoryStatusByKey(b.CategoryKey) != "end" {
		err := errors.New("can not attach products structure to mid or start category")
		return utils.CustomErrorResponse(409, 2, err, "", c)
	}

	if doesCategoriesHasProductStructures(b.CategoryKey) {
		err := errors.New("this category have a product structure")
		return utils.CustomErrorResponse(409, 3, err, "", c)
	}

	psCol := database.GetCollection("productStructures")
	meta, err := psCol.CreateDocument(context.Background(), b)
	if err != nil {
		return utils.CustomErrorResponse(404, 4041, err, "", c)
	}
	cs := addProductsStructureToCategory{ProductStructureKey: meta.Key}
	meta = database.UpdateDocument(b.CategoryKey, cs, "categories")

	menuList := createMenuFromProductStructure(b.ProductFieldList)
	createMenu := menu.CreateMenuDto{
		CategoryKey: b.CategoryKey,
		MenuItems:   menuList,
	}
	err = menu.CreateProductMenu(createMenu)
	if err != nil {
		return utils.CustomErrorResponse(404, 4041, err, "error creating menu from product structure", c)
	}
	return c.JSON(meta)
}

func updateProductStructAndMenu(c *fiber.Ctx)error  {
	b := new(updateIn)
	if err := utils.ParseBodyAndValidate(c, b); err != nil {
		return c.JSON(fmt.Sprintf("error parsing dto:%v", err))
	}

	psCol := database.GetCollection("productStructures")
	flag, err := psCol.DocumentExists(context.Background(), b.ProductStructureKey)
	if err != nil {
		return c.JSON(err)
	}
	if !flag {
		return c.Status(404).JSON("productStructures not found")
	}
	query := fmt.Sprintf("for i in productStructures filter i._key==\"%v\" update i with {productFieldList:PUSH(i.productFieldList,{name:\"%v\",isList:%v}, true)} in productStructures return NEW \n",b.ProductStructureKey,b.Name,b.IsList)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	
	var doc CreateProductStructureDto
	_, err = cursor.ReadDocument(ctx, &doc)

	if b.IsList {
		q:=fmt.Sprintf("for m in menu filter m.categoryKey==\"%v\" update m with {menuItems:PUSH(m.menuItems,{name:\"%v\",items:[]}, true)} in menu\n",doc.CategoryKey,b.Name)
		database.ExecuteGetQuery(q)
	}
	return c.JSON("ok")

}

func getAll(c *fiber.Ctx) error {
	q:=fmt.Sprintf("for i in productStructures \nfor j in categories \nfilter j._key==i.categoryKey return {ps:i,cat:j}")
	return c.JSON(database.ExecuteGetQuery(q))
}





















func getCategoryStatusByKey(key string) string {
	col := database.GetCollection("categories")
	var c Category
	meta, err := col.ReadDocument(context.Background(), key, &c)
	if err != nil {
		panic(fmt.Sprintf("error getting SupplyWorker by key %s", key))
	}
	fmt.Println(meta)
	return c.Status
}

func doesCategoriesHasProductStructures(key string) bool {
	db := database.GetDB()
	ctx := driver.WithQueryCount(context.Background())
	query := fmt.Sprintf("for i in productStructures filter i.categoryKey==\"%v\" return i", key)
	cursor, _ := db.Query(ctx, query, nil)

	defer cursor.Close()
	return cursor.Count() > 0
}

func createMenuFromProductStructure(pf []productField) []menu.Item {
	var items []menu.Item
	for _, v := range pf {
		if v.IsList {
			item := menu.Item{
				Name:  v.Name,
				Items: []string{},
			}
			items = append(items, item)
		}
	}
	return items
}
