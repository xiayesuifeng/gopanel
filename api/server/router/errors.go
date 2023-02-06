package router

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}
