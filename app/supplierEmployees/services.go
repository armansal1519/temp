package supplierEmployees

import (
	"bamachoub-backend-go-v1/app/suppliers"
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

func GetSupplyEmployeeByPhoneNumber(phoneNumber string) (employee, error) {
	db := database.GetDB()
	var sw employee
	query := fmt.Sprintf("for i in supplierEmployee\nfilter i.phoneNumber==\"%v\"\nlimit 1\nreturn i", phoneNumber)
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {

		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	_, err = cursor.ReadDocument(ctx, &sw)
	if err != nil {
		log.Print(1234, err)
		return employee{}, err
	}
	return sw, nil
}

// getSupplierEmployeeByKey get supplier employee  by key
// @Summary get supplier employee  by key
// @Description get supplier employee  by key
// @Tags  supplier
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 400 {object} employee{}
// @Router /supplier-employee/one [get]
func getSupplierEmployeeByKey(key string) (*employee, error) {
	var se employee
	seCol := database.GetCollection("supplierEmployee")
	_, err := seCol.ReadDocument(context.Background(), key, &se)
	if err != nil {
		return nil, err
	}
	return &se, nil

}

func getSupplierAndEmployeeByKey(c *fiber.Ctx) error {
	key := c.Locals("key").(string)
	supplierKey := c.Locals("supplierKey").(string)
	var se employee
	seCol := database.GetCollection("supplierEmployee")
	sCol := database.GetCollection("suppliers")
	_, err := seCol.ReadDocument(context.Background(), key, &se)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	var s suppliers.Supplier
	_, err = sCol.ReadDocument(context.Background(), supplierKey, &s)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.JSON(fiber.Map{
		"supplier": s,
		"employee": se,
	})
}

func createSupplierEmployeeFromSupplierPreview(spKey string) error {
	sp, err := getSupplierPreviewByKey(spKey)
	if err != nil {
		return err
	}

	supplier := suppliers.SupplierIn{
		Address:          sp.Address,
		State:            sp.State,
		City:             sp.City,
		Latitude:         sp.Latitude,
		Longitude:        sp.Longitude,
		Name:             sp.ShopName,
		Code:             "",
		Area:             0,
		AreaWithRoof:     0,
		PhoneNumber:      sp.PhoneNumber,
		Status:           sp.State,
		CreateAt:         time.Now().Unix(),
		CategoriesToSale: sp.CategoriesToSale,
	}
	supplierCol := database.GetCollection("suppliers")
	supplierMeta, err := supplierCol.CreateDocument(context.Background(), supplier)
	if err != nil {
		return err
	}

	tempPass := utils.GenRandomNUmber(1000001, 100000000)
	tempHashPass, _ := password.HashPassword(fmt.Sprintf("%v", tempPass))
	e := employeeIn{
		FirstName:          sp.FirstName,
		LastName:           sp.LastName,
		PhoneNumber:        sp.PhoneNumber,
		ShenasNameCode:     sp.ShenasNameCode,
		Email:              sp.Email,
		PostalCode:         sp.PostalCode,
		NationalCode:       sp.NationalCode,
		BirthDate:          sp.BirthDate,
		ShabaNumber:        sp.ShabaNumber,
		IdCardImage:        sp.IdCardImage,
		IdBookPageOneImage: sp.IdBookPageOneImage,
		IdBookPageTwoImage: sp.IdBookPageTwoImage,
		SalesPermitImage:   sp.SalesPermitImage,
		Access:             []string{"handle-products", "handle-reports", "handle-orders", "handle-wallet"},
		Role:               "manager",
		HashPassword:       tempHashPass,
		HashRefreshToken:   "",
		SupplierKey:        supplierMeta.Key,
		CreateAt:           time.Now().Unix(),
		LastLogin:          0,
		Status:             "ok",
	}
	seCol := database.GetCollection("supplierEmployee")
	_, err = seCol.CreateDocument(context.Background(), e)
	if err != nil {
		return err
	}
	pArr := sms.ParameterArray{
		Parameter:      "VerificationCode",
		ParameterValue: fmt.Sprintf("%v", tempPass),
	}
	sms.SendSms(sp.PhoneNumber, "48985", []sms.ParameterArray{pArr})
	//sms.SendSms(sp.PhoneNumber, fmt.Sprintf("you temp password: %v", tempPass))

	return nil

}

// addToUpdatePool add supplier employee to update pool
// @Summary add supplier employee to update pool
// @Description add supplier employee to update pool
// @Tags  supplier
// @Accept json
// @Produce json
// @Param data body updateEmployee true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 400 {object} string{}
// @Router /supplier-employee/add-update-pool [put]
func addToUpdatePool(c *fiber.Ctx) error {
	data := new(updateEmployee)
	//log.Println(111111,c.Locals("key"))
	if err := utils.ParseBodyAndValidate(c, data); err != nil {
		return c.JSON(err)
	}
	data.CreateAt = time.Now().Unix()
	data.SupplierEmployeeKey = c.Locals("key").(string)
	//updatePoolCol:=database.GetCollection("supplierUpdatePool")
	updatePoolCol := database.GetCollection("supplierEmployeeUpdatePool")
	meta, err := updatePoolCol.CreateDocument(context.Background(), data)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

//func CreateSupplierEmployee(data createSupplierPreview, role string) (*employee, error) {
//
//	e := createSupplierPreviewIn{
//		FirstName:          data.FirstName,
//		LastName:           data.LastName,
//		PhoneNumber:        "",
//		Email:              data.Email,
//		NationalCode:       data.NationalCode,
//		BirthDate:          data.BirthDate,
//		ShabaNumber:        data.ShabaNumber,
//		ShopName:           data.ShopName,
//		Latitude:           data.Latitude,
//		Longitude:          data.Longitude,
//		State:              data.State,
//		City:               data.City,
//		Address:            data.Address,
//		PostalCode:         data.PostalCode,
//		CategoriesToSale:   data.CategoriesToSale,
//		IdCardImage:        data.IdCardImage,
//		IdBookPageOneImage: data.IdBookPageOneImage,
//		IdBookPageTwoImage: data.IdBookPageTwoImage,
//		SalesPermitImage:   data.SalesPermitImage,
//
//		HashPassword:       "",
//		HashRefreshToken:   "",
//		CreateAt:           0,
//		LastLogin:          0,
//	}
//	employeeCol := database.GetCollection("supplierEmployee")
//	var eResp employee
//	ctx := driver.WithReturnNew(context.Background(), &eResp)
//	_, err = employeeCol.CreateDocument(ctx, e)
//	if err != nil {
//		return nil, err
//	}
//	return &eResp, nil
//}

func getAccess(role string, accessList []string) ([]string, error) {

	fullEmployeeAccess := []string{"handle-products", "handle-reports", "handle-orders", "handle-wallet"}
	if role == "admin" || role == "superAdmin" {
		return fullEmployeeAccess, nil
	}
	if len(accessList) == 0 {
		return []string{}, fmt.Errorf("access list is empty")
	}

	return accessList, nil
}

func getRole(role string) (string, error) {
	if role == "admin" || role == "superAdmin" {
		return "manager", nil
	} else if role == "manager" {
		return "employee", nil
	}
	return "", fmt.Errorf("role is not allowed")
}
