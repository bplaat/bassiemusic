import 'package:bassiemusic/pages/genre_page.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
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
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _pageController = PageController(initialPage: 0);
  int _currentPageIndex = 0;

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  final _navigatorKey = GlobalKey<NavigatorState>();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: Theme.of(context).colorScheme.background,

        //---------------------
        body: Navigator(
            key: _navigatorKey,
            initialRoute: '/',
            onGenerateRoute: (RouteSettings settings) {
              late Widget page;
              if (settings.name == "/search") {
                page = const SearchPage();
              } else if (settings.name == "/profile") {
                page = const ProfilePage();
              } else if (settings.name == "/artists") {
                page = const ArtistPage();
              } else if (settings.name == "/genres") {
                page = const GenrePage();
              } else if (settings.name == "/albums") {
                page = const AlbumPage();
              } else {
                page = Stack(children: [
                  Padding(
                      padding:
                          EdgeInsets.only(top: AppBar().preferredSize.height),
                      child: PageView(
                          controller: _pageController,
                          onPageChanged: (index) {
                            setState(() => _currentPageIndex = index);
                          },
                          children: const [
                            HomeHomeTab(),
                            HomeExplorerTab(),
                            HomeLikedTab(),
                            HomeHistoryTab()
                          ])),
                  Positioned(
                      //Place it at the top, and not use the entire screen
                      top: 0.0,
                      left: 0.0,
                      right: 0.0,
                      child: AppBar(
                          title: const Text('BassieMusic'),
                          elevation:
                              _currentPageIndex == 1 || _currentPageIndex == 2
                                  ? 0
                                  : 4,
                          actions: [
                            IconButton(
                              onPressed: () => _navigatorKey.currentState!
                                  .pushNamed('/search'),
                              icon: const Icon(Icons.search),
                            ),
                            SizedBox(
                                width: 56,
                                height: 56,
                                child: Card(
                                    clipBehavior: Clip.antiAliasWithSaveLayer,
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(6),
                                    ),
                                    elevation: 2,
                                    child: InkWell(
                                      onTap: () => _navigatorKey.currentState!
                                          .pushNamed('/profile'),
                                      child: Container(
                                          decoration: BoxDecoration(
                                              image: DecorationImage(
                                                  fit: BoxFit.cover,
                                                  image: CachedNetworkImageProvider(
                                                      "https://bassiemusic-storage.plaatsoft.nl/avatars/c237187f-f029-414d-a682-5fef11afef1b.jpg")))),
                                    ))),
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
        //---------------------
        // floatingActionButton: Card(
        //     shape: RoundedRectangleBorder(
        //       borderRadius: BorderRadius.circular(8),
        //     ),
        //     elevation: 5,
        //     child: Padding(
        //         padding: const EdgeInsets.all(16),
        //         child: Row(
        //           children: [Text("1989"), Text("Taylor Swift")],
        //         ))),
        bottomNavigationBar: Column(
            mainAxisSize: MainAxisSize.min,
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              Padding(
                  padding: EdgeInsets.all(8),
                  child: Row(children: [
                    Container(
                        margin: EdgeInsets.only(right: 16),
                        child: SizedBox(
                            width: 56,
                            height: 56,
                            child: Card(
                              clipBehavior: Clip.antiAliasWithSaveLayer,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(6),
                              ),
                              elevation: 2,
                              child: Container(
                                  decoration: BoxDecoration(
                                      image: DecorationImage(
                                          fit: BoxFit.cover,
                                          image: CachedNetworkImageProvider(
                                              "https://bassiemusic-storage.plaatsoft.nl/albums/medium/b7ba551d-e28e-47a2-8e5e-0bdd80d7bb1b.jpg")))),
                            ))),
                    Expanded(
                        flex: 1,
                        child: Column(children: [
                          Container(
                              margin: EdgeInsets.only(bottom: 4),
                              child: SizedBox(
                                  width: double.infinity,
                                  child: Text("Flower Boy",
                                      style: TextStyle(
                                          fontSize: 16,
                                          fontWeight: FontWeight.w500)))),
                          Container(
                              margin: EdgeInsets.only(bottom: 4),
                              child: SizedBox(
                                  width: double.infinity,
                                  child: Text("Tyler, The Creator",
                                      style: TextStyle(
                                          fontSize: 16, color: Colors.grey)))),
                        ])),
                    IconButton(
                        onPressed: () => {}, icon: Icon(Icons.skip_previous)),
                    IconButton(
                        onPressed: () => {}, icon: Icon(Icons.play_arrow)),
                    IconButton(
                        onPressed: () => {}, icon: Icon(Icons.skip_next)),
                    IconButton(
                        onPressed: () => {}, icon: Icon(Icons.favorite_outline))
                  ])),
              BottomNavigationBar(
                  type: BottomNavigationBarType.fixed,
                  onTap: (index) {
                    if (_navigatorKey.currentState!.canPop()) {
                      _navigatorKey.currentState!.pop();
                    }

                    _pageController.animateToPage(index,
                        duration: const Duration(milliseconds: 200),
                        curve: Curves.ease);

                    setState(() => _currentPageIndex = index);
                  },
                  currentIndex: _currentPageIndex,
                  items: const [
                    BottomNavigationBarItem(
                        icon: Icon(Icons.home), label: 'Home'),
                    BottomNavigationBarItem(
                        icon: Icon(Icons.explore), label: 'Explore'),
                    BottomNavigationBarItem(
                        icon: Icon(Icons.favorite), label: 'Liked'),
                    BottomNavigationBarItem(
                        icon: Icon(Icons.history), label: 'History')
                  ])
            ]));
  }
}
