package commentOnProduct

import (
	"bamachoub-backend-go-v1/app/users"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

// createCommit create product comment
// @Summary create product comment
// @Description create product comment
// @Tags product comment
// @Accept json
// @Produce json
// @Param comment body comment true "comment"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /product-comment [post]
func createCommit(c *fiber.Ctx) error {
	ci := new(comment)
	if err := utils.ParseBodyAndValidate(c, ci); err != nil {
		return c.JSON(err)
	}
	userKey := c.Locals("userKey").(string)
	u, err := users.GetUserByKey(userKey)
	if err != nil {
		return c.JSON(err)
	}
	isBuyer := true
	_, err = getApprovedOrder(u.Key, ci.ProductId)
	if err != nil {
		isBuyer = false
	}
	ci.UserKey = userKey
	ci.IsBuyer = isBuyer
	ci.CreatedAt = time.Now().Unix()
	ci.UserFullName = u.FirstName + " " + u.LastName

	pcCol := database.GetCollection("productComment")
	meta, err := pcCol.CreateDocument(context.Background(), ci)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// getProductComment get all comment for a product
// @Summary return all comment for a product
// @Description return all comment for a product from faq
// @Tags product comment
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Param sort  query string    true  "sort"
// @Param score  query bool     true  "score"
// @Param categoryUrl path string true "categoryUrl"
// @Param productKey path string true "productKey"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /product-comment/{categoryUrl}/{productKey} [get]
func getProductComment(c *fiber.Ctx) error {
	categoryUrl := c.Params("categoryUrl")
	productKey := c.Params("productKey")
	offset := c.Query("offset")
	limit := c.Query("limit")
	sort := c.Query("sort")
	score := c.Query("score")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	sortStr := ""
	if sort == "new" {
		sortStr = "sort k.createdAt desc"
	}
	if sort == "likes" {
		sortStr = "sort length(k.likes ) desc"
	}
	q := fmt.Sprintf("let raw=(for i in productComment  filter i.productId==\"%v/%v\" return i)\nlet r=(for k in raw %v limit %v,%v return k) ", categoryUrl, productKey, sortStr, offset, limit)
	q1 := fmt.Sprintf("let raw=(for i in productComment  filter i.productId==\"%v/%v\" return i) return length(raw)", categoryUrl, productKey)

	if score == "true" {
		q += " let s=(let data=FLATTEN(for i in raw return i.scoreArr)\nfor i in data\ncollect names= i.title into g\nlet scoreArr= (for j in g[*].i return j.score)  return{names:names,scoreArr:AVERAGE(scoreArr)})\nreturn {score:s,comments:r} "
	} else {
		q += " return r "
	}
	log.Println(q)
	res := database.ExecuteGetQuery(q1)
	return c.JSON(fiber.Map{"data": database.ExecuteGetQuery(q), "length": res[0]})
}

// getProductComment get all the images from comment for a product
// @Summary return all the images from comment for a product
// @Description return all the images from comment for a product from faq
// @Tags product comment
// @Accept json
// @Produce json
// @Param categoryUrl path string true "categoryUrl"
// @Param productKey path string true "productKey"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /product-comment/images/{categoryUrl}/{productKey} [get]
func getImagesFromProductComment(c *fiber.Ctx) error {
	categoryUrl := c.Params("categoryUrl")
	productKey := c.Params("productKey")

	q := fmt.Sprintf("let data=(for i in productComment  filter i.productId==\"%v/%v\" return i.imageUrls)\n\nreturn UNIQUE(FLATTEN(data)) ", categoryUrl, productKey)

	return c.JSON(database.ExecuteGetQuery(q))
}

// getAll get all comment
// @Summary return all comment
// @Description return all comment
// @Tags product comment
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Param sort  query string  true  "sort"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /product-comment [get]
func getAll(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	sort := c.Query("sort")

	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	sortStr := ""
	if sort == "new" {
		sortStr = "sort i.createdAt desc"
	}
	if sort == "likes" {
		sortStr = "sort length(i.likes ) desc"
	}
	query := fmt.Sprintf("for i in productComment  %v limit %v,%v return i", sortStr, offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}

// getAll get all comment for one user
// @Summary return all comment for one user
// @Description return all comment for one user
// @Tags product comment
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /product-comment/user [get]
func getByUserKey(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	userKey := c.Locals("userKey").(string)

	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	q := fmt.Sprintf("for i in productComment filter i.userKey==\"%v\" sort i.createdAt desc limit %v,%v return i", userKey, offset, limit)
	ql := fmt.Sprintf("let data=(for i in productComment filter i.userKey==\"%v\" return i) return length(data)", userKey)
	res := database.ExecuteGetQuery(ql)
	return c.JSON(fiber.Map{
		"data":   database.ExecuteGetQuery(q),
		"length": res[0],
	})

}

// updateComment update comment
// @Summary update comment
// @Description update comment
// @Tags product comment
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Param updateCommentType body updateCommentType true "updateCommentType"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /product-comment/{key} [put]
func updateComment(c *fiber.Ctx) error {
	cu := new(updateCommentType)
	if err := utils.ParseBodyAndValidate(c, cu); err != nil {
		return c.JSON(err)
	}
	key := c.Params("key")

	pcCol := database.GetCollection("productComment")
	meta, err := pcCol.UpdateDocument(context.Background(), key, cu)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)

}

// adminUpdateComment update comment
// @Summary update comment
// @Description update comment
// @Tags product comment
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Param adminUpdateCommentType body adminUpdateCommentType true "adminUpdateCommentType"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /product-comment/admin/{key} [put]
func adminUpdateComment(c *fiber.Ctx) error {
	cu := new(adminUpdateCommentType)
	if err := utils.ParseBodyAndValidate(c, cu); err != nil {
		return c.JSON(err)
	}
	key := c.Params("key")
	pcCol := database.GetCollection("productComment")
	meta, err := pcCol.UpdateDocument(context.Background(), key, cu)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)

}

// deleteComment delete comment
// @Summary delete comment
// @Description delete comment
// @Tags product comment
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /product-comment/{key} [delete]
func deleteComment(c *fiber.Ctx) error {
	key := c.Params("key")
	var comment commentOut
	pcCol := database.GetCollection("productComment")
	_, err := pcCol.ReadDocument(context.Background(), key, &comment)
	if err != nil {
		return c.JSON(err)
	}
	isAdmin := c.Locals("isAdmin").(bool)

	if !isAdmin {
		userKey := c.Locals("userKey").(string)
		if userKey != comment.UserKey {
			return c.Status(403).SendString("Unauthorized")
		}
	}

	meta, err := pcCol.RemoveDocument(context.Background(), key)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}
