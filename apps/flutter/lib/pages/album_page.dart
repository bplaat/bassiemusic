import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:intl/intl.dart';
import 'package:http/http.dart' as http;
import '../models/album.dart';
import '../models/track.dart';
import '../components/tracks_list.dart';
import '../config.dart';
import '../utils.dart';

class AlbumPage extends StatefulWidget {
  final Album incompleteAlbum;

  const AlbumPage({super.key, required this.incompleteAlbum});

  @override
  State<AlbumPage> createState() => _AlbumPageState();
}

class _AlbumPageState extends State<AlbumPage> {
  bool _albumLoaded = false;
  Album? _album;

  Future<bool> fetchCompleteAlbum() async {
    if (_albumLoaded) return true;
    _albumLoaded = true;
    try {
      final prefs = await SharedPreferences.getInstance();
      final response = await http.get(
          Uri.parse('$apiUrl/albums/${widget.incompleteAlbum.id}'),
          headers: {
            'User-Agent': userAgent(),
            'Authorization': 'Bearer ${prefs.getString('token')}'
          });
      if (response.statusCode == 200) {
        _album = Album.fromJson(json.decode(utf8.decode(response.bodyBytes)));
        for (Track track in _album!.tracks!) {
          track.album = _album;
        }
        return true;
      }
    } catch (exception) {
      print('Error: ${exception.toString()}');
    }
    return false;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text(widget.incompleteAlbum.title),
        ),
        body: SingleChildScrollView(
            child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Album cover
                    Card(
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                        elevation: 5,
                        clipBehavior: Clip.antiAliasWithSaveLayer,
                        child: AspectRatio(
                            aspectRatio: 1,
                            child: Image(
                              image: CachedNetworkImageProvider(
                                  widget.incompleteAlbum.mediumCoverUrl!),
                              fit: BoxFit.fill,
                            ))),

                    // Album info
                    Container(
                        margin: const EdgeInsets.symmetric(vertical: 8),
                        child: Text(widget.incompleteAlbum.title,
                            style: const TextStyle(
                                fontSize: 24, fontWeight: FontWeight.w500))),
                    Container(
                        margin: const EdgeInsets.symmetric(vertical: 8),
                        child: Text(
                            DateFormat('yyyy-MM-dd')
                                .format(widget.incompleteAlbum.releasedAt),
                            style: const TextStyle(fontSize: 16))),
                    Container(
                        margin: const EdgeInsets.symmetric(vertical: 8),
                        child: Text(
                            widget.incompleteAlbum.genres!
                                .map(
                                  (genre) => genre.name,
                                )
                                .join(', '),
                            style: const TextStyle(fontSize: 16))),
                    Container(
                        margin: const EdgeInsets.symmetric(vertical: 8),
                        child: Text(
                            widget.incompleteAlbum.artists!
                                .map(
                                  (artist) => artist.name,
                                )
                                .join(', '),
                            style: const TextStyle(fontSize: 16))),

                    // Album tracks
                    Container(
                        margin: const EdgeInsets.symmetric(vertical: 8),
                        child: FutureBuilder(
                          future: fetchCompleteAlbum(),
                          builder: (context, snapshot) {
                            if (snapshot.hasError) {
                              return Text('Error: ${snapshot.error}');
                            } else if (snapshot.hasData) {
                              return snapshot.data! && _album != null
                                  ? TracksList(
                                      scrollController: null,
                                      tracks: _album!.tracks!,
                                      onTrackLikedChange: (index, liked) {
                                        setState(() => _album!
                                            .tracks![index].liked = liked);
                                      })
                                  : Column(children: const []);
                            }
                            return const Center(
                                child: CircularProgressIndicator());
                          },
                        ))
                  ],
                ))));
  }
}
