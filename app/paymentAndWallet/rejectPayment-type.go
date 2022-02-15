package paymentAndWallet

type reservedProduct struct {
	TxType      string `json:"txType"`
	PriceId     string `json:"priceId"`
	PaymentKey string `json:"paymentKey"`
	Number      int    `json:"number"`
	FailedCount int    `json:"failedCount"`
	CartKey     string `json:"cartKey"`
	EndTime     int64  `json:"endTime"`
}

type reservedProductOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	reservedProduct
}

type updateReservedProduct struct {
	FailedCount int    `json:"failedCount"`
	EndTime     int64  `json:"endTime"`

}




