import 'package:flutter/material.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'pages/root_page.dart';

class BassieMusicApp extends StatelessWidget {
  const BassieMusicApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: 'BassieMusic',
        debugShowCheckedModeBanner: false,
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        theme: ThemeData(
          brightness: Brightness.light,
          primarySwatch: Colors.blue,
          // accentColor: Colors.blue,
          appBarTheme: const AppBarTheme(
            foregroundColor: Color(0xff121212),
            backgroundColor: Colors.white,
          ),
          bottomNavigationBarTheme: Theme.of(context)
              .bottomNavigationBarTheme
              .copyWith(backgroundColor: Colors.white),
        ),
        darkTheme: ThemeData(
          brightness: Brightness.dark,
          primarySwatch: Colors.blue,
          // accentColor: Colors.blue,
          scaffoldBackgroundColor: const Color(0xff0a0a0a),
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
        ),
        home: const RootPage());
  }
}

void main() {
  runApp(const BassieMusicApp());
}
