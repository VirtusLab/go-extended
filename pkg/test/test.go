package test

import (
	"context"
	"flag"
	"fmt"
	stdlog "log"
	"os"
	"strings"
	"testing"

	"github.com/VirtusLab/go-extended/pkg/cli"
	"github.com/VirtusLab/go-extended/pkg/log"
)

// Context is a golang context used in tests
type Context interface {
	context.Context
	ID() string
	T() *testing.T
	Log() log.StdLogger
	Cleanup()
	AddCleanup(fn cleanupFunc)
}

type cleanupFunc func() error

type testContext struct {
	context.Context
	id         string
	cleanupFns []cleanupFunc
	t          *testing.T
	log        log.StdLogger
	cancel     context.CancelFunc
}

type testLogger struct {
	t *testing.T
}

// NewLogger creates a new test logger with standard interface
func NewLogger(t *testing.T) log.StdLogger {
	return &testLogger{t: t}
}

func (l testLogger) Fatal(args ...interface{}) {
	l.t.Fatal(args...)
}

func (l testLogger) Fatalf(format string, args ...interface{}) {
	l.t.Fatalf(format, args...)
}

func (l testLogger) Fatalln(args ...interface{}) {
	l.t.Fatal(fmt.Sprintln(args...))
}

func (l testLogger) Panic(args ...interface{}) {
	l.t.Fatal(args...)
	panic(fmt.Sprintln(args...))
}

func (l testLogger) Panicf(format string, args ...interface{}) {
	l.t.Fatalf(format, args...)
	panic(fmt.Sprintf(format, args...))
}

func (l testLogger) Panicln(args ...interface{}) {
	l.t.Fatal(fmt.Sprintln(args...))
	panic(fmt.Sprintln(args...))
}

func (l testLogger) Print(args ...interface{}) {
	l.t.Log(args...)
}

func (l testLogger) Printf(format string, args ...interface{}) {
	l.t.Logf(format, args...)
}

func (l testLogger) Println(args ...interface{}) {
	l.t.Log(fmt.Sprintln(args...))
}

// NewContext creates a new test context
func NewContext(parent context.Context, t *testing.T) Context {
	parent, cancel := context.WithCancel(parent)
	id := "<none>"
	logger := NewLogger(t)
	if t != nil {
		id = t.Name()
	} else {
		logger = stdlog.New(os.Stderr, "", stdlog.LstdFlags)
	}

	return &testContext{
		Context: parent,
		id:      id,
		t:       t,
		log:     logger,
		cancel:  cancel,
	}
}

// ID returns a test context identifier
func (ctx *testContext) ID() string {
	return ctx.id
}

// T returns a *testing.T for the test
func (ctx *testContext) T() *testing.T {
	return ctx.t
}

// Log returns log.StdLogger for the test
func (ctx *testContext) Log() log.StdLogger {
	return ctx.log
}

// AddCleanup adds a test cleanup function
func (ctx *testContext) AddCleanup(fn cleanupFunc) {
	ctx.cleanupFns = append(ctx.cleanupFns, fn)
}

// Cleanup must be used with defer to enable proper cleanup after tests
func (ctx *testContext) Cleanup() {
	failed := false
	for i := len(ctx.cleanupFns) - 1; i >= 0; i-- {
		err := ctx.cleanupFns[i]()
		if err != nil {
			failed = true
			if ctx.t != nil {
				ctx.t.Errorf("a cleanup function failed with error: %+v", err)
			} else {
				ctx.log.Fatalf("a cleanup function failed with error: %+v", err)
			}
		}
	}

	if failed {
		if ctx.t != nil {
			ctx.t.Fatalf("a cleanup function failed")
		} else {
			ctx.log.Fatal("a cleanup function failed")
		}
	}
	ctx.cancel()
}

// ProjectRootFlag is the name of the project root directory command line flag
const ProjectRootFlag = "root"

// Test is a typical test definition
type Test struct {
	Name    string
	Context Context
	Fn      func(tt Test)
	More    []Test
}

// Run runs all test case(s) with the standard test framework
func Run(t *testing.T, tests ...Test) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d %s", i, tt.Name), func(t *testing.T) {
			if tt.Fn != nil {
				tt.Fn(tt)
			}
			Run(t, tt.More...)
		})
	}
}

// Main is the main entry point form end-to-end tests
func Main(m *testing.M) {
	projectRoot := flag.String(ProjectRootFlag, "../..", "path to project root")
	flag.Parse()

	ctx := NewContext(context.TODO(), nil)
	ctx.Log().Print("Initializing e2e tests")

	// go test always runs from the test directory; change to project root
	err := os.Chdir(*projectRoot)
	if err != nil {
		ctx.Log().Fatalf("Failed to change directory to project root: %v", err)
	}

	_, _, err = Sh(ctx, nil, nil, "make", "install")
	if err != nil {
		ctx.Log().Fatalf("Can't build: %+v", err)
	}

	defer func() {
		exitCode := m.Run()
		ctx.Cleanup()
		os.Exit(exitCode)
	}()
}

// Sh runs a shell command for tests, see cli.Sh
func Sh(ctx Context, env []string, stdin *string, prog string, args ...string) (stdout, stderr string, err error) {
	allArgs := append([]string{prog}, args...)
	ctx.Log().Print(strings.Join(allArgs, " "))
	stdout, stderr, err = cli.Sh(ctx, ctx.Log(), env, stdin, prog, args...)
	ctx.Log().Printf("stdout: %s", stdout)
	ctx.Log().Printf("stderr: %s", stderr)
	return
}
