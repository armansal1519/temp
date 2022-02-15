package sendingInfo

type sendDayInterval struct {
	Key       string     `json:"_key,omitempty"`
	Intervals []interval `json:"intervals"`
}

type addIntervalRequest struct {
	Key      string   `json:"_key,omitempty"`
	Interval interval `json:"interval"`
}

type interval struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type sendingInfo struct {
	UserKey             string `json:"userKey"`
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

type updateSendingInfoKey struct {
	SendingInfoKey string `json:"sendingInfoKey"`
}

type updateSendingInfo struct {
	AddressKey          string `json:"addressKey"`
	SendIntervalString  string `json:"sendIntervalString" `
	SendingMethod       string `json:"sendingMethod" `
	TransportationPrice int64  `json:"transportationPrice"`
	TransportationType  string `json:"transportationType"`
}
