package faq

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

// createQuestion create questions
// @Summary create questions
// @Description create questions in faq database
// @Tags question
// @Accept json
// @Produce json
// @Param question body questionIn true "question"
// @Success 200 {object} question{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /faq [post]
func createQuestion(c *fiber.Ctx) error {
	q := new(questionIn)
	ctx := context.Background()

	if err := utils.ParseBodyAndValidate(c, q); err != nil {
		return c.JSON(err)
	}

	faqCatCol := database.GetCollection("faqCategory")

	var qCat category
	meta, err := faqCatCol.ReadDocument(context.Background(), q.CategoryKey, &qCat)
	if err != nil {
		log.Println(err)
		return c.Status(404).SendString(fmt.Sprintf("Category Not Founded  : %v", err))
	}

	qu := question{
		Title:         q.Title,
		Info:          q.Info,
		CategoryKey:   q.CategoryKey,
		IsPopular:     q.IsPopular,
		IsForSupplier: qCat.IsForSupplier,
		CreateAt:      time.Now().Unix(),
	}

	faqCol := database.GetCollection("faq")
	meta, err = faqCol.CreateDocument(ctx, qu)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// getAllQuestions get all questions
// @Summary return all questions
// @Description return all questions from faq
// @Tags question
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Param forSupplier  query bool   true  "forSupplier"
// @Param popular  query bool    true  "popular"
// @Success 200 {object} question
// @Failure 404 {object} string{}
// @Router /faq [get]
func getAllQuestions(c *fiber.Ctx) error {
	db := database.GetDB()
	ctx := context.Background()

	offset := c.Query("offset")
	limit := c.Query("limit")
	forSupplier := c.Query("forSupplier")
	popular := c.Query("popular")

	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	var query string
	if popular == "true" && forSupplier == "true" {
		query = fmt.Sprintf("for f in faq filter f.isForSupplier==true and f.isPopular==true  LIMIT %v, %v return f", offset, limit)
	} else if popular != "true" && forSupplier == "true" {
		query = fmt.Sprintf("for f in faq filter f.isForSupplier==true and f.isPopular==false  LIMIT %v, %v return f", offset, limit)
	} else if popular == "true" && forSupplier != "true" {
		query = fmt.Sprintf("for f in faq filter f.isForSupplier==false and f.isPopular==true  LIMIT %v, %v return f", offset, limit)
	} else {
		query = fmt.Sprintf("for f in faq filter f.isForSupplier==false and f.isPopular==false  LIMIT %v, %v return f", offset, limit)
	}
	log.Println(query)

	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return c.JSON(err)
	}
	defer cursor.Close()
	var questionList []getQuestion
	for {
		var doc getQuestion
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return c.JSON(err)
		}
		questionList = append(questionList, doc)
	}
	return c.JSON(questionList)

}

// getQuestionByKey get each question key
// @Summary return question by its key
// @Description return question by its key
// @Tags question
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} question
// @Failure 404 {object} string{}
// @Router /faq/{key} [get]
func getQuestionByKey(c *fiber.Ctx) error {
	key := c.Params("key")

	db := database.GetDB()
	ctx := context.Background()

	col, _ := db.Collection(ctx, "faq")

	var doc getQuestion
	_, err := col.ReadDocument(ctx, key, &doc)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(doc)

}

// updateFAQ update questions
// @Summary update questions
// @Description update questions
// @Tags question
// @Accept json
// @Produce json
// @Param question body question true "question"
// @Param key path string true "key"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /faq/{key} [put]
func updateFAQ(c *fiber.Ctx) error {
	key := c.Params("key")
	q := new(question)
	if err := utils.ParseBodyAndValidate(c, q); err != nil {
		return c.JSON(err)
	}
	faqCol := database.GetCollection("faq")
	meta, err := faqCol.UpdateDocument(context.Background(), key, q)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// deleteQuestion delete questions
// @Summary delete questions
// @Description delete questions
// @Tags question
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} feedBack
// @Failure 404 {object} string{}
// @Router /faq/{key} [delete]
func deleteQuestion(c *fiber.Ctx) error {
	key := c.Params("key")

	db := database.GetDB()
	ctx := context.Background()

	col, _ := db.Collection(ctx, "faq")

	var doc getQuestion
	_, err := col.RemoveDocument(ctx, key)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(doc)

}

// getFAQByCategory get each question by its category
// @Summary return question by its category
// @Description return question by its category key
// @Tags question
// @Accept json
// @Produce json
// @Param catKey path int true "cat Key"
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Success 200 {object} question{}
// @Failure 404 {object} string{}
// @Router /faq/cat/{catKey} [get]
func getFAQByCategory(c *fiber.Ctx) error {
	catKey := c.Params("catKey")
	faqCat := database.GetCollection("faqCategory")
	ctx := context.Background()

	offset := c.Query("offset")
	limit := c.Query("limit")

	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	var exist, _ = faqCat.DocumentExists(ctx, catKey)
	if !exist {
		return errors.New("category not founded")
	}
	k := fmt.Sprintf("for f in faq filter f.categoryKey==\"%v\" LIMIT %v, %v return f", catKey, offset, limit)
	//if forSupplier == "true" {
	//	k = fmt.Sprintf("for f in faq filter f.categoryKey==\"%v\" and f.isForSupplier==true  LIMIT %v, %v return f", catKey, offset, limit)
	//}
	log.Println(k)
	data := database.ExecuteGetQuery(k)
	return c.JSON(data)
}

// searchIntoQuestions search into questions
// @Summary return questions similar with given input
// @Description return questions similar with given input
// @Tags question
// @Accept json
// @Produce json
// @Param input  body  input true  "input"
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Success 200 {object} question{}
// @Failure 404 {object} string{}
// @Router /faq/srch [post]
func searchIntoQuestions(c *fiber.Ctx) error {
	i := new(input)

	// if err := c.BodyParser(i); err != nil {
	// 	return err
	// }
	if err := utils.ParseBodyAndValidate(c, i); err != nil {
		return c.JSON(err)
	}

	offset := c.Query("offset")
	limit := c.Query("limit")

	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	//TODO fix the like structure
	q := fmt.Sprintf("FOR f IN faq Filter f.title like \"%v\" LIMIT %v, %v return f.title", "%"+i.Title+"%", offset, limit)
	data := database.ExecuteGetQuery(q)
	return c.JSON(data)

}

// createFeedBackCollection create new feedback
// @Summary create feedback
// @Description create feedback
// @Tags feedBack
// @Accept json
// @Produce json
// @Param feedBack body feedBack true "feed back"
// @Success 200 {object} feedBack
// @Failure 404 {object} string{}
// @Router /faq/FeedBack [post]
func createFeedBackCollection(c *fiber.Ctx) error {
	a := new(feedBack)
	if err := utils.ParseBodyAndValidate(c, a); err != nil {
		return c.JSON(err)
	}

	faqCat := database.GetCollection("faqComments")

	ctx := context.Background()
	meta, err := faqCat.CreateDocument(ctx, a)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// getFeedBackByKey get feedback by key
// @Summary return feedbacks by its key
// @Description return feedbacks by its key
// @Tags feedBack
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} feedBack
// @Failure 404 {object} string{}
// @Router /faq/FeedBack/{key} [get]
func getFeedBackByKey(c *fiber.Ctx) error {
	colKey := c.Params("key")
	db := database.GetDB()
	ctx := context.Background()
	col, _ := db.Collection(ctx, "faqComments")

	var doc getFeedBack
	_, err := col.ReadDocument(ctx, colKey, &doc)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(doc)
}
