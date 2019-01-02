package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/VirtusLab/go-extended/pkg/log"
)

// Sh executes an operating system shell command
func Sh(ctx context.Context, logger log.StdLogger, env []string, stdin *string, prog string, args ...string) (stdout, stderr string, err error) {
	cmd := exec.CommandContext(ctx, prog, args...)

	if len(env) > 0 {
		cmd.Env = env
	}

	var stdinPipe io.WriteCloser
	if stdin != nil {
		stdinPipe, err = cmd.StdinPipe()
		if err != nil {
			err = fmt.Errorf("can't open stdin pipe: %+v", err)
			return
		}
		defer func() { _ = stdinPipe.Close() }() // just to be sure
	}

	// Set output to Byte Buffers
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	// Propagate POSIX signals
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signals
		err := cmd.Process.Signal(sig)
		if err != nil {
			if err.Error() != "os: process already finished" {
				logger.Printf("Failed to propagate POSIX signal '%s': %+v", sig, err)
			}
		}
	}()

	// Make sure the child process ended
	defer func() {
		err := cmd.Process.Kill()
		if err != nil {
			if err.Error() != "os: process already finished" {
				logger.Printf("Failed to kill the child process: %+v", err)
			}
		}
	}()

	if err = cmd.Start(); err != nil {
		stdout = outb.String()
		stderr = errb.String()
		err = fmt.Errorf("error executing command: %+v", err)
		return
	}

	if stdin != nil {
		if _, err = io.WriteString(stdinPipe, *stdin); err != nil {
			err = fmt.Errorf("error writing to stdin pipe: %+v", err)
			return
		}
		err = stdinPipe.Close() // must be called to flush the buffers
		if err != nil {
			err = fmt.Errorf("error closing the stdin pipe: %+v", err)
			return
		}
	}

	err = cmd.Wait()
	if err != nil {
		err = fmt.Errorf("error waiting for command: %+v", err) // must be called before stdout and stderr
	}
	stdout = outb.String()
	stderr = errb.String()

	return
}
