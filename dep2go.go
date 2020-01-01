package main

import (
	"./app"
	"./config"
	"./factory/deployerFactory"
	"./tools"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	FOLDERNAME string = ".dep2go"
	FILENAME   string = ".dep2go/config"
	SSH_POST   string = "22"
	VERSION    string = "0.4.1-alpha"
	DEVELOPER  string = "developer"
	STAGING    string = "staging"
	PRODUCTION string = "production"
)

func deploy(mode string) {
	ext := getConfigFileExtForDeploy()

	configIns := new(config.Config)

	configIns.ReadConfig(FILENAME+"."+ext, ext)

	deployer := deployerFactory.GetDeployer()
	currentEnv := configIns.GetEnvByType(mode)

	if currentEnv.EnvType == "" {
		errMessage := "That environment is not defined"
		tools.WriteLog(errMessage)
		err := tools.UserError(errMessage)
		fmt.Println(err)
		return
	}
	tools.WriteLog("Started deploy to env: " + currentEnv.EnvType)
	authMethod := currentEnv.AuthType

	var pass string
	if authMethod == "key" {
		pass = currentEnv.KeyFile
	} else {
		pass = currentEnv.Password
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
	tools.WriteLog("Deploy command is finished for env: " + currentEnv.EnvType)
//	return "deploy finished"
}

func getCommandString(afterDeployCommands []config.Command) string {
	output := ""
	for i := 0; i < len(afterDeployCommands); i++ {
		output += afterDeployCommands[i].Item
		if i != len(afterDeployCommands)-1 {
			output += "&&"
		}
	}

	return output
}

func deployToDev() {
	deploy(DEVELOPER)
}

func deployToStaging() {
	deploy(STAGING)
}

func deployToProduction() {
	deploy(PRODUCTION)
}

func main() {
	app := &app.Dep2Go{}

	app.AddCommand("config:init",
		"Make init action",
		func() {
			err := os.Mkdir(FOLDERNAME, os.ModePerm)
			if err != nil {
				tools.WriteLog("Init error: " + err.Error())
			}
		},
	)
	app.AddCommand("config:make",
		"Create config file from command line",
		func() {
			ext := getExtFromCommand()

			configIns := new(config.Config)

			configIns.WriteConfig(FILENAME+"."+ext, ext)
			fmt.Println("Config file created")
		},
	)
	app.AddCommand("config:validate",
		"Validate config file",
		func() {
			fmt.Println("config file is validating...")
			ext := getConfigFileExtForDeploy()
			message := "Config file not valid"

			configIns := new(config.Config)

			if configIns.ValidateConfig(FILENAME+"."+ext, ext) {
				message = "Config file valid"
			}
			fmt.Println(message)
		},
	)
	app.AddCommand("deploy:developer",
		"Deploy project to developer environment",
		deployToDev,
	)

	app.AddCommand("deploy:staging",
		"Deploy project to staging environment",
		deployToStaging,
	)

	app.AddCommand("deploy:production",
		"Deploy project to production environment",
		deployToProduction,
	)

	app.AddCommand("version:show",
		"Show dep2go tool current version",
		func() {
			fmt.Println(VERSION)
		},
	)

	app.Run()
}

func getExtFromCommand() string {
	ext := "xml"
	firstOpt := flag.Arg(1)
	if firstOpt != "" {
		if strings.Contains(firstOpt, "-format") {
			data := strings.Split(firstOpt, "=")
			if len(data) == 2 && (data[1] == "json" || data[1] == "xml") {
				ext = data[1]
			}
		}
	}

	return ext
}

func getConfigFileExtForDeploy() string {
	ext := "json"
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(dir + "/" + FILENAME + ".json"); os.IsNotExist(err) {
		ext = "xml"
	}

	return ext
}
