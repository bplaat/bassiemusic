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
    var body: some View {
        NavigationView {
            MenuView()
            AlbumsView()
        }
        .frame(minWidth: 640, idealWidth: 1280, minHeight: 480, idealHeight: 720)
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
