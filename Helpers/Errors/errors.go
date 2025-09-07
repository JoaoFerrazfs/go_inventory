package errors

const DEFAULT_ERROR_CODE = 500

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(message string, code ...int) *AppError {
	status := DEFAULT_ERROR_CODE
	if len(code) > 0 {
		status = code[0]
	}

	return &AppError{
		Code:    status,
		Message: message,
	}
}
