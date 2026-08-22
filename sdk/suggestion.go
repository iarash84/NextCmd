package sdk

type SuggestionKind string

const (
	Completion   SuggestionKind = "completion"
	NextAction   SuggestionKind = "next-action"
	BestPractice SuggestionKind = "best-practice"
	Recovery     SuggestionKind = "recovery"
)

type Risk string

const (
	Safe        Risk = "safe"
	Mutating    Risk = "mutating"
	Destructive Risk = "destructive"
	Dangerous   Risk = "dangerous"
)

type Placeholder struct {
	Name                 string
	ArgIndex, Start, End int
}
type Suggestion struct {
	Command                    Command
	Title, Description, Reason string
	Kind                       SuggestionKind
	Priority, Relevance        int
	Risk                       Risk
	Source                     string
	Placeholders               []Placeholder
}
