//
//  ArtistsView.swift
//  BassieMusic
//
//  Created by Bastiaan van der Plaat on 03/07/2022.
//

import SwiftUI

struct Artist : Identifiable, Decodable {
    var id: String
    var name: String
    var image: String
    var albums: [Album]?
    var created_at: String
    var updated_at: String
}

class FetchArtists: ObservableObject {
    @Published var artists = [Artist]()
    

    init() {
        loadPage(page: 1)
    }
    
    func loadPage(page: Int) {
        let url = URL(string: "http://localhost:8080/api/artists?page=\(page)")!
        URLSession.shared.dataTask(with: url) {(data, response, error) in
            do {
                let newArtists = try JSONDecoder().decode([Artist].self, from: data!)
                DispatchQueue.main.async {
                    if newArtists.count > 0 {
                        self.artists.append(contentsOf: newArtists)
                        self.loadPage(page: page + 1)
                    }
                }
            } catch {
                print("Error when loading artists")
            }
        }.resume()
    }
}

struct ArtistsView: View {
    @ObservedObject var fetchArtists = FetchArtists()

    var layout = [
        GridItem(.adaptive(minimum: 160))
    ]
    
    var body: some View {
        ScrollView {
            LazyVGrid(columns: layout, spacing: 16) {
                ForEach(fetchArtists.artists, id: \.id) { artist in
                    VStack {
                        AsyncImage(url: URL(string: artist.image)) { image in
                            image.resizable().scaledToFit()
                        } placeholder: {
                            Color.accentColor.opacity(0.1)
                        }
                        .aspectRatio(contentMode: .fit)
                        .mask(RoundedRectangle(cornerRadius: 6))
                        .shadow(radius: 8)
                        
                        Text(artist.name)
                            .bold()
                        
                        Text("Albums: \(artist.albums != nil ? artist.albums!.count : 0)")
                    }
                }
            }.padding(16)
        }
    }
}
