package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/big"
	"net/url"
	"os"
	"path"
	"time"
)

func (u *Utils) ModifyWallpaper(envName, wallpaperPath string) error {
	commands := make(map[string]map[string][]string)

	gnomeCommands := make(map[string][]string)
	gnomeCommands["modify-wallpaper-gnome"] = []string{"gsettings", "set", "org.gnome.desktop.background", "picture-uri", fmt.Sprintf("file://%s", wallpaperPath)}
	gnomeCommands["modify-wallpaper-gnome-dark"] = []string{"gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", fmt.Sprintf("file://%s", wallpaperPath)}

	commands["gnome"] = gnomeCommands

	loopExecuted := false

	for commandName, command := range commands[envName] {
		u.Log.Infof("Running Command for %s", commandName)
		output, err := u.ExecuteCommand(command)

		if err != nil {
			return errors.New(fmt.Sprintf("Error running shell command %s. %s", commandName, err.Error()))
		}

		u.Log.Infof("Output:\n%s", output)
		loopExecuted = true
	}

	if loopExecuted {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Invalid Environment Name: %s", envName))
	}
}

func (u *Utils) RandomPexelsWallpaper(authorization, query, orientation, size, color, locale string) error {
	input := PexelsSearchInput{
		Query:       "people",
		Orientation: PEXELS_ORIENTATION_LANDSCAPE,
		Size:        "large",
		Color:       "",
		Locale:      "",
		Page:        1,
		PerPage:     15,
	}

	if authorization == "" {
		return errors.New(fmt.Sprintf("Can not proceed without a valid Pexels api key. Provided: %s\n", authorization))
	}

	if len(authorization) < 32 {
		return errors.New(fmt.Sprintf("Can not proceed without a valid api key....\n"))
	}

	if query == "" {
		return errors.New(fmt.Sprintf("Can not proceed without a valid search query."))
	}

	if len(query) < 3 {
		return errors.New(fmt.Sprintf("Can not proceed without a valid search query. Must be 3 characters or more. Provided: %s", query))
	}

	input.Query = query

	if orientation == "" {
		log.Infof("Value not provided for orientation. Will use: landscape")
		switch orientation {
		case PEXELS_ORIENTATION_LANDSCAPE:
			input.Orientation = PEXELS_ORIENTATION_LANDSCAPE
		case PEXELS_ORIENTATION_PORTRAIT:
			input.Orientation = PEXELS_ORIENTATION_PORTRAIT
		case PEXELS_ORIENTATION_SQUARE:
			input.Orientation = PEXELS_ORIENTATION_SQUARE
		default:
			input.Orientation = PEXELS_ORIENTATION_LANDSCAPE
		}
	}

	if size == "" {
		u.Log.Infof("Value not provided for size. Will use: large")
		switch size {
		case PEXELS_IMAGE_LARGE:
			input.Size = PEXELS_IMAGE_LARGE
		case PEXELS_IMAGE_MEDIUM:
			input.Size = PEXELS_IMAGE_MEDIUM
		case PEXELS_IMAGE_SMALL:
			input.Size = PEXELS_IMAGE_SMALL
		default:
			input.Size = PEXELS_IMAGE_LARGE
		}
	}

	if color == "" {
		u.Log.Infof("Color not provided. Will not send a color in request")
	}

	// TODO: Check this function to see if color is valid https://go.dev/play/p/rM0e-w7Xfdg
	// TODO: Original function at: https://www.geeksforgeeks.org/check-if-a-given-string-is-a-valid-hexadecimal-color-code-or-not/
	if !u.InArray(color, Colours) {
		u.Log.Infof("Invalid color. Will not send a color in request")
		color = ""
	}

	input.Color = color

	if locale == "" {
		u.Log.Infof("Value not provided for locale. Will use \"en-US\"")
		locale = PEXELS_LOCALE_DEFAULT
	}

	if !u.InArray(locale, Locales) {
		u.Log.Infof("Illegal value for locale. Will use \"en-US\"")
		locale = PEXELS_LOCALE_DEFAULT
	}

	input.Locale = locale

	pexels := u.NewPexels(authorization)
	u.Log.Infof("Sending request to Pexels search API with input...\n%s", input.String())
	output := pexels.Search(input)
	u.Log.Infof("Output Received....\n%s", output.String())
	photoNumber := int64(1)

	if output.TotalResults <= 0 {
		return errors.New(fmt.Sprintf("Nothing found for Search Input\n%s", input.String()))
	}

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
		return errors.New(fmt.Sprintf("Error getting photo number %d. %s", photoNumber, photoErr.Error()))
	}
	if photo != nil {
		u.Log.Infof("Photo Found. Will download Large2X size at URL: %s", photo.Src.Large2X)
		urlParsed, urlErr := url.Parse(photo.Src.Large2X)
		if urlErr != nil {
			u.Log.Printf("Error parsing URL %s. %s", photo.Src.Large2X, urlErr.Error())
			urlParsed = &url.URL{Path: photo.Src.Large2X}
		}

		saveFileName := fmt.Sprintf("SCLI_WALLPAPER_%d_%s_*", time.Now().Nanosecond(), path.Base(urlParsed.Path))
		saveFileHandle, saveFileHandleErr := u.TempFile(saveFileName)
		if saveFileHandleErr != nil {
			return errors.New(fmt.Sprintf("Error creating temporary file to save downloaded wallpaper. %s", saveFileHandleErr.Error()))
		}
		//saveFilePath, saveFilePathErr := filepath.Abs(filepath.Dir(saveFileHandle.Name()))
		//if saveFilePathErr != nil {
		//	u.Log.Errorf("Unable to get filepath of temp file created. %s", saveFilePathErr.Error())
		//	saveFilePath = "Unable to get filepath of temp file created"
		//}
		saveFilePath := saveFileHandle.Name()

		u.Log.Infof("Photo Found!. Saving as %s OR %s", saveFileHandle.Name(), saveFilePath)
		if len(saveFilePath) >= 0 {
			_, responseErr := photo.DownloadLarge2X(saveFilePath)
			if responseErr != nil {
				return errors.New(fmt.Sprintf("Error saving file number %d at %s. %s", photoNumber, saveFilePath, responseErr.Error()))
			}

			return u.ModifyWallpaper("gnome", saveFilePath)

		}
	}

	return errors.New(fmt.Sprintf("Could not get photo number #%d !", photoNumber))

	//else {
	//	log.Println(output.String())
	//
	//	if len(saveFilePath) >= 0 && !u.FileExists(saveFilePath) {
	//		_, responseErr := pexels.SavePhotoByNumberToDisk(output, photoNumber, saveFilePath)
	//		if responseErr != nil {
	//			log.Errorf("Error saving file number 0 at %s. %s", saveFilePath, responseErr.Error())
	//		}
	//	}
	//}
}

func (u *Utils) RandomPexelsWallpaperWithCache(authorization, query, orientation, size, color, locale, cacheDir string) error {
	input := PexelsSearchInput{
		Query:       "people",
		Orientation: PEXELS_ORIENTATION_LANDSCAPE,
		Size:        "large",
		Color:       "",
		Locale:      "",
		Page:        1,
		PerPage:     15,
	}

	if authorization == "" {
		return errors.New(fmt.Sprintf("Can not proceed without a valid Pexels api key. Provided: %s\n", authorization))
	}

	if len(authorization) < 32 {
		return errors.New(fmt.Sprintf("Can not proceed without a valid api key....\n"))
	}

	if query == "" {
		return errors.New(fmt.Sprintf("Can not proceed without a valid search query."))
	}

	if len(query) < 3 {
		return errors.New(fmt.Sprintf("Can not proceed without a valid search query. Must be 3 characters or more. Provided: %s", query))
	}

	input.Query = query

	if orientation == "" {
		log.Infof("Value not provided for orientation. Will use: landscape")
		switch orientation {
		case PEXELS_ORIENTATION_LANDSCAPE:
			input.Orientation = PEXELS_ORIENTATION_LANDSCAPE
		case PEXELS_ORIENTATION_PORTRAIT:
			input.Orientation = PEXELS_ORIENTATION_PORTRAIT
		case PEXELS_ORIENTATION_SQUARE:
			input.Orientation = PEXELS_ORIENTATION_SQUARE
		default:
			input.Orientation = PEXELS_ORIENTATION_LANDSCAPE
		}
	}

	if size == "" {
		u.Log.Infof("Value not provided for size. Will use: large")
		switch size {
		case PEXELS_IMAGE_LARGE:
			input.Size = PEXELS_IMAGE_LARGE
		case PEXELS_IMAGE_MEDIUM:
			input.Size = PEXELS_IMAGE_MEDIUM
		case PEXELS_IMAGE_SMALL:
			input.Size = PEXELS_IMAGE_SMALL
		default:
			input.Size = PEXELS_IMAGE_LARGE
		}
	}

	if color == "" {
		u.Log.Infof("Color not provided. Will not send a color in request")
	}

	// TODO: Check this function to see if color is valid https://go.dev/play/p/rM0e-w7Xfdg
	// TODO: Original function at: https://www.geeksforgeeks.org/check-if-a-given-string-is-a-valid-hexadecimal-color-code-or-not/
	if !u.InArray(color, Colours) {
		u.Log.Infof("Invalid color. Will not send a color in request")
		color = ""
	}

	input.Color = color

	if locale == "" {
		u.Log.Infof("Value not provided for locale. Will use \"en-US\"")
		locale = PEXELS_LOCALE_DEFAULT
	}

	if !u.InArray(locale, Locales) {
		u.Log.Infof("Illegal value for locale. Will use \"en-US\"")
		locale = PEXELS_LOCALE_DEFAULT
	}

	input.Locale = locale

	pexels := u.NewPexels(authorization)
	u.Log.Infof("Sending request to Pexels search API with input...\n%s", input.String())
	output := pexels.Search(input)
	u.Log.Infof("Output Received....\n%s", output.String())
	photoNumber := int64(1)

	if output.TotalResults <= 0 {
		return errors.New(fmt.Sprintf("Nothing found for Search Input\n%s", input.String()))
	}

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
		return errors.New(fmt.Sprintf("Error getting photo number %d. %s", photoNumber, photoErr.Error()))
	}
	if photo != nil {
		u.Log.Infof("Photo Found. Will download Large2X size at URL: %s", photo.Src.Large2X)
		urlParsed, urlErr := url.Parse(photo.Src.Large2X)
		if urlErr != nil {
			u.Log.Printf("Error parsing URL %s. %s", photo.Src.Large2X, urlErr.Error())
			urlParsed = &url.URL{Path: photo.Src.Large2X}
		}

		saveFileName := fmt.Sprintf("SCLI_WALLPAPER_%d_%s_*", time.Now().Nanosecond(), path.Base(urlParsed.Path))
		saveFileNameExtension := path.Ext(urlParsed.Path)

		if cacheDir != "" {
			cacheDir += "/scli/random/wallpaper"
		}

		absCacheDir, absCacheDirErr := u.GetAbsolutePath(cacheDir)
		if absCacheDirErr != nil {
			return absCacheDirErr
		}

		if !u.FileExists(absCacheDir) {
			mkdirAllErr := os.MkdirAll(absCacheDir, 0700)
			if mkdirAllErr != nil {
				return mkdirAllErr
			}
		}

		saveFileHandle, saveFileHandleErr := u.TempFileAtDir(absCacheDir, saveFileName)

		if saveFileHandleErr != nil {
			return errors.New(fmt.Sprintf("Error creating temporary file to save downloaded wallpaper. %s", saveFileHandleErr.Error()))
		}
		//saveFilePath, saveFilePathErr := filepath.Abs(filepath.Dir(saveFileHandle.Name()))
		//if saveFilePathErr != nil {
		//	u.Log.Errorf("Unable to get filepath of temp file created. %s", saveFilePathErr.Error())
		//	saveFilePath = "Unable to get filepath of temp file created"
		//}
		saveFilePath := saveFileHandle.Name() + saveFileNameExtension

		u.Log.Infof("Photo Found!. Saving as %s OR %s", saveFileHandle.Name(), saveFilePath)
		if len(saveFilePath) >= 0 {
			_, responseErr := photo.DownloadLarge2X(saveFilePath)
			if responseErr != nil {
				return errors.New(fmt.Sprintf("Error saving file number %d at %s. %s", photoNumber, saveFilePath, responseErr.Error()))
			}

			return u.ModifyWallpaper("gnome", saveFilePath)

		}
	}

	return errors.New(fmt.Sprintf("Could not get photo number #%d !", photoNumber))

	//else {
	//	log.Println(output.String())
	//
	//	if len(saveFilePath) >= 0 && !u.FileExists(saveFilePath) {
	//		_, responseErr := pexels.SavePhotoByNumberToDisk(output, photoNumber, saveFilePath)
	//		if responseErr != nil {
	//			log.Errorf("Error saving file number 0 at %s. %s", saveFilePath, responseErr.Error())
	//		}
	//	}
	//}
}
