package renderer

import (
	"testing"

	"github.com/VirtusLab/go-extended/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRenderer_NamedRender(t *testing.T) {
	test.Run(t,
		test.Test{
			Name: "empty render",
			Fn: func(_ test.Test) {
				input := ""
				expected := ""

				result, err := New().NamedRender("test", input)

				assert.NoError(t, err)
				assert.Equal(t, expected, result)
			},
		},
		test.Test{
			Name: "simple render",
			Fn: func(_ test.Test) {
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
			},
		},
	)
}

func TestRenderer_Render(t *testing.T) {
	test.Run(t,
		test.Test{
			Name: "parse error",
			Fn: func(_ test.Test) {
				input := "{{ wrong+ }}"
				expected := ""

				result, err := New().NamedRender("test", input)

				assert.Error(t, err)
				assert.Equal(t, expected, result)
			},
		},
		test.Test{
			Name: "validation",
			Fn: func(_ test.Test) {
				err := New().Validate()
				assert.NoError(t, err)
			},
		},
	)
}
