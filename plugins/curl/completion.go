package curl

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

type commandSpec struct {
	args     []string
	title    string
	risk     sdk.Risk
	priority int
}

var commands = []commandSpec{
	{[]string{"--version"}, "Show the installed curl version and protocols", sdk.Safe, 35},
	{[]string{"<url>"}, "Send a GET request", sdk.Safe, 96},
	{[]string{"--head", "<url>"}, "Fetch response headers only", sdk.Safe, 86},
	{[]string{"--show-headers", "<url>"}, "Show response headers with the body", sdk.Safe, 76},
	{[]string{"--location", "<url>"}, "Follow HTTP redirects", sdk.Safe, 82},
	{[]string{"--fail-with-body", "--show-error", "<url>"}, "Fail on HTTP errors and keep the response body", sdk.Safe, 88},
	{[]string{"--silent", "--show-error", "<url>"}, "Hide progress while preserving errors", sdk.Safe, 72},
	{[]string{"--verbose", "<url>"}, "Inspect request, connection, and response details", sdk.Safe, 65},
	{[]string{"--request", "GET", "<url>"}, "Send an explicit GET request", sdk.Safe, 55},
	{[]string{"--request", "POST", "--data", "<data>", "<url>"}, "Send a POST request body", sdk.Mutating, 78},
	{[]string{"--request", "POST", "--header", "Content-Type: application/json", "--data", "<json>", "<url>"}, "Send a JSON POST request", sdk.Mutating, 88},
	{[]string{"--request", "PUT", "--data", "<data>", "<url>"}, "Replace a resource with PUT", sdk.Mutating, 62},
	{[]string{"--request", "PATCH", "--data", "<data>", "<url>"}, "Update a resource with PATCH", sdk.Mutating, 62},
	{[]string{"--request", "DELETE", "<url>"}, "Delete a remote resource", sdk.Destructive, 32},
	{[]string{"--header", "Accept: application/json", "<url>"}, "Request a JSON response", sdk.Safe, 70},
	{[]string{"--header", "<name>: <value>", "<url>"}, "Send a custom request header", sdk.Mutating, 58},
	{[]string{"--data-urlencode", "<name>=<value>", "<url>"}, "Send URL-encoded form data", sdk.Mutating, 60},
	{[]string{"--form", "<name>=@<file>", "<url>"}, "Upload multipart form data", sdk.Mutating, 55},
	{[]string{"--upload-file", "<file>", "<url>"}, "Upload a local file", sdk.Mutating, 55},
	{[]string{"--output", "<file>", "--remove-on-error", "<url>"}, "Download to a local file", sdk.Mutating, 66},
	{[]string{"--remote-name", "--remove-on-error", "<url>"}, "Download using the remote filename", sdk.Mutating, 58},
	{[]string{"--retry", "3", "--connect-timeout", "10", "--max-time", "60", "<url>"}, "Retry transient failures with time limits", sdk.Safe, 76},
	{[]string{"--write-out", "%{http_code} %{time_total}\\n", "--output", "<file>", "<url>"}, "Report status and total transfer time", sdk.Mutating, 52},
	{[]string{"--compressed", "<url>"}, "Request and decode a compressed response", sdk.Safe, 52},
	{[]string{"--config", "<file>", "<url>"}, "Load curl options from a file", sdk.Mutating, 44},
	{[]string{"--netrc-file", "<file>", "<url>"}, "Read credentials from a local netrc file", sdk.Mutating, 42},
	{[]string{"--cacert", "<certificate>", "<url>"}, "Use a trusted CA certificate file", sdk.Safe, 48},
	{[]string{"--insecure", "<url>"}, "Skip TLS certificate verification", sdk.Dangerous, 8},
}

func (*Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	trimmed := strings.TrimSpace(input.Input)
	if !matchesExecutable(trimmed, "curl") {
		return nil, nil
	}
	out := make([]sdk.Suggestion, 0, len(commands))
	for _, spec := range commands {
		out = append(out, suggestion(spec.args, spec.title, sdk.Completion, spec.risk, spec.priority, "Matches the current curl request"))
	}
	state, _ := input.Project.(State)
	return append(out, dynamic(input.Input, state)...), nil
}

func (*Plugin) Help() []sdk.CommandHelp {
	out := make([]sdk.CommandHelp, 0, len(commands))
	for _, spec := range commands {
		out = append(out, sdk.CommandHelp{Command: sdk.Command{Executable: "curl", Args: append([]string(nil), spec.args...)}, Description: spec.title, Risk: spec.risk})
	}
	return out
}

func matchesExecutable(input, executable string) bool {
	if input == "" {
		return true
	}
	first := strings.ToLower(strings.Fields(input)[0])
	return strings.HasPrefix(executable, first) || first == executable
}

func dynamic(input string, state State) []sdk.Suggestion {
	fields := strings.Fields(input)
	if len(fields) < 2 || !strings.EqualFold(fields[0], "curl") {
		return nil
	}
	last, previous := fields[len(fields)-1], ""
	if len(fields) > 2 {
		previous = fields[len(fields)-2]
	}
	out := []sdk.Suggestion{}
	switch {
	case flagValue(last, previous, "--upload-file", "-T"):
		for _, file := range state.Files {
			out = append(out, suggestion([]string{"--upload-file", file, "<url>"}, "Upload "+file, sdk.Completion, sdk.Mutating, 98, "File discovered in the working directory"))
		}
	case flagValue(last, previous, "--config", "-K"):
		for _, file := range state.ConfigFiles {
			out = append(out, suggestion([]string{"--config", file, "<url>"}, "Use curl config "+file, sdk.Completion, sdk.Mutating, 96, "Curl config discovered in the working directory"))
		}
	case flagValue(last, previous, "--cacert"):
		for _, file := range state.Certificates {
			out = append(out, suggestion([]string{"--cacert", file, "<url>"}, "Trust certificate "+file, sdk.Completion, sdk.Safe, 96, "Certificate discovered in the working directory"))
		}
	case flagValue(last, previous, "--data", "--data-binary", "-d"):
		for _, file := range state.Files {
			out = append(out, suggestion([]string{"--data-binary", "@" + file, "<url>"}, "Send data from "+file, sdk.Completion, sdk.Mutating, 94, "File discovered in the working directory"))
		}
	}
	return out
}

func flagValue(last, previous string, flags ...string) bool {
	for _, flag := range flags {
		if flagMatch(last, flag) || flagMatch(previous, flag) || (strings.HasPrefix(flag, "--") && strings.HasPrefix(strings.ToLower(flag), strings.ToLower(last))) {
			return true
		}
	}
	return false
}

func flagMatch(value, flag string) bool {
	if strings.HasPrefix(flag, "--") {
		return strings.EqualFold(value, flag)
	}
	return value == flag
}

func suggestion(args []string, title string, kind sdk.SuggestionKind, risk sdk.Risk, priority int, reason string) sdk.Suggestion {
	copied := append([]string(nil), args...)
	placeholders := []sdk.Placeholder{}
	for i, arg := range copied {
		for offset := 0; offset < len(arg); {
			start := strings.IndexByte(arg[offset:], '<')
			if start < 0 {
				break
			}
			start += offset
			end := strings.IndexByte(arg[start:], '>')
			if end < 0 {
				break
			}
			end += start
			placeholders = append(placeholders, sdk.Placeholder{Name: arg[start+1 : end], ArgIndex: i, Start: start, End: end + 1})
			offset = end + 1
		}
	}
	return sdk.Suggestion{Command: sdk.Command{Executable: "curl", Args: copied}, Title: title, Description: title, Reason: reason, Kind: kind, Priority: priority, Risk: risk, Source: "curl", Placeholders: placeholders}
}
