package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mileusna/useragent"
	"github.com/satori/go.uuid"
)

type User struct {
	ID        string    `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	ReleasedAt time.Time `json:"released_at"`
	Explicit   bool      `json:"explicit"`
	Cover      string    `json:"cover"`
	Genres     []Genre   `json:"genres,omitempty"`
	Artists    []Artist  `json:"artists,omitempty"`
	Tracks     []Track   `json:"tracks,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AlbumType int

const AlbumTypeAlbum AlbumType = 0
const AlbumTypeEP AlbumType = 1
const AlbumTypeSingle AlbumType = 2

type Genre struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Albums    []Album   `json:"albums,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Track struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Disk      int       `json:"disk"`
	Position  int       `json:"position"`
	Duration  int       `json:"duration"`
	Explicit  bool      `json:"explicit"`
	Plays     int64     `json:"plays"`
	Music     string    `json:"music"`
	Album     *Album    `json:"album,omitempty"`
	Artists   []Artist  `json:"artists,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func artistAlbums(artist *Artist, req *http.Request) {
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `album_id` FROM `album_artist` WHERE `artist_id` = UUID_TO_BIN(?)) ORDER BY `released_at` DESC", artist.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()
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
		albumArtists(&album, req)
		artist.Albums = append(artist.Albums, album)
	}
}

func albumGenres(album *Album, req *http.Request) {
	genresQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `genres` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `genre_id` FROM `album_genre` WHERE `album_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer genresQuery.Close()
	for genresQuery.Next() {
		var genre Genre
		genresQuery.Scan(&genre.ID, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
		genre.Image = fmt.Sprintf("http://%s/storage/genres/%s.jpg", req.Host, genre.ID)
		album.Genres = append(album.Genres, genre)
	}
}

func albumArtists(album *Album, req *http.Request) {
	artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `artist_id` FROM `album_artist` WHERE `album_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()
	for artistsQuery.Next() {
		var artist Artist
		artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
		artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", req.Host, artist.ID)
		album.Artists = append(album.Artists, artist)
	}
}

func albumTracks(album *Album, req *http.Request) {
	tracksQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `title`, `disk`, `position`, `duration`, `explicit`, `plays`, `created_at`, `updated_at` FROM `tracks` WHERE `deleted_at` IS NULL AND `album_id` = UUID_TO_BIN(?) ORDER BY `disk`, `position`", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer tracksQuery.Close()
	for tracksQuery.Next() {
		var track Track
		tracksQuery.Scan(&track.ID, &track.Title, &track.Disk, &track.Position, &track.Duration, &track.Explicit, &track.Plays, &track.CreatedAt, &track.UpdatedAt)
		track.Music = fmt.Sprintf("http://%s/storage/tracks/%s.m4a", req.Host, track.ID)
		trackArtists(&track, req)
		album.Tracks = append(album.Tracks, track)
	}
}

func genreAlbums(genre *Genre, req *http.Request) {
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `album_id` FROM `album_genre` WHERE `genre_id` = UUID_TO_BIN(?)) ORDER BY `released_at` DESC", genre.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()
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
		albumArtists(&album, req)
		genre.Albums = append(genre.Albums, album)
	}
}

func trackArtists(track *Track, req *http.Request) {
	artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `artist_id` FROM `track_artist` WHERE `track_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", track.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()
	for artistsQuery.Next() {
		var artist Artist
		artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
		artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", req.Host, artist.ID)
		track.Artists = append(track.Artists, artist)
	}
}

func trackAlbum(track *Track, albumID string, req *http.Request) {
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", albumID)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()

	var album Album
	var albumType AlbumType
	albumsQuery.Next()
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
	track.Album = &album
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
	usersQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `password` FROM `users` WHERE `deleted_at` IS NULL AND `email` = ?", email)
	if err != nil {
		log.Fatalln(err)
	}
	defer usersQuery.Close()

	var user User
	usersQuery.Next()
	if err := usersQuery.Scan(&user.ID, &user.Password); err != nil {
		notFound(res, req)
		return
	}

	// Verify user password
	if !VerifyPassword(password, user.Password) {
		fmt.Fprint(res, "Wrong password\n")
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

	// Write response
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

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	usersJson, _ := json.Marshal(users)
	res.Write(usersJson)
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
	userQuery.Next()
	if err := userQuery.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		notFound(res, req)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	userJson, _ := json.Marshal(user)
	res.Write(userJson)
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

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	artistsJson, _ := json.Marshal(artists)
	res.Write(artistsJson)
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
	artistQuery.Next()
	if err := artistQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt); err != nil {
		notFound(res, req)
		return
	}
	artist.Image = fmt.Sprintf("http://%s/storage/artists/%s.jpg", req.Host, artist.ID)
	artistAlbums(&artist, req)

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	artistJson, _ := json.Marshal(artist)
	res.Write(artistJson)
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

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	albumsJson, _ := json.Marshal(albums)
	res.Write(albumsJson)
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
	albumsQuery.Next()
	if err := albumsQuery.Scan(&album.ID, &albumType, &album.Title, &album.ReleasedAt, &album.Explicit, &album.CreatedAt, &album.UpdatedAt); err != nil {
		notFound(res, req)
		return
	}
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

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	albumJson, _ := json.Marshal(album)
	res.Write(albumJson)
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

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	genresJson, _ := json.Marshal(genres)
	res.Write(genresJson)
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
	genreQuery.Next()
	if err := genreQuery.Scan(&genre.ID, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt); err != nil {
		notFound(res, req)
		return
	}
	genre.Image = fmt.Sprintf("http://%s/storage/genres/%s.jpg", req.Host, genre.ID)
	genreAlbums(&genre, req)

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	genreJson, _ := json.Marshal(genre)
	res.Write(genreJson)
}

func tracksIndex(res http.ResponseWriter, req *http.Request) {
	query, page, limit := parseIndexVars(req)

	// Tracks
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

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	tracksJson, _ := json.Marshal(tracks)
	res.Write(tracksJson)
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
	trackssQuery.Next()
	if err := trackssQuery.Scan(&track.ID, &albumID, &track.Title, &track.Disk, &track.Position, &track.Duration, &track.Explicit, &track.Plays, &track.CreatedAt, &track.UpdatedAt); err != nil {
		notFound(res, req)
		return
	}
	track.Music = fmt.Sprintf("http://%s/storage/tracks/%s.m4a", req.Host, track.ID)
	trackAlbum(&track, albumID, req)
	trackArtists(&track, req)

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	trackJson, _ := json.Marshal(track)
	res.Write(trackJson)
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
	trackssQuery.Next()
	if err := trackssQuery.Scan(&plays); err != nil {
		notFound(res, req)
		return
	}

	db.Exec("UPDATE `tracks` SET `plays` = ? WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", plays+1, vars["id"])

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Write([]byte("{\"success\":true}"))
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
