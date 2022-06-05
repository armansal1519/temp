package productQA

type productQA struct {
	ProductId       string   `json:"productId"`
	ProductTitle    string   `json:"productTitle"`
	ProductImageArr []string `json:"productImageArr"`
	Text            string   `json:"text"`
	CreatedAt       int64    `json:"createdAt"`
	FullName        string   `json:"fullName"`
	UserKey         string   `json:"userKey"`
	QuestionKey     string   `json:"questionKey"`
	Likes           []string `json:"likes"`
	Status          string   `json:"status"`
	RejectionText   string   `json:"rejectionText"`
}

type updateDto struct {
	Text   string `json:"text"`
	Status string `json:"status"`
}
