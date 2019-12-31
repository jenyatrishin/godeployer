package deployer

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"../shellRunner"
	"../tools"
)

const (
	AUTH_PASSWORD string  = "password"
	AUTH_KEY string = "key"
)

type AbstractDeployer interface {
	DeployTo(dir string, repo string, gitUser string, gitPass string, gitBranch string, afterCommands string, beforeCommands string)
	PrepareConfig(host string, port string, user string, authMethodValue string, authMethod string)
}

type SshDeployer struct {
	config *ssh.ClientConfig
	host string
	port string
}

func (d *SshDeployer) PrepareConfig (host string, port string, user string, authMethodValue string, authMethod string) {
	var auth []ssh.AuthMethod
	if authMethod == AUTH_PASSWORD {
		auth = []ssh.AuthMethod{
			ssh.Password(authMethodValue),
		}
	} else if authMethod == AUTH_KEY {
		auth = []ssh.AuthMethod{
			publicKeyFile(authMethodValue),
		}
	}
	d.host = host
	d.port = port
	d.config = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func (d *SshDeployer) DeployTo (dir string, repo string, gitUser string, gitPass string, gitBranch string, afterCommands string, beforeCommands string) {

	d.executeBeforeCommand(beforeCommands)

	addr := d.host+":"+d.port
	//logs
	tools.WriteLog("Started deploy to address: " + addr)
	client, err := ssh.Dial("tcp", addr, d.config)

	//should be changed to tools/error object
	if err != nil{
		//logs
		tools.WriteLog("Connection server error: " + string(err.Error()))
		fmt.Println("Can not connect to server. Please check server connection config")
		return
	}

	session, err := client.NewSession()

	if err != nil {
		//logs
		tools.WriteLog("Server session error: " + string(err.Error()))
		fmt.Println("Can not open session on server.")
		return
	}
	defer session.Close()

	// Enable system stdout
	// Comment these if you uncomment to store in variable
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// StdinPipe for commands
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Start remote shell
	err = session.Shell()

	if err != nil {
		log.Fatal(err)
	}

	if gitUser == "" || gitPass == "" || repo == "" || gitBranch == "" {
		fmt.Println("Git credentials not defined")
		return
	}

	gitDotDir := dir + "/.git"
	gitAddress := "https://"+gitUser+":"+gitPass+"@"+repo
	gitPullCommand := "git pull "+gitAddress+" "+gitBranch
	gitCloneCommand := "git clone " + gitAddress + " ."

	commands := []string{
		"cd "+dir,
		"if test -d "+gitDotDir+"; then "+gitPullCommand+"; else "+gitCloneCommand+"; fi",
		afterCommands,
		"exit",
	}

	commandsString := strings.Join(commands," && ")
	//logs
	tools.WriteLog("Deploy command on server: " + commandsString)

	_, err = fmt.Fprintf(stdin, "%s\n", commandsString)
	if err != nil {
		log.Fatal(err)
		//logs
		tools.WriteLog("Deploy error: " + string(err.Error()))
	}

	// Wait for sess to finish
	err = session.Wait()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deploy is finished")
	//logs
	tools.WriteLog("deploy is finished")
}

func (d *SshDeployer) executeBeforeCommand (beforeCommands string) {
	commands := strings.Split(beforeCommands, "&&")
	if len(commands) > 0 {
		shell := new(shellRunner.Shell)
		for _, one := range commands {
			oneCommand := strings.Split(one, " ")
			command := oneCommand[0]
			args := strings.Replace(one, command+" ", "", 1)
			shell.ExecuteCommand(command, args)
			shell.PrintResult()
		}
	}
}

func publicKeyFile(file string) ssh.AuthMethod {

	buffer, err := ioutil.ReadFile(file)

	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)

	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}
