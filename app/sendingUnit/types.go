package sendingUnit

type sendUnit struct {
	ApprovedOrderKey      string `json:"approvedOrderKey"`
	Number                int    `json:"number"`
	Status                string `json:"status"`
	CreatedAt             int64  `json:"createdAt"`
	TransportationUnitKey string `json:"transportationUnitKey"`
}

type sendUnitOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	sendUnit
}

type updateApprovedOrder struct {
	ConvertToSendUnit int `json:"convertToSendUnit"`
}

type updateTr struct {
	TransportationUnitKey string `json:"transportationUnitKey"`
}

type updateStatus struct {
	Status string `json:"status"`
}
