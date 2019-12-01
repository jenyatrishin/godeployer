package xml

import (
	"../../config"
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type ConfigAdapterXml struct {

}

func (adapter ConfigAdapterXml) ReadConfigFromFile (c interface{}, name string) []byte {
	file, err := os.Open(name)

	if err != nil {
		panic("There ain't file")
	}

	fmt.Println("Config file is opened")

	byteValue, _ := ioutil.ReadAll(file)

	//xml.Unmarshal(byteValue, &c)

	defer file.Close()
	return byteValue
	//return true
}

func (adapter ConfigAdapterXml) WriteConfigToFile (name string) interface{} {

	allowed := []string{"developer","staging","production"}
	var entered []string

	reader := bufio.NewReader(os.Stdin)

	//conf := tools.GetConfigIns()
	var conf config.Config

CreateStep:
	var envIns config.Env

	fmt.Print("Enter environment name ("+ strings.Join(Difference(allowed, entered),",")+" or blank for exit): ")
	env, _ := reader.ReadString('\n')


	if strings.Replace(env, "\n", "", 1) != "" {

		entered = append(entered, strings.Replace(env, "\n", "", 1))

		fmt.Print("Enter server address: ")
		ip, _ := reader.ReadString('\n')

		fmt.Print("Enter server login: ")
		login, _ := reader.ReadString('\n')

		fmt.Print("Enter server password: ")
		password, _ := reader.ReadString('\n')

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

		envIns.EnvType = strings.Replace(env, "\n", "", 1)
		envIns.Server = strings.Replace(ip, "\n", "", 1)
		envIns.Login = strings.Replace(login, "\n", "", 1)
		envIns.Password = strings.Replace(password, "\n", "", 1)
		envIns.HomeDir = strings.Replace(homeDir, "\n", "", 1)

		gitConfigIns := config.GitConfig{
			Repository: prepareStringToSet(gitRepo),
			User: prepareStringToSet(gitUser),
			Password: prepareStringToSet(gitPassword),
			Branch: prepareStringToSet(gitBranch),
		}
		envIns.GitConfig = gitConfigIns

	InputBeforeDeploy:
		fmt.Print("Enter commands before deploy or blank for exit: ")
		beforeCommand, _ := reader.ReadString('\n')

		if strings.Replace(beforeCommand, "\n", "", 1) != "" {
			commandItem := config.Command{Item:prepareStringToSet(beforeCommand)}
			envIns.BeforeDeploy = append(envIns.BeforeDeploy, commandItem)
			goto InputBeforeDeploy
		}

	InputAfterDeploy:
		fmt.Print("Enter commands after deploy or blank for exit: ")
		afterCommand, _ := reader.ReadString('\n')

		if strings.Replace(afterCommand, "\n", "", 1) != "" {
			commandItem := config.Command{Item:strings.Replace(afterCommand, "\n", "", 1)}
			envIns.AfterDeploy = append(envIns.AfterDeploy, commandItem)
			goto InputAfterDeploy
		}

		conf.Envs = append(conf.Envs, envIns)

		//fmt.Println(conf)

		if len(Difference(allowed, entered)) > 0 {
			goto CreateStep
		}
	}


	fmt.Print("Enter project version: ")
	version, _ := reader.ReadString('\n')
	conf.Version = strings.Replace(version,"\n","",1)

	file, _ := xml.MarshalIndent(conf, "", " ")

	_ = ioutil.WriteFile(name, file, 0777)

	return "create config"
}

func prepareStringToSet (str string) string {
	return strings.Replace(str, "\n", "", 1)
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