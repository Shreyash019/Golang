package examples

import (
	"bytes"
	_ "embed" // needed for //go:embed into a string or []byte
	"fmt"
	"image"
	_ "image/png" // imported only for its init(), which registers the PNG decoder
	"io"
	"strings"
)

//go:embed assets/pixel.png
var pixelPNG []byte

//go:embed assets/greeting.txt
var greeting string

// BlankImportRegisters decodes an embedded PNG through the generic image
// package. image knows nothing about PNG by itself: the blank import of
// image/png runs that package's init, which calls image.RegisterFormat.
func BlankImportRegisters(w io.Writer) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(pixelPNG))
	if err != nil {
		fmt.Fprintln(w, "error:", err)
		return
	}
	fmt.Fprintf(w, "format=%s width=%d height=%d (%d bytes embedded)\n", format, cfg.Width, cfg.Height, len(pixelPNG))
}

// EmbedDirective prints a text file that the compiler embedded into the
// binary because of the //go:embed directive (a comment the toolchain reads).
func EmbedDirective(w io.Writer) {
	fmt.Fprintf(w, "%q\n", strings.TrimSpace(greeting))
}
