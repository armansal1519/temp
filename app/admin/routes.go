package admin

import (
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/middleware"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/admin")

	r.Post("", createAdmin)
	r.Get("",
		middleware.CheckAdmin, middleware.AdminHasAccess([]string{"read-admin"}),
		getAll)
	r.Get("by-access",middleware.CheckAdmin,getAdminByAccessToken)
	r.Get("/access" , getAccessArray)
	r.Get("/:key", getAdminByKey)
	r.Put("/:key", updateAdmin)

}

func AuthRoutes(app fiber.Router) {
	r := app.Group("/admin-auth")
	//
	//r.Get("/get-supplier-preview", getSupplierPreview)
	//
	//r.Post("/create-supplier-preview", func(c *fiber.Ctx) error {
	//	ce := new(createSupplierPreview)
	//	if err := utils.ParseBodyAndValidate(c, ce); err != nil {
	//		log.Println(1)
	//		return c.JSON(err)
	//	}
	//	resp, err := CreateSupplierPreview(*ce)
	//	if err != nil {
	//		return c.JSON(err)
	//	}
	//	return c.JSON(resp)
	//})
	//r.Post("/get-validation-code", GetValidationCode)
	//r.Post("/check-validation-code", CheckValidationCode)
	r.Post("/login", func(c *fiber.Ctx) error {
		lr := new(loginRequest)
		if err := utils.ParseBodyAndValidate(c, lr); err != nil {
			return c.JSON(err)
		}
		resp, err := Login(lr.PhoneNumber, lr.Password)
		if err != nil {
			return c.Status(401).JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(resp)
	})
	//r.Post("/get-changePassword-code", GetChangePasswordCode)
	r.Post("/change-password", middleware.CheckAdmin, middleware.IsSuperAdmin, changePassword)
	//r.Post("/change-password-with-login", middleware.SupplierEmployeeAuth([]string{}), func(c *fiber.Ctx) error {
	//	cvc := new(changePasswordWithLoginRequest)
	//	if err := utils.ParseBodyAndValidate(c, cvc); err != nil {
	//		return c.JSON(err)
	//	}
	//	key := c.Locals("key").(string)
	//	err := changePasswordWithLogin(key, cvc.Password)
	//	if err != nil {
	//		return c.JSON(err)
	//	}
	//	return c.JSON(fiber.Map{
	//		"status": "ok",
	//	})
	//})
	r.Get("/get-refresh-token/:token", getRefreshToken)
}
