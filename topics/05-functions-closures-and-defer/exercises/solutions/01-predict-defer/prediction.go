package main

// prediction is the real output, with the reason for each line:
//
//  1. Deferred calls run when the func literal returns, last in first out,
//     and each i was evaluated when its defer statement ran (§5.7).
//  2. Two defers, so LIFO: the closure runs first and reads s when it runs
//     ("b"); the plain call's argument was evaluated at the defer statement
//     ("a"). The closure prints without a newline, so both share a line.
//  3. named: return n + 1 sets the named result n to 2, then the deferred
//     closure multiplies it: 20. unnamed: the result is a separate, unnamed
//     slot set to 2 by return; the closure changes the local n, which is no
//     longer the result: 2 (§5.7).
//  4. Each makeCounter call creates a new c; a was called 3 times, b twice
//     (§5.5).
//  5. The closure captured the variable x, not its value at creation (§5.5).
//  6. zeroFirst(nums...) passes nums itself, so nums[0] becomes 0;
//     zeroFirst(nums[1], nums[2]) builds a new slice from copies of the
//     values, so nothing in nums changes (§5.3).
//  7. grow gets a copy of the header with cap 1; append allocates a new
//     array for the local s, and s[0] = 9 writes that new array. The
//     caller's s still sees [1] (§5.1, topic 06).
//  8. apply runs the functions left to right: (3+1)*2 = 8, 3*2+1 = 7, and
//     no functions leave x unchanged (§5.4).
const prediction = `loop;2;1;0;
closure:b arg: a
20 2
3 2
5
[0 2 3]
[1]
8 7 3
`
