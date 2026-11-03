package main

// prediction is the real output, line by line:
//
//   - -126: int8 120+10 = 130 wraps around: 130 - 256 = -126.
//   - -2 -1 -3: / truncates toward zero (-2.25 → -2), % takes the sign of
//     the dividend (-9 = -2*4 + -1), but >> is an arithmetic shift, which
//     rounds toward minus infinity (-2.25 → -3).
//   - 10 6 2: << binds like *, so 1<<3+2 is (1<<3)+2; ^ is XOR
//     (0010 ^ 0100 = 0110); &^ clears bits (0111 &^ 0101 = 0010).
//   - 3 -3: float-to-int conversion truncates toward zero; it never rounds.
//   - 232 -24: 1000 is 0x3E8; converting keeps the low byte 0xE8, which is
//     232 unsigned and 232-256 = -24 signed.
//   - 8 9: in 9/2*2.0 the integer constants divide first (4), then 4*2.0;
//     9/2.0 is a float division (4.5).
//   - 6 195 n: "naïve" is 6 bytes (ï takes two); s[2] is the first byte of
//     ï (0xC3 = 195); s[0] is 'n'.
//   - true: uint 3-5 wraps around to 2^64-2, which is > 0.
//   - true: 0.5, 0.25 and 0.75 are sums of powers of two, so binary floats
//     represent them exactly. Only decimals like 0.1 are inexact.
const prediction = `-126
-2 -1 -3
10 6 2
3 -3
232 -24
8 9
6 195 n
true
true
`
