//
//  ContentView.swift
//  BassieMusic
//
//  Created by Bastiaan van der Plaat on 03/07/2022.
//

import SwiftUI

struct MenuView : View {
    @State var albumsIsActive = true
    
    var body: some View {
        List {
            Text("BassieMusic")
                .font(.title)
                .padding(.vertical)
            
            NavigationLink(destination: ArtistsView()) {
                Image(systemName: "music.mic").frame(width: 16, height: 16)
                Text("Artists")
            }
            
            NavigationLink(destination: AlbumsView(), isActive: $albumsIsActive) {
                Image(systemName: "rectangle.stack").frame(width: 16, height: 16)
                Text("Albums")
            }
            NavigationLink(destination: TracksView()) {
                Image(systemName: "music.note").frame(width: 16, height: 16)
                Text("Tracks")
            }
        }
    }
}

struct ContentView: View {
    @State var progress = 0.0
    
    var body: some View {
        NavigationView {
            MenuView()
            AlbumsView()
        }
        .toolbar {
            Button(action: { print("test") }) {
                Label("Previous", systemImage: "backward.end.fill")
            }
            Button(action: { print("test") }) {
                Label("Seek back", systemImage: "backward.fill")
            }
            Button(action: { print("test") }) {
                Label("Play", systemImage: "play.fill")
            }
            Button(action: { print("test") }) {
                Label("Seek forward", systemImage: "forward.fill")
            }
            Button(action: { print("test") }) {
                Label("Next", systemImage: "forward.end.fill")
            }
            Text("10:00")
            Slider(value: $progress, in: 0...100, onEditingChanged: { editing in
                print(editing)
            }).frame(width:200)
            Text("10:00")
        }
        .frame(minWidth: 640, idealWidth: 1280, minHeight: 480, idealHeight: 720)
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
