package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math"
	"strings"
)

const (
	PEXELS_URL_SEARCH            = "https://api.pexels.com/v1/search"
	PEXELS_IMAGE_LARGE           = "large"
	PEXELS_IMAGE_MEDIUM          = "medium"
	PEXELS_IMAGE_SMALL           = "small"
	PEXELS_ORIENTATION_LANDSCAPE = "landscape"
	PEXELS_ORIENTATION_PORTRAIT  = "portrait"
	PEXELS_ORIENTATION_SQUARE    = "square"
	PEXELS_LOCALE_DEFAULT        = "en-US"
)

var (
	Locales      = []string{"en-US", "pt-BR", "es-ES", "ca-ES", "de-DE", "it-IT", "fr-FR", "sv-SE", "id-ID", "pl-PL", "ja-JP", "zh-TW", "zh-CN", "ko-KR", "th-TH", "nl-NL", "hu-HU", "vi-VN", "cs-CZ", "da-DK", "fi-FI", "uk-UA", "el-GR", "ro-RO", "nb-NO", "sk-SK", "tr-TR", "ru-RU"}
	Colours      = []string{"red", "orange", "yellow", "green", "turquoise", "blue", "violet", "pink", "brown", "black", "gray", "white"}
	Sizes        = []string{PEXELS_IMAGE_SMALL, PEXELS_IMAGE_MEDIUM, PEXELS_IMAGE_LARGE}
	Orientations = []string{PEXELS_ORIENTATION_LANDSCAPE, PEXELS_ORIENTATION_PORTRAIT, PEXELS_ORIENTATION_SQUARE}
)

type Pexels struct {
	Authorization string
	Utils         *Utils
}

func (u *Utils) NewPexels(authorization string) *Pexels {
	return &Pexels{
		Authorization: authorization,
		Utils:         u,
	}
}

//func (p *Pexels) GetImages() {
//	client := resty.New()
//	resp, respErr := client.R().
//		SetHeader("Authorization", p.Authorization).
//		SetHeader("User-Agent", "SCLI/0.0.1").

//		SetBody(fmt.Sprintf("{\"username\":\"%s\", \"password\":\"%s\"}", c.Configs[clusterID].UserName, c.Configs[clusterID].Password)).
//		//SetResult(&respData).
//		Post(loginURL)
//
//	if respErr != nil {
//		c.Log.Errorf("Error logging-in to cronicle via go-resty. %s", respErr.Error())
//		return j, respErr
//	}
//	sessionData := string(resp.Body())
//}

type PexelsSearchInput struct {
	Query       string `json:"query"`
	Orientation string `json:"orientation"`
	Size        string `json:"size"`
	Color       string `json:"color"`
	Locale      string `json:"locale"`
	Page        uint64 `json:"page"`
	PerPage     uint64 `json:"per_page"`
}

func (p *PexelsSearchInput) String() string {
	return fmt.Sprintf("Query: %s\nOrientation: %s\nSize: %s\nColor: %s\nLocale: %s\nPage: %d\nPerPage: %d", p.Query, p.Orientation, p.Size, p.Color, p.Locale, p.Page, p.PerPage)
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large2X   string `json:"large2x"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type Photo struct {
	ID              int64       `json:"id"`
	Width           int64       `json:"width"`
	Height          int64       `json:"height"`
	URL             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerURL string      `json:"photographer_url"`
	PhotographerID  int64       `json:"photographer_id"`
	AvgColor        string      `json:"avg_color"`
	Src             PhotoSource `json:"src"`
	Liked           bool        `json:"liked"`
	Alt             string      `json:"alt"`
}

type PexelsSearchOutput struct {
	pexels       *Pexels
	SearchInput  PexelsSearchInput `json:"search_input"`
	Page         int64             `json:"page"`
	PerPage      int64             `json:"per_page"`
	Photos       []Photo           `json:"photos"`
	TotalResults int64             `json:"total_results"`
	NextPage     string            `json:"next_page"`
	Errors       []error           `json:"errors"`
	Error        string            `json:"error"`
}

func (p *PexelsSearchOutput) String() string {
	var retval string
	retval += fmt.Sprintf("Page: %d\nPerPage: %d\nPhotos: %d\n", p.Page, p.PerPage, len(p.Photos))
	for k, photo := range p.Photos {
		retval += fmt.Sprintf("Photo #%d\n----------\nID: %d\nWidth: %d\nHeight: %d\nURL: %s\nPhotographer: %s\nPhotographerURL: %s\nPhotographerID: %d\nAvgColor: %s\nSrc.Original: %s\nSrc.Large2X: %s\nSrc.Large: %s\nSrc.Medium: %s\nSrc.Small: %s\nSrc.Portrait: %s\nSrc.Landscape: %s,Src.Tiny: %s\nLiked: %t\nAlt: %s\n",
			k+1, photo.ID, photo.Width, photo.Height, photo.URL, photo.Photographer, photo.PhotographerURL, photo.PhotographerID, photo.AvgColor, photo.Src.Original, photo.Src.Large2X, photo.Src.Large, photo.Src.Medium, photo.Src.Small, photo.Src.Portrait, photo.Src.Landscape, photo.Src.Tiny, photo.Liked, photo.Alt)
	}
	retval += fmt.Sprintf("TotalResults: %d\nNextPage: %s\nErrors \n----------", p.TotalResults, p.NextPage)
	for k, err := range p.Errors {
		retval += fmt.Sprintf("\nError Text #%d: %s\n", k, err.Error())
	}
	return retval
}

func (p *PexelsSearchOutput) GetPhotoByNumber(photoNumber int64) (*Photo, error) {

	if photoNumber <= int64(len(p.Photos)) {
		if photoNumber > 0 {
			photoNumber--
		}
		for k, photo := range p.Photos {
			if int64(k) == photoNumber {
				return &photo, nil
			}
		}
	}

	if photoNumber > int64(len(p.Photos)) && photoNumber <= p.TotalResults {
		return p.SearchPhotoByNumber(photoNumber)
	}

	return nil, errors.New(fmt.Sprintf("PexelsSearchOutput.GetPhotoByNumber::Photo Number (%d) is neither <= %d nor > %d but <= %d", photoNumber, len(p.Photos), len(p.Photos), p.TotalResults))
}

func (p *PexelsSearchOutput) SearchPhotoByNumber(photoNumber int64) (*Photo, error) {
	if photoNumber <= p.TotalResults {
		p.SearchInput.Page = uint64(math.Ceil(float64(uint64(photoNumber) / p.SearchInput.PerPage)))
		return p.pexels.SearchPhotoByNumber(p.SearchInput, photoNumber)
	}
	return nil, errors.New(fmt.Sprintf("PexelsSearchOutput.SearchPhotoByNumber::Photo Number (%d) is > TotalResults (%d)", photoNumber, p.TotalResults))
}

func (p *Pexels) SavePhotoByIDToDisk(searchResults PexelsSearchOutput, photoID int64, targetFilePath string) (*resty.Response, error) {
	if p.Utils.FileExists(targetFilePath) {
		return nil, errors.New(fmt.Sprintf("Target file %s already exists", targetFilePath))
	}

	for _, photo := range searchResults.Photos {
		if photo.ID == photoID {
			return p.Utils.CreateFileWithDataAtURL(photo.Src.Landscape, targetFilePath)
		}
	}

	return nil, errors.New(fmt.Sprintf("Photo with ID %d not found in search results provided", photoID))
}

func (p *Pexels) SavePhotoByNumberToDisk(searchResults PexelsSearchOutput, photoNumber int64, targetFilePath string) (*resty.Response, error) {
	if p.Utils.FileExists(targetFilePath) {
		return nil, errors.New(fmt.Sprintf("Target file %s already exists", targetFilePath))
	}

	photo, photoErr := searchResults.GetPhotoByNumber(photoNumber)

	if photoErr != nil {
		return nil, photoErr
	}

	if photo != nil {
		return p.Utils.CreateFileWithDataAtURL(photo.Src.Original, targetFilePath)
	}

	return nil, errors.New(fmt.Sprintf("Photo at sequence number %d not found in search results provided", photoNumber))
}

func (p *Pexels) GetPhotoURLToDownload(searchResults PexelsSearchOutput, photoNumber int, targetFilePath string) (*resty.Response, error) {
	if p.Utils.FileExists(targetFilePath) {
		return nil, errors.New(fmt.Sprintf("Target file %s already exists", targetFilePath))
	}

	for k, photo := range searchResults.Photos {
		if k == photoNumber {
			return p.Utils.CreateFileWithDataAtURL(photo.Src.Original, targetFilePath)
		}
	}

	return nil, errors.New(fmt.Sprintf("Photo at sequence number %d not found in search results provided", photoNumber))
}

func (p *Pexels) Search(input PexelsSearchInput) PexelsSearchOutput {

	var output PexelsSearchOutput

	output.pexels = p

	if input.PerPage > 80 {
		input.PerPage = 80
	}

	output.SearchInput = input

	client := resty.New()
	resp, respErr := client.R().
		SetHeader("Authorization", p.Authorization).
		SetHeader("User-Agent", "SCLI/0.0.1").
		SetQueryParam("query", input.Query).
		SetQueryParam("orientation", input.Orientation).
		SetQueryParam("size", input.Size).
		SetQueryParam("color", input.Color).
		SetQueryParam("locale", input.Locale).
		SetQueryParam("page", fmt.Sprintf("%d", input.Page)).
		SetQueryParam("per_page", fmt.Sprintf("%d", input.PerPage)).
		Get(PEXELS_URL_SEARCH)

	if respErr != nil {
		p.Utils.Log.Errorf("Error fetching data from pexels API via go-resty. %s", respErr.Error())
		return output
	}

	respCode := resp.StatusCode()

	if respCode > 299 || respCode < 200 {
		output.Errors = append(output.Errors, errors.New(fmt.Sprintf("Response Status Code: %d, Response Status: %s", respCode, resp.String())))
	}

	headers := resp.Header()
	//p.Utils.Log.Infof("Printing Response Headers...")
	//for k, v := range headers {
	//	p.Utils.Log.Infof("Header: %s, Value: %#v", k, v)
	//}

	if strings.Contains(headers["Content-Type"][0], "application/json") {
		outputErr := json.Unmarshal(resp.Body(), &output)
		if outputErr != nil {
			p.Utils.Log.Errorf("Error unmarshalling JSON response. %s", outputErr.Error())
		}
	}

	return output
}

func (p *Pexels) SearchPhotoByNumber(input PexelsSearchInput, photoNumber int64) (*Photo, error) {
	input.Page = uint64(math.Ceil(float64(uint64(photoNumber) / input.PerPage)))
	searchResults := p.Search(input)
	photoNumberInResultsPage := uint64(photoNumber) - (input.PerPage * input.Page)

	for keyPhotoNumber, photo := range searchResults.Photos {
		if uint64(keyPhotoNumber) == photoNumberInResultsPage {
			return &photo, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Pexels.SearchPhotoByNumber::Could not find photo number #%d (#%d in results page #%d) using input:\n %s", photoNumber, photoNumberInResultsPage, input.Page, input.String()))
}

func (p *Photo) Download(sourceFileURL, targetFilePath string) (*resty.Response, error) {
	client := resty.New()
	return client.R().SetOutput(targetFilePath).Get(sourceFileURL)
}

func (p *Photo) DownloadTiny(targetFilePath string) (*resty.Response, error) {
	return p.Download(p.Src.Tiny, targetFilePath)
}

func (p *Photo) DownloadSmall(targetFilePath string) (*resty.Response, error) {
	return p.Download(p.Src.Small, targetFilePath)
}

func (p *Photo) DownloadMedium(targetFilePath string) (*resty.Response, error) {
	return p.Download(p.Src.Medium, targetFilePath)
}

func (p *Photo) DownloadLarge(targetFilePath string) (*resty.Response, error) {
	return p.Download(p.Src.Large, targetFilePath)
}

func (p *Photo) DownloadLarge2X(targetFilePath string) (*resty.Response, error) {
	return p.Download(p.Src.Large2X, targetFilePath)
}

func (p *Photo) DownloadOriginal(targetFilePath string) (*resty.Response, error) {
	return p.Download(p.Src.Original, targetFilePath)
}

func (p *Photo) DownloadLandscape(targetFilePath string) (*resty.Response, error) {
	return p.Download(p.Src.Landscape, targetFilePath)
}

func (p *Photo) DownloadPortrait(targetFilePath string) (*resty.Response, error) {
	return p.Download(p.Src.Portrait, targetFilePath)
}
