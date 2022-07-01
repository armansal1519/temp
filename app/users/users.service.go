package users

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/jwt"
	"bamachoub-backend-go-v1/utils/password"
	"context"
	"encoding/json"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	hashRefreshToken, _ := password.HashPassword(refreshToken)
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

// updateUser update user
// @Summary update user
// @Description update user , phone number is always locked if isAuthenticated firstName lastName and nationalCode become locked as well
// @Tags user
// @Accept json
// @Produce json
// @Param data body updateUserDTO true "data"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} UserOut{}
// @Failure 404 {object} string{}
// @Router /user [patch]
func updateUser(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	b := new(updateUserDTO)
	if err := utils.ParseBodyAndValidate(c, b); err != nil {
		return c.JSON(err)
	}
	oldUser, err := GetUserByKey(userKey)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	if oldUser.IsAuthenticated {
		b.FirstName = oldUser.FirstName
		b.LastName = oldUser.LastName
		b.NationalCode = oldUser.NationalCode
	}

	var newUser UserOut
	userCol := database.GetCollection("users")
	ctx := driver.WithReturnNew(context.Background(), &newUser)
	_, err = userCol.UpdateDocument(ctx, userKey, b)
	if err != nil {
		return c.Status(500).JSON(err)

	}
	return c.JSON(newUser)

}

// addCardInfo add card to user
// @Summary add card to user
// @Description add card to user
// @Tags user
// @Accept json
// @Produce json
// @Param data body cardInfo true "data"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} UserOut{}
// @Failure 404 {object} string{}
// @Router /user/card [patch]
func addCardInfo(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	b := new(cardInfo)
	if err := utils.ParseBodyAndValidate(c, b); err != nil {
		return c.JSON(err)
	}
	var newUser UserOut
	userCol := database.GetCollection("users")
	ctx := driver.WithReturnNew(context.Background(), &newUser)
	newUser.UserCards.CardUserName = b.CardUserName
	newUser.UserCards.Number = b.Number
	newUser.UserCards.BankName = b.BankName
	_, err := userCol.UpdateDocument(ctx, userKey, newUser)
	if err != nil {
		return c.Status(500).JSON(err)

	}
	return c.JSON(newUser)

}

// userAuthentication Authentication user
// @Summary Authentication user
// @Description Authentication user
// @Tags user
// @Accept json
// @Produce json
// @Param data body AuthenticationDto true "data"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} UserOut{}
// @Failure 404 {object} string{}
// @Router /user/auth [post]
func userAuthentication(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	b := new(AuthenticationDto)
	if err := utils.ParseBodyAndValidate(c, b); err != nil {
		return c.JSON(err)
	}
	u, err := GetUserByKey(userKey)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	resp, err := http.Get(fmt.Sprintf("https://inquery.ir/:60?Token=m5uL2Vl2QPISO5pdLrdoehYtB5E&IdCode=%v&BirthDate=%v&Mobile=%v", b.NationalCode, b.BirthDate, u.PhoneNumber))
	if err != nil {
		return c.Status(500).JSON(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	var ar authResponse
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	log.Println(ar.Result.Detail)
	if strings.Contains(ar.Result.Detail, "موفق") {
		u.IsAuthenticated = true
		u.BirthDate = b.BirthDate

		var newUser UserOut
		userCol := database.GetCollection("users")
		ctx := driver.WithReturnNew(context.Background(), &newUser)
		_, err := userCol.UpdateDocument(ctx, userKey, u)
		if err != nil {
			return c.Status(500).JSON(err)
		}
		return c.JSON(newUser)

	}

	return c.Status(403).JSON(ar.Result.Detail)
}
