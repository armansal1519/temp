package commentOnProduct

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/product-comment")

	r.Post("/", middleware.Auth, createCommit)
	r.Get("/:categoryUrl/:productKey", getProductComment)
	r.Get("/user", middleware.Auth, getByUserKey)
	r.Get("/", getAll)
	r.Put("/:key", middleware.Auth, updateComment)
	r.Put("/admin/:key", middleware.IsAdmin, adminUpdateComment)
	r.Delete("/:key", middleware.IsAdmin, middleware.Auth, deleteComment)

}
