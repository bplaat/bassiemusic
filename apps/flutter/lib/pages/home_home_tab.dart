import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../components/artist_card.dart';
import '../components/album_card.dart';
import '../components/tracks_list.dart';
import '../models/user.dart';
import '../models/artist.dart';
import '../models/album.dart';
import '../models/track.dart';
import '../config.dart';
import '../utils.dart';

class HomeHomeTab extends StatefulWidget {
  final User user;

  const HomeHomeTab({super.key, required this.user});

  @override
  State<HomeHomeTab> createState() => _HomeHomeTabState();
}

class _HomeHomeTabState extends State<HomeHomeTab>
    with AutomaticKeepAliveClientMixin<HomeHomeTab> {
  @override
  bool get wantKeepAlive => true;

  Future<List<Track>?> fetchLastPlayedTracks() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final response = await http.get(
          Uri.parse('$apiUrl/users/${widget.user.id}/played_tracks'),
          headers: {
            'User-Agent': userAgent(),
            'Authorization': 'Bearer ${prefs.getString('token')}'
          });

      if (response.statusCode == 200) {
        final tracksJson = json.decode(utf8.decode(response.bodyBytes))['data'];
        return tracksJson.map<Track>((json) => Track.fromJson(json)).toList();
      }
    } catch (exception) {
      print('Error: ${exception.toString()}');
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
      future: fetchLastPlayedTracks(),
      builder: (context, snapshot) {
        if (snapshot.hasError) {
          return Center(child: Text('Error: ${snapshot.error}'));
        } else if (snapshot.hasData) {
          final lang = AppLocalizations.of(context)!;
          final lastPlayedtracks = snapshot.data!;

          Map<String, Artist> lastPlayedArtistsHelper = {};
          for (final track in lastPlayedtracks) {
            for (final artist in track.artists!) {
              lastPlayedArtistsHelper[artist.id] = artist;
            }
          }
          List<Artist> lastPlayedArtists =
              lastPlayedArtistsHelper.values.length > 6
                  ? lastPlayedArtistsHelper.values.toList().sublist(0, 6)
                  : lastPlayedArtistsHelper.values.toList();

          Map<String, Album> lastPlayedAlbumsHelper = {};
          for (final track in lastPlayedtracks) {
            lastPlayedAlbumsHelper[track.album!.id] = track.album!;
          }
          List<Album> lastPlayedAlbums =
              lastPlayedAlbumsHelper.values.length > 6
                  ? lastPlayedAlbumsHelper.values.toList().sublist(0, 6)
                  : lastPlayedAlbumsHelper.values.toList();

          return lastPlayedtracks.isNotEmpty
              ? SingleChildScrollView(
                  child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Container(
                              margin: const EdgeInsets.symmetric(vertical: 8),
                              child: Text(
                                  lang.home_home_hey(widget.user.username),
                                  style: const TextStyle(
                                      fontSize: 20,
                                      fontWeight: FontWeight.w500))),

                          // Last artists
                          Container(
                              margin: const EdgeInsets.symmetric(vertical: 8),
                              child: Text(lang.home_home_last_artists,
                                  style: const TextStyle(
                                      fontSize: 20,
                                      fontWeight: FontWeight.w500))),

                          Container(
                              margin: const EdgeInsets.symmetric(vertical: 8),
                              child: GridView.builder(
                                  gridDelegate:
                                      const SliverGridDelegateWithMaxCrossAxisExtent(
                                    maxCrossAxisExtent: 240,
                                    childAspectRatio: 1 / 1.28,
                                    mainAxisSpacing: 8,
                                    crossAxisSpacing: 8,
                                  ),
                                  shrinkWrap: true,
                                  physics: const NeverScrollableScrollPhysics(),
                                  itemCount: lastPlayedArtists.length,
                                  itemBuilder: (context, index) => ArtistCard(
                                      artist: lastPlayedArtists[index]))),

                          // Last albums
                          Container(
                              margin: const EdgeInsets.symmetric(vertical: 8),
                              child: Text(lang.home_home_last_albums,
                                  style: const TextStyle(
                                      fontSize: 20,
                                      fontWeight: FontWeight.w500))),

                          Container(
                              margin: const EdgeInsets.symmetric(vertical: 8),
                              child: GridView.builder(
                                  gridDelegate:
                                      const SliverGridDelegateWithMaxCrossAxisExtent(
                                    maxCrossAxisExtent: 240,
                                    childAspectRatio: 1 / 1.37,
                                    mainAxisSpacing: 8,
                                    crossAxisSpacing: 8,
                                  ),
                                  shrinkWrap: true,
                                  physics: const NeverScrollableScrollPhysics(),
                                  itemCount: lastPlayedAlbums.length,
                                  itemBuilder: (context, index) => AlbumCard(
                                      album: lastPlayedAlbums[index]))),

                          // Last tracks
                          Container(
                              margin: const EdgeInsets.symmetric(vertical: 8),
                              child: Text(lang.home_home_last_tracks,
                                  style: const TextStyle(
                                      fontSize: 20,
                                      fontWeight: FontWeight.w500))),

                          Container(
                              margin: const EdgeInsets.symmetric(vertical: 8),
                              child: TracksList(
                                  scrollController: null,
                                  tracks: lastPlayedtracks.sublist(0, 5),
                                  onTrackLikedChange: (index, liked) =>
                                      setState(() {
                                        lastPlayedtracks[index].liked = liked;
                                      })))
                        ],
                      )))
              : Center(
                  child: Text(lang.home_home_empty),
                );
        }
        return const Center(child: CircularProgressIndicator());
      },
    );
  }
}
