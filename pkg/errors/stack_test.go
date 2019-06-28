/*
Copyright (c) 2015, Dave Cheney <dave@cheney.net>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package errors

import (
	"fmt"
	"runtime"
	"testing"
)

var initpc = caller()

type X struct{}

// val returns a Frame pointing to itself.
func (x X) val() Frame {
	return caller()
}

// ptr returns a Frame pointing to itself.
func (x *X) ptr() Frame {
	return caller()
}

func TestFrameFormat(t *testing.T) {
	var tests = []struct {
		Frame
		format string
		want   string
	}{{
		initpc,
		"%s",
		"stack_test.go",
	}, {
		initpc,
		"%+s",
		"github.com/VirtusLab/go-extended/pkg/errors.init\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go",
	}, {
		0,
		"%s",
		"unknown",
	}, {
		0,
		"%+s",
		"unknown",
	}, {
		initpc,
		"%d",
		"35",
	}, {
		0,
		"%d",
		"0",
	}, {
		initpc,
		"%n",
		"init",
	}, {
		func() Frame {
			var x X
			return x.ptr()
		}(),
		"%n",
		`\(\*X\).ptr`,
	}, {
		func() Frame {
			var x X
			return x.val()
		}(),
		"%n",
		"X.val",
	}, {
		0,
		"%n",
		"",
	}, {
		initpc,
		"%v",
		"stack_test.go:35",
	}, {
		initpc,
		"%+v",
		"github.com/VirtusLab/go-extended/pkg/errors.init\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:35",
	}, {
		0,
		"%v",
		"unknown:0",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.Frame, tt.format, tt.want)
	}
}

func TestFuncname(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"", ""},
		{"runtime.main", "main"},
		{"github.com/VirtusLab/go-extended/pkg/errors.funcName", "funcName"},
		{"funcName", "funcName"},
		{"io.copyBuffer", "copyBuffer"},
		{"main.(*R).Write", "(*R).Write"},
	}

	for _, tt := range tests {
		got := funcName(tt.name)
		want := tt.want
		if got != want {
			t.Errorf("funcName(%q): want: %q, got %q", tt.name, want, got)
		}
	}
}

func TestStackTrace(t *testing.T) {
	tests := []struct {
		err  error
		want []string
	}{{
		New("ooh1"), []string{
			"github.com/VirtusLab/go-extended/pkg/errors.TestStackTrace\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:147",
		},
	}, {
		Wrapf(New("ooh2"), "ahh"), []string{
			"github.com/VirtusLab/go-extended/pkg/errors.TestStackTrace\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:152", // this is the stack of Wrap, not New
		},
	}, {
		Cause(Wrapf(New("ooh3"), "ahh")), []string{
			"github.com/VirtusLab/go-extended/pkg/errors.TestStackTrace\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:157", // this is the stack of New
		},
	}, {
		func() error { return New("ooh4") }(), []string{
			`github.com/VirtusLab/go-extended/pkg/errors.TestStackTrace.func1` +
				"\n\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:162", // this is the stack of New
			"github.com/VirtusLab/go-extended/pkg/errors.TestStackTrace\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:162", // this is the stack of New's caller
		},
	}, {
		Cause(func() error {
			return func() error {
				return Wrapf(New("hello %s"), fmt.Sprintf("world"))
			}()
		}()), []string{
			`github.com/VirtusLab/go-extended/pkg/errors.TestStackTrace.func2.1` +
				"\n\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:171", // this is the stack of Errorf
			`github.com/VirtusLab/go-extended/pkg/errors.TestStackTrace.func2` +
				"\n\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:172", // this is the stack of Errorf's caller
			"github.com/VirtusLab/go-extended/pkg/errors.TestStackTrace\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:173", // this is the stack of Errorf's caller's caller
		},
	}}
	for i, tt := range tests {
		x, ok := tt.err.(interface {
			StackTrace() StackTrace
		})
		if !ok {
			t.Errorf("expected %#v to implement StackTrace() StackTrace", tt.err)
			continue
		}
		st := x.StackTrace()
		for j, want := range tt.want {
			testFormatRegexp(t, i, st[j], "%+v", want)
		}
	}
}

func stackTrace() StackTrace {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	var st Stack = pcs[0:n]
	return st.StackTrace()
}

func TestStackTraceFormat(t *testing.T) {
	tests := []struct {
		StackTrace
		format string
		want   string
	}{{
		nil,
		"%s",
		`\[\]`,
	}, {
		nil,
		"%v",
		`\[\]`,
	}, {
		nil,
		"%+v",
		"",
	}, {
		nil,
		"%#v",
		`\[\]errors.Frame\(nil\)`,
	}, {
		make(StackTrace, 0),
		"%s",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%v",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%+v",
		"",
	}, {
		make(StackTrace, 0),
		"%#v",
		`\[\]errors.Frame{}`,
	}, {
		stackTrace()[:2],
		"%s",
		`\[stack_test.go stack_test.go\]`,
	}, {
		stackTrace()[:2],
		"%v",
		`\[stack_test.go:200 stack_test.go:247\]`,
	}, {
		stackTrace()[:2],
		"%+v",
		"\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.stackTrace\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:200\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.TestStackTraceFormat\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/stack_test.go:251",
	}, {
		stackTrace()[:2],
		"%#v",
		`\[\]errors.Frame{stack_test.go:200, stack_test.go:259}`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.StackTrace, tt.format, tt.want)
	}
}

// a version of runtime.Caller that returns a Frame, not a uintptr.
func caller() Frame {
	var pcs [3]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	return Frame(frame.PC)
}
