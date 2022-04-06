package homepage

type homepageBase struct {
	Carousel       []Carousel       `json:"carousel"`
	Banners        []Banners        `json:"banners"`
	ProductSlider  []ProductSlider  `json:"productSlider"`
	CategorySlider []CategorySlider `json:"categorySlider"`
	BrandSlider    []BrandSlider    `json:"brandSlider"`
	Text           []Text           `json:"text"`
	Email          Email            `json:"email"`
	BlogContent    BlogContent      `json:"blogContent"`
}
type Carousel struct {
	ImageURL string `json:"imageUrl"`
	Link     string `json:"link"`
}
type Data struct {
	ImageURL string `json:"imageUrl"`
	Link     string `json:"link"`
}
type Banners struct {
	Title           string `json:"title"`
	Position        int    `json:"position"`
	NumberOfBanners int    `json:"numberOfBanners"`
	Data            []Data `json:"data"`
}
type ProductSlider struct {
	Title        string `json:"title"`
	Position     int    `json:"position"`
	CategoryName string `json:"categoryName"`
	Sort         string `json:"sort"`
}
type CategorySlider struct {
	Title        string `json:"title"`
	Position     int    `json:"position"`
	CategoryName string `json:"categoryName"`
	Sort         string `json:"sort"`
}
type BrandSlider struct {
	Title        string `json:"title"`
	Position     int    `json:"position"`
	CategoryName string `json:"categoryName"`
	Sort         string `json:"sort"`
}
type Text struct {
	Title    string `json:"title"`
	Position int    `json:"position"`
	Text     string `json:"text"`
	BtnText  string `json:"btnText"`
	Link     string `json:"link"`
}
type Email struct {
	Show bool `json:"show"`
}
type BlogContent struct {
	Show bool `json:"show"`
}

type saveQuery struct {
	Query string `json:"query"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type finalHomePageData struct {
	Key  string      `json:"_key"`
	Data []saveQuery `json:"data"`
}
