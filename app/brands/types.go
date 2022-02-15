package brands

type brandDto struct {
	Name        string `json:"name" validate:"required"`
	EName        string `json:"eName" validate:"required"`
	ImageUrl    string `json:"imageUrl" validate:"required"`
	Description string `json:"description" validate:"required"`
	CategoryKey string `json:"categoryKey" validate:"required"`
	Seen     string `json:"seen"`
}

type Brand struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	brandDto
}

type editBrand struct {
	Name        string `json:"name" validate:"required"`
	ImageUrl    string `json:"imageUrl" validate:"required"`
	Description string `json:"description" validate:"required"`
	EName        string `json:"eName" validate:"required"`
}
