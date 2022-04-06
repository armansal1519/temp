package graphPayment

type GPayment struct {
	Type            string `json:"type"`
	TotalPrice      int64  `json:"totalPrice"`
	RemainingPrice  int64  `json:"remainingPrice"`
	FromWallet      int64  `json:"fromWallet"`
	Status          string `json:"status"`
	IsRejected      bool   `json:"isRejected"`
	RejectionTime   int64  `json:"rejectionTime"`
	RejectionReason string `json:"rejectionReason"`
	PaymentKey      string `json:"paymentKey"`
	DiscountKey     string `json:"discountKey"`
	DiscountAmount  int64  `json:"discountAmount"`
}

type GPaymentOut struct {
	Key string `json:"_key"`
	Id  string `json:"_id"`
	Rev string `json:"_rev"`
	GPayment
}

type updatePaymentWithDiscount struct {
	DiscountKey    string `json:"discountKey"`
	DiscountAmount int64  `json:"discountAmount"`
	RemainingPrice int64  `json:"remainingPrice"`
}
