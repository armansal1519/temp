package discountCode

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func createDiscountForPhoneNumbers(c *fiber.Ctx) error {
	req := new(DiscountForPhoneNumbersRequest)
	if err := utils.ParseBodyAndValidate(c, req); err != nil {
		return c.JSON(err)
	}
	discountCol := database.GetCollection("discount")
	d := Discount{
		Type:   req.Type,
		Amount: req.Amount,
		EndAt:  req.EndAt,
	}
	meta, err := discountCol.CreateDocument(context.Background(), d)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	phoneNumberStr := " [ "
	for i, number := range req.PhoneNumbers {
		phoneNumberStr += fmt.Sprintf("\"%v\"", number)
		if i < len(req.PhoneNumbers)-1 {
			phoneNumberStr += " , "
		}
	}
	phoneNumberStr += " ] "

	q := fmt.Sprintf("for u in users filter u.phoneNumber in  %v  insert {_from:u._id,_to:\"%v\",isUsed:false} into discountEdge \n", phoneNumberStr, meta.ID.String())
	fmt.Println(q)
	database.ExecuteGetQuery(q)
	return c.JSON("ok")

}

// getDiscountByKey get discount for by key
// @Summary get discount for by key
// @Description get discount for by key
// @Tags discount
// @Accept json
// @Produce json
// @Param   key      path   string     true  "discount key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} DiscountOut{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /discount/{key} [get]
func getDiscountByKey(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	discountKey := c.Params("key")
	q := fmt.Sprintf("for d in discountEdge filter d._from==\"users/%v\" and d._to==\"discount/%v\" and d.isUsed==false\nfor i in discount filter i._id==d._to\nreturn i", userKey, discountKey)
	res := database.ExecuteGetQuery(q)
	return c.JSON(res[0])

}
