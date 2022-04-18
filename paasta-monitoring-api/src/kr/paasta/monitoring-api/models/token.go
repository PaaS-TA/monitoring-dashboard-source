package models

type (
	TokenDetails struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		AccessUuid   string `json:"accessUuid"`
		RefreshUuid  string `json:"refreshUuid"`
		AtExpires    int64  `json:"atExpires"`
		RtExpires    int64  `json:"rtExpires"`
	}
)
