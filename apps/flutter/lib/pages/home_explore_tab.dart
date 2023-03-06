import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../components/artist_card.dart';
import '../components/genre_card.dart';
import '../components/album_card.dart';
import '../components/tracks_list.dart';
import '../models/artist.dart';
import '../models/genre.dart';
import '../models/album.dart';
import '../models/track.dart';
import '../config.dart';
import '../utils.dart';

class HomeExplorerTab extends StatefulWidget {
  const HomeExplorerTab({super.key});

  @override
  State<HomeExplorerTab> createState() => _HomeExplorerTabState();
}

class _HomeExplorerTabState extends State<HomeExplorerTab>
    with AutomaticKeepAliveClientMixin<HomeExplorerTab> {
  final _pageController = PageController(initialPage: 0);
  int _page = 0;

  @override
  bool get wantKeepAlive => true;

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final lang = AppLocalizations.of(context)!;
    return Column(children: [
      Padding(
          padding:
              const EdgeInsets.only(top: 16, left: 16, right: 16, bottom: 8),
          child: Row(children: [
            OutlinedButton(
                style:
                    OutlinedButton.styleFrom(padding: const EdgeInsets.all(16)),
                onPressed: () {
                  _pageController.animateToPage(0,
                      duration: const Duration(milliseconds: 200),
                      curve: Curves.ease);
                  setState(() => _page = 0);
                },
                child: Text(
                  lang.home_explore_artists,
                  style: const TextStyle(color: Colors.grey),
                )),
            const SizedBox(width: 8),
            OutlinedButton(
                style:
                    OutlinedButton.styleFrom(padding: const EdgeInsets.all(16)),
                onPressed: () {
                  _pageController.animateToPage(1,
                      duration: const Duration(milliseconds: 200),
                      curve: Curves.ease);
                  setState(() => _page = 1);
                },
                child: Text(
                  lang.home_explore_genres,
                  style: const TextStyle(color: Colors.grey),
                )),
            const SizedBox(width: 8),
            OutlinedButton(
                style:
                    OutlinedButton.styleFrom(padding: const EdgeInsets.all(16)),
                onPressed: () {
                  _pageController.animateToPage(2,
                      duration: const Duration(milliseconds: 200),
                      curve: Curves.ease);
                  setState(() => _page = 2);
                },
                child: Text(
                  lang.home_explore_albums,
                  style: const TextStyle(color: Colors.grey),
                )),
            const SizedBox(width: 8),
            OutlinedButton(
                style:
                    OutlinedButton.styleFrom(padding: const EdgeInsets.all(16)),
                onPressed: () {
                  _pageController.animateToPage(3,
                      duration: const Duration(milliseconds: 200),
                      curve: Curves.ease);
                  setState(() => _page = 3);
                },
                child: Text(
                  lang.home_explore_tracks,
                  style: const TextStyle(color: Colors.grey),
                )),
          ])),
      Expanded(
          child: PageView(
              controller: _pageController,
              onPageChanged: (index) {
                setState(() => _page = index);
              },
              children: const [
            ArtistsTab(),
            GenresTab(),
            AlbumsTab(),
            TracksTab(),
          ]))
    ]);
  }
}

// Artists
class ArtistsTab extends StatefulWidget {
  const ArtistsTab({super.key});

  @override
  State<ArtistsTab> createState() => _ArtistsTabState();
}

class _ArtistsTabState extends State<ArtistsTab> {
  final ScrollController _scrollController = ScrollController();
  List<Artist> _artists = [];
  int _page = 1;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    loadPage();
    _scrollController.addListener(() {
      if (!_isLoading &&
          _scrollController.position.pixels >
              _scrollController.position.maxScrollExtent * 0.9) {
        loadPage();
      }
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void loadPage() async {
    _isLoading = true;
    try {
      final prefs = await SharedPreferences.getInstance();
      final response = await http
          .get(Uri.parse('$apiUrl/artists?page=${_page++}'), headers: {
        'User-Agent': userAgent(),
        'Authorization': 'Bearer ${prefs.getString('token')}'
      });

      if (response.statusCode == 200) {
        final artistsJson =
            json.decode(utf8.decode(response.bodyBytes))['data'];
        List<Artist> newArtists =
            artistsJson.map<Artist>((json) => Artist.fromJson(json)).toList();
        _artists.addAll(newArtists);
        _isLoading = false;
        setState(() => _artists = _artists);
      }
    } catch (exception) {
      print('Error: ${exception.toString()}');
    }
  }

  @override
  Widget build(BuildContext context) {
    final lang = AppLocalizations.of(context)!;
    return RefreshIndicator(
        onRefresh: () async {
          _artists = [];
          _page = 1;
          _isLoading = false;
          loadPage();
        },
        child: _artists.isNotEmpty
            ? GridView.builder(
                gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
                  maxCrossAxisExtent: 240,
                  childAspectRatio: 1 / 1.28,
                  mainAxisSpacing: 8,
                  crossAxisSpacing: 8,
                ),
                padding: const EdgeInsets.all(8),
                controller: _scrollController,
                itemCount: _artists.length,
                itemBuilder: (context, index) =>
                    ArtistCard(artist: _artists[index]))
            : (_isLoading
                ? const Center(child: CircularProgressIndicator())
                : Center(
                    child: Text(lang.home_explore_artists_empty),
                  )));
  }
}

// Genres
class GenresTab extends StatefulWidget {
  const GenresTab({super.key});

  @override
  State<GenresTab> createState() => _GenresTabState();
}

class _GenresTabState extends State<GenresTab> {
  final ScrollController _scrollController = ScrollController();
  List<Genre> _genres = [];
  int _page = 1;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    loadPage();
    _scrollController.addListener(() {
      if (!_isLoading &&
          _scrollController.position.pixels >
              _scrollController.position.maxScrollExtent * 0.9) {
        loadPage();
      }
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void loadPage() async {
    _isLoading = true;
    try {
      final prefs = await SharedPreferences.getInstance();
      final response = await http
          .get(Uri.parse('$apiUrl/genres?page=${_page++}'), headers: {
        'User-Agent': userAgent(),
        'Authorization': 'Bearer ${prefs.getString('token')}'
      });

      if (response.statusCode == 200) {
        final genresJson = json.decode(utf8.decode(response.bodyBytes))['data'];
        List<Genre> newGenres =
            genresJson.map<Genre>((json) => Genre.fromJson(json)).toList();
        _genres.addAll(newGenres);
        _isLoading = false;
        setState(() => _genres = _genres);
      }
    } catch (exception) {
      print('Error: ${exception.toString()}');
    }
  }

  @override
  Widget build(BuildContext context) {
    final lang = AppLocalizations.of(context)!;
    return RefreshIndicator(
        onRefresh: () async {
          _genres = [];
          _page = 1;
          _isLoading = false;
          loadPage();
        },
        child: _genres.isNotEmpty
            ? GridView.builder(
                gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
                  maxCrossAxisExtent: 240,
                  childAspectRatio: 1 / 1.28,
                  mainAxisSpacing: 8,
                  crossAxisSpacing: 8,
                ),
                padding: const EdgeInsets.all(8),
                controller: _scrollController,
                itemCount: _genres.length,
                itemBuilder: (context, index) =>
                    GenreCard(genre: _genres[index]))
            : (_isLoading
                ? const Center(child: CircularProgressIndicator())
                : Center(
                    child: Text(lang.home_explore_genres_empty),
                  )));
  }
}

// Albums
class AlbumsTab extends StatefulWidget {
  const AlbumsTab({super.key});

  @override
  State<AlbumsTab> createState() => _AlbumsTabState();
}

class _AlbumsTabState extends State<AlbumsTab> {
  final ScrollController _scrollController = ScrollController();
  List<Album> _albums = [];
  int _page = 1;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    loadPage();
    _scrollController.addListener(() {
      if (!_isLoading &&
          _scrollController.position.pixels >
              _scrollController.position.maxScrollExtent * 0.9) {
        loadPage();
      }
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void loadPage() async {
    _isLoading = true;
    try {
      final prefs = await SharedPreferences.getInstance();
      final response = await http
          .get(Uri.parse('$apiUrl/albums?page=${_page++}'), headers: {
        'User-Agent': userAgent(),
        'Authorization': 'Bearer ${prefs.getString('token')}'
      });

      if (response.statusCode == 200) {
        final albumsJson = json.decode(utf8.decode(response.bodyBytes))['data'];
        List<Album> newAlbums =
            albumsJson.map<Album>((json) => Album.fromJson(json)).toList();
        _albums.addAll(newAlbums);
        _isLoading = false;
        setState(() => _albums = _albums);
      }
    } catch (exception) {
      print('Error: ${exception.toString()}');
    }
  }

  @override
  Widget build(BuildContext context) {
    final lang = AppLocalizations.of(context)!;
    return RefreshIndicator(
        onRefresh: () async {
          _albums = [];
          _page = 1;
          _isLoading = false;
          loadPage();
        },
        child: _albums.isNotEmpty
            ? GridView.builder(
                gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
                  maxCrossAxisExtent: 240,
                  childAspectRatio: 1 / 1.37,
                  mainAxisSpacing: 8,
                  crossAxisSpacing: 8,
                ),
                padding: const EdgeInsets.all(8),
                controller: _scrollController,
                itemCount: _albums.length,
                itemBuilder: (context, index) =>
                    AlbumCard(album: _albums[index]))
            : (_isLoading
                ? const Center(child: CircularProgressIndicator())
                : Center(
                    child: Text(lang.home_explore_albums_empty),
                  )));
  }
}

// Tracks
class TracksTab extends StatefulWidget {
  const TracksTab({super.key});

  @override
  State<TracksTab> createState() => _TracksTabState();
}

class _TracksTabState extends State<TracksTab> {
  final ScrollController _scrollController = ScrollController();
  List<Track> _tracks = [];
  int _page = 1;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    loadPage();
    _scrollController.addListener(() {
      if (!_isLoading &&
          _scrollController.position.pixels >
              _scrollController.position.maxScrollExtent * 0.9) {
        loadPage();
      }
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void loadPage() async {
    _isLoading = true;
    try {
      final prefs = await SharedPreferences.getInstance();
      final response = await http
          .get(Uri.parse('$apiUrl/tracks?page=${_page++}'), headers: {
        'User-Agent': userAgent(),
        'Authorization': 'Bearer ${prefs.getString('token')}'
      });

      if (response.statusCode == 200) {
        final tracksJson = json.decode(utf8.decode(response.bodyBytes))['data'];
        List<Track> newTracks =
            tracksJson.map<Track>((json) => Track.fromJson(json)).toList();
        _tracks.addAll(newTracks);
        _isLoading = false;
        setState(() => _tracks = _tracks);
      }
    } catch (exception) {
      print('Error: ${exception.toString()}');
    }
  }

  @override
  Widget build(BuildContext context) {
    final lang = AppLocalizations.of(context)!;
    return RefreshIndicator(
        onRefresh: () async {
          _tracks = [];
          _page = 1;
          _isLoading = false;
          loadPage();
        },
        child: _tracks.isNotEmpty
            ? TracksList(
                scrollController: _scrollController,
                tracks: _tracks,
                onTrackLikedChange: (index, liked) => setState(() {
                      _tracks[index].liked = liked;
                    }))
            : (_isLoading
                ? const Center(child: CircularProgressIndicator())
                : Center(
                    child: Text(lang.home_explore_tracks_empty),
                  )));
  }
}
