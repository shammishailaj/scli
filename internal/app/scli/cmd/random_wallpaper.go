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
	"github.com/spf13/cobra"
	"os"
)

// cleanCmd represents the cleanCmd command
var randomWallpaperCmd = &cobra.Command{
	Use:   "wallpaper",
	Short: "Set random wallpaper",
	Long:  `Set random wallpaper`,
	Run: func(cmd *cobra.Command, args []string) {
		u.Log.Infof("Setting random wallpaper....")

		authorization, authorizationErr := cmd.Flags().GetString("authorization")
		if authorizationErr != nil {
			u.Log.Fatalf("Can not proceed without a valid api key. %s\n", authorizationErr.Error())
		}

		if len(authorization) < 32 {
			u.Log.Fatalf("Can not proceed without a valid api key....\n")
		}

		query, queryErr := cmd.Flags().GetString("query")
		if queryErr != nil {
			u.Log.Fatalf("Can not proceed without a valid search query. %s", queryErr.Error())
		}

		if len(query) < 1 {
			u.Log.Fatalf("Can not proceed without a valid search query. Empty string passed")
		}

		orientation, orientationErr := cmd.Flags().GetString("orientation")
		if orientationErr != nil {
			u.Log.Infof("Error reading value for orientation. %s. Will use: landscape", orientationErr.Error())
		}

		size, sizeErr := cmd.Flags().GetString("size")
		if sizeErr != nil {
			u.Log.Infof("Error reading value for size. %s. Will use: large", sizeErr.Error())
		}

		color, colorErr := cmd.Flags().GetString("color")
		if colorErr != nil {
			u.Log.Errorf("Error parsing value for color. %s. Will not send a color in request", colorErr.Error())
			color = ""
		}

		locale, localeErr := cmd.Flags().GetString("locale")
		if localeErr != nil {
			u.Log.Infof("Error parsing value for locale. %s. Will not send a locale in request", localeErr.Error())
			locale = ""
		}

		dbusSessionBusAddress, dbusSessionBusAddressErr := cmd.Flags().GetString("dbus-session-bus-address")
		if dbusSessionBusAddressErr != nil {
			u.Log.Fatalf("Can not continue without valid value for parameter --dbus-session-bus-address. %s", dbusSessionBusAddressErr.Error())
		}

		if dbusSessionBusAddress == "" {
			u.Log.Fatalf("Can not continue with empty value for parameter --dbus-session-bus-address")
		}

		dbusSessionBusAddressSetErr := os.Setenv("DBUS_SESSION_BUS_ADDRESS", dbusSessionBusAddress)
		if dbusSessionBusAddressSetErr != nil {
			u.Log.Fatalf("Error setting DBUS_SESSION_BUS_ADDRESS environment variable. Can not continue without it")
		}

		cachePath, cachePathErr := cmd.Flags().GetString("cache-path")
		if cachePathErr != nil {
			u.Log.Infof("Cache directory not provided. Will use default %s", utils.RANDOM_WALLPAPER_DEFAULT_CACHE_DIR)
			cachePath = utils.RANDOM_WALLPAPER_DEFAULT_CACHE_DIR
		}

		dbFilePath, dbFilePathErr := cmd.Flags().GetString("db-file")
		if dbFilePathErr != nil {
			u.Log.Fatalf("DB directory not provided. Can not continue.\n")
		}

		//randomWallpaperErr := u.RandomPexelsWallpaper(authorization, query, orientation, size, color, locale)
		//randomWallpaperErr := u.RandomPexelsWallpaperWithCache(authorization, query, orientation, size, color, locale, cachePath)
		randomWallpaperErr := u.RandomPexelsWallpaperWithCacheAndDB(authorization, query, orientation, size, color, locale, cachePath, dbFilePath)
		if randomWallpaperErr != nil {
			u.Log.Errorf("Error setting random wallpaper by downloading from pexels website %s\n", randomWallpaperErr.Error())
			u.Log.Infof("Setting a random wallpaper from cache at %s...", cachePath)
			randomWallpaperFromCacheErr := u.RandomPexelsWallpaperFromCache(cachePath)
			if randomWallpaperFromCacheErr != nil {
				u.Log.Fatalf("scli.random.wallpaper: error setting randome wallpaper from cache: %s", randomWallpaperFromCacheErr)
			}
		}

		u.Log.Infof("Wallpaper changed successfully...\n")
	},
}

func init() {
	randomCmd.AddCommand(randomWallpaperCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// randomWallpaperCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// randomWallpaperCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	randomWallpaperCmd.Flags().StringP("authorization", "a", "", "Pexels API Key used in Authorization request header")
	randomWallpaperCmd.Flags().StringP("color", "c", "", "(Optional) Desired photo color. Supported colors: red, orange, yellow, green, turquoise, blue, violet, pink, brown, black, gray, white or any hexadecimal color code (eg. #ffffff)")
	randomWallpaperCmd.Flags().StringP("dbus-session-bus-address", "d", "", "Value of the gnome variable \"DBUS_SESSION_BUS_ADDRESS\" from your shell. Do an \"echo $DBUS_SESSION_BUS_ADDRESS\" to see its value. On my system (Ubuntu 22.04.1 LTS with Gnome) it was \"unix:path=/run/user/1000/bus\"")
	randomWallpaperCmd.Flags().StringP("locale", "l", "", "The locale of the search you are performing. The current supported locales are: 'en-US' 'pt-BR' 'es-ES' 'ca-ES' 'de-DE' 'it-IT' 'fr-FR' 'sv-SE' 'id-ID' 'pl-PL' 'ja-JP' 'zh-TW' 'zh-CN' 'ko-KR' 'th-TH' 'nl-NL' 'hu-HU' 'vi-VN' 'cs-CZ' 'da-DK' 'fi-FI' 'uk-UA' 'el-GR' 'ro-RO' 'nb-NO' 'sk-SK' 'tr-TR' 'ru-RU'")
	randomWallpaperCmd.Flags().Uint64P("max-results", "m", 15, "The number of results you are requesting per page. Default: 15 Max: 80")
	randomWallpaperCmd.Flags().StringP("orientation", "o", "landscape", "(Optional) Desired photo orientation. The current supported orientations are: landscape, portrait or square")
	randomWallpaperCmd.Flags().StringP("query", "q", "beautiful", "Query to search for")
	randomWallpaperCmd.Flags().StringP("size", "s", "large", "(Optional) Minimum photo size. The current supported sizes are: large(24MP), medium(12MP) or small(4MP)")
	randomWallpaperCmd.Flags().StringP("save-file-path", "t", "./", "Path to the file in which to save the Image. Path should exist, file should not")
	randomWallpaperCmd.Flags().StringP("cache-path", "u", utils.RANDOM_WALLPAPER_DEFAULT_CACHE_DIR, "Path to the directory where the downloaded image files shall be cached without trailing slashes. Default \"~/.cache\". If path does not exist, it shall be created automatically. Files shall be stored inside this path inside \"random/wallpaper\" directory")
	randomWallpaperCmd.Flags().StringP("db-file", "w", "", "Path to the GenjiDB file where the metabadate for the downloaded image files along with wallpaper change logs shall be stored")
}
