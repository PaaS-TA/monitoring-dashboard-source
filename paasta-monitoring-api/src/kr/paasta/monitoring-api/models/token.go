package models

type (
    TokenDetails struct {
        AccessToken  string `json:"accessToken" validate:"required"`
        RefreshToken string `json:"refreshToken" validate:"required"`
        AccessUuid   string `json:"accessUuid"`
        RefreshUuid  string `json:"refreshUuid"`
        AtExpires    int64  `json:"atExpires"`
        RtExpires    int64  `json:"rtExpires"`
    }
)
