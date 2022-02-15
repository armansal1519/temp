package main

import (
	faq "bamachoub-backend-go-v1/app/FAQ"
	"bamachoub-backend-go-v1/app/addBuyMethod"
	"bamachoub-backend-go-v1/app/admin"
	"bamachoub-backend-go-v1/app/brands"
	"bamachoub-backend-go-v1/app/cart"
	"bamachoub-backend-go-v1/app/categories"
	"bamachoub-backend-go-v1/app/commentOnProduct"
	"bamachoub-backend-go-v1/app/contactUs"
	"bamachoub-backend-go-v1/app/driversRegister"
	"bamachoub-backend-go-v1/app/estelam"
	homepage "bamachoub-backend-go-v1/app/homePage"
	"bamachoub-backend-go-v1/app/imageUpload"
	"bamachoub-backend-go-v1/app/massage"
	"bamachoub-backend-go-v1/app/orders"
	"bamachoub-backend-go-v1/app/paymentAndWallet"
	"bamachoub-backend-go-v1/app/productQA"
	"bamachoub-backend-go-v1/app/products"
	"bamachoub-backend-go-v1/app/products/menu"
	"bamachoub-backend-go-v1/app/products/productStructure"
	"bamachoub-backend-go-v1/app/products/productSuggestion"
	"bamachoub-backend-go-v1/app/search"
	"bamachoub-backend-go-v1/app/sendingInfo"
	"bamachoub-backend-go-v1/app/sendingUnit"
	"bamachoub-backend-go-v1/app/similarityGraph"
	"bamachoub-backend-go-v1/app/supplierEmployees"
	"bamachoub-backend-go-v1/app/suppliers"
	"bamachoub-backend-go-v1/app/userAddress"
	"bamachoub-backend-go-v1/app/users"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(v1 fiber.Router) {
	products.Routes(v1)
	categories.Routes(v1)
	users.Routes(v1)
	users.AuthRoutes(v1)
	suppliers.Routes(v1)
	suppliers.WsRoutes(v1)
	supplierEmployees.Routes(v1)
	supplierEmployees.AuthRoutes(v1)
	productStructure.Routes(v1)
	menu.Routes(v1)
	imageUpload.Routes(v1)
	brands.Routes(v1)
	addBuyMethod.Routes(v1)
	productSuggestion.Routes(v1)
	cart.Routes(v1)
	estelam.Routes(v1)
	faq.Routes(v1)
	faq.CategoryRoutes(v1)
	contactUs.Routes(v1)
	driversRegister.Routes(v1)
	driversRegister.DriversRoutes(v1)
	similarityGraph.Routes(v1)
	productQA.Routes(v1)
	userAddress.Routes(v1)
	paymentAndWallet.WalletRoutes(v1)
	paymentAndWallet.PaymentRoutes(v1)
	paymentAndWallet.SupplierConfirmationRoute(v1)
	paymentAndWallet.RejectionRoutes(v1)
	orders.Routes(v1)
	sendingInfo.Routes(v1)
	sendingUnit.Routes(v1)
	commentOnProduct.Routes(v1)

	admin.AuthRoutes(v1)
	admin.Routes(v1)
	massage.Routes(v1)
	search.Routes(v1)
	homepage.Routes(v1)

}
