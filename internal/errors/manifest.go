package errors

import "fmt"

const manifestErrorPrefix = "Manifest error: "

type ManifestError struct {
	Details string
	Err     error
}

func (e *ManifestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", manifestErrorPrefix, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", manifestErrorPrefix, e.Details)
}

func NewManifestError(details string, err error) *ManifestError {
	return &ManifestError{
		Details: details,
		Err:     err,
	}
}
