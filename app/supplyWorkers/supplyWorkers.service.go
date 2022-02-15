package supplyWorkers

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CreateSupplyManager(c *fiber.Ctx) error {
	sm := new(SupplyManagerCreateRequest)
	if err := c.BodyParser(sm); err != nil {
		return err
	}
	errors := utils.Validate(sm)
	if errors != nil {
		c.JSON(errors)
		return nil
	}

	doesWareHouseExist := checkWarehouseExists(sm.SupplierKey)
	if !doesWareHouseExist {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("warehouse whit Id: %v does not exist", sm.SupplierKey))
	}
	newSupplyManager := SupplyWorker{
		FullName:           sm.FullName,
		PhoneNumber:        sm.PhoneNumber,
		Email:              sm.Email,
		NationalCode:       sm.NationalCode,
		Role:               "manager",
		Access:             []string{"manager"},
		SupplierKeyArray:   []string{sm.SupplierKey},
		CurrentSupplierKey: sm.SupplierKey,
		FirstTimeLogin:     true,
	}
	meta := saveSupplyManager(newSupplyManager)
	return c.JSON(meta)

}

func GetSupplyWorkerByKey(key string) getSupplyWorkerByKeyType {
	col := database.GetCollection("supplyWorkers")
	var sw getSupplyWorkerByKeyType
	meta, err := col.ReadDocument(context.Background(), key, &sw)
	if err != nil {
		panic(fmt.Sprintf("error getting SupplyWorker by key %s", key))
	}
	fmt.Println(meta)
	return sw
}

func GetSupplyWorkerByPhoneNumber(phoneNumber string) (*getSupplyWorkerByKeyType, error) {
	db := database.GetDB()
	sw := &getSupplyWorkerByKeyType{}
	//var sw getSupplyWorkerByKeyType
	query := fmt.Sprintf("for i in supplyWorkers\nfilter i.phoneNumber==\"%v\"\nlimit 1\nreturn i", phoneNumber)
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	_, err = cursor.ReadDocument(ctx, sw)
	if err != nil {
		return nil, err
	}
	return sw, nil
}
