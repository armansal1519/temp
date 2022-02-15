package supplierEmployees

import (
	"bamachoub-backend-go-v1/utils"
)

type createSupplierPreviewIn struct {
	FirstName          string   `json:"firstName" validate:"required"`
	LastName           string   `json:"lastName" validate:"required"`
	PhoneNumber        string   `json:"phoneNumber" validate:"required"`
	Email              string   `json:"email"  `
	NationalCode       string   `json:"nationalCode"`
	BirthDate          string   `json:"birthDate" validate:"required"`
	ShabaNumber        string   `json:"shabaNumber" validate:"required"`
	ShopName           string   `json:"shopName" validate:"required"`
	Latitude           float64  `json:"latitude" validate:"required"`
	Longitude          float64  `json:"longitude" validate:"required"`
	State              string   `json:"state" validate:"required"`
	City               string   `json:"city" validate:"required"`
	Address            string   `json:"address" validate:"required"`
	PostalCode         string   `json:"postalCode"`
	CategoriesToSale   []string `json:"categoriesToSale"`
	IdCardImage        string   `json:"idCardImage"`
	IdBookPageOneImage string   `json:"idBookPageOneImage"`
	IdBookPageTwoImage string   `json:"idBookPageTwoImage"`
	SalesPermitImage   string   `json:"salesPermitImage"`
	//Access             []string `json:"access"`
	//Role               string   `json:"role"`
	//HashPassword       string   `json:"hashPassword"`
	//HashRefreshToken   string   `json:"hashRefreshToken"`
	//SupplierKey        string   `json:"supplierKey"`
	CreateAt int64 `json:"createAt"`
	//LastLogin          int64    `json:"lastLogin"`
}

type createSupplierPreview struct {
	FirstName          string   `json:"firstName" validate:"required"`
	LastName           string   `json:"lastName" validate:"required"`
	PhoneNumber        string   `json:"phoneNumber" validate:"required"`
	Email              string   `json:"email"  `
	NationalCode       string   `json:"nationalCode" validate:"required"`
	ShenasNameCode     string   `json:"shenasNameCode" validate:"required"`
	BirthDate          string   `json:"birthDate" validate:"required"`
	ShabaNumber        string   `json:"shabaNumber" validate:"required"`
	ShopName           string   `json:"shopName" validate:"required"`
	Latitude           float64  `json:"latitude" validate:"required"`
	Longitude          float64  `json:"longitude" validate:"required"`
	State              string   `json:"state" validate:"required"`
	City               string   `json:"city" validate:"required"`
	Address            string   `json:"address" validate:"required"`
	PostalCode         string   `json:"postalCode"`
	CategoriesToSale   []string `json:"categoriesToSale"`
	IdCardImage        string   `json:"idCardImage"`
	IdBookPageOneImage string   `json:"idBookPageOneImage"`
	IdBookPageTwoImage string   `json:"idBookPageTwoImage"`
	SalesPermitImage   string   `json:"salesPermitImage"`
	//Access             []string `json:"access"`
}

type supplierPreview struct {
	createSupplierPreviewIn
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
}

type employeeIn struct {
	FirstName          string   `json:"firstName" validate:"required"`
	LastName           string   `json:"lastName" validate:"required"`
	PhoneNumber        string   `json:"phoneNumber" validate:"required"`
	Email              string   `json:"email"  `
	NationalCode       string   `json:"nationalCode" validate:"required"`
	ShenasNameCode     string   `json:"shenasNameCode" validate:"required"`
	BirthDate          string   `json:"birthDate" validate:"required"`
	ShabaNumber        string   `json:"shabaNumber" validate:"required"`
	IdCardImage        string   `json:"idCardImage"`
	IdBookPageOneImage string   `json:"idBookPageOneImage"`
	IdBookPageTwoImage string   `json:"idBookPageTwoImage"`
	SalesPermitImage   string   `json:"salesPermitImage"`
	Access             []string `json:"access"`
	Role               string   `json:"role"`
	HashPassword       string   `json:"hashPassword"`
	HashRefreshToken   string   `json:"hashRefreshToken"`
	SupplierKey        string   `json:"supplierKey"`
	CreateAt           int64    `json:"createAt"`
	LastLogin          int64    `json:"lastLogin"`
	Status             string   `json:"status"`
}

type employee struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	employeeIn
}

type updateEmployee struct {
	SupplierEmployeeKey string `json:"supplierEmployeeKey"`
	UpdateData          struct {
		FirstName          string `json:"firstName" `
		LastName           string `json:"lastName" `
		PhoneNumber        string `json:"phoneNumber" `
		Email              string `json:"email"  `
		NationalCode       string `json:"nationalCode" `
		ShenasNameCode     string `json:"shenasNameCode" `
		BirthDate          string `json:"birthDate" `
		ShabaNumber        string `json:"shabaNumber" `
		IdCardImage        string `json:"idCardImage"`
		IdBookPageOneImage string `json:"idBookPageOneImage"`
		IdBookPageTwoImage string `json:"idBookPageTwoImage"`
		SalesPermitImage   string `json:"salesPermitImage"`
	} `json:"updateData"`
	CreateAt int64 `json:"CreateAt"`
}

type getValidationCodeDto struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,numeric"`
}

type checkValidationCodeDto struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,numeric"`
	Code        string `json:"code" validate:"required,numeric,len=4"`
}

type SaveValidationCode struct {
	Key       string `json:"_key"`
	Code      string `json:"code"`
	CreatedAt int64  `json:"createdAt"`
}

type ResponseHTTP struct {
	Status string `json:"status"`
}

type CustomErrorResponse utils.CError

type loginResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	IsFirstLogin bool     `json:"isFirstLogin"`
	Employee     employee `json:"employee"`
}

type UpdateRefreshToken struct {
	HashRefreshToken string `json:"hashRefreshToken"`
	LastLogin        int64  `json:"lastLogin"`
}

type loginRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type changePassword struct {
	HashPassword string `json:"hashPassword"`
}

type changePasswordWithoutLoginRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Code        string `json:"code"`
}

type changePasswordWithLoginRequest struct {
	Password string `json:"password"`
}

type refreshTokenResponse struct {
	Token string `json:"token"`
}
