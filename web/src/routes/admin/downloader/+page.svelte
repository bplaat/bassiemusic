<script>
    import { formatBytes } from '../../../filters.js';

    export let data;
    const { token, storage } = data;

    let query = '';
    let results = false;
    let albums = [];
    let artists = [];

    async function search() {
        if (query == '') {
            results = false;
            albums = [];
            artists = [];
            return;
        }

        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/deezer_search?${new URLSearchParams({
                q: query,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const data = await response.json();
        results = true;
        albums = data.albums;
        artists = data.artists;
    }

    async function downloadAlbum(album) {
        albums = albums.filter((otherAlbum) => otherAlbum.id != album.id);
        await fetch(
            `${import.meta.env.VITE_API_URL}/download/album?${new URLSearchParams({
                deezer_id: album.id,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
    }

    async function downloadArtist(artist) {
        artists = artists.filter((otherArtist) => otherArtist.id != artist.id);
        await fetch(
            `${import.meta.env.VITE_API_URL}/download/artist?${new URLSearchParams({
                deezer_id: artist.id,
                singles: true,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
    }
</script>

<svelte:head>
    <title>Downloader - Admin - BassieMusic</title>
</svelte:head>

<h1 class="title">Admin Downloader</h1>

<div class="box">
    <h2 class="title is-4">Storage folder size</h2>
    <progress class="progress is-link" value={storage.used} max={storage.max}>
        {((storage.used / storage.max) * 100).toFixed(2)}%
    </progress>
    <p>
        <span class="mr-3">Used: {formatBytes(storage.used)}</span>
        <span>Max: {formatBytes(storage.max)}</span>
    </p>
</div>

<div class="box">
    <h2 class="title is-4">Search and download albums and artists</h2>

    <form on:submit|preventDefault={search} class="field has-addons">
        <div class="control" style="width: 100%;">
            <input class="input" type="text" bind:value={query} placeholder="Find an album or artist..." />
        </div>
        <div class="control">
            <button type="submit" class="button is-link">Search</button>
        </div>
    </form>

    {#if results}
        <div class="columns mt-5">
            <div class="column is-half">
                <h2 class="title is-4">Albums</h2>
                {#each albums as album}
                    <div class="media">
                        <div class="media-left">
                            <div class="box m-0 p-0" style="width: 48px; height: 48px;">
                                <img src={album.cover_medium} alt="Cover of album {album.title}" loading="lazy" />
                            </div>
                        </div>
                        <div class="media-content" style="min-width: 0;">
                            <p class="ellipsis" style="font-weight: 500;">{album.title}</p>
                            <p class="ellipsis">{album.artist.name}</p>
                        </div>
                        <div class="media-right">
                            <button
                                class="button is-link"
                                on:click={() => downloadAlbum(album)}
                                title="Add album to BassieMusic"
                            >
                                <svg class="icon" viewBox="0 0 24 24">
                                    <path d="M19,13H13V19H11V13H5V11H11V5H13V11H19V13Z" />
                                </svg>
                            </button>
                        </div>
                    </div>
                {/each}
                {#if albums.length == 0}
                    <p><i>Can't find any albums on Deezer</i></p>
                {/if}
            </div>

            <div class="column is-half">
                <h2 class="title is-4">Artists</h2>
                {#each artists as artist}
                    <div class="media">
                        <div class="media-left">
                            <div class="box m-0 p-0" style="width: 48px; height: 48px;">
                                <img src={artist.picture_medium} alt="Image of artist {artist.name}" loading="lazy" />
                            </div>
                        </div>
                        <div class="media-content">
                            <p class="ellipsis" style="font-weight: 500;">{artist.name}</p>
                        </div>
                        <button
                            class="button is-link"
                            on:click={() => downloadArtist(artist)}
                            title="Add artist to BassieMusic"
                        >
                            <svg class="icon" viewBox="0 0 24 24">
                                <path d="M19,13H13V19H11V13H5V11H11V5H13V11H19V13Z" />
                            </svg>
                        </button>
                    </div>
                {/each}
                {#if artists.length == 0}
                    <p><i>Can't find any artist on Deezer</i></p>
                {/if}
            </div>
        </div>
    {/if}
</div>
