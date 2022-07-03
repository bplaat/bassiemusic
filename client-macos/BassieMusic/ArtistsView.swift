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
}

class FetchArtists: ObservableObject {
    @Published var artists = [Artist]()

    init() {
        let url = URL(string: "http://localhost:8080/api/artists")!
        URLSession.shared.dataTask(with: url) {(data, response, error) in
            do {
                if let todoData = data {
                    let decodedData = try JSONDecoder().decode([Artist].self, from: todoData)
                    DispatchQueue.main.async {
                        self.artists = decodedData
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
                    }
                }
            }.padding(16)
        }
    }
}

struct ArtistsView_Previews: PreviewProvider {
    static var previews: some View {
        ArtistsView().frame(width: 800, height: 600)
    }
}
