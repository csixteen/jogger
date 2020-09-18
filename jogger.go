/*
Package jogger implements a simple library that allows you to execute a given
command with optional parameters and output the progress to stdout and stderr.
You also have the option of suppressing the output and change the timeout for
the command.

Example:

_, _, err := jogger.Run("ls", []string{"-la"})

// Supressing the output
_, _, err := jogger.Run("ls", []string{"-la"}, NoOutput())
*/
package jogger

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const defaultTimeout = time.Duration(30) * time.Second

type config struct {
	timeout  time.Duration
	noOutput bool
}

// Option is any option that the `Run` method may accept
type Option func(*config)

// NoOutput suppresses the output of the executed command
func NoOutput() Option {
	return func(c *config) {
		c.noOutput = true
	}
}

// Timeout changes the default timeout after which the process
// gets killed.
func Timeout(t time.Duration) Option {
	return func(c *config) {
		c.timeout = t
	}
}

// Run executes the command given with any extra parameters
// in a separate process. It may take extra options.
func Run(command string, args []string, opts ...Option) ([]byte, []byte, error) {
	cfg := &config{timeout: defaultTimeout}
	for _, o := range opts {
		o(cfg)
	}

	cmd := exec.Command(command, args...)

	done := make(chan error)

	// Checks for incoming signals
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGABRT)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		<-c
		done <- cmd.Process.Kill()
	}()

	// Progress capture
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var stdoutBuf, stderrBuf bytes.Buffer
	var stdout, stderr io.Writer

	if cfg.noOutput {
		stdout = io.MultiWriter(&stdoutBuf)
		stderr = io.MultiWriter(&stderrBuf)
	} else {
		stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	}

	// Starts the actual command
	err := cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	var errStdout, errStderr error
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(cfg.timeout):
		if err := cmd.Process.Kill(); err != nil {
			return nil, nil, err
		}

	case err := <-done:
		if err != nil {
			return nil, nil, err
		}

		if errStdout != nil {
			return nil, nil, errStdout
		}

		if errStderr != nil {
			return nil, nil, errStderr
		}
	}

	return stdoutBuf.Bytes(), stderrBuf.Bytes(), nil
}
