package main

import "fmt"

func main() {
	var n1 = 10
	var n2 int = 20
	var f1 float32 = 20
	var f2 float64 = 20.021
	var s1 = "Hello"
	const s2 string = "Bye"
	var b1 bool = true
	fmt.Println(n1, n2, f1, f2, s1, s2, b1)
	sn1 := 10 // Short variable declaration
	sn2 := 1212.212
	ss1 := "Hello"
	fmt.Println(sn1, sn2, ss1)
	fmt.Printf("%s %s %f %.2f", s1, s2, f1, f2)
	fmt.Printf("\n%v %.d %T", n1, n2, s1)
}
