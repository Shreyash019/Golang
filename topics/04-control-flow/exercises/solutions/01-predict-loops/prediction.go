package main

// prediction is the real output, with the reason for each line:
//
//  1. range assigns the next value to i at the start of every iteration, so
//     i *= 10 never changes how many iterations run (§5.3).
//  2. A three-clause loop copies i into the next iteration before i++ runs,
//     so the body's i *= 2 carries over: 0 → 1, 1 → 2 → 3, 3 → 6 → 7,
//     7 → 14 → 15 stops (§5.3).
//  3. continue skips the print for 1; break leaves only the switch, so 3 is
//     still printed (§7.1); fallthrough runs the next body for 5 without
//     testing i > 10 (§5.7).
//  4. continue outer abandons the inner loop as soon as j > i (§5.8).
//  5. range over an array iterates a copy made before the loop: the copy
//     still holds 3, so 1+2+3 + 0+1+2 = 9. Range over a slice copies only the
//     header and shares the elements: 1+2+10 + 0+1+2 = 16 (§5.4).
//  6. The range expression is evaluated once, so the loop runs three times
//     even though nums grows; n is a copy, so nums[0..2] are unchanged (§5.4).
//  7. Ranging over a string yields byte offsets: ñ is two bytes in UTF-8, so
//     b starts at offset 3 and len is 4 (§5.5).
//  8. Since Go 1.22 each iteration has its own i, so the closures return
//     0, 1, 4. With go 1.21 semantics they would all return 9 (§5.3).
const prediction = `0;10;20;30;
0;1;3;7;
0;2;3;4;big;5;
11;21;22;31;32;33;
9 16
[1 2 3 2 4 6]
0a;1ñ;3b;4
0;1;4;
`
