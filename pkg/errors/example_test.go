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
package errors_test

import (
	"fmt"

	"github.com/VirtusLab/go-extended/pkg/errors"
)

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err)

	// Output: whoops
}

func ExampleNew_printf() {
	err := errors.New("whoops")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops
	// github.com/VirtusLab/go-extended/pkg/errors_test.ExampleNew_printf
	//         /home/dfc/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:17
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/VirtusLab/go-extended/pkg/errors/_test/_testmain.go:106
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
}

func ExampleWrapf() {
	cause := errors.New("whoops")
	err := errors.Wrapf(cause, "oh noes")
	fmt.Println(err)

	cause = errors.New("whoops")
	err = errors.Wrapf(cause, "oh noes #%d", 2)
	fmt.Println(err)

	// Example output:
	// oh noes: whoops
	// oh noes #2: whoops
}

func ExampleWrap() {
	cause := errors.New("whoops")
	err := errors.Wrap(cause)
	fmt.Println(err)

	// Output: whoops
}

func ExampleWrap_printf() {
	cause := errors.New("whoops")
	err := errors.Wrap(cause)
	fmt.Printf("%+v", err)

	// Example Output:
	// whoops
	// github.com/VirtusLab/go-extended/pkg/errors_test.ExampleWithStack_printf
	//         /home/fabstu/go/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:55
	// testing.runExample
	//         /usr/lib/go/src/testing/example.go:114
	// testing.RunExamples
	//         /usr/lib/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /usr/lib/go/src/testing/testing.go:744
	// main.main
	//         github.com/VirtusLab/go-extended/pkg/errors/_test/_testmain.go:106
	// runtime.main
	//         /usr/lib/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /usr/lib/go/src/runtime/asm_amd64.s:2086
	// github.com/VirtusLab/go-extended/pkg/errors_test.ExampleWithStack_printf
	//         /home/fabstu/go/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:56
	// testing.runExample
	//         /usr/lib/go/src/testing/example.go:114
	// testing.RunExamples
	//         /usr/lib/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /usr/lib/go/src/testing/testing.go:744
	// main.main
	//         github.com/VirtusLab/go-extended/pkg/errors/_test/_testmain.go:106
	// runtime.main
	//         /usr/lib/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /usr/lib/go/src/runtime/asm_amd64.s:2086
}

func fn() error {
	e1 := errors.New("error")
	e2 := errors.Wrapf(e1, "inner")
	e3 := errors.Wrapf(e2, "middle")
	return errors.Wrapf(e3, "outer")
}

func ExampleCause() {
	err := fn()
	fmt.Println(err)
	fmt.Println(errors.Cause(err))

	// Output: outer: middle: inner: error
	// error
}

func ExampleWrap_extended() {
	err := fn()
	fmt.Printf("%+v\n", err)

	// Example output:
	// error
	// github.com/VirtusLab/go-extended/pkg/errors_test.fn
	//         /home/dfc/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:47
	// github.com/VirtusLab/go-extended/pkg/errors_test.ExampleCause_printf
	//         /home/dfc/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:63
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/VirtusLab/go-extended/pkg/errors/_test/_testmain.go:104
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
	// github.com/VirtusLab/go-extended/pkg/errors_test.fn
	// 	  /home/dfc/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:48: inner
	// github.com/VirtusLab/go-extended/pkg/errors_test.fn
	//        /home/dfc/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:49: middle
	// github.com/VirtusLab/go-extended/pkg/errors_test.fn
	//      /home/dfc/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:50: outer
}

func ExampleErrorf_extended() {
	err := errors.Errorf("whoops: %s", "foo")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops: foo
	// github.com/VirtusLab/go-extended/pkg/errors_test.ExampleErrorf
	//         /home/dfc/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:101
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/VirtusLab/go-extended/pkg/errors/_test/_testmain.go:102
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
}

func Example_stackTrace() {
	err, ok := errors.Cause(fn()).(errors.WithStackTrace)
	if !ok {
		panic("oops, err does not implement stackTracer")
	}

	st := err.StackTrace()
	fmt.Printf("%+v", st[0:2]) // top two frames

	// Example output:
	// github.com/VirtusLab/go-extended/pkg/errors_test.fn
	//	/home/dfc/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:47
	// github.com/VirtusLab/go-extended/pkg/errors_test.Example_stackTrace
	//	/home/dfc/src/github.com/VirtusLab/go-extended/pkg/errors/example_test.go:127
}

func ExampleCause_printf() {
	err := errors.Wrapf(func() error {
		return func() error {
			//lint:ignore S1039 the test introduces it on purpose
			return errors.Errorf("hello %s", fmt.Sprintf("world"))
		}()
	}(), "failed")

	fmt.Printf("%v", err)

	// Output: failed: hello world
}
