package productSuggestion

type productSuggestion struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
	SupplierKey string `json:"supplierKey"`
}

type sampleSuggestion struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ImageUrls   []string `json:"imageUrls"`
	UserKey     string   `json:"userKey"`
}

type productSuggestionOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	productSuggestion
}

type sampleSuggestionOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	sampleSuggestion
}

type betterPrice struct {
	Price       string `json:"price"`
	ProductId   string `json:"productId"`
	ShopName    string `json:"shopName"`
	ShopType    string `json:"shopType"`
	ShopAddress string `json:"shopAddress"`
	ShopPhone   string `json:"shopPhone"`
	State       string `json:"state"`
	UserKey     string `json:"userKey"`
	CreatedAt   string `json:"createdAt"`
}

type betterPriceOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	betterPrice
}
