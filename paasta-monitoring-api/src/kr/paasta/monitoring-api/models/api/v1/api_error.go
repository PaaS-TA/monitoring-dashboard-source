package v1

type ApiError struct {
	OriginError error
	Code string
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}