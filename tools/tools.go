package tools

import (
	"errors"
	"fmt"
)

func UserError(text string) error {
	return errors.New(text)
}

func CommandIsAllowed(command []string, commands map[string]map[string]func()) func() {

	if val,ok := commands[command[0]]; ok {
		if v,is := val[command[1]]; is {
			return v
		}
	}
	return wrongCommand
}

func wrongCommand() {
	fmt.Println(UserError("You set bad command"))
}

func RemoveEmptyElements(data []string) []string {
	output := make([]string,0)
	for _, item := range data {
		if item != "" {
			output = append(output, item)
		}
	}

	return output
}