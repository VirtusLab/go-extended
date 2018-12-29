package matcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_matcher_MatchGroups(t *testing.T) {
	type fields struct {
		matcher Matcher
	}
	type args struct {
		value string
	}
	type test struct {
		name   string
		fields fields
		args   args
		want   map[string]string
		wantOk bool
		f      func(tt test)
	}
	standard := func(tt test) {
		actual, ok := tt.fields.matcher.MatchGroups(tt.args.value)
		assert.Equal(t, tt.wantOk, ok, "test: '%s', entry: '%s'", tt.name, tt.args.value)
		assert.EqualValues(t, tt.want, actual, "test: '%s', entry: '%s'", tt.name, tt.args.value)
	}
	tests := []test{
		{
			name: "empty expression",
			fields: fields{
				matcher: Must(``),
			},
			args: args{
				value: "",
			},
			want:   map[string]string{},
			wantOk: true,
			f:      standard,
		},
		{
			name: "simple expression",
			fields: fields{
				matcher: Must(`^(?P<name>\S+)=(?P<value>\S*)$`),
			},
			args: args{
				value: "test=something",
			},
			want: map[string]string{
				"name":  "test",
				"value": "something",
			},
			wantOk: true,
			f:      standard,
		},
		{
			name: "git url expression - empty",
			fields: fields{
				matcher: Must(`^git@(?P<hostname>[\w\-\.]+):(?P<organisation>[\w\-]+)\/(?P<name>[\w\-]+)\.git$`),
			},
			args: args{
				value: "",
			},
			want:   map[string]string{},
			wantOk: false,
			f:      standard,
		},
		{
			name: "git url expression - invalid",
			fields: fields{
				matcher: Must(`^git@(?P<hostname>[\w\-\.]+):(?P<organisation>[\w\-]+)\/(?P<name>[\w\-]+)\.git$`),
			},
			args: args{
				value: "invalid",
			},
			want:   map[string]string{},
			wantOk: false,
			f:      standard,
		},
		{
			name: "git url expression - missing extension",
			fields: fields{
				matcher: Must(`^git@(?P<hostname>[\w\-\.]+):(?P<organisation>[\w\-]+)\/(?P<name>[\w\-]+)\.git$`),
			},
			args: args{
				value: "git@something.com:anorg/arepo",
			},
			want:   map[string]string{},
			wantOk: false,
			f:      standard,
		},
		{
			name: "git url expression - missing extension",
			fields: fields{
				matcher: Must(`^git@(?P<hostname>[\w\-\.]+):(?P<organisation>[\w\-]+)\/(?P<name>[\w\-]+)\.git$`),
			},
			args: args{
				value: "git@something.com:anorg/arepo.git",
			},
			want: map[string]string{
				"hostname":     "something.com",
				"organisation": "anorg",
				"name":         "arepo",
			},
			wantOk: true,
			f:      standard,
		},
		{
			name: "compile fail",
			f: func(tt test) {
				assert.PanicsWithValue(t,
					"regexp: Compile(`<?:[`): error parsing regexp: missing closing ]: `[`",
					func() {
						Must(`<?:[`)
					},
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { tt.f(tt) })
	}
}

func Test_matcher_Match(t *testing.T) {
	type fields struct {
		matcher Matcher
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "empty",
			fields: fields{matcher: Must("^[a-z]+[0-9]+")},
			args: args{
				value: "",
			},
			want: false,
		}, {
			name:   "simple match",
			fields: fields{matcher: Must("^[a-z]+[0-9]+")},
			args: args{
				value: "asdf1234",
			},
			want: true,
		},
		{
			name:   "no match",
			fields: fields{matcher: Must("^[a-z]+[0-9]+")},
			args: args{
				value: "1234asdf",
			},
			want: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tt.name), func(t *testing.T) {
			m := tt.fields.matcher
			got := m.Match(tt.args.value)
			assert.Equal(t, tt.want, got, tt.name)
		})
	}
}
