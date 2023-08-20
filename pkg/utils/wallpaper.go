package utils

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/shammishailaj/scli/pkg/utils/storage"
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

		saveFileName := fmt.Sprintf("SCLI_WALLPAPER_%d_%s_*.%s", time.Now().Nanosecond(), path.Base(urlParsed.Path), path.Ext(urlParsed.Path))

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

func (u *Utils) RandomPexelsWallpaperWithCacheAndDB(authorization, query, orientation, size, color, locale, cacheDir, dbFilePath string) error {
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

		saveFileName := fmt.Sprintf("SCLI_WALLPAPER_%d_%s_*%s", time.Now().Nanosecond(), path.Base(urlParsed.Path), path.Ext(urlParsed.Path))

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
		saveFilePath := saveFileHandle.Name()

		u.Log.Infof("Photo Found!. Saving as %s OR %s", saveFileHandle.Name(), saveFilePath)

		saveMetaDataErr := photo.SaveMetaDataToGenjiDB(dbFilePath, "scli_random_wallpaper_metadata", saveFilePath, "large2x", input.Query)

		if saveMetaDataErr != nil {
			return saveMetaDataErr
		}

		if len(saveFilePath) >= 0 {
			_, responseErr := photo.DownloadLarge2X(saveFilePath)
			if responseErr != nil {
				return errors.New(fmt.Sprintf("Error saving file number %d at %s. %s", photoNumber, saveFilePath, responseErr.Error()))
			}

			modifyWallpaperErr := u.ModifyWallpaper("gnome", saveFilePath)
			if modifyWallpaperErr != nil {
				return modifyWallpaperErr
			}
			return u.SaveWallpaperChangeLog(dbFilePath, "scli_random_wallpaper_change_log", saveFilePath, photoNumber)
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

func (u *Utils) SaveWallpaperChangeLog(dbPath, tableName, wallpaperPath string, randomNumber int64) error {

	store, storageErr := storage.NewStore(map[string]string{"type": "genji"})
	if storageErr != nil {
		return storageErr
	}
	dbErr := store.Connect(dbPath, context.Background())
	if dbErr != nil {
		return dbErr
	}

	tableDef := `
		yyyy INT NOT NULL,
		mm INT NOT NULL,
		dd INT NOT NULL,
		hh INT NOT NULL,
		mi INT NOT NULL,
		ss INT NOT NULL,
		ms INT NOT NULL,
		us INT NOT NULL,
		ns INT NOT NULL,
		tz TEXT NOT NULL,
		random_number INT NOT NULL,
		filepath TEXT NOT NULL,
		PRIMARY KEY(yyyy,mm,dd,hh,ss,ms,us,ns,tz)
		`

	_, err := store.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s)", tableName, tableDef))
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s(yyyy,mm,dd,hh,mi,ss,ms,us,ns,tz,random_number,filepath) VALUES(", tableName)
	query += fmt.Sprintf("%d,%d,%d,", time.Now().Year(), time.Now().Month(), time.Now().Day())
	query += fmt.Sprintf("%d,%d,%d,", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	query += fmt.Sprintf("%d,%d,%d,%d", time.Now().Second(), time.Now().UnixMilli(), time.Now().UnixMicro(), time.Now().Nanosecond())
	query += fmt.Sprintf("'%s',%d,'%s'", time.Now().Location().String(), randomNumber, wallpaperPath)
	_, err = store.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (u *Utils) RandomPexelsWallpaperFromCache(cacheDir string) error {
	if cacheDir != "" {
		cacheDir += "/scli/random/wallpaper"
	}

	absCacheDir, absCacheDirErr := u.GetAbsolutePath(cacheDir)
	if absCacheDirErr != nil {
		return fmt.Errorf("utils.RandomPexelsWallpaperFromCache: error getting absolute path for %s: %s", cacheDir, absCacheDirErr.Error())
	}

	files, err := os.ReadDir(absCacheDir)
	if err != nil {
		return fmt.Errorf("utils.RandomPexelsWallpaperFromCache: error reading directory: %s: %s", absCacheDir, err.Error())
	}

	filesLen := len(files)

	if filesLen > 0 {
		u.Log.Infof("utils.RandomPexelsWallpaperFromCache: Getting a random photo...")
		u.Log.Infof("utils.RandomPexelsWallpaperFromCache: Generating a random photo number...")
		randomPhotoNumber, randomPhotoNumberErr := rand.Int(rand.Reader, big.NewInt(int64(filesLen)))
		if randomPhotoNumberErr != nil {
			u.Log.Errorf("utils.RandomPexelsWallpaperFromCache: error getting random number between 0 and %d", filesLen)
			u.Log.Infof("utils.RandomPexelsWallpaperFromCache: choosing the 1st image at position 0")
			randomPhotoNumber = big.NewInt(0)
		}

		photoNumber := randomPhotoNumber.Int64()
		u.Log.Infof("utils.RandomPexelsWallpaperFromCache: getting photo number #%d...", photoNumber)

		for i := photoNumber - 1; i < int64(filesLen); i++ {
			fileName := fmt.Sprintf("%s/%s", absCacheDir, files[i].Name())

			if files[i].IsDir() {
				continue
			}

			mimeType, mimeTypeErr := u.GetFileContentType(fileName)
			if mimeTypeErr != nil {
				u.Log.Errorf("utils.RandomPexelsWallpaperFromCache: unable to get content-type for %s: %s", fileName, mimeTypeErr.Error())
				u.Log.Infof("utils.RandomPexelsWallpaperFromCache: skipping to next file...")
				continue
			}

			if u.InArray(mimeType, []string{"image/jpeg", "image/png", "image/gif", "image/bmp", "image/tiff", "image/webp", "image/svg+xml"}) {
				modifyWallpaperErr := u.ModifyWallpaper("gnome", fileName)
				if modifyWallpaperErr != nil {
					return fmt.Errorf("utils.RandomPexelsWallpaperFromCache: error setting wallpaper to %s: %s", fileName, modifyWallpaperErr.Error())
				}
				return nil
			}
		}
	}
	return fmt.Errorf("utils.RandomPexelsWallpaperFromCache: could not find any wallpapers at %s", absCacheDir)
}
