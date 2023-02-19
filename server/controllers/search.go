package controllers

import (
	"fmt"
	"net/url"

	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/structs"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func SearchIndex(c *fiber.Ctx) error {
	query, _, _ := utils.ParseIndexVars(c)

	// Get tracks
	tracks := models.TrackModel(c).With("like", "artists", "album").WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByRaw("`plays` DESC, `updated_at` DESC").Limit("10").Get()

	// Get albums
	albums := models.AlbumModel(c).With("artists", "genres").WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`title`)").Limit("10").Get()

	// Get artists
	artists := models.ArtistModel(c).WhereRaw("`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`name`)").Limit("10").Get()

	// Get Genres
	genres := models.GenreModel(c).WhereRaw("`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`name`)").Limit("10").Get()

	// Get Playlists
	playlists := models.PlaylistModel(c).With("like", "user").Where("public", true).WhereRaw("`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`name`)").Limit("10").Get()

	// Return all values
	return c.JSON(fiber.Map{
		"success":   true,
		"tracks":    tracks,
		"albums":    albums,
		"artists":   artists,
		"genres":    genres,
		"playlists": playlists,
	})
}

func DeezerSearchIndex(c *fiber.Ctx) error {
	// Search deezer artists
	var deezerArtistSearch structs.DeezerArtistSearch
	if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/search/artist?q=%s", url.QueryEscape(c.Query("q"))), &deezerArtistSearch); err != nil {
		return fiber.ErrBadGateway
	}

	// Search deezer albums and filter out what already exists
	var deezerAlbumSearch structs.DeezerAlbumSearch
	if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/search/album?q=%s", url.QueryEscape(c.Query("q"))), &deezerAlbumSearch); err != nil {
		return fiber.ErrBadGateway
	}
	deezerAlbums := []structs.DeezerAlbumSearchItem{}
	for _, deezerAlbum := range deezerAlbumSearch.Data {
		if models.AlbumModel(c).Where("title", deezerAlbum.Title).First() == nil {
			deezerAlbums = append(deezerAlbums, deezerAlbum)
		}
	}

	// Return response
	return c.JSON(fiber.Map{
		"artists": deezerArtistSearch.Data,
		"albums":  deezerAlbums,
	})
}
