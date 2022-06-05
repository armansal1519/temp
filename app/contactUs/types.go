package contactUs

type contactIn struct {
	Title       string   `json:"title"`
	FullName    string   `json:"fullName"`
	PhoneNumber string   `json:"phoneNumber"`
	Text        string   `json:"text"`
	Status      string   `json:"status"`
	ImageArr    []string `json:"imageArr"`
	Website     []string `json:"website"`
	Email       string   `json:"email"`
}

type contact struct {
	Title       string   `json:"title"`
	FullName    string   `json:"fullName"`
	PhoneNumber string   `json:"phoneNumber"`
	Text        string   `json:"text"`
	ImageArr    []string `json:"imageArr"`
	Status      string   `json:"status"`
	Website     []string `json:"website"`
	Email       string   `json:"email"`
	CreatedAt   int64    `json:"createdAt"`
}

type contactOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	contact
}
