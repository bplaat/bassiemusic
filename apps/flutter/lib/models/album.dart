import 'artist.dart';
import 'genre.dart';
import 'track.dart';

class Album {
  final String id;
  final String title;
  final DateTime releasedAt;
  final String? smallCoverUrl;
  final String? mediumCoverUrl;
  final String? largeCoverUrl;
  final bool explicit;
  final bool? liked;
  final DateTime createdAt;
  final List<Artist>? artists;
  final List<Genre>? genres;
  final List<Track>? tracks;

  Album({
    required this.id,
    required this.title,
    required this.releasedAt,
    required this.explicit,
    required this.smallCoverUrl,
    required this.mediumCoverUrl,
    required this.largeCoverUrl,
    required this.liked,
    required this.createdAt,
    required this.artists,
    required this.genres,
    required this.tracks,
  });

  factory Album.fromJson(Map<String, dynamic> json) {
    return Album(
        id: json['id'],
        title: json['title'],
        releasedAt: DateTime.parse(json['released_at']),
        smallCoverUrl: json['small_cover'],
        mediumCoverUrl: json['medium_cover'],
        largeCoverUrl: json['large_cover'],
        explicit: json['explicit'],
        liked: json['liked'],
        createdAt: DateTime.parse(json['created_at']),
        artists: json['artists']
            ?.map<Artist>((json) => Artist.fromJson(json))
            .toList(),
        genres:
            json['genres']?.map<Genre>((json) => Genre.fromJson(json)).toList(),
        tracks: json['tracks']
            ?.map<Track>((json) => Track.fromJson(json))
            .toList());
  }
}
