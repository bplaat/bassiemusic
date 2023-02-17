import 'package:flutter/material.dart';

class AlbumPage extends StatefulWidget {
  const AlbumPage({super.key});

  @override
  State<AlbumPage> createState() => _AlbumPageState();
}

class _AlbumPageState extends State<AlbumPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: Theme.of(context).colorScheme.background,
        appBar: AppBar(
            title: const Text('Album'),
        ),
        body: Center(child: Text("Album"))
    );
  }
}
