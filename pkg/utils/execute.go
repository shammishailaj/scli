package utils

import (
	"os/exec"
	"strings"
)

func (u *Utils) ExecuteCommand(command []string) ([]byte, error) {
	var output []byte
	cmd0 := command[0]
	cmds := command[1:]
	cmd := exec.Command(cmd0, cmds...)
	err := cmd.Run()
	if err != nil {
		u.Log.Errorf("Error executing command %s %s. %s", cmd0, strings.Join(command, " "), err.Error())
		return output, err
	}
	output, err = cmd.Output()
	return output, err
}

func (u *Utils) SudoExecuteCommand(command []string) ([]byte, error) {
	var output []byte
	path, pathErr := exec.LookPath("sudo")
	if pathErr != nil {
		return output, pathErr
	}

	cmd0 := path
	cmd := exec.Command(cmd0, command...)
	err := cmd.Run()
	if err != nil {
		u.Log.Errorf("Error executing command %s %s. %s", cmd0, strings.Join(command, " "), err.Error())
		return output, err
	}
	output, err = cmd.Output()
	return output, err
}
