package main

import (
	"./adapter"
	"./config"
	"./factory/xmlAdapterFactory"
	"./factory/deployerFactory"
	"./tools"
	"flag"
	"fmt"
	"strings"
)

const FILENAME string = "config"
const SSH_POST string = "22"
const VERSION string = "0.2-alpha"

func run (f func() string) bool {
	f()
	return true
}

func createConfig () string {
	//should be taken from command line
	ext := "xml"
	adapterIns := GetAdapterByType(ext)
	configIns := new(config.Config)
	//deprecated - gonna be removed
	configIns.ConfigFileType = ext
	configIns.WriteConfig(FILENAME+"."+ext,adapterIns)
	fmt.Println("Config file created")
	return "OK"
}

func validateConfig () string {
	fmt.Println("validate config")
	return "create config"
}

func deploy (mode string) string {
	//should be taken from command line
	ext := "xml"
	adapterIns := GetAdapterByType(ext)
	configIns := new(config.Config)
	//deprecated - gonna be removed on 0.3-alpha version
	configIns.ConfigFileType = ext
	configIns.ReadConfig(FILENAME+"."+ext, adapterIns)

	deployer := deployerFactory.GetDeployer()
	currentEnv := configIns.GetEnvByType(mode)
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
		getAfterCommandString(currentEnv.AfterDeploy),
		)

	return "read config"
}

func getAfterCommandString (afterDeployCommands []config.Command) string {
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
	deploy("developer")
	return "deployed to dev"
}

func deployToStaging () string {
	deploy("staging")
	return "deployed to staging"
}

func deployToProduction () string {
	deploy("production")
	return "deployed to prod"
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

func GetAdapterByType (configType string) adapter.ConfigAdapter {
	if configType == "xml" {
		return xmlAdapterFactory.GetXmlAdapter()
	}
	return xmlAdapterFactory.GetXmlAdapter()
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