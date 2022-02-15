package paymentAndWallet

type paymentHistory struct {
	PayerKey              string `json:"payerKey"`
	OrderKey              string `json:"orderKey"`
	TxType                string `json:"txType"`
	Amount                int64  `json:"amount"`
	Status                string `json:"status"`
	CardHolder            string `json:"cardHolder"`
	ShaparakRefId         string `json:"ShaparakRefId"`
	TransId               string `json:"transId"`
	ImageUrl              string `json:"imageUrl"`
	CheckNumber           string `json:"checkNumber"`
	IncludeTransportation bool   `json:"includeTransportation"`
	CreatedAt             int64  `json:"createdAt"`
}

type paymentOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	paymentHistory
}

type updatePaymentHistory struct {
	CardHolder    string `json:"cardHolder"`
	ShaparakRefId string `json:"ShaparakRefId"`
	Status        string `json:"status"`
}

type updateTransId struct {
	TransId string `json:"transId"`
}

type UpdateOrderWithPaymentKey struct {
	PaymentKey string `json:"paymentKey"`
}

type createPaymentByPortal struct {
	OrderKey              string `json:"orderKey"`
	Amount                int64  `json:"amount"`
	Status                string `json:"status"`
	IncludeTransportation bool   `json:"includeTransportation"`
}

type PaymentByImage struct {
	OrderKey              string `json:"orderKey"`
	Type                  string `json:"type"`
	IncludeTransportation bool   `json:"includeTransportation"`
	ImageUrl              string `json:"imageUrl"`
	OverwritePaymentKey   bool   `json:"overwritePaymentKey"`
}

type checkByImage struct {
	OrderKey    string `json:"orderKey"`
	Amount      int64  `json:"amount"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	CheckNumber string `json:"checkNumber"`
	ImageUrl    string `json:"imageUrl"`
}

type filter struct {
	TxType        string `json:"txType"`
	OrderKey      string `json:"orderKey"`
	PayerKey      string `json:"payerKey"`
	Status        string `json:"status"`
	ShaparakRefId string `json:"shaparakRefId"`
	CheckNumber   string `json:"checkNumber"`
}
