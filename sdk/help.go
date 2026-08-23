package sdk

// CommandHelp is UI-independent command documentation supplied by a plugin.
type CommandHelp struct {
	Command     Command
	Description string
	Risk        Risk
}

// HelpProvider exposes a plugin command catalog without coupling it to a UI.
type HelpProvider interface {
	Help() []CommandHelp
}
