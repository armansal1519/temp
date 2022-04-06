package search

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/search")

	r.Post("/ms", SetMostSearch)
	r.Post("/", Search)

}
