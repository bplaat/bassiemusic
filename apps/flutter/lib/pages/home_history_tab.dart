import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../components/tracks_list.dart';
import '../models/user.dart';
import '../models/track.dart';
import '../config.dart';
import '../utils.dart';

class HomeHistoryTab extends StatefulWidget {
  final User user;

  const HomeHistoryTab({super.key, required this.user});

  @override
  State<HomeHistoryTab> createState() => _HomeHistoryTabState();
}

class _HomeHistoryTabState extends State<HomeHistoryTab>
    with AutomaticKeepAliveClientMixin<HomeHistoryTab> {
  final ScrollController _scrollController = ScrollController();
  List<Track> _tracks = [];
  int _page = 1;
  bool _isLoading = false;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    loadPage();
    _scrollController.addListener(() {
      if (!_isLoading &&
          _scrollController.position.pixels >
              _scrollController.position.maxScrollExtent * 0.9) {
        loadPage();
      }
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void loadPage() async {
    _isLoading = true;
    try {
      final prefs = await SharedPreferences.getInstance();
      final response = await http.get(
          Uri.parse(
              '$apiUrl/users/${widget.user.id}/played_tracks?page=${_page++}'),
          headers: {
            'User-Agent': userAgent(),
            'Authorization': 'Bearer ${prefs.getString('token')}'
          });

      if (response.statusCode == 200) {
        final tracksJson = json.decode(utf8.decode(response.bodyBytes))['data'];
        List<Track> newTracks =
            tracksJson.map<Track>((json) => Track.fromJson(json)).toList();
        _tracks.addAll(newTracks);
        _isLoading = false;
        setState(() => _tracks = _tracks);
      }
    } catch (exception) {
      print('Error: ${exception.toString()}');
    }
  }

  @override
  Widget build(BuildContext context) {
    final lang = AppLocalizations.of(context)!;
    return RefreshIndicator(
        onRefresh: () async {
          _tracks = [];
          _page = 1;
          _isLoading = false;
          loadPage();
        },
        child: _tracks.isNotEmpty
            ? TracksList(
                scrollController: _scrollController,
                tracks: _tracks,
                onTrackLikedChange: (index, liked) => setState(() {
                      _tracks[index].liked = liked;
                    }))
            : (_isLoading
                ? const Center(child: CircularProgressIndicator())
                : Center(
                    child: Text(lang.home_history_empty),
                  )));
  }
}
