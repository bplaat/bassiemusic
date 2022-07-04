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
            NavigationLink(destination: TracksView()) {
                Image(systemName: "music.note").frame(width: 16, height: 16)
                Text("Tracks")
            }
        }.listStyle(SidebarListStyle())
        
    }
}

struct ContentView: View {
    @State var progress = 0.0
    @State var volume = 100.0
    
    private func toggleSidebar() {
        NSApp.keyWindow?.firstResponder?
            .tryToPerform(#selector(NSSplitViewController.toggleSidebar(_:)), with: nil)
    }
    
    var body: some View {
        NavigationView {
            MenuView()
            AlbumsView()
        }
        .toolbar {
            
            ToolbarItemGroup(placement: .primaryAction) {
                Button(action: toggleSidebar) {
                    Image(systemName: "sidebar.left")
                        .help("Toggle Sidebar")
                }

                HStack {
                    VStack {
                        Text("Track")
                        Text("Artist")
                    }
                }
                
            }
            
            ToolbarItemGroup(placement: .principal) {
                Button(action: { print("Previous track") }) {
                    Image(systemName: "backward.end.fill")
                        .help("Previous track")
                }
                Button(action: { print("Seek backward") }) {
                    Image(systemName: "backward.fill")
                        .help("Seek backward")
                }
                Button(action: { print("Play") }) {
                    Image(systemName: "play.fill")
                        .help("Play")
                }
                Button(action: { print("Seek forward") }) {
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
            
            ToolbarItemGroup(placement: .principal) {
                Slider(value: $volume, in: 0...100, onEditingChanged: { editing in
                    print(editing)
                }).frame(width:100)

                Button(action: { print("Volume") }) {
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
