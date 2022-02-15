package userComparisonList

type userComparisonList struct {
	UserKey        string   `json:"userKey"`
	ImageList      []string `json:"imageList"`
	ProductKeyList []string `json:"productKeyList"`
	Title          string   `json:"title"`
	Status         string   `json:"status"`
	CreatedAt      int64    `json:"createdAt"`
	CreatedBy      string   `json:"createdBy"`
	Seen           int      `json:"seen"`
}
