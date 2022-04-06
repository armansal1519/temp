package search

type search struct {
	SearchString string `json:"searchString" validate:"required"`
}

type mostSearch struct {
	Id          string `json:"id" validate:"required"`
	Url         string `json:"url"`
	NameOrTitle string `json:"nameOrTitle" validate:"required"`
	SearchCount int    `json:"searchCount"`
	UserKey     string `json:"userKey"`
}

type mostSearchOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	mostSearch
}

type searchResponse struct {
	Categories []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"categories"`
	MostSearch []mostSearchOut `json:"mostSearch"`
	Products   []struct {
		Id    string `json:"id"`
		Title string `json:"title"`
	} `json:"products"`
}
