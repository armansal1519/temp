package sendingUnit

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
)

func acceptSendUnit(c *fiber.Ctx) error {
	key := c.Params("key")
	sendUnitCol := database.GetCollection("sendUnit")
	var su sendUnitOut
	ctx := driver.WithReturnNew(context.Background(), &su)
	u := updateStatus{Status: "received"}
	meta, err := sendUnitCol.UpdateDocument(ctx, key, &u)
	if err != nil {
		return c.JSON(err)
	}
	query := fmt.Sprintf("for i in sendingUnit filter i.approvedOrderKey ==\"%v\" return  i", su.ApprovedOrderKey)
	db := database.GetDB()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []sendUnitOut
	for {
		var doc sendUnitOut
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}

	flag := true
	for _, datum := range data {
		if datum.Status != "received" {
			flag = false
		}
	}
	if flag {
		addToSupplierWallet()
	}
	return c.JSON(meta)
}

func addToSupplierWallet() {
	log.Println("money add to supplier wallet")
}
