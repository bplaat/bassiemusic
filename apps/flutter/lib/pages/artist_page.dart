import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:http/http.dart' as http;
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../models/artist.dart';
import '../components/tracks_list.dart';
import '../components/album_card.dart';
import '../config.dart';
import '../utils.dart';

class ArtistPage extends StatefulWidget {
  final Artist incompleteArtist;

  const ArtistPage({super.key, required this.incompleteArtist});

  @override
  State<ArtistPage> createState() => _ArtistPageState();
}

class _ArtistPageState extends State<ArtistPage> {
  bool _artistLoaded = false;
  Artist? _artist;

  Future<bool> fetchCompleteArtist() async {
    if (_artistLoaded) return true;
    _artistLoaded = true;
    try {
      final prefs = await SharedPreferences.getInstance();
      final response = await http.get(
          Uri.parse('$apiUrl/artists/${widget.incompleteArtist.id}'),
          headers: {
            'User-Agent': userAgent(),
            'Authorization': 'Bearer ${prefs.getString('token')}'
          });
      if (response.statusCode == 200) {
        _artist = Artist.fromJson(json.decode(utf8.decode(response.bodyBytes)));
        return true;
      }
    } catch (exception) {
      print('Error: ${exception.toString()}');
    }
    return false;
  }

  @override
  Widget build(BuildContext context) {
    final lang = AppLocalizations.of(context)!;
    return Scaffold(
        appBar: AppBar(
          title: Text(widget.incompleteArtist.name),
        ),
        body: SingleChildScrollView(
            child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Artist cover
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
                                    widget.incompleteArtist.mediumImageUrl!),
                                fit: BoxFit.fill,
                              ))),

                      // Artist info
                      Container(
                          margin: const EdgeInsets.symmetric(vertical: 8),
                          child: Text(widget.incompleteArtist.name,
                              style: const TextStyle(
                                  fontSize: 24, fontWeight: FontWeight.w500))),

                      // Artist stuff
                      FutureBuilder(
                        future: fetchCompleteArtist(),
                        builder: (context, snapshot) {
                          if (snapshot.hasError) {
                            return Text("Error: ${snapshot.error}");
                          } else if (snapshot.hasData) {
                            return snapshot.data! && _artist != null
                                ? Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                        // Artist top tracks
                                        Container(
                                            margin: const EdgeInsets.symmetric(
                                                vertical: 8),
                                            child: Text(lang.artist_top_tracks,
                                                style: const TextStyle(
                                                    fontSize: 20,
                                                    fontWeight:
                                                        FontWeight.w500))),

                                        Container(
                                            margin: const EdgeInsets.symmetric(
                                                vertical: 8),
                                            child: TracksList(
                                                scrollController: null,
                                                tracks: _artist!.topTracks!,
                                                onTrackLikedChange:
                                                    (index, liked) {
                                                  setState(() => _artist!
                                                      .topTracks![index]
                                                      .liked = liked);
                                                })),

                                        // Artist albums
                                        Container(
                                            margin: const EdgeInsets.symmetric(
                                                vertical: 8),
                                            child: Text(lang.artist_albums,
                                                style: const TextStyle(
                                                    fontSize: 20,
                                                    fontWeight:
                                                        FontWeight.w500))),

                                        Container(
                                            margin: const EdgeInsets.symmetric(
                                                vertical: 8),
                                            child: GridView.builder(
                                                gridDelegate:
                                                    const SliverGridDelegateWithMaxCrossAxisExtent(
                                                  maxCrossAxisExtent: 240,
                                                  childAspectRatio: 1 / 1.37,
                                                  mainAxisSpacing: 8,
                                                  crossAxisSpacing: 8,
                                                ),
                                                shrinkWrap: true,
                                                itemCount:
                                                    _artist!.albums!.length,
                                                itemBuilder: (context, index) =>
                                                    AlbumCard(
                                                        album: _artist!
                                                            .albums![index]))),
                                      ])
                                : Column(children: const []);
                          }
                          return const Center(
                              child: CircularProgressIndicator());
                        },
                      )
                    ]))));
  }
}
