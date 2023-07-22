package errors

import "fmt"

const argumentOrInputErr = "Argument or input error: "

type ArgumentOrInputError struct {
	Details string
	Err     error
}

func (e *ArgumentOrInputError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", argumentOrInputErr, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", argumentOrInputErr, e.Details)
}

func NewArgumentError(details string, err error) *ArgumentOrInputError {
	return &ArgumentOrInputError{
		Details: details,
		Err:     err,
	}
}
