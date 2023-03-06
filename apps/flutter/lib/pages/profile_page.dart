import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../models/user.dart';
import '../config.dart';
import '../utils.dart';

class ProfilePage extends StatefulWidget {
  final User user;
  final Function(User? user) onAuthChange;

  const ProfilePage(
      {super.key, required this.user, required this.onAuthChange});

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  void logout() async {
    final prefs = await SharedPreferences.getInstance();
    await http.put(Uri.parse('$apiUrl/auth/logout'), headers: {
      'User-Agent': userAgent(),
      'Authorization': 'Bearer ${prefs.getString('token')}'
    });
    await prefs.remove('token');
    widget.onAuthChange(null);
  }

  @override
  Widget build(BuildContext context) {
    final lang = AppLocalizations.of(context)!;
    return Scaffold(
        appBar: AppBar(
          title: Text(lang.profile_header),
        ),
        body: SingleChildScrollView(
            child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(children: [
                  // Logout button
                  Container(
                      margin: const EdgeInsets.symmetric(vertical: 8),
                      child: SizedBox(
                          width: double.infinity,
                          child: ElevatedButton(
                              onPressed: logout,
                              style: ElevatedButton.styleFrom(
                                  primary: Colors.blue,
                                  shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(48)),
                                  padding: const EdgeInsets.symmetric(
                                      horizontal: 24, vertical: 16)),
                              child: Text(lang.profile_logout,
                                  style: const TextStyle(
                                      color: Colors.white, fontSize: 18))))),
                ]))));
  }
}
