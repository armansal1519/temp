package massage

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/msg")

	r.Post("by-phone",
		//middleware.CheckAdmin, middleware.AdminHasAccess([]string{"write-massageAndFAQ"}),
		middleware.TestAdmin,
		sendMsgByPhoneNumberUsers)
	r.Get("user", middleware.Auth, getMassageByUserKey)
	r.Get("supplier", middleware.GetSupplierByEmployee, getMassageBySupplierKey)

}
