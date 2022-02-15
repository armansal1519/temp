package estelam

type onlineSuppliers struct {
	Uuid        string `json:"uuid"`
	SupplierKey string `json:"supplierKey"`
}

type CreateEstelamRequest struct {
	UserKey         string `json:"userKey"`
	Variant         string `json:"variant" validate:"required"`
	ProductId       string `json:"productId" validate:"required"`
	Price           bool   `json:"price"`
	OneMonthPrice   bool   `json:"oneMoundPrice"`
	TwoMonthPrice   bool   `json:"twoMoundPrice"`
	ThreeMonthPrice bool   `json:"threeMoundPrice"`
	Number          int    `json:"number"`
	CreatedAt       int64  `json:"createdAt"`
}

type addToEstelamCart struct {
	Key              string `json:"_key"`
	UserKey          string `json:"userKey"`
	Variant          string `json:"variant" validate:"required"`
	ProductId        string `json:"productId" validate:"required"`
	ImageUrl         string `json:"imageUrl"`
	ProductTitle     string `json:"productTitle"`
	Price            bool   `json:"price"`
	OneMonthPrice    bool   `json:"oneMoundPrice"`
	TwoMonthPrice    bool   `json:"twoMoundPrice"`
	ThreeMonthPrice  bool   `json:"threeMoundPrice"`
	Number           int    `json:"number"`
	CreatedAt        int64  `json:"createdAt"`
	WillExpireAt     int64  `json:"willExpireAt"`
	TimeOfResponse   int64  `json:"timeOfResponse"`
	NumberOfResponse int    `json:"numberOfResponse"`
}

type estelamCartOut struct {
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	addToEstelamCart
}

type estelamSupplier struct {
	SupplierKey     string `json:"supplierKey"`
	EstelamCartKey  string `json:"estelamCartKey"`
	Variant         string `json:"variant" `
	ProductId       string `json:"productId" `
	ImageUrl        string `json:"imageUrl"`
	ProductTitle    string `json:"productTitle"`
	Price           bool   `json:"price"`
	OneMonthPrice   bool   `json:"oneMoundPrice"`
	TwoMonthPrice   bool   `json:"twoMoundPrice"`
	ThreeMonthPrice bool   `json:"threeMoundPrice"`
	Number          int    `json:"number"`
	CreatedAt       int64  `json:"createdAt"`
	WillExpireAt    int64  `json:"willExpireAt"`
}

type estelamSupplierOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	estelamSupplier
}

type responseToEstelamIn struct {
	FormNumber      int    `json:"formNumber"`
	ToNumber        int    `json:"ToNumber" validate:"required"`
	EstelamCartKey  string `json:"estelamCartKey" validate:"required"`
	Price           int64  `json:"price"`
	OneMonthPrice   int64  `json:"oneMoundPrice"`
	TwoMonthPrice   int64  `json:"twoMoundPrice"`
	ThreeMonthPrice int64  `json:"threeMoundPrice"`
	CreatedAt       int64  `json:"createdAt"`
}
