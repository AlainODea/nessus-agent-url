package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/coreos/go-semver/semver"
)

func main() {
	substr := os.Args[1]
	downloadPageUrl := "https://www.tenable.com/downloads/nessus-agents?loginAttempted=true"
	downloadUrl, err := GetMatchingDownloadUrl(substr, LoadDownloadPage(downloadPageUrl))
	if err != nil {
		panic(err)
	}
	fmt.Println(downloadUrl)
}

func LoadDownloadPage(url string) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != 200 {
			return nil, StatusCodeError(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
		}
		return res.Body, nil
	}
}

type StatusCodeError string

func (s StatusCodeError) Error() string {
	return string(s)
}

func GetMatchingDownloadUrl(substr string, pageLoader func() (io.ReadCloser, error)) (string, error) {
	doc, err := GetDownloadsDocument(pageLoader)
	if err != nil {
		return "", err
	}
	agentDownloads := GetDownloadsFromDocument(doc)
	latestMatchingDownload := GetLatestMatchingDownload(agentDownloads, substr)
	downloadUrl := GetDownloadUrl(latestMatchingDownload)
	return downloadUrl, nil
}

func GetDownloadsDocument(pageLoader func() (io.ReadCloser, error)) (*goquery.Document, error) {
	reader, err := pageLoader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func GetDownloadsFromDocument(doc *goquery.Document) *[]DownloadItem {
	var nextData NextData
	doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		nextDataJson := s.Text()
		err := json.Unmarshal([]byte(nextDataJson), &nextData)
		if err != nil {
			panic(fmt.Sprintf("Could not unmarshal the downloads JSON: \"%v\"", err))
		}
	})
	return nextData.Props.PageProps.Page.Downloads
}

func GetLatestMatchingDownload(agentDownloads *[]DownloadItem, substr string) DownloadItem {
	var amazonLinuxDownloads []DownloadItem
	for _, downloadItem := range *agentDownloads {
		if strings.Contains(downloadItem.Name, substr) {
			amazonLinuxDownloads = append(amazonLinuxDownloads, downloadItem)
		}
	}
	sort.Sort(DownloadItemsByVersion(amazonLinuxDownloads))
	return amazonLinuxDownloads[len(amazonLinuxDownloads)-1]
}

func GetDownloadUrl(download DownloadItem) string {
	return fmt.Sprintf("https://www.tenable.com/downloads/api/v1/public/pages/nessus-agents/downloads/%v/download?i_agree_to_tenable_license_agreement=true", download.Id)
}

type DownloadItem struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MetaData    MetaData `json:"meta_data"`
}

type MetaData struct {
	Md5         string `json:"md5"`
	Sha256      string `json:"sha256"`
	Product     string `json:"product"`
	Version     string `json:"version"`
	ReleaseDate string `json:"release_date"`
}

type NextData struct {
	Props struct {
		PageProps struct {
			Page struct {
				Downloads *[]DownloadItem `json:"downloads"`
			} `json:"page"`
		} `json:"pageProps"`
	} `json:"props"`
}

type DownloadItemsByVersion []DownloadItem

func (e DownloadItemsByVersion) Len() int {
	return len(e)
}

func (e DownloadItemsByVersion) Less(i, j int) bool {
	iVersion, err := semver.NewVersion(e[i].MetaData.Version)
	if err != nil {
		panic(err)
	}
	jVersion, err := semver.NewVersion(e[j].MetaData.Version)
	if err != nil {
		panic(err)
	}
	return iVersion.LessThan(*jVersion)
}

func (e DownloadItemsByVersion) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
