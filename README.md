# BassieMusic
A new online streaming music service thingy

## Dependencies
- Go Compiler
- MySQL
- Node.js
- yt-dlp in PATH
- ffmpeg in PATH

## TODO
In progress:
- Fix many bugs
- bplaat: Popup shows music queue and you can add, reorder and remove songs from queue

Backlog:
- lplaat: Hover track play button but dblclick still needs to work
- lplaat: User account sessions management UI like GitHub
- Discover other peoples public playlists via global search and like them to add them to your sidebar
- bplaat: Users can create, edit, delete find public and private playlists
- Change and remove artists, albums, genres, tracks in web interface when admin
- Move track/play messages to a websocket connection to reduce load
- Better downloader experience with live progress via websockets
- Deploy 0.1 version on local server at bassiemusic.ml and api.bassiemusic.ml

Long term:
- native Android streaming app 100% online
- native Android app offline download features
- macOS webview wrapper app
- basic native macOS app written in objc with cocoa
