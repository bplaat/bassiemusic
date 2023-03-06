import 'track.dart';
import 'album.dart';

class Artist {
  final String id;
  final String name;
  final String? smallImageUrl;
  final String? mediumImageUrl;
  final String? largeImageUrl;
  final bool? liked;
  final DateTime createdAt;
  final List<Track>? topTracks;
  final List<Album>? albums;

  Artist(
      {required this.id,
      required this.name,
      required this.smallImageUrl,
      required this.mediumImageUrl,
      required this.largeImageUrl,
      required this.liked,
      required this.createdAt,
      required this.topTracks,
      required this.albums});

  factory Artist.fromJson(Map<String, dynamic> json) {
    return Artist(
        id: json['id'],
        name: json['name'],
        smallImageUrl: json['small_image'],
        mediumImageUrl: json['medium_image'],
        largeImageUrl: json['large_image'],
        liked: json['liked'],
        createdAt: DateTime.parse(json['created_at']),
        topTracks: json['top_tracks']
            ?.map<Track>((json) => Track.fromJson(json))
            .toList(),
        albums: json['albums']
            ?.map<Album>((json) => Album.fromJson(json))
            .toList());
  }
}
