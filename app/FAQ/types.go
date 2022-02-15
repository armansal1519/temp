package faq

import "github.com/arangodb/go-driver"

type questionIn struct {
	Title        string `json:"title" validate:"required,min=3"`
	Info         []info `json:"info"`
	CategoryKey  string `json:"categoryKey" validate:"required"`
	IsPopular    bool   `json:"isPopular"`
	MoreInfoLink string `json:"moreInfoLink" `
}

type question struct {
	Title         string `json:"title" validate:"required,min=3"`
	Info          []info `json:"info"`
	CategoryKey   string `json:"categoryKey" validate:"required"`
	MoreInfoLink  string `json:"moreInfoLink" `
	IsPopular     bool   `json:"isPopular"`
	IsForSupplier bool   `json:"isForSupplier"`
	CreateAt      int64  `json:"createAt"`
}

type info struct {
	Text     string `json:"text" validate:"required,min=3"`
	ImageUrl string `json:"imageUrl" validate:"required,min=3"`
}

type getQuestion struct {
	question
	driver.DocumentMeta
}

type category struct {
	Title         string `json:"title" validate:"required,min=3"`
	ImageUrl      string `json:"imageUrl" `
	IsForSupplier bool   `json:"isForSupplier"`
	CreatedAt     int64  `json:"createdAt"`
}

type getCategory struct {
	category
	driver.DocumentMeta
}

type input struct {
	Title string `json:"title" validate:"required"`
}

type feedBack struct {
	IsUseful    bool   `json:"isUseful"`
	QuestionKey string `json:"questionKey" `
	MultiSelect string `json:"multiSelect" `
	Comment     string `json:"comment"`
}

type getFeedBack struct {
	feedBack
	driver.DocumentMeta
}
