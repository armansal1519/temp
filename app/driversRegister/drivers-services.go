package driversRegister

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

// CreateDriversInfo  create drivers
// @Summary  create drivers information
// @Description create drivers and their information
// @Tags driverInfo
// @Accept json
// @Produce json
// @Param data body driverInfo true "data"
// @Success 200 {object} driverInfo{}
// @Failure 404 {object} driverInfo{}
// @Router /drivers [post]
func CreateDriversInfo(data driverInfo) (*driverInfo, error) {
	d := driverInfo{
		FirstName:            data.FirstName,
		LastName:             data.LastName,
		Gender:               data.Gender,
		NationalNo:           data.NationalNo,
		BirthdayDate:         data.BirthdayDate,
		City:                 data.City,
		Province:             data.Province,
		Address:              data.Address,
		Latitude:             data.Latitude,
		Longitude:            data.Longitude,
		PostCode:             data.PostCode,
		PhoneNo:              data.PhoneNo,
		CardNo:               data.CardNo,
		AccountNo:            data.AccountNo,
		CarType:              data.CarType,
		PlateNo:              data.PlateNo,
		AppearanceStatus:     data.AppearanceStatus,
		CarTechnical:         data.CarTechnical,
		InsuranceStatus:      data.InsuranceStatus,
		BarbandStatus:        data.BarbandStatus,
		CarColor:             data.CarColor,
		CarTools:             data.CarTools,
		DriverEthics:         data.DriverEthics,
		Sabeghe:              data.Sabeghe,
		PhysicalCondition:    data.PhysicalCondition,
		Panctuality:          data.Panctuality,
		IdCardImage:          data.IdCardImage,
		IdBookPageOneImage:   data.IdBookPageOneImage,
		IdBookPageTwoImage:   data.IdBookPageTwoImage,
		LicenseImage:         data.LicenseImage,
		FaceImage:            data.FaceImage,
		InsurancePolicyImage: data.InsurancePolicyImage,
		CreatedAt:            time.Now().Unix(),
		Description:          data.Description,
	}
	driverCol := database.GetCollection("driverInfo")
	var drInf driverInfo
	ctx := driver.WithReturnNew(context.Background(), &drInf)
	_, err := driverCol.CreateDocument(ctx, d)
	if err != nil {
		return nil, err
	}
	return &drInf, nil
}

// getAllDriversInfo get all drivers information
// @Summary return all informations of drivers
// @Description return all drivers information
// @Tags driverInfo
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Success 200 {object} driverInfo
// @Failure 404 {object} string{}
// @Router /drivers [get]
func getAllDriversInfo(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	query := fmt.Sprintf("FOR d IN driverInfo SORT d.createdAt LIMIT %v, %v RETURN d", offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}

// getDriversInfoByKey get driver information by key
// @Summary return information of driver by its key
// @Description return driver's information by given key
// @Tags driverInfo
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} driverInfo
// @Failure 404 {object} string{}
// @Router /drivers/{key} [get]
func getDriversInfoByKey(c *fiber.Ctx) error {
	key := c.Params("key")
	ctx := context.Background()

	driverCol := database.GetCollection("driverInfo")
	var doc getDriverInfo
	_, err := driverCol.ReadDocument(ctx, key, &doc)
	if err != nil {
		return c.Status(404).SendString("Document Not Founded")
	}
	return c.JSON(doc)
}

// deleteDriversInfoByKey  delete driver by its key
// @Summary  delete driver information
// @Description  delete driver's information by its key
// @Tags driverInfo
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} driverInfo
// @Failure 404 {object} string{}
// @Router /drivers/{key} [delete]
func deleteDriversInfoByKey(c *fiber.Ctx) error {
	key := c.Params("key")
	ctx := context.Background()

	driverCol := database.GetCollection("driverInfo")
	var doc getDriverInfo
	_, err := driverCol.RemoveDocument(ctx, key)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(doc)
	// return c.Status(204).SendString("Successfully Deleted")
}

// updateDriversInfoByKey  update drivers information
// @Summary update information of each driver by its key
// @Description update informations of driver by its key
// @Tags driverInfo
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} driverInfo{}
// @Failure 404 {object} string{}
// @Router /drivers/{key} [put]
func updateDriversInfoByKey(c *fiber.Ctx) error {
	key := c.Params("key")
	edit := new(editDriverInfo)
	if err := utils.ParseBodyAndValidate(c, edit); err != nil {
		c.JSON(err)
	}
	driverCol := database.GetCollection("driverInfo")
	var update driverInfo
	ctx := driver.WithReturnNew(context.Background(), &update)
	_, err := driverCol.UpdateDocument(ctx, key, edit)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(update)
}

// searchIntoDrivers search into drivers
// @Summary return drivers by its name
// @Description return drivers with the given name in search box
// @Tags driverInfo
// @Accept json
// @Produce json
// @Param input  body  input true  "input"
// @Param offset query int   true  "Offset"
// @Param limit  query int   true  "limit"
// @Success 200 {object} driverInfo{}
// @Failure 404 {object} string{}
// @Router /drivers/search [post]
func searchIntoDrivers(c *fiber.Ctx) error {
	input := new(input)
	if err := utils.ParseBodyAndValidate(c, input); err != nil {
		return c.JSON(err)
	}

	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}
	searchString1 := "%" + input.FirstName + "%"
	searchString2 := "%" + input.LastName + "%"
	query := fmt.Sprintf("FOR d IN driverInfo FILTER d.firstName LIKE \"%v\"  ||  d.lastName LIKE \"%v\" SORT d.createdAt LIMIT %v, %v RETURN d ", searchString1, searchString2, offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}

// filterDriversBaseOnCarType filter drivers base on their car types
// @Summary return drivers that has the specific car type
// @Description return drivers that has the car type same as given
// @Tags driverInfo
// @Accept json
// @Produce json
// @Param input  body  input true  "input"
// @Param offset query int   true  "Offset"
// @Param limit  query int   true  "limit"
// @Success 200 {object} driverInfo{}
// @Failure 404 {object} string{}
// @Router /drivers/filter [post]
func filterDriversBaseOnCarType(c *fiber.Ctx) error {
	input := new(carType)
	if err := utils.ParseBodyAndValidate(c, input); err != nil {
		return c.JSON(err)
	}

	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	query := fmt.Sprintf("FOR d IN driverInfo FILTER d.carType == \"%v\" SORT d.createdAt  LIMIT %v, %v RETURN d", input.Type, offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}
