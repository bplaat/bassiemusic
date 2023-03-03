import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import '../models/album.dart';

class AlbumCard extends StatelessWidget {
  final Album album;

  const AlbumCard({Key? key, required this.album}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        elevation: 5,
        clipBehavior: Clip.antiAliasWithSaveLayer,
        child: InkWell(
          onTap: () => Navigator.pushNamed(context, '/albums'),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Image(
                image: CachedNetworkImageProvider(album.mediumCoverUrl),
                fit: BoxFit.fill,
              ),
              Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Container(
                          margin: const EdgeInsets.only(bottom: 8),
                          child: Text(album.title,
                              overflow: TextOverflow.ellipsis,
                              style: const TextStyle(
                                  fontWeight: FontWeight.w500))),
                      Text(
                          album.artists!
                              .map(
                                (artist) => artist.name,
                              )
                              .join(', '),
                          overflow: TextOverflow.ellipsis,
                          style: const TextStyle(color: Colors.grey)),
                    ],
                  ))
            ],
          ),
        ));
  }
}
