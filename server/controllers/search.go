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
	authUser := c.Locals("authUser").(*models.User)
	query, _, _ := utils.ParseIndexVars(c)

	// Get tracks
	tracks := models.TrackModel.WithArgs("liked", c.Locals("authUser")).With("artists", "album").WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByRaw("`plays` DESC, LOWER(`title`)").Limit(10).Get()

	// Get albums
	albums := models.AlbumModel.With("artists").WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`title`)").Limit(10).Get()

	// Get artists
	artists := models.ArtistModel.WhereRaw("`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`name`)").Limit(10).Get()

	// Get Genres
	genres := models.GenreModel.WhereRaw("`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`name`)").Limit(10).Get()

	// Get Playlists
	playlistsQuery := models.PlaylistModel.WhereRaw("`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`name`)")
	if authUser.Role != models.UserRoleAdmin {
		playlistsQuery = playlistsQuery.Where("public", true)
	}
	playlists := playlistsQuery.Limit(10).Get()

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
	query, _, _ := utils.ParseIndexVars(c)

	// Search deezer artists
	var deezerArtistSearch structs.DeezerArtistSearch
	if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/search/artist?q=%s", url.QueryEscape(query)), &deezerArtistSearch); err != nil {
		return fiber.ErrBadGateway
	}

	// Search deezer albums and filter out what already exists
	var deezerAlbumSearch structs.DeezerAlbumSearch
	if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/search/album?q=%s", url.QueryEscape(query)), &deezerAlbumSearch); err != nil {
		return fiber.ErrBadGateway
	}
	deezerAlbums := []structs.DeezerAlbumSearchItem{}
	for _, deezerAlbum := range deezerAlbumSearch.Data {
		if models.AlbumModel.Where("title", deezerAlbum.Title).First() == nil {
			deezerAlbums = append(deezerAlbums, deezerAlbum)
		}
	}

	// Return response
	return c.JSON(fiber.Map{
		"artists": deezerArtistSearch.Data,
		"albums":  deezerAlbums,
	})
}
