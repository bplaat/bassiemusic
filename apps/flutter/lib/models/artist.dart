class Artist {
  final String id;
  final String name;
  final String smallImageUrl;
  final String mediumImageUrl;
  final String largeImageUrl;
  final bool? liked;
  final DateTime createdAt;

  Artist({
    required this.id,
    required this.name,
    required this.smallImageUrl,
    required this.mediumImageUrl,
    required this.largeImageUrl,
    required this.liked,
    required this.createdAt,
  });

  factory Artist.fromJson(Map<String, dynamic> json) {
    return Artist(
        id: json['id'],
        name: json['name'],
        smallImageUrl: json['small_image'],
        mediumImageUrl: json['medium_image'],
        largeImageUrl: json['large_image'],
        liked: json['liked'],
        createdAt: DateTime.parse(json['created_at']));
  }
}
