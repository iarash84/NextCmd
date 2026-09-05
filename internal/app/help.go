package app

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"nextcmd/internal/completion"
)

func parseHelpCommand(input string) (plugin string, help bool) {
	fields := strings.Fields(strings.TrimSpace(input))
	if len(fields) == 0 || (fields[0] != ":?" && fields[0] != ":؟") {
		return "", false
	}
	if len(fields) > 1 {
		return fields[1], true
	}
	return "", true
}

func printHelp(writer io.Writer, engine *completion.Engine, pluginName string) {
	if pluginName == "" {
		printGeneralHelp(writer, engine)
		return
	}
	if printBuiltinHelp(writer, pluginName) {
		return
	}
	info, commands, found := engine.Help(pluginName)
	if !found {
		fmt.Fprintf(writer, "Unknown plugin %q. Use :? to list loaded plugins.\n", pluginName)
		return
	}
	fmt.Fprintf(writer, "%s plugin commands (%s):\n", info.Name, info.ID)
	if len(commands) == 0 {
		fmt.Fprintln(writer, "  This plugin does not expose a command catalog.")
		return
	}
	sort.SliceStable(commands, func(i, j int) bool { return commands[i].Command.Display() < commands[j].Command.Display() })
	for _, command := range commands {
		fmt.Fprintf(writer, "  %-52s [%s]\n      %s\n", command.Command.Display(), command.Risk, command.Description)
	}
}

func printGeneralHelp(writer io.Writer, engine *completion.Engine) {
	fmt.Fprintln(writer, "NextCmd help")
	fmt.Fprintln(writer, "Type : in the command editor to show all built-in commands.")
	fmt.Fprintln(writer, "Keyboard:")
	fmt.Fprintln(writer, "  Up/Down       Highlight suggestions")
	fmt.Fprintln(writer, "  Tab/Right     Accept the highlighted suggestion")
	fmt.Fprintln(writer, "  Left/Right    Move the caret inside the command editor")
	fmt.Fprintln(writer, "  Ctrl+P/Ctrl+N Recall the previous or next executed command")
	fmt.Fprintln(writer, "  Ctrl+A/Ctrl+E Move to the start or end of the line")
	fmt.Fprintln(writer, "  Ctrl+U        Clear the current command line")
	fmt.Fprintln(writer, "  Enter         Accept a suggestion, then execute on the next press")
	fmt.Fprintln(writer, "  Backspace     Delete the previous character")
	fmt.Fprintln(writer, "  Escape        Clear the current command line")
	fmt.Fprintln(writer, "  Ctrl+C/Ctrl+D Exit and clean the terminal UI")
	fmt.Fprintln(writer, "\nShell execution:")
	fmt.Fprintln(writer, "  ! <command>   Run through cmd.exe on Windows or /bin/sh elsewhere")
	fmt.Fprintln(writer, "                Supports shell built-ins, pipes, redirects, and variables")
	fmt.Fprintln(writer, "Commands:")
	fmt.Fprintln(writer, "  :? or :؟             Show this help")
	fmt.Fprintln(writer, "  :? <plugin>           Show every command exposed by a plugin")
	fmt.Fprintln(writer, "  cd/:cd <path>         Change the working directory")
	fmt.Fprintln(writer, "  pwd/:pwd              Print the current working directory")
	fmt.Fprintln(writer, "  :ls [path]            List files and directories")
	fmt.Fprintln(writer, "  :mkdir <path>         Create a directory (parents included)")
	fmt.Fprintln(writer, "  :del <path>           Move a file or directory to trash")
	fmt.Fprintln(writer, "  :trash <path>         Move a file or directory to trash")
	fmt.Fprintln(writer, "  :undo                 Restore the last trashed item")
	fmt.Fprintln(writer, "  :history [count]      Show recent redacted command history")
	fmt.Fprintln(writer, "  :plugins              Show enabled plugins")
	fmt.Fprintln(writer, "  :clear                Clear the terminal screen")
	fmt.Fprintln(writer, "  :config               Show effective configuration")
	fmt.Fprintln(writer, "  :which <command>      Locate an executable in PATH")
	fmt.Fprintln(writer, "  :version              Show version and build information")
	fmt.Fprintln(writer, "  exit, quit, or :q     Exit NextCmd")
	fmt.Fprintln(writer, "Loaded plugins:")
	for _, plugin := range engine.Plugins() {
		fmt.Fprintf(writer, "  %-12s %s — use :? %s\n", plugin.Info().ID, plugin.Info().Description, plugin.Info().ID)
	}
}
