import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../models/user.dart';
import '../models/artist.dart';
import '../models/album.dart';
import '../models/genre.dart';
import 'search_page.dart';
import 'profile_page.dart';
import 'artist_page.dart';
import 'genre_page.dart';
import 'album_page.dart';
import 'home_home_tab.dart';
import 'home_explore_tab.dart';
import 'home_liked_tab.dart';
import 'home_history_tab.dart';

class HomePage extends StatefulWidget {
  final User user;
  final Function(User? user) onAuthChange;

  const HomePage({super.key, required this.user, required this.onAuthChange});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _navigatorKey = GlobalKey<NavigatorState>();
  final _pageController = PageController(initialPage: 0);
  int _page = 0;

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final lang = AppLocalizations.of(context)!;
    return Scaffold(
        body: Navigator(
            key: _navigatorKey,
            initialRoute: '/',
            onGenerateRoute: (RouteSettings settings) {
              late Widget page;
              if (settings.name == "/search") {
                page = const SearchPage();
              } else if (settings.name == "/profile") {
                page = ProfilePage(
                    user: widget.user, onAuthChange: widget.onAuthChange);
              } else if (settings.name == "/artist") {
                page =
                    ArtistPage(incompleteArtist: settings.arguments as Artist);
              } else if (settings.name == "/album") {
                page = AlbumPage(incompleteAlbum: settings.arguments as Album);
              } else if (settings.name == "/genre") {
                page = GenrePage(genre: settings.arguments as Genre);
              } else {
                page = Stack(children: [
                  Padding(
                      padding:
                          EdgeInsets.only(top: AppBar().preferredSize.height),
                      child: PageView(
                          controller: _pageController,
                          onPageChanged: (index) {
                            setState(() => _page = index);
                          },
                          children: [
                            HomeHomeTab(user: widget.user),
                            const HomeExplorerTab(),
                            HomeLikedTab(user: widget.user),
                            HomeHistoryTab(user: widget.user)
                          ])),
                  Positioned(
                      //Place it at the top, and not use the entire screen
                      top: 0.0,
                      left: 0.0,
                      right: 0.0,
                      child: AppBar(
                          title: Text(lang.app_name),
                          elevation: 4,
                          actions: [
                            IconButton(
                              onPressed: () => _navigatorKey.currentState!
                                  .pushNamed('/search'),
                              icon: const Icon(Icons.search),
                            ),
                            Container(
                                margin: const EdgeInsets.all(4),
                                child: SizedBox(
                                    width: 48,
                                    height: 48,
                                    child: Card(
                                        clipBehavior:
                                            Clip.antiAliasWithSaveLayer,
                                        shape: RoundedRectangleBorder(
                                          borderRadius:
                                              BorderRadius.circular(6),
                                        ),
                                        elevation: 2,
                                        child: InkWell(
                                          onTap: () => _navigatorKey
                                              .currentState!
                                              .pushNamed('/profile'),
                                          child: Container(
                                              decoration: BoxDecoration(
                                                  image: DecorationImage(
                                                      fit: BoxFit.cover,
                                                      image: CachedNetworkImageProvider(
                                                          widget.user
                                                              .smallAvatarUrl!)))),
                                        )))),
                          ])),
                ]);
              }
              return MaterialPageRoute<dynamic>(
                builder: (context) {
                  return page;
                },
                settings: settings,
              );
            }),
        bottomNavigationBar: Material(
            elevation: 8,
            child: Column(
                mainAxisSize: MainAxisSize.min,
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  Container(
                      color: Theme.of(context)
                          .bottomNavigationBarTheme
                          .backgroundColor,
                      child: Padding(
                          padding: const EdgeInsets.all(8),
                          child: Row(children: [
                            SizedBox(
                                width: 56,
                                height: 56,
                                child: Card(
                                  clipBehavior: Clip.antiAliasWithSaveLayer,
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(6),
                                  ),
                                  elevation: 2,
                                  child: Container(
                                      decoration: const BoxDecoration(
                                          image: DecorationImage(
                                              fit: BoxFit.cover,
                                              image: CachedNetworkImageProvider(
                                                  "https://bassiemusic-storage.plaatsoft.nl/albums/medium/b7ba551d-e28e-47a2-8e5e-0bdd80d7bb1b.jpg")))),
                                )),
                            const SizedBox(width: 16),
                            Expanded(
                                flex: 1,
                                child: Column(children: const [
                                  SizedBox(
                                      width: double.infinity,
                                      child: Text("Flower Boy",
                                          style: TextStyle(
                                              fontSize: 16,
                                              fontWeight: FontWeight.w500))),
                                  SizedBox(height: 4),
                                  SizedBox(
                                      width: double.infinity,
                                      child: Text("Tyler, The Creator",
                                          style: TextStyle(
                                              fontSize: 16,
                                              color: Colors.grey))),
                                ])),
                            const SizedBox(width: 16),
                            IconButton(
                                onPressed: () => {},
                                icon: const Icon(Icons.skip_previous)),
                            IconButton(
                                onPressed: () => {},
                                icon: const Icon(Icons.play_arrow)),
                            IconButton(
                                onPressed: () => {},
                                icon: const Icon(Icons.skip_next)),
                            IconButton(
                                onPressed: () => {},
                                icon: const Icon(Icons.favorite_outline))
                          ]))),
                  BottomNavigationBar(
                      elevation: 0,
                      type: BottomNavigationBarType.fixed,
                      onTap: (index) {
                        if (_navigatorKey.currentState!.canPop()) {
                          _navigatorKey.currentState!.pop();
                        }
                        _pageController.animateToPage(index,
                            duration: const Duration(milliseconds: 200),
                            curve: Curves.ease);
                        setState(() => _page = index);
                      },
                      currentIndex: _page,
                      items: [
                        BottomNavigationBarItem(
                            icon: const Icon(Icons.home),
                            label: lang.home_home),
                        BottomNavigationBarItem(
                            icon: const Icon(Icons.explore),
                            label: lang.home_explore),
                        BottomNavigationBarItem(
                            icon: const Icon(Icons.favorite),
                            label: lang.home_liked),
                        BottomNavigationBarItem(
                            icon: const Icon(Icons.history),
                            label: lang.home_history)
                      ])
                ])));
  }
}
