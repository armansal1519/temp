package productQA

import (
	"bamachoub-backend-go-v1/app/users"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

// getQAByProductKey get each question by its product key
// @Summary return question by its product key
// @Description return question by its product key
// @Tags productQA
// @Accept json
// @Produce json
// @Param categoryUrl path string true "categoryUrl"
// @Param productKey path string true "productKey"
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Param sort  query string    true  "sort"
// @Success 200 {object} productQA{}
// @Failure 404 {object} string{}
// @Router /products-q-a/{categoryUrl}/{productKey} [get]
func getQAByProductKey(c *fiber.Ctx) error {
	categoryUrl := c.Params("categoryUrl")
	productKey := c.Params("productKey")
	sort := c.Query("sort")
	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	if sort != "likes" && sort != "new" {
		return c.Status(400).SendString("sort value is not acceptable only likes or new")
	}
	if sort == "new" {
		sort = "i.createdAt desc"
	}
	if sort == "likes" {
		sort = "length(i.likes ) desc"
	}

	query := fmt.Sprintf("let data=(for i in productQA filter i.productId==\"%v/%v\" filter i.status==\"valid\" filter i.questionKey==\"\"  sort %v return i) \nlet docs=(for j in data  limit %v,%v \nlet q= (for k in productQA filter j._key==k.questionKey  && k.status==\"valid\"  return k)\nreturn {question:j,answers:q})\nreturn {len:LENGTH(data),data:docs}", categoryUrl, productKey, sort, offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))

}

// createQA create products questions
// @Summary create questions
// @Description create questions
// @Tags productQA
// @Accept json
// @Produce json
// @Param question body productQA true "question"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /products-q-a [post]
func createQA(c *fiber.Ctx) error {
	pqa := new(productQA)
	if err := utils.ParseBodyAndValidate(c, pqa); err != nil {
		return c.JSON(err)
	}
	userKey := c.Locals("userKey").(string)
	userCol := database.GetCollection("users")
	var u users.UserOut
	_, err := userCol.ReadDocument(context.Background(), userKey, &u)

	pqa.FullName = u.FirstName + " " + u.LastName
	pqa.UserKey = userKey
	pqa.CreatedAt = time.Now().Unix()
	pqa.Status = "wait"

	col := database.GetCollection("productQA")
	meta, err := col.CreateDocument(context.Background(), pqa)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// updateQA update products questions
// @Summary update questions
// @Description update questions
// @Tags productQA
// @Accept json
// @Produce json
// @Param question body updateDto true "question"
// @Param key path string true "key"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /products-q-a/{key} [put]
func updateQA(c *fiber.Ctx) error {
	key := c.Params("key")
	uqa := new(updateDto)
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

// getAll return all questions
// @Summary return all questions
// @Description return all questions
// @Tags productQA
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Param valid  query string    true  "valid or not"
// @Success 200 {object} productQA{}
// @Failure 404 {object} string{}
// @Router /products-q-a [get]
func getAll(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	valid := c.Params("valid")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	if valid == "not" {
		valid = "not"
	} else if valid == "wait" {
		valid = "wait"
	} else {
		valid = "valid"
	}
	query := fmt.Sprintf("for i in productQA filter i.status==\"%v\" sort i.createdAt limit %v,%v return i", valid, offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}

// createQA create products questions
// @Summary create questions
// @Description create questions
// @Tags productQA
// @Accept json
// @Produce json
// @Param op path string true "operation only add or remove"
// @Param questionKey path string true "questionKey"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /products-q-a/likes/{op}/{questionKey} [post]
func likes(c *fiber.Ctx) error {
	op := c.Params("op")
	qaKey := c.Params("qaKey")
	userKey := c.Locals("userKey").(string)
	if op == "add" {
		query := fmt.Sprintf("for i in productQA\nfilter i._key==\"%v\"\nupdate i with {likes:append(i.likes,\"%v\",true)} in productQA", qaKey, userKey)
		database.ExecuteGetQuery(query)
		return c.SendString("added")
	} else if op == "remove" {
		query := fmt.Sprintf("for i in productQA\nfilter i._key==\"%v\"\nupdate i with {likes:REMOVE_VALUE(i.likes,\"%v\")} in productQA", qaKey, userKey)
		database.ExecuteGetQuery(query)
		return c.SendString("removed")
	}
	return c.Status(400).SendString("op can only be add or remove")
}

// getQAForUser return all questions or answers for user
// @Summary return all questions or answers for user
// @Description return all questions or answers for user
// @Tags productQA
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Param qa path string true "qa only q or a"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} productQA{}
// @Failure 404 {object} string{}
// @Router /products-q-a/user/{qa} [get]
func getQAForUser(c *fiber.Ctx) error {
	qa := c.Params("qa")
	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	userKey := c.Locals("userKey").(string)
	if qa == "a" {
		query := fmt.Sprintf("let data=(for i in productQA filter i.userKey==\"%v\" filter i.questionKey!=\"\"  sort i.createdAt desc return i) \nlet docs=(for j in data  limit %v,%v return j)\nreturn {len:LENGTH(data),data:docs}", userKey, offset, limit)
		return c.JSON(database.ExecuteGetQuery(query))
	}
	if qa == "q" {
		query := fmt.Sprintf("let data=(for i in productQA filter i.userKey==\"%v\" filter i.questionKey==\"\"  sort i.createdAt desc return i) \nlet docs=(for j in data  limit %v,%v return j)\nreturn {len:LENGTH(data),data:docs}", userKey, offset, limit)
		return c.JSON(database.ExecuteGetQuery(query))
	}
	return c.Status(400).SendString("op must be q or a")

}

// removeQA remove QA
// @Summary remove QA
// @Description remove QA
// @Tags productQA
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /products-q-a/{key} [delete]
func removeQA(c *fiber.Ctx) error {
	key := c.Params("key")
	col := database.GetCollection("productQA")
	meta, err := col.RemoveDocument(context.Background(), key)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}
