package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Input and Output in Golang")
	var name string
	fmt.Print("Enter your name: ")
	// fmt.Scan(&name)
	// fmt.Scanln(&name)
	scanner := bufio.NewReader(os.Stdin)
	name, err := scanner.ReadString('\n')
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Hello", name)
	// if len(os.Args) < 2 {
	// 	fmt.Println("Please provide a greeting message.")
	// 	return
	// }
	// greeting := os.Args[1]
	// fmt.Printf("%s, world!\n", greeting)
}
