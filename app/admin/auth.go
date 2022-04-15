package admin

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/jwt"
	"bamachoub-backend-go-v1/utils/password"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
	"time"
)

// Login login admin
// @Summary login admin
// @Description login admin  by phoneNumber and password
// @Tags adminAuth
// @Accept json
// @Produce json
// @Param data body loginRequest true "data"
// @Security ApiKeyAuth
// @Failure 400 {object} string{}
// @Failure 401 {object} string{}
// @Router /admin-auth/login [post]
func Login(phoneNumber string, pass string) (*loginResponse, error) {
	admin, err := GetAdminByPhoneNumber(phoneNumber)

	log.Println(1111)

	if err != nil {

		return nil, fmt.Errorf("phoneNumber or password is wrong  \n error: %v", err)
	}
	match := password.CheckPasswordHash(pass, admin.HashPassword)
	if !match {
		return nil, fmt.Errorf("wrong password or phoneNumber")
	}
	p := jwt.AdminTokenPayload{
		Key:    admin.Key,
		Access: strings.Join(admin.Access[:], ","),
	}
	log.Println(1111)

	accessToken := jwt.GenerateAdminToken(&p, false)
	log.Println(accessToken)
	refreshToken := jwt.GenerateAdminToken(&p, true)
	hashRefreshToken, _ := password.HashPassword(refreshToken)
	urt := UpdateRefreshToken{
		HashRefreshToken: hashRefreshToken,
		LastLogin:        time.Now().Unix(),
	}
	seCol := database.GetCollection("admin")
	_, err = seCol.UpdateDocument(context.Background(), admin.Key, urt)
	if err != nil {
		return nil, err
	}
	lr := loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Admin:        *admin,
	}
	return &lr, nil

}

// changePassword changePassword admin
// @Summary changePassword admin
// @Description changePassword admin  by oldPassword and newPassword
// @Tags adminAuth
// @Accept json
// @Produce json
// @Param data body changePasswordIn true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []string{}
// @Failure 400 {object} string{}
// @Failure 401 {object} string{}
// @Router /admin-auth/change-password [post]
func changePassword(c *fiber.Ctx) error {
	p := new(changePasswordIn)
	if err := utils.ParseBodyAndValidate(c, p); err != nil {
		return c.JSON(err)
	}
	adminKey := c.Locals("adminKey").(string)
	adminCol := database.GetCollection("admin")
	isSuperAdmin := c.Locals("isSuperAdmin").(bool)
	var a AdminOut
	_, err := adminCol.ReadDocument(context.Background(), adminKey, &a)
	if err != nil {
		return c.JSON(err)
	}
	if !isSuperAdmin {
		if !password.CheckPasswordHash(p.OldPassword, a.HashPassword) {
			return c.Status(403).SendString("wrong password")
		}
	}
	newHashPassword, _ := password.HashPassword(p.NewPassword)
	u := updatePassword{HashPassword: newHashPassword}
	meta, err := adminCol.UpdateDocument(context.Background(), adminKey, u)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)

}

// getRefreshToken  get access token by sending refresh token
// @Summary get access token by sending refresh token
// @Description get access token by sending refresh token
// @Tags adminAuth
// @Accept json
// @Produce json
// @Success 200 {object} string{}
// @Failure 401 {object} string{}
// @Router /admin-auth/get-refresh-token/{token} [get]
func getRefreshToken(c *fiber.Ctx) error {
	t := c.Params("token")
	//log.Println(11111, t)
	// CheckPasswordHash the token which is in the chunks
	payload, err := jwt.VerifyAdmin(t, true)
	if err != nil {
		log.Println(1, err)
		return fiber.ErrUnauthorized
	}

	adminCol := database.GetCollection("admin")
	var a AdminOut
	_, err = adminCol.ReadDocument(context.Background(), payload.Key, &a)

	if err != nil {
		log.Println(2)
		return fiber.ErrUnauthorized
	}
	match := password.CheckPasswordHash(t, a.HashRefreshToken)

	accessToken := jwt.GenerateAdminToken(payload, false)
	if match {
		return c.JSON(fiber.Map{"Token": accessToken})
	}
	log.Println(3)
	return fiber.ErrUnauthorized
}
