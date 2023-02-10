package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func Fetch(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

func FetchJson(url string, data any) {
	if err := json.Unmarshal(Fetch(url), data); err != nil {
		log.Fatalln(err)
	}
}

func FetchFile(url string, path string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()
	if _, err = io.Copy(out, resp.Body); err != nil {
		log.Fatalln(err)
	}
}
