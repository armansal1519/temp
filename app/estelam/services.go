package estelam

import (
	"bamachoub-backend-go-v1/app/cart"
	"bamachoub-backend-go-v1/app/products"
	"bamachoub-backend-go-v1/app/suppliers"
	"bamachoub-backend-go-v1/app/users"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"bamachoub-backend-go-v1/utils/sms"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
	"time"
)

// createEstelamRequest  create estelam by user
// @Summary create estelam by user
// @Description create estelam by user
// @Tags estelam
// @Accept json
// @Produce json
// @Param data body CreateEstelamRequest true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /estelam/create [post]
func createEstelamRequest(c *fiber.Ctx) error {
	cer := new(CreateEstelamRequest)

	if err := utils.ParseBodyAndValidate(c, cer); err != nil {
		return c.JSON(err)
	}

	productColName := strings.Split(cer.ProductId, "/")[0]
	productKey := strings.Split(cer.ProductId, "/")[1]
	userKey := c.Locals("userKey").(string)

	if !cer.Price && !cer.ThreeMonthPrice && !cer.TwoMonthPrice && !cer.OneMonthPrice {
		return c.Status(400).SendString("all options is false")
	}

	query := fmt.Sprintf("let a=(for j in supplier_%v_estelam \nfilter j._to==\"%v\" \nfilter j.variant==\"%v\"\nfilter  ", productColName, cer.ProductId, cer.Variant)

	if cer.Price {
		query += "j.price==true ||"
	}
	if cer.OneMonthPrice {
		query += "j.oneMoundPrice==true ||"
	}
	if cer.TwoMonthPrice {
		query += "j.twoMoundPrice==true ||"
	}
	if cer.ThreeMonthPrice {
		query += "j.threeMoundPrice ==true ||"
	}

	chars := []rune(query)
	newChars := chars[:len(chars)-2]

	query = string(newChars)

	query += "  return j._from) return UNIQUE(a)[0]"

	log.Println(query)
	var p products.Product
	productCol := database.GetCollection(productColName)
	_, err := productCol.ReadDocument(context.Background(), productKey, &p)

	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var supplierIds []string
	for {
		var doc string
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		supplierIds = append(supplierIds, doc)
	}

	if len(supplierIds) <= 0 {
		return c.Status(409).SendString("no supplier was found")
	}

	aec := addToEstelamCart{
		Key:              userKey + productKey,
		UserKey:          userKey,
		Variant:          cer.Variant,
		ProductId:        cer.ProductId,
		ProductTitle:     p.Title,
		ImageUrl:         p.ImageArr[0],
		Price:            cer.Price,
		OneMonthPrice:    cer.OneMonthPrice,
		TwoMonthPrice:    cer.TwoMonthPrice,
		ThreeMonthPrice:  cer.ThreeMonthPrice,
		Number:           cer.Number,
		CreatedAt:        time.Now().Unix(),
		WillExpireAt:     time.Now().Add(2 * time.Hour).Unix(),
		TimeOfResponse:   -1,
		NumberOfResponse: 0,
	}
	estelamCartCol := database.GetCollection("estelamCart")

	meta, err := estelamCartCol.CreateDocument(context.Background(), aec)
	if err != nil {
		return c.JSON(err)
	}
	supplierEstelamArr := make([]estelamSupplier, 0)
	supplierEstelamCol := database.GetCollection("supplierEstelam")
	log.Print(supplierIds)
	supplierKeys := make([]string, 0)
	for _, si := range supplierIds {
		temp := estelamSupplier{
			SupplierKey:     strings.Split(si, "/")[1],
			EstelamCartKey:  meta.Key,
			Variant:         cer.Variant,
			ProductId:       cer.ProductId,
			ImageUrl:        p.ImageArr[0],
			ProductTitle:    p.Title,
			Price:           cer.Price,
			OneMonthPrice:   cer.OneMonthPrice,
			TwoMonthPrice:   cer.TwoMonthPrice,
			ThreeMonthPrice: cer.ThreeMonthPrice,
			Number:          cer.Number,
			CreatedAt:       time.Now().Unix(),
			WillExpireAt:    time.Now().Add(2 * time.Hour).Unix(),
		}
		supplierEstelamArr = append(supplierEstelamArr, temp)
		supplierKeys = append(supplierKeys, strings.Split(si, "/")[1])
	}
	suppliers.NewEstelam(supplierKeys)

	metaArr, errArr, err := supplierEstelamCol.CreateDocuments(context.Background(), supplierEstelamArr)
	if err != nil {
		return c.Status(409).JSON(fiber.Map{
			"error":    err,
			"errorArr": errArr,
		})
	}

	return c.JSON(metaArr)

}

// getEstelamForSupplier  get estelam request for supplier
// @Summary get estelam request for supplier
// @Description get estelam request for supplier
// @Tags estelam
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []estelamSupplierOut{}
// @Failure 404 {object} string{}
// @Router /estelam/supplier [get]
func getEstelamForSupplier(supplierKey string) (*[]estelamSupplierOut, error) {
	query := fmt.Sprintf("for i in supplierEstelam filter i.supplierKey==\"%v\" return i", supplierKey)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error while running query:%v", query)
	}
	defer cursor.Close()
	var data []estelamSupplierOut
	for {
		var doc estelamSupplierOut
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error while running query:%v", query)

		}
		data = append(data, doc)
	}
	return &data, err
}

// getEstelamForUser  get estelam request for user
// @Summary get estelam request for user
// @Description get estelam request for user
// @Tags estelam
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []getEstelamForUserResp{}
// @Failure 404 {object} string{}
// @Router /estelam/user [get]
func getEstelamCart(userKey string) (*[]getEstelamForUserResp, error) {
	query := fmt.Sprintf("for i in estelamCart filter i.userKey==\"%v\"\nlet resp=(for j in estelamResponse filter j.estelamCartKey==i._key return j )\nreturn {estelamItem:i,SupplierResponse:resp}", userKey)
	log.Println(query)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error while running query:%v", query)
	}
	defer cursor.Close()
	var data []getEstelamForUserResp
	for {
		var doc getEstelamForUserResp
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error while running query:%v", query)

		}
		data = append(data, doc)
	}
	return &data, err
}

// responseToEstelam  get estelam request for supplier
// @Summary get estelam request for supplier
// @Description get estelam request for supplier
// @Tags estelam
// @Accept json
// @Produce json
// @Param data body responseToEstelamIn true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /estelam/supplier/response [post]
func responseToEstelam(c *fiber.Ctx) error {
	rte := new(responseToEstelamIn)

	if err := utils.ParseBodyAndValidate(c, rte); err != nil {
		return c.JSON("validation err : " + err.Error())
	}

	var ec estelamCartOut
	query := fmt.Sprintf("for i in estelamCart filter i._key==\"%v\"\nupdate i with {numberOfResponse: i.numberOfResponse +1,timeOfResponse:%v} in estelamCart\nreturn NEW", rte.EstelamCartKey, time.Now().Unix())
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return c.JSON(fmt.Errorf("error while running query:%v", query))
	}
	defer cursor.Close()

	_, err = cursor.ReadDocument(ctx, &ec)
	if err != nil {
		return c.JSON(err)
	}

	supplierId := c.Locals("supplierId").(string)
	supplierKey := supplierId
	flag, err := isSupplierValid(supplierKey, rte.EstelamCartKey)
	fmt.Println("flag=", flag, err)
	if !flag {
		return c.Status(403).SendString(fmt.Sprintf("supplier is not allowed . supplierKey: %v", supplierKey))

	}

	//if (ec.Price && rte.Price == 0) || (!ec.Price && rte.Price != 0) {
	//	return c.Status(409).SendString(fmt.Sprintf("price does not match %v != %v", ec.Price, rte.Price))
	//}
	//if (ec.OneMonthPrice && rte.OneMonthPrice == 0) || (!ec.OneMonthPrice && rte.OneMonthPrice != 0) {
	//	return c.Status(409).SendString(fmt.Sprintf("price does not match %v != %v", ec.OneMonthPrice, rte.OneMonthPrice))
	//}
	//if (ec.TwoMonthPrice && rte.TwoMonthPrice == 0) || (!ec.TwoMonthPrice && rte.TwoMonthPrice != 0) {
	//	return c.Status(409).SendString(fmt.Sprintf("price does not match %v != %v", ec.TwoMonthPrice, rte.TwoMonthPrice))
	//}
	//if (ec.ThreeMonthPrice && rte.ThreeMonthPrice == 0) || (!ec.ThreeMonthPrice && rte.ThreeMonthPrice != 0) {
	//	return c.Status(409).SendString(fmt.Sprintf("price does not match %v != %v", ec.ThreeMonthPrice, rte.ThreeMonthPrice))
	//}

	seKey := c.Locals("supplierEmployeeKey").(string)

	responseArr := make([]createResponseToEstelam, 0)
	if ec.Price {
		cr := createResponseToEstelam{
			FromNumber:          rte.FromNumber,
			ToNumber:            rte.ToNumber,
			EstelamCartKey:      rte.EstelamCartKey,
			SupplierKey:         supplierId,
			SupplierEmployeeKey: seKey,
			Price:               rte.Price,
			PricingType:         "price",
			CreatedAt:           time.Now().Unix(),
			ExpireAt:            time.Now().Add(2 * time.Hour).Unix(),
		}
		responseArr = append(responseArr, cr)
	}
	if ec.OneMonthPrice {
		cr := createResponseToEstelam{
			FromNumber:          rte.FromNumber,
			ToNumber:            rte.ToNumber,
			EstelamCartKey:      rte.EstelamCartKey,
			SupplierKey:         supplierId,
			SupplierEmployeeKey: seKey,
			Price:               rte.Price,
			PricingType:         "one",
			CreatedAt:           time.Now().Unix(),
			ExpireAt:            time.Now().Add(2 * time.Hour).Unix(),
		}
		responseArr = append(responseArr, cr)
	}
	if ec.TwoMonthPrice {
		cr := createResponseToEstelam{
			FromNumber:          rte.FromNumber,
			ToNumber:            rte.ToNumber,
			EstelamCartKey:      rte.EstelamCartKey,
			SupplierKey:         supplierId,
			SupplierEmployeeKey: seKey,
			Price:               rte.Price,
			PricingType:         "two",
			CreatedAt:           time.Now().Unix(),
			ExpireAt:            time.Now().Add(2 * time.Hour).Unix(),
		}
		responseArr = append(responseArr, cr)
	}
	if ec.ThreeMonthPrice {
		cr := createResponseToEstelam{
			FromNumber:          rte.FromNumber,
			ToNumber:            rte.ToNumber,
			EstelamCartKey:      rte.EstelamCartKey,
			SupplierKey:         supplierId,
			SupplierEmployeeKey: seKey,
			Price:               rte.Price,
			PricingType:         "three",
			CreatedAt:           time.Now().Unix(),
			ExpireAt:            time.Now().Add(2 * time.Hour).Unix(),
		}
		responseArr = append(responseArr, cr)
	}

	estelamResponseCol := database.GetCollection("estelamResponse")
	meta, errArr, err := estelamResponseCol.CreateDocuments(context.Background(), responseArr)
	if err != nil {
		return c.JSON(fiber.Map{
			"error":  err,
			"errArr": errArr,
		})
	}
	if ec.NumberOfResponse == 1 {
		u, _ := users.GetUserByKey(ec.UserKey)
		pArr := sms.ParameterArray{
			Parameter:      "Full name",
			ParameterValue: fmt.Sprintf("%v %v", u.FirstName, u.LastName),
		}
		sms.SendSms(u.PhoneNumber, "57117", []sms.ParameterArray{pArr})

	}

	return c.JSON(meta)

}

// createCartFromEstelam  create cart by estelam
// @Summary create cart by estelam
// @Description create cart by estelam
// @Tags estelam
// @Accept json
// @Produce json
// @Param data body cartFromEstelam true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} cart.CartOut{}
// @Failure 404 {object} string{}
// @Router /estelam/to-cart [post]
func createCartFromEstelam(c *fiber.Ctx) error {
	cfs := new(cartFromEstelam)

	if err := utils.ParseBodyAndValidate(c, cfs); err != nil {
		return c.JSON("validation err : " + err.Error())
	}
	var sr responseToEstelamOut
	estelamResponseCol := database.GetCollection("estelamResponse")
	_, err := estelamResponseCol.ReadDocument(context.Background(), cfs.SupplierResponseKey, &sr)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	var ec estelamCartOut
	estelamCartCol := database.GetCollection("estelamCart")
	_, err = estelamCartCol.ReadDocument(context.Background(), cfs.EstelamCartKey, &ec)
	if err != nil {
		return c.Status(500).JSON(err)
	}
	//var p int64
	//if cfs.PricingType == "price" {
	//	p = sr.Price
	//} else if cfs.PricingType == "one" {
	//	p = sr.OneMonthPrice
	//} else if cfs.PricingType == "two" {
	//	p = sr.TwoMonthPrice
	//} else if cfs.PricingType == "three" {
	//	p = sr.ThreeMonthPrice
	//} else {
	//	return c.Status(400).SendString(fmt.Sprintf("pricing type canonly be price,one ,two,three but it was : %v", cfs.PricingType))
	//}

	if cfs.Number > sr.ToNumber || cfs.Number < sr.FromNumber {
		return c.Status(400).SendString(fmt.Sprintf("number must be between %v and %v but it was : %v", sr.FromNumber, sr.ToNumber, cfs.Number))
	}

	newCart := cart.Cart{
		PriceId:                ec.ID,
		SupplierKey:            sr.SupplierKey,
		ProductId:              ec.ProductId,
		PricePerNumber:         sr.Price,
		Number:                 cfs.Number,
		Variant:                ec.Variant,
		ProductTitle:           ec.ProductTitle,
		ProductImageUrl:        ec.ImageUrl,
		PricingType:            sr.PricingType,
		CreatedAt:              time.Now().Unix(),
		UserKey:                ec.UserKey,
		UserAuthType:           "authenticated",
		CommissionPercent:      1.5,
		CheckCommissionPercent: 1.5,
		UniqueString:           fmt.Sprintf("%v_%v_%v", ec.ID, sr.PricingType, ec.UserKey),
	}
	var co cart.CartOut
	ctx := driver.WithReturnNew(context.Background(), &co)
	cartCol := database.GetCollection("cart")
	_, err = cartCol.CreateDocument(ctx, newCart)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(co)
}
