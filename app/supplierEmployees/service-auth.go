package supplierEmployees

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/jwt"
	"bamachoub-backend-go-v1/utils/password"
	"bamachoub-backend-go-v1/utils/sms"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
	"time"
)

// GetValidationCode send validation code by sms to given phoneNumber
// @Summary send validation code by sms to given phoneNumber
// @Description send validation code by sms to given phoneNumber for supplier registration
// @Tags supplierAuth
// @Accept json
// @Produce json
// @Param data body getValidationCodeDto true "phoneNumber"
// @Success 200 {object} ResponseHTTP{}
// @Failure 400 {object} ResponseHTTP{}
// @Router /supplier-employee-auth/get-validation-code   [post]
func GetValidationCode(c *fiber.Ctx) error {
	gvc := new(getValidationCodeDto)
	if err := utils.ParseBodyAndValidate(c, gvc); err != nil {
		return c.JSON(err)
	}
	rn := utils.GenRandomNUmber(1001, 9999)
	gvCol := database.GetCollection("EmployeeValidationCode")
	svc := SaveValidationCode{
		Key:       gvc.PhoneNumber,
		Code:      fmt.Sprintf("%v", rn),
		CreatedAt: time.Now().Add(time.Second * 120).Unix(),
	}
	log.Println(time.Now().Unix())
	_, err := gvCol.CreateDocument(context.Background(), svc)
	if err != nil {
		//return c.JSON(err)
		return utils.CustomErrorResponse(409, 1, err, "", c)
	}
	pArr := sms.ParameterArray{
		Parameter:      "VerificationCode",
		ParameterValue: fmt.Sprintf("%v", rn),
	}
	sms.SendSms(gvc.PhoneNumber, "48985", []sms.ParameterArray{pArr})
	//sms.SendSms(gvc.PhoneNumber, fmt.Sprintf("code:%v\nبا استفاده از کد زیر وارد شوید", rn))
	log.Println(rn)

	time.AfterFunc(120*time.Second, func() {

		_, err = gvCol.RemoveDocument(context.Background(), gvc.PhoneNumber)
	})

	return c.JSON(fiber.Map{
		"phoneNumber": gvc.PhoneNumber,
	})
}

// CheckValidationCode  check validation code
// @Summary check validation code
// @Description check validation code and if it is valid, send status:ok
// @Tags supplierAuth
// @Accept json
// @Produce json
// @Param data body checkValidationCodeDto true "phoneNumber and validation code"
// @Success 200 {object} ResponseHTTP{}
// @Failure 400 {object} CustomErrorResponse{}
// @Failure 401 {object} CustomErrorResponse{}
// @Router /supplier-employee-auth/check-validation-code  [post]
func CheckValidationCode(c *fiber.Ctx) error {
	cvc := new(checkValidationCodeDto)
	if err := utils.ParseBodyAndValidate(c, cvc); err != nil {
		return c.JSON(err)
	}
	gvCol := database.GetCollection("EmployeeValidationCode")
	var validationData SaveValidationCode

	_, err := gvCol.ReadDocument(context.Background(), cvc.PhoneNumber, &validationData)
	if err != nil {
		ce := utils.CError{
			Code:    1,
			Error:   fmt.Sprintf("%v", err),
			DevInfo: "problem while reading validation code from database",
			UserMsg: "کد تایید مورد نظر یافت نشد",
		}

		return c.JSON(ce)
	}
	if validationData.Code != cvc.Code {
		ce := utils.CError{
			Code:    2,
			Error:   fmt.Sprint("validation code does not match"),
			DevInfo: "validation code != sent Code",
			UserMsg: "کد تایید اشتباه است",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(ce)
	}
	if validationData.CreatedAt < time.Now().Unix() {
		ce := utils.CError{
			Code:    3,
			Error:   fmt.Sprint("validation code is expired"),
			DevInfo: "",
			UserMsg: "کد تایید فاسد شده است",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(ce)
	}

	_, err = gvCol.RemoveDocument(context.Background(), cvc.PhoneNumber)

	if err != nil {
		return c.JSON(err)
	}
	//t := jwt.TokenPayload{
	//	Id: cvc.PhoneNumber,
	//}
	//token := jwt.GenerateAccessToken(&t)

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

// CreateSupplierPreview create supplier preview
// @Summary create supplier preview
// @Description create supplier preview for management
// @Tags supplierAuth
// @Accept json
// @Produce json
// @Param data body createSupplierPreview true "data"
// @Success 200 {object} createSupplierPreview{}
// @Failure 400 {object} CustomErrorResponse{}
// @Failure 401 {object} CustomErrorResponse{}
// @Router /supplier-employee-auth/create-supplier-preview [post]
func CreateSupplierPreview(data createSupplierPreview) (*createSupplierPreviewIn, error) {

	s := createSupplierPreviewIn{
		FirstName:          data.FirstName,
		LastName:           data.LastName,
		PhoneNumber:        data.PhoneNumber,
		Email:              data.Email,
		NationalCode:       data.NationalCode,
		BirthDate:          data.BirthDate,
		ShabaNumber:        data.ShabaNumber,
		ShopName:           data.ShopName,
		Latitude:           data.Latitude,
		Longitude:          data.Longitude,
		State:              data.State,
		City:               data.City,
		Address:            data.Address,
		PostalCode:         data.PostalCode,
		CategoriesToSale:   data.CategoriesToSale,
		IdCardImage:        data.IdCardImage,
		IdBookPageOneImage: data.IdBookPageOneImage,
		IdBookPageTwoImage: data.IdBookPageTwoImage,
		SalesPermitImage:   data.SalesPermitImage,
		CreateAt:           time.Now().Unix(),
	}
	spCol := database.GetCollection("supplierPreview")
	var eResp createSupplierPreviewIn
	ctx := driver.WithReturnNew(context.Background(), &eResp)
	_, err := spCol.CreateDocument(ctx, s)
	if err != nil {
		return nil, err
	}
	return &eResp, nil
}

func getSupplierPreview(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	q := fmt.Sprintf("for s in supplierPreview limit %v,%v return s", offset, limit)
	return c.JSON(database.ExecuteGetQuery(q))
}

func getSupplierPreviewByKey(key string) (*supplierPreview, error) {
	var sp supplierPreview
	spCol := database.GetCollection("supplierPreview")
	_, err := spCol.ReadDocument(context.Background(), key, &sp)
	if err != nil {
		return nil, err
	}
	return &sp, nil
}

// supplierEmployeeLogin login supplieremployee
// @Summary login supplieremployee
// @Description login supplieremployee  by phoneNumber and password
// @Tags supplierAuth
// @Accept json
// @Produce json
// @Param data body loginRequest true "data"
// @Success 200 {object} loginResponse{}
// @Failure 400 {object} CustomErrorResponse{}
// @Failure 401 {object} CustomErrorResponse{}
// @Router /supplier-employee-auth/login [post]
func supplierEmployeeLogin(phoneNumber string, pass string) (*loginResponse, error) {
	employee, err := GetSupplyEmployeeByPhoneNumber(phoneNumber)
	fmt.Println(3)
	if err != nil {

		return nil, fmt.Errorf("phoneNumber or password is wrong  \n error: %v", err)
	}
	fmt.Println(4)
	fmt.Println(employee.HashPassword, pass)

	match := password.CheckPasswordHash(pass, employee.HashPassword)
	fmt.Println(5)
	if !match {
		return nil, fmt.Errorf("wrong password or phoneNumber")
	}
	fmt.Println(6)
	p := jwt.SETokenPayload{
		Key:         employee.Key,
		SupplierKey: employee.SupplierKey,
		Access:      strings.Join(employee.Access[:], ","),
	}
	fmt.Println(7)
	accessToken := jwt.GenerateSupplierEmployeeToken(&p, false)
	fmt.Println(8)
	refreshToken := jwt.GenerateSupplierEmployeeToken(&p, true)
	fmt.Println(9)
	hashRefreshToken, _ := password.HashPassword(refreshToken)
	fmt.Println(10)
	urt := UpdateRefreshToken{
		HashRefreshToken: hashRefreshToken,
		LastLogin:        time.Now().Unix(),
	}
	fmt.Println(11)
	seCol := database.GetCollection("supplierEmployee")
	fmt.Println(12)
	_, err = seCol.UpdateDocument(context.Background(), employee.Key, urt)
	fmt.Println(13)
	if err != nil {
		return nil, err
	}
	lr := loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsFirstLogin: employee.LastLogin == 0,
		Employee:     employee,
	}
	fmt.Println(14)
	return &lr, nil

}

// changePasswordWithLogin  change password by sending code
// @Summary change password by sending code
// @Description change password by sending code
// @Tags supplierAuth
// @Accept json
// @Produce json
// @Param data body changePasswordWithLoginRequest true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} ResponseHTTP{}
// @Failure 400 {object} ResponseHTTP{}
// @Router /supplier-employee-auth/change-password-with-login [post]
func changePasswordWithLogin(key string, pass string) error {
	hp, _ := password.HashPassword(pass)
	cp := changePassword{HashPassword: hp}
	seCol := database.GetCollection("supplierEmployee")
	_, err := seCol.UpdateDocument(context.Background(), key, cp)
	if err != nil {
		return err
	}
	return nil

}

// GetChangePasswordCode  change password by sending code
// @Summary change password by sending code
// @Description change password by sending code
// @Tags supplierAuth
// @Accept json
// @Produce json
// @Param data body changePasswordWithoutLoginRequest true "data"
// @Success 200 {object} ResponseHTTP{}
// @Failure 400 {object} ResponseHTTP{}
// @Router /supplier-employee-auth/changePassword-without-login [post]
func changePasswordWithoutLogin(c *fiber.Ctx) error {
	cvc := new(changePasswordWithoutLoginRequest)
	if err := utils.ParseBodyAndValidate(c, cvc); err != nil {
		return c.JSON(err)
	}
	gvCol := database.GetCollection("changePasswordCode")
	var validationData SaveValidationCode

	_, err := gvCol.ReadDocument(context.Background(), cvc.PhoneNumber, &validationData)
	if err != nil {
		ce := utils.CError{
			Code:    1,
			Error:   fmt.Sprintf("%v", err),
			DevInfo: "problem while reading validation code from database",
			UserMsg: "کد تایید مورد نظر یافت نشد",
		}

		return c.JSON(ce)
	}
	if validationData.Code != cvc.Code {
		ce := utils.CError{
			Code:    2,
			Error:   fmt.Sprint("validation code does not match"),
			DevInfo: "validation code != sent Code",
			UserMsg: "کد تایید اشتباه است",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(ce)
	}
	if validationData.CreatedAt < time.Now().Unix() {
		ce := utils.CError{
			Code:    3,
			Error:   fmt.Sprint("validation code is expired"),
			DevInfo: "",
			UserMsg: "کد تایید فاسد شده است",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(ce)
	}

	_, err = gvCol.RemoveDocument(context.Background(), cvc.PhoneNumber)
	se, err := GetSupplyEmployeeByPhoneNumber(cvc.PhoneNumber)
	if err != nil {
		ce := utils.CError{
			Code:    4,
			Error:   fmt.Sprintf("%v", err),
			DevInfo: "",
			UserMsg: "کاربر یافت نشد",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(ce)
	}
	err = changePasswordWithLogin(se.Key, cvc.Password)
	if err != nil {

		return c.JSON(err)
	}
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

// GetChangePasswordCode  send change password code by sms to given phoneNumber
// @Summary send change password  code by sms to given phoneNumber
// @Description send change password  code by sms to given phoneNumber for supplier
// @Tags supplierAuth
// @Accept json
// @Produce json
// @Param data body getValidationCodeDto true "phoneNumber"
// @Success 200 {object} ResponseHTTP{}
// @Failure 400 {object} ResponseHTTP{}
// @Router /supplier-employee-auth/get-changePassword-code  [post]
func GetChangePasswordCode(c *fiber.Ctx) error {
	gvc := new(getValidationCodeDto)
	if err := utils.ParseBodyAndValidate(c, gvc); err != nil {
		return c.JSON(err)
	}
	rn := utils.GenRandomNUmber(1001, 9999)
	gvCol := database.GetCollection("changePasswordCode")
	svc := SaveValidationCode{
		Key:       gvc.PhoneNumber,
		Code:      fmt.Sprintf("%v", rn),
		CreatedAt: time.Now().Add(time.Second * 120).Unix(),
	}
	log.Println(time.Now().Unix())
	_, err := gvCol.CreateDocument(context.Background(), svc)
	if err != nil {
		//return c.JSON(err)
		return utils.CustomErrorResponse(409, 1, err, "", c)
	}
	pArr := sms.ParameterArray{
		Parameter:      "VerificationCode",
		ParameterValue: fmt.Sprintf("%v", rn),
	}
	sms.SendSms(gvc.PhoneNumber, "48985", []sms.ParameterArray{pArr})
	//sms.SendSms(gvc.PhoneNumber, fmt.Sprintf("code:%v\nکد برای عوض کردن رمزعبور", rn))
	//log.Println(rn)

	time.AfterFunc(120*time.Second, func() {

		_, err = gvCol.RemoveDocument(context.Background(), gvc.PhoneNumber)
	})

	return c.JSON(fiber.Map{
		"phoneNumber": gvc.PhoneNumber,
	})
}

// getRefreshToken  get access token by sending refresh token
// @Summary get access token by sending refresh token
// @Description get access token by sending refresh token
// @Tags supplierAuth
// @Accept json
// @Produce json
// @Success 200 {object} refreshTokenResponse{}
// @Failure 401 {object} ResponseHTTP{}
// @Router /supplier-employee-auth/get-refresh-token/{token} [get]
func getRefreshToken(c *fiber.Ctx) error {
	t := c.Params("token")
	log.Println(11111, t)
	// CheckPasswordHash the token which is in the chunks
	payload, err := jwt.VerifySupplierEmployee(t, true)
	if err != nil {
		log.Println(1, err)
		return fiber.ErrUnauthorized
	}

	se, err := getSupplierEmployeeByKey(payload.Key)
	if err != nil {
		log.Println(2)
		return fiber.ErrUnauthorized
	}
	match := password.CheckPasswordHash(t, se.HashRefreshToken)
	//log.Println(se.HashPassword)
	//log.Println(t)
	if match {
		return c.JSON(refreshTokenResponse{Token: t})
	}
	log.Println(3)
	return fiber.ErrUnauthorized
}
