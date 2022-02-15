package massage

import (
	"bamachoub-backend-go-v1/app/admin"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func sendMsgByPhoneNumber(c *fiber.Ctx	)error{
	m := new(sendMsgByPhoneNumberReq)
	if err := utils.ParseBodyAndValidate(c, m); err != nil {
		return c.JSON(err)
	}
	adminKey:=c.Locals("adminKey").(string)
	adminCol := database.GetCollection("admin")
	var a admin.AdminOut
	_, err := adminCol.ReadDocument(context.Background(), adminKey, &a)
	if err != nil {
		return c.JSON(err)
	}

	msg:=massage{
		Title:            m.Title,
		ImageUrl:         m.ImageUrl,
		Text:             m.Text,
		Link:             m.Link,
		Importence:       m.Importence,
		CreatedAt:        time.Now().Unix(),
		CreatedBy:        fmt.Sprintf("%v %v",a.FirstName,a.LastName),
		AdminDescription: m.AdminDescription,
	}
	msgCol:=database.GetCollection("massage")

	meta,err:=msgCol.CreateDocument(context.Background(),msg)

	pnStr := "["
	for i, key := range m.PhoneNumberArray {
		pnStr += fmt.Sprintf("\"%v\"", key)
		if i < len(m.PhoneNumberArray)-1 {
			pnStr += " , "
		}
	}
	pnStr += "] "
	query := fmt.Sprintf("let userIds=(for u in users filter u.phoneNumber in %v return u._id)\nfor i in userIds\nINSERT { _from: \"%v\", _to: i , seen:false } INTO massageUserEdge OPTIONS { ignoreErrors: true }\n ",pnStr,meta.ID.String())
	database.GetCollection(query)
	return c.JSON(fiber.Map{"msg":"پیام ها با موفقیت ثبت شد"})
}

func sendSupplierMessage(c *fiber.Ctx) error {
	m := new(SendSupplierMassageReq)
	if err := utils.ParseBodyAndValidate(c, m); err != nil {
		return c.JSON(err)
	}

	if m.For!="user" && m.For=="supplier" {
		return c.Status(400).SendString("For only can be user of supplier")
	}

	adminKey:=c.Locals("adminKey").(string)
	adminCol := database.GetCollection("admin")
	var a admin.AdminOut
	_, err := adminCol.ReadDocument(context.Background(), adminKey, &a)
	if err != nil {
		return c.JSON(err)
	}
	sm:=SupplierMassage{
		Title:            m.Title,
		ImageUrl:         m.ImageUrl,
		Text:             m.Text,
		State:            m.State,
		For:              m.For,
		Link:             m.Link,
		Importence:       m.Importence,
		CreatedAt:        time.Now().Unix(),
		CreatedBy:        fmt.Sprintf("%v %v", a.FirstName, a.LastName),
		AdminDescription: m.AdminDescription,
	}
	msgCol:=database.GetCollection("massage")

	meta,err:=msgCol.CreateDocument(context.Background(),sm)
	if m.For=="user" {
		query := fmt.Sprintf("let s=(for i in suppliers filter i.state==\"%v\" return i._key )\nlet seIds=(for j in supplierEmployee filter j.supplierKey in s return j._id)\nlet seu=UNIQUE(seIds)\nfor k in seu\nINSERT { _from: \"%v\", _to: k , seen:false } INTO massageUserEdge  ",sm.State,meta.ID.String())
		database.ExecuteGetQuery(query)
		return c.JSON(fiber.Map{"msg":"پیام ها با موفقیت ثبت شد"})
	}
	query := fmt.Sprintf("let u=(for i in userAddress filter i.state==\"%v\" return i.userKey )\nlet seu=UNIQUE(u)\nfor k in seu\nINSERT { _from: \"%v\", _to: concat(\"users/\",k) , seen:false } INTO massageUserEdge  ",m.State,meta.ID.String())
	database.ExecuteGetQuery(query)
	return c.JSON(fiber.Map{"msg":"پیام ها با موفقیت ثبت شد"})

}

func getMassageByUserKey(c *fiber.Ctx	)error  {
	userKey := c.Locals("userKey").(string)
	query := fmt.Sprintf("for i in users filter i._key==\"%v\" for v in 1..1 inbound i graph \"userMsgGraph\" return v" , userKey)
	return c.JSON(database.ExecuteGetQuery(query))
}

func getMassageBySupplierKey(c *fiber.Ctx	)error   {
	supplierEmployeeKey:=c.Locals("supplierEmployeeKey").(string)
	query := fmt.Sprintf("for i in supplierEmployee filter i._key==\"%v\" for v in 1..1 inbound i graph \"userMsgGraph\" return v",supplierEmployeeKey)
	return c.JSON(database.ExecuteGetQuery(query))

}


