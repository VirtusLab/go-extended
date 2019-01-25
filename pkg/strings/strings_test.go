package strings

import (
	"testing"

	"github.com/VirtusLab/go-extended/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestEllipsis(t *testing.T) {
	test.Run(t,
		test.Test{
			Name: "empty",
			Fn: func(tt test.Test) {
				s := ""
				max := 0
				want := ""
				got := Ellipsis(s, max)
				assert.Equal(t, want, got)
			},
		},
		test.Test{
			Name: "small",
			Fn: func(tt test.Test) {
				s := "12345"
				max := 4
				want := "1234"
				got := Ellipsis(s, max)
				assert.Equal(t, want, got)
			},
		},
		test.Test{
			Name: "edge",
			Fn: func(tt test.Test) {
				s := "1234567890"
				max := 10
				want := "1234567890"
				got := Ellipsis(s, max)
				assert.Equal(t, want, got)
			},
		},
		test.Test{
			Name: "happy",
			Fn: func(tt test.Test) {
				s := "1234"
				max := 10
				want := "1234"
				got := Ellipsis(s, max)
				assert.Equal(t, want, got)
			},
		},
		test.Test{
			Name: "normal+1",
			Fn: func(tt test.Test) {
				s := "aaaaaaaaaa1"
				max := 10
				want := "aaaaaaa..."
				got := Ellipsis(s, max)
				assert.Equal(t, want, got)
			},
		},
		test.Test{
			Name: "normal+2",
			Fn: func(tt test.Test) {
				s := "aaaaaaaaaa12"
				max := 10
				want := "aaaaaaa..."
				got := Ellipsis(s, max)
				assert.Equal(t, want, got)
			},
		},
		test.Test{
			Name: "normal+3",
			Fn: func(tt test.Test) {
				s := "aaaaaaaaaa123"
				max := 10
				want := "aaaaaaa..."
				got := Ellipsis(s, max)
				assert.Equal(t, want, got)
			},
		},
		test.Test{
			Name: "normal+4",
			Fn: func(tt test.Test) {
				s := "aaaaaaaaaa1234"
				max := 10
				want := "aaaaaaa..."
				got := Ellipsis(s, max)
				assert.Equal(t, want, got)
			},
		},
	)
}
