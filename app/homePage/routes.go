package homepage

import "github.com/gofiber/fiber/v2"




func Routes(app fiber.Router) {
	r := app.Group("/homepage")

	r.Post("/base",setBaseData)

	r.Get("/",getHomePage)

}
