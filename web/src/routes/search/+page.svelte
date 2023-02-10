<script>
    import { onMount, onDestroy } from 'svelte';
    import GenreCard from '../../components/genre-card.svelte';
    import AlbumCard from '../../components/album-card.svelte';
    import ArtistCard from '../../components/artist-card.svelte';
    import TracksTable from '../../components/tracks-table.svelte';

    export let data;
    let { token, genres: allGenres } = data;

    // Lazy load all genres
    async function fetchGenresPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/genres?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newGenres } = await response.json();
        allGenres.push(...newGenres);
        allGenres = allGenres;
    }

    let bottom;
    if (allGenres.length != data.total) {
        let observer;
        onMount(() => {
            let page = 2;
            observer = new IntersectionObserver(
                (entries, observer) => {
                    for (const entry of entries) {
                        if (allGenres.length >= data.total) {
                            observer.unobserve(entry.target);
                        } else {
                            fetchGenresPage(page++);
                        }
                    }
                },
                {
                    root: document.body,
                }
            );
            observer.observe(bottom);
        });
        onDestroy(() => {
            if (observer) observer.unobserve(bottom);
        });
    }

    // Search results
    let searchTerm = '';
    let hasResult = false;
    let genres = [];
    let albums = [];
    let artists = [];
    let tracks = [];

    // Perform search
    async function search() {
        if (searchTerm != '') {
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

            // Show response
            const data = await response.json();
            hasResult = true;
            artists = data.artists;
            albums = data.albums;
            tracks = data.tracks;
            genres = data.genres;
        } else {
            genres = allGenres;
            hasResult = false;
            albums = [];
            artists = [];
            tracks = [];
        }
    }
</script>

<svelte:head>
    <title>Search - BassieMusic</title>
</svelte:head>

<h2 class="title">Search</h2>

<form on:submit|preventDefault={search} class="field has-addons mb-5">
    <div class="control" style="width: 100%;">
        <input
            class="input"
            type="text"
            bind:value={searchTerm}
            placeholder="Find a track, album, artist or genre..."
        />
    </div>
    <div class="control">
        <button type="submit" class="button is-link">Search</button>
    </div>
</form>

{#if hasResult}
    {#if tracks.length > 0}
        <h2 class="title is-5">Tracks</h2>
        <TracksTable {token} {tracks} />
    {/if}

    {#if artists.length > 0}
        <h2 class="title is-5 mt-5">Artists</h2>
        <div class="columns is-multiline mb-5">
            {#each artists as artist}
                <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                    <ArtistCard {artist} />
                </div>
            {/each}
        </div>
    {/if}

    {#if albums.length > 0}
        <h2 class="title is-5">Albums</h2>
        <div class="columns is-multiline mb-5">
            {#each albums as album}
                <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                    <AlbumCard {album} />
                </div>
            {/each}
        </div>
    {/if}

    {#if genres.length > 0}
        <h2 class="title is-5">Genres</h2>
        <div class="columns is-multiline is-mobile">
            {#each genres as genre}
                <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                    <GenreCard {genre} />
                </div>
            {/each}
        </div>
    {/if}
{:else}
    <h2 class="title is-5">Genres</h2>
    <div class="columns is-multiline is-mobile">
        {#each allGenres as genre}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <GenreCard {genre} />
            </div>
        {/each}
    </div>
{/if}

<div bind:this={bottom} />
