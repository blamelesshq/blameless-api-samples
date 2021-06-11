package utils

import (
	"fmt"
	"log"
	"strconv"

	"github.com/manifoldco/promptui"
)

// IntPrompt provides you a simple interface to execute integer prompt request
func IntPrompt(label string) int {
	validateInt := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return fmt.Errorf("unable to parse integer for %s", label)
		}
		return nil
	}

	intP := promptui.Prompt{
		Label:    label,
		Validate: validateInt,
	}
	result, err := intP.Run()
	if err != nil {
		log.Fatal(err)
	}

	id, err := strconv.Atoi(result)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

// StringPrompt provides you a simple interface to execute string prompt request
func StringPrompt(label string) string {
	validateString := func(input string) error {
		if len(input) < 1 {
			return fmt.Errorf("must provide at least one character for %s", label)
		}
		return nil
	}
	stringP := promptui.Prompt{
		Label:    label,
		Validate: validateString,
	}
	result, err := stringP.Run()
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// BooleanPrompt provides you a simple interface to execute boolean prompt request
func BooleanPrompt(label string) bool {
	validateBool := func(input string) error {
		if input != "true" || input != "false" {
			return fmt.Errorf("must provide either true or false for %s", label)
		}
		return nil
	}
	boolP := promptui.Prompt{
		Label:    label,
		Validate: validateBool,
	}
	result, err := boolP.Run()
	if err != nil {
		log.Fatal(err)
	}
	bl, err := strconv.ParseBool(result)
	if err != nil {
		log.Fatal(err)
	}
	return bl
}
