package productSuggestion

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

// AddNewSuggestion add product suggestion
// @Summary add product suggestion
// @Description add product suggestion
// @Tags product suggestion
// @Accept json
// @Produce json
// @Param data body productSuggestion true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} productSuggestionOut{}
// @Failure 400 {object} string
// @Router /product-suggestion [post]
func AddNewSuggestion(c *fiber.Ctx) error {
	ps := new(productSuggestion)

	if err := utils.ParseBodyAndValidate(c, ps); err != nil {
		return c.JSON(err)
	}
	psCol := database.GetCollection("productSuggestion")
	key := c.Locals("supplierId").(string)
	ps.SupplierKey = key
	var psOut productSuggestionOut
	ctx := driver.WithReturnNew(context.Background(), &psOut)
	_, err := psCol.CreateDocument(ctx, ps)
	if err != nil {
		return c.JSON(err)

	}
	return c.JSON(psOut)
}

// AddFromSample add product suggestion from sample
// @Summary add product suggestion from sample
// @Description add product suggestion from sample
// @Tags product suggestion
// @Accept json
// @Produce json
// @Param data body sampleSuggestion true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} sampleSuggestionOut{}
// @Failure 400 {object} string
// @Router /product-suggestion/sample [post]
func AddFromSample(c *fiber.Ctx) error {
	ps := new(sampleSuggestion)

	if err := utils.ParseBodyAndValidate(c, ps); err != nil {
		return c.JSON(err)
	}
	psCol := database.GetCollection("productSuggestion")
	userKey := c.Locals("userKey").(string)
	ps.UserKey = userKey
	var psOut sampleSuggestionOut
	ctx := driver.WithReturnNew(context.Background(), &psOut)
	_, err := psCol.CreateDocument(ctx, ps)
	if err != nil {
		return c.JSON(err)

	}
	return c.JSON(psOut)
}

// getProductSuggestions get product suggestion
// @Summary get product suggestion
// @Description get product suggestion
// @Tags product suggestion
// @Accept json
// @Produce json
// @Success 200 {object} productSuggestionOut{}
// @Failure 400 {object} string
// @Router /product-suggestion [get]
func getProductSuggestions(c *fiber.Ctx) error {
	q := fmt.Sprintf("for i in productSuggestion return i")
	return c.JSON(database.ExecuteGetQuery(q))
}

// deleteProductSuggestions delete product suggestion
// @Summary delete product suggestion
// @Description delete product suggestion
// @Tags product suggestion
// @Accept json
// @Produce json
// @Param   key      path   string     true  "key"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /product-suggestion/{key} [delete]
func deleteProductSuggestions(c *fiber.Ctx) error {
	key := c.Params("key")
	psCol := database.GetCollection("productSuggestion")
	meta, err := psCol.RemoveDocument(context.Background(), key)
	if err != nil {
		return c.JSON(err)

	}
	return c.JSON(meta)

}

// AddBetterPrice add product suggestion with better price
// @Summary add product suggestion with better price
// @Description add product suggestion with better price
// @Tags product suggestion
// @Accept json
// @Produce json
// @Param data body betterPrice true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} betterPriceOut{}
// @Failure 400 {object} string
// @Router /product-suggestion/better-price [post]
func AddBetterPrice(c *fiber.Ctx) error {
	ps := new(betterPrice)

	if err := utils.ParseBodyAndValidate(c, ps); err != nil {
		return c.JSON(err)
	}
	psCol := database.GetCollection("betterPriceSuggestion")
	userKey := c.Locals("userKey").(string)
	ps.UserKey = userKey
	var psOut betterPriceOut
	ctx := driver.WithReturnNew(context.Background(), &psOut)
	_, err := psCol.CreateDocument(ctx, ps)
	if err != nil {
		return c.JSON(err)

	}
	return c.JSON(psOut)
}
