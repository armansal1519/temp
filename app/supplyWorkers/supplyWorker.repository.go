package supplyWorkers

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
)

func saveSupplyManager(sm SupplyWorker) driver.DocumentMeta {
	col := database.GetCollection("supplyWorkers")
	ctx := context.Background()
	meta, err := col.CreateDocument(ctx, sm)
	if err != nil {
		panic(fmt.Sprintf("error creating document at saveSupplyManager:%v", err))
	}
	return meta
}

func checkWarehouseExists(supplierKey string) bool {
	col := database.GetCollection("suppliers")
	flag, err := col.DocumentExists(context.Background(), supplierKey)
	if err != nil {
		panic(fmt.Sprintf("error while varifying document exists :%v \n%v", supplierKey, err))
	}
	return flag

}
