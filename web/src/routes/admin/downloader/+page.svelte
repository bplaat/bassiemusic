<script>
    export let data;
    let { token } = data;

    // Artists
    let artistQuery = "";
    let artists = [];

    async function searchArtists() {
        const response = await fetch(
            `${
                import.meta.env.VITE_API_URL
            }/deezer/artists?${new URLSearchParams({
                q: artistQuery,
            })}`
        );
        const { data } = await response.json();
        artists = data;
    }

    async function downloadArtist(artist) {
        await fetch(
            `${
                import.meta.env.VITE_API_URL
            }/download/artist?${new URLSearchParams({
                deezer_id: artist.id,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
    }

    // Albums
    let albumQuery = "";
    let albums = [];

    async function searchAlbums() {
        const response = await fetch(
            `${
                import.meta.env.VITE_API_URL
            }/deezer/albums?${new URLSearchParams({
                q: albumQuery,
            })}`
        );
        const { data } = await response.json();
        albums = data;
    }

    async function downloadAlbum(album) {
        await fetch(
            `${
                import.meta.env.VITE_API_URL
            }/download/album?${new URLSearchParams({
                deezer_id: album.id,
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

<div class="columns">
    <!-- Artists -->
    <div class="column">
        <div class="box">
            <h2 class="title is-4">Download Artist</h2>

            <form
                on:submit|preventDefault={searchArtists}
                class="field has-addons"
            >
                <div class="control">
                    <input
                        class="input"
                        type="text"
                        bind:value={artistQuery}
                        placeholder="Search an artist..."
                    />
                </div>
                <div class="control">
                    <button type="submit" class="button is-link">
                        Search
                    </button>
                </div>
            </form>

            {#each artists as artist}
                <div style="display: flex;">
                    <img
                        src={artist.picture_medium}
                        alt="Image of artist {artist.name}"
                        class="mr-3"
                        style="width: 64px; height: 64px;"
                    />
                    <div style="flex: 1">
                        <b>{artist.name}</b>
                    </div>
                    <button
                        class="button is-link"
                        on:click={() => downloadArtist(artist)}
                        >Add to BassieMusic</button
                    >
                </div>
                <hr />
            {/each}
        </div>
    </div>

    <!-- Albums -->
    <div class="column">
        <div class="box">
            <h2 class="title is-4">Download Album</h2>

            <form
                on:submit|preventDefault={searchAlbums}
                class="field has-addons"
            >
                <div class="control">
                    <input
                        class="input"
                        type="text"
                        bind:value={albumQuery}
                        placeholder="Search an album..."
                    />
                </div>
                <div class="control">
                    <button type="submit" class="button is-link">
                        Search
                    </button>
                </div>
            </form>

            {#each albums as album}
                <div style="display: flex;">
                    <img
                        src={album.cover_medium}
                        alt="Cover of album {album.title}"
                        class="mr-3"
                        style="width: 64px; height: 64px;"
                    />
                    <div style="flex: 1">
                        <b>{album.title}</b><br />
                        By {album.artist.name}
                    </div>
                    <button
                        class="button is-link"
                        on:click={() => downloadAlbum(album)}
                        >Add to BassieMusic</button
                    >
                </div>
                <hr />
            {/each}
        </div>
    </div>
</div>
