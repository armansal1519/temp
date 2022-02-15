package orders

import "bamachoub-backend-go-v1/app/cart"

type OrderItem struct {
	Type              string         `json:"type"`
	Cart              []cart.CartOut `json:"cart"`
	TotalPrice        int64          `json:"totalPrice"`
	RemainingPrice    int64          `json:"remainingPrice"`
	PaymentKey        string         `json:"paymentKey"`
	FromWallet        int64          `json:"fromWallet"`
	StatusForEachItem []string       `json:"statusForEachItem"`
	Status            string         `json:"status"`
}

type Order struct {
	Key                          string      `json:"_key,omitempty"`
	ID                           string      `json:"_id,omitempty"`
	Rev                          string      `json:"_rev,omitempty"`
	UserKey                      string      `json:"userKey"`
	SendingInfoKey               string      `json:"sendingInfoKey"`
	OrderItems                   []OrderItem `json:"orderItems"`
	Status                       string      `json:"status"`
	TransportationPrice          int64       `json:"transportationPrice"`
	IsTransportationPriceIsPayed bool        `json:"isTransportationPriceIsPayed"`
	TransportationPriceWithPrice bool        `json:"transportationPriceWithPrice"`
	UseWalletForTransportation   bool        `json:"useWalletForTransportation"`
	TransportationPaymentId      string      `json:"transportationPaymentId"`
}
