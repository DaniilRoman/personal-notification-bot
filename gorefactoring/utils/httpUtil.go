package modules

import (
	"encoding/json"
	"io"
	"net/http"
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