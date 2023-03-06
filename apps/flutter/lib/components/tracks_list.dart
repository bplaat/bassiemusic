import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;
import 'package:cached_network_image/cached_network_image.dart';
import '../models/track.dart';
import '../config.dart';
import '../utils.dart';

class TracksList extends StatelessWidget {
  static const playerChannel = MethodChannel('bassiemusic.plaatsoft.nl/player');

  final ScrollController? scrollController;
  final List<Track> tracks;
  final Function(int, bool) onTrackLikedChange;

  const TracksList(
      {Key? key,
      required this.tracks,
      required this.scrollController,
      required this.onTrackLikedChange})
      : super(key: key);

  void likeTrack(Track track) async {
    final prefs = await SharedPreferences.getInstance();
    await http.put(Uri.parse('$apiUrl/tracks/${track.id}/like'), headers: {
      'User-Agent': userAgent(),
      'Authorization': 'Bearer ${prefs.getString('token')}'
    });
  }

  void deleteTrackLike(Track track) async {
    final prefs = await SharedPreferences.getInstance();
    await http.delete(Uri.parse('$apiUrl/tracks/${track.id}/like'), headers: {
      'User-Agent': userAgent(),
      'Authorization': 'Bearer ${prefs.getString('token')}'
    });
  }

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
        padding: const EdgeInsets.all(8),
        controller: scrollController,
        shrinkWrap: scrollController == null,
        physics: scrollController == null
            ? const NeverScrollableScrollPhysics()
            : const ScrollPhysics(),
        itemCount: tracks.length,
        itemBuilder: (context, index) => InkWell(
            onTap: () async {
              try {
                final bool result = await playerChannel.invokeMethod(
                    'start', tracks[index].musicUrl);
                print('Result: $result');
              } on PlatformException catch (e) {
                print('Error: ${e.message}');
              }
            },
            child: Row(children: [
              SizedBox(
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
                                    tracks[index].album!.smallCoverUrl!)))),
                  )),
              const SizedBox(width: 16),
              Expanded(
                  flex: 1,
                  child: Column(children: [
                    SizedBox(
                        width: double.infinity,
                        child: Text(tracks[index].title,
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(
                                fontSize: 16, fontWeight: FontWeight.w500))),
                    const SizedBox(height: 4),
                    SizedBox(
                        width: double.infinity,
                        child: Text(
                            tracks[index]
                                .artists!
                                .map(
                                  (artist) => artist.name,
                                )
                                .join(', '),
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(
                                fontSize: 16, color: Colors.grey))),
                  ])),
              const SizedBox(width: 16),
              IconButton(
                  onPressed: () {
                    if (tracks[index].liked!) {
                      deleteTrackLike(tracks[index]);
                      onTrackLikedChange(index, false);
                    } else {
                      likeTrack(tracks[index]);
                      onTrackLikedChange(index, true);
                    }
                  },
                  color: tracks[index].liked! ? Colors.red : null,
                  icon: Icon(tracks[index].liked!
                      ? Icons.favorite
                      : Icons.favorite_outline)),
            ])));
  }
}
