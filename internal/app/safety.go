package app

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"nextcmd/sdk"
)

type commandSafety struct {
	requiresConfirmation bool
	reason               string
}

func assessCommandSafety(command sdk.Command) commandSafety {
	if command.ShellCommand != "" {
		normalized := strings.NewReplacer("&&", ";", "||", ";", "|", ";", "&", ";").Replace(command.ShellCommand)
		for _, segment := range strings.Split(normalized, ";") {
			fields := strings.Fields(segment)
			if len(fields) == 0 {
				continue
			}
			if assessment := assessCommandTokens(fields[0], fields[1:]); assessment.requiresConfirmation {
				return assessment
			}
		}
		return commandSafety{}
	}
	return assessCommandTokens(command.Executable, command.Args)
}

func assessCommandTokens(executable string, args []string) commandSafety {
	name := commandName(executable)
	if name == "sudo" && len(args) > 0 {
		return assessCommandTokens(args[0], args[1:])
	}
	lowerArgs := make([]string, len(args))
	for index, arg := range args {
		lowerArgs[index] = strings.ToLower(strings.Trim(arg, `"'`))
	}
	switch name {
	case "rm":
		if hasShortOrLongFlag(lowerArgs, 'r', "--recursive") && hasShortOrLongFlag(lowerArgs, 'f', "--force") {
			return commandSafety{true, "recursive forced removal can permanently delete files"}
		}
	case "git":
		if len(lowerArgs) > 0 && lowerArgs[0] == "reset" && containsArg(lowerArgs, "--hard") {
			return commandSafety{true, "git reset --hard discards uncommitted changes"}
		}
		if len(lowerArgs) > 0 && lowerArgs[0] == "clean" && hasShortOrLongFlag(lowerArgs[1:], 'f', "--force") {
			return commandSafety{true, "git clean permanently removes untracked files"}
		}
		if len(lowerArgs) > 0 && lowerArgs[0] == "push" && (containsArg(lowerArgs, "--force") || containsArg(lowerArgs, "--force-with-lease") || containsArg(lowerArgs, "-f")) {
			return commandSafety{true, "force-push can overwrite remote history"}
		}
	case "docker":
		if containsArg(lowerArgs, "prune") {
			return commandSafety{true, "docker prune permanently removes unused resources"}
		}
	case "kubectl":
		if len(lowerArgs) > 0 && lowerArgs[0] == "delete" {
			return commandSafety{true, "kubectl delete removes resources from the selected cluster"}
		}
	case "terraform":
		if len(lowerArgs) > 0 && (lowerArgs[0] == "apply" || lowerArgs[0] == "destroy") {
			return commandSafety{true, "Terraform can materially change or destroy managed infrastructure"}
		}
	case "remove-item":
		if (containsArg(lowerArgs, "-recurse") || containsArg(lowerArgs, "-r")) && (containsArg(lowerArgs, "-force") || containsArg(lowerArgs, "-fo")) {
			return commandSafety{true, "recursive forced removal can permanently delete files"}
		}
	case "del", "erase", "rd", "rmdir":
		if containsArg(lowerArgs, "/s") && containsArg(lowerArgs, "/q") {
			return commandSafety{true, "quiet recursive removal can permanently delete files"}
		}
	case "powershell", "pwsh":
		joined := strings.Join(lowerArgs, " ")
		if strings.Contains(joined, "remove-item") && strings.Contains(joined, "-recurse") && strings.Contains(joined, "-force") {
			return commandSafety{true, "recursive forced removal can permanently delete files"}
		}
	}
	return commandSafety{}
}

func commandName(executable string) string {
	name := strings.ToLower(strings.Trim(filepath.Base(strings.Trim(executable, `"'`)), `"'`))
	return strings.TrimSuffix(name, ".exe")
}

func containsArg(args []string, wanted string) bool {
	for _, arg := range args {
		if arg == wanted {
			return true
		}
	}
	return false
}

func hasShortOrLongFlag(args []string, short byte, long string) bool {
	for _, arg := range args {
		if arg == long {
			return true
		}
		if len(arg) > 1 && arg[0] == '-' && arg[1] != '-' && strings.ContainsRune(arg[1:], rune(short)) {
			return true
		}
	}
	return false
}

func stripSafetyApproval(command sdk.Command) (sdk.Command, bool) {
	if command.ShellCommand != "" {
		trimmed := strings.TrimSpace(command.ShellCommand)
		if strings.HasSuffix(trimmed, " --yes") {
			command.ShellCommand = strings.TrimSpace(strings.TrimSuffix(trimmed, " --yes"))
			return command, true
		}
		return command, false
	}
	if len(command.Args) > 0 && command.Args[len(command.Args)-1] == "--yes" {
		command.Args = append([]string(nil), command.Args[:len(command.Args)-1]...)
		return command, true
	}
	return command, false
}

func confirmUnsafeCommand(reader io.Reader, writer io.Writer, command sdk.Command, safety commandSafety) (bool, error) {
	fmt.Fprintf(writer, "Dangerous command: %s\n", command.Display())
	fmt.Fprintf(writer, "Reason: %s\n", safety.reason)
	fmt.Fprint(writer, "Continue? [y/N]: ")
	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, err
		}
		return false, io.EOF
	}
	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}
