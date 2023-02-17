import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import '../models/genre.dart';

class GenreCard extends StatelessWidget {
  final Genre genre;

  const GenreCard({Key? key, required this.genre}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        elevation: 5,
        clipBehavior: Clip.antiAliasWithSaveLayer,
        child: InkWell(
          onTap: () => Navigator.pushNamed(context, '/genres'),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Image(
                image: CachedNetworkImageProvider(genre.mediumImageUrl),
                fit: BoxFit.fill,
              ),
              Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Container(
                          margin: const EdgeInsets.only(bottom: 8),
                          child: Text(genre.name,
                              overflow: TextOverflow.ellipsis,
                              style: const TextStyle(
                                  fontWeight: FontWeight.w500))),
                    ],
                  ))
            ],
          ),
        ));
  }
}
