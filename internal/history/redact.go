package history

import (
	"strings"
	"unicode"

	"nextcmd/sdk"
)

const redacted = "<redacted>"

var sensitiveOptions = map[string]bool{
	"--api-key": true, "--apikey": true, "--access-token": true,
	"--auth-token": true, "--client-secret": true, "--password": true,
	"--passwd": true, "--private-key": true, "--secret": true, "--token": true,
}

var credentialOptions = map[string]bool{
	"-u": true, "--user": true, "--proxy-user": true,
}

var headerOptions = map[string]bool{"-h": true, "--header": true}

// Redact returns a detached command safe to persist in history. It handles both
// structured arguments and explicitly requested shell text without evaluating it.
func Redact(command sdk.Command) sdk.Command {
	out := command
	out.Args = redactArguments(command.Args)
	out.ShellCommand = redactShell(command.ShellCommand)
	if command.Environment != nil {
		out.Environment = make(map[string]string, len(command.Environment))
		for name, value := range command.Environment {
			if sensitiveName(name) {
				value = redacted
			}
			out.Environment[name] = value
		}
	}
	return out
}

func redactArguments(args []string) []string {
	out := append([]string(nil), args...)
	for index := range out {
		value := out[index]
		lower := strings.ToLower(value)
		if sensitiveOptions[lower] || credentialOptions[lower] {
			if index+1 < len(out) {
				out[index+1] = redacted
			}
			continue
		}
		if headerOptions[lower] && index+1 < len(out) && sensitiveHeader(out[index+1]) {
			out[index+1] = redacted
			continue
		}
		if name, _, ok := strings.Cut(value, "="); ok && (sensitiveOptions[strings.ToLower(name)] || credentialOptions[strings.ToLower(name)] || sensitiveName(name)) {
			out[index] = name + "=" + redacted
			continue
		}
		if name, header, ok := strings.Cut(value, "="); ok && headerOptions[strings.ToLower(name)] && sensitiveHeader(header) {
			out[index] = name + "=" + redacted
			continue
		}
		out[index] = redactURLCredentials(value)
	}
	return out
}

type shellWord struct {
	start int
	end   int
	value string
}

func redactShell(command string) string {
	if command == "" {
		return ""
	}
	words := shellWords(command)
	values := make([]string, len(words))
	for index, word := range words {
		values[index] = word.value
	}
	clean := redactArguments(values)
	for index := 0; index+2 < len(values); index++ {
		if sensitiveName(values[index]) && values[index+1] == "=" {
			clean[index+2] = redacted
		}
	}
	var result strings.Builder
	last := 0
	for index, word := range words {
		replacement := clean[index]
		if replacement == word.value {
			continue
		}
		result.WriteString(command[last:word.start])
		result.WriteString(replacement)
		last = word.end
	}
	if last == 0 {
		return command
	}
	result.WriteString(command[last:])
	return result.String()
}

// shellWords finds shell-like words while preserving their byte spans. Operators
// are boundaries; quotes and escapes only affect tokenization and are not executed.
func shellWords(input string) []shellWord {
	var words []shellWord
	for index := 0; index < len(input); {
		for index < len(input) && (unicode.IsSpace(rune(input[index])) || strings.ContainsRune("|&;()<>", rune(input[index]))) {
			index++
		}
		if index >= len(input) {
			break
		}
		start := index
		var value strings.Builder
		var quote byte
		for index < len(input) {
			current := input[index]
			if quote == 0 && (unicode.IsSpace(rune(current)) || strings.ContainsRune("|&;()<>", rune(current))) {
				break
			}
			if current == '\\' && quote != '\'' && index+1 < len(input) {
				index++
				value.WriteByte(input[index])
				index++
				continue
			}
			if current == '\'' || current == '"' {
				if quote == 0 {
					quote = current
					index++
					continue
				}
				if quote == current {
					quote = 0
					index++
					continue
				}
			}
			value.WriteByte(current)
			index++
		}
		words = append(words, shellWord{start: start, end: index, value: value.String()})
	}
	return words
}

func sensitiveName(name string) bool {
	normalized := strings.ToLower(strings.TrimSpace(name))
	normalized = strings.TrimPrefix(normalized, "$env:")
	normalized = strings.TrimPrefix(normalized, "$")
	normalized = strings.NewReplacer("-", "_", ".", "_").Replace(normalized)
	for _, part := range strings.FieldsFunc(normalized, func(r rune) bool { return r == '_' }) {
		if part == "password" || part == "passwd" || part == "token" || part == "secret" || part == "credential" || part == "apikey" {
			return true
		}
	}
	if strings.Contains(normalized, "api_key") || strings.Contains(normalized, "private_key") {
		return true
	}
	return false
}

func sensitiveHeader(header string) bool {
	name, _, ok := strings.Cut(header, ":")
	if !ok {
		return false
	}
	name = strings.ToLower(strings.TrimSpace(name))
	return name == "authorization" || name == "proxy-authorization" || name == "x-api-key" || name == "api-key"
}

func redactURLCredentials(value string) string {
	scheme := strings.Index(value, "://")
	if scheme < 0 {
		return value
	}
	authorityStart := scheme + 3
	rest := value[authorityStart:]
	authorityEnd := len(rest)
	if separator := strings.IndexAny(rest, "/?#"); separator >= 0 {
		authorityEnd = separator
	}
	at := strings.LastIndex(rest[:authorityEnd], "@")
	if at < 0 {
		return value
	}
	return value[:authorityStart] + redacted + rest[at:]
}
