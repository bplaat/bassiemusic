import 'package:flutter/material.dart';
import 'pages/home_page.dart';

class BassieMusicApp extends StatelessWidget {
  const BassieMusicApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: 'BassieMusic',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
            brightness: Brightness.dark,
            primarySwatch: Colors.blue,
            accentColor: Colors.blue,
            appBarTheme: const AppBarTheme(
              foregroundColor: Colors.white,
              backgroundColor: Color(0xff121212),
            ),
            cardTheme: Theme.of(context)
                .cardTheme
                .copyWith(color: const Color(0xff121212)),
            bottomNavigationBarTheme: Theme.of(context)
                .bottomNavigationBarTheme
                .copyWith(backgroundColor: const Color(0xff121212)),
            colorScheme: Theme.of(context).colorScheme.copyWith(
                brightness: Brightness.dark,
                background: const Color(0xff0a0a0a))),
        home: const HomePage());
  }
}

void main() {
  runApp(const BassieMusicApp());
}
