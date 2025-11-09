package excel

import (
	"errors"
	"fmt"
)

type Error interface {
	error
	Args(...any) Error
}

type excelError struct {
	text string
	args []any
}

func (r *excelError) Error() string {
	if len(r.args) > 0 {
		return fmt.Sprintf(r.text, r.args...)
	}
	return r.text
}

func (r *excelError) Args(args ...any) Error {
	r.args = args
	return r
}

func newError(text string) Error {
	return &excelError{text: text}
}

var (
	ErrInvalidTarget              = newError("excel: target must be a pointer to a slice")
	ErrTargetMustBeSliceOfStructs = newError("excel: target must be a slice of structs")
	ErrNoHeaders                  = newError("excel: headers not found or empty")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}
