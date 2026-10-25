package examples

import "os"

func ExampleImportPathVsPackageName() {
	ImportPathVsPackageName(os.Stdout)
	// Output:
	// rand.IntN: 76 61 78
	// pathpkg.Join: cmd/api/main.go
}

func ExampleBlankImportRegisters() {
	BlankImportRegisters(os.Stdout)
	// Output:
	// format=png width=3 height=2 (78 bytes embedded)
}

func ExampleEmbedDirective() {
	EmbedDirective(os.Stdout)
	// Output:
	// "Hello from an embedded file."
}

// BuildInfo depends on the machine and toolchain, so it has no // Output:
// block: it is compiled but not run by go test.
func ExampleBuildInfo() {
	BuildInfo(os.Stdout)
}
