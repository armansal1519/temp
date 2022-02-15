package admin

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/password"
	"bamachoub-backend-go-v1/utils/sms"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)


// createAdmin  create admin
// @Summary create admin
// @Description create admin
// @Tags admin
// @Accept json
// @Produce json
// @Param data body adminIn true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /admin [post]
func createAdmin(c *fiber.Ctx) error {
	ai := new(adminIn)
	if err := utils.ParseBodyAndValidate(c, ai); err != nil {
		return c.JSON(err)
	}
	adminCol := database.GetCollection("admin")
	rn := utils.GenRandomNUmber(100001, 999999)
	hashPassword := password.Generate(fmt.Sprintf("%v", rn))
	ai.HashPassword = hashPassword
	ai.CreateAt = time.Now().Unix()
	ai.Status = "ok"
	meta, err := adminCol.CreateDocument(context.Background(), ai)
	if err != nil {
		return c.JSON(err)
	}
	pArr := sms.ParameterArray{
		Parameter:      "VerificationCode",
		ParameterValue: fmt.Sprintf("%v", rn),
	}
	sms.SendSms(ai.PhoneNumber, "48985", []sms.ParameterArray{pArr})
	//sms.SendSms(ai.PhoneNumber, fmt.Sprintf("password: %v", rn))
	return c.JSON(meta)
}


// getAdminByKey  get admin by key
// @Summary get admin by key
// @Description get admin by key
// @Tags admin
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /admin/{key} [get]
func getAdminByKey(c *fiber.Ctx) error {
	key := c.Params("key")
	adminCol := database.GetCollection("admin")
	var a AdminOut
	_, err := adminCol.ReadDocument(context.Background(), key, &a)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(a)

}

func getAdminByAccessToken(c *fiber.Ctx) error {
	key:=c.Locals("adminKey").(string)
	adminCol := database.GetCollection("admin")
	var a AdminOut
	_, err := adminCol.ReadDocument(context.Background(), key, &a)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(a)
}


// getAll  get all admin
// @Summary get all admin
// @Description get all admin
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /admin [get]
func getAll(c *fiber.Ctx) error {
	query := fmt.Sprintf("for i in admin return i")
	return c.JSON(database.ExecuteGetQuery(query))
}


// updateAdmin  update admin
// @Summary update admin
// @Description update admin
// @Tags admin
// @Accept json
// @Produce json
// @Param data body updateAdminIn true "data"
// @Param key path string true "key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /admin/{key} [put]
func updateAdmin(c *fiber.Ctx) error {
	uai := new(updateAdminIn)
	key := c.Params("key")
	if err := utils.ParseBodyAndValidate(c, uai); err != nil {
		return c.JSON(err)
	}
	adminCol := database.GetCollection("admin")
	meta, err := adminCol.UpdateDocument(context.Background(), key, uai)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

func GetAdminByPhoneNumber(phoneNumber string) (*AdminOut, error) {
	db := database.GetDB()
	sw := &AdminOut{}
	query := fmt.Sprintf("for i in admin filter i.phoneNumber==\"%v\" limit 1 return i", phoneNumber)
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {

		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	_, err = cursor.ReadDocument(ctx, sw)
	if err != nil {
		log.Print(1234, err)
		return nil, err
	}
	return sw, nil
}


// getAccessArray  get all accesses
// @Summary get all accesses
// @Description get all accesses
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /admin/access [get]
func getAccessArray(c *fiber.Ctx)error{
	return c.JSON(getAllAdminAccess())
}