package sdk

import "time"

type CompletionContext struct {
	Input, WorkingDirectory string
	Project                 any
	History                 []HistoryEntry
}
type CommandContext struct {
	Input, WorkingDirectory string
	Project                 any
	Previous                *ExecutionResult
}
type ExecutionContext struct {
	WorkingDirectory string
	Project          any
	Result           ExecutionResult
}
type ProjectContext struct{ WorkingDirectory string }
type DetectionResult struct {
	Detected bool
	Project  any
	CacheFor time.Duration
}
type HistoryEntry struct {
	Command          Command       `json:"command"`
	WorkingDirectory string        `json:"workingDirectory"`
	Timestamp        time.Time     `json:"timestamp"`
	ExitCode         int           `json:"exitCode"`
	Duration         time.Duration `json:"duration"`
	Plugin           string        `json:"plugin,omitempty"`
}
