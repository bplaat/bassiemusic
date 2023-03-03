import 'dart:convert';
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:http/http.dart' as http;
import '../components/artist_card.dart';
import '../components/genre_card.dart';
import '../components/album_card.dart';
import '../models/artist.dart';
import '../models/genre.dart';
import '../models/album.dart';
import '../models/track.dart';
import '../config.dart';

class HomeExplorerTab extends StatefulWidget {
  const HomeExplorerTab({super.key});

  @override
  State<HomeExplorerTab> createState() => _HomeExplorerTabState();
}

class _HomeExplorerTabState extends State<HomeExplorerTab>
    with
        SingleTickerProviderStateMixin,
        AutomaticKeepAliveClientMixin<HomeExplorerTab> {
  late TabController _tabController;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(vsync: this, length: 4, initialIndex: 0);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: Theme.of(context).colorScheme.background,
        appBar: AppBar(
            title: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(
              text: "Artists",
            ),
            Tab(
              text: "Genres",
            ),
            Tab(
              text: "Albums",
            ),
            Tab(
              text: "Tracks",
            ),
          ],
        )),
        body: TabBarView(controller: _tabController, children: [
          const ArtistsTab(),
          const GenresTab(),
          const AlbumsTab(),
          const TracksTab(),
        ]));
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
      final response = await http.get(
          Uri.parse(
              'https://bassiemusic-api.plaatsoft.nl/artists?page=${_page++}'),
          headers: {HttpHeaders.authorizationHeader: 'Bearer ${token}'});

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
                  childAspectRatio: 1 / 1.5,
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
                : const Center(
                    child: Text('No artists'),
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
      final response = await http.get(
          Uri.parse(
              'https://bassiemusic-api.plaatsoft.nl/genres?page=${_page++}'),
          headers: {HttpHeaders.authorizationHeader: 'Bearer ${token}'});

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
                  childAspectRatio: 1 / 1.5,
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
                : const Center(
                    child: Text('No genres'),
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
      final response = await http.get(
          Uri.parse(
              'https://bassiemusic-api.plaatsoft.nl/albums?page=${_page++}'),
          headers: {HttpHeaders.authorizationHeader: 'Bearer ${token}'});

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
                  childAspectRatio: 1 / 1.5,
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
                : const Center(
                    child: Text('No albums'),
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
      final response = await http.get(
          Uri.parse(
              'https://bassiemusic-api.plaatsoft.nl/tracks?page=${_page++}'),
          headers: {HttpHeaders.authorizationHeader: 'Bearer ${token}'});

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
    return RefreshIndicator(
        onRefresh: () async {
          _tracks = [];
          _page = 1;
          _isLoading = false;
          loadPage();
        },
        child: _tracks.isNotEmpty
            ? ListView.builder(
                padding: const EdgeInsets.all(8),
                controller: _scrollController,
                itemCount: _tracks.length,
                itemBuilder: (context, index) => InkWell(
                    onTap: () => {},
                    child: Row(children: [
                      Container(
                          margin: EdgeInsets.only(right: 16),
                          child: SizedBox(
                              width: 56,
                              height: 56,
                              child: Card(
                                clipBehavior: Clip.antiAliasWithSaveLayer,
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(6),
                                ),
                                elevation: 2,
                                child: Container(
                                    decoration: BoxDecoration(
                                        image: DecorationImage(
                                            fit: BoxFit.cover,
                                            image: CachedNetworkImageProvider(
                                                _tracks[index]
                                                    .album!
                                                    .smallCoverUrl)))),
                              ))),
                      Expanded(
                          flex: 1,
                          child: Column(children: [
                            Container(
                                margin: EdgeInsets.only(bottom: 4),
                                child: SizedBox(
                                    width: double.infinity,
                                    child: Text(_tracks[index].title,
                                        style: TextStyle(
                                            fontSize: 16,
                                            fontWeight: FontWeight.w500)))),
                            Container(
                                margin: EdgeInsets.only(bottom: 4),
                                child: SizedBox(
                                    width: double.infinity,
                                    child: Text(
                                        _tracks[index]
                                            .artists!
                                            .map(
                                              (artist) => artist.name,
                                            )
                                            .join(', '),
                                        style: TextStyle(
                                            fontSize: 16,
                                            color: Colors.grey)))),
                          ])),
                      IconButton(
                          onPressed: () => {},
                          color: _tracks[index].liked! ? Colors.red : null,
                          icon: Icon(_tracks[index].liked! ? Icons.favorite : Icons.favorite_outline)),
                    ])))
            : (_isLoading
                ? const Center(child: CircularProgressIndicator())
                : const Center(
                    child: Text('No tracks'),
                  )));
  }
}
