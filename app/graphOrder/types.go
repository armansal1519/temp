package graphOrder

import (
	"bamachoub-backend-go-v1/app/graphPayment"
	"bamachoub-backend-go-v1/app/users"
)

type GOrderItem struct {
	PriceId                string  `json:"priceId"`
	SupplierKey            string  `json:"supplierKey"`
	ProductId              string  `json:"productId"`
	PricePerNumber         int64   `json:"pricePerNumber"`
	Number                 int     `json:"number"`
	Variant                string  `json:"variant"`
	ProductTitle           string  `json:"productTitle"`
	ProductImageUrl        string  `json:"productImageUrl"`
	UserKey                string  `json:"userKey"`
	UserAuthType           string  `json:"userAuthType"`
	CommissionPercent      float64 `json:"commissionPercent"`
	CheckCommissionPercent float64 `json:"checkCommissionPercent"`
	IsWaitingForPayment    bool    `json:"isWaitingForPayment"`

	IsApprovedBySupplier bool   `json:"isApprovedBySupplier"`
	SupplierEmployeeId   string `json:"supplierEmployeeId"`

	IsProcessing bool `json:"isProcessing"`

	IsArrived bool `json:"isArrived"`

	IsCancelled   bool   `json:"isCancelled"`
	CancelledById string `json:"cancelledById"`

	IsReferred     bool   `json:"isReferred"`
	ReferredReason string `json:"referredReason"`

	CreatedAt int64 `json:"createdAt"`
}

type GOrderItemOut struct {
	Key string `json:"_key"`
	Id  string `json:"_id"`
	Rev string `json:"_rev"`
	GOrderItem
}

type GOrder struct {
	SendingInfoKey               string `json:"sendingInfoKey"`
	UserKey                      string `json:"userKey"`
	TransportationPrice          int64  `json:"transportationPrice"`
	IsTransportationPriceIsPayed bool   `json:"isTransportationPriceIsPayed"`
	TransportationPriceWithPrice bool   `json:"transportationPriceWithPrice"`
	UseWalletForTransportation   bool   `json:"useWalletForTransportation"`
	TransportationPaymentId      string `json:"transportationPaymentId"`
	TotalAmount                  int64  `json:"totalAmount"`
	Status                       string `json:"status"`
	CreateAt                     int64  `json:"createAt"`
}

type GOrderOut struct {
	Key string `json:"_key"`
	Id  string `json:"_id"`
	Rev string `json:"_rev"`
	GOrder
}

type edgeData struct {
	Payment    graphPayment.GPayment
	OrderItems []GOrderItem
}

type edgeIdData struct {
	PaymentIds    string
	OrderItemIdes []string
}

type orderEdge struct {
	From string `json:"_from"`
	To   string `json:"_to"`
}

type GOrderResponseOut struct {
	Order GOrderOut `json:"order"`
	Items []struct {
		Payment    graphPayment.GPaymentOut `json:"payment"`
		OrderItems []GOrderItemOut          `json:"orderItems"`
	} `json:"items"`
	Reserved reservedInfo `json:"reserved"`
}

type sendingInfo struct {
	OrderKey       string `json:"orderKey"`
	SendingInfoKey string `json:"sendingInfoKey"`
}

type updateOrderBySendingInfo struct {
	TransportationPrice int64  `json:"transportationPrice"`
	SendingInfoKey      string `json:"sendingInfoKey"`
}

type OrderPaymentAndOrderItem struct {
	Order     GOrderOut
	Payment   graphPayment.GPaymentOut
	OrderItem GOrderItem
}

type OrderItemsAndPayment struct {
	Payment    graphPayment.GPaymentOut
	OrderItems []GOrderItemOut
}

type reservedInfo struct {
	IsReserved bool  `json:"isReserved"`
	TimeToEnd  int64 `json:"timeToEnd"`
}

type getOrderForAdminDto struct {
	OrderStatus  []string `json:"orderStatus"`
	PaymentTypes []string `json:"paymentTypes"`
	States       []string `json:"states"`
	Time         int64    `json:"time"`
}

type GOrderResponseForAdminOut struct {
	User  users.UserOut `json:"user"`
	Order GOrderOut     `json:"order"`
	Items []struct {
		Payment    graphPayment.GPaymentOut `json:"payment"`
		OrderItems []GOrderItemOut          `json:"orderItems"`
	} `json:"items"`
	Reserved reservedInfo `json:"reserved"`
}
