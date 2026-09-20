package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	line1, _ := reader.ReadString('\n')
	line1 = strings.TrimSpace(line1)

	line2, _ := reader.ReadString('\n')
	line2 = strings.TrimSpace(line2)

	line3, _ := reader.ReadString('\n')
	line3 = strings.TrimSpace(line3)

	a, err := strconv.Atoi(line1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	b, err := strconv.Atoi(line2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	switch line3 {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(a / b)
	default:
		fmt.Println("Invalid operation")
	}
}
