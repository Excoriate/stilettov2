package utils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"text/template"
)

func TestCompileTemplate(t *testing.T) {
	t.Run("should return error when both template content and template file path are empty", func(t *testing.T) {
		opts := TemplateCompilationOpts{}
		_, err := CompileTemplate(opts)

		assert.Errorf(t, err, "both template content and template file path cannot be empty concurrently")
	})

	t.Run("should return error when FuncMap is nil", func(t *testing.T) {
		opts := TemplateCompilationOpts{
			TemplateContent: "some content",
		}
		_, err := CompileTemplate(opts)

		assert.Errorf(t, err, "FuncMap cannot be nil")
	})

	t.Run("should return error when template file path is invalid", func(t *testing.T) {
		opts := TemplateCompilationOpts{
			TemplateFilePath: "some invalid path",
			FuncMap:          template.FuncMap{},
		}
		_, err := CompileTemplate(opts)

		assert.Errorf(t, err, "could not read the template file: open some invalid path: no such file or directory")
	})

	t.Run("should return error when data is nil", func(t *testing.T) {
		opts := TemplateCompilationOpts{
			TemplateContent: "some content",
			FuncMap:         template.FuncMap{},
		}
		_, err := CompileTemplate(opts)

		assert.Errorf(t, err, "data cannot be nil")
	})

	t.Run("should parse a template with an getEnv funcMap", func(t *testing.T) {
		_ = os.Setenv("SHELL", "/bin/bash")

		opts := TemplateCompilationOpts{
			TemplateContent: "{{ getEnv \"SHELL\" }}",
			FuncMap: template.FuncMap{
				"getEnv": os.Getenv,
			},
			Data: make(map[string]string),
		}

		tplCompiled, err := CompileTemplate(opts)
		tplCompiledStr := tplCompiled.String()

		assert.NoErrorf(t, err, "could not compile the template: %s", err)
		assert.NotNilf(t, tplCompiled, "template is nil")
		assert.NotEmptyf(t, tplCompiledStr, "template is empty")
	})

	t.Run("should parse a more complex yaml file, with a getEnv funcMap", func(t *testing.T) {
		_ = os.Setenv("SHELL", "/bin/bash")
		_ = os.Setenv("HOME", "/home/user")
		_ = os.Setenv("USER", "user")

		// Define the template content with placeholders for environment variables
		templateContent := `{
		"shell": "{{ getEnv "SHELL" }}",
		"home": "{{ getEnv "HOME" }}",
		"user": "{{ getEnv "USER" }}"
	}`

		// Create the template compilation options
		opts := TemplateCompilationOpts{
			TemplateContent: templateContent,
			FuncMap: template.FuncMap{
				"getEnv": os.Getenv,
			},
			Data: make(map[string]string),
		}

		tpmCompiled, err := CompileTemplate(opts)
		tplCompiledStr := tpmCompiled.String()

		// Unmarshal the compiled template into the result struct
		var resultStruct struct {
			Shell string `json:"shell"`
			Home  string `json:"home"`
			User  string `json:"user"`
		}

		err = json.Unmarshal([]byte(tplCompiledStr), &resultStruct)

		assert.NotNilf(t, tpmCompiled, "template is nil")
		assert.NotEmptyf(t, tplCompiledStr, "template is empty")
		assert.NoError(t, err, "could not compile the template: %s")
		assert.Equal(t, "/bin/bash", resultStruct.Shell)
		assert.Equal(t, "/home/user", resultStruct.Home)
		assert.Equal(t, "user", resultStruct.User)
	})
}
