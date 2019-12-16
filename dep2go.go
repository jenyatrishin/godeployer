package main

import (
	"./config"
	"./factory/deployerFactory"
	"./tools"
	"flag"
	"fmt"
	"strings"
)

const FILENAME string = "config"
const SSH_POST string = "22"
const VERSION string = "0.3.1-alpha"
const DEVELOPER string = "developer"
const STAGING string = "staging"
const PRODUCTION string = "production"

func run (f func() string) bool {
	f()
	return true
}

func createConfig () string {
	ext := getExtFromCommand()

	configIns := new(config.Config)

	configIns.WriteConfig(FILENAME+"."+ext, ext)
	fmt.Println("Config file created")
	return "OK"
}

func validateConfig () string {
	fmt.Println("config file is validating...")
	ext := getExtFromCommand()
	message := "Config file not valid"

	configIns := new(config.Config)

	if configIns.ValidateConfig(FILENAME+"."+ext, ext) {
		message = "Config file valid "
	}
	fmt.Println(message)
	return "validate config"
}

func deploy (mode string) string {
	ext := getExtFromCommand()

	configIns := new(config.Config)

	configIns.ReadConfig(FILENAME+"."+ext, ext)

	deployer := deployerFactory.GetDeployer()
	currentEnv := configIns.GetEnvByType(mode)

	if currentEnv.EnvType == "" {
		err := tools.UserError("That environment is not defined")
		fmt.Println(err)
		return ""
	}

	authMethod := currentEnv.AuthType

	var pass string
	if authMethod == "key" {
		pass = currentEnv.KeyFile
	} else {
		pass =  currentEnv.Password
	}

	deployer.PrepareConfig(currentEnv.Server, SSH_POST, currentEnv.Login, pass, authMethod)
	deployer.DeployTo(
		currentEnv.HomeDir,
		currentEnv.GitConfig.Repository,
		currentEnv.GitConfig.User,
		currentEnv.GitConfig.Password,
		currentEnv.GitConfig.Branch,
		getCommandString(currentEnv.AfterDeploy),
		getCommandString(currentEnv.BeforeDeploy),
		)

	return "deploy finished"
}

func getCommandString (afterDeployCommands []config.Command) string {
	output := ""
	for i := 0; i < len(afterDeployCommands); i++ {
		output += afterDeployCommands[i].Item
		if i != len(afterDeployCommands)-1 {
			output += "&&"
		}
	}

	return output
}

func deployToDev () string {
	return deploy(DEVELOPER)
}

func deployToStaging () string {
	return deploy(STAGING)
}

func deployToProduction () string {
	return deploy(PRODUCTION)
}

func getVersion () string {
	fmt.Println(VERSION)
	return VERSION
}

func main () {
	commands := getCommands()

	flag.Parse()

	command := flag.Args()

	if len(command) > 0 {
		s := strings.Split(command[0], ":")
		if len(s) < 2 {
			fmt.Println(tools.UserError("You set bad command"))
			return
		}
		closure := tools.CommandIsAllowed(s, commands)
		run(closure)
	} else {
		fmt.Println("Dep2go deployment tool")
		fmt.Println("Usage: dep2go [--command]")
		getCommandsToString()
		return
	}
}

func getCommands () map[string]map[string]func()string {
	commands := map[string]map[string]func()string{
		"config": {
			"make": createConfig,
			"validate": validateConfig,
		},
		"deploy": {
			"developer": deployToDev,
			"staging": deployToStaging,
			"production": deployToProduction,
		},
		"version": {
			"show": getVersion,
		},
	}

	return commands
}

func getCommandsToString () {
	commands := getCommands()
	fmt.Println("Commands that can be used:")
	for k, v := range commands {
		for q, _ := range v {
			fmt.Println(k + ":" + q)
		}
	}
}

func getExtFromCommand () string {
	ext := "xml"
	firstOpt := flag.Arg(1)
	if firstOpt != "" {
		if strings.Contains(firstOpt, "-ext") {
			data := strings.Split(firstOpt, "=")
			if len(data) == 2 && (data[1] == "json" || data[1] == "xml"){
				ext = data[1]
			}
		}
	}

	return ext
}
