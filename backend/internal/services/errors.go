package services

type Error interface {
	// To satisfy the generic error interface.
	error

	// Used to classify the error - should be an HTTP status code.
	Code() int

	// Returns the error details message.
	Message() string
}

func NewError(code int, message string, origErr error) Error {
	var errs []error
	if origErr != nil {
		errs = append(errs, origErr)
	}
	return newBaseError(code, message, errs)
}

func FormatError(code int, message string, errs []error) string {
	if len(errs) > 0 {
		message = message + ": " + errorList(errs)
	}
	return message
}

func errorList(errs []error) string {
	if len(errs) == 0 {
		return ""
	}
	var str string
	for i, err := range errs {
		if i > 0 {
			str = str + ", "
		}
		str = str + err.Error()
	}
	return str
}

type baseError struct {
	code    int
	message string
	errs    []error
}

func newBaseError(code int, message string, errs []error) *baseError {
	return &baseError{
		code:    code,
		message: message,
		errs:    errs,
	}
}

func (b baseError) Error() string {
	size := len(b.errs)
	if size > 0 {
		return FormatError(b.code, b.message, b.errs)
	}

	return FormatError(b.code, b.message, nil)
}

func (b baseError) String() string {
	return b.Error()
}

func (b baseError) Code() int {
	return b.code
}

func (b baseError) Message() string {
	return b.message
}
