package suppliers

import "github.com/arangodb/go-driver"

type SupplierIn struct {
	Address          string   `json:"address" validate:"required"`
	State              string   `json:"state"`
	City                 string        `json:"city"`
	Latitude         float64  `json:"latitude" validate:"required"`
	Longitude        float64  `json:"longitude" validate:"required"`
	Name             string   `json:"name" validate:"required"`
	Code             string   `json:"code" validate:"required"`
	Area             float64  `json:"area" validate:"required"`
	AreaWithRoof     float64  `json:"areaWithRoof" validate:"required"`
	PhoneNumber      string   `json:"phoneNumber" validate:"required"`
	CategoriesToSale []string `json:"categoriesToSale"`
	WalletAmount     int64    `json:"walletAmount"`
	Status           string   `json:"status" validate:"required"`
	CreateAt         int64    `json:"createAt" validate:"required"`
}

type updateSupplier struct {
	SupplierKey string `json:"supplierKey"`
	UpdateData  struct {
		Address          string   `json:"address" `
		Latitude         float64  `json:"latitude" `
		Longitude        float64  `json:"longitude" `
		Name             string   `json:"name" `
		Code             string   `json:"code" `
		Area             float64  `json:"area" `
		AreaWithRoof     float64  `json:"areaWithRoof" `
		PhoneNumber      string   `json:"phoneNumber" `
		CategoriesToSale []string `json:"categoriesToSale"`
	} `json:"updateData"`
	CreateAt int64 `json:"CreateAt"`
}

type supplier struct {
	driver.DocumentMeta
	SupplierIn
}

type Fav struct {
	Key         string `json:"_key,omitempty"`
	SupplierKey string `json:"supplierKey"  validate:"required"`
	CategoryUrl string `json:"categoryUrl"  validate:"required"`
	ProductKey  string `json:"productKey"  validate:"required"`
	Status      string `json:"status"`
}

//type FavIn struct {
//	CategoryUrl string `json:"categoryUrl"  validate:"required"`
//	PriceKey  string `json:"productKey"  validate:"required"`
//	Status            string `json:"status"`
//}

type onlineSuppliers struct {
	Key  string `json:"_key""`
	Uuid string `json:"uuid"`
}
