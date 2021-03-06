package userFav

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// removeFromUserFav add productId from fav field in user
// @Summary add productId from fav field in user
// @Description add productId from fav field in user
// @Tags  user-fav
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body uf true "data"
// @param Authorization header string true "Authorization"
// @Success 200 {object} users.UserOut{}
// @Failure 400 {object} string
// @Router /user-fav/add [post]
func addToUserFav(c *fiber.Ctx) error {
	lr := new(uf)
	if err := utils.ParseBodyAndValidate(c, lr); err != nil {
		return c.JSON(err)
	}
	db := database.GetDB()
	ctx := context.Background()
	if !strings.Contains(lr.ProductId, "/") {
		return c.Status(400).JSON(" you must send document Id like sheet/1234")
	}

	colName := strings.Split(lr.ProductId, "/")[0]
	key := strings.Split(lr.ProductId, "/")[1]

	col, err := db.Collection(ctx, colName)
	if err != nil {
		return c.Status(404).JSON("collection not found you must send document Id")
	}
	flag, err := col.DocumentExists(ctx, key)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	if !flag {
		return c.Status(404).JSON("document not found you must send document Id")
	}

	userKey := c.Locals("userKey").(string)
	q := fmt.Sprintf("for u in users\nfilter u._key==\"%v\"\nupdate u with {fav:PUSH(u.fav, \"%v\", true)} in users\nreturn NEW", userKey, lr.ProductId)
	res := database.ExecuteGetQuery(q)
	return c.JSON(res[0])

}

// removeFromUserFav remove productId from fav field in user
// @Summary remove productId from fav field in user
// @Description remove productId from fav field in user
// @Tags  user-fav
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body uf true "data"
// @param Authorization header string true "Authorization"
// @Success 200 {object} users.UserOut{}
// @Failure 400 {object} string
// @Router /user-fav/remove [post]
func removeFromUserFav(c *fiber.Ctx) error {
	lr := new(uf)
	if err := utils.ParseBodyAndValidate(c, lr); err != nil {
		return c.JSON(err)
	}
	userKey := c.Locals("userKey").(string)
	q := fmt.Sprintf("for u in users\nfilter u._key==\"%v\"\nupdate u with {fav:REMOVE_VALUE(u.fav, \"%v\")} in users\nreturn NEW", userKey, lr.ProductId)
	res := database.ExecuteGetQuery(q)
	return c.JSON(res[0])

}

// getUserFav get user fav
// @Summary get user fav
// @Description get user fav by jwt
// @Tags  user-fav
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []products.Product{}
// @Failure 400 {object} string
// @Router /user-fav [get]
func getUserFav(c *fiber.Ctx) error {

	offset := c.Query("offset")
	limit := c.Query("limit")

	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	userKey := c.Locals("userKey").(string)

	q := fmt.Sprintf("for i in users filter i._key==\"%v\" \nfor p in productSearch filter p._id in i.fav limit %v,%v return p", userKey, offset, limit)
	ql := fmt.Sprintf("let data=(for i in users filter i._key==\"%v\" \nfor p in productSearch filter p._id in i.fav return p ) return length(data)", userKey)
	res := database.ExecuteGetQuery(ql)
	return c.JSON(fiber.Map{
		"data": database.ExecuteGetQuery(q),
		"l":    res[0],
	})
}
