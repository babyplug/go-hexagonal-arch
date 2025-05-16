package domain

var (
	ErrInvalidCredentials  = NewError(401, "invalid credentials")
	ErrDataNotFound        = NewError(404, "data not found")
	ErrTokenCreationFailed = NewError(500, "token creation failed")
	ErrInternal            = NewError(500, "internal server error")
	ErrDuplicateEmail      = NewError(409, "duplicate email")
)

type Error struct {
	Code    int
	Message string
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
