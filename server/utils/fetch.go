package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func Fetch(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func FetchJson(url string, data any) error {
	body, err := Fetch(url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, data); err != nil {
		return err
	}
	return nil
}

func FetchFile(url string, path string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()
	if _, err = io.Copy(out, response.Body); err != nil {
		log.Fatalln(err)
	}
}
