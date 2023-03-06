import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../models/user.dart';
import '../config.dart';
import '../utils.dart';

class LoginPage extends StatefulWidget {
  final Function(User? user) onAuthChange;

  const LoginPage({super.key, required this.onAuthChange});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _logonController = TextEditingController();
  final _passwordController = TextEditingController();

  bool _loading = false;
  bool _errors = false;

  @override
  void dispose() {
    _logonController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  void login() async {
    setState(() => _loading = true);

    final response = await http.post(Uri.parse('$apiUrl/auth/login'), headers: {
      'User-Agent': userAgent(),
    }, body: {
      'logon': _logonController.text,
      'password': _passwordController.text
    });
    final data = json.decode(response.body);
    if (!data.containsKey('token')) {
      setState(() {
        _loading = false;
        _errors = true;
      });
      return;
    }

    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('token', data['token']);

    User user = User.fromJson(data['user']);
    widget.onAuthChange(user);
  }

  @override
  Widget build(BuildContext context) {
    final lang = AppLocalizations.of(context)!;
    return Scaffold(
        appBar: AppBar(
          title: Text(lang.app_name),
        ),
        body: Center(
            child: SingleChildScrollView(
                child: Container(
                    constraints:
                        const BoxConstraints(maxWidth: double.infinity),
                    padding: const EdgeInsets.all(16),
                    child: Column(children: [
                      // Header and info
                      Container(
                          margin: const EdgeInsets.symmetric(vertical: 8),
                          child: Text(lang.login_header,
                              style: const TextStyle(
                                  fontSize: 32, fontWeight: FontWeight.w500))),
                      Container(
                          margin: const EdgeInsets.symmetric(vertical: 8),
                          child: Text(lang.login_info,
                              style: const TextStyle(
                                  color: Colors.grey, fontSize: 16),
                              textAlign: TextAlign.center)),

                      // Error
                      if (_errors) ...[
                        Container(
                            margin: const EdgeInsets.symmetric(vertical: 8),
                            child: Text(lang.login_error,
                                style: const TextStyle(
                                    fontSize: 16, color: Colors.red)))
                      ],

                      // Logon input
                      Container(
                          margin: const EdgeInsets.symmetric(vertical: 8),
                          child: TextFormField(
                              controller: _logonController,
                              onFieldSubmitted: (value) {
                                if (!_loading) login();
                              },
                              autocorrect: false,
                              style: const TextStyle(fontSize: 18),
                              decoration: InputDecoration(
                                  border: OutlineInputBorder(
                                      borderRadius: BorderRadius.circular(48)),
                                  contentPadding: const EdgeInsets.symmetric(
                                      horizontal: 24, vertical: 16),
                                  labelText: lang.login_logon))),

                      // Password input
                      Container(
                          margin: const EdgeInsets.symmetric(vertical: 8),
                          child: TextFormField(
                              controller: _passwordController,
                              onFieldSubmitted: (value) {
                                if (!_loading) login();
                              },
                              autocorrect: false,
                              obscureText: true,
                              style: const TextStyle(fontSize: 18),
                              decoration: InputDecoration(
                                  border: OutlineInputBorder(
                                      borderRadius: BorderRadius.circular(48)),
                                  contentPadding: const EdgeInsets.symmetric(
                                      horizontal: 24, vertical: 16),
                                  labelText: lang.login_password))),

                      // Login button
                      Container(
                          margin: const EdgeInsets.symmetric(vertical: 8),
                          child: SizedBox(
                              width: double.infinity,
                              child: ElevatedButton(
                                  onPressed: _loading ? null : login,
                                  style: ElevatedButton.styleFrom(
                                      primary: Colors.blue,
                                      shape: RoundedRectangleBorder(
                                          borderRadius:
                                              BorderRadius.circular(48)),
                                      padding: const EdgeInsets.symmetric(
                                          horizontal: 24, vertical: 16)),
                                  child: Text(lang.login_login,
                                      style: const TextStyle(
                                          color: Colors.white,
                                          fontSize: 18))))),

                      // Login footer
                      Container(
                          margin: const EdgeInsets.symmetric(vertical: 8),
                          child: Text(lang.login_footer,
                              style: const TextStyle(
                                  color: Colors.grey, fontSize: 16),
                              textAlign: TextAlign.center)),
                    ])))));
  }
}
