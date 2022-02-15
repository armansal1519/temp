package faq

import "github.com/gofiber/fiber/v2"

func Routes(app fiber.Router) {
	r := app.Group("/faq")

	r.Post("/", createQuestion)
	r.Get("/", getAllQuestions)
	r.Get("/:key", getQuestionByKey)
	r.Delete("/:key", deleteQuestion)

	//r.Post("/cat", createCategory)
	//r.Get("/cat", GetCategory)
	r.Get("/cat/:catKey", getFAQByCategory)

	r.Put("/:key", updateFAQ)

	//r.Delete("/cat/:catKey", deleteCategory)

	r.Post("/srch", searchIntoQuestions)

	r.Post("/FeedBack", createFeedBackCollection)
	r.Get("/FeedBack/:key", getFeedBackByKey)
}

func CategoryRoutes(app fiber.Router) {
	r := app.Group("/faq-category")

	r.Post("/", createCategory)
	r.Get("/", GetCategory)
	r.Delete("/:catKey", deleteCategory)

}
