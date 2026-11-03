// Package examples holds the runnable worked examples from this topic's notes.
//
// Each example writes its output to an io.Writer, so the same code is run by
// cmd/examples and checked by the Example tests in examples_test.go: the
// output shown in the notes is the output the tests verify.
package examples

import "io"

// Example is one worked example from the notes.
type Example struct {
	// Name is the kebab-case name used on the command line, e.g. "slice-aliasing".
	Name string
	// Section is the notes section the example belongs to, e.g. "5.2".
	Section string
	// Summary is a one-line description shown when listing examples.
	Summary string
	// Crashes marks examples that panic, deadlock or race on purpose.
	// They run only when named explicitly, never as part of "all".
	Crashes bool
	// Run executes the example, writing its output to w.
	Run func(w io.Writer)
}

// All returns the examples in the order they appear in the notes.
func All() []Example {
	return []Example{
		{Name: "zero-values", Section: "5.1", Summary: "every kind of variable starts at its zero value", Run: ZeroValues},
		{Name: "declaration-forms", Section: "5.2", Summary: "var with a type, var with an initializer, :=, grouped var", Run: DeclarationForms},
		{Name: "default-types", Section: "5.2", Summary: "the type an untyped constant gets when nothing else decides", Run: DefaultTypes},
		{Name: "untyped-constants", Section: "5.3", Summary: "exact constant arithmetic; the type comes at the point of use", Run: UntypedConstants},
		{Name: "constant-float-exactness", Section: "5.3", Summary: "0.1+0.2 == 0.3 for constants but not for float64 variables", Run: ConstantFloatExactness},
		{Name: "typed-vs-untyped", Section: "5.4", Summary: "typed time.Duration constants mix with untyped ones, not with int variables", Run: TypedVsUntyped},
		{Name: "iota-basics", Section: "5.5", Summary: "iota is the line index in a const block; expressions repeat", Run: IotaBasics},
		{Name: "scope-levels", Section: "5.6", Summary: "the same name in nested blocks resolves inside-out", Run: ScopeLevels},
		{Name: "scope-starts-after-declaration", Section: "5.6", Summary: "in x := x + 10 the right side still sees the outer x", Run: ScopeStartsAfterDeclaration},
		{Name: "short-redeclaration", Section: "5.7", Summary: ":= reuses a variable in the same scope, shadows in an inner one", Run: ShortRedeclaration},
		{Name: "useful-zero-values", Section: "6.1", Summary: "strings.Builder, bytes.Buffer, sync.Mutex, nil slices work with no constructor", Run: UsefulZeroValues},
		{Name: "multiple-assignment", Section: "6.2", Summary: "swap, two-phase assignment order, the blank identifier", Run: MultipleAssignment},
		{Name: "iota-bit-flags", Section: "6.3", Summary: "1 << iota bit sets with |, & and &^", Run: IotaBitFlags},
		{Name: "iota-sizes", Section: "6.4", Summary: "KiB, MiB, GiB, TiB with iota, used as untyped constants", Run: IotaSizes},
		{Name: "shadowed-err", Section: "7.1", Summary: ":= inside an if hides the outer variable: the parsed value is lost", Run: ShadowedErr},
		{Name: "shadowed-package-var", Section: "7.2", Summary: ":= in a function declares a local instead of setting the package variable", Run: ShadowedPackageVar},
		{Name: "constant-int-division", Section: "7.7", Summary: "1/2 is integer division even when stored in a float64", Run: ConstantIntDivision},
		{Name: "iota-zero-value-trap", Section: "7.9", Summary: "an iota enum whose 0 is a real state makes unset values look valid", Run: IotaZeroValueTrap},
		{Name: "iota-reorder", Section: "7.10", Summary: "inserting a constant renumbers the rest and breaks stored values", Run: IotaReorder},
	}
}
