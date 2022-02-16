package cart

type cartIn struct {
	PriceId     string `json:"priceId" validate:"required"`
	Number      int    `json:"number" validate:"required"`
	PricingType string `json:"pricingType" validate:"required"`
}

type cart struct {
	PriceId                string  `json:"priceId"`
	SupplierKey            string  `json:"supplierKey"`
	ProductId              string  `json:"productId"`
	PricePerNumber         int64   `json:"pricePerNumber"`
	Number                 int     `json:"number"`
	ProductTitle           string  `json:"productTitle"`
	ProductImageUrl        string  `json:"productImageUrl"`
	PricingType            string  `json:"PricingType"`
	CreatedAt              int64   `json:"createdAt"`
	UserKey                string  `json:"userKey"`
	UserAuthType           string  `json:"userAuthType"`
	CommissionPercent      float64 `json:"commissionPercent"`
	CheckCommissionPercent float64 `json:"checkCommissionPercent"`
	UniqueString           string
}

type CartOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	cart
}

type cErr struct {
	Status    int    `json:"code"`
	ErrorCode int    `json:"error"`
	DevInfo   string `json:"devInfo"`
	UserMsg   string `json:"userMsg"`
}

type GroupedCart struct {
	Type string    `json:"type"`
	Cart []CartOut `json:"cart"`
}

type updateCart struct {
	Number int `json:"number"`
}
