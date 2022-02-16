package users

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils/jwt"
	"bamachoub-backend-go-v1/utils/password"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func GetUsers(c *fiber.Ctx) error {
	return c.JSON(database.GetAll("users", 0, 32))
}

type createUserType struct {
	PhoneNumber string `json:"phoneNumber" `
	Password    string `json:"password" `
	Role        string
}

func CreateUser(c *fiber.Ctx) error {
	u := new(createUserType)

	if err := c.BodyParser(u); err != nil {
		return err
	}

	log.Println(u.PhoneNumber)
	log.Println(u.Password)
	return c.JSON(u)
}
func getUserByPhoneNumber(phoneNumber string) (*UserOut, error) {
	db := database.GetDB()
	ctx := context.Background()
	query := fmt.Sprintf("for u in users filter u.phoneNumber==\"%v\" return u", phoneNumber)
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var user UserOut
	_, err = cursor.ReadDocument(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByKey get user by accessToken
// @Summary get user by accessToken
// @Description get user by accessToken
// @Tags  user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []UserOut{}
// @Failure 400 {object} string
// @Router /user/one [get]
func GetUserByKey(key string) (*UserOut, error) {
	userCol := database.GetCollection("users")
	var u UserOut
	_, err := userCol.ReadDocument(context.Background(), key, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func getUserByJwt(c *fiber.Ctx) error {
	from := c.Get("service")
	if from != "samples" {
		return c.JSON(fiber.ErrUnauthorized)
	}
	token := c.Params("token")
	payload, err := jwt.Verify(token)
	userKey := payload.Key
	if err != nil {
		return fiber.ErrUnauthorized
	}

	p := jwt.TokenPayload{
		Key: userKey,
	}
	accessToken := jwt.GenerateAccessToken(&p)
	refreshToken := jwt.GenerateRefreshToken(&p)
	hashRefreshToken := password.Generate(refreshToken)
	urt := UpdateRefreshTokenServices{
		HashRefreshTokenServices: hashRefreshToken,
		LastLogin:                time.Now().Unix(),
	}
	userCol := database.GetCollection("users")
	var uOut UserOut
	ctx := driver.WithReturnNew(context.Background(), &uOut)
	_, err = userCol.UpdateDocument(ctx, userKey, urt)
	if err != nil {
		return c.JSON(err)
	}
	lr := loginAndRegistrationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         uOut,
	}
	return c.JSON(lr)
}

func CreateHeadlessUser() (string, error) {
	data := headlessUser{State: "headless", PhoneNumber: fmt.Sprintf("%v", time.Now().Nanosecond())}
	userCol := database.GetCollection("users")
	meta, err := userCol.CreateDocument(context.Background(), data)
	if err != nil {
		return "", err
	}
	return meta.Key, nil
}
