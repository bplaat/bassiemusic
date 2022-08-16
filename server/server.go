package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mileusna/useragent"
	"github.com/satori/go.uuid"
)

type AuthLoginResponse struct {
	Success bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
	User    *User   `json:"user,omitempty"`
}

func authLogin(res http.ResponseWriter, req *http.Request) {
	// Parse vars
	queryVars := req.URL.Query()

	email := ""
	if emailVar, ok := queryVars["email"]; ok {
		email = emailVar[0]
	}

	password := ""
	if passwordVar, ok := queryVars["password"]; ok {
		password = passwordVar[0]
	}

	// Get user by email
	userQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `password` FROM `users` WHERE `deleted_at` IS NULL AND `email` = ?", email)
	if err != nil {
		log.Fatalln(err)
	}
	defer userQuery.Close()

	var user User
	if !userQuery.Next() {
		var response AuthLoginResponse
		response.Success = false
		response.Message = "Wrong email or password"
		jsonResponse(res, response)
		return
	}
	userQuery.Scan(&user.ID, &user.Password)

	// Verify user password
	if !VerifyPassword(password, user.Password) {
		var response AuthLoginResponse
		response.Success = false
		response.Message = "Wrong email or password"
		jsonResponse(res, response)
		return
	}

	// Create new session
	ua := useragent.Parse(req.Header.Get("User-Agent"))

	token, err := HashPassword(email)
	if err != nil {
		log.Fatalln(err)
		return
	}

	sessionId := uuid.NewV4()
	_, err = db.Exec("INSERT INTO `sessions` (`id`, `user_id`, `token`, `os`, `platform`, `version`, `expires_at`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?, ?)",
		sessionId.String(), user.ID, token, ua.OS, ua.Name, ua.Version, time.Now().Add(365*24*60*60*time.Second).Format(time.RFC3339))
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Response
	var response AuthLoginResponse
	response.Success = true
	response.Token = base64.StdEncoding.EncodeToString([]byte(token))

	// User
	userQuery, err = db.Query("SELECT BIN_TO_UUID(`id`), `firstname`, `lastname`, `email`, `created_at`, `updated_at` FROM `users` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", user.ID)
	defer userQuery.Close()
	userQuery.Next()
	userQuery.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	response.User = &user
	jsonResponse(res, response)
}

func authLogout(res http.ResponseWriter, req *http.Request) {
	// TODO
}

func usersIndex(res http.ResponseWriter, req *http.Request) {
	query, page, limit := parseIndexVars(req)

	// Users
	usersQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `firstname`, `lastname`, `email`, `created_at`, `updated_at` FROM `users` WHERE `deleted_at` IS NULL AND (`firstname` LIKE ? OR `lastname` LIKE ? OR `email` LIKE ?) ORDER BY LOWER(`lastname`) LIMIT ?, ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer usersQuery.Close()

	users := []User{}
	for usersQuery.Next() {
		var user User
		usersQuery.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		users = append(users, user)
	}
	jsonResponse(res, users)
}

func usersShow(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// User
	userQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `firstname`, `lastname`, `email`, `created_at`, `updated_at` FROM `users` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", vars["id"])
	if err != nil {
		notFound(res, req)
		return
	}
	defer userQuery.Close()

	var user User
	if !userQuery.Next() {
		notFound(res, req)
		return
	}
	userQuery.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	jsonResponse(res, user)
}

func artistsIndex(res http.ResponseWriter, req *http.Request) {
	query, page, limit := parseIndexVars(req)

	// Artists
	artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()

	artists := []Artist{}
	for artistsQuery.Next() {
		var artist Artist
		artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
		artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", req.Host, artist.ID)
		artistAlbums(&artist, req)
		artists = append(artists, artist)
	}
	jsonResponse(res, artists)
}

func artistsShow(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// Artist
	artistQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", vars["id"])
	if err != nil {
		notFound(res, req)
		return
	}
	defer artistQuery.Close()

	var artist Artist
	if !artistQuery.Next() {
		notFound(res, req)
		return
	}
	artistQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
	artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", req.Host, artist.ID)
	artistAlbums(&artist, req)
	jsonResponse(res, artist)
}

func albumsIndex(res http.ResponseWriter, req *http.Request) {
	query, page, limit := parseIndexVars(req)

	// Albums
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `title` LIKE ? ORDER BY LOWER(`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()

	albums := []Album{}
	for albumsQuery.Next() {
		var album Album
		var albumType AlbumType
		albumsQuery.Scan(&album.ID, &albumType, &album.Title, &album.ReleasedAt, &album.Explicit, &album.CreatedAt, &album.UpdatedAt)
		if albumType == AlbumTypeAlbum {
			album.Type = "album"
		}
		if albumType == AlbumTypeEP {
			album.Type = "ep"
		}
		if albumType == AlbumTypeSingle {
			album.Type = "single"
		}
		album.Cover = fmt.Sprintf("http://%s/storage/albums/%s.jpg", req.Host, album.ID)
		albumGenres(&album, req)
		albumArtists(&album, req)
		albums = append(albums, album)
	}
	jsonResponse(res, albums)
}

func albumsShow(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// Album
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", vars["id"])
	if err != nil {
		notFound(res, req)
		return
	}
	defer albumsQuery.Close()

	var album Album
	var albumType AlbumType
	if !albumsQuery.Next() {
		notFound(res, req)
		return
	}
	albumsQuery.Scan(&album.ID, &albumType, &album.Title, &album.ReleasedAt, &album.Explicit, &album.CreatedAt, &album.UpdatedAt)
	if albumType == AlbumTypeAlbum {
		album.Type = "album"
	}
	if albumType == AlbumTypeEP {
		album.Type = "ep"
	}
	if albumType == AlbumTypeSingle {
		album.Type = "single"
	}
	album.Cover = fmt.Sprintf("http://%s/storage/albums/%s.jpg", req.Host, album.ID)
	albumGenres(&album, req)
	albumArtists(&album, req)
	albumTracks(&album, req)
	jsonResponse(res, album)
}

func genresIndex(res http.ResponseWriter, req *http.Request) {
	query, page, limit := parseIndexVars(req)

	// Genres
	genresQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `genres` WHERE `deleted_at` IS NULL AND `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer genresQuery.Close()

	genres := []Genre{}
	for genresQuery.Next() {
		var genre Genre
		genresQuery.Scan(&genre.ID, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
		genre.Image = fmt.Sprintf("http://%s/storage/genres/%s.jpg", req.Host, genre.ID)
		genreAlbums(&genre, req)
		genres = append(genres, genre)
	}
	jsonResponse(res, genres)
}

func genresShow(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// Genre
	genreQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `genres` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", vars["id"])
	if err != nil {
		notFound(res, req)
		return
	}
	defer genreQuery.Close()

	var genre Genre
	if !genreQuery.Next() {
		notFound(res, req)
		return
	}
	genreQuery.Scan(&genre.ID, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
	genre.Image = fmt.Sprintf("http://%s/storage/genres/%s.jpg", req.Host, genre.ID)
	genreAlbums(&genre, req)
	jsonResponse(res, genre)
}

func tracksIndex(res http.ResponseWriter, req *http.Request) {
	query, page, limit := parseIndexVars(req)

	trackssQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `plays`, `created_at`, `updated_at` FROM `tracks` WHERE `deleted_at` IS NULL AND `title` LIKE ? ORDER BY `plays` DESC, LOWER(`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer trackssQuery.Close()

	tracks := []Track{}
	for trackssQuery.Next() {
		var track Track
		var albumID string
		trackssQuery.Scan(&track.ID, &albumID, &track.Title, &track.Disk, &track.Position, &track.Duration, &track.Explicit, &track.Plays, &track.CreatedAt, &track.UpdatedAt)
		track.Music = fmt.Sprintf("http://%s/storage/tracks/%s.m4a", req.Host, track.ID)
		trackAlbum(&track, albumID, req)
		trackArtists(&track, req)
		tracks = append(tracks, track)
	}
	jsonResponse(res, tracks)
}

func tracksShow(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// Track
	trackssQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `plays`, `created_at`, `updated_at` FROM `tracks` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", vars["id"])
	if err != nil {
		notFound(res, req)
		return
	}
	defer trackssQuery.Close()

	var track Track
	var albumID string
	if !trackssQuery.Next() {
		notFound(res, req)
		return
	}
	trackssQuery.Scan(&track.ID, &albumID, &track.Title, &track.Disk, &track.Position, &track.Duration, &track.Explicit, &track.Plays, &track.CreatedAt, &track.UpdatedAt)
	track.Music = fmt.Sprintf("http://%s/storage/tracks/%s.m4a", req.Host, track.ID)
	trackAlbum(&track, albumID, req)
	trackArtists(&track, req)
	jsonResponse(res, track)
}

type TracksPlayResponse struct {
	Success bool `json:"success"`
}

func tracksPlay(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// Track
	trackssQuery, err := db.Query("SELECT `plays` FROM `tracks` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", vars["id"])
	if err != nil {
		notFound(res, req)
		return
	}
	defer trackssQuery.Close()

	var plays int64
	if !trackssQuery.Next() {
		notFound(res, req)
		return
	}
	trackssQuery.Scan(&plays)

	db.Exec("UPDATE `tracks` SET `plays` = ? WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", plays+1, vars["id"])
	var response TracksPlayResponse
	response.Success = true
	jsonResponse(res, response)
}

func notFound(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprint(res, "404 page not found\n")
}

func startServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "BassieMusic API")
	})
	router.NotFoundHandler = http.HandlerFunc(notFound)

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/login", authLogin)
	api.HandleFunc("/auth/logout", authLogout)

	// api.HandleFunc("/auth/sessions", authSessionIndex)
	// api.HandleFunc("/auth/sessions/{id}", authSessionShow)
	// api.HandleFunc("/auth/sessions/{id}/revoke", authSessionRevoke)

	api.HandleFunc("/users", usersIndex)
	api.HandleFunc("/users/{id}", usersShow)
	// POST api.HandleFunc("/users/{id}", usersEdit)

	api.HandleFunc("/artists", artistsIndex)
	api.HandleFunc("/artists/{id}", artistsShow)

	api.HandleFunc("/albums", albumsIndex)
	api.HandleFunc("/albums/{id}", albumsShow)

	api.HandleFunc("/genres", genresIndex)
	api.HandleFunc("/genres/{id}", genresShow)

	api.HandleFunc("/tracks", tracksIndex)
	api.HandleFunc("/tracks/{id}", tracksShow)
	api.HandleFunc("/tracks/{id}/play", tracksPlay)

	fileServer := http.FileServer(NeuteredFileSystem{http.Dir("./storage")})
	router.PathPrefix("/storage/").Handler(http.StripPrefix("/storage", fileServer))

	fmt.Printf("The server is listening on: http://localhost:8080/\n")
	http.ListenAndServe(":8080", router)
}
