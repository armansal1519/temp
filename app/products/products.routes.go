package products

import (
	"bamachoub-backend-go-v1/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Routes(app fiber.Router) {
	r := app.Group("/products")
	r.Get("Product-body/:key", GetProductBodyByCategoryKey)
	//r.Get("Product-base/:key",getBaseProduct)
	r.Get("/:productName", GetProductByUrl)
	r.Post("/body", createProductBody)
	r.Post("/color/:spId/:productKey", getProductWithColorCode)

	r.Post("/", func(c *fiber.Ctx) error {
		p := new(productInfo)

		if err := utils.ParseBodyAndValidate(c, p); err != nil {
			return c.JSON(err)
		}
		res, err := createTheFuckingProduct(*p)
		if err != nil {
			c.JSON(err)
		}
		return c.JSON(res)
	})

	r.Get("/all/:dbName", func(c *fiber.Ctx) error {
		dbName := c.Params("dbName")
		s := strings.Split(dbName, ",")
		data, err := getProducts(s)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(data)
	})

	r.Get("/cat/:dbName/:key", func(c *fiber.Ctx) error {
		dbName := c.Params("dbName")
		key := c.Params("key")
		offset := c.Query("offset")
		limit := c.Query("limit")
		if offset == "" || limit == "" {
			return c.Status(400).JSON("offset or limit query can not be empty")
		}
		data, err := getProductFromCategory(dbName, key, offset, limit)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(data)
	})

	r.Get("/one/:dbName/:key", func(c *fiber.Ctx) error {
		dbName := c.Params("dbName")
		key := c.Params("key")
		data, err := getProductByKey(dbName, key)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(data)
	})

	r.Put("/:categoryUrl/:key", updateProduct)
	r.Delete("/:categoryUrl/:key", deleteProduct)

	r.Post("/basic-search/:dbName", basicSearchProducts)
	r.Post("/basic-filter/:dbName", basicFilter)
	r.Post("/advance-filter/:dbName/:categoryKey", AdvanceFilter)

}
