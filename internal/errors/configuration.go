package errors

import "fmt"

type ConfigurationError struct {
	Details string
	Err     error
}

func (e *ConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("PipelineCfg configuration error: %s: %s", e.Details, e.Err.Error())
	}
	return fmt.Sprintf("PipelineCfg configuration error: %s", e.Details)
}

func NewConfigurationError(details string, err error) *ConfigurationError {
	return &ConfigurationError{
		Details: fmt.Sprintf("Unable to start pipeline instance %s", details),
		Err:     err,
	}
}
