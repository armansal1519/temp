package userRequests

type userWalletRequest struct {
	UserKey string `json:"userKey"`

	Amount           int64  `json:"amount"`
	Type             string `json:"type"`
	WalletHistoryKey string `json:"walletHistoryKey"`
}

type updateWalletStatus struct {
	TxStatus string `json:"txStatus"`
}
