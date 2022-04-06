package paymentAndWallet

import "bamachoub-backend-go-v1/app/cart"

type supplierInfoForConfirmation struct {
	SupplierKey  string `json:"supplierKey"`
	OrderKey     string `json:"orderKey"`
	OrderItemKey string `json:"cartKey"`
}

type supplierInfoForConfirmationOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	supplierInfoForConfirmation
}

type getSupplierConfirmationResponse struct {
	Cart cart.CartOut                   `json:"cart"`
	Info supplierInfoForConfirmationOut `json:"info"`
}

//type ApprovedOrder struct {
//	UserKey           string  `json:"userKey"`
//	ProductId         string  `json:"productId"`
//	ProductTitle      string  `json:"productTitle"`
//	ProductImageUrl   string  `json:"productImageUrl"`
//	SupplierKey       string  `json:"supplierKey"`
//	PaymentKey        string  `json:"paymentKey"`
//	TxType            string  `json:"txType"`
//	Price             int64   `json:"price"`
//	CommissionPercent float64 `json:"commissionPercent"`
//	ConvertToSendUnit int     `json:"convertToSendUnit"`
//	Number            int     `json:"number"`
//	CreatedAt         int64   `json:"createdAt"`
//	SendInfoKey       string  `json:"sendInfoKey"`
//	Status            string  `json:"status"`
//}

type updateOrder struct {
	IsApprovedBySupplier bool   `json:"isApprovedBySupplier"`
	SupplierEmployeeId   string `json:"supplierEmployeeId"`
}

//type ApprovedOrderOut struct {
//	Key string `json:"_key,omitempty"`
//	ID  string `json:"_id,omitempty"`
//	Rev string `json:"_rev,omitempty"`
//	ApprovedOrder
//}

type rejectionPoolItem struct {
	UserKey         string `json:"userKey"`
	ProductId       string `json:"productId"`
	ProductTitle    string `json:"productTitle"`
	ProductImageUrl string `json:"productImageUrl"`
	RejectBy        string `json:"rejectBy"`
	PaymentKey      string `json:"paymentKey"`
	TxType          string `json:"txType"`
	Price           int64  `json:"price"`
	Number          int    `json:"number"`
	CreatedAt       int64  `json:"createdAt"`
	SendInfoKey     string `json:"sendInfoKey"`
	Status          string `json:"status"`
}

type rejectionPoolItemOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	rejectionPoolItem
}
