package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func yesOrNo(override bool, prompt string) (bool, error) {
	if override {
		fmt.Print(prompt)
		fmt.Println("Y (override)")
		return true, nil
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	return strings.ToLower(text) == "y\n" || strings.ToLower(text) == "yes\n", nil
}
