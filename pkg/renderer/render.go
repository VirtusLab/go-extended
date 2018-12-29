package renderer

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

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

// Renderer structure holds parameters and options
type Renderer struct {
	parameters     map[string]interface{}
	options        []string
	leftDelim      string
	rightDelim     string
	extraFunctions template.FuncMap
}

// New creates a new renderer with the specified parameters and zero or more options
func New() *Renderer {
	r := &Renderer{
		parameters:     map[string]interface{}{},
		options:        []string{MissingKeyErrorOption},
		leftDelim:      LeftDelim,
		rightDelim:     RightDelim,
		extraFunctions: template.FuncMap{},
	}
	return r
}

// Delim mutates Renderer with new left and right delimiters
func (r *Renderer) Delim(left, right string) *Renderer {
	r.leftDelim = left
	r.rightDelim = right
	return r
}

// Functions mutates Renderer with new template functions
func (r *Renderer) Functions(extraFunctions template.FuncMap) *Renderer {
	r.extraFunctions = extraFunctions
	return r
}

// Options mutates Renderer with new template functions
func (r *Renderer) Options(options ...string) *Renderer {
	r.options = options
	return r
}

// Parameters mutates Renderer with new template parameters
func (r *Renderer) Parameters(parameters map[string]interface{}) *Renderer {
	r.parameters = parameters
	return r
}

// Render is a simple rendering function, also used as a custom template function
// to allow in-template recursive rendering, see also NamedRender
func (r *Renderer) Render(rawTemplate string) (string, error) {
	return r.NamedRender("nameless", rawTemplate)
}

// NamedRender is the main rendering function, see also Render, Parameters and ExtraFunctions
func (r *Renderer) NamedRender(templateName, rawTemplate string) (string, error) {
	err := r.Validate()
	if err != nil {
		return "", err
	}
	t, err := r.Parse(templateName, rawTemplate, r.extraFunctions)
	if err != nil {
		return "", err
	}
	out, err := r.Execute(t)
	if err != nil {
		return "", err
	}
	return out, nil
}

// Validate checks the internal state and returns error if necessary
func (r *Renderer) Validate() error {
	if r.parameters == nil {
		return errors.New("unexpected 'nil' parameters")
	}

	if len(r.leftDelim) == 0 {
		return errors.New("unexpected empty leftDelim")
	}
	if len(r.rightDelim) == 0 {
		return errors.New("unexpected empty rightDelim")
	}

	for _, o := range r.options {
		switch o {
		case MissingKeyErrorOption:
		case MissingKeyInvalidOption:
		default:
			return fmt.Errorf("unexpected option: '%s', option must be in: '%s'",
				o, strings.Join([]string{MissingKeyInvalidOption, MissingKeyErrorOption}, ", "))
		}
	}
	return nil
}

// Parse is a basic template parsing function
func (r *Renderer) Parse(templateName, rawTemplate string, extraFunctions template.FuncMap) (*template.Template, error) {
	return template.New(templateName).
		Delims(r.leftDelim, r.rightDelim).
		Funcs(extraFunctions).
		Option(r.options...).
		Parse(rawTemplate)
}

// Execute is a basic template execution function
func (r *Renderer) Execute(t *template.Template) (string, error) {
	var buffer bytes.Buffer
	err := t.Execute(&buffer, r.parameters)
	if err != nil {
		retErr := err
		if e, ok := err.(template.ExecError); ok {
			retErr = fmt.Errorf("error (ExecError) evaluating the template named '%s': %s", e.Name, err)
		}
		return "", retErr
	}
	return buffer.String(), nil
}
