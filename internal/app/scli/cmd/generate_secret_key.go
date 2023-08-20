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
	"github.com/shammishailaj/scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// generateSecretKeyCmd represents the cleanCmd command
var generateSecretKeyCmd = &cobra.Command{
	Use:   "appsecret [options] /path/to/new/wallpaper",
	Short: "Used to generate APP_SECRET for a framework",
	Long:  `Used to generate APP_SECRET for a framework`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Running appsecret...")

		appName, appNameErr := cmd.Flags().GetString("framework")
		if appNameErr != nil {
			u.Log.Errorf("Error getting value for framework parameter. %s", appNameErr.Error())
			u.Log.Infof("Using default value: symfony")
			appName = "symfony"
		}
		appName = strings.ToLower(appName)

		secretLength := 32

		secretStr, secretStrErr := u.RandomBytes(secretLength)
		if secretStrErr != nil {
			u.Log.Errorf("Error reading random bytes: %s", secretStrErr.Error())
			os.Exit(utils.GENERATE_APPSECRET_ERROR_GETTING_RANDOM_BYTES)
		}

		appSecret, appSecretErr := u.Bin2hex(secretStr)
		if appSecretErr != nil {
			u.Log.Errorf("Error converting secret string %s to app secret: %s", secretStr, appSecretErr.Error())
			os.Exit(utils.GENERATE_APPSECRET_ERROR_CONVERTING_RANDOM_BYTES_TO_APPSECRET)
		}

		u.Log.Infof("APP_SECRET=\"%s\"", appSecret)
	},
}

func init() {
	generateCmd.AddCommand(generateSecretKeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modifyWallpaperCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	generateSecretKeyCmd.Flags().StringP("framework", "a", "symfony", "Generate APP_SECRET for a framework. Valid values are: symfony")
}
