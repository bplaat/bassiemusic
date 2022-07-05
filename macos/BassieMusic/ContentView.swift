//
//  ContentView.swift
//  BassieMusic
//
//  Created by Bastiaan van der Plaat on 03/07/2022.
//

import SwiftUI
import AVFoundation

struct ContentView: View {
    @State var track: Track?
    @State var progress = 0.0
    @State var volume = 100.0
    @State var albumsIsActive = true
    @State var audioPlayer: AVPlayer?
    
    private func toggleSidebar() {
        NSApp.keyWindow?.firstResponder?
            .tryToPerform(#selector(NSSplitViewController.toggleSidebar(_:)), with: nil)
    }
    
    private func playTrack(track: Track) {
        self.track = track
        self.audioPlayer?.pause();
        self.audioPlayer = AVPlayer(url: URL(string: self.track!.music)!)
        self.audioPlayer?.play()
    }
    
    private func seekBackward() {
        self.audioPlayer!.seek(to: self.audioPlayer!.currentTime() - CMTime(seconds: 10, preferredTimescale: 1))
    }
    
    private func playPause() {
        if self.audioPlayer!.rate > 0 {
            self.audioPlayer!.pause();
        } else {
            self.audioPlayer!.play();
        }
    }
    
    private func seekForward() {
        self.audioPlayer!.seek(to: self.audioPlayer!.currentTime() + CMTime(seconds: 10, preferredTimescale: 1))
    }
    
    
    var body: some View {
        NavigationView {
            
                List {
                    Text("BassieMusic")
                        .font(.title2)
                        .padding(EdgeInsets(top: 0, leading: 0, bottom: 8, trailing: 0))
                    
                    NavigationLink(destination: ArtistsView()) {
                        Image(systemName: "music.mic").frame(width: 16, height: 16)
                        Text("Artists")
                    }
                    
                    NavigationLink(destination: AlbumsView(), isActive: $albumsIsActive) {
                        Image(systemName: "rectangle.stack").frame(width: 16, height: 16)
                        Text("Albums")
                    }
                    NavigationLink(destination: TracksView(playTrack: { track in
                        self.playTrack(track: track)
                    })) {
                        Image(systemName: "music.note").frame(width: 16, height: 16)
                        Text("Tracks")
                    }
                }.listStyle(SidebarListStyle())
            
            AlbumsView()
        }
        .toolbar {
            ToolbarItemGroup(placement: .primaryAction) {
                Button(action: toggleSidebar) {
                    Image(systemName: "sidebar.left")
                        .help("Toggle Sidebar")
                }

                HStack {
                    if track != nil {
                        VStack {
                            Text(track!.title)
                            Text(track!.artists.map(\.name).joined(separator: ", "))
                        }
                    }
                }
                
            }
            
            ToolbarItemGroup(placement: .principal) {
                if self.audioPlayer != nil {
                    Button(action: { print("Previous track") }) {
                        Image(systemName: "backward.end.fill")
                            .help("Previous track")
                    }
                    Button(action: seekBackward) {
                        Image(systemName: "backward.fill")
                            .help("Seek backward")
                    }
                    Button(action: playPause) {
                        Image(systemName: self.audioPlayer!.rate > 0 ? "pause.fill" : "play.fill")
                            .help(self.audioPlayer!.rate > 0 ? "Pause" : "Play")
                    }
                    Button(action: seekForward) {
                        Image(systemName: "forward.fill")
                            .help("Seek forward")
                    }
                    Button(action: { print("Next track") }) {
                        Image(systemName: "forward.end.fill")
                            .help("Next track")
                    }
                    Slider(value: $progress, in: 0...100) {

                    } minimumValueLabel: {
                        Text("00:00")
                    } maximumValueLabel: {
                        Text("10:00")
                    } onEditingChanged: { isEditing in

                    }
                    .frame(minWidth:200)
                }
            }
            
            ToolbarItemGroup(placement: .principal) {
                Slider(value: $volume, in: 0...100, onEditingChanged: { editing in
                    print(editing)
                }).frame(width:100)

                Button(action: {
                    self.volume = 0
                }) {
                    Image(systemName: "speaker.wave.3.fill")
                        .help("Volume")
                }
            }
        }
        .frame(minWidth: 640, idealWidth: 1280, minHeight: 480, idealHeight: 720)
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
