package matcher

import (
	"testing"

	"github.com/VirtusLab/go-extended/pkg/test"
	"github.com/stretchr/testify/assert"
)

func Test_matcher_MatchGroups(t *testing.T) {
	test.Run(t,
		test.Test{
			Name: "empty expression",
			Fn: func(tt test.Test) {
				m := Must(``)
				value := ""
				want := map[string]string{}

				got, ok := m.MatchGroups(value)

				assert.True(t, ok, "test: '%s', entry: '%s'", tt.Name, value)
				assert.EqualValues(t, want, got, "test: '%s', entry: '%s'", tt.Name, value)
			},
		},
		test.Test{
			Name: "simple expression",
			Fn: func(tt test.Test) {
				m := Must(`^(?P<name>\S+)=(?P<value>\S*)$`)
				value := "test=something"
				want := map[string]string{
					"name":  "test",
					"value": "something",
				}

				got, ok := m.MatchGroups(value)

				assert.True(t, ok, "test: '%s', entry: '%s'", tt.Name, value)
				assert.EqualValues(t, want, got, "test: '%s', entry: '%s'", tt.Name, value)
			},
		},
		test.Test{
			Name: "git url expression - empty",
			Fn: func(tt test.Test) {
				m := Must(`^git@(?P<hostname>[\w\-\.]+):(?P<organisation>[\w\-]+)\/(?P<name>[\w\-]+)\.git$`)
				value := ""
				want := map[string]string{}

				got, ok := m.MatchGroups(value)

				assert.False(t, ok, "test: '%s', entry: '%s'", tt.Name, value)
				assert.EqualValues(t, want, got, "test: '%s', entry: '%s'", tt.Name, value)
			},
		},
		test.Test{
			Name: "git url expression - invalid",
			Fn: func(tt test.Test) {
				m := Must(`^git@(?P<hostname>[\w\-\.]+):(?P<organisation>[\w\-]+)\/(?P<name>[\w\-]+)\.git$`)
				value := "invalid"
				want := map[string]string{}

				got, ok := m.MatchGroups(value)

				assert.False(t, ok, "test: '%s', entry: '%s'", tt.Name, value)
				assert.EqualValues(t, want, got, "test: '%s', entry: '%s'", tt.Name, value)
			},
		},
		test.Test{
			Name: "git url expression - missing extension",
			Fn: func(tt test.Test) {
				m := Must(`^git@(?P<hostname>[\w\-\.]+):(?P<organisation>[\w\-]+)\/(?P<name>[\w\-]+)\.git$`)
				value := "git@something.com:anorg/arepo"
				want := map[string]string{}

				got, ok := m.MatchGroups(value)

				assert.False(t, ok, "test: '%s', entry: '%s'", tt.Name, value)
				assert.EqualValues(t, want, got, "test: '%s', entry: '%s'", tt.Name, value)
			},
		},
		test.Test{
			Name: "git url expression - missing extension",
			Fn: func(tt test.Test) {
				m := Must(`^git@(?P<hostname>[\w\-\.]+):(?P<organisation>[\w\-]+)\/(?P<name>[\w\-]+)\.git$`)
				value := "git@something.com:anorg/arepo.git"
				want := map[string]string{
					"hostname":     "something.com",
					"organisation": "anorg",
					"name":         "arepo",
				}
				got, ok := m.MatchGroups(value)

				assert.True(t, ok, "test: '%s', entry: '%s'", tt.Name, value)
				assert.EqualValues(t, want, got, "test: '%s', entry: '%s'", tt.Name, value)
			},
		},
		test.Test{
			Name: "compile fail",
			Fn: func(tt test.Test) {
				assert.PanicsWithValue(t,
					"regexp: Compile(`<?:[`): error parsing regexp: missing closing ]: `[`",
					func() {
						Must(`<?:[`)
					},
					tt.Name,
				)
			},
		},
	)
}

func Test_matcher_Match(t *testing.T) {
	test.Run(t,
		test.Test{
			Name: "empty",
			Fn: func(tt test.Test) {
				m := Must("^[a-z]+[0-9]+")
				value := ""

				got := m.Match(value)

				assert.False(t, got, tt.Name)
			},
		},
		test.Test{
			Name: "simple match",
			Fn: func(tt test.Test) {
				m := Must("^[a-z]+[0-9]+")
				value := "asdf1234"
				got := m.Match(value)

				assert.True(t, got, tt.Name)

			},
		},
		test.Test{
			Name: "no match",
			Fn: func(tt test.Test) {
				m := Must("^[a-z]+[0-9]+")
				value := "1234asdf"

				got := m.Match(value)

				assert.False(t, got, tt.Name)

			},
		},
	)
}
