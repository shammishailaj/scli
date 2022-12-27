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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// cleanCmd represents the cleanCmd command
var modifyWallpaperCmd = &cobra.Command{
	Use:   "wallpaper",
	Short: "Used to modify the current system wallpaper",
	Long:  `Used to modify the current system wallpaper`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Running wallpaper...")
		//gsettings set org.gnome.desktop.background picture-uri 'file://PathToImage'
		//gsettings get org.gnome.desktop.background picture-uri 'file:///usr/share/backgrounds/warty-final-ubuntu.png' as per https://linuxconfig.org/set-wallpaper-on-ubuntu-20-04-using-command-line

		commands := make(map[string][]string, 1)

		commands["modify-wallpaper"] = []string{"gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file:///home/sam/Pictures/4k_1.jpeg"}

		for commandName, command := range commands {
			log.Infof("Running Command for %s", commandName)
			output, err := u.SudoExecuteCommand(command)

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
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
