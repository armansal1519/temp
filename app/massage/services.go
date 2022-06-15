package massage

import (
	"bamachoub-backend-go-v1/app/admin"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

// addToCart  add to massage
// @Summary adds to massage
// @Description adds to massage
// @Tags massage
// @Accept json
// @Produce json
// @Param data body sendMsgByPhoneNumberReq true "data"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /msg/by-phone [post]
func sendMsgByPhoneNumberUsers(c *fiber.Ctx) error {
	m := new(sendMsgByPhoneNumberReq)
	if err := utils.ParseBodyAndValidate(c, m); err != nil {
		return c.JSON(err)
	}
	adminKey := c.Locals("adminKey").(string)
	adminCol := database.GetCollection("admin")
	var a admin.AdminOut
	_, err := adminCol.ReadDocument(context.Background(), adminKey, &a)
	if err != nil {
		return c.JSON(err)
	}

	msg := massage{
		Title:            m.Title,
		ImageUrl:         m.ImageUrl,
		Text:             m.Text,
		Link:             m.Link,
		Importence:       m.Importence,
		CreatedAt:        time.Now().Unix(),
		CreatedBy:        fmt.Sprintf("%v %v", a.FirstName, a.LastName),
		AdminDescription: m.AdminDescription,
	}
	msgCol := database.GetCollection("massage")

	meta, err := msgCol.CreateDocument(context.Background(), msg)

	pnStr := "["
	for i, key := range m.PhoneNumberArray {
		pnStr += fmt.Sprintf("\"%v\"", key)
		if i < len(m.PhoneNumberArray)-1 {
			pnStr += " , "
		}
	}
	pnStr += "] "
	query := fmt.Sprintf("let userIds=(for u in users filter u.phoneNumber in %v return u._id) for i in userIds\nINSERT { _from: \"%v\", _to: i , seen:false } INTO massageUserEdge  return NEW ", pnStr, meta.ID.String())
	database.ExecuteGetQuery(query)
	return c.JSON(fiber.Map{"msg": "پیام ها با موفقیت ثبت شد"})
}

// sendMsgByPhoneNumberSuppliers  send message for suppliers
// @Summary send message for suppliers
// @Description send message for suppliers
// @Tags massage
// @Accept json
// @Produce json
// @Param data body sendMsgByPhoneNumberReq true "data"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /msg/by-phone-supplier [post]
func sendMsgByPhoneNumberSuppliers(c *fiber.Ctx) error {
	m := new(sendMsgByPhoneNumberReq)
	if err := utils.ParseBodyAndValidate(c, m); err != nil {
		return c.JSON(err)
	}
	adminKey := c.Locals("adminKey").(string)
	adminCol := database.GetCollection("admin")
	var a admin.AdminOut
	_, err := adminCol.ReadDocument(context.Background(), adminKey, &a)
	if err != nil {
		return c.JSON(err)
	}

	msg := massage{
		Title:            m.Title,
		ImageUrl:         m.ImageUrl,
		Text:             m.Text,
		Link:             m.Link,
		Importence:       m.Importence,
		CreatedAt:        time.Now().Unix(),
		CreatedBy:        fmt.Sprintf("%v %v", a.FirstName, a.LastName),
		AdminDescription: m.AdminDescription,
	}
	msgCol := database.GetCollection("massage")

	meta, err := msgCol.CreateDocument(context.Background(), msg)

	pnStr := "["
	for i, key := range m.PhoneNumberArray {
		pnStr += fmt.Sprintf("\"%v\"", key)
		if i < len(m.PhoneNumberArray)-1 {
			pnStr += " , "
		}
	}
	pnStr += "] "
	query := fmt.Sprintf("let userIds=(for u in supplierEmployee filter u.phoneNumber in %v return u._id)\nfor i in userIds INSERT { _from: \"%v\", _to: i , seen:false } INTO massageSupplierEdge \n ", pnStr, meta.ID.String())
	database.ExecuteGetQuery(query)
	return c.JSON(fiber.Map{"msg": "پیام ها با موفقیت ثبت شد"})
}

//func sendMessageAll(c *fiber.Ctx) error {
//	m := new(SendSupplierMassageReq)
//	if err := utils.ParseBodyAndValidate(c, m); err != nil {
//		return c.JSON(err)
//	}
//
//	if m.For != "user" && m.For != "supplier" {
//		return c.Status(400).SendString("For only can be user of supplier")
//	}
//
//	adminKey := c.Locals("adminKey").(string)
//	adminCol := database.GetCollection("admin")
//	var a admin.AdminOut
//	_, err := adminCol.ReadDocument(context.Background(), adminKey, &a)
//	if err != nil {
//		return c.JSON(err)
//	}
//	sm := SupplierMassage{
//		Title:            m.Title,
//		ImageUrl:         m.ImageUrl,
//		Text:             m.Text,
//		State:            m.State,
//		For:              m.For,
//		Link:             m.Link,
//		Importence:       m.Importence,
//		CreatedAt:        time.Now().Unix(),
//		CreatedBy:        fmt.Sprintf("%v %v", a.FirstName, a.LastName),
//		AdminDescription: m.AdminDescription,
//	}
//	msgCol := database.GetCollection("massage")
//
//	meta, err := msgCol.CreateDocument(context.Background(), sm)
//	if m.For == "user" {
//		query := fmt.Sprintf("let s=(for i in suppliers filter i.state==\"%v\" return i._key )\nlet seIds=(for j in supplierEmployee filter j.supplierKey in s return j._id)\nlet seu=UNIQUE(seIds)\nfor k in seu\nINSERT { _from: \"%v\", _to: k , seen:false } INTO massageUserEdge  ", sm.State, meta.ID.String())
//		database.ExecuteGetQuery(query)
//		return c.JSON(fiber.Map{"msg": "پیام ها با موفقیت ثبت شد"})
//	}
//	query := fmt.Sprintf("let u=(for i in userAddress filter i.state==\"%v\" return i.userKey )\nlet seu=UNIQUE(u)\nfor k in seu\nINSERT { _from: \"%v\", _to: concat(\"users/\",k) , seen:false } INTO massageUserEdge  ", m.State, meta.ID.String())
//	database.ExecuteGetQuery(query)
//	return c.JSON(fiber.Map{"msg": "پیام ها با موفقیت ثبت شد"})
//
//}

// getAllMsg get all msg for admin
// @Summary gets all msg for admin
// @Description gets all msg for admin
// @Tags massage
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /msg [get]
func getAllMsg(c *fiber.Ctx) error {
	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" || offset == "" {
		return c.Status(400).SendString("offset or limit is empty")
	}

	q := fmt.Sprintf("for i in massage limit %v,%v return i", offset, limit)
	return c.JSON(database.ExecuteGetQuery(q))

}

// getMassageByUserKey get massages by jwt
// @Summary get massages by jwt
// @Description get massages by jwt , set seen to true if to mark massage as seen
// @Tags massage
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Param   seen      query    bool     false        "seen"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /msg/user [get]
func getMassageByUserKey(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	limit := c.Query("limit")
	offset := c.Query("offset")
	update := c.Query("seen")
	//fmt.Println(offset, limit)

	if limit == "" || offset == "" {
		return c.Status(400).SendString("offset or limit is empty")
	}
	query := ""
	if update == "true" {
		query = fmt.Sprintf("for i in users filter i._key==\"%v\" for v,e in 1..1 inbound i graph \"userMsgGraph\" limit %v,%v update e with {seen:true} in massageUserEdge return {msg:v,seen:e}", userKey, offset, limit)
	} else {
		query = fmt.Sprintf("for i in users filter i._key==\"%v\" for v,e in 1..1 inbound i graph \"userMsgGraph\" filter e.seen==false return {msg:v,seen:e}", userKey)
	}
	fmt.Println(query)
	return c.JSON(database.ExecuteGetQuery(query))
}

// getMassageBySupplierKey get massages by jwt
// @Summary get massages by jwt
// @Description get massages by jwt , set seen to true if to mark massage as seen
// @Tags massage
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Param   seen      query    bool     false        "seen"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /msg/supplier [get]
func getMassageBySupplierKey(c *fiber.Ctx) error {
	supplierEmployeeKey := c.Locals("supplierEmployeeKey").(string)
	limit := c.Query("limit")
	offset := c.Query("offset")
	update := c.Query("seen")
	if limit == "" || offset == "" {
		return c.Status(400).SendString("offset or limit is empty")
	}
	query := ""
	if update == "true" {
		query = fmt.Sprintf("for i in supplierEmployee filter i._key==\"%v\" for v,e in 1..1 inbound i graph \"userMsgGraph\" limit %v,%v update e with {seen:true} in massageSupplierEdge return {msg:v,seen:e}", supplierEmployeeKey, offset, limit)
	} else {
		query = fmt.Sprintf("for i in supplierEmployee filter i._key==\"%v\" for v,e in 1..1 inbound i graph \"userMsgGraph\" filter e.seen==false return {msg:v,seen:e}", supplierEmployeeKey)
	}

	log.Println(query)
	//query := fmt.Sprintf("for i in supplierEmployee filter i._key==\"%v\" for v in 1..1 inbound i graph \"userMsgGraph\" return v", supplierEmployeeKey)
	return c.JSON(database.ExecuteGetQuery(query))

}
