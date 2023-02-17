import 'package:flutter/material.dart';

class HomeLikedTab extends StatefulWidget {
  const HomeLikedTab({super.key});

  @override
  State<HomeLikedTab> createState() => _HomeLikedTabState();
}

class _HomeLikedTabState extends State<HomeLikedTab> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(vsync: this, length: 3);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }


  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: Theme.of(context).colorScheme.background,
        appBar: AppBar(
            title: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(
              text: "Artists",
            ),
            Tab(
              text: "Albums",
            ),
            Tab(
              text: "Tracks",
            ),
          ],
        )), body: Center(child: Text('Liked')));
  }
}
