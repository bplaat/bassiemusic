class User {
  final String id;
  final String username;
  final String email;
  final String? smallAvatarUrl;
  final String? mediumAvatarUrl;
  final bool allowExplicit;
  final String role;
  final String language;
  final String theme;
  final DateTime createdAt;

  User({
    required this.id,
    required this.username,
    required this.email,
    required this.smallAvatarUrl,
    required this.mediumAvatarUrl,
    required this.allowExplicit,
    required this.role,
    required this.language,
    required this.theme,
    required this.createdAt,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
        id: json['id'],
        username: json['username'],
        email: json['email'],
        smallAvatarUrl: json['small_avatar'],
        mediumAvatarUrl: json['medium_avatar'],
        allowExplicit: json['allow_explicit'],
        role: json['role'],
        language: json['language'],
        theme: json['theme'],
        createdAt: DateTime.parse(json['created_at']));
  }
}
