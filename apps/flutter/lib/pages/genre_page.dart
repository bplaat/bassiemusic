import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../components/album_card.dart';
import '../models/album.dart';
import '../models/genre.dart';
import '../config.dart';
import '../utils.dart';

class GenrePage extends StatefulWidget {
  final Genre genre;

  const GenrePage({super.key, required this.genre});

  @override
  State<GenrePage> createState() => _GenrePageState();
}

class _GenrePageState extends State<GenrePage> {
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
          .get(Uri.parse('$apiUrl/genres/${widget.genre.id}/albums?page=${_page++}'), headers: {
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
    return Scaffold(
        appBar: AppBar(
          title: Text(widget.genre.name),
        ),
        body: RefreshIndicator(
            onRefresh: () async {
              _albums = [];
              _page = 1;
              _isLoading = false;
              loadPage();
            },
            child: _albums.isNotEmpty
                ? GridView.builder(
                    gridDelegate:
                        const SliverGridDelegateWithMaxCrossAxisExtent(
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
                        child: Text(lang.genre_albums_empty),
                      ))));
  }
}
