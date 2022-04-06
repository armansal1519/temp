package productStructure

import "github.com/gofiber/fiber/v2"

func Routes(app fiber.Router) {
	r := app.Group("/products-structure")
	r.Get("/", getAll)
	r.Get("/:key", getProductStructureByCategoryKey)
	r.Post("", createProductStructure)
	r.Put("", updateProductStructAndMenu)

}
