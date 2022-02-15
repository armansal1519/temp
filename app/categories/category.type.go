package categories

type CreateCategoryType struct {
	Name                   string   `validate:"required,min=3,max=32" json:"Name"`
	Url                    string   `json:"url" validate:"required"`
	ImageUrl               string   `json:"imageUrl" `
	Text                   string   `json:"text" `
	From                   string   `validate:"required,min=3,max=32" json:"from"`
	CustomerReviewItems    []string `json:"customerReviewItems"`
	CommissionPercent      float64  `json:"commissionPercent"`
	CheckCommissionPercent float64  `json:"checkCommissionPercent"`
}

type BaseCategoryDto struct {
	Name     string `validate:"required,min=3,max=32" json:"name"`
	Url      string `validate:"required,min=3,max=32" json:"url"`
	Text     string `json:"text" validate:"required"`
	ImageUrl string `json:"imageUrl" `
}
type BaseCategorySave struct {
	Name      string `validate:"required,min=3,max=32" json:"name"`
	Url       string `validate:"required,min=3,max=32" json:"url"`
	Text      string `json:"text" `
	ImageUrl  string `json:"imageUrl" `
	GraphPath string `json:"graphPath"`
	Status    string `json:"status"`
}

type BaseCategoryOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	BaseCategorySave
}

type Category struct {
	Name                   string   `json:"name"`
	Url                    string   `json:"url" validate:"required"`
	GraphPath              string   `json:"graphPath"`
	Status                 string   `json:"status"`
	ImageUrl               string   `json:"imageUrl" `
	Text                   string   `json:"text"`
	CustomerReviewItems    []string `json:"customerReviewItems"`
	CommissionPercent      float64  `json:"commissionPercent"`
	CheckCommissionPercent float64  `json:"checkCommissionPercent"`
}

type CategoryOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	Category
}

type createNewCategoryEdge struct {
	From string `json:"_from"`
	To   string `json:"_to"`
}

type updateStatus struct {
	Status string `json:"status"`
}

type ResponseHTTP struct {
	Status string `json:"status"`
}
