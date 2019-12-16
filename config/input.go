package config

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type Input struct {
	Conf Config
}

func (i *Input) callUserInput () {

	allowed := []string{"developer","staging","production"}
	var entered []string

	reader := bufio.NewReader(os.Stdin)

	var conf Config

CreateStep:
	var envIns Env

	fmt.Print("Enter environment name ("+ strings.Join(Difference(allowed, entered),",")+" or blank for exit): ")
	env, _ := reader.ReadString('\n')


	if prepareStringToSet(env) != "" {

		password := ""
		key := ""

		entered = append(entered, prepareStringToSet(env))

		fmt.Print("Enter server address: ")
		ip, _ := reader.ReadString('\n')

		fmt.Print("Enter server login: ")
		login, _ := reader.ReadString('\n')

		fmt.Print("Enter server auth type (key or password only): ")
		authType, _ := reader.ReadString('\n')

		if prepareStringToSet(authType) == "key" {
			fmt.Print("Enter key location: ")
			key, _ = reader.ReadString('\n')
		} else {
			fmt.Print("Enter server password: ")
			password, _ = reader.ReadString('\n')
		}

		fmt.Print("Enter project home dir on server: ")
		homeDir, _ := reader.ReadString('\n')

		fmt.Print("Enter git repository https address: ")
		gitRepo, _ := reader.ReadString('\n')

		fmt.Print("Enter git user name: ")
		gitUser, _ := reader.ReadString('\n')

		fmt.Print("Enter git user password: ")
		gitPassword, _ := reader.ReadString('\n')

		fmt.Print("Enter git branch name: ")
		gitBranch, _ := reader.ReadString('\n')

		envIns.EnvType = prepareStringToSet(env)
		envIns.Server = prepareStringToSet(ip)
		envIns.Login = prepareStringToSet(login)
		envIns.Password = prepareStringToSet(password)
		envIns.HomeDir = prepareStringToSet(homeDir)
		envIns.AuthType = prepareStringToSet(authType)
		envIns.KeyFile = prepareStringToSet(key)

		gitConfigIns := GitConfig{
			Repository: prepareStringToSet(gitRepo),
			User: prepareStringToSet(gitUser),
			Password: prepareStringToSet(gitPassword),
			Branch: prepareStringToSet(gitBranch),
		}
		envIns.GitConfig = gitConfigIns

	InputBeforeDeploy:
		fmt.Print("Enter commands before deploy or blank for exit: ")
		beforeCommand, _ := reader.ReadString('\n')
		if prepareStringToSet(beforeCommand) != "" {
			commandItem := Command{Item:prepareStringToSet(beforeCommand)}
			envIns.BeforeDeploy = append(envIns.BeforeDeploy, commandItem)
			goto InputBeforeDeploy
		}

	InputAfterDeploy:
		fmt.Print("Enter commands after deploy or blank for exit: ")
		afterCommand, _ := reader.ReadString('\n')

		if prepareStringToSet(afterCommand) != "" {
			commandItem := Command{Item:prepareStringToSet(afterCommand)}
			envIns.AfterDeploy = append(envIns.AfterDeploy, commandItem)
			goto InputAfterDeploy
		}

		conf.Envs = append(conf.Envs, envIns)

		if len(Difference(allowed, entered)) > 0 {
			goto CreateStep
		}
	}

	i.Conf = conf
	//fmt.Print("Enter project version: ")
	//version, _ := reader.ReadString('\n')
	//conf.Version = prepareStringToSet(version)
}

func prepareStringToSet (str string) string {
	sep := "\n"
	if runtime.GOOS == "windows" {
		sep = "\r\n"
	}
	return strings.Replace(str, sep, "", 1)
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
