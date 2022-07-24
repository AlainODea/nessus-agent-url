package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"os"
	"sort"
	"testing"
)

func TestLoadDownloadPage_InvalidDomain(t *testing.T) {
	loadDownloadPage := LoadDownloadPage("https://invalid.")
	page, err := loadDownloadPage()
	if page != nil {
		t.Errorf("got %v, wanted nil", page)
	}

	got := err.Error()
	want := "Get \"https://invalid.\": dial tcp: lookup invalid.: no such host"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestLoadDownloadPage_404NotFound(t *testing.T) {
	loadDownloadPage := LoadDownloadPage("https://httpbin.org/status/404")
	page, err := loadDownloadPage()
	if page != nil {
		t.Errorf("got %v, wanted nil", page)
	}

	got := err.Error()
	want := "status code error: 404 404 Not Found"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestLoadDownloadPage_200OK(t *testing.T) {
	loadDownloadPage := LoadDownloadPage("https://httpbin.org/status/200")
	page, err := loadDownloadPage()
	if page == nil {
		t.Errorf("got nil, wanted non-nil page")
	}
	if err != nil {
		t.Errorf("got %v, wanted nil", err)
	}
}

func TestGetMatchingDownloadUrl(t *testing.T) {
	fileReader := func() (io.ReadCloser, error) {
		file, err := os.Open("test/nessus-agents.html")
		if err != nil {
			fmt.Println("Error opening file!!!")
		}
		return file, nil
	}
	downloadUrl, err := GetMatchingDownloadUrl("-amzn.x86_64.rpm", fileReader)
	if err != nil {
		panic(err)
	}

	got := downloadUrl
	want := "https://www.tenable.com/downloads/api/v1/public/pages/nessus-agents/downloads/16733/download?i_agree_to_tenable_license_agreement=true"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetMatchingDownloadUrl_AlreadyClosed(t *testing.T) {
	fileReader := func() (io.ReadCloser, error) {
		file, err := os.Open("test/nessus-agents.html")
		if err != nil {
			fmt.Println("Error opening file!!!")
		}
		err = file.Close()
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	downloadUrl, err := GetMatchingDownloadUrl("-amzn.x86_64.rpm", fileReader)
	if downloadUrl != "" {
		t.Errorf("got %q, wanted \"\"", downloadUrl)
	}

	got := err.Error()
	want := "read test/nessus-agents.html: file already closed"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetDownloadUrl(t *testing.T) {
	newestDownloadItem := DownloadItem{
		Id:          16733,
		Name:        "NessusAgent-10.1.4-amzn.x86_64.rpm",
		Description: "Amazon Linux 2015.03, 2015.09, 2017.09, 2018.03, Amazon Linux 2 (x86_64)",
		MetaData: MetaData{
			Md5:         "b25c5e1a033eed78f7f2366254748bfb",
			Sha256:      "47dabe3c313c4026ad318cac25f88f31c4f55ec613a980e614f9f1223f355b40",
			Product:     "Nessus Agents - 10.1.4",
			Version:     "10.1.4",
			ReleaseDate: "2022-06-11T00:00:00.000Z",
		},
	}

	downloadUrl := GetDownloadUrl(newestDownloadItem)

	got := downloadUrl
	want := "https://www.tenable.com/downloads/api/v1/public/pages/nessus-agents/downloads/16733/download?i_agree_to_tenable_license_agreement=true"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetDownloadsDocument(t *testing.T) {
	fileReader := func() (io.ReadCloser, error) {
		file, err := os.Open("test/nessus-agents.html")
		if err != nil {
			fmt.Println("Error opening file!!!")
		}
		return file, nil
	}
	doc, err := GetDownloadsDocument(fileReader)

	if doc == nil {
		t.Errorf("got nil, wanted non-nil doc")
	}

	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
}

func TestGetDownloadsDocument_DoesNotExist(t *testing.T) {
	fileReader := func() (io.ReadCloser, error) {
		return os.Open("test/does-not-exist.html")
	}
	doc, err := GetDownloadsDocument(fileReader)

	if doc != nil {
		t.Errorf("got %v, wanted nil", doc)
	}

	got := err.Error()
	want := "open test/does-not-exist.html: no such file or directory"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetLatestMatchingDownload(t *testing.T) {
	newestDownloadItem := DownloadItem{
		Id:          16733,
		Name:        "NessusAgent-10.1.4-amzn.x86_64.rpm",
		Description: "Amazon Linux 2015.03, 2015.09, 2017.09, 2018.03, Amazon Linux 2 (x86_64)",
		MetaData: MetaData{
			Md5:         "b25c5e1a033eed78f7f2366254748bfb",
			Sha256:      "47dabe3c313c4026ad318cac25f88f31c4f55ec613a980e614f9f1223f355b40",
			Product:     "Nessus Agents - 10.1.4",
			Version:     "10.1.4",
			ReleaseDate: "2022-06-11T00:00:00.000Z",
		},
	}

	fileReader := func() (io.ReadCloser, error) {
		file, err := os.Open("test/nessus-agents.html")
		if err != nil {
			fmt.Println("Error opening file!!!")
		}
		return file, nil
	}
	doc, err := GetDownloadsDocument(fileReader)
	if err != nil {
		panic(err)
	}
	downloadItems := GetDownloadsFromDocument(doc)
	download := GetLatestMatchingDownload(downloadItems, "-amzn.x86_64.rpm")

	got := download
	want := newestDownloadItem

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetDownloadsDocument_AlreadyClosed(t *testing.T) {
	fileReader := func() (io.ReadCloser, error) {
		file, err := os.Open("test/nessus-agents.html")
		if err != nil {
			fmt.Println("Error opening file!!!")
		}
		err = file.Close()
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	_, err := GetDownloadsDocument(fileReader)

	got := err.Error()
	want := "read test/nessus-agents.html: file already closed"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestParseDownloadsPage(t *testing.T) {
	file, err := os.Open("test/nessus-agents.html")
	if err != nil {
		fmt.Println("Error opening file!!!")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
	}
	downloadItems := GetDownloadsFromDocument(doc)

	got := len(*downloadItems)
	want := 39

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestDownloadItemsByVersion_Reversed(t *testing.T) {
	newestDownloadItem := DownloadItem{
		Id:          16733,
		Name:        "NessusAgent-10.1.4-amzn.x86_64.rpm",
		Description: "Amazon Linux 2015.03, 2015.09, 2017.09, 2018.03, Amazon Linux 2 (x86_64)",
		MetaData: MetaData{
			Md5:         "b25c5e1a033eed78f7f2366254748bfb",
			Sha256:      "47dabe3c313c4026ad318cac25f88f31c4f55ec613a980e614f9f1223f355b40",
			Product:     "Nessus Agents - 10.1.4",
			Version:     "10.1.4",
			ReleaseDate: "2022-06-11T00:00:00.000Z",
		},
	}
	oldestDownloadItem := DownloadItem{
		Id:          16733,
		Name:        "NessusAgent-10.1.4-amzn.x86_64.rpm",
		Description: "Amazon Linux 2015.03, 2015.09, 2017.09, 2018.03, Amazon Linux 2 (x86_64)",
		MetaData: MetaData{
			Md5:         "b25c5e1a033eed78f7f2366254748bfb",
			Sha256:      "47dabe3c313c4026ad318cac25f88f31c4f55ec613a980e614f9f1223f355b40",
			Product:     "Nessus Agents - 10.1.4",
			Version:     "10.1.4",
			ReleaseDate: "2022-06-11T00:00:00.000Z",
		},
	}
	downloadItems := []DownloadItem{
		newestDownloadItem,
		oldestDownloadItem,
	}
	sort.Sort(DownloadItemsByVersion(downloadItems))
	got := downloadItems[len(downloadItems)-1]
	want := newestDownloadItem

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestDownloadItemsByVersion_AlreadySorted(t *testing.T) {
	newestDownloadItem := DownloadItem{
		Id:          16733,
		Name:        "NessusAgent-10.1.4-amzn.x86_64.rpm",
		Description: "Amazon Linux 2015.03, 2015.09, 2017.09, 2018.03, Amazon Linux 2 (x86_64)",
		MetaData: MetaData{
			Md5:         "b25c5e1a033eed78f7f2366254748bfb",
			Sha256:      "47dabe3c313c4026ad318cac25f88f31c4f55ec613a980e614f9f1223f355b40",
			Product:     "Nessus Agents - 10.1.4",
			Version:     "10.1.4",
			ReleaseDate: "2022-06-11T00:00:00.000Z",
		},
	}
	oldestDownloadItem := DownloadItem{
		Id:          16733,
		Name:        "NessusAgent-10.1.4-amzn.x86_64.rpm",
		Description: "Amazon Linux 2015.03, 2015.09, 2017.09, 2018.03, Amazon Linux 2 (x86_64)",
		MetaData: MetaData{
			Md5:         "b25c5e1a033eed78f7f2366254748bfb",
			Sha256:      "47dabe3c313c4026ad318cac25f88f31c4f55ec613a980e614f9f1223f355b40",
			Product:     "Nessus Agents - 10.1.4",
			Version:     "10.1.4",
			ReleaseDate: "2022-06-11T00:00:00.000Z",
		},
	}
	downloadItems := []DownloadItem{
		oldestDownloadItem,
		newestDownloadItem,
	}
	sort.Sort(DownloadItemsByVersion(downloadItems))
	got := downloadItems[len(downloadItems)-1]
	want := newestDownloadItem

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
