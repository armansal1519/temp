package sendingInfo

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// addInterval create sending interval
// @Summary create sending interval
// @Description create sending interval in faq database
// @Tags sending info
// @Accept json
// @Produce json
// @Param addIntervalRequest body addIntervalRequest true "addIntervalRequest"
// @Success 200 {object} []string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /sending-info/add-interval [post]
func addInterval(c *fiber.Ctx) error {
	s := new(addIntervalRequest)
	if err := utils.ParseBodyAndValidate(c, s); err != nil {
		return c.JSON(err)
	}
	days := []string{"monday", "tuesday", "wednesday",
		"thursday", "friday", "saturday", "sunday"}
	if !utils.Contains(days, s.Key) {
		return c.Status(400).SendString(fmt.Sprintf("key must be on if this %v", days))
	}
	Col := database.GetCollection("sendDayInterval")
	flag, err := Col.DocumentExists(context.Background(), s.Key)
	if !flag {
		sdi := sendDayInterval{
			Key:       s.Key,
			Intervals: []interval{s.Interval},
		}
		meta, err := Col.CreateDocument(context.Background(), sdi)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(meta)
	}

	var data sendDayInterval
	_, err = Col.ReadDocument(context.Background(), s.Key, &data)
	if err != nil {

		return c.JSON(err)
	}

	newArr, err := addToInterval(data.Intervals, s.Interval)
	if err != nil {
		return c.JSON(err)
	}
	data.Intervals = newArr
	meta, err := Col.UpdateDocument(context.Background(), data.Key, data)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)

}

// removeInterval remove sending interval
// @Summary remove sending interval
// @Description remove sending interval
// @Tags sending info
// @Accept json
// @Produce json
// @Param addIntervalRequest body addIntervalRequest true "addIntervalRequest"
// @Success 200 {object} []string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /sending-info/remove-interval [post]
func removeInterval(c *fiber.Ctx) error {
	s := new(addIntervalRequest)
	if err := utils.ParseBodyAndValidate(c, s); err != nil {
		return c.JSON(err)
	}
	days := []string{"monday", "tuesday", "wednesday",
		"thursday", "friday", "saturday", "sunday"}
	if !utils.Contains(days, s.Key) {
		return c.Status(400).SendString(fmt.Sprintf("key must be on if this %v", days))
	}
	Col := database.GetCollection("sendDayInterval")
	var data sendDayInterval
	_, err := Col.ReadDocument(context.Background(), s.Key, &data)
	if err != nil {
		return c.JSON(err)
	}
	for i, interval := range data.Intervals {
		if s.Interval.To == interval.To && s.Interval.From == interval.From {
			data.Intervals = remove(data.Intervals, i)
		}
	}
	meta, err := Col.UpdateDocument(context.Background(), data.Key, data)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)

}

// getSendDayInterval get all sending interval
// @Summary return all sending interval
// @Description return all sending interval
// @Tags sending info
// @Accept json
// @Produce json
// @Success 200 {object} []sendDayInterval{}
// @Failure 404 {object} string{}
// @Router /sending-info/interval [get]
func getSendDayInterval(c *fiber.Ctx) error {
	query := fmt.Sprintf("for i in sendDayInterval return i")
	return c.JSON(database.ExecuteGetQuery(query))
}

// CreateSendingInfo create sending info
// @Summary create sending info
// @Description create sending info
// @Tags sending info
// @Accept json
// @Produce json
// @Param sendingInfo body sendingInfo true "sendingInfo"
// @Param key path string true "key"
// @Success 200 {object} []string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /sending-info/info/{key} [post]
func CreateSendingInfo(c *fiber.Ctx) error {
	si := new(sendingInfo)
	userKey := c.Locals("userKey").(string)
	si.UserKey = userKey
	orderKey := c.Params("orderKey")
	if err := utils.ParseBodyAndValidate(c, si); err != nil {
		return c.JSON(err)
	}

	if si.TransportationType != "user-address" && si.TransportationType != "bamachoub" {
		return c.Status(400).SendString("TransportationType must be bamachoub or user-address but is : " + si.TransportationType)
	}
	if si.SendingMethod != "fast" && si.SendingMethod != "normal" {
		return c.Status(400).SendString("SendingMethod must be fast or normal but is : " + si.SendingMethod)
	}
	col := database.GetCollection("sendingInfo")
	meta, err := col.CreateDocument(context.Background(), si)
	if err != nil {
		return c.JSON(err)
	}
	u := updateSendingInfoKey{SendingInfoKey: meta.Key}
	orderCol := database.GetCollection("order")
	orderMeta, err := orderCol.UpdateDocument(context.Background(), orderKey, u)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(orderMeta)
}

// GetSendingInfoByKey get  sending info by key
// @Summary return  sending info by key
// @Description return  sending info by key
// @Tags sending info
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} sendingInfoOut{}
// @Failure 404 {object} string{}
// @Router /sending-info/info/{key} [get]
func GetSendingInfoByKey(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	key := c.Params("key")

	query := fmt.Sprintf("for i in sendingInfo filter i._key==\"%v\" && i.userKey==\"%v\" return i", key, userKey)
	return c.JSON(database.ExecuteGetQuery(query))
}

func getSendingInfoByKey(key string) (*sendingInfoOut, error) {
	col := database.GetCollection("sendingInfo")
	var doc sendingInfoOut
	_, err := col.ReadDocument(context.Background(), key, &doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// editSendingInfo update sending info
// @Summary update sending info
// @Description update sending info
// @Tags sending info
// @Accept json
// @Produce json
// @Param updateSendingInfo body updateSendingInfo true "sendingInfo"
// @Param key path string true "key"
// @Success 200 {object} []string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /sending-info/info/{key} [put]
func editSendingInfo(c *fiber.Ctx) error {
	u := new(updateSendingInfo)
	key := c.Params("key")
	if err := utils.ParseBodyAndValidate(c, u); err != nil {
		return c.JSON(err)
	}
	doc, err := getSendingInfoByKey(key)
	if err != nil {
		return c.JSON(err)
	}
	userKey := c.Locals("userKey").(string)
	if userKey != doc.UserKey {
		return c.Status(403).SendString("not your data")
	}
	col := database.GetCollection("sendingInfo")
	meta, err := col.UpdateDocument(context.Background(), key, u)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}
