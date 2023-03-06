import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import '../models/user.dart';
import '../config.dart';
import '../utils.dart';
import 'home_page.dart';
import 'login_page.dart';

class RootPage extends StatefulWidget {
  const RootPage({super.key});

  @override
  State<RootPage> createState() => _RootPageState();
}

class _RootPageState extends State<RootPage> {
  bool _userLoaded = false;
  User? _user;

  Future<bool> checkAuth() async {
    if (_userLoaded) return _user != null;
    _userLoaded = true;
    final prefs = await SharedPreferences.getInstance();
    if (prefs.getString('token') != null) {
      final response = await http.get(Uri.parse('$apiUrl/auth/validate'),
          headers: {
            'User-Agent': userAgent(),
            'Authorization': 'Bearer ${prefs.getString('token')}'
          });
      final data = json.decode(response.body);
      if (data.containsKey('user')) {
        _user = User.fromJson(data['user']);
        return true;
      }
    }
    return false;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: FutureBuilder(
          future: checkAuth(),
          builder: (context, snapshot) {
            if (snapshot.hasError) {
              return Center(child: Text("Error: ${snapshot.error}"));
            } else if (snapshot.hasData) {
              return snapshot.data! && _user != null
                  ? HomePage(
                      user: _user!,
                      onAuthChange: (User? user) {
                        setState(() => _user = user);
                      })
                  : LoginPage(onAuthChange: (User? user) {
                      setState(() => _user = user);
                    });
            }
            return const Center(child: CircularProgressIndicator());
          },
        ));
  }
}
