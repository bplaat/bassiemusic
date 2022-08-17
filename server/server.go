package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mileusna/useragent"
	"github.com/satori/go.uuid"
)

func authLogin(c *fiber.Ctx) error {
	email := c.Query("email")
	password := c.Query("password")

	// Get user by email
	userQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `password` FROM `users` WHERE `deleted_at` IS NULL AND `email` = ?", email)
	if err != nil {
		log.Fatalln(err)
	}
	defer userQuery.Close()

	var user User
	if !userQuery.Next() {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Wrong email or password",
		})
	}
	userQuery.Scan(&user.ID, &user.Password)

	// Verify user password
	if !VerifyPassword(password, user.Password) {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Wrong email or password",
		})
	}

	// Create new session
	ua := useragent.Parse(c.Get("User-Agent"))

	token, err := HashPassword(email)
	if err != nil {
		log.Fatalln(err)
	}

	sessionId := uuid.NewV4()
	_, err = db.Exec("INSERT INTO `sessions` (`id`, `user_id`, `token`, `os`, `platform`, `version`, `expires_at`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?, ?)",
		sessionId.String(), user.ID, token, ua.OS, ua.Name, ua.Version, time.Now().Add(365*24*60*60*time.Second).Format(time.RFC3339))
	if err != nil {
		log.Fatalln(err)
	}

	// Write response
	userQuery, err = db.Query("SELECT BIN_TO_UUID(`id`), `firstname`, `lastname`, `email`, `created_at`, `updated_at` FROM `users` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", user.ID)
	defer userQuery.Close()
	userQuery.Next()
	userQuery.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	return c.JSON(fiber.Map{
		"success": true,
		"token":   base64.StdEncoding.EncodeToString([]byte(token)),
		"user":    user,
	})
}

func authLogout(c *fiber.Ctx) error {
	// TODO
	return fiber.ErrNotFound
}

func usersIndex(c *fiber.Ctx) error {
	query, page, limit := parseIndexVars(c)
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
	return c.JSON(users)
}

func usersShow(c *fiber.Ctx) error {
	userQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `firstname`, `lastname`, `email`, `created_at`, `updated_at` FROM `users` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", c.Params("id"))
	if err != nil {
		log.Fatalln(err)
	}
	defer userQuery.Close()

	var user User
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}
	userQuery.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	return c.JSON(user)
}

func artistsIndex(c *fiber.Ctx) error {
	query, page, limit := parseIndexVars(c)
	artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()

	artists := []Artist{}
	for artistsQuery.Next() {
		var artist Artist
		artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
		artist.Image = fmt.Sprintf("%s/storage/artists/%s.jpg", c.BaseURL(), artist.ID)
		artistAlbums(&artist, c)
		artists = append(artists, artist)
	}
	return c.JSON(artists)
}

func artistsShow(c *fiber.Ctx) error {
	artistQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", c.Params("id"))
	if err != nil {
		log.Fatalln(err)
	}
	defer artistQuery.Close()

	var artist Artist
	if !artistQuery.Next() {
		return fiber.ErrNotFound
	}
	artistQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
	artist.Image = fmt.Sprintf("%s/storage/artists/%s.jpg", c.BaseURL(), artist.ID)
	artistAlbums(&artist, c)
	return c.JSON(artist)
}

func albumsIndex(c *fiber.Ctx) error {
	query, page, limit := parseIndexVars(c)
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
		album.Cover = fmt.Sprintf("%s/storage/albums/%s.jpg", c.BaseURL(), album.ID)
		albumGenres(&album, c)
		albumArtists(&album, c)
		albums = append(albums, album)
	}
	return c.JSON(albums)
}

func albumsShow(c *fiber.Ctx) error {
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", c.Params("id"))
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()

	var album Album
	var albumType AlbumType
	if !albumsQuery.Next() {
		return fiber.ErrNotFound
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
	album.Cover = fmt.Sprintf("%s/storage/albums/%s.jpg", c.BaseURL(), album.ID)
	albumGenres(&album, c)
	albumArtists(&album, c)
	albumTracks(&album, c)
	return c.JSON(album)
}

func genresIndex(c *fiber.Ctx) error {
	query, page, limit := parseIndexVars(c)
	genresQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `genres` WHERE `deleted_at` IS NULL AND `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer genresQuery.Close()

	genres := []Genre{}
	for genresQuery.Next() {
		var genre Genre
		genresQuery.Scan(&genre.ID, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
		genre.Image = fmt.Sprintf("%s/storage/genres/%s.jpg", c.BaseURL(), genre.ID)
		genreAlbums(&genre, c)
		genres = append(genres, genre)
	}
	return c.JSON(genres)
}

func genresShow(c *fiber.Ctx) error {
	genreQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `genres` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", c.Params("id"))
	if err != nil {
		log.Fatalln(err)
	}
	defer genreQuery.Close()

	var genre Genre
	if !genreQuery.Next() {
		return fiber.ErrNotFound
	}
	genreQuery.Scan(&genre.ID, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
	genre.Image = fmt.Sprintf("%s/storage/genres/%s.jpg", c.BaseURL(), genre.ID)
	genreAlbums(&genre, c)
	return c.JSON(genre)
}

func tracksIndex(c *fiber.Ctx) error {
	query, page, limit := parseIndexVars(c)
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
		track.Music = fmt.Sprintf("%s/storage/tracks/%s.m4a", c.BaseURL(), track.ID)
		trackAlbum(&track, albumID, c)
		trackArtists(&track, c)
		tracks = append(tracks, track)
	}
	return c.JSON(tracks)
}

func tracksShow(c *fiber.Ctx) error {
	trackssQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `plays`, `created_at`, `updated_at` FROM `tracks` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", c.Params("id"))
	if err != nil {
		log.Fatalln(err)
	}
	defer trackssQuery.Close()

	var track Track
	var albumID string
	if !trackssQuery.Next() {
		return fiber.ErrNotFound
	}
	trackssQuery.Scan(&track.ID, &albumID, &track.Title, &track.Disk, &track.Position, &track.Duration, &track.Explicit, &track.Plays, &track.CreatedAt, &track.UpdatedAt)
	track.Music = fmt.Sprintf("%s/storage/tracks/%s.m4a", c.BaseURL(), track.ID)
	trackAlbum(&track, albumID, c)
	trackArtists(&track, c)
	return c.JSON(track)
}

func tracksPlay(c *fiber.Ctx) error {
	trackssQuery, err := db.Query("SELECT `plays` FROM `tracks` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", c.Params("id"))
	if err != nil {
		log.Fatalln(err)
	}
	defer trackssQuery.Close()

	var plays int64
	if !trackssQuery.Next() {
		return fiber.ErrNotFound
	}
	trackssQuery.Scan(&plays)
	db.Exec("UPDATE `tracks` SET `plays` = ? WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", plays+1, c.Params("id"))
	return c.JSON(fiber.Map{
		"success": true,
	})
}

func startServer() {
	app := fiber.New()
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("BassieMusic API")
	})

	storage := app.Group("/storage")
	storage.Use(etag.New())
	storage.Static("/", "./storage")

	api := app.Group("/api")
	api.Use(limiter.New(limiter.Config{
		Expiration: 60 * time.Second,
		Max:        100,
	}))

	api.Get("/auth/login", authLogin)
	api.Get("/auth/logout", authLogout)

	// api.Get("/auth/sessions", authSessionIndex)
	// api.Get("/auth/sessions/:id", authSessionShow)
	// api.Get("/auth/sessions/:id/revoke", authSessionRevoke)

	api.Get("/users", usersIndex)
	api.Get("/users/:id", usersShow)
	// api.Post("/users/:id", usersEdit)

	api.Get("/artists", artistsIndex)
	api.Get("/artists/:id", artistsShow)

	api.Get("/albums", albumsIndex)
	api.Get("/albums/:id", albumsShow)

	api.Get("/genres", genresIndex)
	api.Get("/genres/:id", genresShow)

	api.Get("/tracks", tracksIndex)
	api.Get("/tracks/:id", tracksShow)
	api.Get("/tracks/:id/play", tracksPlay)

	log.Fatal(app.Listen(":8080"))
}
