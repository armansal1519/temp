package transportation

type sendingInfo struct {
	AddressKey          string `json:"addressKey" validate:"required"`
	SendIntervalString  string `json:"sendIntervalString" validate:"required"`
	SendingMethod       string `json:"sendingMethod" validate:"required"`
	TransportationPrice int64  `json:"transportationPrice"`
	TransportationType  string `json:"transportationType"`
}
type sendingInfoOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	sendingInfo
}
