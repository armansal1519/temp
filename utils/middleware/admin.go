package middleware

import (
	"bamachoub-backend-go-v1/utils/jwt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

func IsAdmin(c *fiber.Ctx) error {
	h := c.Get("Authorization")
	if h == "" {
		c.Locals("isAdmin", false)
		return c.Next()

	}
	chunks := strings.Split(h, " ")
	if len(chunks) < 2 {
		c.Locals("isAdmin", false)
		return c.Next()

	}

	user, err := jwt.VerifyAdmin(chunks[1], false)
	if err != nil {
		c.Locals("isAdmin", false)
		return c.Next()

	}
	c.Locals("isAdmin", true)
	c.Locals("adminKey", user.Key)
	c.Locals("adminAccess", user.Access)
	return c.Next()

}

func CheckAdmin(c *fiber.Ctx) error {
	h := c.Get("Authorization")
	if h == "" {
		log.Println(1)
		return fiber.ErrUnauthorized

	}
	chunks := strings.Split(h, " ")
	if len(chunks) < 2 {
		log.Println(2)
		return fiber.ErrUnauthorized
	}

	user, err := jwt.VerifyAdmin(chunks[1], false)
	if err != nil {

		return fiber.ErrUnauthorized
	}
	c.Locals("isAdmin", true)
	c.Locals("adminKey", user.Key)
	c.Locals("adminAccess", user.Access)
	return c.Next()

}

func AdminHasAccess(accessArr []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		access := c.Locals("adminAccess").(string)
		aa := strings.Split(access, ",")
		for _, s := range aa {
			for _, s2 := range accessArr {
				if s == s2 {
					return c.Next()
				}
			}
		}
		return fiber.ErrUnauthorized

	}
}

func IsSuperAdmin(c *fiber.Ctx) error {
	access := c.Locals("adminAccess").(string)
	aa := strings.Split(access, ",")
	for _, s := range aa {
		if s == "sa" {
			c.Locals("isSuperAdmin", true)
			return c.Next()
		}
	}
	return c.Next()

}
