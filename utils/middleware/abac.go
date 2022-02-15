package middleware

import (
	"bamachoub-backend-go-v1/app/supplyWorkers"
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/jwt"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Abac(accesses []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")

		if h == "" {
			err := fmt.Errorf("header not found")
			return utils.CustomErrorResponse(400, 4001, err, "", c)
		}

		// Spliting the header
		chunks := strings.Split(h, " ")

		// If header signature is not like `Bearer <token>`, then throw
		// This is also required, otherwise chunks[1] will throw out of bound error
		if len(chunks) < 2 {
			err := fmt.Errorf("header is not acceptable: %v", h)
			return utils.CustomErrorResponse(400, 4001, err, "", c)
		}

		// Verify the token which is in the chunks
		payload, err := jwt.Verify(chunks[1])

		if err != nil {
			return utils.CustomErrorResponse(401, 4011, err, "", c)
		}

		swAccess := getAccess(payload.Key)
		for i := range accesses {
			for j := range swAccess {
				if accesses[i] == swAccess[j] {

					return c.Next()
				}
			}
		}

		return fiber.NewError(fiber.StatusForbidden,
			fmt.Sprintf("you can access this endpoint %v", swAccess),
		)
	}
}

func getAccess(key string) []string {
	//splitedId := strings.Split(id, "/")
	//key := splitedId[1]
	sw := supplyWorkers.GetSupplyWorkerByKey(key)
	return sw.Access
}
