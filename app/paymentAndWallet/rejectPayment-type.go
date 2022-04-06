package paymentAndWallet

type reservedProduct struct {
	PriceId      string `json:"priceId"`
	PaymentKey   string `json:"paymentKey"`
	OrderItemKey string `json:"orderItemKey"`
	Number       int    `json:"number"`
	FailedCount  int    `json:"failedCount"`
	CartKey      string `json:"cartKey"`
	EndTime      int64  `json:"endTime"`
}

type reservedProductOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	reservedProduct
}

type updateReservedProduct struct {
	FailedCount int   `json:"failedCount"`
	EndTime     int64 `json:"endTime"`
}

type updatePaymentForRejection struct {
	IsRejected      bool   `json:"isRejected"`
	RejectionTime   int64  `json:"rejectionTime"`
	RejectionReason string `json:"rejectionReason"`
}

type rejectRequest struct {
	RejectionReason string `json:"rejectionReason"`
}
