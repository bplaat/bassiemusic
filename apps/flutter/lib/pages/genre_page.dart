import 'package:flutter/material.dart';

class GenrePage extends StatefulWidget {
  const GenrePage({super.key});

  @override
  State<GenrePage> createState() => _GenrePageState();
}

class _GenrePageState extends State<GenrePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: Theme.of(context).colorScheme.background,
        appBar: AppBar(
            title: const Text('Genre'),
        ),
        body: Center(child: Text("Genre"))
    );
  }
}
