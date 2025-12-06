package response

type AppError struct {
	Code    int
	Key     string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, key string, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Key:     key,
		Message: msg,
		Err:     err,
	}
}
