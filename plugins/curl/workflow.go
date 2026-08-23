package curl

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "curl") {
		return nil, nil
	}
	args := input.Result.Command.Args
	url := lastURL(args)
	if url == "" {
		url = "<url>"
	}
	out := []sdk.Suggestion{
		suggestion([]string{"--head", url}, "Inspect response headers", sdk.NextAction, sdk.Safe, 80, "Inspect metadata for the completed transfer"),
		suggestion([]string{"--write-out", "%{http_code} %{time_total}\\n", "--output", "<file>", url}, "Measure status and transfer time", sdk.NextAction, sdk.Mutating, 62, "Record request status and timing"),
	}
	if hasArg(args, "--request", "-X") && (hasArg(args, "POST") || hasArg(args, "PUT") || hasArg(args, "PATCH") || hasArg(args, "DELETE")) {
		out = append(out, suggestion([]string{"--verbose", url}, "Inspect the request exchange", sdk.NextAction, sdk.Safe, 58, "Review the mutating request if its result was unexpected"))
	}
	return out, nil
}

func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(input.Input)), "curl") {
		return nil, nil
	}
	return []sdk.Suggestion{
		suggestion([]string{"--fail-with-body", "--show-error", "<url>"}, "Preserve HTTP errors and a non-zero exit", sdk.BestPractice, sdk.Safe, 72, "Curl otherwise treats many HTTP error responses as successful transfers"),
		suggestion([]string{"--retry", "3", "--connect-timeout", "10", "--max-time", "60", "<url>"}, "Bound and retry transient transfers", sdk.BestPractice, sdk.Safe, 64, "Avoid unbounded waits and handle transient failures"),
	}, nil
}

func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "curl") {
		return nil, nil
	}
	message := strings.ToLower(input.Result.Stdout + input.Result.Stderr)
	url := lastURL(input.Result.Command.Args)
	if url == "" {
		url = "<url>"
	}
	state, _ := input.Project.(State)
	switch {
	case strings.Contains(message, "could not resolve host"):
		return []sdk.Suggestion{suggestion([]string{"--verbose", url}, "Inspect DNS and connection details", sdk.Recovery, sdk.Safe, 92, "Curl could not resolve the host name")}, nil
	case strings.Contains(message, "operation timed out"), strings.Contains(message, "timeout was reached"):
		return []sdk.Suggestion{suggestion([]string{"--retry", "3", "--connect-timeout", "10", "--max-time", "60", url}, "Retry with explicit time limits", sdk.Recovery, sdk.Safe, 94, "The transfer timed out")}, nil
	case strings.Contains(message, "ssl certificate problem"), strings.Contains(message, "certificate verify failed"):
		out := []sdk.Suggestion{}
		for _, certificate := range state.Certificates {
			out = append(out, suggestion([]string{"--cacert", certificate, url}, "Trust certificate "+certificate, sdk.Recovery, sdk.Safe, 98, "Use a certificate found in the working directory"))
		}
		return append(out, suggestion([]string{"--insecure", url}, "Retry without certificate verification", sdk.Recovery, sdk.Dangerous, 8, "Diagnostic fallback only; TLS identity would not be verified")), nil
	case strings.Contains(message, "failed to open/read local data"), strings.Contains(message, "error reading a file"):
		out := make([]sdk.Suggestion, 0, len(state.Files))
		for _, file := range state.Files {
			out = append(out, suggestion([]string{"--data-binary", "@" + file, url}, "Send data from "+file, sdk.Recovery, sdk.Mutating, 90, "Use a file found in the working directory"))
		}
		return out, nil
	case strings.Contains(message, "the requested url returned error"), input.Result.ExitCode == 22:
		return []sdk.Suggestion{suggestion([]string{"--show-headers", "--fail-with-body", url}, "Inspect the HTTP error response", sdk.Recovery, sdk.Safe, 92, "Curl reported an HTTP response error")}, nil
	}
	return nil, nil
}

func hasArg(args []string, values ...string) bool {
	for _, arg := range args {
		for _, value := range values {
			if strings.EqualFold(arg, value) {
				return true
			}
		}
	}
	return false
}

func lastURL(args []string) string {
	for i := len(args) - 1; i >= 0; i-- {
		lower := strings.ToLower(args[i])
		if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
			return args[i]
		}
	}
	return ""
}
