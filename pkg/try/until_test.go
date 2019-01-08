package try

import (
	"testing"
	"time"

	"github.com/VirtusLab/go-extended/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestUntil(t *testing.T) {
	test.Run(t,
		test.Test{
			Name: "immediate",
			Fn: func(tt test.Test) {
				something := func() (b bool, e error) {
					return true, nil
				}
				tick := 1 * time.Nanosecond
				timeout := 1 * time.Nanosecond

				ok, err := Until(something, tick, timeout)

				assert.NoError(t, err)
				assert.True(t, ok)
			},
		},
		test.Test{
			Name: "timout",
			Fn: func(tt test.Test) {
				something := func() (b bool, e error) {
					return false, nil
				}
				tick := 1 * time.Nanosecond
				timeout := 3 * time.Nanosecond

				ok, err := Until(something, tick, timeout)

				assert.EqualError(t, err, "timed out after: 3ns, tries: 1")
				assert.False(t, ok)
			},
		},
		test.Test{
			Name: "retired",
			Fn: func(tt test.Test) {
				wantCounter := 3
				counter := 0
				something := func() (b bool, e error) {
					counter = counter + 1
					if counter < wantCounter {
						return false, nil
					}
					return true, nil
				}
				tick := 1 * time.Microsecond
				timeout := 1000 * time.Microsecond

				ok, err := Until(something, tick, timeout)

				assert.NoError(t, err)
				assert.Equal(t, wantCounter, counter)
				assert.True(t, ok)
			},
		},
	)
}
