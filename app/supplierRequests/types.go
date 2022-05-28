package supplierRequests

type sRequest struct {
	SupplierKey string `json:"supplierKey"`
	Amount      int64  `json:"amount"`
	Type        string `json:"type"`
}
