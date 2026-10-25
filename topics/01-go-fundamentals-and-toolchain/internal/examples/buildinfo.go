package examples

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
)

// BuildInfo prints what the toolchain recorded about the running binary:
// the Go version it was built with, the target platform, and the main
// module. `go version -m <binary>` shows the same data from outside.
func BuildInfo(w io.Writer) {
	fmt.Fprintln(w, "go version:", runtime.Version())
	fmt.Fprintln(w, "platform:  ", runtime.GOOS+"/"+runtime.GOARCH)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprintln(w, "no build info (binary built without module support)")
		return
	}
	fmt.Fprintln(w, "main module:", info.Main.Path, info.Main.Version)
	for _, s := range info.Settings {
		switch s.Key {
		case "-trimpath", "CGO_ENABLED", "GOOS", "GOARCH", "vcs", "vcs.modified":
			fmt.Fprintf(w, "  %s=%s\n", s.Key, s.Value)
		}
	}
}
