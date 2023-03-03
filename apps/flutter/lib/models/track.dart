import 'album.dart';
import 'artist.dart';

class Track {
  final String id;
  final String title;
  final int disk;
  final int position;
  final double duration;
  final bool explicit;
  final bool? liked;
  final DateTime createdAt;
  final Album? album;
  final List<Artist>? artists;

  Track({
    required this.id,
    required this.title,
    required this.disk,
    required this.position,
    required this.duration,
    required this.explicit,
    required this.liked,
    required this.createdAt,
    required this.album,
    required this.artists,
  });

  factory Track.fromJson(Map<String, dynamic> json) {
    return Track(
        id: json['id'],
        title: json['title'],
        disk: json['disk'],
        position: json['position'],
        duration: json['duration'].toDouble(),
        explicit: json['explicit'],
        liked: json['liked'],
        createdAt: DateTime.parse(json['created_at']),
        album: json['album'] != null ? Album.fromJson(json['album']) : null,
        artists: json['artists']
            ?.map<Artist>((json) => Artist.fromJson(json))
            .toList());
  }
}
