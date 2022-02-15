package userAddress

type AddressIn struct {
	State        string  `json:"state"`
	City         string  `json:"city"`
	AddressText  string  `json:"addressText"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Pelak        string  `json:"pelak"`
	PostalCode   string  `json:"postalCode"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	PhoneNumber  string  `json:"phoneNumber"`
	NationalCode string  `json:"nationalCode"`
	UserKey      string  `json:"userKey"`
	IsForMySelf  bool    `json:"isForMySelf"`
}

type AddressOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	AddressIn
}
