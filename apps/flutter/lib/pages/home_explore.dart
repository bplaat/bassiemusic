import 'dart:convert';
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:cached_network_image/cached_network_image.dart';
import '../models/album.dart';
import '../config.dart';

class HomeExplorerTab extends StatefulWidget {
  const HomeExplorerTab({super.key});

  @override
  State<HomeExplorerTab> createState() => _HomeExplorerTabState();
}

class _HomeExplorerTabState extends State<HomeExplorerTab>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  final ScrollController _scrollController = ScrollController();
  List<Album> _albums = [];
  int _page = 1;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(vsync: this, length: 4, initialIndex: 2);

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
    _tabController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void loadPage() async {
    _isLoading = true;
    try {
      final response = await http.get(
          Uri.parse(
              'https://bassiemusic-api.plaatsoft.nl/albums?page=${_page++}'),
          headers: {
            HttpHeaders.authorizationHeader:
                'Bearer ${token}'
          });

      if (response.statusCode == 200) {
        final albumsJson = json.decode(response.body)['data'];
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
                      ))));
  }
}

class AlbumCard extends StatelessWidget {
  final Album album;

  const AlbumCard({Key? key, required this.album}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
      ),
      elevation: 5,
      clipBehavior: Clip.antiAliasWithSaveLayer,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Image(
            image: CachedNetworkImageProvider(album.mediumCoverUrl),
            fit: BoxFit.fill,
          ),
          Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Container(
                      margin: const EdgeInsets.only(bottom: 8),
                      child: Text(album.title,
                          overflow: TextOverflow.ellipsis,
                          style: const TextStyle(fontWeight: FontWeight.w500))),
                  Text(
                      album.artists!
                          .map(
                            (artist) => artist.name,
                          )
                          .join(', '),
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(color: Colors.grey)),
                ],
              ))
        ],
      ),
    );
  }
}
