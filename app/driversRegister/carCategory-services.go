package driversRegister

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

// createCarTypeCategory  create car category
// @Summary  create category for each car types
// @Description  category for each car types
// @Tags  carCategory
// @Accept  json
// @Produce  json
// @Param  carCategory  body  carCategory  true  "carCategory"
// @Success  200  {object}  carCategory{}
// @Failure  400  {object}  string
// @Router /car-category [post]
func createCarTypeCategory(c *fiber.Ctx) error {
	car := new(carCategory)
	ctx := context.Background()

	if err := utils.ParseBodyAndValidate(c, car); err != nil {
		return c.JSON(err)
	}
	carCol := database.GetCollection("carCat")

	carT := carCategory{
		CarTypes:   car.CarTypes,
		CarTonnage: car.CarTonnage,
		CarVolume:  car.CarVolume,
		CreatedAt:  time.Now().Unix(),
	}

	meta, err := carCol.CreateDocument(ctx, carT)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// getAllCarTypes get all of car categories
// @Summary return all types
// @Description return all categories
// @Tags carCategory
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Success 200 {object} carCategory
// @Failure 404 {object} string{}
// @Router /car-category [get]
func getAllCarTypes(c *fiber.Ctx) error {
	db := database.GetDB()
	ctx := context.Background()

	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("offset and limit must have value")
	}

	query := fmt.Sprintf("FOR c IN carCat LIMIT %v, %v RETURN c", offset, limit)

	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return c.JSON(err)
	}
	defer cursor.Close()
	var categoryList []getCategory
	for {
		var doc getCategory
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return c.JSON(err)
		}
		categoryList = append(categoryList, doc)
	}
	return c.JSON(categoryList)
}

// deleteCarCategory  delete car category
// @Summary  delete car category
// @Description  delete car category
// @Tags carCategory
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} carCategory
// @Failure 404 {object} string{}
// @Router /car-category/{key} [delete]
func deleteCarCategory(c *fiber.Ctx) error {
	key := c.Params("key")

	db := database.GetDB()
	ctx := context.Background()

	col, _ := db.Collection(ctx, "carCat")

	var doc getCategory
	_, err := col.RemoveDocument(ctx, key)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(doc)
}
