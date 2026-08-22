package sdk

import "context"

type CompletionProvider interface {
	Complete(context.Context, CompletionContext) ([]Suggestion, error)
}
type NextActionProvider interface {
	NextActions(context.Context, ExecutionContext) ([]Suggestion, error)
}
type BestPracticeProvider interface {
	BestPractices(context.Context, CommandContext) ([]Suggestion, error)
}
type RecoveryProvider interface {
	Recover(context.Context, ExecutionContext) ([]Suggestion, error)
}
type ProjectDetector interface {
	Detect(context.Context, ProjectContext) (DetectionResult, error)
}
