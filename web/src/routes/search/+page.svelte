<script>
    import GenreCard from "../../components/genre-card.svelte";
    import AlbumCard from "../../components/album-card.svelte";
    import ArtistCard from "../../components/artist-card.svelte";
    import TracksTable from "../../components/tracks-table.svelte";

    export let data;
    let token = data.token;
    let allGenres = data.genres;

    let genres = allGenres;
    let albums = [];
    let artists = [];
    let tracks = [];

    let searchTerm = "";
    let showGenres = true;
    let showAlbums = false;
    let showArtists = false;
    let showTracks = false;

    // Perform search
    async function search() {
        if (searchTerm != "") {
            // Clear genres
            showGenres = false;

            // Load data from database
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/search?${new URLSearchParams({
                    q: searchTerm,
                })}`,
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );

            const data = await response.json();

            artists = data.artists;
            showArtists = artists.length != 0;

            albums = data.albums;
            showAlbums = albums.length != 0;

            tracks = data.tracks;
            showTracks = tracks.length != 0;

            let tempGenres = data.genres;
            showGenres = tempGenres.length != 0;
            if (showGenres) genres = tempGenres;
        } else {
            showGenres = true;
            showAlbums = false;
            showArtists = false;
            showTracks = false;

            genres = allGenres;

            albums = [];
            artists = [];
            tracks = [];
        }
    }

    // Handles when press enter in search field
    function handleKeyDown(event) {
        if (event.key === "Enter") {
            search();
        }
    }
</script>

<svelte:head>
    <title>Search - BassieMusic</title>
</svelte:head>

<h2 class="title">Search</h2>

<div class="field has-addons mb-5">
    <div class="control">
        <input
            class="input"
            type="text"
            bind:value={searchTerm}
            on:keydown={handleKeyDown}
            placeholder="Find a track, album, artist or genre"
        />
    </div>
    <div class="control">
        <button class="button is-link" on:click={search}> Search </button>
    </div>
</div>

{#if showTracks}
    <h2 class="title is-5">Tracks</h2>

    <TracksTable {token} {tracks} />
{/if}

{#if showArtists}
    <h2 class="title is-5 mt-5">Artists</h2>

    <div class="columns is-multiline mb-5">
        {#each artists as artist}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <ArtistCard {artist} />
            </div>
        {/each}
    </div>
{/if}

{#if showAlbums}
    <h2 class="title is-5">Albums</h2>

    <div class="columns is-multiline mb-5">
        {#each albums as album}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <AlbumCard {album} />
            </div>
        {/each}
    </div>
{/if}

{#if showGenres}
    <h2 class="title is-5">Genres</h2>

    <div class="columns is-multiline is-mobile">
        {#each genres as genre}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <GenreCard {genre} />
            </div>
        {/each}
    </div>
{/if}
