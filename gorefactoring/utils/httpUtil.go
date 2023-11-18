package modules

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func DoGet(url string, response any) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, response)
	return nil
}

func GetHtml(url string) (string, error) {
	html := ""
	res, err := http.Get(url)
	if err != nil {
		return html, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return html, err
	}

	html = string(bytes)
	return html, nil
}

func GetDoc(url string) (*goquery.Document, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
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
	defer res.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
	  return nil, err
	}
	return doc, nil
}