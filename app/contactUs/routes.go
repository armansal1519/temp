package contactUs

import "github.com/gofiber/fiber/v2"

func Routes(app fiber.Router) {
	r := app.Group("/contact-us")

	r.Post("/", create)
	r.Get("/", getAll)
	r.Get("/:key", getByKey)

	//TODO only admin can delete that - fix it later
	r.Put("/:key", update)
	r.Delete("/:key", delete)

}
