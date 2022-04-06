package discountCode

type DiscountForPhoneNumbersRequest struct {
	PhoneNumbers []string `json:"phoneNumbers"`
	Type         string   `json:"type"`
	Amount       int64    `json:"amount"`
	EndAt        int64    `json:"endAt"`
}

type Discount struct {
	Type   string `json:"type"`
	Amount int64  `json:"amount"`
	EndAt  int64  `json:"endAt"`
}

type DiscountOut struct {
	Key string `json:"_key"`
	Id  string `json:"_id"`
	Rev string `json:"_rev"`
	Discount
}
