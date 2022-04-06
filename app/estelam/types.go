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
	OneMonthPrice   bool   `json:"oneMonthPrice"`
	TwoMonthPrice   bool   `json:"twoMonthPrice"`
	ThreeMonthPrice bool   `json:"threeMonthPrice"`
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
	FromNumber      int    `json:"fromNumber"`
	ToNumber        int    `json:"ToNumber"`
	EstelamCartKey  string `json:"estelamCartKey" validate:"required"`
	Price           int64  `json:"price"`
	OneMonthPrice   int64  `json:"oneMonthPrice"`
	TwoMonthPrice   int64  `json:"twoMonthPrice"`
	ThreeMonthPrice int64  `json:"threeMonthPrice"`
}

type createResponseToEstelam struct {
	FromNumber          int    `json:"fromNumber"`
	ToNumber            int    `json:"ToNumber"`
	EstelamCartKey      string `json:"estelamCartKey" validate:"required"`
	SupplierKey         string `json:"supplierKey"`
	SupplierEmployeeKey string `json:"supplierEmployeeKey"`
	Price               int64  `json:"price"`
	PricingType         string `json:"pricingType"`
	CreatedAt           int64  `json:"createdAt"`
	ExpireAt            int64  `json:"expireAt"`
}
type responseToEstelamOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	createResponseToEstelam
}

//type getEstelamForUserResp struct {
//	EstelamItem struct {
//		Key              string `json:"_key"`
//		Id               string `json:"_id"`
//		Rev              string `json:"_rev"`
//		UserKey          string `json:"userKey"`
//		Variant          string `json:"variant"`
//		ProductId        string `json:"productId"`
//		ImageUrl         string `json:"imageUrl"`
//		ProductTitle     string `json:"productTitle"`
//		Price            bool   `json:"price"`
//		OneMonthPrice    bool   `json:"oneMonthPrice"`
//		TwoMonthPrice    bool   `json:"twoMonthPrice"`
//		ThreeMonthPrice  bool   `json:"threeMonthPrice"`
//		Number           int    `json:"number"`
//		CreatedAt        int    `json:"createdAt"`
//		WillExpireAt     int    `json:"willExpireAt"`
//		TimeOfResponse   int    `json:"timeOfResponse"`
//		NumberOfResponse int    `json:"numberOfResponse"`
//	} `json:"estelamItem"`
//	SupplierResponse []responseToEstelamOut `json:"SupplierResponse"`
//}

type getEstelamForUserResp struct {
	EstelamItem struct {
		Key              string `json:"_key"`
		Id               string `json:"_id"`
		Rev              string `json:"_rev"`
		UserKey          string `json:"userKey"`
		Variant          string `json:"variant"`
		ProductId        string `json:"productId"`
		ImageUrl         string `json:"imageUrl"`
		ProductTitle     string `json:"productTitle"`
		Price            bool   `json:"price"`
		OneMoundPrice    bool   `json:"oneMoundPrice"`
		TwoMoundPrice    bool   `json:"twoMoundPrice"`
		ThreeMoundPrice  bool   `json:"threeMoundPrice"`
		Number           int    `json:"number"`
		CreatedAt        int    `json:"createdAt"`
		WillExpireAt     int    `json:"willExpireAt"`
		TimeOfResponse   int    `json:"timeOfResponse"`
		NumberOfResponse int    `json:"numberOfResponse"`
	} `json:"estelamItem"`
	SupplierResponse []struct {
		Key                 string `json:"_key"`
		Id                  string `json:"_id"`
		Rev                 string `json:"_rev"`
		FromNumber          int    `json:"fromNumber"`
		ToNumber            int    `json:"ToNumber"`
		EstelamCartKey      string `json:"estelamCartKey"`
		SupplierKey         string `json:"supplierKey"`
		SupplierEmployeeKey string `json:"supplierEmployeeKey"`
		Price               int    `json:"price"`
		PricingType         string `json:"pricingType"`
		CreatedAt           int    `json:"createdAt"`
		ExpireAt            int    `json:"expireAt"`
	} `json:"SupplierResponse"`
}

type cartFromEstelam struct {
	SupplierResponseKey string `json:"supplierResponseKey"`
	EstelamCartKey      string `json:"estelamCartKey"`
	Number              int    `json:"number"`
}
