package helpers

import "strings"

func CleanInput(text string) []string {
	input := text

	input = strings.TrimSpace(text)
	input = strings.ToLower(input)

	if input == "" {
		return []string{}
	}

	separatedInputs := strings.Split(input, " ")

	return separatedInputs
}
