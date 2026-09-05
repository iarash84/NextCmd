package sdk

import (
	"context"
	"io"
)

// Runner executes a structured command and returns its captured result.
// Implementations may invoke the operating system or provide deterministic
// results for tests and embedding applications.
type Runner interface {
	Run(context.Context, Command) ExecutionResult
}

// StreamingRunner is a Runner that can mirror output while retaining it in the
// returned result. Core uses this capability for interactive execution.
type StreamingRunner interface {
	Runner
	RunStreaming(context.Context, Command, io.Writer, io.Writer) ExecutionResult
}
