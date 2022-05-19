package v1

type (
	TokenDetails struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		AccessUuid   string `json:"accessUuid"`
		RefreshUuid  string `json:"refreshUuid"`
		AtExpires    int64  `json:"atExpires"`
		RtExpires    int64  `json:"rtExpires"`
	}

	CreateToken struct {
		Username string `json:"username" example:"username" validate:"required"`
		Password string `json:"password" example:"password" validate:"required"`
	}

	RefreshToken struct {
		RefreshToken string `json:"refreshToken" example:"refreshToken" validate:"required"`
	}
)
