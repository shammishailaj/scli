package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
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
	Query       string
	Orientation string
	Size        string
	Color       string
	Locale      string
	Page        uint
	PerPage     uint
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
	ID              int         `json:"id"`
	Width           int         `json:"width"`
	Height          int         `json:"height"`
	URL             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerURL string      `json:"photographer_url"`
	PhotographerID  int         `json:"photographer_id"`
	AvgColor        string      `json:"avg_color"`
	Src             PhotoSource `json:"src"`
	Liked           bool        `json:"liked"`
	Alt             string      `json:"alt"`
}

type PexelsSearchOutput struct {
	Page         int     `json:"page"`
	PerPage      int     `json:"per_page"`
	Photos       []Photo `json:"photos"`
	TotalResults int     `json:"total_results"`
	NextPage     string  `json:"next_page"`
}

func (p *PexelsSearchOutput) String() string {
	var retval string
	retval += fmt.Sprintf("Page: %d\nPerPage: %d\nPhotos: %d\n", p.Page, p.PerPage, len(p.Photos))
	for k, photo := range p.Photos {
		retval += fmt.Sprintf("Photo #%d\n----------\nID: %d\nWidth: %d\nHeight: %d\nURL: %s\nPhotographer: %s\nPhotographerURL: %s\nPhotographerID: %d\nAvgColor: %s\nSrc.Original: %s\nSrc.Large2X: %s\nSrc.Large: %s\nSrc.Medium: %s\nSrc.Small: %s\nSrc.Portrait: %s\nSrc.Landscape: %s,Src.Tiny: %s\nLiked: %t\nAlt: %s\n",
			k+1, photo.ID, photo.Width, photo.Height, photo.URL, photo.Photographer, photo.PhotographerURL, photo.PhotographerID, photo.AvgColor, photo.Src.Original, photo.Src.Large2X, photo.Src.Large, photo.Src.Medium, photo.Src.Small, photo.Src.Portrait, photo.Src.Landscape, photo.Src.Tiny, photo.Liked, photo.Alt)
	}
	retval += fmt.Sprintf("TotalResults: %d\nNextPage: %s\n", p.TotalResults, p.NextPage)
	return retval
}

func (p *Pexels) SavePhotoByIDToDisk(searchResults PexelsSearchOutput, photoID int, targetFilePath string) (*resty.Response, error) {
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

func (p *Pexels) SavePhotoByNumberToDisk(searchResults PexelsSearchOutput, photoNumber int, targetFilePath string) (*resty.Response, error) {
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

	if input.PerPage > 80 {
		input.PerPage = 80
	}

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
