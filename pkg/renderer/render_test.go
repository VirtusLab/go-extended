package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderer_NamedRender_Empty(t *testing.T) {
	Run(t, Test{
		name: "empty render",
		f: func(tt Test) {
			input := ""
			expected := ""

			result, err := New().NamedRender(tt.name, input)

			assert.NoError(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
		},
	})
}

func TestRenderer_NamedRender_Simple(t *testing.T) {
	Run(t, Test{
		name: "simple render",
		f: func(tt Test) {
			input := `key: {{ .value }}
something:
  nested: {{ .something.nested }}`

			expected := `key: some
something:
  nested: val`

			params := map[string]interface{}{
				"value":     "some",
				"something": map[string]string{"nested": "val"},
			}

			result, err := New(WithParameters(params)).NamedRender(tt.name, input)

			assert.NoError(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
		},
	})
}

func TestRenderer_Render_Error(t *testing.T) {
	Run(t, Test{
		name: "parse error",
		f: func(tt Test) {
			input := "{{ wrong+ }}"
			expected := ""

			result, err := New().NamedRender(tt.name, input)

			assert.Error(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
		},
	})
}

func TestRenderer_Render_Validate_Default(t *testing.T) {
	Run(t, Test{
		name: "validation",
		f: func(tt Test) {
			err := New().Validate()
			assert.NoError(t, err, tt.name)
		},
	})
}

type Test struct {
	name string
	f    func(tt Test)
}

func Run(t *testing.T, tt Test) {
	t.Run(tt.name, func(t *testing.T) { tt.f(tt) })
}
