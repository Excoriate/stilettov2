package errors

import "fmt"

const runnerConfigErrorPrefix = "Runner configuration error: "
const runnerExecutionErrorPrefix = "Runner execution error: "

type RunnerConfigurationError struct {
	Details string
	Err     error
}

func (e *RunnerConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", runnerConfigErrorPrefix, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", runnerExecutionErrorPrefix, e.Details)
}

func NewRunnerConfigurationError(details string, err error) *RunnerConfigurationError {
	return &RunnerConfigurationError{
		Details: details,
		Err:     err,
	}
}

type RunnerExecutionError struct {
	Details string
	Err     error
}

func (e *RunnerExecutionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", runnerExecutionErrorPrefix, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", runnerExecutionErrorPrefix, e.Details)
}

func NewRunnerExecutionError(details string, err error) *RunnerExecutionError {
	return &RunnerExecutionError{
		Details: details,
		Err:     err,
	}
}
