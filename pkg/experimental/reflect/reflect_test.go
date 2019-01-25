package reflect

import (
	"testing"

	"github.com/VirtusLab/go-extended/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestCheckEmpty(t *testing.T) {
	test.Run(t,
		test.Test{
			Name: "happy",
			Fn: func(tt test.Test) {
				err := CheckEmpty("not empty")
				assert.NoError(t, err)
			},
		},
		test.Test{
			Name: "empty",
			More: []test.Test{
				{
					Name: "nil",
					Fn: func(tt test.Test) {
						err := CheckEmpty(nil)
						assert.EqualError(t, err, "nil value")
						assert.IsType(t, &UnexpectedNilValue{}, err)
					},
				},
				{
					Name: "string",
					Fn: func(tt test.Test) {
						err := CheckEmpty("")
						assert.EqualError(t, err, "empty value")
						assert.IsType(t, &UnexpectedEmptyValue{}, err)
					},
				},
				{
					Name: "slice",
					Fn: func(tt test.Test) {
						err := CheckEmpty([]int{})
						assert.EqualError(t, err, "empty value")
						assert.IsType(t, &UnexpectedEmptyValue{}, err)
					},
				},
				{
					Name: "array",
					Fn: func(tt test.Test) {
						err := CheckEmpty([0]int{})
						assert.EqualError(t, err, "empty value")
						assert.IsType(t, &UnexpectedEmptyValue{}, err)
					},
				},
				{
					Name: "map",
					Fn: func(tt test.Test) {
						err := CheckEmpty(map[string]interface{}{})
						assert.EqualError(t, err, "empty value")
						assert.IsType(t, &UnexpectedEmptyValue{}, err)
					},
				},
				{
					Name: "chan",
					Fn: func(tt test.Test) {
						err := CheckEmpty(make(chan string))
						assert.EqualError(t, err, "empty value")
						assert.IsType(t, &UnexpectedEmptyValue{}, err)
					},
				},
				{
					Name: "int",
					Fn: func(tt test.Test) {
						err := CheckEmpty(0)
						assert.EqualError(t, err, "unsupported type of 'int' with value '0'")
						assert.IsType(t, &UnsupportedValueType{}, err)
					},
				},
			},
		},
	)
}

func TestLen(t *testing.T) {

}
