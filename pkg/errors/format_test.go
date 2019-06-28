package errors

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type CustomError struct {
	stack *Stack
	s     string
}

func (e *CustomError) Error() string {
	return e.s
}

func (e *CustomError) Format(s fmt.State, verb rune) {
	FormatCauseAndStack(e, e.stack, s, verb)
}

func (e *CustomError) StackTrace() StackTrace {
	return e.stack.StackTrace()
}

func ErrCustom(s string) *CustomError {
	return &CustomError{
		stack: Callers(),
		s:     s,
	}
}

/*
 * Mostly taken from "github.com/pkg/errors"
 */

func TestNew(t *testing.T) {
	tests := []struct {
		err  string
		want error
	}{
		{"", fmt.Errorf("")},
		{"foo", fmt.Errorf("foo")},
		{"foo", New("foo")},
		{"string with format specifiers: %v", errors.New("string with format specifiers: %v")},
	}

	for _, tt := range tests {
		got := New(tt.err)
		if got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestWrapNil(t *testing.T) {
	got := Wrapf(nil, "no error")
	if got != nil {
		t.Errorf("Wrap(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrapf(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		got := Wrapf(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrap(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

func TestCause(t *testing.T) {
	x := New("error")
	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	}, {
		// explicit nil error is nil
		err:  (error)(nil),
		want: nil,
	}, {
		// typed nil is nil
		err:  (*nilError)(nil),
		want: (*nilError)(nil),
	}, {
		// uncaused error is unaffected
		err:  io.EOF,
		want: io.EOF,
	}, {
		// caused error returns cause
		err:  Wrapf(io.EOF, "ignored"),
		want: io.EOF,
	}, {
		err:  x, // return from errors.New
		want: x,
	}, {
		Wrapf(nil, "whoops"),
		nil,
	}, {
		Wrapf(io.EOF, "whoops"),
		io.EOF,
	}, {
		Wrap(nil),
		nil,
	}, {
		Wrap(io.EOF),
		io.EOF,
	}}

	for i, tt := range tests {
		got := Cause(tt.err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, tt.want)
		}
	}
}

func TestWrapfNil(t *testing.T) {
	got := Wrapf(nil, "no error")
	if got != nil {
		t.Errorf("Wrapf(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWrapf(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrapf(io.EOF, "read error without format specifiers"), "client error", "client error: read error without format specifiers: EOF"},
		{Wrapf(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
	}

	for _, tt := range tests {
		got := Wrapf(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrapf(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Errorf("read error without format specifiers"), "read error without format specifiers"},
		{Errorf("read error with %d format specifier", 1), "read error with 1 format specifier"},
	}

	for _, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("Errorf(%v): got: %q, want %q", tt.err, got, tt.want)
		}
	}
}

func TestWithStackNil(t *testing.T) {
	got := Wrap(nil)
	if got != nil {
		t.Errorf("Wrap(nil): got %#v, expected nil", got)
	}
}

func TestWithStack(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{io.EOF, "EOF"},
		{Wrap(io.EOF), "EOF"},
	}

	for _, tt := range tests {
		got := Wrap(tt.err).Error()
		if got != tt.want {
			t.Errorf("Wrap(%v): got: %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestWithMessageNil(t *testing.T) {
	got := Wrapf(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessage(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrapf(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		got := Wrapf(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestWithMessagefNil(t *testing.T) {
	got := Wrapf(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessagef(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrapf(io.EOF, "read error without format specifier"), "client error", "client error: read error without format specifier: EOF"},
		{Wrapf(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
	}

	for _, tt := range tests {
		got := Wrapf(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestErrorEquality(t *testing.T) {
	vals := []error{
		nil,
		io.EOF,
		errors.New("EOF"),
		New("EOF"),
		Errorf("EOF"),
		Wrapf(io.EOF, "EOF"),
		Wrapf(io.EOF, "EOF%d", 2),
		Wrapf(nil, "whoops"),
		Wrapf(io.EOF, "whoops"),
		Wrap(io.EOF),
		Wrap(nil),
	}

	for i := range vals {
		for j := range vals {
			_ = vals[i] == vals[j] // mustn't panic
		}
	}
}

func TestFormatNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		New("error"),
		"%s",
		"error",
	}, {
		New("error"),
		"%v",
		"error",
	}, {
		New("error"),
		"%+v",
		"error\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatNew\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:293",
	}, {
		New("error"),
		"%q",
		`"error"`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatErrorf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Errorf("%s", "error"),
		"%s",
		"error",
	}, {
		Errorf("%s", "error"),
		"%v",
		"error",
	}, {
		Errorf("%s", "error"),
		"%+v",
		"error\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatErrorf\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:323",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrap(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrapf(New("error"), "error2"),
		"%s",
		"error2: error",
	}, {
		Wrapf(New("error"), "error2"),
		"%v",
		"error2: error",
	}, {
		Wrapf(New("error"), "error2"),
		"%+v",
		"error\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWrap\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:349",
	}, {
		Wrapf(io.EOF, "error"),
		"%s",
		"error: EOF",
	}, {
		Wrapf(io.EOF, "error"),
		"%v",
		"error: EOF",
	}, {
		Wrapf(io.EOF, "error"),
		"%+v",
		"EOF\n" +
			"error\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWrap\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:363",
	}, {
		Wrapf(Wrapf(io.EOF, "error1"), "error2"),
		"%+v",
		"EOF\n" +
			"error1\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWrap\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:370\n",
	}, {
		Wrapf(New("error with space"), "context"),
		"%q",
		`"context: error with space"`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrapf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrapf(io.EOF, "error%d", 2),
		"%s",
		"error2: EOF",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%v",
		"error2: EOF",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%+v",
		"EOF\n" +
			"error2\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWrapf\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:401",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%s",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%v",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%+v",
		"error\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWrapf\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:416",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWithStack(t *testing.T) {
	tests := []struct {
		error
		format string
		want   []string
	}{{
		Wrap(io.EOF),
		"%s",
		[]string{"EOF"},
	}, {
		Wrap(io.EOF),
		"%v",
		[]string{"EOF"},
	}, {
		Wrap(io.EOF),
		"%+v",
		[]string{
			"EOF",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:442"},
	}, {
		Wrap(New("error")),
		"%s",
		[]string{"error"},
	}, {
		Wrap(New("error")),
		"%v",
		[]string{"error"},
	}, {
		Wrap(New("error")),
		"%+v",
		[]string{"error",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:457",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:457"},
	}, {
		Wrap(Wrap(io.EOF)),
		"%+v",
		[]string{"EOF",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:465",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:465"},
	}, {
		Wrap(Wrap(Wrapf(io.EOF, "message"))),
		"%+v",
		[]string{"EOF",
			"message",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:473",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:473",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:473"},
	}, {
		Wrap(Errorf("error%d", 1)),
		"%+v",
		[]string{"error1",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:484",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:484"},
	}, {
		ErrCustom("custom error"),
		"%+v",
		[]string{"custom error",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithStack\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:492"},
	}}

	for i, tt := range tests {
		testFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func TestFormatWithMessage(t *testing.T) {
	tests := []struct {
		error
		format string
		want   []string
	}{{
		Wrapf(New("error"), "error2"),
		"%s",
		[]string{"error2: error"},
	}, {
		Wrapf(New("error"), "error2"),
		"%v",
		[]string{"error2: error"},
	}, {
		Wrapf(New("error"), "error2"),
		"%+v",
		[]string{
			"error",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:518",
			"error2",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:518"},
	}, {
		Wrapf(io.EOF, "addition1"),
		"%s",
		[]string{"addition1: EOF"},
	}, {
		Wrapf(io.EOF, "addition1"),
		"%v",
		[]string{"addition1: EOF"},
	}, {
		Wrapf(io.EOF, "addition1"),
		"%+v",
		[]string{
			"EOF",
			"addition1",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:536"},
	}, {
		Wrapf(Wrapf(io.EOF, "addition1"), "addition2"),
		"%v",
		[]string{"addition2: addition1: EOF"},
	}, {
		Wrapf(Wrapf(io.EOF, "addition1"), "addition2"),
		"%+v",
		[]string{
			"EOF",
			"addition1",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:548",
			"addition2",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:548"},
	}, {
		Wrapf(Wrapf(io.EOF, "error1"), "error2"),
		"%+v",
		[]string{
			"EOF",
			"error1",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:559",
			"error2",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:559"},
	}, {
		Wrapf(Errorf("error%d", 1), "error2"),
		"%+v",
		[]string{
			"error1",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:570",
			"error2",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:570"},
	}, {
		Wrapf(Wrap(io.EOF), "error"),
		"%+v",
		[]string{
			"EOF",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:580",
			"error",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:580"},
	}, {
		Wrapf(Wrapf(Wrap(io.EOF), "inside-error"), "outside-error"),
		"%+v",
		[]string{
			"EOF",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:590",
			"inside-error",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:590",
			"outside-error",
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWithMessage\n" +
				"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:590"},
	}}

	for i, tt := range tests {
		testFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func wrappedNew(message string) error { // This function will be mid-stack inlined in go 1.12+
	return New(message)
}

func TestFormatWrappedNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		wrappedNew("error"),
		"%+v",
		"error\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.wrappedNew\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:610\n" +
			"github.com/VirtusLab/go-extended/pkg/errors.TestFormatWrappedNew\n" +
			"\t.+/github.com/VirtusLab/go-extended/pkg/errors/format_test.go:619",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func testFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	t.Helper()
	got := fmt.Sprintf(format, arg)
	gotLines := strings.SplitN(got, "\n", -1)
	wantLines := strings.SplitN(want, "\n", -1)

	if len(wantLines) > len(gotLines) {
		t.Errorf("test %d: wantLines(%d) > gotLines(%d):\n got: %q\nwant: %q", n+1, len(wantLines), len(gotLines), got, want)
		return
	}

	for i, w := range wantLines {
		match, err := regexp.MatchString(w, gotLines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("test %d: line %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, i+1, format, got, want)
		}
	}
}

var stackLineR = regexp.MustCompile(`\.`)

// parseBlocks parses input into a slice, where:
//  - incase entry contains a newline, its a stacktrace
//  - incase entry contains no newline, its a solo line.
//
// Detecting stack boundaries only works incase the Wrap-calls are
// to be found on the same line, thats why it is optionally here.
//
// Example use:
//
// for _, e := range blocks {
//   if strings.ContainsAny(e, "\n") {
//     // Match as stack
//   } else {
//     // Match as line
//   }
// }
//
func parseBlocks(input string, detectStackboundaries bool) ([]string, error) {
	var blocks []string

	stack := ""
	wasStack := false
	lines := map[string]bool{} // already found lines

	for _, l := range strings.Split(input, "\n") {
		isStackLine := stackLineR.MatchString(l)

		switch {
		case !isStackLine && wasStack:
			blocks = append(blocks, stack, l)
			stack = ""
			lines = map[string]bool{}
		case isStackLine:
			if wasStack {
				// Detecting two stacks after another, possible cause lines match in
				// our tests due to Wrap(Wrap(io.EOF)) on same line.
				if detectStackboundaries {
					if lines[l] {
						if len(stack) == 0 {
							return nil, New("len of block must not be zero here")
						}

						blocks = append(blocks, stack)
						stack = l
						lines = map[string]bool{l: true}
						continue
					}
				}

				stack = stack + "\n" + l
			} else {
				stack = l
			}
			lines[l] = true
		case !isStackLine && !wasStack:
			blocks = append(blocks, l)
		default:
			return nil, New("must not happen")
		}

		wasStack = isStackLine
	}

	// Use up stack
	if stack != "" {
		blocks = append(blocks, stack)
	}
	return blocks, nil
}

func testFormatCompleteCompare(t *testing.T, n int, arg interface{}, format string, want []string, detectStackBoundaries bool) {
	gotStr := fmt.Sprintf(format, arg)

	got, err := parseBlocks(gotStr, detectStackBoundaries)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != len(want) {
		t.Fatalf("test %d: fmt.Sprintf(%s, err) -> wrong number of blocks: got(%d) want(%d)\n got: %s\nwant: %s\ngotStr: %q",
			n+1, format, len(got), len(want), prettyBlocks(got), prettyBlocks(want), gotStr)
	}

	for i := range got {
		if strings.ContainsAny(want[i], "\n") {
			// Match as stack
			match, err := regexp.MatchString(want[i], got[i])
			if err != nil {
				t.Fatal(err)
			}
			if !match {
				t.Fatalf("test %d: block %d: fmt.Sprintf(%q, err):\ngot:\n%q\nwant:\n%q\nall-got:\n%s\nall-want:\n%s\n",
					n+1, i+1, format, got[i], want[i], prettyBlocks(got), prettyBlocks(want))
			}
		} else {
			// Match as message
			if got[i] != want[i] {
				t.Fatalf("test %d: fmt.Sprintf(%s, err) at block %d got != want:\n got: %q\nwant: %q", n+1, format, i+1, got[i], want[i])
			}
		}
	}
}

func prettyBlocks(blocks []string) string {
	var out []string

	for _, b := range blocks {
		out = append(out, fmt.Sprintf("%v", b))
	}

	return "   " + strings.Join(out, "\n   ")
}
