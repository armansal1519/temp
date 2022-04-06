package sendingUnit

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

// createSendUnit create sending send unit
// @Summary create sending send unit
// @Description create sending send unit
// @Tags sending unit
// @Accept json
// @Produce json
// @Param sendUnit body sendUnit true "sendUnit"
// @Success 200 {object} []string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /sending-unit [post]
//func createSendUnit(c *fiber.Ctx) error {
//	su := new(sendUnit)
//	if err := utils.ParseBodyAndValidate(c, su); err != nil {
//		return c.JSON(err)
//	}
//	var ao graphOrder.GOrderItemOut
//	aoCol := database.GetCollection("gOrderItem")
//	_, err := aoCol.ReadDocument(context.Background(), su.ApprovedOrderKey, &ao)
//	if err != nil {
//		return c.JSON(err)
//	}
//	if su.Number+ao.ConvertToSendUnit > ao.Number {
//		return c.Status(400).SendString("number of send unit is bigger than order number")
//	}
//
//	su.Status = "processing"
//
//	col := database.GetCollection("sendingUnit")
//	meta, err := col.CreateDocument(context.Background(), su)
//	if err != nil {
//		return c.JSON(err)
//	}
//
//	u := updateApprovedOrder{ConvertToSendUnit: ao.ConvertToSendUnit + su.Number}
//	meta, err = col.UpdateDocument(context.Background(), ao.Key, u)
//	if err != nil {
//		return c.JSON(err)
//	}
//	return c.JSON(meta)
//}

// getSendUnits get all sending unit
// @Summary return all sending unit
// @Description return all sending unit
// @Tags sending unit
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Param not-send  query bool    true  "true or false"
// @Success 200 {object} []sendUnitOut{}
// @Failure 404 {object} string{}
// @Router /sending-unit [get]
func getSendUnits(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")

	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	notSend := c.Query("not-send")
	nsString := ""
	if notSend == "true" {
		nsString = "filter i.transportationUnitKey==\"\""
	}
	query := fmt.Sprintf("for i in sendingUnit %v sort i.createdAt limit %v,%v return i", nsString, offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}

func getSendUnitsByUserKey(c *fiber.Ctx) error {
	userKey := ""
	isAdmin := c.Locals("isAdmin").(bool)
	if isAdmin {
		userKey = c.Query("user-key")
	} else {
		userKey = c.Locals("userKey").(string)

	}
	query := fmt.Sprintf("for i in approvedOrder filter i.userKey==\"%v\"\nlet su = (for j in sendingUnit filter j.approvedOrderKey==i.key return j)\nreturn {ao:i,su:su}", userKey)
	return c.JSON(database.ExecuteGetQuery(query))
}

// addOrRemoveSendUnits add or remove transportationKey to a send unit
// @Summary add or remove transportationKey to a send unit
// @Description add or remove transportationKey to a send unit
// @Tags sending unit
// @Accept json
// @Produce json
// @Param op path string true "op"
// @Param unitKey path string true "unitKey"
// @Param trKey path string true "trKey"
// @Success 200 {object} []string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /sending-unit/{op}/{unitKey}/{trKey} [put]
func addOrRemoveSendUnits(c *fiber.Ctx) error {
	op := c.Params("op")
	unitKey := c.Params("unitKey")
	trKey := c.Params("trKey")
	var new sendUnitOut
	ctx := driver.WithReturnNew(context.Background(), &new)
	col := database.GetCollection("sendingUnit")
	if op == "add" {
		u := updateTr{TransportationUnitKey: trKey}
		_, err := col.UpdateDocument(ctx, unitKey, u)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(new)
	} else if op == "remove" {
		u := updateTr{TransportationUnitKey: ""}
		_, err := col.UpdateDocument(ctx, unitKey, u)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(new)
	} else {
		return c.Status(400).SendString("op only can be add or remove")
	}
}

// removeSendUnit delete send unit
// @Summary delete send unit
// @Description delete send unit
// @Tags sending unit
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} []string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /sending-unit/{key} [delete]
//func removeSendUnit(c *fiber.Ctx) error {
//	key := c.Params("key")
//
//	var su sendUnitOut
//	col := database.GetCollection("sendingUnit")
//	_, err := col.ReadDocument(context.Background(), key, &su)
//	if err != nil {
//		return c.JSON(err)
//	}
//	if su.TransportationUnitKey != "" {
//		return c.Status(409).SendString("send to user")
//	}
//	var ao paymentAndWallet.ApprovedOrderOut
//	aoCol := database.GetCollection("approvedOrder")
//	_, err = aoCol.ReadDocument(context.Background(), su.ApprovedOrderKey, &ao)
//	if err != nil {
//		return c.JSON(err)
//	}
//	u := updateApprovedOrder{ConvertToSendUnit: ao.ConvertToSendUnit - su.Number}
//	_, err = col.UpdateDocument(context.Background(), ao.Key, u)
//	if err != nil {
//		return c.JSON(err)
//	}
//	meta, err := col.RemoveDocument(context.Background(), key)
//	if err != nil {
//		return c.JSON(err)
//	}
//	return c.JSON(meta)
//}

// removeTr delete transportation unit from send unit
// @Summary delete transportation unit from send unit
// @Description delete transportation unit from send unit
// @Tags sending unit
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} []string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /sending-unit/tr/{key} [delete]
func removeTr(c *fiber.Ctx) error {
	key := c.Params("key")
	query := fmt.Sprintf("for i in sendingUnit filter i.transportationUnitKey==\"%v\" update i with {transportationUnitKey:\"\"} in sendingUnit return New", key)
	return c.JSON(database.ExecuteGetQuery(query))

}
