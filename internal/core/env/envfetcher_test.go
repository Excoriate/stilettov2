package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnvVarsByType(t *testing.T) {
	t.Run("FailIfNotSet", func(t *testing.T) {
		args := VarsByTypeOpt{
			FailIfNotSet:          true,
			Prefix:                "_TEST",
			RequiredEnvVars:       []string{},
			IgnoreIfNotSetOrEmpty: []string{},
		}

		_ = os.Setenv("_TEST_VAR1", "value1")
		_ = os.Setenv("_TEST_VAR2", "value2")

		envVars, err := GetEnvVarsByType(args)

		assert.NoError(t, err, "The GetEnvVarsByType should not return an error")
		assert.NotNil(t, envVars, "The GetEnvVarsByType should return a map")
		assert.Equal(t, "value1", envVars["_TEST_VAR1"],
			"The GetEnvVarsByType should return a valid map")
		assert.Equal(t, "value2", envVars["_TEST_VAR2"],
			"The GetEnvVarsByType should return a valid map")
	})

	t.Run("EnvVarsNotSetButAreIgnored", func(t *testing.T) {
		args := VarsByTypeOpt{
			FailIfNotSet:          true,
			Prefix:                "_TEST",
			RequiredEnvVars:       []string{"_TEST_VAR1"},
			IgnoreIfNotSetOrEmpty: []string{"_TEST_VAR2"},
		}

		_ = os.Setenv("_TEST_VAR1", "value1")
		_ = os.Setenv("_TEST_VAR2", "")

		envVars, err := GetEnvVarsByType(args)

		assert.NoError(t, err, "The GetEnvVarsByType should not return an error")
		assert.NotNil(t, envVars, "The GetEnvVarsByType should return a map")
		assert.Equal(t, "value1", envVars["_TEST_VAR1"], "The GetEnvVarsByType should return a valid map")
		assert.Equal(t, "", envVars["_TEST_VAR2"], "The GetEnvVarsByType should return a valid map")
	})

	t.Run("EnvVarsAreRequiredAndNotSetOrEmpty", func(t *testing.T) {
		args := VarsByTypeOpt{
			FailIfNotSet:          true,
			Prefix:                "_TEST",
			RequiredEnvVars:       []string{"_TEST_VAR1", "_TEST_VAR2"},
			IgnoreIfNotSetOrEmpty: []string{},
		}

		_ = os.Setenv("_TEST_VAR1", "value1")

		envVars, err := GetEnvVarsByType(args)

		assert.Error(t, err, "The GetEnvVarsByType should return an error")
		assert.Nil(t, envVars, "The GetEnvVarsByType should return a nil map")
	})

	t.Run("RemoveEnvVarsIfFound", func(t *testing.T) {
		args := VarsByTypeOpt{
			FailIfNotSet:          true,
			Prefix:                "_TEST",
			RequiredEnvVars:       []string{},
			IgnoreIfNotSetOrEmpty: []string{},
			RemoveEnvVarsIfFound:  []string{"_TEST_VAR1"},
		}

		_ = os.Setenv("_TEST_VAR1", "value1")
		_ = os.Setenv("_TEST_VAR2", "value2")

		envVars, err := GetEnvVarsByType(args)

		assert.NoError(t, err, "The GetEnvVarsByType should not return an error")
		assert.NotNil(t, envVars, "The GetEnvVarsByType should return a map")
		assert.Equal(t, "value2", envVars["_TEST_VAR2"], "The GetEnvVarsByType should return a valid map")
		assert.Equal(t, "", envVars["_TEST_VAR1"], "The GetEnvVarsByType should return a valid map")
	})
}
