package utils

import (
	"bytes"
	"os/exec"
	"strings"
)

func (u *Utils) ExecuteCommandV2(cmd string, args []string) ([]byte, error) {
	var out bytes.Buffer
	// Execute the command
	command := exec.Command(cmd, args...)
	u.Log.Infof("cmd: %s, args: %s", cmd, strings.Join(args, " "))
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		u.Log.Errorf("Error running command: %s", err.Error())
		return out.Bytes(), err
	}
	// TODO need to figure-out a way to log this output
	u.Log.Infoln(out.String())
	return out.Bytes(), nil
}

func (u *Utils) ExecuteCommand(command []string) ([]byte, error) {
	u.Log.Infof("Executing command %s", strings.Join(command, " "))
	var output, outputErrors bytes.Buffer
	cmd0 := command[0]
	cmds := command[1:]
	cmd := exec.Command(cmd0, cmds...)
	cmd.Stdout = &output
	cmd.Stderr = &outputErrors
	err := cmd.Run()
	if err != nil {
		u.Log.Errorf("Error executing command %s %s. %s", cmd0, strings.Join(command, " "), err.Error())
		return output.Bytes(), err
	}
	u.Log.Infoln("Command output is:", output.String())
	u.Log.Infoln("Command output error is:", outputErrors.String())
	return output.Bytes(), err
}

func (u *Utils) SudoExecuteCommand(command []string) ([]byte, error) {
	u.Log.Infof("Executing command %s", strings.Join(command, " "))
	var output, outputErrors bytes.Buffer
	path, pathErr := exec.LookPath("sudo")
	if pathErr != nil {
		u.Log.Infof("Error looking-up path for \"sudo\" %s", pathErr.Error())
		return output.Bytes(), pathErr
	}

	u.Log.Infof("Found path for \"sudo\" %s", path)
	cmd0 := path
	u.Log.Infof("Command being executed will be %s %s", cmd0, strings.Join(command, " "))
	cmd := exec.Command(cmd0, command...)
	cmd.Stdout = &output
	cmd.Stderr = &outputErrors
	err := cmd.Run()
	if err != nil {
		u.Log.Errorf("Error executing command %s %s. %s", cmd0, strings.Join(command, " "), err.Error())
		return output.Bytes(), err
	}
	u.Log.Infoln("Command output is:", output.String())
	u.Log.Infoln("Command output error is:", outputErrors.String())
	return output.Bytes(), err
}
