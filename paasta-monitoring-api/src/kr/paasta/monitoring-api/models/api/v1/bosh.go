package v1

type (
	BoshSummary struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	Bosh struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	BoshProcess struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	BoshChart struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	BoshLog struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)
