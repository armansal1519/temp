package productQA

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/products-q-a")

	r.Get("", getAll)
	r.Get("/user/:op", middleware.Auth, getQAForUser)
	r.Get("/:categoryUrl/:productKey", getQAByProductKey)
	r.Post("", middleware.Auth, createQA)
	r.Post("/likes/:op/:qaKey", middleware.Auth, likes)
	r.Put("/:key", updateQA)
	r.Put("/admin/:key", adminUpdateQA)

	r.Delete("/key", middleware.Auth, removeQA)

}
