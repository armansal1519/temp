package brands

import (
	"bamachoub-backend-go-v1/utils"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/brands")

	r.Get("/", getAllBrands)
	r.Get("/:key", GetBrandByKey)
	r.Get("/url/:categoryurl", getAllBrandsByCategoryUrl)

	r.Get("/used/:dbName/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")
		dbName := c.Params("dbName")

		res, err := getBrandsUsedUnderCategoryByCategoryKey(key, dbName)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(res)
	})

	r.Post("/category", getBrandsByCategoryName)

	r.Post("", func(c *fiber.Ctx) error {
		b := new(brandDto)
		if err := utils.ParseBodyAndValidate(c, b); err != nil {
			return c.JSON(err)
		}
		res, err := CreateBrand(*b)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(res)
	})

	r.Put("/:key", updateBrand)
	r.Delete("/:key", removeBrand)

}

type cat struct {
	Cat string `json:"cat"`
}
