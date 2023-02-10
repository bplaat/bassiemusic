<script>
    import { page } from '$app/stores';
    import { browser } from '$app/environment';
    import TracksTable from '../../../components/tracks-table.svelte';
    import AlbumCard from '../../../components/album-card.svelte';

    export let data;
    let { token, artist } = data;

    if (browser) {
        page.subscribe(async (page) => {
            if (page.url.pathname.startsWith('/artists/') && page.url.pathname != `/artists/${artist.id}`) {
                const response = await fetch(`${import.meta.env.VITE_API_URL}/artists/${page.params.id}`, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });
                artist = await response.json();
            }
        });
    }

    let topTracksTable;

    function likeArtist() {
        fetch(`${import.meta.env.VITE_API_URL}/artists/${artist.id}/like${artist.liked ? '/delete' : ''}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        artist.liked = !artist.liked;
    }

    let albumType = 'all';
    $: filteredAlbums = (artist.albums || []).filter((album) => {
        if (albumType == 'all') return true;
        if (albumType == 'album') return album.type == 'album';
        if (albumType == 'ep') return album.type == 'ep';
        if (albumType == 'single') return album.type == 'single';
    });
</script>

<svelte:head>
    <title>{artist.name} - Artists - BassieMusic</title>
</svelte:head>

<div class="buttons">
    <button class="button" on:click={() => history.back()}>
        <svg class="icon" viewBox="0 0 24 24">
            <path d="M20,11V13H8L13.5,18.5L12.08,19.92L4.16,12L12.08,4.08L13.5,5.5L8,11H20Z" />
        </svg>
    </button>
</div>

<div class="columns">
    <div class="column is-one-quarter mr-5">
        <div class="box p-0">
            <img style="aspect-ratio: 1;" src={artist.large_image} alt="Image of artist {artist.name}" loading="lazy" />
        </div>
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title">{artist.name}</h2>

        <div class="buttons">
            <button class="button is-large" on:click={topTracksTable.playFirstTrack}>
                <svg class="icon" viewBox="0 0 24 24">
                    <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                </svg>
            </button>

            <button class="button is-large" on:click={likeArtist}>
                {#if artist.liked}
                    <svg class="icon is-colored" viewBox="0 0 24 24">
                        <path
                            fill="#f14668"
                            d="M12,21.35L10.55,20.03C5.4,15.36 2,12.27 2,8.5C2,5.41 4.42,3 7.5,3C9.24,3 10.91,3.81 12,5.08C13.09,3.81 14.76,3 16.5,3C19.58,3 22,5.41 22,8.5C22,12.27 18.6,15.36 13.45,20.03L12,21.35Z"
                        />
                    </svg>
                {:else}
                    <svg class="icon" viewBox="0 0 24 24">
                        <path
                            d="M12.1,18.55L12,18.65L11.89,18.55C7.14,14.24 4,11.39 4,8.5C4,6.5 5.5,5 7.5,5C9.04,5 10.54,6 11.07,7.36H12.93C13.46,6 14.96,5 16.5,5C18.5,5 20,6.5 20,8.5C20,11.39 16.86,14.24 12.1,18.55M16.5,3C14.76,3 13.09,3.81 12,5.08C10.91,3.81 9.24,3 7.5,3C4.42,3 2,5.41 2,8.5C2,12.27 5.4,15.36 10.55,20.03L12,21.35L13.45,20.03C18.6,15.36 22,12.27 22,8.5C22,5.41 19.58,3 16.5,3Z"
                        />
                    </svg>
                {/if}
            </button>
        </div>
    </div>
</div>

<h2 class="title mt-5">Top Tracks</h2>
{#if artist.top_tracks.length > 0}
    <TracksTable bind:this={topTracksTable} {token} tracks={artist.top_tracks} />
{:else}
    <p><i>This artist doens't have any top tracks</i></p>
{/if}

<h2 class="title mt-5">Albums</h2>
{#if artist.albums != undefined}
    <div class="tabs is-toggle">
        <ul>
            <li class:is-active={albumType == 'all'}>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a href="#" on:click|preventDefault={() => (albumType = 'all')}>All</a>
            </li>
            <li class:is-active={albumType == 'album'}>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a href="#" on:click|preventDefault={() => (albumType = 'album')}>Albums</a>
            </li>
            <li class:is-active={albumType == 'ep'}>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a href="#" on:click|preventDefault={() => (albumType = 'ep')}>EPs</a>
            </li>
            <li class:is-active={albumType == 'single'}>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a href="#" on:click|preventDefault={() => (albumType = 'single')}>Singles</a>
            </li>
        </ul>
    </div>

    {#if filteredAlbums.length > 0}
        <div class="columns is-multiline is-mobile">
            {#each filteredAlbums as album}
                <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                    <AlbumCard {album} />
                </div>
            {/each}
        </div>
    {:else}
        <p><i>This artist has no albums of the selected type</i></p>
    {/if}
{:else}
    <p><i>This artist has no albums</i></p>
{/if}
