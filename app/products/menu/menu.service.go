package menu

import (
	"bamachoub-backend-go-v1/app/categories"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

func CreateProductMenu(b CreateMenuDto) error {
	//b := new(CreateMenuDto)
	//
	//if err := utils.ParseBodyAndValidate(c, b); err != nil {
	//	panic(fmt.Sprintf("error parsing dto:%v", err))
	//}

	cCol := database.GetCollection("categories")

	doesExist, err := cCol.DocumentExists(context.Background(), b.CategoryKey)
	if err != nil || !doesExist {
		return err
	}

	menuCol := database.GetCollection("menu")
	meta, err := menuCol.CreateDocument(context.Background(), b)
	if err != nil {
		return err
	}
	m := addMenuKeyToCategory{MenuKey: meta.Key}
	meta = database.UpdateDocument(b.CategoryKey, m, "categories")

	return nil
}

//func getMenuByCategoryKey(c *fiber.Ctx) error {
//	key := c.Params("key")
//	query := fmt.Sprintf("for i in categories\nfilter i._key==\"%v\"\nfor j in menu\nfilter i.menuKey==j._key\nreturn j\n", key)
//	data := database.ExecuteGetQuery(query)
//	return c.JSON(data[0])
//}
func getMenuByCategoryKey(c *fiber.Ctx) error {
	key := c.Params("key")
	var cat categories.CategoryOut
	categoryCol := database.GetCollection("categories")
	_, err := categoryCol.ReadDocument(context.Background(), key, &cat)
	if err != nil {
		return c.JSON(err)
	}
	if cat.Status == "end" {
		q := fmt.Sprintf("for i in menu filter i.categoryKey==\"%v\" return i", key)
		res := database.ExecuteGetQuery(q)
		return c.JSON(res[0])
	}

	query := fmt.Sprintf("let ck=(for c in categories filter c._key==\"%v\"  for v in 0..10 outbound c graph \"categoryGraph\" filter v.status==\"end\" return v._key) \nfor i in menu filter i.categoryKey in ck\nreturn i", key)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []ReturnMenu
	for {
		var doc ReturnMenu
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	union := getUnion(data)
	final := order(data, union)
	return c.JSON(fiber.Map{
		"categoryKey": key,
		"menuItems":   final,
	})
}

func addToMenu(c *fiber.Ctx) error {
	key := c.Params("key")

	b := new(CreateMenuDto)

	if err := utils.ParseBodyAndValidate(c, b); err != nil {
		panic(fmt.Sprintf("error parsing dto:%v", err))
	}
	menuCol := database.GetCollection("menu")
	var m ReturnMenu
	ctx := driver.WithReturnNew(context.Background(), &m)
	fmt.Println(b)
	_, err := menuCol.UpdateDocument(ctx, key, b)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(m)
}
