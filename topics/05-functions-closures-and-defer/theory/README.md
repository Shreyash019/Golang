# 05 · Functions, Closures & Defer

> **Level:** Foundation
> **Prerequisites:** [02 Variables & constants](../../02-variables-and-constants/theory/README.md) (scope, zero values) · [04 Control flow](../../04-control-flow/theory/README.md) (loops, per-iteration loop variables §5.3, `defer` in a loop §7.10)
> **Related:** 06 Arrays & slices (what a slice argument shares) · 09 Pointers & value semantics · 10 Structs & methods (methods, method values) · 14 Generics (generic functions) · 15 Iterators (functions as iterators) · 18 Errors, panic & recover (errors in depth, `recover`) · 24 sync (`defer mu.Unlock()`) · 28 Memory management (escape analysis in depth) · 35 Functional patterns (functional options, middleware)
> **Go version:** go1.27.1 (darwin/arm64) · everything here applies from Go 1.22 on; open-coded defers since 1.14, the register-based calling convention since 1.17 (amd64) and 1.18 (arm64)
> **Targets stuck patterns:** none recorded yet (no diagnostic). Likely targets: closures and `defer` "behaving randomly" (hypothesis 3, missing mental models: what is captured, when arguments are evaluated), and Java/Python habits: exceptions as control flow, one-method interfaces and factories where a function would do, out-parameters (hypothesis 4)
> **Started:** 2026-11-04 · **Status:** 📖 In progress
>
> **Syllabus items covered:** 6 (functions: declaration, parameters, multiple & named returns, variadic functions) and 7 (function concepts: anonymous functions, first-class & higher-order functions, closures, recursion, `defer`). Plus a first look at returning `error` (home: 18). `recover` is only mentioned (home: 18). Methods are topic 10, generic functions topic 14.
> **Run the code:** `cd topics/05-functions-closures-and-defer && make examples` · exercises in [`exercises/`](../exercises/README.md)

---

## 1. TL;DR

- **Everything is passed by value.** Each argument is copied into a fresh parameter. A slice, map, channel, func or pointer *value* is small and refers to shared data, so the callee can change that data but never the caller's variable. There are no default parameters, no keyword arguments and no overloading.
- **Functions return several values.** The idiom is `(result, error)` with `error` last, checked right away with `if err != nil`. Named results document the meaning, start at their zero values, and can be changed by a deferred closure.
- **Functions are values.** They have types (`func(int) int`), a zero value (`nil`, which panics when called), and can be stored, passed and returned. A **closure** captures *variables*, not values: it shares them with the enclosing function and keeps them alive, usually on the heap.
- **`defer f(args)` evaluates `f` and `args` immediately and runs the call when the surrounding function returns,** last in, first out, even during a panic. It's Go's `finally`/try-with-resources. Put the `defer` right after acquiring the resource. A `defer` inside a loop waits for the *function*, not the iteration.
- **Recursion works, but there's no tail-call elimination.** Goroutine stacks grow on demand up to 1 GB, and blowing past that is a fatal error, not a recoverable panic.
- **Habits to drop:** exceptions emulated with `panic`/`recover`, out-parameters, one-method "strategy" interfaces and factories where a `func` does the job, overloading emulated with `interface{}` arguments, and huge parameter lists instead of a struct.

## 2. Why It Exists

### The problems Go's design solves here

**Explicit, cheap, predictable calls.** Go has no overloading, no default arguments and no keyword arguments. One name means one function with one signature. The Go FAQ ("Why does Go not support overloading of methods and operators?") says experience showed that having several functions with the same name but different signatures was "occasionally useful but … could also be confusing and fragile". The price is more names (`strings.Split`, `strings.SplitN`) and struct parameters or functional options (35) for optional settings. The benefit: you always know which function runs, and a call site can be read without knowing argument types.

**Multiple results instead of exceptions and out-parameters.** C returns one value and writes the rest through pointers. Java and Python return one value and throw exceptions for failure. Go returns several values, so a function can return its result *and* an error as an ordinary value. Failure becomes part of the signature, visible at every call site (Go FAQ, "Why does Go not have exceptions?": coupling exceptions to a control structure "results in convoluted code"). The cost is the famous `if err != nil` repetition. The benefit is that every failure path is visible in the code a reviewer reads.

**First-class functions and closures without the ceremony.** Java added lambdas in Java 8 on top of "functional interfaces", and captured variables must be *effectively final*. Python has full closures, but its `lambda` is limited to one expression and rebinding a captured name needs `nonlocal`. Go has had full closures since 1.0: any function literal can read and assign any variable in scope, and the compiler moves those variables to the heap when they outlive the call. That's what makes `http.HandlerFunc`, `sort.Slice`, middleware and iterators lightweight.

**`defer`: cleanup next to acquisition.** Java's `try/finally` and try-with-resources, and Python's `with`, put the cleanup in a block structure around the code. Go's `defer` puts it on the line after the acquisition (`f, err := os.Open(p)` … `defer f.Close()`). You see "open" and "close" together, and every return path, including a panic, runs it. The trade-off: `defer` is function-scoped, not block-scoped, which surprises people in loops (§7.4). Since Go 1.14, most defers are compiled inline ("open-coded"), so the old advice to avoid `defer` in hot code is mostly obsolete (§10).

### What Python and Java do instead

| Concern | Python | Java | Go |
|---|---|---|---|
| Argument passing | Object references passed by value ("pass by assignment") | Primitives by value, objects by reference value | Everything by value: copies of ints, structs, arrays, headers, pointers |
| Optional / named arguments | Defaults, `*args`, `**kwargs` | Overloading, builders | None: a struct parameter, variadics, functional options (35) |
| Several results | Return a tuple | Return an object / `record` | Multiple results (not a tuple value) |
| Failure | Raise an exception | Throw (checked or unchecked) | Return `error` last; `panic` only for bugs (18) |
| Anonymous functions | `lambda` (one expression) or nested `def` | Lambdas, anonymous classes | Function literals: full bodies, any statements |
| Captured variables | Shared; rebinding needs `nonlocal` | Must be effectively final (copied) | Shared and assignable; the compiler moves them to the heap if needed |
| Function types | Any callable; `Callable[[int], int]` hints | Functional interfaces (`Function<T,R>`) | `func(int) int` is a type |
| Cleanup | `with`, `try/finally` | try-with-resources, `finally` | `defer`, function-scoped, LIFO |
| Recursion limit | ~1000 frames (`sys.getrecursionlimit`) | `StackOverflowError` (catchable); a fixed thread stack, set with `-Xss` | Growable goroutine stacks up to 1 GB; overflow is fatal |
| Tail calls | Not optimised | Not optimised | Not optimised |

## 3. Mental Model

**A call creates a new set of variables (the parameters and results) and assigns copies of the arguments to them. A function value is a pointer to code plus the variables it captured. `defer` pushes a fully evaluated call onto the function's own stack of "things to do on the way out".**

```
caller                                 callee frame (new variables for every call)
┌──────────────────────┐               ┌────────────────────────────────────┐
│ n   = 1              │── copy ──────►│ n   = 1     (changing it: local)    │
│ sl  = {ptr,len,cap} ─┼── copy ──────►│ sl  = {ptr,len,cap} ─┐              │
│ p   = &x ────────────┼── copy ──────►│ p   = &x             │              │
└──────────┬───────────┘               └──────────────────────┼─────────────┘
           │                                                  │
           ▼                                                  ▼
       x (caller's)                     shared backing array [ 1 | 2 ]
```

```
closure value                         captured variable (heap, if it outlives the call)
┌──────────────┬──────────────┐       ┌──────────┐
│ code pointer │ &count ──────┼──────►│ count: 3 │◄── every closure made by THIS call shares it
└──────────────┴──────────────┘       └──────────┘

defer stack of the current call (runs bottom-up at return):
  defer fmt.Println("c")   ← runs 1st
  defer fmt.Println("b")   ← runs 2nd
  defer fmt.Println("a")   ← runs 3rd
```

**Where it breaks:**
- **"Pass by value" doesn't mean "the callee can't change my data".** A copied slice header, map or pointer still points at shared data (§5.1). And `f(s...)` passes the caller's slice itself, with no copy at all (§5.3).
- **The compiler sometimes captures by value.** If a captured variable is never reassigned after capture (and is small), the compiler copies it into the closure instead of moving it to the heap (§5.5). You can't observe the difference, but it explains why some closures don't allocate.
- **A deferred *closure* doesn't evaluate anything early.** Only the function value and arguments in the `defer` statement are evaluated immediately. A closure body reads variables when it finally runs (§5.7).
- **Not every call is a real call.** Small functions are inlined, and inlined closures never allocate. Recursion depth is limited by stack memory, not by a counter.

## 4. Core Concepts & Terminology

| Term | Precise meaning (spec wording where possible) | Common confusion |
|---|---|---|
| **Signature** | The parameter and result types of a function; `func(a, b int) (int, error)` | Parameter *names* aren't part of the type |
| **Function type** | "the set of all functions with the same parameter and result types"; zero value `nil` | Func values are only comparable to `nil` |
| **Parameter / argument** | Parameters are the callee's variables; arguments are the values passed and *assigned* to them | "Go passes slices/maps by reference": no, it copies their small header/handle |
| **Multiple results** | A function can return several values; a call with several results is not a value you can store | No tuple type; `x := f()` needs as many variables as results |
| **Named results** | Result parameters with names; they're variables initialised to zero values | A bare `return` returns their current values |
| **Variadic parameter** | Final parameter `...T`; inside the function its type is `[]T` | `f(s...)` passes `s` itself, not a copy |
| **Function literal** | An anonymous function expression: `func(x int) int { return x * 2 }` | Not limited to one expression (unlike Python's `lambda`) |
| **Closure** | "Function literals are closures: they may refer to variables declared in a surrounding function. Those variables are then shared…" (spec) | Captures the variable, not its value |
| **Higher-order function** | Takes and/or returns functions | Common in the standard library: `sort.Slice`, `strings.Map`, `http.HandlerFunc` |
| **Recursion** | A function calling itself (directly or through others) | No tail-call elimination in Go |
| **`defer`** | Evaluates the function value and arguments now, invokes the call just before the surrounding function returns, LIFO | Function-scoped, not block-scoped |
| **Open-coded defer** | (Go 1.14+) The compiler inlines deferred calls at each return when the function has at most 8 defers, none in loops | Why `defer` is nearly free in typical code |
| **Escape analysis** | The compiler's decision whether a variable can live on the stack or must be on the heap (28) | Captured variables that outlive the call escape |
| **Inlining** | Replacing a call with the callee's body; `-gcflags=-m` prints `can inline` | Inlined calls have no frame and no call overhead |
| **`error`** | A predeclared interface type: `interface{ Error() string }`; `nil` means "no error" (18) | Not an exception; just a value returned last |

## 5. How It Works Inside

```mermaid
sequenceDiagram
  participant C as caller
  participant F as callee
  C->>C: evaluate the function value and the arguments (left to right)
  C->>F: assign copies to fresh parameters (in registers when they fit; Go 1.17+ ABI)
  F->>F: run the body; each "defer" pushes an evaluated call
  F->>F: "return x" assigns the result parameters
  F->>F: run deferred calls, last first (they may change named results)
  F-->>C: results (in registers)
```

The spec (*Calls*): "In a function call, the function value and arguments are evaluated in the usual order. After they are evaluated, new storage is allocated for the function's variables, which includes its parameters and results. Then, the arguments of the call are passed to the function, which means that they are assigned to their corresponding function parameters."

**Under the hood (go1.27, 64-bit).** Since Go 1.17 (amd64) and 1.18 (arm64 and others), the compiler passes arguments and results in registers when they fit ("ABIInternal"), with the stack as a fallback. Each goroutine has its own stack, which starts small (a few KiB) and is *copied* to a bigger one when a function's prologue sees it's about to run out. That's why deep recursion works, and why Go has no fixed thread-stack limit like Java's `-Xss`. The ceiling is 1 GB on 64-bit platforms (§5.6). Stacks and the heap in depth: topic 28.

### 5.1 Parameters are copies

Every parameter is a new variable holding a copy of the argument. What that copy *shares* depends on what the value contains:

| Argument type | What's copied | Can the callee change what the caller sees? |
|---|---|---|
| `int`, `string`, `bool`, struct, **array** | The whole value | No |
| pointer `*T` | The address | Yes: `*p = …` writes the caller's variable |
| slice `[]T` | The 24-byte header `{ptr, len, cap}` | Existing elements yes; `append` and re-slicing no (topic 06) |
| map, channel | A pointer to the runtime's map/channel | Yes: it's the same map/channel |
| func value | A pointer to the closure | It can call it, and the closure's captured state is shared |

**Predict the output** ([`internal/examples/basics.go`](../internal/examples/basics.go) · `make example N=pass-by-value`):
```go
func mutate(n int, arr [2]int, sl []int, m map[string]int, p *int) {
	n = 100            // changes the local copy
	arr[0] = 100       // changes the local copy of the array
	sl[0] = 100        // writes through the copied header into the shared array
	sl = append(sl, 3) // grows the local header only
	m["b"] = 2         // the map is shared
	*p = 42            // writes the caller's variable
}

n, arr, sl, m, x := 1, [2]int{1, 2}, []int{1, 2}, map[string]int{"a": 1}, 1
mutate(n, arr, sl, m, &x)
```
<details><summary>Output</summary>

```
n:   1     (the int was copied)
x:   42    (written through the pointer)
arr: [1 2] (the whole array was copied)
sl:  [100 2] (element write shared; the append was not)
m:   map[a:1 b:2] (same map)
```
Why: "Go passes slices and maps by reference" is a common but wrong shortcut. Everything is copied. A slice header and a map handle are just values that *point at* shared data. That's why `sl[0] = 100` is visible but `append` (which changed the callee's own header) isn't. Value vs pointer semantics in depth: topic 09; the slice header: topic 06.
</details>

### 5.2 Multiple results and named results

A function can return any number of results. A call with several results can be assigned with `:=`/`=`, partly discarded with `_`, or passed straight into a call whose parameters match exactly (`add(divmod(17, 5))`). It can't be stored in one variable: there's no tuple type.

**Predict the output** ([`internal/examples/basics.go`](../internal/examples/basics.go) · `make example N=multiple-returns`):
```go
func divmod(a, b int) (int, int) { return a / b, a % b }
func add(a, b int) int { return a + b }

q, r := divmod(17, 5)
_, r = divmod(-17, 5)
if n, err := strconv.Atoi("12a"); err != nil {
	fmt.Fprintln(w, "Atoi failed:", err, "| n is the zero value:", n)
}
fmt.Fprintln(w, "sum of divmod:", add(divmod(17, 5)))
```
<details><summary>Output</summary>

```
divmod(17, 5) = 3 2
remainder only: -2
Atoi failed: strconv.Atoi: parsing "12a": invalid syntax | n is the zero value: 0
sum of divmod: 5
```
Why: `-17 % 5` is `-2` (the remainder takes the dividend's sign, [03 §5.3](../../03-basic-types-and-operators/theory/README.md#53-division-remainder-and-shifts)). On failure, the convention is to return zero values for the other results, so callers never use them.
</details>

**Named results** are ordinary variables declared in the signature and initialised to their zero values. `return` with no operands (a *bare* or *naked* return) returns their current values:

**Predict the output** ([`internal/examples/basics.go`](../internal/examples/basics.go) · `make example N=named-results`):
```go
func minMax(xs []int) (lo, hi int) {
	if len(xs) == 0 {
		return // lo and hi are still 0
	}
	lo, hi = xs[0], xs[0]
	for _, x := range xs[1:] {
		lo = min(lo, x)
		hi = max(hi, x)
	}
	return lo, hi // explicit, even though "return" alone would work
}
```
<details><summary>Output</summary>

```
minMax: -2 9
minMax(nil): 0 0 (zero values)
```
Why: named results earn their place when they **document** results of the same type (`(lo, hi int)` beats `(int, int)`) or when a deferred closure must change a result (§5.7). Bare returns in anything longer than a few lines hurt readability (Go Code Review Comments, *Named Result Parameters* and *Naked Returns*).
</details>

### 5.3 Variadic parameters

A final parameter `...T` accepts zero or more `T` arguments. Inside the function it's a `[]T`. Called with individual values, Go builds a new slice for them, and with none the slice is `nil`. Called with `s...`, "it is passed unchanged as the value for a ...T parameter. In this case no new slice is created" (spec).

```
sum(1, 2, 3)  →  new array [1|2|3], nums = {ptr→it, len 3, cap 3}
sum()         →  nums = nil
sum(xs...)    →  nums = xs   (same header, same backing array: no copy)
```

**Predict the output** ([`internal/examples/basics.go`](../internal/examples/basics.go) · `make example N=variadic` and `N=variadic-aliasing`):
```go
func sum(nums ...int) int { /* adds them up */ }
func describe(label string, nums ...int) string {
	return fmt.Sprintf("%s: len=%d nil=%t", label, len(nums), nums == nil)
}
func double(nums ...int) {
	for i := range nums {
		nums[i] *= 2
	}
}

fmt.Fprintln(w, sum(), sum(1), sum(1, 2, 3))
fmt.Fprintln(w, sum([]int{10, 20, 30}...))
fmt.Fprintln(w, describe("ids"), describe("ids", 7, 8))

scores := []int{1, 2, 3}
double(scores...)
fmt.Fprintln(w, "after double(scores...):", scores)
```
<details><summary>Output</summary>

```
0 1 6
60
ids: len=0 nil=true ids: len=2 nil=false
```
```
after double(scores...): [2 4 6]
after double(4, 5, 6): nothing visible to change
```
Why: `double(scores...)` received the caller's slice, so its writes changed `scores`. A variadic function that keeps or modifies its argument (stores it in a struct, sorts it) should copy it first (`slices.Clone`). You can't mix the two forms: `sum(0, xs...)` is a compile error (§7.8).
</details>

### 5.4 Functions are values

A function type like `func(int, int) int` describes every function with those parameter and result types. Function values can be assigned, stored in structs/slices/maps, passed and returned. The zero value is `nil`, and **calling a nil func panics**. Func values can only be compared to `nil`.

**Predict the output** ([`internal/examples/funcvalues.go`](../internal/examples/funcvalues.go) · `make example N=function-values`):
```go
var op func(int, int) int // nil
fmt.Fprintln(w, "zero value is nil:", op == nil)

op = add
fmt.Fprintln(w, "op = add:", op(2, 3))

op = func(a, b int) int { return a * b } // an anonymous function
fmt.Fprintln(w, "op = multiply:", op(2, 3))

type transform func(string) string
steps := []transform{strings.TrimSpace, strings.ToUpper, func(s string) string { return s + "!" }}
s := "  hello "
for _, step := range steps {
	s = step(s)
}
fmt.Fprintf(w, "pipeline: %q\n", s)
fmt.Fprintf(w, "type of add: %T\n", add)
```
<details><summary>Output</summary>

```
zero value is nil: true
op = add: 5
op = multiply: 6
pipeline: "HELLO!"
type of add: func(int, int) int
```
Why: `strings.TrimSpace` is assignable to `transform` because its type `func(string) string` has the same underlying type (defined types: topic 13). Named function types are how the standard library turns functions into interface implementations, e.g. `http.HandlerFunc` (34).
</details>

**Higher-order functions** take or return functions:

**Predict the output** ([`internal/examples/funcvalues.go`](../internal/examples/funcvalues.go) · `make example N=higher-order`):
```go
func filter(xs []int, keep func(int) bool) []int {
	var out []int
	for _, x := range xs {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

// compose returns f∘g: a function that applies g, then f.
func compose(f, g func(int) int) func(int) int {
	return func(x int) int { return f(g(x)) }
}

nums := []int{1, 2, 3, 4, 5, 6}
fmt.Fprintln(w, "filter even:", filter(nums, func(n int) bool { return n%2 == 0 }))
inc := func(n int) int { return n + 1 }
sq := func(n int) int { return n * n }
fmt.Fprintln(w, "compose(sq, inc)(3) =", compose(sq, inc)(3), "(inc first, then sq)")
fmt.Fprintln(w, "compose(inc, sq)(3) =", compose(inc, sq)(3))
```
<details><summary>Output</summary>

```
filter even: [2 4 6]
filter > 3:  [4 5 6]
compose(sq, inc)(3) = 16 (inc first, then sq)
compose(inc, sq)(3) = 10
```
Why: in Go, a plain loop is usually clearer than chains of higher-order helpers. Use them where they name an idea (a predicate, a strategy, a callback), not to imitate streams (§8.4). The standard library's `slices` package covers the common ones generically (06, 14).
</details>

**A nil func value** ([`internal/examples/funcvalues.go`](../internal/examples/funcvalues.go) · `make example N=nil-func-call`, crashes on purpose). A real run:
```go
var onDone func(string)
fmt.Fprintln(w, "calling a nil func value...")
onDone("finished")
```
```
calling a nil func value...
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x2 addr=0x0 pc=0x100d0ef48]

goroutine 1 [running]:
github.com/Shreyash019/Golang/topics/05-functions-closures-and-defer/internal/examples.NilFuncCall({0x100de8090?, 0x6900b3482010?})
	…/topics/05-functions-closures-and-defer/internal/examples/funcvalues.go:67 +0x48
```
Optional callbacks (`OnRetry func(error)` in a config struct) need an `if cb != nil` check, or a no-op default.

### 5.5 Closures: functions that capture variables

A function literal can use variables from the function around it. "Those variables are then shared between the surrounding function and the function literal, and they survive as long as they are accessible" (spec). The closure holds a reference to the **variable**, so it sees later assignments and its own assignments are seen outside.

```
counter() call #1                         counter() call #2
  count₁ (heap)  ◄── a's closure            count₂ (heap)  ◄── b's closure
  a := counter()                            b := counter()
```

**How it's compiled.** A func value is a pointer to a small block: the code pointer, then the captured variables. A captured variable that's **reassigned** after capture is captured *by reference*: the variable moves to the heap and the closure stores its address. A variable that's never reassigned is captured *by value* (copied in). Real output from a scratch package (`go build -gcflags=-m=2`):

```go
func Counter() func() int {          func Greeter(prefix string) func(string) string {
	count := 0                           return func(name string) string { return prefix + name }
	return func() int {              }
		count++
		return count
	}
}
```
```
esc/cap.go:3:14: Greeter capturing by value: prefix (addr=false assign=false width=16)
esc/esc.go:4:2: Counter capturing by ref: count (addr=false assign=true width=8)
esc/esc.go:4:2: moved to heap: count
```
A closure that doesn't outlive its call (used locally, e.g. `add := func(x int) { total += x }` in a loop) is usually inlined and allocates nothing. The heap only gets involved when the closure *escapes*. Escape analysis in depth: topic 28.

**Predict the output** ([`internal/examples/closures.go`](../internal/examples/closures.go) · `make example N=closure-counter`):
```go
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

a := counter()
b := counter()
fmt.Fprintln(w, a(), a(), a(), "| a has its own count")
fmt.Fprintln(w, b(), b(), "| b has another")
fmt.Fprintln(w, a(), "| a continues")
```
<details><summary>Output</summary>

```
1 2 3 | a has its own count
1 2 | b has another
4 | a continues
```
Why: each *call* to `counter` creates a new `count`, and each closure keeps its own alive.
</details>

**Predict the output** ([`internal/examples/closures.go`](../internal/examples/closures.go) · `make example N=closure-shares-variable`):
```go
greeting := "hello"
say := func() string { return greeting + ", world" }
greeting = "goodbye" // changed after the closure was created
fmt.Fprintln(w, say())

total := 0
addTo := func(n int) { total += n }
read := func() int { return total }
addTo(5)
addTo(10)
fmt.Fprintln(w, "read():", read(), "| total:", total)

v := 1
snapshot := v // to capture the current value, copy it first
get := func() int { return snapshot }
v = 2
fmt.Fprintln(w, "snapshot:", get(), "| v:", v)
```
<details><summary>Output</summary>

```
goodbye, world
read(): 15 | total: 15
snapshot: 1 | v: 2
```
Why: closures share variables. Two closures over the same variable see each other's changes. To freeze a value, copy it into a variable nobody reassigns. Before Go 1.22, loops had *one* shared loop variable, so every closure created in a loop saw its last value ([04 §5.3](../../04-control-flow/theory/README.md#53-loop-variables-one-per-iteration-go-122)).
</details>

**Shared private state.** Two closures from the same call share their captured variables. It's a lightweight alternative to a small struct with methods (10):

**Predict the output** ([`internal/examples/closures.go`](../internal/examples/closures.go) · `make example N=closure-state`):
```go
func idGenerator(prefix string) (next func() string, reset func()) {
	id := 0
	next = func() string {
		id++
		return fmt.Sprintf("%s%d", prefix, id)
	}
	reset = func() { id = 0 }
	return next, reset
}

next, reset := idGenerator("order-")
fmt.Fprintln(w, next(), next(), next())
reset()
fmt.Fprintln(w, next())
```
<details><summary>Output</summary>

```
order-1 order-2 order-3
order-1
```
Why: once there are more than two or three operations, or the state needs to be inspected or tested, a struct with methods reads better. Closures aren't safe for concurrent use either: `id++` from two goroutines is a data race (24, 25).
</details>

**A recursive function literal** needs its variable declared before the literal is assigned, because in `fib := func…` the name `fib` isn't in scope inside its own initialiser (§7.10):

**Predict the output** ([`internal/examples/closures.go`](../internal/examples/closures.go) · `make example N=recursive-closure`):
```go
var fib func(n int) int
fib = func(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
fmt.Fprintln(w, "fib(10) =", fib(10))
```
<details><summary>Output</summary>

```
fib(10) = 55
```
E5's parser uses exactly this for three mutually recursive closures.
</details>

### 5.6 Recursion and the growable stack

A recursive function needs a **base case** that stops, and a **recursive case** that moves toward it. Each call gets its own frame (its own parameters and locals) on the goroutine's stack. Go does **not** eliminate tail calls, so even `return gcd(b, a%b)` uses a new frame per call. Trees, nested data (JSON, ASTs, file systems) and divide-and-conquer algorithms are where recursion reads best.

**Predict the output** ([`internal/examples/recursion.go`](../internal/examples/recursion.go) · `make example N=recursion`):
```go
func factorial(n int) int {
	if n <= 1 { // base case
		return 1
	}
	return n * factorial(n-1) // recursive case: n shrinks toward 1
}

func printTree(w io.Writer, n *node, depth int) {
	fmt.Fprintf(w, "%s%s\n", strings.Repeat("  ", depth), n.name)
	for _, c := range n.children {
		printTree(w, c, depth+1)
	}
}
```
<details><summary>Output</summary>

```
factorial(10) = 3628800
gcd(84, 36) = 12
/
  etc
    hosts
  usr
    bin
      go
    lib
files: 7
```
</details>

**How the stack grows.** Almost every function starts with a check (the compiler skips it only for small leaf functions): is there room for my frame? If not, the runtime allocates a stack twice the size, copies the old one over, adjusts pointers into it, and continues. Growth is amortised and invisible, so a goroutine can start with a few KiB and still recurse deeply. The limit is 1 GB on 64-bit platforms. Exceeding it is a **fatal error** that `recover` can't catch, and it kills the whole process (`make example N=stack-overflow`, crashes on purpose). A real run, first lines:
```
recursing forever...
runtime: goroutine stack exceeds 1000000000-byte limit
runtime: sp=0x449b99160390 stack=[0x449b99160000, 0x449bb9160000]
fatal error: stack overflow
```
staticcheck flags the obvious case (`SA5007: infinite recursive call`). For recursion driven by **untrusted input** (deeply nested JSON, a user-supplied tree), bound the depth or use an explicit stack (a slice) instead.

### 5.7 `defer`

The spec (*Defer statements*): "Each time a "defer" statement executes, the function value and parameters to the call are evaluated as usual and saved anew but the actual function is not invoked. Instead, deferred functions are invoked immediately before the surrounding function returns, in the reverse order they were deferred. That is, if the surrounding function returns through an explicit return statement, deferred functions are executed after any result parameters are set by that return statement but before the function returns to its caller."

Four rules follow from that, one example each:

**1. LIFO, at return** ([`internal/examples/defers.go`](../internal/examples/defers.go) · `make example N=defer-order`):
```go
fmt.Fprintln(w, "start")
defer fmt.Fprintln(w, "deferred 1")
defer fmt.Fprintln(w, "deferred 2")
defer fmt.Fprintln(w, "deferred 3")
fmt.Fprintln(w, "end of function body")
```
<details><summary>Output</summary>

```
start
end of function body
deferred 3
deferred 2
deferred 1
```
Why LIFO: resources are released in the reverse order they were acquired (unlock before closing the connection that the locked data came from), just like nested `with` blocks.
</details>

**2. Arguments are evaluated at the `defer` statement; closures read variables later** (`make example N=defer-args-evaluated`):
```go
x := 1
defer fmt.Fprintln(w, "deferred call, argument evaluated at defer time: x =", x)
defer func() { fmt.Fprintln(w, "deferred closure, reads x when it runs:     x =", x) }()
x = 2
fmt.Fprintln(w, "x set to", x)
```
<details><summary>Output</summary>

```
x set to 2
deferred closure, reads x when it runs:     x = 2
deferred call, argument evaluated at defer time: x = 1
```
Why: `defer fmt.Fprintln(w, …, x)` copied `x` (1) into the saved call. The closure has no arguments; its body reads `x` when it runs. This is the root of §7.3's `defer log(err)` bug.
</details>

**3. Deferred closures can read and change named results** (`make example N=defer-modifies-result`):
```go
func triple(n int) (result int) {
	defer func() { result *= 3 }() // runs after "return n" sets result = n
	return n
}

func loadUser(name string) (user string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("load user %q: %w", name, err)
		}
	}()
	if name == "" {
		return "", errors.New("empty name")
	}
	return strings.ToUpper(name), nil
}
```
<details><summary>Output</summary>

```
triple(5) = 15
loadUser(""): load user "": empty name
loadUser("gopher"): <nil>
```
Why: `return n` first assigns `result = n`, then the deferred closure runs and can change `result`, then the function returns. It only works with **named** results: with unnamed results the deferred closure has no name to assign (E1 line 3). Use it for annotating every error a function returns, and for capturing a `Close` error (§7.6). Don't use it for everyday logic; it's surprising.
</details>

**4. Deferred calls run during a panic** (`make example N=defer-runs-on-panic`, crashes on purpose). A real run:
```go
defer fmt.Fprintln(w, "deferred: cleanup ran during the panic")
fmt.Fprintln(w, "about to index out of range")
var xs []int
_ = xs[3]
```
```
about to index out of range
deferred: cleanup ran during the panic
panic: runtime error: index out of range [3] with length 0

goroutine 1 [running]:
github.com/Shreyash019/Golang/topics/05-functions-closures-and-defer/internal/examples.DeferRunsOnPanic({0x100fe0090?, 0x2053fae8a010?})
	…/topics/05-functions-closures-and-defer/internal/examples/defers.go:66 +0x7c
```
That's why `defer mu.Unlock()` and `defer f.Close()` are safe even when the code between panics. A deferred call can also stop a panic with `recover()`. That's topic 18, and it isn't an exception mechanism (§8.2).

**`defer` is function-scoped** (`make example N=defer-in-loop`). In a loop, every iteration adds a pending call, and they all run when the *function* returns ([04 §7.10](../../04-control-flow/theory/README.md#710-defer-in-a-loop-runs-at-function-end-not-at-iteration-end)). Move the body into a function so each iteration's `defer` runs at the end of that call:
```go
for _, name := range []string{"a", "b", "c"} {
	func() {
		release := acquire(name)
		defer release() // runs at the end of this iteration
	}()
}
```
<details><summary>Output</summary>

```
-- defer directly in the loop:
open a (open now: 1)
open b (open now: 2)
open c (open now: 3)
loop finished with 3 open
close c
close b
close a
-- body moved into a function:
open a (open now: 1)
close a
open b (open now: 1)
close b
open c (open now: 1)
close c
loop finished with 0 open
```
In real code, prefer a **named** helper (`processFile(name) error`) over an anonymous func called in place. It can return an error and be tested on its own. E2 is this bug.
</details>

**How `defer` is implemented.** Since Go 1.14, if a function has at most 8 `defer` statements and none of them is in a loop, the compiler **open-codes** them: it stores the deferred function and arguments in the frame, sets a bit per executed `defer`, and calls them directly at each return. That costs about one or two nanoseconds (§10). Otherwise (a `defer` in a loop, or too many), the runtime keeps a chain of defer records, which is slower. When a panic unwinds, the runtime finds open-coded defers through metadata the compiler emitted and runs them.

### 5.8 A first look at returning `error`

`error` is a predeclared interface type with one method, `Error() string`. A `nil` error means success. The convention, everywhere in Go:
- `error` is the **last** result: `func Parse(s string) (Config, error)`.
- On failure, return zero values for the other results plus a non-nil error. On success, return `nil`.
- Check immediately, handle or return, and add context with `fmt.Errorf("doing X: %w", err)`. The `%w` verb *wraps* the error so callers can still identify it with `errors.Is`.
- Package-level error values (`var ErrNotFound = errors.New("not found")`) are **sentinel errors**, compared with `errors.Is`, not `==`, because they may be wrapped.

**Predict the output** ([`internal/examples/errors_first.go`](../internal/examples/errors_first.go) · `make example N=error-first-look`):
```go
var ErrNegative = errors.New("quantity is negative")

func parseQuantity(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parse quantity %q: %w", s, err)
	}
	if n < 0 {
		return 0, fmt.Errorf("parse quantity %q: %w", s, ErrNegative)
	}
	return n, nil
}

for _, in := range []string{"3", "x", "-2"} {
	q, err := parseQuantity(in)
	if err != nil {
		fmt.Fprintf(w, "%-4q error: %v (negative? %t)\n", in, err, errors.Is(err, ErrNegative))
		continue
	}
	fmt.Fprintf(w, "%-4q ok: %d\n", in, q)
}
```
<details><summary>Output</summary>

```
"3"  ok: 3
"x"  error: parse quantity "x": strconv.Atoi: parsing "x": invalid syntax (negative? false)
"-2" error: parse quantity "-2": quantity is negative (negative? true)
```
Why: each layer adds what it was doing, so the message reads like a stack trace in words. Error types, `errors.As`, joining, `panic` vs `error` and `recover` are all topic 18.
</details>

## 6. Worked Examples

Every example in §5 lives in [`internal/examples/`](../internal/examples/) and is checked by an Example test (crashing ones are marked and were run by hand). `make examples` lists them:

| Concept | Example | Section |
|---|---|---|
| Parameters are copies | `pass-by-value` | 5.1 |
| Multiple and named results | `multiple-returns`, `named-results` | 5.2 |
| Variadic parameters, `s...` aliasing | `variadic`, `variadic-aliasing` | 5.3 |
| Function values, higher-order functions, nil funcs | `function-values`, `higher-order`, `nil-func-call` (crashes) | 5.4 |
| Closures | `closure-counter`, `closure-shares-variable`, `closure-state`, `recursive-closure` | 5.5 |
| Recursion and the stack limit | `recursion`, `stack-overflow` (crashes) | 5.6 |
| `defer` | `defer-order`, `defer-args-evaluated`, `defer-modifies-result`, `defer-runs-on-panic` (crashes), `defer-in-loop` | 5.7 |
| Returning `error` | `error-first-look` | 5.8 |
| The trace idiom | `defer-trace` | 6.1 |
| `os.Exit` skips deferred calls | `go run ./cmd/exit-skips-defer` | 7.5 |

### 6.1 The trace idiom: `defer trace(name)()`

Effective Go (*Defer*) shows a trick that follows from "the arguments are evaluated now": `defer un(trace("a"))` runs `trace` immediately, on entry, and defers `un`. The closure form below goes one step further: `trace` returns the function to defer, so the call on entry and the cleanup on exit share state (here, `depth`).

**Predict the output** ([`internal/examples/defers.go`](../internal/examples/defers.go) · `make example N=defer-trace`):
```go
depth := 0
trace := func(name string) func() {
	fmt.Fprintf(w, "%senter %s\n", strings.Repeat("  ", depth), name)
	depth++
	return func() {
		depth--
		fmt.Fprintf(w, "%sleave %s\n", strings.Repeat("  ", depth), name)
	}
}

inner := func() {
	defer trace("inner")()
	fmt.Fprintf(w, "%s(inner's work)\n", strings.Repeat("  ", depth))
}
outer := func() {
	defer trace("outer")()
	inner()
	inner()
}
outer()
```
<details><summary>Output</summary>

```
enter outer
  enter inner
    (inner's work)
  leave inner
  enter inner
    (inner's work)
  leave inner
leave outer
```
Why: in `defer trace("inner")()`, the function value to defer is `trace("inner")`, which is evaluated **now** (so "enter" prints now). The returned closure is what runs at return. The same shape times a function: `defer timer("load")()` where `timer` captures `time.Now()` and logs `time.Since(start)`. Forgetting the final `()` is a classic slip: `defer trace("inner")` defers the *whole* trace call and prints "enter" at exit. E5 builds its parse trace on this idea.
</details>

### 6.2 A function's shape, from signature to return

```go
// ParseConfig reads key=value lines from r. Unknown keys are an error.
func ParseConfig(ctx context.Context, r io.Reader, opts ParseOptions) (Config, error) {
	//            ^ context first        ^ inputs     ^ options as a struct   ^ error last
	var cfg Config
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if err := ctx.Err(); err != nil { // cancellation: topic 26
			return Config{}, err
		}
		if err := cfg.set(sc.Text(), opts); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err) // zero value + context
		}
	}
	if err := sc.Err(); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}
```
*(Illustrative: `Config.set` is a method (10) and `context` is topic 26.)* Context first, `error` last, an options struct instead of a long parameter list, guard clauses, zero values on failure, one wrap per layer. Nearly every exported Go function you write for a service has this shape.

## 7. Gotchas

Compile errors and runs below were produced with go1.27.1 (and golangci-lint 2.14) in a scratch module (`example.com/gotchas`) and copied from the terminal, unless they're examples in this module.

| # | Gotcha | Symptom | Fix |
|---|---|---|---|
| 7.1 | Closures in loops before Go 1.22 | Every closure sees the last value | `go 1.22`+ in `go.mod` (04 §7.2) |
| 7.2 | Closures share variables | "Random" values; data races between goroutines | Copy the value first; synchronise |
| 7.3 | `defer` arguments evaluated immediately | `defer log(err)` logs `<nil>` | Defer a closure |
| 7.4 | `defer` in a loop | Resources held until the function returns | Move the body into a function (E2) |
| 7.5 | `os.Exit` / `log.Fatal` skip deferred calls | Unflushed buffers, unclosed files | Return from `main` via `run() error` |
| 7.6 | Ignoring the error from `Close` on a written file | Silent data loss | Named result + deferred close check |
| 7.7 | `defer x.Close()` before checking the error | Nil pointer panic | Check `err` first, then defer |
| 7.8 | Variadic surprises | Caller's slice modified; `[]string` → `...any` doesn't compile | Copy; build a `[]any` |
| 7.9 | Shadowed named results with a bare return | `result parameter … not in scope at return` | Use `=`, or return explicitly |
| 7.10 | Recursive closure declared with `:=` | `undefined: fib` | `var fib func…` first |
| 7.11 | Unbounded recursion | `fatal error: stack overflow` (not recoverable) | Base case; depth limit for untrusted input |
| 7.12 | Calling a nil func value | Nil pointer panic | Check for nil or use a no-op default |
| 7.13 | Overloading, default parameters, comparing funcs | Compile errors | Different names, option structs, compare to nil only |
| 7.14 | Missing return / mismatched results | `missing return`; `assignment mismatch` | Every path returns; take every result |
| 7.15 | Silently ignored results | `errors.New("x")` alone: vet; an ignored `error`: nothing | `go vet`, errcheck |

### 7.1 Closures created in a loop before Go 1.22

Recap of [04 §7.2](../../04-control-flow/theory/README.md#72-loop-variables-captured-before-go-122): in a module whose `go.mod` says `go 1.21` or lower, a loop has one variable for all iterations, so closures (and goroutines) created in the loop all see its final value. From `go 1.22` on, each iteration has its own variable. You'll still see `v := v` lines in older code. They're harmless and can be deleted after upgrading.

### 7.2 Closures share variables, including between goroutines

See §5.5 (`make example N=closure-shares-variable`): a closure reads the variable's **current** value when it runs. Two common forms of the bug:
- A callback registered now and run later (an HTTP handler, a `time.AfterFunc`, a deferred closure) reads a variable that was changed in between.
- Closures run as goroutines (`go func() { total += n }()`) write a shared captured variable at the same time: a **data race**, which `go test -race` reports and which loses updates (22, 25).

**Fix:** copy the value you mean into a new variable (or pass it as an argument: `go func(n int) { … }(n)`), and protect shared state with a mutex or channel (24).

### 7.3 `defer` arguments are evaluated immediately

A real run:
```go
func save() (err error) {
	defer fmt.Println("save finished, err =", err)
	err = errors.New("disk full")
	return err
}

func main() {
	fmt.Println("returned:", save())
}
```
```
save finished, err = <nil>
returned: disk full
```
**Why:** `err` was `nil` when the `defer` statement ran, and that's the value saved with the call (§5.7 rule 2). `go vet` and the default linters say nothing. **Fix:** defer a closure, which reads `err` when it runs: `defer func() { log.Printf("save finished: %v", err) }()`. The same trap hits `defer timer(time.Now())` (fine: you *want* the start time now) versus `defer log.Println(time.Since(start))` (wrong: computed at the `defer`, so about zero).

### 7.4 `defer` in a loop

See §5.7 (`make example N=defer-in-loop`) and [04 §7.10](../../04-control-flow/theory/README.md#710-defer-in-a-loop-runs-at-function-end-not-at-iteration-end). Opening files, rows or locks in a loop with `defer` in the loop body holds all of them until the function returns: "too many open files", exhausted connection pools, a lock held far too long. A `defer` in a loop is also never open-coded, so each one costs more (§10). **Fix:** one function per iteration (E2).

### 7.5 `os.Exit` and `log.Fatal` skip deferred calls

A real run of [`cmd/exit-skips-defer`](../cmd/exit-skips-defer/main.go):
```go
func main() {
	defer fmt.Println("deferred: flush and close") // never printed
	fmt.Println("work done, exiting with status 1")
	os.Exit(1)
}
```
```
$ go run ./cmd/exit-skips-defer
work done, exiting with status 1
exit status 1
```
**Why:** `os.Exit` ends the process immediately. Deferred calls are a feature of *returning* from a function, and `os.Exit` never returns. `log.Fatal*` call `os.Exit(1)`. **Fix:** keep `main` tiny: `if err := run(); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }`, with all the work and all the `defer`s inside `run`, which returns normally. Every `main` in these exercises has that shape. Don't call `log.Fatal` below `main`.

### 7.6 Ignoring the `Close` error on a file you wrote

For a file you **wrote**, `Close` can be the call that reports the error (for example a failed flush to a network file system), so `defer f.Close()` silently drops it. For a file you only **read**, ignoring it is fine (E2's solution says so). The idiomatic pattern uses a named result and a deferred closure:
```go
func writeReport(path string, data []byte) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr // report Close's error only if nothing failed earlier
		}
	}()
	_, err = f.Write(data)
	return err
}
```
Tooling: errcheck (on by default in golangci-lint) reports ``Error return value of `f.Close` is not checked`` for `defer f.Close()` on an `*os.File`. It didn't report it for the `io.ReadCloser` in E2 with the default settings. `errors.Join(err, f.Close())` (Go 1.20) is a shorter alternative when you want both errors (18).

### 7.7 Deferring `Close` before checking the error

A real run (the `conn` type is illustrative; methods are topic 10):
```go
c, err := dial("db:5432") // returns nil, error
defer c.Close()           // registered before the error check
if err != nil {
	fmt.Println("dial failed:", err)
	return
}
```
```
dial failed: connection refused
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x2 addr=0x0 pc=0x104473568]

goroutine 1 [running]:
main.(*conn).Close(0x2bde13d36000?)
```
**Why:** on failure, `c` is `nil`, and the deferred `Close` dereferences it at return. The real-world version is `resp, err := http.Get(url); defer resp.Body.Close()` (34). **Fix:** always check `err` first, then `defer`.

### 7.8 Variadic surprises

Passing `s...` gives the callee your slice (§5.3, `make example N=variadic-aliasing`): a variadic function that sorts, modifies or stores its argument changes or keeps the caller's data. And the `...` form must match the element type exactly, alone:
```go
names := []string{"a", "b"}
fmt.Println(names...)
```
```
./main.go:7:14: cannot use names (variable of type []string) as []any value in argument to fmt.Println
```
```go
func sum(nums ...int) int
_ = sum(0, xs...)
```
```
./main.go:7:13: too many arguments in call to sum
	have (number, []int...)
	want (...int)
```
**Why:** `[]string` and `[]any` have different memory layouts (a string vs an interface value per element), so there's no conversion between them, and `s...` must be the *only* argument for the variadic part. **Fix:** build the slice you need (`args := make([]any, len(names)); for i, n := range names { args[i] = n }`, or `fmt.Println(strings.Join(names, " "))`), or `append([]int{0}, xs...)`.

### 7.9 A bare return with shadowed named results

```go
func parse(s string) (n int, err error) {
	if s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return
		}
		_ = n
	}
	return
}
```
```
./main.go:9:4: result parameter n not in scope at return
	./main.go:7:3: inner declaration of var n int
./main.go:9:4: result parameter err not in scope at return
	./main.go:7:6: inner declaration of var err error
```
**Why:** `:=` in the inner block declared *new* `n` and `err` ([02 §7.1](../../02-variables-and-constants/theory/README.md#71--in-an-inner-block-shadows-the-outer-variable)). A bare return there would return the outer, unset ones, so the compiler refuses. Helpful, but only for bare returns: with `return n, err` the shadowing compiles and may return the wrong variables. **Fix:** use `=` to assign the named results, or return values explicitly.

### 7.10 A recursive closure declared with `:=`

```go
fib := func(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
```
```
./main.go:8:10: undefined: fib
```
**Why:** a variable declared with `:=` is in scope only after its declaration ends, which is after the function literal. **Fix:** `var fib func(int) int` first, then `fib = func…` (§5.5).

### 7.11 Unbounded recursion is fatal

See §5.6: past 1 GB of stack the runtime prints `fatal error: stack overflow` and the process exits. Unlike Java's `StackOverflowError`, it can't be caught. Recursion on attacker-controlled depth (nested JSON arrays, a recursive data structure from a request) is a denial-of-service risk. **Fix:** a correct base case, plus a depth limit or an explicit stack for untrusted input. staticcheck catches only the obvious "no base case at all" (SA5007).

### 7.12 Calling a nil func value

See §5.4 (`make example N=nil-func-call`): a nil pointer dereference panic. **Fix:** optional callbacks get a nil check (`if s.onEvict != nil { s.onEvict(k) }`) or a no-op default set in the constructor. A `defer`red nil func panics when it's *invoked* at return, not at the `defer` statement (spec), which makes the stack trace point at the end of the function.

### 7.13 No overloading, no default parameters, no func comparison

```go
func area(w, h int) int { return w * h }
func area(side int) int { return side * side }
```
```
./main.go:4:6: area redeclared in this block
	./main.go:3:6: other declaration of area
```
```go
func greet(name string = "world") string { return "hi " + name }
```
```
./main.go:3:24: syntax error: unexpected = in parameter list; possibly missing comma or )
```
```go
f := func() {}
g := func() {}
println(f == g)
```
```
./main.go:6:10: invalid operation: f == g (func can only be compared to nil)
```
**Fix:** different names (`Area`, `SquareArea`), an options struct whose zero value means "defaults", functional options (35), or variadics. To tell callbacks apart, keep an ID next to them.

### 7.14 Missing returns and result-count mismatches

```go
func sign(n int) int {
	if n < 0 {
		return -1
	} else if n > 0 {
		return 1
	}
}
```
```
./main.go:9:1: missing return
```
```go
n := strconv.Atoi("42")
```
```
./main.go:6:7: assignment mismatch: 1 variable but strconv.Atoi returns 2 values
```
**Why:** a function with results must end in a *terminating statement* (04 §4), and the compiler doesn't reason about whether your conditions are exhaustive. Every result must be received or explicitly discarded with `_`. **Fix:** end with an unconditional `return`, and write `n, err := …` (then handle `err`).

### 7.15 Ignored results compile

```go
errors.New("x")
```
compiles, and `go vet` reports `result of errors.New call not used`. But calling a function and ignoring its `error` result (`f.Write(b)`, `json.Unmarshal(data, &v)`) compiles **and passes vet**. errcheck (default in golangci-lint) catches it: ``Error return value of `f.Write` is not checked``. Run `make lint`. Where ignoring is deliberate, say so: `_ = f.Close() // read-only`.

## 8. Idiomatic vs Non-Idiomatic

| ❌ Non-idiomatic | ✅ Idiomatic | Why (source) |
|---|---|---|
| `func Parse(s string, out *Config) bool` | `func Parse(s string) (Config, error)` | Effective Go, *Multiple return values* |
| `panic(err)` for bad input, `recover()` to "catch" it | Return `error`; `panic` only for programmer bugs | Effective Go, *Errors*/*Panic*; Go Code Review Comments, *Don't Panic* |
| `type Strategy interface{ Apply(int) int }` + `XStrategyImpl` + `StrategyFactory` | `type apply func(int) int` and plain functions/closures | Go proverb "The bigger the interface, the weaker the abstraction"; commonly recommended |
| `func New(a, b, c, d, e string, f, g int, h bool)` | An options struct, or functional options (35) | Google Go Style Guide, *Option structure* |
| `func Do(retry bool)` called as `Do(true)` | Two functions, or a named type/option | Commonly recommended: boolean parameters hide meaning at call sites |
| Naked `return` in a 40-line function | Explicit `return v, err` | Go Code Review Comments, *Naked Returns* |
| `(int, int, error)` with no names, meanings unclear | `(lo, hi int, err error)` when the names add information | Go Code Review Comments, *Named Result Parameters* |
| `f, err := os.Open(p); … ; f.Close()` at the end | `defer f.Close()` right after the error check | Effective Go, *Defer* |
| `log.Fatal(err)` deep in a library | Return the error to the caller | Go Code Review Comments, *Don't Panic*; §7.5 |
| `doStuff(ctx)` with `ctx` as the 3rd parameter | `ctx context.Context` first | Go Code Review Comments, *Contexts*; `context` package docs |

### 8.1 Return values, not out-parameters

```go
// ❌ the C/Java habit: write through a pointer, report success with a bool
func parsePrice(s string, result *int64) bool {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v < 0 {
		return false
	}
	*result = v
	return true
}

// ✅ the value and the reason it failed
func parsePrice(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v < 0 {
		return 0, fmt.Errorf("invalid price %q", s)
	}
	return v, nil
}
```
Effective Go (*Multiple return values*) presents multiple results as the fix for "in-band error returns" and "a pointer to a return value". Callers can't confuse which argument is written, and the error explains itself. Out-parameters do have one legitimate use, decoding *into* caller-owned memory to avoid allocation (`json.Unmarshal(data, &v)`, `io.Reader.Read(p []byte)`), and there the API is built around it.

### 8.2 Errors are returned, not thrown

```go
// ❌ exceptions in disguise
func run(args []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("checkout: %v", r)
		}
	}()
	if len(args) == 0 {
		panic("no items")
	}
	// …
}

// ✅
func run(args []string) error {
	if len(args) == 0 {
		return errors.New("no items")
	}
	// …
}
```
Go Code Review Comments (*Don't Panic*): "Don't use panic for normal error handling. Use error and multiple return values." `panic` is for bugs and impossible states. `recover` belongs at a few boundaries (an HTTP server recovering a handler's panic, 34) and is topic 18. E4 is this refactor.

### 8.3 A function type instead of a one-method interface and a factory

```go
// ❌ Java's strategy pattern, transliterated
type DiscountStrategy interface{ CalculateDiscount(prices []int64) int64 }
type PercentDiscountStrategyImpl struct{ percent int64 }
func (s *PercentDiscountStrategyImpl) CalculateDiscount(p []int64) int64 { … }
type DiscountStrategyFactory struct{}
func (f *DiscountStrategyFactory) CreateStrategy(code string) DiscountStrategy { … }

// ✅ the strategy is a function; closures hold its parameters
type discount func(prices []int64) int64

func percentOff(percent int64) discount {
	return func(prices []int64) int64 { return subtotal(prices) * percent / 100 }
}

func discountFor(code string) (discount, error) {
	switch code {
	case "SAVE10":
		return percentOff(10), nil
	// …
	default:
		return nil, fmt.Errorf("unknown coupon %q", code)
	}
}
```
When the "strategy" has no state worth naming and one operation, it's a function. Interfaces earn their place when a type has several related methods, or when a *consumer* needs to accept many implementations it doesn't know about (12, 17). Both forms are idiomatic in the standard library, and they combine: `http.Handler` is an interface, and `http.HandlerFunc` adapts a plain function to it.

### 8.4 Plain loops over functional chains

```go
// ❌ streams-style helpers stacked up (from Java streams / Python comprehensions)
total := reduce(mapInts(filterInts(orders, isPaid), amount), add, 0)

// ✅
total := 0
for _, o := range orders {
	if o.Paid {
		total += o.Amount
	}
}
```
The loop is easier to read, debug (one place for a breakpoint) and change, and it allocates nothing. Higher-order functions are idiomatic where they **parameterise behaviour**: a comparison for `slices.SortFunc`, a handler, a retry policy (E3's `Retry`), a visitor (E3's `Walk`).

### 8.5 `defer` right after acquisition; small `main`

```go
func run() error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close() // read-only: ignoring Close's error is fine

	mu.Lock()
	defer mu.Unlock() // 24
	// …
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1) // the only exit; every defer in run has already run
	}
}
```

## 9. Coming from Python / Java

| Concept | Python (FastAPI) | Java (Spring Boot) | Go | What doesn't translate |
|---|---|---|---|---|
| Declaring | `def f(a: int) -> int:` | `int f(int a)` in a class | `func f(a int) int` (package-level) | No methods required: plain functions are normal |
| Optional args | `def f(a, b=2, *, c=3)` | Overloads, builders | Options struct / functional options (35) | No defaults, no keywords |
| Varargs | `*args`, `**kwargs` | `int... xs` | `xs ...int` | No `**kwargs`; `f(s...)` aliases |
| Several results | `return a, b` (a tuple) | Return a `record` | `return a, b` (two values, no tuple) | Can't store "the result" in one variable |
| Errors | `raise` / `try` | `throw` / `try` | `return …, err` / `if err != nil` | No stack unwinding for errors |
| Lambdas | `lambda x: x * 2` | `x -> x * 2` | `func(x int) int { return x * 2 }` | Full bodies; verbose but explicit types |
| Closures | Shared cells; `nonlocal` to assign | Effectively-final copies | Shared variables, assignable | Java code relying on copies behaves differently |
| Function types | `Callable[[int], int]` | `Function<Integer, Integer>` | `func(int) int` | No boxing; nil is the zero value |
| Decorators | `@retry(3)` | AOP annotations | Wrapper functions: `withRetry(f)`, `func(http.Handler) http.Handler` | Explicit wrapping, no magic (35) |
| Generators | `yield` | `Iterator`/`Stream` | Closures, or range-over-func iterators (15) | |
| Cleanup | `with open(p) as f:` | try-with-resources | `defer f.Close()` | Function-scoped, not block-scoped |
| `finally` | `finally:` | `finally {}` | `defer` (runs on return and on panic) | |
| Static/utility classes | module functions | `final class Utils { static … }` | Package-level functions | No `Utils` class |
| Recursion limit | `RecursionError` at ~1000 | `StackOverflowError` | Growable stack; fatal at 1 GB | Can't catch it |

**Habits to unlearn here (AGENTS §1.1):**
- **Exceptions** → return `error`; check it right away; wrap with `%w`. `panic`/`recover` are not try/catch (§8.2, E4).
- **Out-parameters and "result objects"** → multiple results (§8.1, E4).
- **Strategy interfaces, `…Impl` types and factories** for a single behaviour → a func type and closures (§8.3, E4).
- **Utility classes / static helper classes** → package-level functions in a well-named package (16).
- **Decorators and AOP** → explicit wrapper functions, visible at the call site (35).
- **Streams and comprehensions everywhere** → a plain `for` loop (§8.4).

## 10. Production Lens

- **`defer` is cheap now: measure before removing it.** Real benchmark on go1.27.1, Apple M4 Pro, functions marked `//go:noinline` so the call isn't optimised away:
  ```
  BenchmarkWithDefer-12         	342359823	         3.489 ns/op	       0 B/op	       0 allocs/op
  BenchmarkWithoutDefer-12      	583545004	         2.050 ns/op	       0 B/op	       0 allocs/op
  BenchmarkDeferInLoop10-12     	15265904	        78.78 ns/op	       0 B/op	       0 allocs/op
  BenchmarkDirectInLoop10-12    	251183089	         4.948 ns/op	       0 B/op	       0 allocs/op
  ```
  `mu.Lock(); defer mu.Unlock()` costs about 1.4 ns more than unlocking by hand. That's open-coded defer (§5.7), and it's never worth the risk of a missed unlock on an early return or a panic. Ten `defer`s **in a loop** cost about 75 ns more than ten direct calls (about 7.5 ns each, on the slower defer-record path): one more reason to move loop bodies into functions. Writing benchmarks: topic 19.
- **Closures that escape allocate.** `go build -gcflags=-m` shows `moved to heap: count` and `func literal escapes to heap` for a returned closure (§5.5). In a hot path (per request, per item), returning closures can show up in allocation profiles (29). A closure used locally and called directly is usually inlined (`can inline SumWith.func1`, `inlining call to SumWith.func1` in the same scratch build) and costs nothing.
- **Inlining.** Small functions (under the compiler's cost budget, 80 nodes) are inlined. `-gcflags=-m` prints `can inline f`. Calls through a func value or an interface are usually *not* inlined, because the target isn't known at compile time (the compiler can sometimes devirtualise, and PGO helps; 29). Don't contort code for this; measure first.
- **API design.** Context first, error last, options as a struct, no boolean flags, no out-parameters. Accept small function types for callbacks (`func(attempt int) error`). Name the type when it appears in several signatures (`type Opener func(string) (io.ReadCloser, error)` in E2). Document whether a callback may be nil, may be called concurrently, and when it stops being called.
- **Errors at boundaries.** Wrap with context once per layer (`fmt.Errorf("load user %q: %w", id, err)`), log once, at the top (18, 43). The deferred-annotation pattern (§5.7) keeps that consistent in functions with many returns.
- **Testing.** Functions with injected dependencies as func values (an `Opener`, a clock `func() time.Time`, a `sleep func(time.Duration)`) are trivial to fake in tests, as E2 does: no mocking framework needed (20). Test the error paths and that cleanup happened (E2 checks no file is left open).
- **Concurrency.** A closure passed to `go` must not share mutable captured variables without synchronisation (§7.2). Callbacks you accept may be invoked from other goroutines: say so in the doc comment, and protect their state. `defer mu.Unlock()` is the default unlock style (24).
- **Shutdown.** Deferred cleanups run only on normal return or panic, not on `os.Exit` or a `SIGKILL`. Flushing logs, closing DB pools and finishing in-flight requests belong in a graceful shutdown path started from `main` via `run` (38).

### Senior review: a report exporter

```go
// Export writes one CSV per customer and returns how many were written.
func Export(customers []Customer, dir string, verbose bool, skipEmpty bool, maxRows int) int {
	written := 0
	for _, c := range customers {
		f, err := os.Create(filepath.Join(dir, c.ID+".csv"))
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		rows, err := loadRows(c.ID)
		defer log.Printf("exported %s: %d rows, err=%v", c.ID, len(rows), err)
		if err != nil {
			panic(err)
		}
		if skipEmpty == true && len(rows) == 0 {
			continue
		}
		for i, r := range rows {
			if i >= maxRows {
				break
			}
			f.WriteString(r.CSV() + "\n")
		}
		written++
	}
	return written
}
```

Review comments a senior Go engineer would leave:
1. **`defer f.Close()` in the loop** (§7.4): every file stays open until `Export` returns. With 50 000 customers that's "too many open files". Extract `exportOne(dir string, c Customer, opts Options) error` and defer inside it.
2. **`defer f.Close()` before the error check** (§7.7): if `os.Create` fails, `f` is nil. `(*os.File).Close` happens to return `ErrInvalid` on a nil file instead of panicking, but the pattern panics for most types. Check `err` first.
3. **`log.Fatal` in a library function** (§7.5) kills the whole service and skips every deferred call, including the `Close`s. Return an error.
4. **`panic(err)` for a failed load** (§8.2): a database hiccup isn't a programmer bug. Return `fmt.Errorf("export %s: %w", c.ID, err)`.
5. **`defer log.Printf(…, len(rows), err)`**: the arguments are evaluated at the `defer` (§7.3), so it logs the values from *before* anything happened, and only when the whole export finishes, once per customer, all at the end. Log after the work, once.
6. **Write and Close errors are ignored** (§7.6, §7.15): a full disk produces silently truncated CSVs, and the function reports them as `written`. Check `WriteString`'s error, and capture `Close`'s error with a named result in `exportOne`. errcheck flags both.
7. **The parameter list.** `Export(cs, dir, true, false, 100)` is unreadable at the call site: two booleans and a magic number. Use an `ExportOptions{SkipEmpty: true, MaxRows: 100}` struct. `verbose` is unused: delete it (logging is the caller's choice, 43).
8. **The signature.** It should be `func Export(ctx context.Context, customers []Customer, dir string, opts ExportOptions) (int, error)`: context first (26) so a long export can be cancelled, and an error last so callers know it failed.
9. **`skipEmpty == true`** → `opts.SkipEmpty`. Building each line with `r.CSV() + "\n"` allocates. `encoding/csv` handles quoting correctly (32).

## 11. Interview Lens

| Question | What a strong answer shows | Red flags |
|---|---|---|
| Does Go pass slices and maps by reference? | Everything is by value; a slice header/map handle is copied; element writes visible, `append` not; `f(s...)` passes the same slice | "Yes, by reference" |
| What does this print? `x := 1; defer fmt.Println(x); x = 2` | `1`: arguments evaluated at the `defer`; a deferred closure would print 2 | Says 2 |
| In what order do deferred calls run, and when? | LIFO, after `return` sets the results, before the caller resumes; also during a panic; not on `os.Exit` | "At the end of the block" |
| How can a deferred function change a return value? | Named results; `return` assigns them first; the classic error-annotation and `Close`-error patterns | Thinks it works with unnamed results |
| What's a closure? What does it capture? | Function literal + the variables it references, shared and kept alive; heap if it escapes; per-iteration loop variables since 1.22 | "A copy of the values" |
| What happened to the loop-variable/closure bug in Go 1.22? | Per-iteration variables, gated by the `go` line; `v := v` obsolete | Unaware of the change |
| Is `defer` expensive? | Open-coded since 1.14 (~1–2 ns); defers in loops are slower; measured, not guessed | "Avoid defer for performance" as a blanket rule |
| Why no exceptions in Go? When is `panic` OK? | Errors are values, explicit control flow; `panic` for bugs/impossible states; `recover` at boundaries (18) | Uses `panic` for validation |
| Why doesn't Go have default parameters/overloading? | Simplicity, one name = one signature; options structs, functional options | — |
| What happens on infinite recursion? | Stack grows to 1 GB, then a fatal, unrecoverable error; no TCO | "StackOverflowError you can catch" |

**What separates a senior answer:** connecting the mechanics to failure modes and tooling, e.g. "I check `err` before deferring `Close`, capture `Close` errors on writes with a named result, keep `os.Exit` in a three-line `main`, and let errcheck and vet run in CI. When a closure shows up in an allocation profile, `-gcflags=-m` tells me why."

## 12. Exercises

[`exercises/`](../exercises/README.md) · run from this topic folder (`make` lists every command). Don't open `solutions/` until you've tried.

| # | Exercise | Trains | Command | Status | Score |
|---|---|---|---|---|---|
| E1 | [Predict: defer & closures](../exercises/01-predict-defer/) | `defer` order and argument evaluation, deferred closures and named results, closure state, variadic aliasing, `append` in a callee, composing funcs (§5.1–5.7) | `make ex E=01-predict-defer` | ⬜ | – |
| E2 | [Fix the bug: line counter](../exercises/02-fix-linecount/) | `defer` in a loop holds every file open (§5.7, §7.4) | `make ex E=02-fix-linecount` | ⬜ | – |
| E3 | [Implement: funcs](../exercises/03-implement-funcs/) | Multiple results + sentinel error, variadics, returning closures, shared closure state, a higher-order retry with `%w`, recursion (§5.2–5.8) | `make ex E=03-implement-funcs` | ⬜ | – |
| E4 | [Refactor: Java-style checkout](../exercises/04-refactor-checkout/) | Strategy interface + `Impl` types + factory → func values; `panic`/`recover` → errors; out-parameter → `(value, error)` (§8) | `make ex E=04-refactor-checkout` | ⬜ | – |
| E5 | [Mini-project: `calc`](../exercises/05-project-calc/) | A recursive-descent expression evaluator (`cmd/calc` + `internal/expr`): mutually recursive closures, positioned errors, a sentinel error, a `defer`-based trace | `make ex E=05-project-calc` | ⬜ | – |

Before asking the tutor to `check En`: `make check E=<folder>` (gofmt + vet + race tests).

Exercise status: ⬜ not started · 🛠️ attempted · ✅ passed + reviewed · 🔁 redo (solved after seeing the full solution)

## 13. Cheat Sheet

1. **Everything is passed by value.** Pointers, slices, maps, channels and funcs are small values that *refer to* shared data. `f(s...)` passes `s` itself.
2. **No overloading, defaults or keyword args:** different names, an options struct, variadics, functional options.
3. **`(T, error)`, error last;** zero values on failure; check right away; wrap with `fmt.Errorf("…: %w", err)`; `errors.Is` for sentinels.
4. **Named results** for documentation or deferred modification; bare `return` only in short functions.
5. **Func values:** typed, nil by default (calling nil panics), comparable only to nil.
6. **Closures capture variables, not values;** shared, kept alive, on the heap if they escape (`-gcflags=-m`). Per-iteration loop variables since 1.22.
7. **Recursive literal:** `var f func(…)…; f = func…`. No tail calls; stacks grow to 1 GB, then a fatal error.
8. **`defer`: arguments evaluated now, call runs at return, LIFO,** also on panic, never on `os.Exit`/`log.Fatal`. Deferred closures read variables at run time.
9. **`defer` right after the error check;** function-scoped, so loop bodies with `defer` become functions. Capture `Close` errors on writes with a named result.
10. **Habits:** no out-params, no `panic` for errors, no strategy interfaces/factories for one function, plain loops over stream chains, tiny `main` + `run() error`.

## 14. References & Further Reading

- **The Go Programming Language Specification**: *Function types*, *Function declarations*, *Function literals*, *Calls* (*Passing arguments to ... parameters*), *Return statements*, *Defer statements*, *Terminating statements*, *Handling panics*: go.dev/ref/spec
- **Effective Go**: *Functions* (*Multiple return values*, *Named result parameters*, *Defer*), *Errors*, *Panic*, *Recover*: go.dev/doc/effective_go
- **Go Code Review Comments**: *Contexts*, *Don't Panic*, *Error Strings*, *Handle Errors*, *Named Result Parameters*, *Naked Returns*: go.dev/wiki/CodeReviewComments
- **Go FAQ**: *Why does Go not support overloading of methods and operators?*, *Why does Go not have exceptions?*: go.dev/doc/faq
- **Go blog**: Andrew Gerrand, *Defer, Panic, and Recover* (2010): go.dev/blog/defer-panic-and-recover · Rob Pike, *Errors are values* (2015): go.dev/blog/errors-are-values · Andrew Gerrand, *Error handling and Go* (2011): go.dev/blog/error-handling-and-go
- **Release notes**: Go 1.14 (open-coded defers): go.dev/doc/go1.14 · Go 1.17 (register-based calling convention): go.dev/doc/go1.17 · Go 1.20 (`errors.Join`): go.dev/doc/go1.20 · Go 1.22 (per-iteration loop variables): go.dev/doc/go1.22
- **Google Go Style Guide**, best practices: *Function argument lists* (*Option structure*, *Variadic options*): google.github.io/styleguide/go
- Rob Pike, **Go Proverbs** (2015): go-proverbs.github.io
- Teiva Harsanyi, **100 Go Mistakes and How to Avoid Them** (Manning, 2022): #43 *Never using named result parameters*, #44 *Unintended side effects with named result parameters*, #45 *Returning a nil receiver*, #46 *Using a filename as a function input*, #47 *Ignoring how defer arguments and receivers are evaluated*, #54 *Not handling defer errors*
- Package docs and tools: `go doc errors`, `go doc fmt.Errorf`, `go build -gcflags=-m` (escape analysis and inlining), `go vet` (`unusedresult`), errcheck and staticcheck (SA5007) via golangci-lint

---

## 🗣️ Q&A Discussion

Sessions, exercise reviews and stuck entries are logged in [`functions-closures-and-defer-qa.md`](functions-closures-and-defer-qa.md).
