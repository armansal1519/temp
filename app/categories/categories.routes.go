package categories

import (
	"bamachoub-backend-go-v1/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/categories")
	r.Get("", getBaseCategories)
	r.Get("/:key", getCategoryByKey)
	r.Get("/price-range/:dbName/:key", getPriceRangeUnderOnCategory)
	r.Post("/base", middleware.CheckAdmin, middleware.AdminHasAccess([]string{"nakon mohammad"}), CreateBaseCategory)
	r.Post("", CreateCategory)
	r.Put("/:key", update)

}
