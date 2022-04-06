package userFav

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/user-fav")

	r.Post("/add", middleware.Auth, addToUserFav)
	r.Post("/remove", middleware.Auth, removeFromUserFav)
	r.Get("", middleware.Auth, getUserFav)

}
