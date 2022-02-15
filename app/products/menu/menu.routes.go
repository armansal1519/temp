package menu

import "github.com/gofiber/fiber/v2"

func Routes(app fiber.Router) {
	r := app.Group("/menu")
	//r.Post("", createProductMenu)
	r.Get("/category/:key", getMenuByCategoryKey)
	r.Post("/add/:key", addToMenu)

}
