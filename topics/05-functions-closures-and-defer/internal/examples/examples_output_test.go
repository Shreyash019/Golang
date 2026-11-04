package examples

import "os"

func ExamplePassByValue() {
	PassByValue(os.Stdout)
	// Output:
	// n:   1     (the int was copied)
	// x:   42    (written through the pointer)
	// arr: [1 2] (the whole array was copied)
	// sl:  [100 2] (element write shared; the append was not)
	// m:   map[a:1 b:2] (same map)
}

func ExampleMultipleReturns() {
	MultipleReturns(os.Stdout)
	// Output:
	// divmod(17, 5) = 3 2
	// remainder only: -2
	// Atoi failed: strconv.Atoi: parsing "12a": invalid syntax | n is the zero value: 0
	// sum of divmod: 5
}

func ExampleNamedResults() {
	NamedResults(os.Stdout)
	// Output:
	// minMax: -2 9
	// minMax(nil): 0 0 (zero values)
}

func ExampleVariadic() {
	Variadic(os.Stdout)
	// Output:
	// 0 1 6
	// 60
	// ids: len=0 nil=true ids: len=2 nil=false
}

func ExampleVariadicAliasing() {
	VariadicAliasing(os.Stdout)
	// Output:
	// after double(scores...): [2 4 6]
	// after double(4, 5, 6): nothing visible to change
}

func ExampleFunctionValues() {
	FunctionValues(os.Stdout)
	// Output:
	// zero value is nil: true
	// op = add: 5
	// op = multiply: 6
	// pipeline: "HELLO!"
	// type of add: func(int, int) int
}

func ExampleHigherOrder() {
	HigherOrder(os.Stdout)
	// Output:
	// filter even: [2 4 6]
	// filter > 3:  [4 5 6]
	// compose(sq, inc)(3) = 16 (inc first, then sq)
	// compose(inc, sq)(3) = 10
}

func ExampleClosureCounter() {
	ClosureCounter(os.Stdout)
	// Output:
	// 1 2 3 | a has its own count
	// 1 2 | b has another
	// 4 | a continues
}

func ExampleClosureSharesVariable() {
	ClosureSharesVariable(os.Stdout)
	// Output:
	// goodbye, world
	// read(): 15 | total: 15
	// snapshot: 1 | v: 2
}

func ExampleClosureState() {
	ClosureState(os.Stdout)
	// Output:
	// order-1 order-2 order-3
	// order-1
}

func ExampleRecursiveClosure() {
	RecursiveClosure(os.Stdout)
	// Output:
	// fib(10) = 55
}

func ExampleRecursion() {
	Recursion(os.Stdout)
	// Output:
	// factorial(10) = 3628800
	// gcd(84, 36) = 12
	// /
	//   etc
	//     hosts
	//   usr
	//     bin
	//       go
	//     lib
	// files: 7
}

func ExampleDeferOrder() {
	DeferOrder(os.Stdout)
	// Output:
	// start
	// end of function body
	// deferred 3
	// deferred 2
	// deferred 1
}

func ExampleDeferArgsEvaluated() {
	DeferArgsEvaluated(os.Stdout)
	// Output:
	// x set to 2
	// deferred closure, reads x when it runs:     x = 2
	// deferred call, argument evaluated at defer time: x = 1
}

func ExampleDeferModifiesResult() {
	DeferModifiesResult(os.Stdout)
	// Output:
	// triple(5) = 15
	// loadUser(""): load user "": empty name
	// loadUser("gopher"): <nil>
}

func ExampleDeferInLoop() {
	DeferInLoop(os.Stdout)
	// Output:
	// -- defer directly in the loop:
	// open a (open now: 1)
	// open b (open now: 2)
	// open c (open now: 3)
	// loop finished with 3 open
	// close c
	// close b
	// close a
	// -- body moved into a function:
	// open a (open now: 1)
	// close a
	// open b (open now: 1)
	// close b
	// open c (open now: 1)
	// close c
	// loop finished with 0 open
}

func ExampleDeferTrace() {
	DeferTrace(os.Stdout)
	// Output:
	// enter outer
	//   enter inner
	//     (inner's work)
	//   leave inner
	//   enter inner
	//     (inner's work)
	//   leave inner
	// leave outer
}

func ExampleErrorFirstLook() {
	ErrorFirstLook(os.Stdout)
	// Output:
	// "3"  ok: 3
	// "x"  error: parse quantity "x": strconv.Atoi: parsing "x": invalid syntax (negative? false)
	// "-2" error: parse quantity "-2": quantity is negative (negative? true)
}
