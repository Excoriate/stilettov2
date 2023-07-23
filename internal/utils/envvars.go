package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type EnvVars map[string]string

func FetchAllEnvVarsFromHost() (EnvVars, error) {
	result := make(EnvVars)

	for _, env := range os.Environ() {
		keyValue := strings.Split(env, "=")
		result[keyValue[0]] = RemoveDoubleQuotes(keyValue[1])
	}

	return result, nil
}

// FetchEnvVarsWithPrefix fetches environment variables that start with the specified prefix
// and returns an error if any of the variables either do not exist or have an empty value.
func FetchEnvVarsWithPrefix(prefix string) (EnvVars, error) {
	result := make(EnvVars)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]

		if strings.HasPrefix(key, prefix) {
			value := pair[1]
			if value == "" {
				return nil, errors.New(fmt.Sprintf("Environment variable %s has an empty value", key))
			}
			result[key] = RemoveDoubleQuotes(value)
		}
	}

	if len(result) == 0 {
		return nil, errors.New(fmt.Sprintf("No environment variables with the prefix %s found", prefix))
	}

	return result, nil
}

func FetchEnvVarsWithPrefixIncludeEmptyValues(prefix string) (EnvVars, error) {
	result := make(EnvVars)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]

		if strings.HasPrefix(key, prefix) {
			value := pair[1]
			result[key] = RemoveDoubleQuotes(value)
		}
	}

	if len(result) == 0 {
		return nil, errors.New(fmt.Sprintf("No environment variables with the prefix %s found", prefix))
	}

	return result, nil
}

func MergeEnvVars(envVars ...EnvVars) EnvVars {
	result := make(EnvVars)

	for _, env := range envVars {
		for key, value := range env {
			if key != "" && value != "" {
				result[key] = RemoveDoubleQuotes(value)
			}
		}
	}

	return result
}

func GetEnvVarsFromDotFile(filepath string) (map[string]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf(" .env file %s is not in the correct format", filepath)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = RemoveDoubleQuotes(value)

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(env) == 0 {
		return nil, fmt.Errorf(" .env file %s is empty", filepath)
	}

	return env, nil
}
