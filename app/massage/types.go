package massage

type sendMsgByPhoneNumberReq struct {
	PhoneNumberArray []string `json:"phoneNumberArray"`
	Title       string `json:"title"`
	ImageUrl    string `json:"imageUrl"`
	Text string `json:"text"`
	Link string `json:"link"`
	Importence string `json:"importance"`
	CreatedAt int64 `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	AdminDescription string `json:"adminDescription"`
}

type SendSupplierMassageReq struct {
	Title       string `json:"title"`
	ImageUrl    string `json:"imageUrl"`
	Text string `json:"text"`
	State string `json:"state"`
	For  string `json:"for"`
	Link string `json:"link"`
	Importence string `json:"importance"`
	AdminDescription string `json:"adminDescription"`
}



type massageEdge struct {
	Form string `json:"_form"`
	To string `json:"_to"`
	IsSeen bool `json:"isSeen"`
}

type massage struct {
	Title       string `json:"title"`
	ImageUrl    string `json:"imageUrl"`
	Text string `json:"text"`
	Link string `json:"link"`
	Importence string `json:"importance"`
	CreatedAt int64 `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	AdminDescription string `json:"adminDescription"`

}

type SupplierMassage struct {
	Title       string `json:"title"`
	ImageUrl    string `json:"imageUrl"`
	Text string `json:"text"`
	State string `json:"state"`
	For  string `json:"for"`
	Link string `json:"link"`
	Importence string `json:"importance"`
	CreatedAt int64 `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	AdminDescription string `json:"adminDescription"`
}