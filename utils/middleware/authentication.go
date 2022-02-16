package middleware

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils/jwt"
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Auth is the authentication middleware
func Auth(c *fiber.Ctx) error {
	isAdmin := c.Locals("isAdmin")

	if isAdmin != nil {
		b := isAdmin.(bool)
		if b {
			return c.Next()
		}
	}

	h := c.Get("Authorization")

	if h == "" {
		return fiber.ErrUnauthorized
	}

	// Spliting the header
	chunks := strings.Split(h, " ")

	// If header signature is not like `Bearer <token>`, then throw
	// This is also required, otherwise chunks[1] will throw out of bound error
	if len(chunks) < 2 {
		return fiber.ErrUnauthorized
	}

	//_,err:=jwt.VerifyAdmin(chunks[1],false)
	//if err==nil {
	//	c.Next()
	//}

	// Verify the token which is in the chunks
	user, err := jwt.Verify(chunks[1])

	if err != nil {
		return fiber.ErrUnauthorized
	}

	c.Locals("userKey", user.Key)
	//fmt.Println(user)
	return c.Next()
}

func AddToCartAuth(c *fiber.Ctx) error {
	h := c.Get("Authorization")
	isLogin := true

	if h == "" {
		c.Locals("userKey", "")
		c.Locals("isLogin", false)
		return c.Next()
	}
	chunks := strings.Split(h, " ")
	if len(chunks) < 2 {
		isLogin = false
	}

	user, err := jwt.Verify(chunks[1])

	if err != nil {
		isLogin = false
	}
	if isLogin {
		c.Locals("userKey", user.Key)
		c.Locals("isLogin", true)

		return c.Next()
	}

	c.Locals("userKey", "")
	c.Locals("isLogin", false)

	return c.Next()
}

func IsAuthenticated(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	if userKey == "" {
		c.Locals("isAuthenticated", false)
		return c.Next()
	}
	log.Println("in isAuthenticated")
	u, err := GetUserByKey(userKey)
	if err != nil {
		c.JSON(err)
	}

	if u.IsAuthenticated {
		c.Locals("isAuthenticated", true)
		return c.Next()
	}
	c.Locals("isAuthenticated", false)
	log.Println("isAuthenticated==false")
	return c.Next()

}

func SupplierEmployeeAuth(access []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")

		if h == "" {
			log.Println(3)

			return fiber.ErrUnauthorized
		}

		// Spliting the header
		chunks := strings.Split(h, " ")

		// If header signature is not like `Bearer <token>`, then throw
		// This is also required, otherwise chunks[1] will throw out of bound error
		if len(chunks) < 2 {
			log.Println(1)
			return fiber.ErrUnauthorized
		}

		// Verify the token which is in the chunks
		payload, err := jwt.VerifySupplierEmployee(chunks[1], false)

		if err != nil {
			log.Println(2)

			return fiber.ErrUnauthorized
		}

		c.Locals("key", payload.Key)
		c.Locals("supplierKey", payload.SupplierKey)
		//if payload.Role == "admin" || payload.Role == "superAdmin" {
		//	return c.Next()
		//}

		//c.Locals("userId", payload.Id)
		//c.Locals("role", payload.Role)
		if len(access) == 0 {
			return c.Next()
		}

		for _, v := range access {
			for _, s := range strings.Split(payload.Access, ",") {
				if v == s {
					return c.Next()
				}
			}

		}
		log.Println(4)

		return fiber.ErrUnauthorized
	}
}

func GetUserByKey(key string) (*UserOut, error) {
	userCol := database.GetCollection("users")
	var u UserOut
	_, err := userCol.ReadDocument(context.Background(), key, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

type UserOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	user
}

type user struct {
	PhoneNumber     string `json:"phoneNumber"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	BirthDate       string `json:"birthDate"`
	NationalCode    string `json:"nationalCode"`
	Level           string `json:"level"`
	CreatedAt       int64  `json:"createdAt"`
	LastLogin       int64  `json:"lastLogin"`
	IsAuthenticated bool   `json:"isAuthenticated"`
}
