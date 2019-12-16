package shellRunner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Shell struct {
	Result string
	Status string
}

func (shell *Shell) ExecuteCommand (command string, arguments string) bool {
	cmd := exec.Command(command, arguments)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		shell.Status = "error"
		shell.Result = err.Error()
		os.Stderr.WriteString(err.Error())
		return false
	}
	shell.Status = "success"
	shell.Result = string(cmdOutput.Bytes())

	return true
}

func (shell *Shell) PrintResult () bool {
	fmt.Println(shell.Result)
	return true
}
