package modules

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func doGet(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req , err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0"},
	}

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