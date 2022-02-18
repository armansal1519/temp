package users

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/middleware"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func Routes(app fiber.Router) {
	r := app.Group("/user")
	r.Get("", GetUsers)
	r.Get("/one", middleware.Auth, func(c *fiber.Ctx) error {
		userKey := c.Locals("userKey").(string)
		resp, err := GetUserByKey(userKey)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(resp)
	})
	r.Post("", CreateUser)

}

func AuthRoutes(app fiber.Router) {
	r := app.Group("/user-auth")

	r.Post("/get-validation-code", func(c *fiber.Ctx) error {
		lr := new(checkForLoginReq)
		if err := utils.ParseBodyAndValidate(c, lr); err != nil {
			return c.JSON(err)
		}
		resp, err := checkUserPhoneNumberForLogin(lr.PhoneNumber)
		if err != nil {
			return c.Status(401).JSON(fmt.Sprintf("%v", err))
		}
		return c.JSON(resp)
	})

	r.Post("/login", func(c *fiber.Ctx) error {
		ld := new(LoginDto)
		if err := utils.ParseBodyAndValidate(c, ld); err != nil {
			return c.JSON(err)
		}

		resp, err := loginWithValidationCode(ld.PhoneNumber, ld.Code)
		if err != nil {
			return c.Status(401).JSON(fmt.Sprintf("%v", err))
		}

		//handling temp user cart
		tempUserKey := c.Get("temp-user-key", "")
		if tempUserKey != "" {
			a := ""
			if resp.User.IsAuthenticated {
				a = "authenticated"
			} else {
				a = "login"
			}
			q := fmt.Sprintf("for i in cart filter i.userKey==\"%v\"  and i.userAuthType==\"headless\" update i with {userKey:\"%v\",userAuthType:\"%v\"} in cart", tempUserKey, resp.User.Key, a)
			log.Println(q)

			database.ExecuteGetQuery(q)
		}

		return c.JSON(resp)
	})
	r.Post("/register", func(c *fiber.Ctx) error {
		ld := new(LoginDto)
		if err := utils.ParseBodyAndValidate(c, ld); err != nil {
			return c.JSON(err)
		}
		resp, err := registerWithValidationCode(ld.PhoneNumber, ld.Code)
		if err != nil {
			return c.Status(401).JSON(fmt.Sprintf("%v", err))
		}
		//handling temp user cart
		tempUserKey := c.Get("temp-user-key", "")
		if tempUserKey != "" {
			a := ""
			if resp.User.IsAuthenticated {
				a = "authenticated"
			} else {
				a = "login"
			}
			q := fmt.Sprintf("for i in cart filter i.userKey==\"%v\"  and i.userAuthType==\"headless\" update i with {userKey:\"%v\",userAuthType:\"%v\"} in cart", tempUserKey, resp.User.Key, a)
			database.ExecuteGetQuery(q)
		}
		return c.JSON(resp)
	})

	r.Get("/get-refresh-token/:token", getRefreshToken)

	r.Get("/services/:token", getUserByJwt)
}
