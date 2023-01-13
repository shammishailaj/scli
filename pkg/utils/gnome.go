package utils

import (
	"errors"
	"fmt"
	"os"
)

func (u *Utils) GnomeShowSeconds() ([]byte, error) {
	cmd, cmdErr := u.GnomeCommands("enable-seconds-gnome")
	if cmdErr != nil {
		return []byte{}, cmdErr
	}

	return u.ExecuteCommand(cmd)
}

func (u *Utils) GnomeCommands(commandName string) ([]string, error) {
	commands := make(map[string]map[string][]string)

	gnomeCommands := make(map[string][]string)
	gnomeCommands["modify-wallpaper-gnome"] = []string{"gsettings", "set", "org.gnome.desktop.background", "picture-uri", fmt.Sprintf("file://__WALLPAPER_PATH__")}
	gnomeCommands["modify-wallpaper-gnome-dark"] = []string{"gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", fmt.Sprintf("file://__WALLPAPER_PATH__")}
	gnomeCommands["enable-seconds-gnome"] = []string{"gsettings", "set", "org.gnome.desktop.interface", "clock-show-seconds", "true"}
	gnomeCommands["disable-seconds-gnome"] = []string{"gsettings", "set", "org.gnome.desktop.interface", "clock-show-seconds", "false"}
	gnomeCommands["enable-weekday-gnome"] = []string{"gsettings", "set", "org.gnome.desktop.interface", "clock-show-weekday", "true"}
	gnomeCommands["disable-weekday-gnome"] = []string{"gsettings", "set", "org.gnome.desktop.interface", "clock-show-weekday", "false"}

	commands["gnome"] = gnomeCommands
	retVal, ok := commands["gnome"][commandName]
	if ok {
		return retVal, nil
	}

	return []string{}, errors.New(fmt.Sprintf("Command for %s not found", commandName))
}

func (u *Utils) GnomeGetDBUSSessionBusAddress() string {
	return os.Getenv("DBUS_SESSION_BUS_ADDRESS")
}
