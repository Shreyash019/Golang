package main

// prediction is the real output, line by line:
//
//   - 0 5 15: iota is the line index (0, 1, 2, 3) and the omitted
//     expressions repeat iota * 5; the third line is discarded by _.
//   - if: 20: in `level := level + 1` the right side still sees the package
//     level (1), because a local's scope starts after its declaration; the
//     if's own level is 2 * 10.
//   - function: 2 and package: 1: the if's level and the function's level
//     are separate variables, and nothing ever assigns the package one.
//   - 12 0 true: `m, err :=` declares m and reuses err (same scope), which
//     now holds Atoi's error for "x"; m is 0.
//   - 3 3.5 3.5: 7/2 divides two untyped integer constants; a float operand
//     makes the whole constant expression a float.
//   - int32 int: the default types of an untyped rune and an untyped integer.
//   - true false: constant arithmetic is exact, float64 arithmetic isn't.
const prediction = `0 5 15
if: 20
function: 2
package: 1
12 0 true
3 3.5 3.5
int32 int
true false
`
