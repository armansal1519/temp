package admin

type adminIn struct {
	FirstName        string   `json:"firstName"`
	LastName         string   `json:"lastName"`
	PhoneNumber      string   `json:"phoneNumber"`
	NationalCode     string   `json:"nationalCode"`
	ShenasNameCode   string   `json:"shenasNameCode"`
	BirthDate        string   `json:"birthDate"`
	ShabaNumber      string   `json:"shabaNumber" `
	Access           []string `json:"access"`
	HashPassword     string   `json:"hashPassword"`
	HashRefreshToken string   `json:"hashRefreshToken"`
	CreateAt         int64    `json:"createAt"`
	LastLogin        int64    `json:"lastLogin"`
	Status           string   `json:"status"`
}

type updateAdminIn struct {
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	PhoneNumber    string   `json:"phoneNumber"`
	NationalCode   string   `json:"nationalCode"`
	ShenasNameCode string   `json:"shenasNameCode"`
	BirthDate      string   `json:"birthDate"`
	ShabaNumber    string   `json:"shabaNumber" `
	Access         []string `json:"access"`
	Status         string   `json:"status"`
}

type AdminOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	adminIn
}

type loginResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	Admin        AdminOut `json:"admin"`
}

type UpdateRefreshToken struct {
	HashRefreshToken string `json:"hashRefreshToken"`
	LastLogin        int64  `json:"lastLogin"`
}

type loginRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

func getAllAdminAccess() []string {
	return []string{
		"read-interval",
		"write-interval",
		"delete-interval",
		"read-carAndTransportation",
		"write-carAndTransportation",
		"delete-carAndTransportation",
		"read-userAndUserInput",
		"write-userAndUserInput",
		"delete-userAndUserInput",
		"read-massageAndFAQ",
		"write-massageAndFAQ",
		"delete-massageAndFAQ",
		"read-productAndBrand",
		"write-productAndBrand",
		"delete-productAndBrand",
		"read-CategoryAndCollection",
		"write-CategoryAndCollection",
		"delete-CategoryAndCollection",
		"read-discount",
		"write-discount",
		"delete-discount",
		"read-supplierAuth",
		"write-supplierAuth",
		"delete-supplierAuth",
		"read-estelamAndOrder",
		"write-estelamAndOrder",
		"delete-estelamAndOrder",
		"read-sendUnit",
		"write-sendUnit",
		"delete-sendUnit",
		"read-wallet",
		"write-wallet",
		"delete-wallet",
		"read-admin",
		"write-admin",
	}
}

type changePasswordIn struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type updatePassword struct {
	HashPassword string `json:"hashPassword"`
}
