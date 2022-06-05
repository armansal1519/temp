package users

type changePasswordDto struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"NewPassword" validate:"required"`
}

type checkForLoginReq struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}
type checkForLoginRes struct {
	PhoneNumber  string `json:"phoneNumber" validate:"required"`
	IsRegistered bool   `json:"isRegistered" validate:"required"`
}

type UserOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	user
}

type user struct {
	PhoneNumber      string   `json:"phoneNumber"`
	FirstName        string   `json:"firstName"`
	LastName         string   `json:"lastName"`
	Email            string   `json:"email"`
	BirthDate        string   `json:"birthDate"`
	NationalCode     string   `json:"nationalCode"`
	Level            string   `json:"level"`
	CreatedAt        int64    `json:"createdAt"`
	LastLogin        int64    `json:"lastLogin"`
	IsAuthenticated  bool     `json:"isAuthenticated"`
	HashRefreshToken string   `json:"hashRefreshToken"`
	WalletAmount     int64    `json:"walletAmount"`
	UserCards        cardInfo `json:"userCards"`
}

type cardInfo struct {
	Number       string `json:"number"`
	BankName     string `json:"bankName"`
	CardUserName string `json:"cardUserName"`
}
type SaveValidationCode struct {
	Key       string `json:"_key"`
	Code      string `json:"code"`
	CreatedAt int64  `json:"createdAt"`
}

type loginAndRegistrationResponse struct {
	AccessToken  string  `json:"accessToken"`
	RefreshToken string  `json:"refreshToken"`
	User         UserOut `json:"user"`
}

type UpdateRefreshToken struct {
	HashRefreshToken string `json:"hashRefreshToken"`
	LastLogin        int64  `json:"lastLogin"`
}
type UpdateRefreshTokenServices struct {
	HashRefreshTokenServices string `json:"HashRefreshTokenServices"`
	LastLogin                int64  `json:"lastLogin"`
}

type LoginDto struct {
	PhoneNumber string `json:"phoneNumber"`
	Code        string `json:"code"`
}

type headlessUser struct {
	State       string `json:"state"`
	PhoneNumber string `json:"phoneNumber"`
}

type updateLastLogin struct {
	LastLogin int64 `json:"lastLogin"`
}

type updateUserDTO struct {
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
	Email        string `json:"email" validate:"required"`
	BirthDate    string `json:"birthDate" validate:"required"`
	NationalCode string `json:"nationalCode" validate:"required"`
	HomeNumber   string `json:"homeNumber" validate:"required"`
}

type AuthenticationDto struct {
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
	BirthDate    string `json:"birthDate" validate:"required"`
	NationalCode string `json:"nationalCode" validate:"required"`
}

type authResponse struct {
	Information struct {
		Type      string `json:"Type"`
		Applicant string `json:"Applicant"`
		Token     string `json:"Token"`
		Credit    string `json:"Credit"`
		Access    string `json:"Access"`
	} `json:"Information"`
	Result struct {
		ID     string `json:"ID"`
		Detail string `json:"Detail"`
	} `json:"Result"`
}
