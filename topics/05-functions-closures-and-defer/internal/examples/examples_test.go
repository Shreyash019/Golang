package examples

import (
	"regexp"
	"testing"
)

var kebab = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func TestAllIsWellFormed(t *testing.T) {
	seen := make(map[string]bool)
	for _, ex := range All() {
		if !kebab.MatchString(ex.Name) {
			t.Errorf("example name %q is not kebab-case", ex.Name)
		}
		if seen[ex.Name] {
			t.Errorf("duplicate example name %q", ex.Name)
		}
		seen[ex.Name] = true
		if ex.Section == "" || ex.Summary == "" || ex.Run == nil {
			t.Errorf("example %q needs a Section, a Summary and a Run func", ex.Name)
		}
	}
}

// Each deterministic example gets an Example test whose // Output: block is
// the exact output pasted into the notes, e.g.
//
//	func ExampleSliceAliasing() {
//		SliceAliasing(os.Stdout)
//		// Output:
//		// [1 2 3]
//	}
//
// Use "// Unordered output:" for map iteration, and no output check at all
// for examples whose output depends on scheduling.
