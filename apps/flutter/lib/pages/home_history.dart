import 'package:flutter/material.dart';

class HomeHistoryTab extends StatefulWidget {
  const HomeHistoryTab({super.key});

  @override
  State<HomeHistoryTab> createState() => _HomeHistoryTabState();
}

class _HomeHistoryTabState extends State<HomeHistoryTab> {
  @override
  Widget build(BuildContext context) {
    return Center(child: Text('History'));
  }
}
