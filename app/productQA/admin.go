package productQA

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"github.com/gofiber/fiber/v2"
)

// adminUpdateQA update products questions
// @Summary update questions
// @Description update questions
// @Tags productQA
// @Accept json
// @Produce json
// @Param question body adminUpdateDto true "question"
// @Param key path string true "key"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /products-q-a/admin/{key} [put]
func adminUpdateQA(c *fiber.Ctx) error {
	key := c.Params("key")
	uqa := new(adminUpdateDto)
	if err := utils.ParseBodyAndValidate(c, uqa); err != nil {
		return c.JSON(err)
	}
	uqa.Status = "wait"
	col := database.GetCollection("productQA")
	meta, err := col.UpdateDocument(context.Background(), key, uqa)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

type adminUpdateDto struct {
	Text          string `json:"text"`
	RejectionText string `json:"rejectionText"`
	Status        string `json:"status"`
}
