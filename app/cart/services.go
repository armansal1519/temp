package cart

import (
	"bamachoub-backend-go-v1/app/addBuyMethod"
	"bamachoub-backend-go-v1/app/products"
	"bamachoub-backend-go-v1/app/users"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"log"
	"strings"
	"time"
)





func addToCart(cartIn cartIn, isLogin bool, userKey string, isAuthenticated bool) (*CartOut, cErr) {
	edgeName := strings.Split(cartIn.PriceId, "/")[0]
	edgeKey := strings.Split(cartIn.PriceId, "/")[1]

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
			UserMsg:   "تعداد انتخاب شده از تعداد مچاز بیشتر است",
		}
	}

	if cartIn.PricingType == "price" {
		if !isLogin {
			if cartIn.UserKey == "" {
				userKey, err = users.CreateHeadlessUser()
				log.Print(userKey)
				if err != nil {
					return nil, cErr{
						Status:    409,
						ErrorCode: 1,
						DevInfo:   err.Error(),
						UserMsg:   "مشکل در ایجاد کاربر اولیه",
					}
				}
			} else {
				userKey = cartIn.UserKey
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

	c := cart{
		PriceId:           cartIn.PriceId,
		Number:            cartIn.Number,
		PricingType:       cartIn.PricingType,
		PricePerNumber:    pricePerNumber,
		ProductTitle:      p.Title,
		ProductImageUrl:   p.ImageArr[0],
		CreatedAt:         time.Now().Unix(),
		CommissionPercent: cp,
		UserKey:           userKey,
		SupplierKey:       supplierKey,
		ProductId:         price.To,
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
