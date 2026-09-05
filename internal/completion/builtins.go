package completion

import "nextcmd/sdk"

func builtinSuggestions() []sdk.Suggestion {
	return []sdk.Suggestion{
		builtin(":?", "Show help", "Show keyboard controls, built-in commands, and loaded plugins", 100),
		builtin(":q", "Exit NextCmd", "Exit the application and restore the terminal", 95),
		builtin(":ls", "List directory", "List files and folders in the active working directory", 90),
		builtinWithArgument(":mkdir", "<path>", "Make directory", "Create a directory (and parents) in the active working directory", 88),
		builtinWithArgumentRisk(":del", "<path>", "Delete file or directory", "Move a file or directory to trash so it can be restored with :undo", 87, sdk.Destructive),
		builtinWithArgumentRisk(":trash", "<path>", "Move to trash", "Move a file or directory to NextCmd trash", 86, sdk.Mutating),
		builtin(":undo", "Undo delete", "Restore the last file or directory moved to trash in this session", 84),
		builtin(":plugins", "Show plugins", "List every enabled plugin with its version and description", 85),
		builtin(":history", "Show history", "Show the 20 most recent redacted commands", 80),
		builtin(":clear", "Clear screen", "Clear the terminal without changing application state", 75),
		builtin(":config", "Show configuration", "Show effective settings and configuration paths", 70),
		builtinWithArgument(":which", "<command>", "Locate executable", "Locate a command's executable in PATH", 65),
		builtin(":version", "Show version", "Show NextCmd, Go, platform, and build information", 60),
		builtinWithArgument(":cd", "<path>", "Change directory", "Change the active working directory", 55),
		builtin(":pwd", "Print directory", "Print the active working directory", 50),
	}
}

func builtinWithArgument(executable, argument, title, description string, priority int) sdk.Suggestion {
	return builtinWithArgumentRisk(executable, argument, title, description, priority, sdk.Safe)
}

func builtinWithArgumentRisk(executable, argument, title, description string, priority int, risk sdk.Risk) sdk.Suggestion {
	suggestion := builtin(executable, title, description, priority)
	suggestion.Risk = risk
	suggestion.Command.Args = []string{argument}
	suggestion.Placeholders = []sdk.Placeholder{{Name: argument[1 : len(argument)-1], ArgIndex: 0, Start: 0, End: len(argument)}}
	return suggestion
}

func builtin(executable, title, description string, priority int) sdk.Suggestion {
	return sdk.Suggestion{
		Command:     sdk.Command{Executable: executable},
		Title:       title,
		Description: description,
		Reason:      "Built-in NextCmd command",
		Kind:        sdk.Completion,
		Priority:    priority,
		Relevance:   100,
		Risk:        sdk.Safe,
		Source:      "nextcmd",
	}
}
