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
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/shammishailaj/scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math/big"
	"net/url"
	"os"
	"path"
	"strings"
)

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

// cleanCmd represents the cleanCmd command
var pexelsGetimagesCmd = &cobra.Command{
	Use:   "getimages",
	Short: "Used to get a new image from pexels",
	Long:  `Used to get a new image from pexels`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("getting pexels image")
		input := utils.PexelsSearchInput{
			Query:       "people",
			Orientation: "landscape",
			Size:        "large",
			Color:       "",
			Locale:      "",
			Page:        1,
			PerPage:     15,
		}

		authorization, authorizationErr := cmd.Flags().GetString("authorization")
		if authorizationErr != nil {
			log.Fatalf("Can not proceed without a valid api key. %s\n", authorizationErr.Error())
		}

		if len(authorization) < 32 {
			log.Fatalf("Can not proceed without a valid api key....\n")
		}

		query, queryErr := cmd.Flags().GetString("query")
		if queryErr != nil {
			log.Fatalf("Can not proceed without a valid search query. %s", queryErr.Error())
		}

		if len(query) < 1 {
			log.Fatalf("Can not proceed without a valid search query. Empty string passed")
		}

		input.Query = query

		orientation, orientationErr := cmd.Flags().GetString("orientation")
		if orientationErr != nil {
			log.Infof("Error reading value for orientation. %s. Will use: landscape", orientationErr.Error())
		}

		switch orientation {
		case utils.PEXELS_ORIENTATION_LANDSCAPE:
			input.Orientation = utils.PEXELS_ORIENTATION_LANDSCAPE
		case utils.PEXELS_ORIENTATION_PORTRAIT:
			input.Orientation = utils.PEXELS_ORIENTATION_PORTRAIT
		case utils.PEXELS_ORIENTATION_SQUARE:
			input.Orientation = utils.PEXELS_ORIENTATION_SQUARE
		default:
			input.Orientation = utils.PEXELS_ORIENTATION_LANDSCAPE
		}

		size, sizeErr := cmd.Flags().GetString("size")
		if sizeErr != nil {
			log.Infof("Error reading value for size. %s. Will use: large", sizeErr.Error())
		}

		switch size {
		case utils.PEXELS_IMAGE_LARGE:
			input.Size = utils.PEXELS_IMAGE_LARGE
		case utils.PEXELS_IMAGE_MEDIUM:
			input.Size = utils.PEXELS_IMAGE_MEDIUM
		case utils.PEXELS_IMAGE_SMALL:
			input.Size = utils.PEXELS_IMAGE_SMALL
		default:
			input.Size = utils.PEXELS_IMAGE_LARGE
		}

		color, colorErr := cmd.Flags().GetString("color")
		if colorErr != nil {
			log.Infof("Error parsing value for color. %s. Will not send a color in request", colorErr.Error())
		}

		if len(color) < 0 {
			log.Infof("Invalid color. Will not send a color in request")
		}

		// Check this function to see if color is valid https://go.dev/play/p/rM0e-w7Xfdg
		// Original function at: https://www.geeksforgeeks.org/check-if-a-given-string-is-a-valid-hexadecimal-color-code-or-not/

		locale, localeErr := cmd.Flags().GetString("locale")

		if localeErr != nil {
			log.Infof("Error parsing value for locale. %s. Will not send a locale in request", localeErr.Error())
		}

		if !u.InArray(locale, utils.Locales) {
			log.Infof("Illegal value for locale. Will not send a locale in request")
		}

		page, pageErr := cmd.Flags().GetUint64("page")
		if pageErr != nil {
			log.Infof("Error parsing value for page. %s. Will use the defaul value of 1", pageErr.Error())
		}

		input.Page = page

		maxResults, maxResultsErr := cmd.Flags().GetUint64("max-results")

		if maxResultsErr != nil {
			log.Errorf("Error parsing value for maxResults. %s. Will use the default value of 15", maxResultsErr.Error())
		}

		input.PerPage = maxResults

		saveFilePath, saveFilePathErr := cmd.Flags().GetString("save-file-path")
		if saveFilePathErr != nil {
			u.Log.Errorf("Error getting save file path. %s", saveFilePathErr.Error())
			saveFilePath = ""
		}

		randomize, randomizeErr := cmd.Flags().GetBool("randomize")
		if randomizeErr != nil {
			u.Log.Errorf("Error getting value for parameter --randomize. %s", randomizeErr.Error())
			u.Log.Infof("Using default value for --randomize: false")
			randomize = false
		}

		pexels := u.NewPexels(authorization)
		log.Infof("Sending request to Pexels search API...")
		output := pexels.Search(input)
		photoNumber := int64(1)

		if randomize {
			u.Log.Infof("Getting a random photo...")
			u.Log.Infof("Generating a random photo number...")
			randomPhotoNumber, randomPhotoNumberErr := rand.Int(rand.Reader, big.NewInt(int64(output.TotalResults)))
			if randomPhotoNumberErr != nil {
				u.Log.Errorf("Error getting random number between 0 and %d", output.TotalResults)
				u.Log.Infof("Choosing the 1st image at position 0")
				randomPhotoNumber = big.NewInt(0)
			}

			photoNumber = randomPhotoNumber.Int64()
			u.Log.Infof("Getting photo number #%d...", photoNumber)

			photo, photoErr := output.GetPhotoByNumber(photoNumber)
			if photoErr != nil {
				u.Log.Fatalf("Error getting photo number %d. %s", photoNumber, photoErr.Error())
			}
			if photo != nil {
				u.Log.Infof("Photo Found. Will download Large2X size at URL: %s", photo.Src.Large2X)
				urlParsed, urlErr := url.Parse(photo.Src.Large2X)
				if urlErr != nil {
					log.Printf("Error parsing URL %s. %s", u, urlErr.Error())
					urlParsed = &url.URL{Path: photo.Src.Large2X}
				}

				if saveFilePath == "" {
					cwd, cwdErr := os.Getwd()
					if cwdErr != nil {
						u.Log.Errorf("Error getting current working directory. %s", cwdErr.Error())
						cwd = os.TempDir()
						u.Log.Infof("Setting directory to: %s", cwd)
					}
					u.Log.Infof("Will use current working directory for output: %s", cwd)
					if strings.HasSuffix(cwd, "/") {
						saveFilePath = fmt.Sprintf("%s%s", cwd, path.Base(urlParsed.Path))
					} else {
						saveFilePath = fmt.Sprintf("%s/%s", cwd, path.Base(urlParsed.Path))
					}
				}
				if strings.HasSuffix(saveFilePath, "/") {
					saveFilePath = fmt.Sprintf("%s%s", saveFilePath, path.Base(urlParsed.Path))
				} else {
					saveFilePath = fmt.Sprintf("%s/%s", saveFilePath, path.Base(urlParsed.Path))
				}

				u.Log.Infof("Photo Found!. Saving as %s", saveFilePath)
				if len(saveFilePath) >= 0 && !u.FileExists(saveFilePath) {
					_, responseErr := photo.DownloadLarge2X(saveFilePath)
					if responseErr != nil {
						u.Log.Errorf("Error saving file number %d at %s. %s", photoNumber, saveFilePath, responseErr.Error())
					}
				}
			} else {
				u.Log.Errorf("Could not get photo number #%d !", photoNumber)
			}

		} else {
			log.Println(output.String())

			if len(saveFilePath) >= 0 && !u.FileExists(saveFilePath) {
				_, responseErr := pexels.SavePhotoByNumberToDisk(output, photoNumber, saveFilePath)
				if responseErr != nil {
					log.Errorf("Error saving file number 0 at %s. %s", saveFilePath, responseErr.Error())
				}
			}
		}
	},
}

func init() {
	pexelsCmd.AddCommand(pexelsGetimagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	pexelsGetimagesCmd.Flags().StringP("authorization", "a", "", "Pexels API Key used in Authorization request header")
	pexelsGetimagesCmd.Flags().StringP("query", "q", "beautiful", "Query to search for")
	pexelsGetimagesCmd.Flags().StringP("orientation", "o", "landscape", "(Optional) Desired photo orientation. The current supported orientations are: landscape, portrait or square")
	pexelsGetimagesCmd.Flags().StringP("size", "s", "large", "(Optional) Minimum photo size. The current supported sizes are: large(24MP), medium(12MP) or small(4MP)")
	pexelsGetimagesCmd.Flags().StringP("color", "c", "", "(Optional) Desired photo color. Supported colors: red, orange, yellow, green, turquoise, blue, violet, pink, brown, black, gray, white or any hexadecimal color code (eg. #ffffff)")
	pexelsGetimagesCmd.Flags().StringP("locale", "l", "", "The locale of the search you are performing. The current supported locales are: 'en-US' 'pt-BR' 'es-ES' 'ca-ES' 'de-DE' 'it-IT' 'fr-FR' 'sv-SE' 'id-ID' 'pl-PL' 'ja-JP' 'zh-TW' 'zh-CN' 'ko-KR' 'th-TH' 'nl-NL' 'hu-HU' 'vi-VN' 'cs-CZ' 'da-DK' 'fi-FI' 'uk-UA' 'el-GR' 'ro-RO' 'nb-NO' 'sk-SK' 'tr-TR' 'ru-RU'")
	pexelsGetimagesCmd.Flags().Uint64P("page", "p", 1, "The page number you are requesting. Default: 1")
	pexelsGetimagesCmd.Flags().Uint64P("max-results", "m", 15, "The number of results you are requesting per page. Default: 15 Max: 80")
	pexelsGetimagesCmd.Flags().StringP("save-file-path", "t", "./", "Path to the file in which to save the Image. Path should exist, file should not")
	pexelsGetimagesCmd.Flags().BoolP("randomize", "r", false, "Whether to randomize the output or not. Only outputs a single file at path")
}
