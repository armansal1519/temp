package productSuggestion

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/product-suggestion")
	r.Get("", getProductSuggestions)
	r.Post("", middleware.GetSupplierByEmployee, AddNewSuggestion)
	r.Post("/sample", middleware.Auth, AddFromSample)
	r.Post("/better-price", middleware.Auth, AddBetterPrice)
	r.Delete("/:key", deleteProductSuggestions)

}
