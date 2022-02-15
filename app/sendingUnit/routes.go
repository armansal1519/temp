package sendingUnit

import "github.com/gofiber/fiber/v2"

func Routes(app fiber.Router) {
	r := app.Group("/send-unit")

	//TODO admin
	r.Post("/", createSendUnit)
	r.Post("/:key", acceptSendUnit)
	r.Put("/:op/:unitKey/:trKey", addOrRemoveSendUnits)
	r.Get("/", getSendUnits)
	r.Get("/by-user-key", getSendUnitsByUserKey)
	r.Delete("/:key", removeSendUnit)
	r.Delete("/tr/:key", removeTr)

}
