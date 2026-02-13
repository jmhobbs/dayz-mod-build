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
	trimmed := strings.ToLower(strings.TrimSpace(text))
	return trimmed == "y" || trimmed == "yes", nil
}
