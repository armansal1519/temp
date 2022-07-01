package addBuyMethod

import (
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/middleware"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/add-buy-method")

	r.Post("/get-price", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {

		offset := c.Query("offset")
		limit := c.Query("limit")
		if offset == "" || limit == "" {
			return c.Status(400).SendString("offset or query must have value")
		}
		supplierId := c.Locals("supplierId").(string)
		b := new(brandFilter)
		fmt.Println(111, string(c.Body()))
		if err := utils.ParseBodyAndValidate(c, b); err != nil {
			return c.JSON(err)
		}

		resp, err := getAllPricesWithProductsBySupplierKey(supplierId, b.Brand, b.Search, offset, limit)

		if err != nil {
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(resp)
	})

	r.Get("/price/brand", middleware.GetSupplierByEmployee, getPriceBrandsBySupplierKey)

	r.Get("/price/:categoryUrl", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		offset := c.Query("offset")
		limit := c.Query("limit")
		if offset == "" || limit == "" {
			return c.Status(400).SendString("offset or query must have value")
		}
		categoryUrl := c.Params("categoryUrl")
		supplierId := c.Locals("supplierId").(string)
		// supplierKey := strings.Split(supplierId, "/")[1]

		resp, l, err := getPriceWithProductBySupplierKey(categoryUrl, supplierId, offset, limit)

		if err != nil {
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(fiber.Map{
			"data":   resp,
			"length": l,
		})
	})

	r.Post("/price", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		pgi := new(PriceIn)
		supplierId := c.Locals("supplierId").(string)
		if err := utils.ParseBodyAndValidate(c, pgi); err != nil {
			return c.JSON(err)
		}
		resp, err := AddPriceToProduct(*pgi, supplierId)

		if err != nil {
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(resp)
	})

	r.Put("/price/group_update/:priceCol", func(c *fiber.Ctx) error {
		gue := new(groupUpdatePriceIn)
		priceCol := c.Params("priceCol")
		if err := utils.ParseBodyAndValidate(c, gue); err != nil {
			log.Println(err)
			return c.JSON(err)
		}
		resp, err := groupUpdatePrice(*gue, priceCol)

		if err != nil {
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(resp)
	})

	r.Put("/price/:priceCol/:priceKey", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		pgi := new(updatePrice)
		priceCol := c.Params("priceCol")
		priceKey := c.Params("priceKey")
		//supplierId := c.Locals("supplierId").(string)
		if err := utils.ParseBodyAndValidate(c, pgi); err != nil {
			log.Println(err)
			return c.JSON(err)
		}
		resp, err := updatePriceOfProduct(*pgi, priceCol, priceKey)

		if err != nil {
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(resp)
	})

	r.Delete("/price/:col/:key",
		//middleware.GetSupplierByEmployee,
		func(c *fiber.Ctx) error {
			productKey := c.Params("key")
			productCol := c.Params("col")
			err := deletePrice(productKey, productCol)
			if err != nil {
				return c.JSON(fmt.Sprintf("%v", err))
			}
			return c.Status(204).SendString("document deleted")
		})

	//******************************************************************************************************************
	//estelam
	//******************************************************************************************************************

	r.Post("/get-estelam", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {

		offset := c.Query("offset")
		limit := c.Query("limit")
		if offset == "" || limit == "" {
			return c.Status(400).SendString("offset or query must have value")
		}
		supplierId := c.Locals("supplierId").(string)
		b := new(brandFilter)
		fmt.Println(111, string(c.Body()))
		if err := utils.ParseBodyAndValidate(c, b); err != nil {
			return c.JSON(err)
		}

		resp, l, err := getAllEstelamsWithProductsBySupplierKey(supplierId, b.Brand, b.Search, offset, limit)

		if err != nil {
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(fiber.Map{
			"data":   resp,
			"length": l,
		})
	})

	r.Get("/estelam/brand", middleware.GetSupplierByEmployee, getEstelamBrandsBySupplierKey)

	r.Get("/estelam/:categoryUrl", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		offset := c.Query("offset")
		limit := c.Query("limit")
		if offset == "" || limit == "" {
			return c.Status(400).SendString("offset or query must have value")
		}
		categoryUrl := c.Params("categoryUrl")
		supplierId := c.Locals("supplierId").(string)
		// fmt.Println("33333333333333", supplierId)
		// supplierKey := strings.Split(supplierId, "/")[1]

		resp, l, err := getEstelamWithProductBySupplierKey(categoryUrl, supplierId, offset, limit)

		if err != nil {
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(fiber.Map{
			"data":   resp,
			"length": l,
		})
	})

	r.Post("/estelam", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		est := new(estelamIn)
		log.Println(3)
		supplierId := c.Locals("supplierId").(string)
		if err := utils.ParseBodyAndValidate(c, est); err != nil {
			log.Println(1, err)
			return c.JSON(err)
		}
		fmt.Println(string(c.Body()))
		if est.ProductId == "" {
			log.Println("product id" + est.ProductId)
			return c.Status(400).JSON("productId cant be empty")
		}
		resp, err := AddEstelamToProduct(*est, supplierId)

		if err != nil {
			log.Println(2)
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(resp)
	})

	r.Put("/estelam/group_update/:priceCol", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		gue := new(groupUpdateEstelamIn)
		priceCol := c.Params("priceCol")
		if err := utils.ParseBodyAndValidate(c, gue); err != nil {
			log.Println(err)
			return c.JSON(err)
		}
		resp, err := groupUpdateEstelam(*gue, priceCol)

		if err != nil {
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(resp)
	})
	r.Put("/estelam/:priceCol/:priceKey", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		pgi := new(updateEstelam)
		priceCol := c.Params("priceCol")
		priceKey := c.Params("priceKey")
		//supplierId := c.Locals("supplierId").(string)
		log.Println(priceCol, priceKey)
		if err := utils.ParseBodyAndValidate(c, pgi); err != nil {
			log.Println(err)
			return c.JSON(err)
		}
		resp, err := updateEstelamOfProduct(*pgi, priceCol, priceKey)

		if err != nil {
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(resp)
	})

	r.Delete("/estelam/:col/:key", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		productKey := c.Params("key")
		productCol := c.Params("col")
		err := deleteEstelam(productKey, productCol)
		if err != nil {
			return c.JSON(fmt.Sprintf("%v", err))
		}
		return c.Status(204).SendString("document deleted")
	})

}
