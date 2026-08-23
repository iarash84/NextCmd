package curl

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"nextcmd/sdk"
)

type State struct {
	Files        []string
	ConfigFiles  []string
	Certificates []string
}

func (*Plugin) Detect(ctx context.Context, input sdk.ProjectContext) (sdk.DetectionResult, error) {
	entries, err := os.ReadDir(input.WorkingDirectory)
	if err != nil {
		return sdk.DetectionResult{}, err
	}
	state := State{}
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return sdk.DetectionResult{}, ctx.Err()
		default:
		}
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		state.Files = append(state.Files, filepath.ToSlash(name))
		lower := strings.ToLower(name)
		extension := strings.ToLower(filepath.Ext(name))
		if lower == ".curlrc" || extension == ".curl" || extension == ".conf" {
			state.ConfigFiles = append(state.ConfigFiles, filepath.ToSlash(name))
		}
		if extension == ".pem" || extension == ".crt" || extension == ".cer" {
			state.Certificates = append(state.Certificates, filepath.ToSlash(name))
		}
	}
	sort.Strings(state.Files)
	sort.Strings(state.ConfigFiles)
	sort.Strings(state.Certificates)
	return sdk.DetectionResult{Detected: true, Project: state, CacheFor: 2 * time.Second}, nil
}
