package utils

import (
	"errors"
	"fmt"
)

func (u *Utils) DPKGReconfigureUnattendedUpgrades(disable bool) ([]byte, error) {
	var cmd []string
	var cmdErr error

	if disable {
		cmd, cmdErr = u.DPKGCommands("disable-unattended-upgrades")
		if cmdErr != nil {
			return []byte{}, cmdErr
		}
	} else {
		cmd, cmdErr = u.DPKGCommands("enable-unattended-upgrades")
		if cmdErr != nil {
			return []byte{}, cmdErr
		}
	}

	return u.ExecuteCommand(cmd)
}

func (u *Utils) DPKGCommands(commandName string) ([]string, error) {
	commands := make(map[string]map[string][]string)

	dpkgCommands := make(map[string][]string)
	dpkgCommands["enable-unattended-upgrades"] = []string{"dpkg-reconfigure", "-pmedium", "unattended-upgrades"}  // As per https://unix.stackexchange.com/a/694582/91242
	dpkgCommands["disable-unattended-upgrades"] = []string{"dpkg-reconfigure", "-pmedium", "unattended-upgrades"} // As per https://unix.stackexchange.com/a/694582/91242

	commands["dpkg"] = dpkgCommands
	retVal, ok := commands["dpkg"][commandName]
	if ok {
		return retVal, nil
	}

	return []string{}, errors.New(fmt.Sprintf("Command for %s not found", commandName))
}
