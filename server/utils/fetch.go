package utils

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/structs"
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

func DeezerFetch(url string, data any) error {
	tries := 0
	for {
		body, err := Fetch(url)

		if err == nil {
			var deezerError structs.DeezerError
			err = json.Unmarshal(body, &deezerError)
			if deezerError.Error.Code != 0 {
				err = errors.New("Api Rate limit")
			}
		}

		if err == nil {
			err = json.Unmarshal(body, data)
		}

		if err != nil {
			tries += 1
			time.Sleep(2 * time.Second)
			log.Println("waiting")
			if tries == 5 {
				return err
			}
		} else {
			return err
		}
	}
}

func DeezerFetchFile(url string, path string) {
	tries := 0
	if url != "" {
		for {
			response, err := http.Get(url)
			if err == nil {
				defer response.Body.Close()
				out, err := os.Create(path)
				if err == nil {
					defer out.Close()
					if _, err = io.Copy(out, response.Body); err == nil {
						return
					}
				}
			}

			tries += 1
			time.Sleep(2 * time.Second)
			if tries == 5 {
				log.Fatalln(err)
			}
		}
	}
}
