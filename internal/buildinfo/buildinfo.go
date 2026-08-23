package buildinfo

import (
	"runtime"
	"runtime/debug"
	"strings"
)

// Version is replaced at build time with -ldflags. Development builds keep
// the explicit "dev" value instead of reporting a misleading release number.
var Version = "dev"

type Info struct {
	Version, GoVersion, OS, Architecture, Revision string
}

func Current() Info {
	info := Info{Version: Version, GoVersion: runtime.Version(), OS: runtime.GOOS, Architecture: runtime.GOARCH}
	if build, ok := debug.ReadBuildInfo(); ok {
		if info.Version == "dev" && build.Main.Version != "" && build.Main.Version != "(devel)" {
			info.Version = build.Main.Version
		}
		for _, setting := range build.Settings {
			if setting.Key == "vcs.revision" {
				info.Revision = setting.Value
				if len(info.Revision) > 12 {
					info.Revision = info.Revision[:12]
				}
				break
			}
		}
	}
	info.Version = strings.TrimSpace(info.Version)
	return info
}
