package gostrings

import (
	"fmt"
	"unicode/utf8"
)

func GoStrings() {
	/*
		A string is an immutable sequence of bytes.

		This definition is the foundation for understanding Go strings.

		Let's break it down:
		- Sequence: Bytes are arranged in a specific order.
		- Bytes: Each byte represents an 8-bit value, from 0 to 255.
		- Immutable: Once a string value exists, you cannot modify its individual bytes through string indexing.
		- Not necessarily text: A string can contain arbitrary bytes, including bytes that do not form valid UTF-8 text.

		Each character in this ASCII string occupies one byte.
	*/
	fmt.Println("===== STRINGS =====")
	message := "We learning GO string."
	fmt.Println(message)
	fmt.Println("Message Length: ", len(message))

	first := "First String"
	second := first
	fmt.Println("SECOND", second)
	second = "Second String"
	fmt.Println(first)
	fmt.Println("SECOND", second)

	fmt.Println(len("Go"))
	fmt.Println(len("café"))
	fmt.Println(len("नमस्ते"))
	fmt.Println(len("😊"))
	fmt.Println(len("A"))
	fmt.Println(len("é"))
	fmt.Println(len("न"))
	fmt.Println(len("😊"))
	s := "é"
	fmt.Println(len(s))
	fmt.Printf("% x\n", []byte(s))
	s1 := string([]byte{0xff, 0xfe})

	fmt.Println(utf8.ValidString(s1)) // false
	fmt.Println(message[1])
	fmt.Printf("Decoded string: %c\n", message[1])

	// Verifying string
	s1v := "Hello, How are you?"
	s2v := string([]byte{0xff})
	r1v, size1v := utf8.DecodeRuneInString(s1v)
	r2v, size2v := utf8.DecodeRuneInString(s2v)
	fmt.Printf("%U, %U\n", r1v, r2v)
	fmt.Println(size1v, size2v)
	s11s := "नमस्ते"

	for i, r := range s11s {
		fmt.Printf("index=%d rune=%c codepoint=%U\n", i, r, r)
	}
	s12s := "A😊B"

	for i, r := range s12s {
		fmt.Printf("%d %c %U\n", i, r, r)
	}
}
