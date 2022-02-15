package workerAuth

type LoginDto struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Password    string `json:"password"`
}

type refreshTokenDto struct {
	RefreshToken string `json:"refreshToken"`
}
