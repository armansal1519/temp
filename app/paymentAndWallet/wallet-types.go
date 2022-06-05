package paymentAndWallet

type SupplierWalletHistory struct {
	Amount      int64  `json:"amount"`
	SupplierKey string `json:"supplierKey"`
	PaymentKey  string `json:"paymentKey"`
	CreatedAt   int64  `json:"createdAt"`
	Income      bool   `json:"income"`
	TxType      string `json:"txType"`
	TxStatus    string `json:"txStatus"`
}

type supplierWalletOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	SupplierWalletHistory
}

type UserWalletHistory struct {
	Amount     int64  `json:"amount"`
	UserKey    string `json:"userKey"`
	PaymentKey string `json:"paymentKey"`
	CreatedAt  int64  `json:"createdAt"`
	Income     bool   `json:"income"`
	TxType     string `json:"txType"`
	TxStatus   string `json:"txStatus"`
}

type userWalletOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	SupplierWalletHistory
}

type supplierPageResponse struct {
	WalletAmount int64               `json:"walletAmount"`
	TotalIn      int64               `json:"totalIn"`
	TotalOut     int64               `json:"totalOut"`
	History      []supplierWalletOut `json:"history"`
}

type temp1 struct {
	Income  int64 `json:"income"`
	Outcome int64 `json:"outcome"`
}

type addToWallet struct {
	Amount int64 `json:"amount"`
}
