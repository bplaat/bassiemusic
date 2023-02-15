<script>
    import { musicPlayer } from '../stores.js';

    export let artist;
    export let token;

    async function fetchAndPlayTracks() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/artists/${artist.id}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        const completeArtist = await response.json();
        $musicPlayer.playTracks(completeArtist.top_tracks, completeArtist.top_tracks[0]);
    }
</script>

<a class="card has-image-play-button" href="/artists/{artist.id}">
    <div class="card-image" style="aspect-ratio: 1;">
        <img src={artist.medium_image} alt="Image of artist {artist.name}" loading="lazy" />
        <button class="button image-play-button" on:click|preventDefault={fetchAndPlayTracks}>
            <svg class="icon" viewBox="0 0 24 24">
                <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
            </svg>
        </button>
    </div>
    <div class="card-content">
        <h3 class="title is-6 ellipsis">{artist.name}</h3>
    </div>
</a>
