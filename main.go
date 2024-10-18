package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"time"
)

func ShellExecWithArgs(ctx context.Context, cmdName string, args []string, dir string, timeout time.Duration) error {
	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = dir

	slog.Info(fmt.Sprintf("cmd: %s %v, path: %s", cmdName, args, cmd.Dir))

	return executeCommand(ctx, cmd, timeout)
}

func executeCommand(ctx context.Context, cmd *exec.Cmd, timeout time.Duration) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	timer := time.AfterFunc(timeout, cancel)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("StdoutPipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("StderrPipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start: %w", err)
	}

	stdoutScanner := bufio.NewScanner(stdout)
	stdoutChan := make(chan string)
	stdoutDone := make(chan bool)
	stderrScanner := bufio.NewScanner(stderr)
	stderrChan := make(chan string)
	stderrDone := make(chan bool)

	go streamReader(stdoutScanner, stdoutChan, stdoutDone, timer, timeout)
	go streamReader(stderrScanner, stderrChan, stderrDone, timer, timeout)

	isRunning := true
	for isRunning {
		select {
		case <-stdoutDone:
			isRunning = false
		case line := <-stdoutChan:
			slog.Info(line)
		case line := <-stderrChan:
			slog.Error(line)
		case <-ctx.Done():
			return fmt.Errorf("command cancelled due to timeout or context cancellation")
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}

func streamReader(scanner *bufio.Scanner, outputChan chan string, doneChan chan bool, timer *time.Timer, timeout time.Duration) {
	buf := make([]byte, 4096)
	scanner.Buffer(buf, 65536)

	scanner.Split(splitFunc)
	defer close(outputChan)
	defer close(doneChan)
	for scanner.Scan() {
		outputChan <- scanner.Text()
		timer.Reset(timeout)
	}
	doneChan <- true
}

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	for i := 0; i < len(data); i++ {
		switch data[i] {
		case '\n':
			if i > 0 && data[i-1] == '\r' {
				return i + 1, data[:i-1], nil // CRLF
			}
			return i + 1, data[:i], nil // LF
		case '\r':
			if i == len(data)-1 || data[i+1] != '\n' {
				return i + 1, data[:i], nil // CR
			}
		}
	}

	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func main() {
	ctx := context.Background()
	timeout := 5 * time.Second

	err := ShellExecWithArgs(ctx, "bash", []string{"-c", "for i in {1..5}; do echo \"output $i\"; sleep 3; done"}, ".", timeout)
	// err := ShellExecWithArgs(ctx, "bash", []string{"-c", "for i in {1..3}; do echo \"error $i\" 1>&2; sleep 2; done"}, ".", timeout)
	// err := ShellExecWithArgs(ctx, "bash", []string{"-c", "echo \"start\"; sleep 10; echo \"end\""}, ".", timeout)
	if err != nil {
		slog.Error(fmt.Sprintf("Error: %v", err))
	}
}

// func main() {
// 	ctx := context.Background()
// 	cmdName := "bash"
// 	args := []string{"-c", "for i in {1..5}; do echo \"output $i\"; sleep 3; done"}
// 	cmd := exec.CommandContext(ctx, cmdName, args...)
// 	cmd.Dir = "."

// 	slog.Info(fmt.Sprintf("cmd: %s %v, path: %s", cmdName, args, cmd.Dir))
// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		slog.Error(fmt.Sprintf("Error: %v", err))
// 	}
// 	slog.Info(string(out))
// }
