//
//  AlbumsView.swift
//  BassieMusic
//
//  Created by Bastiaan van der Plaat on 03/07/2022.
//

import SwiftUI

struct Album : Identifiable, Decodable {
    var id: String
    var title: String
    var released_at: String
    var cover: String
    var artists: [Artist]?
    var created_at: String
    var updated_at: String
}

class FetchAlbums: ObservableObject {
    @Published var albums = [Album]()

    init() {
        let url = URL(string: "http://localhost:8080/api/albums")!
        URLSession.shared.dataTask(with: url) {(data, response, error) in
            do {
                if let todoData = data {
                    let decodedData = try JSONDecoder().decode([Album].self, from: todoData)
                    DispatchQueue.main.async {
                        self.albums = decodedData
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

struct AlbumsView: View {
    @ObservedObject var fetchAlbums = FetchAlbums()

    var layout = [
        GridItem(.adaptive(minimum: 160))
    ]
    
    var body: some View {
        ScrollView {
            LazyVGrid(columns: layout, spacing: 16) {
                ForEach(fetchAlbums.albums, id: \.id) { album in
                    VStack(alignment: .leading) {
                        AsyncImage(url: URL(string: album.cover)) { image in
                            image.resizable().scaledToFit()
                        } placeholder: {
                            Color.accentColor.opacity(0.1)
                        }
                        .aspectRatio(contentMode: .fit)
                        .mask(RoundedRectangle(cornerRadius: 6))
                        .shadow(radius: 8)
                        
                        Text(album.title)
                            .bold()
                        
                        Text(album.artists!.map(\.name).joined(separator: ", "))
                    }
                }
            }.padding(16)
        }
    }
}

struct AlbumsView_Previews: PreviewProvider {
    static var previews: some View {
        AlbumsView().frame(width: 800, height: 600)
    }
}
