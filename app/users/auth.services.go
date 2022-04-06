package users

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/jwt"
	"bamachoub-backend-go-v1/utils/password"
	"bamachoub-backend-go-v1/utils/sms"
	"context"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

// checkUserPhoneNumberForLogin send validation code by sms to given phoneNumber
// @Summary send validation code by sms to given phoneNumber
// @Description send validation code by sms to given phoneNumber for user registration or login
// @Tags user Auth
// @Accept json
// @Produce json
// @Param data body checkForLoginReq true "phoneNumber"
// @Success 200 {object} checkForLoginRes{}
// @Failure 400 {object} string{}
// @Router /user-auth/get-validation-code   [post]
func checkUserPhoneNumberForLogin(phoneNumber string) (checkForLoginRes, error) {
	_, err := getUserByPhoneNumber(phoneNumber)

	if !errors.Is(err, driver.NoMoreDocumentsError{}) && err != nil {
		return checkForLoginRes{}, err
	}

	if errors.Is(err, driver.NoMoreDocumentsError{}) {
		rn := utils.GenRandomNUmber(1001, 9999)
		gvCol := database.GetCollection("userValidationCode")
		svc := SaveValidationCode{
			Key:       phoneNumber,
			Code:      fmt.Sprintf("%v", rn),
			CreatedAt: time.Now().Add(time.Second * 120).Unix(),
		}
		log.Println(time.Now().Unix())
		_, err = gvCol.CreateDocument(context.Background(), svc)
		if err != nil {
			return checkForLoginRes{}, err
		}
		pArr := sms.ParameterArray{
			Parameter:      "VerificationCode",
			ParameterValue: fmt.Sprintf("%v", rn),
		}
		sms.SendSms(phoneNumber, "48985", []sms.ParameterArray{pArr})
		log.Println(rn)

		time.AfterFunc(120*time.Second, func() {

			_, err = gvCol.RemoveDocument(context.Background(), phoneNumber)
		})
		res := checkForLoginRes{
			PhoneNumber:  phoneNumber,
			IsRegistered: false,
		}
		return res, nil
	} else {
		rn := utils.GenRandomNUmber(1001, 9999)
		gvCol := database.GetCollection("userValidationCode")
		svc := SaveValidationCode{
			Key:       phoneNumber,
			Code:      fmt.Sprintf("%v", rn),
			CreatedAt: time.Now().Add(time.Second * 120).Unix(),
		}
		log.Println(time.Now().Unix())
		_, err = gvCol.CreateDocument(context.Background(), svc)
		if err != nil {
			return checkForLoginRes{}, err
		}
		pArr := sms.ParameterArray{
			Parameter:      "VerificationCode",
			ParameterValue: fmt.Sprintf("%v", rn),
		}
		sms.SendSms(phoneNumber, "48985", []sms.ParameterArray{pArr})
		log.Println(rn)

		time.AfterFunc(120*time.Second, func() {

			_, err = gvCol.RemoveDocument(context.Background(), phoneNumber)
		})

		return checkForLoginRes{
			PhoneNumber:  phoneNumber,
			IsRegistered: true,
		}, nil
	}

}

// loginWithValidationCode user login with validation code
// @Summary user login with validation code
// @Description user login with validation code
// @Tags user Auth
// @Accept json
// @Produce json
// @Param data body LoginDto true "data"
// @Success 200 {object} loginAndRegistrationResponse{}
// @Failure 400 {object} string{}
// @Failure 401 {object} string{}
// @Router /user-auth/login [post]
func loginWithValidationCode(phoneNumber string, validationCode string) (loginAndRegistrationResponse, error) {
	var validationData SaveValidationCode

	uvCol := database.GetCollection("userValidationCode")

	_, err := uvCol.ReadDocument(context.Background(), phoneNumber, &validationData)
	if err != nil {
		return loginAndRegistrationResponse{}, err
	}
	if validationData.Code != validationCode {

		return loginAndRegistrationResponse{}, fmt.Errorf("wrong validation code or phone number")

	}
	if validationData.CreatedAt < time.Now().Unix() {
		return loginAndRegistrationResponse{}, fmt.Errorf("validation code is expired")

	}
	_, err = uvCol.RemoveDocument(context.Background(), phoneNumber)

	if err != nil {
		return loginAndRegistrationResponse{}, err

	}
	u, err := getUserByPhoneNumber(phoneNumber)
	if err != nil {
		return loginAndRegistrationResponse{}, err

	}
	fmt.Println(11111)
	p := jwt.TokenPayload{
		Key: u.Key,
	}
	accessToken := jwt.GenerateAccessToken(&p)
	refreshToken := jwt.GenerateRefreshToken(&p)
	hashRefreshToken, _ := password.HashPassword(refreshToken)
	urt := UpdateRefreshToken{
		HashRefreshToken: hashRefreshToken,
		LastLogin:        time.Now().Unix(),
	}
	userCol := database.GetCollection("users")
	var uOut UserOut
	ctx := driver.WithReturnNew(context.Background(), &uOut)
	_, err = userCol.UpdateDocument(ctx, u.Key, urt)
	if err != nil {
		log.Println(err)
		return loginAndRegistrationResponse{}, err
	}
	lr := loginAndRegistrationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         uOut,
	}
	return lr, nil

}

// registerWithValidationCode user registration with validation code
// @Summary user registration with validation code
// @Description user registration with validation code
// @Tags user Auth
// @Accept json
// @Produce json
// @Param data body LoginDto true "data"
// @Success 200 {object} loginAndRegistrationResponse{}
// @Failure 400 {object} string{}
// @Failure 401 {object} string{}
// @Router /user-auth/register [post]
func registerWithValidationCode(phoneNumber string, validationCode string) (loginAndRegistrationResponse, error) {
	var validationData SaveValidationCode

	uvCol := database.GetCollection("userValidationCode")

	_, err := uvCol.ReadDocument(context.Background(), phoneNumber, &validationData)
	if err != nil {
		return loginAndRegistrationResponse{}, err
	}
	if validationData.Code != validationCode {
		return loginAndRegistrationResponse{}, fmt.Errorf("wrong validation code or phone number")

	}
	if validationData.CreatedAt < time.Now().Unix() {
		return loginAndRegistrationResponse{}, fmt.Errorf("validation code is expired")

	}
	_, err = uvCol.RemoveDocument(context.Background(), phoneNumber)

	if err != nil {
		return loginAndRegistrationResponse{}, err

	}
	userCol := database.GetCollection("users")
	u := user{PhoneNumber: phoneNumber}
	meta, err := userCol.CreateDocument(context.Background(), u)
	if err != nil {
		return loginAndRegistrationResponse{}, err

	}
	fmt.Println(11111)
	p := jwt.TokenPayload{
		Key: meta.Key,
	}
	accessToken := jwt.GenerateAccessToken(&p)
	refreshToken := jwt.GenerateRefreshToken(&p)
	hashRefreshToken, _ := password.HashPassword(refreshToken)
	urt := UpdateRefreshToken{
		HashRefreshToken: hashRefreshToken,
		LastLogin:        time.Now().Unix(),
	}

	var uOut UserOut
	ctx := driver.WithReturnNew(context.Background(), &uOut)
	_, err = userCol.UpdateDocument(ctx, meta.Key, urt)
	if err != nil {
		log.Println(err)
		return loginAndRegistrationResponse{}, err
	}
	lr := loginAndRegistrationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         uOut,
	}
	return lr, nil
}

// getRefreshToken  get access token by sending refresh token
// @Summary get access token by sending refresh token
// @Description get access token by sending refresh token
// @Tags user Auth
// @Accept json
// @Produce json
// @Param   token      path   string     true  "token"
// @Success 200 {object} string{}
// @Failure 401 {object} string{}
// @Router /user-auth/get-refresh-token/{token} [get]
func getRefreshToken(c *fiber.Ctx) error {
	t := c.Params("token")
	log.Println(11111, t)
	// CheckPasswordHash the token which is in the chunks
	payload, err := jwt.VerifyRefreshToken(t)
	if err != nil {
		log.Println(1, err)
		return fiber.ErrUnauthorized
	}

	se, err := GetUserByKey(payload.Key)
	if err != nil {
		log.Println(2)
		return fiber.ErrUnauthorized
	}
	match := password.CheckPasswordHash(t, se.HashRefreshToken)
	p := jwt.TokenPayload{
		Key: payload.Key,
	}
	accessToken := jwt.GenerateAccessToken(&p)
	if match {
		return c.JSON(fiber.Map{"token": accessToken})
	}
	log.Println(3)
	return fiber.ErrUnauthorized
}
