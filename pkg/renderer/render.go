package renderer

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/VirtusLab/go-extended/pkg/renderer/config"
)

// Renderer allows for parameterised text template rendering
type Renderer interface {
	Configuration() config.Config
	Reconfigure(configurators ...func(*config.Config))

	Render(rawTemplate string) (string, error)
	NamedRender(templateName, rawTemplate string) (string, error)
	Validate() error
	Parse(templateName, rawTemplate string, extraFunctions template.FuncMap) (*template.Template, error)
	Execute(t *template.Template) (string, error)
}

type renderer struct {
	config *config.Config
}

// New creates a new renderer with the specified parameters and zero or more options
func New(configurators ...func(*config.Config)) Renderer {
	r := &renderer{
		config: &config.Config{
			Parameters:     map[string]interface{}{},
			Options:        []string{config.MissingKeyErrorOption},
			LeftDelim:      config.LeftDelim,
			RightDelim:     config.RightDelim,
			ExtraFunctions: template.FuncMap{},
		},
	}
	r.Reconfigure(configurators...)
	return r
}

// Reconfigure mutates the configuration state with the given configurators
func (r *renderer) Reconfigure(configurators ...func(*config.Config)) {
	for _, c := range configurators {
		c(r.config)
	}
}

// Configuration returns current configuration
func (r *renderer) Configuration() config.Config {
	return *r.config
}

// WithParameters mutates Renderer configuration with new template parameters
func WithParameters(parameters map[string]interface{}) func(*config.Config) {
	return func(c *config.Config) {
		c.Parameters = parameters
	}
}

// WithOptions mutates Renderer configuration with new template functions
func WithOptions(options ...string) func(*config.Config) {
	return func(c *config.Config) {
		c.Options = options
	}
}

// WithDelim mutates Renderer configuration with new left and right delimiters
func WithDelim(left, right string) func(*config.Config) {
	return func(c *config.Config) {
		c.LeftDelim = left
		c.RightDelim = right
	}
}

// WithFunctions mutates Renderer configuration with new template functions
func WithFunctions(extraFunctions template.FuncMap) func(*config.Config) {
	return func(c *config.Config) {
		c.ExtraFunctions = extraFunctions
	}
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
	t, err := r.Parse(templateName, rawTemplate, r.config.ExtraFunctions)
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
	if r.config.Parameters == nil {
		return errors.New("unexpected 'nil' parameters")
	}

	if len(r.config.LeftDelim) == 0 {
		return errors.New("unexpected empty leftDelim")
	}
	if len(r.config.RightDelim) == 0 {
		return errors.New("unexpected empty rightDelim")
	}

	for _, o := range r.config.Options {
		switch o {
		case config.MissingKeyErrorOption:
		case config.MissingKeyInvalidOption:
		default:
			return fmt.Errorf("unexpected option: '%s', option must be in: '%s'",
				o, strings.Join([]string{
					config.MissingKeyInvalidOption, config.MissingKeyErrorOption,
				}, ", "))
		}
	}
	return nil
}

// Parse is a basic template parsing function
func (r *renderer) Parse(templateName, rawTemplate string, extraFunctions template.FuncMap) (*template.Template, error) {
	return template.New(templateName).
		Delims(r.config.LeftDelim, r.config.RightDelim).
		Funcs(extraFunctions).
		Option(r.config.Options...).
		Parse(rawTemplate)
}

// Execute is a basic template execution function
func (r *renderer) Execute(t *template.Template) (string, error) {
	var buffer bytes.Buffer
	err := t.Execute(&buffer, r.config.Parameters)
	if err != nil {
		retErr := err
		if e, ok := err.(template.ExecError); ok {
			retErr = fmt.Errorf("error (ExecError) evaluating the template named '%s': %s", e.Name, err)
		}
		return "", retErr
	}
	return buffer.String(), nil
}
