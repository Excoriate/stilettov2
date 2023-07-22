package errors

import "fmt"

const taskConfigErrorPrefix = "Task configuration error: "
const taskExecutionErrorPrefix = "Task execution error: "

type TaskConfigurationError struct {
	Details string
	Err     error
}

func (e *TaskConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", taskConfigErrorPrefix, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", taskExecutionErrorPrefix, e.Details)
}

func NewTaskConfigurationError(details string, err error) *TaskConfigurationError {
	return &TaskConfigurationError{
		Details: details,
		Err:     err,
	}
}

type TaskExecutionError struct {
	Details string
	Err     error
}

func (e *TaskExecutionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", taskExecutionErrorPrefix, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", taskExecutionErrorPrefix, e.Details)
}

func NewTaskExecutionError(details string, err error) *TaskExecutionError {
	return &TaskExecutionError{
		Details: details,
		Err:     err,
	}
}
