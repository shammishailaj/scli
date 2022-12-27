/*
Copyright © 2022  <>

Licensed under the HLT License, Version 0.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/shammishailaj/scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// cleanCmd represents the cleanCmd command
var modifyWallpaperCmd = &cobra.Command{
	Use:   "wallpaper [options] /path/to/new/wallpaper",
	Short: "Used to modify the current system wallpaper",
	Long:  `Used to modify the current system wallpaper`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Running wallpaper...")
		//gsettings set org.gnome.desktop.background picture-uri 'file://PathToImage'
		//gsettings get org.gnome.desktop.background picture-uri 'file:///usr/share/backgrounds/warty-final-ubuntu.png' as per https://linuxconfig.org/set-wallpaper-on-ubuntu-20-04-using-command-line

		envName, envNameErr := cmd.Flags().GetString("environment")
		if envNameErr != nil {
			u.Log.Errorf("Error getting value for environment parameter. %s", envNameErr.Error())
			u.Log.Infof("Using default value: gnome")
			envName = "gnome"
		}
		envName = strings.ToLower(envName)

		wallpaperPath := args[0]

		if !u.FileExists(args[0]) {
			u.Log.Errorf("Error finding absolute file path for %s\n", args[0])
			os.Exit(utils.MODIFY_WALLPAPER_FILE_DOES_NOT_EXIST)
		}

		absWallpaperPath, absWallpaperPathErr := filepath.Abs(args[0])
		if absWallpaperPathErr != nil {
			u.Log.Errorf("Error finding absolute file path for %s. %s", args[0], absWallpaperPathErr.Error())
			os.Exit(utils.MODIFY_WALLPAPER_ABS_FILEPATH_NOT_FOUND)
		}

		wallpaperPath = absWallpaperPath

		commands := make(map[string]map[string][]string)

		gnomeCommands := make(map[string][]string)
		gnomeCommands["modify-wallpaper-gnome"] = []string{"gsettings", "set", "org.gnome.desktop.background", "picture-uri", fmt.Sprintf("file://%s", wallpaperPath)}
		gnomeCommands["modify-wallpaper-gnome-dark"] = []string{"gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", fmt.Sprintf("file://%s", wallpaperPath)}

		commands["gnome"] = gnomeCommands

		for commandName, command := range commands[envName] {
			log.Infof("Running Command for %s", commandName)
			output, err := u.ExecuteCommand(command)

			if err != nil {
				log.Errorf("Error running shell command %s. %s", commandName, err.Error())
			}

			log.Infof("Output:\n%s", output)
		}
	},
}

func init() {
	modifyCmd.AddCommand(modifyWallpaperCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modifyWallpaperCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	modifyWallpaperCmd.Flags().StringP("environment", "a", "gnome", "If set, will execute commands for specified desktop environment. Valid values are: gnome")
}
