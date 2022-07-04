using System.Text.Json;
using System.Text.Json.Serialization;

namespace BassieMusic
{
    public class Artist
    {
        [JsonPropertyName("id")]
        public string ID { get; set; } = string.Empty;
        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;
    }

    public class Album
    {
        [JsonPropertyName("id")]
        public string ID { get; set; } = string.Empty;
        [JsonPropertyName("title")]
        public string Title { get; set; } = string.Empty;
        [JsonPropertyName("artists")]
        public Artist[]? Artists { get; set; }
        [JsonPropertyName("tracks")]
        public Track[]? Tracks { get; set; }
    }

    public class Track
    {
        [JsonPropertyName("id")]
        public string ID { get; set; } = string.Empty;
        [JsonPropertyName("title")]
        public string Title { get; set; } = string.Empty;
        [JsonPropertyName("disk")]
        public int Disk { get; set; } = 0;
        [JsonPropertyName("position")]
        public int Position { get; set; } = 0;
        [JsonPropertyName("duration")]
        public int Duration { get; set; } = 0;
        [JsonPropertyName("artists")]
        public Artist[]? Artists { get; set; }
    }
}
