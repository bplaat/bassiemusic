//
//  ContentView.swift
//  BassieMusic
//
//  Created by Bastiaan van der Plaat on 03/07/2022.
//

import SwiftUI

struct Album : Identifiable, Decodable {
    var id: String
    var title: String
    var cover: String
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

struct MenuView : View {
    var body: some View {
        VStack {
            HStack {
                Text("BassieMusic")
                Spacer()
            }
            
            HStack {
                Image(systemName: "photo")
                Text("Artists")
                Spacer()
            }
            HStack {
                Image(systemName: "photo")
                Text("Albums")
                Spacer()
            }
            HStack {
                Image(systemName: "photo")
                Text("Tracks")
                Spacer()
            }
            Spacer()
        }.padding(16)
    }
}

struct MainView : View {
    @ObservedObject var fetchAlbums = FetchAlbums()

    var layout = [
        GridItem(.adaptive(minimum: 120))
    ]
    
    var body: some View {
        ScrollView {
            LazyVGrid(columns: layout, spacing: 16) {
                ForEach(fetchAlbums.albums, id: \.id) { album in
                    VStack {
                        AsyncImage(url: URL(string: album.cover)) { image in
                            image.resizable().scaledToFit()
                        } placeholder: {
                            Image(systemName: "photo").resizable()
                        }
                        .mask(RoundedRectangle(cornerRadius: 4))
                        .shadow(radius: 8)
                        Text(album.title)
                    }
                }
            }.padding(16)
        }
    }
}

struct ContentView: View {
    var body: some View {
        NavigationView {
            MenuView()
            MainView()
        }
        .frame(minWidth: 320, idealWidth: 1280, minHeight: 240, idealHeight: 720)
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
            .frame(height: 600.0)
    }
}
