package tools

import (
	"errors"
	"fmt"
)

func UserError(text string) error {
	return errors.New(text)
}

func CommandIsAllowed (command []string, commands map[string]map[string]func()string) func() string {

	if val,ok := commands[command[0]]; ok {
		if v,is := val[command[1]]; is {
			return v
		}
	}
	return wrongCommand
}

func wrongCommand () string{
	fmt.Println(UserError("You set bad command"))
	return ""
}