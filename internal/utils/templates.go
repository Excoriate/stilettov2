package utils

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type TemplateCompilationOpts struct {
	TemplateContent  string
	TemplateFilePath string
	Data             interface{}
	TemplateName     string
	FuncMap          template.FuncMap
}

// CompileTemplate compiles a text template with the given data and function map.
// It returns the compiled template as a bytes.Buffer.
func CompileTemplate(opts TemplateCompilationOpts) (bytes.Buffer, error) {
	if opts.TemplateContent == "" && opts.TemplateFilePath == "" {
		return bytes.Buffer{},
			fmt.Errorf("both template content and template file path cannot be empty concurrently")
	}

	if opts.FuncMap == nil {
		return bytes.Buffer{},
			fmt.Errorf("FuncMap cannot be nil")
	}

	var content string
	if opts.TemplateFilePath != "" {
		file, err := os.ReadFile(opts.TemplateFilePath)
		if err != nil {
			return bytes.Buffer{}, fmt.Errorf("could not read the template file: %v", err)
		}
		content = string(file)
	} else {
		content = opts.TemplateContent
	}

	if opts.TemplateName == "" {
		opts.TemplateName = "tmpl" // default name
	}

	if opts.Data == nil {
		return bytes.Buffer{}, fmt.Errorf("data cannot be nil")
	}

	tmpl := template.New(opts.TemplateName).Funcs(opts.FuncMap)

	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("could not parse the template: %s", err)
	}

	var tpl bytes.Buffer

	if err = tmpl.Execute(&tpl, opts.Data); err != nil {
		return bytes.Buffer{}, fmt.Errorf("could not execute the template: %s", err)
	}

	return tpl, nil
}
