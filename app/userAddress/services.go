package userAddress

import (
	"bamachoub-backend-go-v1/app/users"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

// addAddress create user address
// @Summary create user address
// @Description create user address
// @Tags user address
// @Accept json
// @Produce json
// @Param AddressIn body AddressIn true "Address"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} AddressOut{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /user-address [post]
func addAddress(c *fiber.Ctx) error {
	a := new(AddressIn)
	if err := utils.ParseBodyAndValidate(c, a); err != nil {
		return c.JSON(err)
	}
	aCol := database.GetCollection("userAddress")
	userKey := c.Locals("userKey").(string)
	a.UserKey = userKey
	if a.IsForMySelf {
		var u users.UserOut
		userCol := database.GetCollection("users")
		_, err := userCol.ReadDocument(context.Background(), userKey, &u)
		if err != nil {
			return c.JSON(err)
		}
		a.LastName = u.LastName
		a.FirstName = u.FirstName
		a.NationalCode = u.NationalCode
		a.PhoneNumber = u.PhoneNumber
	}

	meta, err := aCol.CreateDocument(context.Background(), a)
	if err != nil {
		return c.JSON(err)

	}
	return c.JSON(meta)
}

// getAddressByUserKey get each user address key
// @Summary return user address by user key
// @Description return user address by user key
// @Tags user address
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} AddressOut{}
// @Failure 404 {object} string{}
// @Router /user-address/user [get]
func getAddressByUserKey(c *fiber.Ctx) error {
	userKey := c.Locals("userKey")
	query := fmt.Sprintf("for i in userAddress filter i.userKey == \"%v\" return i", userKey)
	return c.JSON(database.ExecuteGetQuery(query))
}

// getAddressByKey get each user address key
// @Summary return user address by user key
// @Description return user address by user key
// @Tags user address
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} AddressOut{}
// @Failure 404 {object} string{}
// @Router /user-address/{key} [get]
func getAddressByKey(key string) (*AddressOut, error) {
	var a AddressOut
	aCol := database.GetCollection("userAddress")
	_, err := aCol.ReadDocument(context.Background(), key, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// editAddress update user address
// @Summary update user address
// @Description update user address
// @Tags user address
// @Accept json
// @Produce json
// @Param AddressIn body AddressIn true "Address"
// @Param key path string true "key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} AddressOut{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /user-address/{key} [put]
func editAddress(c *fiber.Ctx) error {
	a := new(AddressIn)
	key := c.Params("key")
	if err := utils.ParseBodyAndValidate(c, a); err != nil {
		return c.JSON(err)
	}
	oldAddress, err := getAddressByKey(key)
	if err != nil {
		return c.JSON(err)
	}
	if oldAddress.UserKey != key {
		return c.Status(403).SendString("Unauthorized")
	}
	aCol := database.GetCollection("userAddress")
	var newAddress AddressOut
	ctx := driver.WithReturnNew(context.Background(), &newAddress)
	_, err = aCol.UpdateDocument(ctx, key, a)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(newAddress)
}

// removeAddress delete user address
// @Summary delete user address
// @Description delete user address
// @Tags user address
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /user-address/{key} [delete]
func removeAddress(c *fiber.Ctx) error {
	key := c.Params("key")
	aCol := database.GetCollection("userAddress")
	meta, err := aCol.RemoveDocument(context.Background(), key)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}
