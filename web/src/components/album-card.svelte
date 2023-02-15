<script>
    import { musicPlayer } from '../stores.js';

    export let album;
    export let token;

    async function fetchAndPlayTracks() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/albums/${album.id}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        const completeAlbum = await response.json();
        const tracks = completeAlbum.tracks.map((track) => {
            track.album = album;
            return track;
        });
        $musicPlayer.playTracks(tracks, tracks[0]);
    }
</script>

<a class="card has-image-play-button" href="/albums/{album.id}">
    <div class="card-image has-image-tags" style="aspect-ratio: 1;">
        <img src={album.medium_cover} alt="Cover of album {album.name}" loading="lazy" />
        <div class="image-tags">
            {#if album.type == 'album'}
                <span class="tag">ALBUM</span>
            {/if}
            {#if album.type == 'ep'}
                <span class="tag">EP</span>
            {/if}
            {#if album.type == 'single'}
                <span class="tag">SINGLE</span>
            {/if}
            {#if album.explicit}
                <span class="tag is-danger" title="Explicit lyrics">E</span>
            {/if}
        </div>
        <button class="button image-play-button" on:click|preventDefault={fetchAndPlayTracks}>
            <svg class="icon" viewBox="0 0 24 24">
                <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
            </svg>
        </button>
    </div>
    <div class="card-content">
        <h3 class="title is-6 mb-2 ellipsis">{album.title}</h3>
        <p class="ellipsis">
            {#each album.artists as artist}
                <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
            {/each}
        </p>
    </div>
</a>
