package contactUs

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

// create create contact us
// @Summary create contact us
// @Description create contact us
// @Tags contact us
// @Accept json
// @Produce json
// @Param contactIn body contactIn true "contact"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /contact-us  [post]
func create(c *fiber.Ctx) error {
	ci := new(contactIn)
	if err := utils.ParseBodyAndValidate(c, ci); err != nil {
		return c.JSON(err)
	}
	contactCol := database.GetCollection("contactUs")
	contact := contact{
		Title:       ci.Title,
		FullName:    ci.FullName,
		PhoneNumber: ci.PhoneNumber,
		Text:        ci.Text,
		ImageArr:    ci.ImageArr,
		CreatedAt:   time.Now().Unix(),
	}
	meta, err := contactCol.CreateDocument(context.Background(), contact)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)

}

// getAll get all contact
// @Summary return all contact
// @Description return all contact
// @Tags contact us
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Success 200 {object} []contactOut{}
// @Failure 404 {object} string{}
// @Router /contact-us [get]
func getAll(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	query := fmt.Sprintf("for i in contactUs limit %v,%v return i", offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))

}

// getByKey get one contact
// @Summary return one contact
// @Description return one contact
// @Tags contact us
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} contactOut{}
// @Failure 404 {object} string{}
// @Router /contact-us/{key} [get]
func getByKey(c *fiber.Ctx) error {
	key := c.Params("key")
	contactCol := database.GetCollection("contactUs")
	var contact contactOut
	_, err := contactCol.ReadDocument(context.Background(), key, &contact)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(contact)
}

// update create contact us
// @Summary update contact us
// @Description update contact us
// @Tags contact us
// @Accept json
// @Produce json
// @Param contactIn body contactIn true "contact"
// @Param key path string true "key"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /contact-us/{key}  [put]
func update(c *fiber.Ctx) error {
	ci := new(contactIn)
	if err := utils.ParseBodyAndValidate(c, ci); err != nil {
		return c.JSON(err)
	}
	key := c.Params("key")
	contactCol := database.GetCollection("contactUs")
	meta, err := contactCol.UpdateDocument(context.Background(), key, ci)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// delete get one contact
// @Summary delete one contact
// @Description delete one contact
// @Tags contact us
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} contactOut{}
// @Failure 404 {object} string{}
// @Router /contact-us/{key} [delete]
func delete(c *fiber.Ctx) error {
	key := c.Params("key")
	contactCol := database.GetCollection("contactUs")
	_, err := contactCol.RemoveDocument(context.Background(), key)
	if err != nil {
		return c.JSON(err)
	}
	return c.Status(204).SendString("deleted")
}
