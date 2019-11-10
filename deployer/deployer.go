package deployer

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
)

type AbstractDeployer interface {
	DeployTo(dir string)
	PrepareConfig(host string, port string, user string, password string)
}

type SshDeployer struct {
	config *ssh.ClientConfig
	host string
	port string
}

func (d *SshDeployer) PrepareConfig (host string, port string, user string, password string) {
	d.host = host
	d.port = port
	d.config = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

}

func (d *SshDeployer) DeployTo (dir string) {
	//addr := fmt.Sprintf("%s:%d", d.host, d.port)
	addr := d.host+":"+d.port

	client, err := ssh.Dial("tcp", addr, d.config)

	//should be changed to tools/error object
	if err != nil{
		panic(err)
	}

	session, err := client.NewSession()

	if err != nil {
		panic(err)
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

	//session.Run("cd "+ dir + "&& pwd")

	_, err = fmt.Fprintf(stdin, "%s\n", "cd "+ dir + "&& pwd")
	if err != nil {
		log.Fatal(err)
	}

	// Wait for sess to finish
	err = session.Wait()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("finish")
}