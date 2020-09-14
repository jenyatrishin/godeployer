package config

import (
	"fmt"
	"runtime"
	"strings"
)

type Input struct {
	Conf Config
}

func (i *Input) callUserInput () {
	allowed := []string{"developer","staging","production"}
	var entered []string
	var conf Config
	encoder := Encoder{}

CreateStep:
	var envIns Env

	fmt.Print("Enter environment name ("+ strings.Join(Difference(allowed, entered),",")+" or blank for continue): ")
	env := ""
	fmt.Scanf("%s", &env)

	if prepareStringToSet(env) != "" {

		password := ""
		key := ""
		ip := ""
		login := ""
		authType := ""
		homeDir := ""
		gitRepo := ""
		gitUser := ""
		gitPassword := ""
		gitBranch := ""

		entered = append(entered, prepareStringToSet(env))

		fmt.Print("Enter server address: ")
		fmt.Scanf("%s", &ip)

		fmt.Print("Enter server login: ")
		fmt.Scanf("%s", &login)

		fmt.Print("Enter server auth type (key or password only): ")
		fmt.Scanf("%s", &authType)

		if prepareStringToSet(authType) == "key" {
			fmt.Print("Enter key location: ")
			fmt.Scanf("%s", &key)
		} else {
			fmt.Print("Enter server password: ")
			fmt.Scanf("%s", &password)
		}

		fmt.Print("Enter project home dir on server: ")
		fmt.Scanf("%s", &homeDir)

		fmt.Print("Enter git repository https address: ")
		fmt.Scanf("%s", &gitRepo)

		fmt.Print("Enter git user name: ")
		fmt.Scanf("%s", &gitUser)

		fmt.Print("Enter git user password: ")
		fmt.Scanf("%s", &gitPassword)

		fmt.Print("Enter git branch name: ")
		fmt.Scanf("%s", &gitBranch)

		envIns.EnvType = prepareStringToSet(env)
		envIns.Server = prepareStringToSet(ip)
		envIns.Login = prepareStringToSet(login)
		if len(password) > 0 {
			envIns.Password = encoder.Encode(prepareStringToSet(password))
		} else {
			envIns.Password = []byte(prepareStringToSet(password))
		}

		envIns.HomeDir = prepareStringToSet(homeDir)
		envIns.AuthType = prepareStringToSet(authType)
		envIns.KeyFile = prepareStringToSet(key)

		gitConfigIns := GitConfig{
			Repository: prepareStringToSet(gitRepo),
			User: prepareStringToSet(gitUser),
			Password: encoder.Encode(prepareStringToSet(gitPassword)),
			Branch: prepareStringToSet(gitBranch),
		}
		envIns.GitConfig = gitConfigIns

	InputBeforeDeploy:
		fmt.Print("Enter commands before deploy or blank for continue: ")
		beforeCommand := ""
		fmt.Scanf("%s", &beforeCommand)
		if prepareStringToSet(beforeCommand) != "" {
			commandItem := Command{Item:prepareStringToSet(beforeCommand)}
			envIns.BeforeDeploy = append(envIns.BeforeDeploy, commandItem)
			goto InputBeforeDeploy
		}

	InputAfterDeploy:
		fmt.Print("Enter commands after deploy or blank for continue: ")
		afterCommand := ""
		fmt.Scanf("%s", &afterCommand)

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

	fmt.Print("Enter project name: ")
	projectName := ""
	fmt.Scanf("%s", &projectName)
	conf.ProjectName = prepareStringToSet(projectName)

	i.Conf = conf
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
