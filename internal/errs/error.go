package errs

import "fmt"

type AppError struct {
	Code    string // used for programmatic checking
	Message string // client-friendly message
	Err     error  // original error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", e.Code, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
