package supplierRequests

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/supplier-request")

	r.Post("/remove-from-wallet/:amount", middleware.GetSupplierByEmployee, removeFromWallet)

}
