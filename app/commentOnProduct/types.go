package commentOnProduct

type comment struct {
	ProductId    string   `json:"productId"`
	Title        string   `json:"title"`
	Text         string   `json:"text"`
	ImageUrls    []string `json:"imageUrls"`
	UserFullName string   `json:"userFullName"`
	UserKey      string   `json:"userKey"`
	IsAnonymous  bool     `json:"isAnonymous"`
	IsBuyer      bool     `json:"isBuyer"`
	ScoreArr     []score  `json:"scoreArr"`
	Likes         []string `json:"likes"`
	CreatedAt    int64    `json:"createdAt"`
	Status       string   `json:"status"`
}

type commentOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	comment
}

type updateCommentType struct {
	Title       string   `json:"title"`
	Text        string   `json:"text"`
	ImageUrls   []string `json:"imageUrls"`
	IsAnonymous bool     `json:"isAnonymous"`
	Likes         []string `json:"likes"`
	ScoreArr    []score  `json:"scoreArr"`
}

type adminUpdateCommentType struct {
	Title       string   `json:"title"`
	Text        string   `json:"text"`
	ImageUrls   []string `json:"imageUrls"`
	IsAnonymous bool     `json:"isAnonymous"`
	ScoreArr    []score  `json:"scoreArr"`
	Likes         []string `json:"likes"`
	Status      string   `json:"status"`
}

type score struct {
	Title string  `json:"title"`
	Score float32 `json:"score"`
}
