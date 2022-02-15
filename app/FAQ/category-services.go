package faq

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

// createCategory create new category
// @Summary create category
// @Description create category
// @Tags faq category
// @Accept json
// @Produce json
// @Param category body category true "category"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /faq-category [post]
func createCategory(c *fiber.Ctx) error {
	cat := new(category)
	ctx := context.Background()
	if err := utils.ParseBodyAndValidate(c, cat); err != nil {
		return c.JSON(err)
	}

	faqCat := database.GetCollection("faqCategory")

	cat.CreatedAt = time.Now().Unix()

	meta, err := faqCat.CreateDocument(ctx, cat)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// GetCategory get category by key
// @Summary return each category by its key
// @Description return category
// @Tags faq category
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Param forSupplier  query bool    true  "forSupplier"
// @Success 200 {object} []category{}
// @Failure 404 {object} string{}
// @Router /faq-category [get]
func GetCategory(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	forSupplier := c.Query("forSupplier")
	if forSupplier != "true" {
		forSupplier = "false"
	}

	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	q := fmt.Sprintf("for f in faqCategory filter f.isForSupplier==%v sort f.createdAt LIMIT %v, %v return f", forSupplier, offset, limit)

	resp := database.ExecuteGetQuery(q)
	return c.JSON(resp)
}

// deleteCategory delete category
// @Summary delete category by its key
// @Description delete category
// @Tags faq category
// @Accept json
// @Produce json
// @Param catKey path string true "cat key"
// @Success 200 {object} category
// @Failure 404 {object} string{}
// @Router /faq-category/{catKey} [delete]
func deleteCategory(c *fiber.Ctx) error {
	catKey := c.Params("catKey")
	db := database.GetDB()
	ctx := context.Background()

	col, _ := db.Collection(ctx, "faqCategory")

	var doc getCategory
	_, err := col.RemoveDocument(ctx, catKey)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(&doc)
}
