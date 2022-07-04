using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Linq;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;

namespace BassieMusic
{
    public partial class MainWindow : Window
    {
        private static readonly HttpClient client = new HttpClient();

        private Label OutputLabel;

        private static async Task<List<Album>?> fetchAlbums()
        {
            var albumsJson = await client.GetStreamAsync("http://localhost:8080/api/albums");
            return await JsonSerializer.DeserializeAsync<List<Album>>(albumsJson);
        }

        private static async Task<Album?> fetchAlbum(string id)
        {
            var albumJson = await client.GetStreamAsync($"http://localhost:8080/api/albums/{id}");
            return await JsonSerializer.DeserializeAsync<Album>(albumJson);
        }

        async Task FetchStuff()
        {
            var albums = await fetchAlbums();
            foreach (var otherAlbum in albums)
            {
                var album = await fetchAlbum(otherAlbum.ID);

                OutputLabel.Content += $"\n# {album.Title} by {String.Join(", ", album.Artists!.Select(artist => artist.Name))}\n";
                foreach (var track in album.Tracks)
                {
                    OutputLabel.Content += $"{track.Position}. {track.Title} by {String.Join(", ", track.Artists!.Select(artist => artist.Name))}\n";
                }
            }
        }

        public MainWindow()
        {
            InitializeComponent();
            OutputLabel = (Label)FindName("Output");
            FetchStuff();
        }
    }
}
