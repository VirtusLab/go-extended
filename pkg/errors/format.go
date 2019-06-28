// Package errors is heavily inspired by "github.com/pkg/errors"
package errors

import (
	"fmt"
	"io"
)

// WithCause represents an error that was caused by another error
type WithCause interface {
	// Cause returns the error that caused this error
	Cause() error
}

// WithStackTrace represents an error with a stack trace
type WithStackTrace interface {
	// StackTrace returns a stack trace for this error
	StackTrace() StackTrace
}

type tracedError struct {
	cause   error
	stack   *Stack
	message string
}

// Wrapf wraps a given error with a formatted message and a stack trace.
// Please note the error type will change when wrapped.
// It records the stack trace at the point it was called.
func Wrapf(e error, format string, args ...interface{}) error {
	if e == nil {
		return nil
	}
	return &tracedError{
		cause:   e,
		message: fmt.Sprintf(format, args...),
		stack:   Callers(),
	}
}

// Wrap wraps a given error with a stack trace.
// Please note the error type will change when wrapped.
// It records the stack trace at the point it was called.
func Wrap(e error) error {
	if e == nil {
		return nil
	}
	return &tracedError{
		cause:   e,
		message: "",
		stack:   Callers(),
	}
}

// New returns a new error with the supplied message and a stack trace.
// It records the stack trace at the point it was called.
func New(message string) error {
	return &tracedError{
		cause:   nil,
		message: message,
		stack:   Callers(),
	}
}

// Errorf returns a new error with the supplied formatted message and stack trace.
// It records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return &tracedError{
		cause:   fmt.Errorf(format, args...),
		message: "",
		stack:   Callers(),
	}
}

func (t *tracedError) Error() string {
	hasMessage := len(t.message) > 0
	hasCause := t.cause != nil
	if hasMessage && hasCause {
		return t.message + ": " + t.cause.Error()
	} else if hasMessage && !hasCause {
		return t.message
	} else if !hasMessage && hasCause {
		return t.cause.Error()
	}
	return ""
}

func (t *tracedError) Cause() error {
	return t.cause
}

func (t *tracedError) StackTrace() StackTrace {
	return t.stack.StackTrace()
}

// Format implements fmt.Formatter used by Sprint(f) or Fprint(f) etc.
func (t *tracedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			hasMessage := len(t.message) > 0
			hasCause := t.cause != nil
			if hasMessage && hasCause {
				_, _ = fmt.Fprintf(s, "%+v", t.cause)
				_, _ = fmt.Fprintf(s, "\n%s", t.message)
			} else if hasMessage && !hasCause {
				_, _ = fmt.Fprintf(s, "%s", t.message)
			} else if !hasMessage && hasCause {
				_, _ = fmt.Fprintf(s, "%+v", t.cause)
			}
			t.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, t.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", t.Error())
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(WithCause)
		if !ok {
			break
		}
		if cause.Cause() == nil {
			break
		}
		err = cause.Cause()
	}
	return err
}

// FormatCauseAndStack helps to implement fmt.Formatter used by Sprint(f) or Fprint(f) etc.
// Use for custom error implementations with Cause and StackFormatter
func FormatCauseAndStack(err error, f *Stack, s fmt.State, verb rune) {
	_, _ = io.WriteString(s, err.Error())
	f.Format(s, verb)
}
