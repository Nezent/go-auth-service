package shared

// Structured common domain errors
var (
	ErrNotFound = &DomainError{
		Code:       "NOT_FOUND",
		StatusCode: 404,
		Message:    "Resource not found",
	}
	ErrUnauthorized = &DomainError{
		Code:       "UNAUTHORIZED",
		StatusCode: 401,
		Message:    "Unauthorized access",
	}
	ErrValidationFailed = &DomainError{
		Code:       "VALIDATION_FAILED",
		StatusCode: 400,
		Message:    "Validation failed",
	}
)

// DomainError represents a domain-specific error with code, status, and message.
type DomainError struct {
	Code       string `json:"code"`        // e.g. "USER_NOT_FOUND"
	StatusCode int    `json:"status_code"` // e.g. 404
	Message    string `json:"message"`     // e.g. "User not found"
}

// Error implements the error interface.
func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new DomainError with code, status, and message.
func NewDomainError(code string, status int, message string) *DomainError {
	return &DomainError{
		Code:       code,
		StatusCode: status,
		Message:    message,
	}
}
