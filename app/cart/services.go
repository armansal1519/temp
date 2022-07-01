package cart

import (
	"bamachoub-backend-go-v1/app/addBuyMethod"
	"bamachoub-backend-go-v1/app/products"
	"bamachoub-backend-go-v1/app/users"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
	"time"
)

// addToCart  add to Cart
// @Summary adds to Cart
// @Description adds to Cart , if there is jwt or temp-user-key adds to Cart of user else create new temp-user
// @Tags cart
// @Accept json
// @Produce json
// @Param data body cartIn true "data"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @param temp-user-key header string false "temp-user-key"
// @Success 200 {object} CartOut{}
// @Failure 404 {object} cErr{}
// @Router /cart [post]
func addToCart(cartIn cartIn, isLogin bool, userKey string, tempUserKey string, isAuthenticated bool) (*CartOut, cErr) {
	fmt.Println(cartIn.PriceId)
	edgeName := strings.Split(cartIn.PriceId, "/")[0]
	edgeKey := strings.Split(cartIn.PriceId, "/")[1]
	log.Println(edgeKey)

	var price addBuyMethod.PriceOut
	edgeCol := database.GetCollection(edgeName)
	_, err := edgeCol.ReadDocument(context.Background(), edgeKey, &price)
	if err != nil {
		return nil, cErr{
			Status:    404,
			ErrorCode: 1,
			DevInfo:   err.Error(),
			UserMsg:   "قیمت پیدا نشد",
		}
	}
	if cartIn.Number > int(price.TotalNumber) || cartIn.Number > int(price.TotalNumberInCart) {
		return nil, cErr{
			Status:    400,
			ErrorCode: 1,
			DevInfo:   fmt.Sprintf("%v  %v  %v", cartIn.Number, price.TotalNumber, price.TotalNumberInCart),
			UserMsg:   "تعداد انتخاب شده از تعداد مجاز بیشتر است",
		}
	}

	if cartIn.PricingType == "price" {
		if !isLogin {
			if tempUserKey == "" {
				userKey, err = users.CreateHeadlessUser()
				if err != nil {
					return nil, cErr{
						Status:    409,
						ErrorCode: 1,
						DevInfo:   err.Error(),
						UserMsg:   "مشکل در ایجاد کاربر اولیه",
					}
				}
			} else {
				userKey = tempUserKey
			}
		}

	} else if cartIn.PricingType == "one" || cartIn.PricingType == "two" || cartIn.PricingType == "three" {
		if !isLogin {
			return nil, cErr{
				Status:    403,
				ErrorCode: 1,
				DevInfo:   fmt.Sprintf("user is not logged in"),
				UserMsg:   "برای این فعالیت کاربر باید ورود کرده باشد",
			}
		}
		if !isAuthenticated {
			return nil, cErr{
				Status:    403,
				ErrorCode: 2,
				DevInfo:   fmt.Sprintf("user is not authenticated"),
				UserMsg:   "برای این فعالیت کاربر بایداحراز هویت کرده باشد",
			}
		}

	} else {
		return nil, cErr{
			Status:    409,
			ErrorCode: 2,
			DevInfo:   fmt.Sprintf("pricing type is unaccepted %v", cartIn.PricingType),
			UserMsg:   "",
		}
	}

	//check if buy method is valid
	if price.Price == 0 && cartIn.PricingType == "price" {
		return nil, cErr{
			Status:    409,
			ErrorCode: 3,
			DevInfo:   fmt.Sprintf("pricing type is not available %v", cartIn.PricingType),
			UserMsg:   "تامین کننده این نوع خرید را اجازه نداده",
		}
	}

	if price.OneMonthPrice == 0 && cartIn.PricingType == "one" {
		return nil, cErr{
			Status:    409,
			ErrorCode: 3,
			DevInfo:   fmt.Sprintf("pricing type is not available %v", cartIn.PricingType),
			UserMsg:   "تامین کننده این نوع خرید را اجازه نداده",
		}
	}
	if price.TwoMonthPrice == 0 && cartIn.PricingType == "two" {
		return nil, cErr{
			Status:    409,
			ErrorCode: 3,
			DevInfo:   fmt.Sprintf("pricing type is not available %v", cartIn.PricingType),
			UserMsg:   "تامین کننده این نوع خرید را اجازه نداده",
		}
	}
	if price.ThreeMonthPrice == 0 && cartIn.PricingType == "three" {
		return nil, cErr{
			Status:    409,
			ErrorCode: 3,
			DevInfo:   fmt.Sprintf("pricing type is not available %v", cartIn.PricingType),
			UserMsg:   "تامین کننده این نوع خرید را اجازه نداده",
		}
	}

	var pricePerNumber int64
	if cartIn.PricingType == "price" {
		pricePerNumber = price.Price
	} else if cartIn.PricingType == "one" {
		pricePerNumber = price.OneMonthPrice
	} else if cartIn.PricingType == "two" {
		pricePerNumber = price.TwoMonthPrice
	} else if cartIn.PricingType == "three" {
		pricePerNumber = price.ThreeMonthPrice
	} else {
		return nil, cErr{
			Status:    409,
			ErrorCode: 3,
			DevInfo:   fmt.Sprintf("unknwon first price"),
			UserMsg:   "",
		}
	}
	idSlice := strings.Split(price.To, "/")
	productColName := idSlice[0]
	productKey := idSlice[1]
	productCol := database.GetCollection(productColName)
	var p products.Product
	_, err = productCol.ReadDocument(context.Background(), productKey, &p)
	if err != nil {
		return nil, cErr{
			Status:    409,
			ErrorCode: 3,
			DevInfo:   fmt.Sprintf("error in reading product \n error: %v ", err),
			UserMsg:   "",
		}
	}

	cartCol := database.GetCollection("cart")

	supplierKey := strings.Split(price.From, "/")[1]
	var cp float64
	if cartIn.PricingType == "price" {
		cp = p.CommissionPercent
	} else {
		cp = p.CheckCommissionPercent
	}
	uat := ""
	if isAuthenticated && isLogin {
		uat = "authenticated"
	} else if !isAuthenticated && isLogin {
		uat = "login"
	} else {
		uat = "headless"
	}

	c := Cart{
		PriceId:           cartIn.PriceId,
		Number:            cartIn.Number,
		PricingType:       cartIn.PricingType,
		PricePerNumber:    pricePerNumber,
		ProductTitle:      p.Title,
		ProductImageUrl:   p.ImageArr[0],
		CreatedAt:         time.Now().Unix(),
		CommissionPercent: cp,
		UserAuthType:      uat,
		UserKey:           userKey,
		Variant:           price.Variant,
		SupplierKey:       supplierKey,
		ProductId:         price.To,
		UniqueString:      fmt.Sprintf("%v_%v_%v", cartIn.PriceId, cartIn.PricingType, userKey),
	}

	var co CartOut
	ctx := driver.WithReturnNew(context.Background(), &co)
	_, err = cartCol.CreateDocument(ctx, c)
	if err != nil {
		return nil, cErr{
			Status:    409,
			ErrorCode: 3,
			DevInfo:   err.Error(),
			UserMsg:   "مشکل در ایجاد کالا در سبد خرید",
		}
	}
	return &co, cErr{Status: -1}
}

// getCartByUserKey  get Cart by user key
// @Summary get Cart by user key
// @Description get Cart by user key , by jwt or by temp-user-key
// @Tags cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @param temp-user-key header string false "temp-user-key"
// @Success 200 {object} []CartOut{}
// @Failure 404 {object} cErr{}
// @Router /cart [get]
func getCartByUserKey(isLogin bool, userKey string) (*[]CartOut, cErr) {
	q := ""
	if !isLogin {
		q = fmt.Sprintf("for i in cart filter i.userKey==\"%v\"  and i.userAuthType==\"headless\" return i", userKey)

	} else {
		q = fmt.Sprintf("for i in cart filter i.userKey==\"%v\"  and i.userAuthType!=\"headless\" return i", userKey)

	}

	db := database.GetDB()
	ctx := context.Background()
	log.Println(q)
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		return nil, cErr{
			Status:    409,
			ErrorCode: 1,
			DevInfo:   err.Error(),
			UserMsg:   "مشکل در گرفتن کالا در سبد خرید",
		}
	}
	defer cursor.Close()
	var data []CartOut
	for {
		var doc CartOut
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}

	return &data, cErr{Status: -1}
}

// update edit Cart by key
// @Summary edit Cart by key
// @Description edit Cart by key , jwt or  temp-user-key must exist
// @Tags cart
// @Accept json
// @Produce json
// @Param data body updateCart true "data"
// @Param   key      path   string     true  "key"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @param temp-user-key header string false "temp-user-key"
// @Success 200 {object} CartOut{}
// @Failure 404 {object} cErr{}
// @Router /cart/{key} [patch]
func update(c *fiber.Ctx) error {
	cartKey := c.Params("key")
	userKey := c.Locals("userKey").(string)
	tempUserKey := c.Get("temp-user-key", "")

	if tempUserKey == "" && userKey == "" {
		return c.Status(400).SendString("user key missing ")
	}
	if userKey == "" && tempUserKey != "" {
		userKey = tempUserKey
	}
	uc := new(updateCart)
	if err := utils.ParseBodyAndValidate(c, uc); err != nil {
		return c.JSON(err)
	}

	var co CartOut
	cartCol := database.GetCollection("cart")
	_, err := cartCol.ReadDocument(context.Background(), cartKey, &co)
	if err != nil {
		return c.JSON(err)
	}
	edgeName := strings.Split(co.PriceId, "/")[0]
	edgeKey := strings.Split(co.PriceId, "/")[1]

	var price addBuyMethod.PriceOut
	edgeCol := database.GetCollection(edgeName)
	_, err = edgeCol.ReadDocument(context.Background(), edgeKey, &price)
	if err != nil {
		return c.JSON(err)
	}
	log.Println(co.PriceId)
	if uc.Number > int(price.TotalNumber) || uc.Number > int(price.TotalNumberInCart) {
		return c.Status(400).JSON(cErr{
			Status:    400,
			ErrorCode: 1,
			DevInfo:   fmt.Sprintf("%v  %v  %v", uc.Number, price.TotalNumber, price.TotalNumberInCart),
			UserMsg:   "تعداد انتخاب شده از تعداد مجاز بیشتر است",
		})
	}

	q := fmt.Sprintf("for i in cart filter i._key==\"%v\" and i.userKey==\"%v\" update i with {number: %v} in cart return NEW", cartKey, userKey, uc.Number)

	data := database.ExecuteGetQuery(q)
	if data == nil {
		return c.Status(404).SendString("Cart not found")
	}
	return c.JSON(data[0])

}

// remove delete Cart by key
// @Summary delete Cart by key
// @Description delete Cart by key , jwt or  temp-user-key must exist
// @Tags cart
// @Accept json
// @Produce json
// @Param   key      path   string     true  "key"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @param temp-user-key header string false "temp-user-key"
// @Success 200 {object} CartOut{}
// @Failure 404 {object} cErr{}
// @Router /cart/{key} [delete]
func remove(c *fiber.Ctx) error {
	cartKey := c.Params("key")
	userKey := c.Locals("userKey").(string)
	tempUserKey := c.Get("temp-user-key", "")

	if tempUserKey == "" && userKey == "" {
		return c.Status(400).SendString("user key missing ")
	}
	if userKey == "" && tempUserKey != "" {
		userKey = tempUserKey
	}
	q := fmt.Sprintf("for i in cart filter i._key==\"%v\" and i.userKey==\"%v\" remove i in cart", cartKey, userKey)
	database.ExecuteGetQuery(q)
	return c.Status(204).JSON("ok")
}

//
//func removeOne(c *fiber.Ctx) error {
//
//}
