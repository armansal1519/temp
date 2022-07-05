package paymentAndWallet

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func WalletRoutes(app fiber.Router) {
	r := app.Group("/wallet")

	r.Get("/data", middleware.GetSupplierByEmployee, getDataForSupplierWalletPage)
	r.Get("/user/history/", middleware.Auth, func(c *fiber.Ctx) error {
		userKey := c.Locals("userKey").(string)
		offset := c.Query("offset")
		limit := c.Query("limit")
		if limit == "" || offset == "" {
			return c.Status(400).SendString("offset or limit is empty")
		}
		res, t, err := GetUserWalletHistoryByUserKey(userKey, offset, limit)
		if err != nil {
			return c.Status(500).JSON(err)
		}
		return c.JSON(fiber.Map{
			"history": res,
			"income":  t.Income,
			"outcome": t.Outcome,
		})
	})
	r.Post("/url", middleware.GetSupplierByEmployee, getPaymentUrlForSupplierWallet)
	r.Post("/user/url", middleware.Auth, getPaymentUrlForUserWallet)
	r.Post("/verify/:key", middleware.GetSupplierByEmployee, VerifyPaymentAndAddToWallet)
	r.Post("/user/verify/:key", middleware.Auth, VerifyUserPaymentAndAddToWallet)

}

func PaymentRoutes(app fiber.Router) {
	r := app.Group("/payment")

	r.Post("/by-url", middleware.Auth, getPaymentUrl)
	r.Post("/by-image", middleware.Auth, createPaymentByImage)
	r.Post("/verify-url/:key", middleware.Auth, verifyPaymentUrl)
	//
	r.Post("/by-check", middleware.Auth, createCheckPayment)
	//TODO admin
	r.Post("/verify-image/:key", verifyPaymentImage)
	r.Post("/verify-check/:key", verifyCheckImage)
	//
	r.Get("/user/:orderKey?", middleware.Auth, getPaymentByUserKey)
	r.Get("/info", paymentConst)
	r.Get("/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")
		resp, err := getPaymentByKey(key)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(resp)
	})
	//
	//r.Post("/filter", filterPayment)

}

func SupplierConfirmationRoute(app fiber.Router) {
	r := app.Group("/suppliers-confirmation")
	r.Get("", middleware.GetSupplierByEmployee, GetOrderConfirmationBySupplierKey)

	r.Post("/approve/:infoKey", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		infoKey := c.Params("infoKey")
		supplierId := c.Locals("supplierId").(string)
		supplierEmployeeId := c.Locals("supplierEmployeeKey").(string)
		err := approveOrder(infoKey, supplierId, supplierEmployeeId)
		if err != nil {
			return c.JSON(err)
		}
		return c.SendString("ok")

	})
	r.Post("/reject/:infoKey", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		infoKey := c.Params("infoKey")
		supplierKey := c.Locals("supplierId").(string)

		err := rejectOrder(infoKey, supplierKey, "user")
		if err != nil {
			return c.Status(500).JSON(err)
		}
		return c.SendString("ok")

	})

	r.Post("/g-reject/:orderItemKey", middleware.GetSupplierByEmployee, func(c *fiber.Ctx) error {
		orderItemKey := c.Params("orderItemKey")
		supplierKey := c.Locals("supplierId").(string)

		newOrderItem, err := graphRejectOrder(orderItemKey, supplierKey, false)
		if err != nil {
			return c.Status(500).JSON(err)
		}
		return c.JSON(newOrderItem)

	})
	r.Get("/g-reject", middleware.IsAdmin, middleware.GetSupplierByEmployee, getGraphRejectedOrder)
	r.Post("/g-accept/:orderItemKey", middleware.GetSupplierByEmployee, acceptARejectedOrder)

}

func RejectionRoutes(app fiber.Router) {
	r := app.Group("/rejection")

	r.Get("/pool", middleware.GetSupplierByEmployee, getRejectionPool)
	r.Post("/by-image/:key", rejectPaymentImage)

}
