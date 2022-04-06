package gOrderItem

type orderItemCancelRequest struct {
	ProductId    string `json:"productId"`
	Number       int    `json:"number"`
	CancelReason string `json:"cancelReason"`
	CancelAll    bool   `json:"cancelAll"`
}

type orderItemCancel struct {
	ProductId    string `json:"productId"`
	Number       int    `json:"number"`
	CancelReason string `json:"cancelReason"`
	UserKey      string `json:"userKey"`
	CreatedAt    int64  `json:"createdAt"`
	CancelAll    bool   `json:"cancelAll"`
}

type OrderItemReferRequest struct {
	ProductId          string   `json:"productId"`
	Number             int      `json:"number"`
	ReferReason        string   `json:"referReason"`
	ReferReasonDetails string   `json:"referReasonDetails"`
	ImageArr           []string `json:"imageArr"`
	CancelAll          bool     `json:"cancelAll"`
}

type OrderItemRefer struct {
	ProductId          string   `json:"productId"`
	Number             int      `json:"number"`
	ReferReason        string   `json:"referReason"`
	ReferReasonDetails string   `json:"referReasonDetails"`
	UserKey            string   `json:"userKey"`
	CreatedAt          int64    `json:"createdAt"`
	ImageArr           []string `json:"imageArr"`
	CancelAll          bool     `json:"cancelAll"`
}
