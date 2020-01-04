package app

import (
	"../tools"
	"flag"
	"fmt"
	"gopkg.in/gookit/color.v1"
	"strings"
)

const (
	SPACES   int = 20
	VERSION  string = "0.4.2-alpha"
)

type Dep2Go struct {
	commands map[string]map[string]func()
	commandsDesc map[string]map[string]string
}

func (app *Dep2Go) Run() {

	app.AddCommand(
		"help:show",
		"show help information",
		app.useBaseCommand,
		)
	app.AddCommand("version:show",
		"Show dep2go tool current version",
		func() {
			fmt.Println(VERSION)
		},
	)

	flag.Parse()

	command := flag.Args()

	if len(command) > 0 {
		s := strings.Split(command[0], ":")
		if len(s) < 2 {
			fmt.Println(tools.UserError("You set bad command"))
			return
		}
		closure := tools.CommandIsAllowed(s, app.commands)
		closure()
	} else {
		app.useBaseCommand()
		return
	}
}

func (app *Dep2Go) AddCommand(name string, desc string, callback func()) {
	nameParts := strings.Split(name,":")
	if len(nameParts) < 2 {
		color.Red.Println("You're using not correct name format -" + name)
		return
	}
	app.insertCommand(nameParts[0], nameParts[1], callback)
	app.insertDescription(nameParts[0], nameParts[1], desc)
}

func (app *Dep2Go) insertCommand(namePart0 string, namePart1 string, callback func()) {
	if len(app.commands) == 0 {
		app.commands = make(map[string]map[string]func())
	}

	var tempMap map[string]func()

	if app.commands[namePart0] != nil {
		app.commands[namePart0][namePart1] = callback
		tempMap = app.commands[namePart0]
	} else {
		tempMap = map[string]func(){
			namePart1: callback,
		}
	}
	app.commands[namePart0] = tempMap
}

func (app *Dep2Go) insertDescription(namePart0 string, namePart1 string, desc string) {
	if len(app.commandsDesc) == 0 {
		app.commandsDesc = make(map[string]map[string]string)
	}

	var tempMap map[string]string

	if app.commandsDesc[namePart0] != nil {
		app.commandsDesc[namePart0][namePart1] = desc
		tempMap = app.commandsDesc[namePart0]
	} else {
		tempMap = map[string]string{
			namePart1: desc,
		}
	}
	app.commandsDesc[namePart0] = tempMap
}

func (app *Dep2Go) useBaseCommand() {
	color.Green.Print("Dep2go deployment tool ")
	fmt.Print("version ")
	color.Cyan.Print(VERSION)
	fmt.Println(" (c) 2019")
	fmt.Println("")
	color.Cyan.Println("Usage:")
	fmt.Println("    dep2go [--command] [-options]")
	fmt.Println("")
	color.Cyan.Println("Available options:")
	spacesOptions := strings.Repeat(" ", SPACES - 7)
	color.Green.Println("    " + "-format" + spacesOptions + " - " + "set format for new config file")
	fmt.Println("")
	color.Cyan.Println("Available commands:")
	for k, v := range app.commands {
		fmt.Println(k)
		for q, _ := range v {
			desc := app.commandsDesc[k][q]
			spacesDiff := SPACES - len(k) - len(q) - 1
			spacesString := strings.Repeat(" ", spacesDiff)
			color.Green.Println("    "+k + ":" + q + spacesString + " - " + desc)
		}
	}
}