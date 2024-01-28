package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func doGet(url string) (*http.Response, error) {
	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
		Timeout: 5 * time.Second,
	}
	req , err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0")
	req.Header.Add("Accept", "application/json, text/javascript")
	req.Header.Add("Content-Type", "application/json")

	res , err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DoGet(url string, response any) error {
	res, err := doGet(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	  }
	return nil
}

func GetStatusCode(url string) (int, error) {
	res , err := doGet(url)
	if err != nil {
		return 500, err
	}
	return res.StatusCode, nil
}

func GetDoc(url string) (*goquery.Document, error) {
	res , err := doGet(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
	  return nil, err
	}
	return doc, nil
}