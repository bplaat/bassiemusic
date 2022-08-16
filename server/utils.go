package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func fetch(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

func fetchJson(url string, data any) {
	json.Unmarshal(fetch(url), data)
}

func fetchFile(url string, path string) {
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
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
}

type NeuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func parseIndexVars(req *http.Request) (string, int, int) {
	queryVars := req.URL.Query()

	query := ""
	if queryVar, ok := queryVars["query"]; ok {
		query = queryVar[0]
	}

	page := 1
	if pageVar, ok := queryVars["page"]; ok {
		if pageInt, err := strconv.Atoi(pageVar[0]); err == nil {
			page = pageInt
			if page < 1 {
				page = 1
			}
		}
	}

	limit := 20
	if limitVar, ok := queryVars["limit"]; ok {
		if limitInt, err := strconv.Atoi(limitVar[0]); err == nil {
			limit = limitInt
			if limit < 1 {
				limit = 1
			}
			if limit > 50 {
				limit = 50
			}
		}
	}

	return query, page, limit
}

func jsonResponse(res http.ResponseWriter, data any) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	dataJson, _ := json.Marshal(data)
	res.Write(dataJson)
}
