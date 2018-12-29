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

// Renderer allows for parameterised text template rendering
type Renderer interface {
	Parameters() map[string]interface{}
	Options() []string
	Delim() (string, string)
	Functions() template.FuncMap

	WithParameters(parameters map[string]interface{}) Renderer
	WithOptions(options ...string) Renderer
	WithDelim(left, right string) Renderer
	WithFunctions(extraFunctions template.FuncMap) Renderer

	Render(rawTemplate string) (string, error)
	NamedRender(templateName, rawTemplate string) (string, error)
	Validate() error
	Parse(templateName, rawTemplate string, extraFunctions template.FuncMap) (*template.Template, error)
	Execute(t *template.Template) (string, error)
}

type renderer struct {
	parameters     map[string]interface{}
	options        []string
	leftDelim      string
	rightDelim     string
	extraFunctions template.FuncMap
}

// New creates a new renderer with the specified parameters and zero or more options
func New() Renderer {
	return &renderer{
		parameters:     map[string]interface{}{},
		options:        []string{MissingKeyErrorOption},
		leftDelim:      LeftDelim,
		rightDelim:     RightDelim,
		extraFunctions: template.FuncMap{},
	}
}

// Parameters returns current parameters
func (r *renderer) Parameters() map[string]interface{} {
	return r.parameters
}

// Options returns current options
func (r *renderer) Options() []string {
	return r.options
}

// Delim returns current delimiters
func (r *renderer) Delim() (string, string) {
	return r.leftDelim, r.rightDelim
}

// Functions returns current extra functions
func (r *renderer) Functions() template.FuncMap {
	return r.extraFunctions
}

// WithParameters mutates Renderer with new template parameters
func (r *renderer) WithParameters(parameters map[string]interface{}) Renderer {
	r.parameters = parameters
	return r
}

// WithOptions mutates Renderer with new template functions
func (r *renderer) WithOptions(options ...string) Renderer {
	r.options = options
	return r
}

// WithDelim mutates Renderer with new left and right delimiters
func (r *renderer) WithDelim(left, right string) Renderer {
	r.leftDelim = left
	r.rightDelim = right
	return r
}

// WithFunctions mutates Renderer with new template functions
func (r *renderer) WithFunctions(extraFunctions template.FuncMap) Renderer {
	r.extraFunctions = extraFunctions
	return r
}

// Render is a simple rendering function, also used as a custom template function
// to allow in-template recursive rendering, see also NamedRender
func (r renderer) Render(rawTemplate string) (string, error) {
	return r.NamedRender("nameless", rawTemplate)
}

// NamedRender is the main rendering function, see also Render, WithParameters and Functions
func (r *renderer) NamedRender(templateName, rawTemplate string) (string, error) {
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
func (r *renderer) Validate() error {
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
func (r *renderer) Parse(templateName, rawTemplate string, extraFunctions template.FuncMap) (*template.Template, error) {
	return template.New(templateName).
		Delims(r.leftDelim, r.rightDelim).
		Funcs(extraFunctions).
		Option(r.options...).
		Parse(rawTemplate)
}

// Execute is a basic template execution function
func (r *renderer) Execute(t *template.Template) (string, error) {
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
