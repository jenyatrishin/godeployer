package tools

import (
	//"../config"
	//adapterXml "../adapter/xml"
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

func Difference(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m :=map [string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal==1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}

