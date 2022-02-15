package supplierEmployees

import (
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/middleware"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func Routes(app fiber.Router) {
	r := app.Group("/supplier-employee")

	r.Get("/one", middleware.SupplierEmployeeAuth([]string{}), func(c *fiber.Ctx) error {
		key := c.Locals("key").(string)
		res, err := getSupplierEmployeeByKey(key)
		if err != nil {
			return c.JSON(err)
		}
		res.HashRefreshToken = ""
		res.HashPassword = ""
		return c.JSON(res)
	})
	r.Put("/add-update-pool", middleware.SupplierEmployeeAuth([]string{}), addToUpdatePool)

	r.Post("/create-by-admin/:key", func(c *fiber.Ctx) error {
		spKey := c.Params("key")
		err := createSupplierEmployeeFromSupplierPreview(spKey)
		if err != nil {
			ce := utils.CError{
				Code:    1,
				Error:   fmt.Sprintf("%v", err),
				DevInfo: "",
				UserMsg: "مشکل در ساختن تامین کننده",
			}

			return c.JSON(ce)
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

}

func AuthRoutes(app fiber.Router) {
	r := app.Group("/supplier-employee-auth")

	r.Get("/get-supplier-preview", getSupplierPreview)

	r.Post("/create-supplier-preview", func(c *fiber.Ctx) error {
		ce := new(createSupplierPreview)
		if err := utils.ParseBodyAndValidate(c, ce); err != nil {
			log.Println(1)
			return c.JSON(err)
		}
		resp, err := CreateSupplierPreview(*ce)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(resp)
	})
	r.Post("/get-validation-code", GetValidationCode)
	r.Post("/check-validation-code", CheckValidationCode)
	r.Post("/login", func(c *fiber.Ctx) error {
		lr := new(loginRequest)
		if err := utils.ParseBodyAndValidate(c, lr); err != nil {
			return c.JSON(err)
		}
		resp, err := supplierEmployeeLogin(lr.PhoneNumber, lr.Password)
		if err != nil {
			return c.Status(401).JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(resp)
	})
	r.Post("/get-changePassword-code", GetChangePasswordCode)
	r.Post("/changePassword-without-login", changePasswordWithoutLogin)
	r.Post("/change-password-with-login", middleware.SupplierEmployeeAuth([]string{}), func(c *fiber.Ctx) error {
		cvc := new(changePasswordWithLoginRequest)
		if err := utils.ParseBodyAndValidate(c, cvc); err != nil {
			return c.JSON(err)
		}
		key := c.Locals("key").(string)
		err := changePasswordWithLogin(key, cvc.Password)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
	r.Get("/get-refresh-token/:token", getRefreshToken)
}
