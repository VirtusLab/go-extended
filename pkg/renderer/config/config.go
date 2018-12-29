package config

import "text/template"

const (
	// MissingKeyInvalidOption is the renderer option to continue execution on missing key and print "<no value>"
	MissingKeyInvalidOption = "missingkey=invalid"
	// MissingKeyErrorOption is the renderer option to stops execution immediately with an error on missing key
	MissingKeyErrorOption = "missingkey=error"
	// LeftDelim is the default left template delimiter
	LeftDelim = "{{"
	// RightDelim is the default right template delimiter
	RightDelim = "}}"
)

// Config holds the renderer configuration
type Config struct {
	Parameters     map[string]interface{}
	Options        []string
	LeftDelim      string
	RightDelim     string
	ExtraFunctions template.FuncMap
}
