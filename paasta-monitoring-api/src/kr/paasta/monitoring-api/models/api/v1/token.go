package v1

type (
    TokenDetails struct {
        AccessToken  string `json:"accessToken" validate:"required"`
        RefreshToken string `json:"refreshToken" validate:"required"`
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
        AccessToken  string `json:"accessToken" example:"accessToken" validate:"required"`
        RefreshToken string `json:"refreshToken" example:"refreshToken" validate:"required"`
    }
)
