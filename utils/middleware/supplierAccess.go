package middleware

import (
	"bamachoub-backend-go-v1/utils/jwt"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

func SupplierEmployeeAccess(c *fiber.Ctx) error {
	c.Locals("role", "manager")
	return c.Next()
}

func GetSupplierByEmployee(c *fiber.Ctx) error {
	h := c.Get("Authorization")
	if h == "" {
		return fiber.ErrUnauthorized
	}
	chunks := strings.Split(h, " ")
	if len(chunks) < 2 {

		return fiber.ErrUnauthorized
	}
	payload, err := jwt.VerifySupplierEmployee(chunks[1], false)

	if err != nil {
		log.Println(err)
		return c.Status(401).JSON(fmt.Sprintf("problem with token %v :", err))
	}

	//key := strings.Split(payload.Key, "/")[1]
	//q := fmt.Sprintf("for se in supplierEmployee filter se._key==\"%v\" for s in suppliers filter s._key==se.supplierKey return s", payload.Key)
	//db := database.GetDB()
	//ctx := context.Background()
	//cursor, err := db.Query(ctx, q, nil)
	//if err != nil {
	//	return c.JSON(fmt.Errorf("error while excuting query: %v \n error:%v", q, err))
	//}
	//defer cursor.Close()
	//
	//var doc supplier
	//_, err = cursor.ReadDocument(ctx, &doc)
	//if err != nil {
	//	return c.JSON(fmt.Errorf("error while reading document: %v \n error:%v", q, err))
	//
	//}
	c.Locals("supplierId", payload.SupplierKey)
	c.Locals("supplierEmployeeKey",payload.Key)

	return c.Next()

}

type SupplierIn struct {
	Address      string  `json:"address" validate:"required"`
	Latitude     float64 `json:"latitude" validate:"required"`
	Longitude    float64 `json:"longitude" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Code         string  `json:"code" validate:"required"`
	Area         float64 `json:"area" validate:"required"`
	AreaWithRoof float64 `json:"areaWithRoof" validate:"required"`
	PhoneNumber  string  `json:"phoneNumber" validate:"required"`
	Status       string  `json:"status" validate:"required"`
	CreateAt     float64 `json:"createAt" validate:"required"`
}

type supplier struct {
	driver.DocumentMeta
	SupplierIn
}
