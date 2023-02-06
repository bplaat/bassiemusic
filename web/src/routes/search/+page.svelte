<svelte:head>
    <title>Search - BassieMusic</title>
</svelte:head>

<script>
    import Cookies from "js-cookie";
    import GenreCard from "../../components/genre-card.svelte";
    import AlbumCard from "../../components/album-card.svelte";
    import ArtistCard from "../../components/artist-card.svelte";
    import TracksTable from "../../components/tracks-table.svelte";

    export let data;
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
                        Authorization: `Bearer ${Cookies.get("token")}`,
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
            if(showGenres) genres = tempGenres

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

<h2 class="title">Search</h2>

<div class="field has-addons">
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
        <button class="button is-info" on:click={search}> Search </button>
    </div>
</div>

<br>

{#if showTracks}
    <div>
        <h2 class="title is-size-5">Tracks</h2>

        <TracksTable {tracks} />
    </div>
    <br>
{/if}

{#if showArtists}
    <div>
        <h2 class="title is-size-5">Artists</h2>

        <div class="columns is-multiline">
            {#each artists as artist}
                <div class="column is-one-fifth">
                    <ArtistCard {artist} />
                </div>
            {/each}
        </div>        
    </div>
    <br>
{/if}

{#if showAlbums}
    <div>
        <h2 class="title is-size-5">Albums</h2>

        <div class="columns is-multiline">
            {#each albums as album}
                <div class="column is-one-fifth">
                    <AlbumCard {album} />
                </div>
            {/each}
        </div>
    </div>
    <br>
{/if}

{#if showGenres}
    <div>
        <h2 class="title is-size-5">Genres</h2>

        <div class="columns is-multiline">
            {#each genres as genre}
                <div class="column is-one-fifth">
                    <GenreCard {genre} />
                </div>
            {/each}
        </div>
    </div>
{/if}
