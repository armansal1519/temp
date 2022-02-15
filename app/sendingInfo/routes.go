package sendingInfo

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/sending-info")

	r.Post("/add-interval", addInterval)
	r.Post("/remove-interval", removeInterval)
	r.Post("/info/:orderKey", middleware.Auth, CreateSendingInfo)
	r.Get("/interval", getSendDayInterval)
	r.Get("/info/:key", middleware.Auth, GetSendingInfoByKey)

	r.Put("/info/:key", middleware.Auth, editSendingInfo)

}
