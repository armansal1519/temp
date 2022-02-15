package similarityGraph

type similarityEdge struct {
	Key      string `json:"_key"`
	Type     string `json:"type" validate:"required"`
	CoreEdge bool   `json:"coreEdge"`
	From     string `json:"_from" validate:"required"`
	To       string `json:"_to" validate:"required"`
}

type SimilarityNodeRequest struct {
	ProductsKeyArray []string `json:"productsKeyArray"`
	Title            string   `json:"title"`
	Tags             []string `json:"tags"`
	Description      string   `json:"description"`
	Color            string   `json:"color"`
	Pattern          string   `json:"pattern"`
	UserMade         bool     `json:"userMade"`
	Public           bool     `json:"public"`
	IsCollection     bool     `json:"isCollection"`
}

type similarityNode struct {
	Title            string   `json:"title" validate:"required"`
	ProductsKeyArray []string `json:"productsKeyArray"`
	ImageUrls        []string `json:"imageUrl"`
	UserKey          string   `json:"userKey"`
	Tags             []string `json:"tags"`
	Description      string   `json:"description"`
	Color            string   `json:"color"`
	Pattern          string   `json:"pattern"`
	UserMade         bool     `json:"userMade"`
	Public           bool     `json:"public"`
	IsCollection     bool     `json:"isCollection"`
	Status           string   `json:"status"`
	CreatedAt        int64    `json:"createdAt"`
	UpdatedAt        int64    `json:"updatedAt"`
	CreatedBy        string   `json:"createdBy"`
	UpdatedBy        string   `json:"updatedBy"`
	Seen             int      `json:"seen"`
}

type similarityNodeOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	similarityNode
}

type productEdge struct {
	Key   string `json:"_key"`
	From  string `json:"_from"`
	To    string `json:"_to"`
	Score int    `json:"score"`
}

type addProductToNode struct {
	ProductsKeyArray []string `json:"productsKeyArray"`
	ImageUrls        []string `json:"imageUrl"`
}
