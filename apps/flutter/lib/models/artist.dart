class Artist {
  final String id;
  final String name;
  final DateTime createdAt;

  Artist({
    required this.id,
    required this.name,
    required this.createdAt,
  });

  factory Artist.fromJson(Map<String, dynamic> json) {
    return Artist(
      id: json['id'],
      name: json['name'],
      createdAt: DateTime.parse(json['created_at'])
    );
  }
}
