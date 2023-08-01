package entities

import (
	"os"
	"strings"
	"text/template"
)

const ClientTypeCli = "CLI"
const ManifestTypeTask = "MANIFEST_TASK"
const ManifestTypeJob = "MANIFEST_JOB"
const ManifestTypeWorkflow = "MANIFEST_WORKFLOW"

type FunctionMap map[string]interface{}

// TmplCfgFuncMaps is a map of template keywords and their respective functions.
// It's used to map a specific word or keyword in a given manifes into
// a function that will be executed by the template engine.
var TmplCfgFuncMaps = map[string]template.FuncMap{
	"readEnv": {
		"readEnv": os.Getenv,
	},
	"readPWD": {
		"getPwd": os.Getwd,
	},
	"readHome": {
		"getHome": os.UserHomeDir,
	},
	"doReplace": {
		"replace": func(s, old, new string) string {
			return strings.ReplaceAll(s, old, new)
		},
	},
	"doTrimspace": {
		"trimspace": func(s string) string {
			return strings.TrimSpace(s)
		},
	},
}
