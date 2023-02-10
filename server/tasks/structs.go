package tasks

type DeezerArtistSearch struct {
	Data []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Link          string `json:"link"`
		Picture       string `json:"picture"`
		PictureSmall  string `json:"picture_small"`
		PictureMedium string `json:"picture_medium"`
		PictureBig    string `json:"picture_big"`
		PictureXl     string `json:"picture_xl"`
		NbAlbum       int    `json:"nb_album"`
		NbFan         int    `json:"nb_fan"`
		Radio         bool   `json:"radio"`
		Tracklist     string `json:"tracklist"`
		Type          string `json:"type"`
	} `json:"data"`
	Total int    `json:"total"`
	Next  string `json:"next"`
}

type DeezerArtist struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Link          string `json:"link"`
	Share         string `json:"share"`
	Picture       string `json:"picture"`
	PictureSmall  string `json:"picture_small"`
	PictureMedium string `json:"picture_medium"`
	PictureBig    string `json:"picture_big"`
	PictureXl     string `json:"picture_xl"`
	NbAlbum       int    `json:"nb_album"`
	NbFan         int    `json:"nb_fan"`
	Radio         bool   `json:"radio"`
	Tracklist     string `json:"tracklist"`
	Type          string `json:"type"`
}

type DeezerArtistAlbums struct {
	Data []struct {
		ID             int    `json:"id"`
		Title          string `json:"title"`
		Link           string `json:"link"`
		Cover          string `json:"cover"`
		CoverSmall     string `json:"cover_small"`
		CoverMedium    string `json:"cover_medium"`
		CoverBig       string `json:"cover_big"`
		CoverXl        string `json:"cover_xl"`
		Md5Image       string `json:"md5_image"`
		GenreID        int    `json:"genre_id"`
		Fans           int    `json:"fans"`
		ReleaseDate    string `json:"release_date"`
		RecordType     string `json:"record_type"`
		Tracklist      string `json:"tracklist"`
		ExplicitLyrics bool   `json:"explicit_lyrics"`
		Type           string `json:"type"`
	} `json:"data"`
	Total int    `json:"total"`
	Next  string `json:"next"`
}

type DeezerAlbumSearch struct {
	Data []struct {
		ID             int    `json:"id"`
		Title          string `json:"title"`
		Link           string `json:"link"`
		Cover          string `json:"cover"`
		CoverSmall     string `json:"cover_small"`
		CoverMedium    string `json:"cover_medium"`
		CoverBig       string `json:"cover_big"`
		CoverXl        string `json:"cover_xl"`
		Md5Image       string `json:"md5_image"`
		GenreID        int    `json:"genre_id"`
		NbTracks       int    `json:"nb_tracks"`
		RecordType     string `json:"record_type"`
		Tracklist      string `json:"tracklist"`
		ExplicitLyrics bool   `json:"explicit_lyrics"`
		Artist         struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Link          string `json:"link"`
			Picture       string `json:"picture"`
			PictureSmall  string `json:"picture_small"`
			PictureMedium string `json:"picture_medium"`
			PictureBig    string `json:"picture_big"`
			PictureXl     string `json:"picture_xl"`
			Tracklist     string `json:"tracklist"`
			Type          string `json:"type"`
		} `json:"artist"`
		Type string `json:"type"`
	} `json:"data"`
	Total int    `json:"total"`
	Next  string `json:"next"`
}

type DeezerAlbum struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Upc         string `json:"upc"`
	Link        string `json:"link"`
	Share       string `json:"share"`
	Cover       string `json:"cover"`
	CoverSmall  string `json:"cover_small"`
	CoverMedium string `json:"cover_medium"`
	CoverBig    string `json:"cover_big"`
	CoverXl     string `json:"cover_xl"`
	Md5Image    string `json:"md5_image"`
	GenreID     int    `json:"genre_id"`
	Genres      struct {
		Data []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
			Type    string `json:"type"`
		} `json:"data"`
	} `json:"genres"`
	Label                 string `json:"label"`
	NbTracks              int    `json:"nb_tracks"`
	Duration              int    `json:"duration"`
	Fans                  int    `json:"fans"`
	ReleaseDate           string `json:"release_date"`
	RecordType            string `json:"record_type"`
	Available             bool   `json:"available"`
	Tracklist             string `json:"tracklist"`
	ExplicitLyrics        bool   `json:"explicit_lyrics"`
	ExplicitContentLyrics int    `json:"explicit_content_lyrics"`
	ExplicitContentCover  int    `json:"explicit_content_cover"`
	Contributors          []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Link          string `json:"link"`
		Share         string `json:"share"`
		Picture       string `json:"picture"`
		PictureSmall  string `json:"picture_small"`
		PictureMedium string `json:"picture_medium"`
		PictureBig    string `json:"picture_big"`
		PictureXl     string `json:"picture_xl"`
		Radio         bool   `json:"radio"`
		Tracklist     string `json:"tracklist"`
		Type          string `json:"type"`
		Role          string `json:"role"`
	} `json:"contributors"`
	Artist struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		PictureSmall  string `json:"picture_small"`
		PictureMedium string `json:"picture_medium"`
		PictureBig    string `json:"picture_big"`
		PictureXl     string `json:"picture_xl"`
		Tracklist     string `json:"tracklist"`
		Type          string `json:"type"`
	} `json:"artist"`
	Type   string `json:"type"`
	Tracks struct {
		Data []struct {
			ID                    int    `json:"id"`
			Readable              bool   `json:"readable"`
			Title                 string `json:"title"`
			TitleShort            string `json:"title_short"`
			TitleVersion          string `json:"title_version"`
			Link                  string `json:"link"`
			Duration              int    `json:"duration"`
			Rank                  int    `json:"rank"`
			ExplicitLyrics        bool   `json:"explicit_lyrics"`
			ExplicitContentLyrics int    `json:"explicit_content_lyrics"`
			ExplicitContentCover  int    `json:"explicit_content_cover"`
			Preview               string `json:"preview"`
			Md5Image              string `json:"md5_image"`
			Artist                struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				Tracklist string `json:"tracklist"`
				Type      string `json:"type"`
			} `json:"artist"`
			Album struct {
				ID          int    `json:"id"`
				Title       string `json:"title"`
				Cover       string `json:"cover"`
				CoverSmall  string `json:"cover_small"`
				CoverMedium string `json:"cover_medium"`
				CoverBig    string `json:"cover_big"`
				CoverXl     string `json:"cover_xl"`
				Md5Image    string `json:"md5_image"`
				Tracklist   string `json:"tracklist"`
				Type        string `json:"type"`
			} `json:"album"`
			Type string `json:"type"`
		} `json:"data"`
	} `json:"tracks"`
}

type DeezerGenre struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	PictureSmall  string `json:"picture_small"`
	PictureMedium string `json:"picture_medium"`
	PictureBig    string `json:"picture_big"`
	PictureXl     string `json:"picture_xl"`
	Type          string `json:"type"`
}

type DeezerTrack struct {
	ID                    int      `json:"id"`
	Readable              bool     `json:"readable"`
	Title                 string   `json:"title"`
	TitleShort            string   `json:"title_short"`
	TitleVersion          string   `json:"title_version"`
	Isrc                  string   `json:"isrc"`
	Link                  string   `json:"link"`
	Share                 string   `json:"share"`
	Duration              int      `json:"duration"`
	TrackPosition         int      `json:"track_position"`
	DiskNumber            int      `json:"disk_number"`
	Rank                  int      `json:"rank"`
	ReleaseDate           string   `json:"release_date"`
	ExplicitLyrics        bool     `json:"explicit_lyrics"`
	ExplicitContentLyrics int      `json:"explicit_content_lyrics"`
	ExplicitContentCover  int      `json:"explicit_content_cover"`
	Preview               string   `json:"preview"`
	Bpm                   float64  `json:"bpm"`
	Gain                  float64  `json:"gain"`
	AvailableCountries    []string `json:"available_countries"`
	Contributors          []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Link          string `json:"link"`
		Share         string `json:"share"`
		Picture       string `json:"picture"`
		PictureSmall  string `json:"picture_small"`
		PictureMedium string `json:"picture_medium"`
		PictureBig    string `json:"picture_big"`
		PictureXl     string `json:"picture_xl"`
		Radio         bool   `json:"radio"`
		Tracklist     string `json:"tracklist"`
		Type          string `json:"type"`
		Role          string `json:"role"`
	} `json:"contributors"`
	Md5Image string `json:"md5_image"`
	Artist   struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Link          string `json:"link"`
		Share         string `json:"share"`
		Picture       string `json:"picture"`
		PictureSmall  string `json:"picture_small"`
		PictureMedium string `json:"picture_medium"`
		PictureBig    string `json:"picture_big"`
		PictureXl     string `json:"picture_xl"`
		Radio         bool   `json:"radio"`
		Tracklist     string `json:"tracklist"`
		Type          string `json:"type"`
	} `json:"artist"`
	Album struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Link        string `json:"link"`
		Cover       string `json:"cover"`
		CoverSmall  string `json:"cover_small"`
		CoverMedium string `json:"cover_medium"`
		CoverBig    string `json:"cover_big"`
		CoverXl     string `json:"cover_xl"`
		Md5Image    string `json:"md5_image"`
		ReleaseDate string `json:"release_date"`
		Tracklist   string `json:"tracklist"`
		Type        string `json:"type"`
	} `json:"album"`
	Type string `json:"type"`
}

type YoutubeVideo struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Formats []struct {
		FormatID   string  `json:"format_id"`
		FormatNote string  `json:"format_note"`
		Ext        string  `json:"ext"`
		Protocol   string  `json:"protocol"`
		Acodec     string  `json:"acodec"`
		Vcodec     string  `json:"vcodec"`
		URL        string  `json:"url"`
		Width      int     `json:"width"`
		Height     int     `json:"height"`
		Fps        float64 `json:"fps"`
		Rows       int     `json:"rows,omitempty"`
		Columns    int     `json:"columns,omitempty"`
		Fragments  []struct {
			URL      string  `json:"url"`
			Duration float64 `json:"duration"`
		} `json:"fragments,omitempty"`
		AudioExt    string  `json:"audio_ext"`
		VideoExt    string  `json:"video_ext"`
		Format      string  `json:"format"`
		Resolution  string  `json:"resolution"`
		AspectRatio float64 `json:"aspect_ratio"`
		HTTPHeaders struct {
			UserAgent      string `json:"User-Agent"`
			Accept         string `json:"Accept"`
			AcceptLanguage string `json:"Accept-Language"`
			SecFetchMode   string `json:"Sec-Fetch-Mode"`
		} `json:"http_headers"`
		Asr                int         `json:"asr,omitempty"`
		Filesize           int         `json:"filesize,omitempty"`
		SourcePreference   int         `json:"source_preference,omitempty"`
		AudioChannels      int         `json:"audio_channels,omitempty"`
		Quality            float64     `json:"quality,omitempty"`
		HasDrm             bool        `json:"has_drm,omitempty"`
		Tbr                float64     `json:"tbr,omitempty"`
		Language           interface{} `json:"language,omitempty"`
		LanguagePreference int         `json:"language_preference,omitempty"`
		Preference         interface{} `json:"preference,omitempty"`
		DynamicRange       interface{} `json:"dynamic_range,omitempty"`
		Abr                float64     `json:"abr,omitempty"`
		DownloaderOptions  struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options,omitempty"`
		Container      string  `json:"container,omitempty"`
		Vbr            float64 `json:"vbr,omitempty"`
		FilesizeApprox int     `json:"filesize_approx,omitempty"`
	} `json:"formats"`
	Thumbnails []struct {
		URL        string `json:"url"`
		Preference int    `json:"preference"`
		ID         string `json:"id"`
		Height     int    `json:"height,omitempty"`
		Width      int    `json:"width,omitempty"`
		Resolution string `json:"resolution,omitempty"`
	} `json:"thumbnails"`
	Thumbnail         string      `json:"thumbnail"`
	Description       string      `json:"description"`
	Uploader          string      `json:"uploader"`
	UploaderID        string      `json:"uploader_id"`
	UploaderURL       string      `json:"uploader_url"`
	ChannelID         string      `json:"channel_id"`
	ChannelURL        string      `json:"channel_url"`
	Duration          int         `json:"duration"`
	ViewCount         int         `json:"view_count"`
	AverageRating     interface{} `json:"average_rating"`
	AgeLimit          int         `json:"age_limit"`
	WebpageURL        string      `json:"webpage_url"`
	Categories        []string    `json:"categories"`
	Tags              []string    `json:"tags"`
	PlayableInEmbed   bool        `json:"playable_in_embed"`
	LiveStatus        string      `json:"live_status"`
	ReleaseTimestamp  interface{} `json:"release_timestamp"`
	FormatSortFields  []string    `json:"_format_sort_fields"`
	AutomaticCaptions struct {
	} `json:"automatic_captions"`
	Subtitles struct {
	} `json:"subtitles"`
	CommentCount         interface{} `json:"comment_count"`
	Chapters             interface{} `json:"chapters"`
	LikeCount            int         `json:"like_count"`
	Channel              string      `json:"channel"`
	ChannelFollowerCount int         `json:"channel_follower_count"`
	UploadDate           string      `json:"upload_date"`
	Availability         string      `json:"availability"`
	OriginalURL          string      `json:"original_url"`
	WebpageURLBasename   string      `json:"webpage_url_basename"`
	WebpageURLDomain     string      `json:"webpage_url_domain"`
	Joinctor             string      `json:"joinctor"`
	JoinctorKey          string      `json:"joinctor_key"`
	PlaylistCount        int         `json:"playlist_count"`
	Playlist             string      `json:"playlist"`
	PlaylistID           string      `json:"playlist_id"`
	PlaylistTitle        string      `json:"playlist_title"`
	PlaylistUploader     interface{} `json:"playlist_uploader"`
	PlaylistUploaderID   interface{} `json:"playlist_uploader_id"`
	NEntries             int         `json:"n_entries"`
	PlaylistIndex        int         `json:"playlist_index"`
	LastPlaylistIndex    int         `json:"__last_playlist_index"`
	PlaylistAutonumber   int         `json:"playlist_autonumber"`
	DisplayID            string      `json:"display_id"`
	Fulltitle            string      `json:"fulltitle"`
	DurationString       string      `json:"duration_string"`
	IsLive               bool        `json:"is_live"`
	WasLive              bool        `json:"was_live"`
	RequestedSubtitles   interface{} `json:"requested_subtitles"`
	HasDrm               interface{} `json:"_has_drm"`
	RequestedFormats     []struct {
		Asr                interface{} `json:"asr"`
		Filesize           int         `json:"filesize"`
		FormatID           string      `json:"format_id"`
		FormatNote         string      `json:"format_note"`
		SourcePreference   int         `json:"source_preference"`
		Fps                int         `json:"fps"`
		AudioChannels      interface{} `json:"audio_channels"`
		Height             int         `json:"height"`
		Quality            float64     `json:"quality"`
		HasDrm             bool        `json:"has_drm"`
		Tbr                float64     `json:"tbr"`
		URL                string      `json:"url"`
		Width              int         `json:"width"`
		Language           interface{} `json:"language"`
		LanguagePreference int         `json:"language_preference"`
		Preference         interface{} `json:"preference"`
		Ext                string      `json:"ext"`
		Vcodec             string      `json:"vcodec"`
		Acodec             string      `json:"acodec"`
		DynamicRange       string      `json:"dynamic_range"`
		Vbr                float64     `json:"vbr,omitempty"`
		DownloaderOptions  struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options"`
		Container   string  `json:"container"`
		Protocol    string  `json:"protocol"`
		VideoExt    string  `json:"video_ext"`
		AudioExt    string  `json:"audio_ext"`
		Format      string  `json:"format"`
		Resolution  string  `json:"resolution"`
		AspectRatio float64 `json:"aspect_ratio"`
		HTTPHeaders struct {
			UserAgent      string `json:"User-Agent"`
			Accept         string `json:"Accept"`
			AcceptLanguage string `json:"Accept-Language"`
			SecFetchMode   string `json:"Sec-Fetch-Mode"`
		} `json:"http_headers"`
		Abr float64 `json:"abr,omitempty"`
	} `json:"requested_formats"`
	Format         string      `json:"format"`
	FormatID       string      `json:"format_id"`
	Ext            string      `json:"ext"`
	Protocol       string      `json:"protocol"`
	Language       interface{} `json:"language"`
	FormatNote     string      `json:"format_note"`
	FilesizeApprox int         `json:"filesize_approx"`
	Tbr            float64     `json:"tbr"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	Resolution     string      `json:"resolution"`
	Fps            int         `json:"fps"`
	DynamicRange   string      `json:"dynamic_range"`
	Vcodec         string      `json:"vcodec"`
	Vbr            float64     `json:"vbr"`
	StretchedRatio interface{} `json:"stretched_ratio"`
	AspectRatio    float64     `json:"aspect_ratio"`
	Acodec         string      `json:"acodec"`
	Abr            float64     `json:"abr"`
	Asr            int         `json:"asr"`
	AudioChannels  int         `json:"audio_channels"`
	Epoch          int         `json:"epoch"`
	Filename       string      `json:"_filename"`
	Filename0      string      `json:"filename"`
	Urls           string      `json:"urls"`
	Type           string      `json:"_type"`
	Version        struct {
		Version        string      `json:"version"`
		CurrentGitHead interface{} `json:"current_git_head"`
		ReleaseGitHead string      `json:"release_git_head"`
		Repository     string      `json:"repository"`
	} `json:"_version"`
}
