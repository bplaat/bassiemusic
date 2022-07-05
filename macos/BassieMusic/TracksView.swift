//
//  TracksView.swift
//  BassieMusic
//
//  Created by Bastiaan van der Plaat on 03/07/2022.
//

import SwiftUI

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
        loadPage(page: 1)
    }

    func loadPage(page: Int) {
        let url = URL(string: "http://localhost:8080/api/tracks?page=\(page)")!
        URLSession.shared.dataTask(with: url) {(data, response, error) in
            do {
                let newTracks = try JSONDecoder().decode([Track].self, from: data!)
                DispatchQueue.main.async {
                    if newTracks.count > 0 {
                        self.tracks.append(contentsOf: newTracks)
                        self.loadPage(page: page + 1)
                    }
                }
            } catch {
                print("Error when loading tracks")
            }
        }.resume()
    }
}

struct TracksView: View {
    var playTrack: (_ track: Track) -> Void
    
    @ObservedObject var fetchTracks = FetchTracks()

    @State var selectedTrackId: Track.ID?

    var body: some View {
        Table(fetchTracks.tracks, selection: $selectedTrackId) {
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
        .onChange(of: selectedTrackId) { selectedTrackId in
            let track = fetchTracks.tracks.first(where: { $0.id == selectedTrackId! })
            if track == nil {
                return
            }
            self.playTrack(track!)
        }
    }
}

