package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderer_NamedRender_Empty(t *testing.T) {
	t.Run("empty render", func(_ *testing.T) {
		input := ""
		expected := ""

		result, err := New().NamedRender("test", input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestRenderer_NamedRender_Simple(t *testing.T) {
	t.Run("simple render", func(_ *testing.T) {
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

		result, err := New(WithParameters(params)).NamedRender("test", input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestRenderer_Render_Error(t *testing.T) {
	t.Run("parse error", func(_ *testing.T) {
		input := "{{ wrong+ }}"
		expected := ""

		result, err := New().NamedRender("test", input)

		assert.Error(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestRenderer_Render_Validate_Default(t *testing.T) {
	t.Run("validation", func(_ *testing.T) {
		err := New().Validate()
		assert.NoError(t, err)
	})
}
