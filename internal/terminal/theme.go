package terminal

import (
	"fmt"
	"io"
	"os"
	"strings"

	"nextcmd/sdk"
)

const (
	ansiReset   = "\x1b[0m"
	ansiBold    = "\x1b[1m"
	ansiDim     = "\x1b[2m"
	ansiCyan    = "\x1b[36m"
	ansiGreen   = "\x1b[32m"
	ansiYellow  = "\x1b[33m"
	ansiRed     = "\x1b[31m"
	ansiMagenta = "\x1b[35m"
	ansiWhite   = "\x1b[97m"
)

func supportsColor(writer io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" || strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	file, ok := writer.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func paint(enabled bool, color, text string) string {
	if !enabled {
		return text
	}
	return color + text + ansiReset
}

func suggestionBadge(kind sdk.SuggestionKind) string {
	switch kind {
	case sdk.NextAction:
		return "NEXT"
	case sdk.BestPractice:
		return "TIP"
	case sdk.Recovery:
		return "FIX"
	default:
		return "COMP"
	}
}

func kindColor(kind sdk.SuggestionKind) string {
	switch kind {
	case sdk.NextAction:
		return ansiCyan
	case sdk.BestPractice:
		return ansiMagenta
	case sdk.Recovery:
		return ansiYellow
	default:
		return ansiDim
	}
}

func riskColor(risk sdk.Risk) string {
	switch risk {
	case sdk.Mutating:
		return ansiYellow
	case sdk.Destructive:
		return ansiRed
	case sdk.Dangerous:
		return ansiMagenta
	default:
		return ansiGreen
	}
}

// PrintWelcome renders the compact application banner without affecting plugins.
func PrintWelcome(writer io.Writer) {
	color := supportsColor(writer)
	fmt.Fprintf(writer, "%s %s\n", paint(color, ansiBold+ansiCyan, "❯ NextCmd"), paint(color, ansiDim, "programming command assistant"))
	fmt.Fprintf(writer, "  %s  %s  %s\n\n",
		paint(color, ansiDim, "↑/↓ select · Tab/→ accept · Enter run"),
		paint(color, ansiCyan, ":? help"),
		paint(color, ansiDim, ":q exit"))
}

// PrintExecutionSummary makes command completion visible without changing captured output.
func PrintExecutionSummary(writer io.Writer, result sdk.ExecutionResult) {
	color := supportsColor(writer)
	label, tone := "done", ansiGreen
	if result.Canceled {
		label, tone = "canceled", ansiYellow
	} else if result.ExitCode != 0 || result.Err != nil {
		label, tone = "failed", ansiRed
	}
	fmt.Fprintf(writer, "%s %s\n", paint(color, tone+ansiBold, "● "+label), paint(color, ansiDim, fmt.Sprintf("exit %d · %s", result.ExitCode, result.Duration.Round(1e6))))
}
