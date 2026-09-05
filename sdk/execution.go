package sdk

import "time"

type ExecutionResult struct {
	Command        Command
	Stdout, Stderr string
	ExitCode       int
	Duration       time.Duration
	Canceled       bool
	Err            error
}
