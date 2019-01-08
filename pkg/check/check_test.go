package check

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/VirtusLab/go-extended/pkg/test"
)

func TestIsValidJavaScriptIdentifier(t *testing.T) {
	test.Run(t,
		test.Test{
			Name: "simple",
			Fn: func(tt test.Test) {
				value := "simpleName"
				err := IsValidJavaScriptIdentifier(value)
				assert.NoError(t, err, tt.Name)
			},
		},
		test.Test{
			Name: "dollar",
			Fn: func(tt test.Test) {
				value := "$"
				err := IsValidJavaScriptIdentifier(value)
				assert.NoError(t, err, tt.Name)
			},
		},
		test.Test{
			Name: "underscore",
			Fn: func(tt test.Test) {
				value := "_"
				err := IsValidJavaScriptIdentifier(value)
				assert.NoError(t, err, tt.Name)
			},
		},
		test.Test{
			Name: "dollar underscore",
			Fn: func(tt test.Test) {
				value := "$_"
				err := IsValidJavaScriptIdentifier(value)
				assert.NoError(t, err, tt.Name)
			},
		},
		test.Test{
			Name: "underscore dollar",
			Fn: func(tt test.Test) {
				value := "_$"
				err := IsValidJavaScriptIdentifier(value)
				assert.NoError(t, err, tt.Name)
			},
		},
		test.Test{
			Name: "complex",
			Fn: func(tt test.Test) {
				value := "$$_1cOmPlEx0_$$_"
				err := IsValidJavaScriptIdentifier(value)
				assert.NoError(t, err, tt.Name)
			},
		},
		test.Test{
			Name: "empty",
			Fn: func(tt test.Test) {
				value := ""
				wantErr := "must be a valid JavaScript identifier, '' does not match pattern '^[a-zA-Z_$][a-zA-Z0-9_$]*$'"
				err := IsValidJavaScriptIdentifier(value)
				assert.EqualError(t, err, wantErr, tt.Name)
			},
		},
		test.Test{
			Name: "invalid with dash",
			Fn: func(tt test.Test) {
				value := "invalid-with-dash"
				wantErr := "must be a valid JavaScript identifier, 'invalid-with-dash' does not match pattern '^[a-zA-Z_$][a-zA-Z0-9_$]*$'"
				err := IsValidJavaScriptIdentifier(value)
				assert.EqualError(t, err, wantErr, tt.Name)
			},
		},
		test.Test{
			Name: "starts with numbers",
			Fn: func(tt test.Test) {
				value := "123notSoGood"
				wantErr := "must be a valid JavaScript identifier, '123notSoGood' does not match pattern '^[a-zA-Z_$][a-zA-Z0-9_$]*$'"
				err := IsValidJavaScriptIdentifier(value)
				assert.EqualError(t, err, wantErr, tt.Name)
			},
		},
	)
}
