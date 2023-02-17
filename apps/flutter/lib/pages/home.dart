import 'package:flutter/material.dart';
import 'home_home.dart';
import 'home_explore.dart';
import 'home_liked.dart';
import 'home_history.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _pageController = PageController(initialPage: 1);
  int _currentPageIndex = 1;

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: Theme.of(context).colorScheme.background,
        appBar: AppBar(
            title: const Text('BassieMusic'),
            elevation: _currentPageIndex == 1 || _currentPageIndex == 2 ? 0 : 4,
            actions: [
              IconButton(
                onPressed: () => {},
                icon: const Icon(Icons.search),
              ),
              CircleAvatar(
                backgroundColor: Colors.brown.shade800,
                child: const Text('BP'),
              )
            ]),
        body: PageView(
            controller: _pageController,
            onPageChanged: (index) {
              setState(() => _currentPageIndex = index);
            },
            children: const [
              HomeHomeTab(),
              HomeExplorerTab(),
              HomeLikedTab(),
              HomeHistoryTab()
            ]),
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
        bottomNavigationBar: BottomNavigationBar(
            type: BottomNavigationBarType.fixed,
            onTap: (index) {
              _pageController.animateToPage(index,
                  duration: const Duration(milliseconds: 300),
                  curve: Curves.ease);
              setState(() => _currentPageIndex = index);
            },
            currentIndex: _currentPageIndex,
            items: const [
              BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
              BottomNavigationBarItem(
                  icon: Icon(Icons.explore), label: 'Explore'),
              BottomNavigationBarItem(
                  icon: Icon(Icons.thumb_up), label: 'Liked'),
              BottomNavigationBarItem(
                  icon: Icon(Icons.history), label: 'History')
            ]));
  }
}
