import 'package:flutter/material.dart';

class HomeLikedTab extends StatefulWidget {
  const HomeLikedTab({super.key});

  @override
  State<HomeLikedTab> createState() => _HomeLikedTabState();
}

class _HomeLikedTabState extends State<HomeLikedTab>
    with SingleTickerProviderStateMixin, AutomaticKeepAliveClientMixin<HomeLikedTab> {
  late TabController _tabController;

  @override
  bool get wantKeepAlive => true;

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
        )),
        body: TabBarView(controller: _tabController, children: [
          Center(child: Text('Liked Artists')),
          Center(child: Text('Liked Albums')),
          Center(child: Text('Liked Tracks'))
        ]));
  }
}
