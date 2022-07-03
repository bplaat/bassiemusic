package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Artist struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Albums    []Album   `json:"albums,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Album struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	ReleasedAt time.Time `json:"released_at"`
	Cover      string    `json:"cover"`
	Artists    []Artist  `json:"artists,omitempty"`
	Tracks     []Track   `json:"tracks,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Track struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Disk      int       `json:"disk"`
	Position  int       `json:"position"`
	Duration  int       `json:"duration"`
	Music     string    `json:"music"`
	Artists   []Artist  `json:"artists,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func parseIndexVars(request *http.Request) (string, int, int) {
	queryVars := request.URL.Query()

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

func artistsIndex(response http.ResponseWriter, request *http.Request) {
	query, page, limit := parseIndexVars(request)
	artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()

	artists := []Artist{}
	for artistsQuery.Next() {
		var artist Artist
		artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
		artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", request.Host, artist.ID)

		albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `title`, `released_at`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `album_id` FROM `album_artist` WHERE `artist_id` = UUID_TO_BIN(?)) ORDER BY `released_at`", artist.ID)
		if err != nil {
			log.Fatalln(err)
		}
		defer albumsQuery.Close()
		for albumsQuery.Next() {
			var album Album
			albumsQuery.Scan(&album.ID, &album.Title, &album.ReleasedAt, &album.CreatedAt, &album.UpdatedAt)
			album.Cover = fmt.Sprintf("http://%s/storage/albums/%s.jpg", request.Host, album.ID)
			artist.Albums = append(artist.Albums, album)
		}

		artists = append(artists, artist)
	}

	response.Header().Set("Content-Type", "application/json")
	artistsJson, _ := json.Marshal(artists)
	response.Write(artistsJson)
}

func artistsShow(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	artistQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", vars["id"])
	if err != nil {
		notFound(response, request)
		return
	}
	defer artistQuery.Close()

	var artist Artist
	artistQuery.Next()
	if err := artistQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt); err != nil {
		notFound(response, request)
		return
	}
	artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", request.Host, artist.ID)

	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `title`, `released_at`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `album_id` FROM `album_artist` WHERE `artist_id` = UUID_TO_BIN(?)) ORDER BY `released_at`", artist.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()
	for albumsQuery.Next() {
		var album Album
		albumsQuery.Scan(&album.ID, &album.Title, &album.ReleasedAt, &album.CreatedAt, &album.UpdatedAt)
		album.Cover = fmt.Sprintf("http://%s/storage/albums/%s.jpg", request.Host, album.ID)
		artist.Albums = append(artist.Albums, album)
	}

	response.Header().Set("Content-Type", "application/json")
	artistJson, _ := json.Marshal(artist)
	response.Write(artistJson)
}

func albumsIndex(response http.ResponseWriter, request *http.Request) {
	query, page, limit := parseIndexVars(request)
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `title`, `released_at`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `title` LIKE ? ORDER BY LOWER(`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()

	albums := []Album{}
	for albumsQuery.Next() {
		var album Album
		albumsQuery.Scan(&album.ID, &album.Title, &album.ReleasedAt, &album.CreatedAt, &album.UpdatedAt)
		album.Cover = fmt.Sprintf("http://%s/storage/albums/%s.jpg", request.Host, album.ID)

		artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `artist_id` FROM `album_artist` WHERE `album_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", album.ID)
		if err != nil {
			log.Fatalln(err)
		}
		defer artistsQuery.Close()
		for artistsQuery.Next() {
			var artist Artist
			artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
			artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", request.Host, artist.ID)
			album.Artists = append(album.Artists, artist)
		}

		albums = append(albums, album)
	}

	response.Header().Set("Content-Type", "application/json")
	albumsJson, _ := json.Marshal(albums)
	response.Write(albumsJson)
}

func albumsShow(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `title`, `released_at`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", vars["id"])
	if err != nil {
		notFound(response, request)
		return
	}
	defer albumsQuery.Close()

	var album Album
	albumsQuery.Next();
	if err := albumsQuery.Scan(&album.ID, &album.Title, &album.ReleasedAt, &album.CreatedAt, &album.UpdatedAt); err != nil {
		notFound(response, request)
		return
	}
	album.Cover = fmt.Sprintf("http://%s/storage/albums/%s.jpg", request.Host, album.ID)

	artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `artist_id` FROM `album_artist` WHERE `album_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()
	for artistsQuery.Next() {
		var artist Artist
		artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
		artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", request.Host, artist.ID)
		album.Artists = append(album.Artists, artist)
	}

	tracksQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `title`, `disk`, `position`, `duration`, `created_at`, `updated_at` FROM `tracks` WHERE `deleted_at` IS NULL AND `album_id` = UUID_TO_BIN(?) ORDER BY `disk`, `position`", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer tracksQuery.Close()
	for tracksQuery.Next() {
		var track Track
		tracksQuery.Scan(&track.ID, &track.Title, &track.Disk, &track.Position, &track.Duration, &track.CreatedAt, &track.UpdatedAt)
		track.Music = fmt.Sprintf("http://%s/storage/tracks/%s.m4a", request.Host, track.ID)

		artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `artist_id` FROM `track_artist` WHERE `track_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", track.ID)
		if err != nil {
			log.Fatalln(err)
		}
		defer artistsQuery.Close()
		for artistsQuery.Next() {
			var artist Artist
			artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
			artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", request.Host, artist.ID)
			track.Artists = append(track.Artists, artist)
		}

		album.Tracks = append(album.Tracks, track)
	}

	response.Header().Set("Content-Type", "application/json")
	albumJson, _ := json.Marshal(album)
	response.Write(albumJson)
}

func tracksIndex(response http.ResponseWriter, request *http.Request) {
	query, page, limit := parseIndexVars(request)
	trackssQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `title`, `disk`, `position`, `duration`, `created_at`, `updated_at` FROM `tracks` WHERE `deleted_at` IS NULL AND `title` LIKE ? ORDER BY LOWER(`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer trackssQuery.Close()

	tracks := []Track{}
	for trackssQuery.Next() {
		var track Track
		trackssQuery.Scan(&track.ID, &track.Title, &track.Disk, &track.Position, &track.Duration, &track.CreatedAt, &track.UpdatedAt)
		track.Music = fmt.Sprintf("http://%s/storage/tracks/%s.m4a", request.Host, track.ID)

		artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `artist_id` FROM `track_artist` WHERE `track_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", track.ID)
		if err != nil {
			log.Fatalln(err)
		}
		defer artistsQuery.Close()
		for artistsQuery.Next() {
			var artist Artist
			artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
			artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", request.Host, artist.ID)
			track.Artists = append(track.Artists, artist)
		}

		tracks = append(tracks, track)
	}

	response.Header().Set("Content-Type", "application/json")
	tracksJson, _ := json.Marshal(tracks)
	response.Write(tracksJson)
}

func tracksShow(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	trackssQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `title`, `disk`, `position`, `duration`, `created_at`, `updated_at` FROM `tracks` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", vars["id"])
	if err != nil {
		notFound(response, request)
		return
	}
	defer trackssQuery.Close()

	var track Track
	trackssQuery.Next()
	if err := trackssQuery.Scan(&track.ID, &track.Title, &track.Disk, &track.Position, &track.Duration, &track.CreatedAt, &track.UpdatedAt); err != nil {
		notFound(response, request)
		return
	}
	track.Music = fmt.Sprintf("http://%s/storage/tracks/%s.m4a", request.Host, track.ID)

	artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `artist_id` FROM `track_artist` WHERE `track_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", track.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()
	for artistsQuery.Next() {
		var artist Artist
		artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
		artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", request.Host, artist.ID)
		track.Artists = append(track.Artists, artist)
	}

	response.Header().Set("Content-Type", "application/json")
	trackJson, _ := json.Marshal(track)
	response.Write(trackJson)
}

func notFound(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusNotFound)
	fmt.Fprint(response, "404 Not Found")
}

func startServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "BassieMusic API")
	})
	router.NotFoundHandler = http.HandlerFunc(notFound)

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/artists", artistsIndex)
	api.HandleFunc("/artists/{id}", artistsShow)
	api.HandleFunc("/albums", albumsIndex)
	api.HandleFunc("/albums/{id}", albumsShow)
	api.HandleFunc("/tracks", tracksIndex)
	api.HandleFunc("/tracks/{id}", tracksShow)

	fileServer := http.FileServer(NeuteredFileSystem{http.Dir("./storage")})
	router.PathPrefix("/storage/").Handler(http.StripPrefix("/storage", fileServer))

	fmt.Printf("The server is listening on: http://localhost:8080/\n")
	http.ListenAndServe(":8080", router)
}
