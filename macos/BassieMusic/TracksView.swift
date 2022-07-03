//
//  TracksView.swift
//  BassieMusic
//
//  Created by Bastiaan van der Plaat on 03/07/2022.
//

import SwiftUI
import AVFoundation

struct Track : Identifiable, Decodable {
    var id: String
    var title: String
    var disk: Int
    var position: Int
    var duration: Int
    var music: String
    var album: Album
    var artists: [Artist]
    var created_at: String
    var updated_at: String
}

class FetchTracks: ObservableObject {
    @Published var tracks = [Track]()

    init() {
        let url = URL(string: "http://localhost:8080/api/tracks?limit=50")!
        URLSession.shared.dataTask(with: url) {(data, response, error) in
            do {
                if let todoData = data {
                    let decodedData = try JSONDecoder().decode([Track].self, from: todoData)
                    DispatchQueue.main.async {
                        self.tracks = decodedData
                    }
                } else {
                    print("No data")
                }
            } catch {
                print("Error")
            }
        }.resume()
    }
}

struct TracksView: View {
    @ObservedObject var fetchTracks = FetchTracks()
    
    @State var audioPlayer: AVPlayer?
    @State var selectedTracks = Set<Track.ID>()
    
    var body: some View {
        Table(fetchTracks.tracks, selection: $selectedTracks) {
            TableColumn("Artists") {
                Text($0.artists.map(\.name).joined(separator: ", "))
            }
            TableColumn("Album") {
                Text($0.album.title)
            }
            TableColumn("Title") {
                Text($0.title)
            }
            TableColumn("Duration") {
                Text(String(format: "%d:%02d", $0.duration / 60, $0.duration % 60))
            }
        }
        .onChange(of: selectedTracks) { selectedTracks in
            if selectedTracks.first == nil {
                return
            }
            
            let track = fetchTracks.tracks.first(where: { $0.id == selectedTracks.first! })
            if track == nil {
                return
            }
            
            if (self.audioPlayer != nil) {
                self.audioPlayer?.pause();
            }
            self.audioPlayer = AVPlayer(url: URL(string: track!.music)!)
            self.audioPlayer?.play()
        }
    }
}

struct TracksView_Previews: PreviewProvider {
    static var previews: some View {
        TracksView().frame(width: 800, height: 600)
    }
}
