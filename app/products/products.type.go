package products

import (
	"bamachoub-backend-go-v1/app/addBuyMethod"
	"github.com/arangodb/go-driver"
)

type productBodyDto struct {
	CategoryKey string `json:"categoryKey" validate:"required"`

	TitleMaker struct {
		Title string   `json:"title"`
		Data  []string `json:"data"`
	} `json:"titleMaker"`
	CompleteSpec []struct {
		Name  string   `json:"name"`
		Items []string `json:"items"`
	} `json:"completeSpec"`
	MainSpecs    []string `json:"mainSpecs"`
	VariationObj v        `json:"variationsObj"`
	//ImageArr     []string `json:"imageArr"`
	//Status       string   `json:"status" `
	//Description  string   `json:"description"`
	//Brand        string   `json:"brand"`
	//Tags         []string `json:"tags"`
}

type v struct {
	Title      string `json:"title"`
	Variations []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"variations"`
}

type productBodyIn struct {
	CategoryKey string `json:"categoryKey" validate:"required"`
	//CategoryPath string `json:"categoryName" validate:"required"`
	TitleMaker   titleMaker `json:"titleMaker"`
	CompleteSpec []csType   `json:"completeSpec"`
	MainSpecs    []f        `json:"mainSpecs"`
	VariationObj v          `json:"variationsObj"`
	ImageArr     []string   `json:"imageArr"`
	Status       string     `json:"status" `
	Description  string     `json:"description"`
	Brand        string     `json:"brand"`
	Tags         []string   `json:"tags"`
}

type ProductBodyOut struct {
	driver.DocumentMeta
	productBodyIn
}

type csType struct {
	Name  string `json:"name"`
	Items []f    `json:"items"`
}

type f struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//type baseProductDto struct {
//	CategoryKey  string            `json:"categoryKey" validate:"required"`
//	Fields       []f 			   `json:"fields" `
//	VariationObj v                 `json:"variationsObj"`
//	ImageArr     []string          `json:"imageArr"`
//	Status       string            `json:"status" `
//	Description  string            `json:"description"`
//	Brand        string            `json:"brand"`
//	Tags         []string          `json:"tags"`
//}

type productInfo struct {
	CategoryKey   string `json:"categoryKey"`
	Fields        []f    `json:"fields"`
	VariationsObj []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"variationsObj"`
	ImageArr    []string `json:"imageArr"`
	Description string   `json:"description"`
	Brand       string   `json:"brand"`
	SpId        string   `json:"spId"`
	Tags        []string `json:"tags"`
}

type Product struct {
	CategoryKey            string     `json:"categoryKey" `
	CategoryPath           string     `json:"categoryName" `
	Title                  string     `json:"title"`
	CompleteSpec           []csType   `json:"completeSpec"`
	MainSpecs              []f        `json:"mainSpecs"`
	VariationObj           v          `json:"variationsObj"`
	ImageArr               []string   `json:"imageArr"`
	SpId                   string     `json:"spId"`
	Status                 string     `json:"status" `
	Description            string     `json:"description"`
	Brand                  string     `json:"brand"`
	Tags                   []string   `json:"tags"`
	CreatedAt              int64      `json:"createdAt"`
	SeenNumber             int        `json:"seenNumber"`
	BuyNumber              int        `json:"buyNumber"`
	CommissionPercent      float64    `json:"commissionPercent"`
	CheckCommissionPercent float64    `json:"checkCommissionPercent"`
	LowestPrice            int64      `json:"lowestPrice"`
	LowestCheckPrice       checkPrice `json:"lowestCheckPrice"`
	FilterArr              []string   `json:"filterArr"`
}

type checkPrice struct {
	Type  string `json:"type"`
	Price int64  `json:"price"`
}

type productOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	Product
}

type getProductByKeyResponse struct {
	Product    productOut                `json:"Product"`
	EstelamArr []addBuyMethod.EstelamOut `json:"estelamArr"`
	PriceArr   []addBuyMethod.PriceOut   `json:"priceArr"`
}

type titleMaker struct {
	Title string   `json:"title"`
	Data  []string `json:"data"`
}

type colorOut struct {
	Main productOut   `json:"main"`
	Sub  []productOut `json:"sub"`
}
