package errors_test

import (
	"fmt"

	"github.com/VirtusLab/go-extended/pkg/errors"
)

type ErrCustom struct {
	stack *errors.Stack
	s     string
}

func (e *ErrCustom) Error() string {
	return e.s
}

func (e *ErrCustom) Format(s fmt.State, verb rune) {
	errors.FormatCauseAndStack(e, e.stack, s, verb)
}

func (e *ErrCustom) StackTrace() errors.StackTrace {
	return e.stack.StackTrace()
}

func NewErrCustom(s string) *ErrCustom {
	return &ErrCustom{
		stack: errors.Callers(),
		s:     s,
	}
}

func ExampleFormatCauseAndStack() {
	err := NewErrCustom("whoops")

	fmt.Println(err)

	// Output: whoops
}

func ExampleFormatCauseAndStack_printf() {
	err := NewErrCustom("whoops")

	fmt.Printf("%+v", err)

	// Example output:
	// whoops
	// github.com/VirtusLab/go-extended/pkg/errors_test.ExampleFormatCauseAndStack_prints
	//         /home/dfc/go/src/github.com/VirtusLab/go-extended/pkg/errors/custom_test.go:41
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:121
	// testing.runExamples
	//         /home/dfc/go/src/testing/example.go:45
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:1073
	// main.main
	//        _testmain.go:108
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:200
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:1337
}
