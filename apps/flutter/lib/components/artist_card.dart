import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import '../models/artist.dart';

class ArtistCard extends StatelessWidget {
  final Artist artist;

  const ArtistCard({Key? key, required this.artist}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        elevation: 5,
        clipBehavior: Clip.antiAliasWithSaveLayer,
        child: InkWell(
          onTap: () =>
              Navigator.pushNamed(context, '/artist', arguments: artist),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              AspectRatio(
                  aspectRatio: 1,
                  child: Image(
                    image: CachedNetworkImageProvider(artist.mediumImageUrl!),
                    fit: BoxFit.fill,
                  )),
              Padding(
                  padding: const EdgeInsets.all(16),
                  child:  Text(artist.name,
                          overflow: TextOverflow.ellipsis,
                          style: const TextStyle(fontWeight: FontWeight.w500)),
                    )
            ],
          ),
        ));
  }
}
