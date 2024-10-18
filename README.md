# Shell Command Executor with Timeout and Context Support

This Go project provides a utility to execute shell commands with arguments, directory specification, and a timeout using a context. It captures the standard output and error streams of the command and logs them in real-time. The project leverages context to handle command cancellation, making it useful for tasks that require precise timeout control.

## Features
- Execute shell commands with arguments in a specified directory.
- Stream standard output and error in real-time.
- Timeout handling with context-based cancellation.
- Custom buffer and split function for efficient output parsing.

## Usage

The main function demonstrates how to execute a simple bash command that prints output in a loop. The command is cancelled if it exceeds the specified timeout.

```go
func main() {
    ctx := context.Background()
    timeout := 5 * time.Second

    err := ShellExecWithArgs(ctx, "bash", []string{"-c", "for i in {1..5}; do echo \"output $i\"; sleep 3; done"}, ".", timeout)
    if err != nil {
        slog.Error(fmt.Sprintf("Error: %v", err))
    }
}
```

In this example:
- A bash script is executed with a timeout of 5 seconds.
- The script outputs five lines with a 3-second interval between them. If the command exceeds 5 seconds, it is cancelled.

## Functions
- ShellExecWithArgs: Executes a command with arguments, a specified directory, and a timeout.
- executeCommand: Handles the execution and cancellation of the command, as well as logging the output and errors.
- streamReader: Reads output from the command's stdout and stderr streams and logs it in real-time.
- splitFunc: Custom function to handle output splitting for better control over line endings.

## Requirements
Go 1.16 or later

## Installation
Clone the repository and run the example:
```go
git clone <repository_url>
cd <repository_directory>
go run main.go
```

## License
This project is licensed under the MIT License.