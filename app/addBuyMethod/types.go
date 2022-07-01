package addBuyMethod

type PriceIn struct {
	ProductId         string `json:"productId" validate:"required"`
	Price             int64  `json:"price" `
	OneMonthPrice     int64  `json:"oneMonthPrice" `
	TwoMonthPrice     int64  `json:"twoMonthPrice" `
	ThreeMonthPrice   int64  `json:"threeMonthPrice" `
	Variant           string `json:"variant" validate:"required"`
	TotalNumber       int64  `json:"totalNumber" `
	TotalNumberInCart int64  `json:"totalNumberInCart" `
	CodeForSupplier   string `json:"codeForSupplier" `
	Show              bool   `json:"show" `
}

type updatePrice struct {
	Price             int64  `json:"price" `
	OneMonthPrice     int64  `json:"oneMonthPrice" `
	TwoMonthPrice     int64  `json:"twoMonthPrice" `
	ThreeMonthPrice   int64  `json:"threeMonthPrice" `
	TotalNumber       int64  `json:"totalNumber" `
	TotalNumberInCart int64  `json:"totalNumberInCart" `
	CodeForSupplier   string `json:"codeForSupplier" `
	Show              bool   `json:"show" `
}

type PriceGroupCreate struct {
	From              string `json:"_from"`
	To                string `json:"_to"`
	Price             int64  `json:"price"`
	OneMonthPrice     int64  `json:"oneMonthPrice"`
	TwoMonthPrice     int64  `json:"twoMonthPrice"`
	ThreeMonthPrice   int64  `json:"threeMonthPrice"`
	Variant           string `json:"variant"`
	TotalNumber       int64  `json:"totalNumber"`
	TotalNumberInCart int64  `json:"totalNumberInCart"`
	CodeForSupplier   string `json:"codeForSupplier"`
	PriceRepetition   int    `json:"priceRepetition"`

	Show      bool  `json:"show"`
	CreatedAt int64 `json:"createdAt"`
}

type PriceOut struct {
	PriceGroupCreate
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
}

type estelamIn struct {
	ProductId       string `json:"productId" `
	CodeForSupplier string `json:"codeForSupplier"`
	Variant         string `json:"variant"`
	Price           bool   `json:"price"`
	OneMonthPrice   bool   `json:"oneMonthPrice"`
	TwoMonthPrice   bool   `json:"twoMonthPrice"`
	ThreeMonthPrice bool   `json:"threeMonthPrice"`
	Show            bool   `json:"show" `
}

//type estelamCreate struct {
//	From            string        `json:"_from"`
//	To              string        `json:"_to"`
//	Key             string        `json:"_key"`
//	CodeForSupplier string        `json:"codeForSupplier"`
//	VariantArr      []baseEstelam `json:"variantArr"`
//
//	CreatedAt int64 `json:"createdAt"`
//}

type baseEstelam struct {
	Variant         string `json:"variant"`
	Price           bool   `json:"price"`
	OneMonthPrice   bool   `json:"oneMonthPrice"`
	TwoMonthPrice   bool   `json:"twoMonthPrice"`
	ThreeMonthPrice bool   `json:"threeMonthPrice"`
	Show            bool   `json:"show" `
}

type CreateEstelam struct {
	From            string `json:"_from"`
	To              string `json:"_to"`
	Key             string `json:"_key"`
	CodeForSupplier string `json:"codeForSupplier"`
	Variant         string `json:"variant"`
	Price           bool   `json:"price"`
	OneMonthPrice   bool   `json:"oneMonthPrice"`
	TwoMonthPrice   bool   `json:"twoMonthPrice"`
	ThreeMonthPrice bool   `json:"threeMonthPrice"`
	Show            bool   `json:"show" `
	CreatedAt       int64  `json:"createdAt"`
}

type updateEstelam struct {
	CodeForSupplier string `json:"codeForSupplier"`
	Price           bool   `json:"price"`
	OneMonthPrice   bool   `json:"oneMonthPrice"`
	TwoMonthPrice   bool   `json:"twoMonthPrice"`
	ThreeMonthPrice bool   `json:"threeMonthPrice"`
	Show            bool   `json:"show" `
}

type groupUpdateEstelamIn struct {
	PriceKeys       []string `json:"priceKeys"`
	ChangeBuyMode   bool     `json:"changeBuyMode"`
	ChangeStatus    bool     `json:"changeStatus"`
	Price           bool     `json:"price"`
	OneMonthPrice   bool     `json:"oneMonthPrice"`
	TwoMonthPrice   bool     `json:"twoMonthPrice"`
	ThreeMonthPrice bool     `json:"threeMonthPrice"`
	Show            bool     `json:"show" `
}

type groupUpdateEstelam1 struct {
	Price           bool `json:"price"`
	OneMonthPrice   bool `json:"oneMonthPrice"`
	TwoMonthPrice   bool `json:"twoMonthPrice"`
	ThreeMonthPrice bool `json:"threeMonthPrice"`
}
type groupUpdateEstelam2 struct {
	Show bool `json:"show" `
}

type EstelamOut struct {
	CreateEstelam
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
}

type updateLowestPriceInProduct struct {
	LowestPrice int64 `json:"lowestPrice"`
}

type Product struct {
	CategoryKey       string     `json:"categoryKey" `
	CategoryPath      string     `json:"categoryName" `
	Title             string     `json:"title"`
	CompleteSpec      []csType   `json:"completeSpec"`
	MainSpecs         []f        `json:"mainSpecs"`
	VariationObj      v          `json:"variationsObj"`
	ImageArr          []string   `json:"imageArr"`
	Status            string     `json:"status" `
	Description       string     `json:"description"`
	Brand             string     `json:"brand"`
	Tags              []string   `json:"tags"`
	CreatedAt         int64      `json:"createdAt"`
	SeenNumber        int        `json:"seenNumber"`
	BuyNumber         int        `json:"buyNumber"`
	CommissionPercent float32    `json:"commissionPercent"`
	LowestPrice       int64      `json:"lowestPrice"`
	LowestCheckPrice  checkPrice `json:"lowestCheckPrice"`

	FilterString string `json:"filterString"`
}

type checkPrice struct {
	Type  string `json:"type"`
	Price int64  `json:"price"`
}
type updateCheckPrice struct {
	LowestCheckPrice checkPrice `json:"lowestCheckPrice"`
}

type updatePriceAndCheckPrice struct {
	LowestPrice      int64      `json:"lowestPrice"`
	LowestCheckPrice checkPrice `json:"lowestCheckPrice"`
}

type priceAndProduct struct {
	Product Product  `json:"product"`
	Price   PriceOut `json:"price"`
}

type estelamAndProduct struct {
	Product Product    `json:"product"`
	Estelam EstelamOut `json:"estelam"`
}

type csType struct {
	Name  string `json:"name"`
	Items []f    `json:"items"`
}

type f struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type v struct {
	Title      string `json:"title"`
	Variations []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"variations"`
}

type groupUpdatePriceIn struct {
	PriceKeys          []string `json:"priceKeys"`
	ChangePrice        bool     `json:"changePrice"`
	Price              bool     `json:"price"`
	OneMonthPrice      bool     `json:"oneMonthPrice"`
	TwoMonthPrice      bool     `json:"twoMonthPrice"`
	ThreeMonthPrice    bool     `json:"threeMonthPrice"`
	ChangePriceMethod  string   `json:"changePriceMethod"`
	ChangePriceValue   int64    `json:"changePriceValue"`
	ChangeNumber       bool     `json:"changeNumber"`
	ChangeNumberMethod string   `json:"changeNumberMethod"`
	ChangeNumberValue  int64    `json:"changeNumberValue"`
	ChangeStatus       bool     `json:"changeStatus"`
	Show               bool     `json:"show" `
}

type groupUpdatePrice1 struct {
	Show bool `json:"show" `
}

type brandFilter struct {
	Brand  string `json:"brand" `
	Search string `json:"search" `
}
