<script>
    import { page } from '$app/stores';
    import { browser } from '$app/environment';
    import TracksTable from '../../../components/tracks-table.svelte';

    export let data;
    let { token, album } = data;
    album.tracks = album.tracks.slice().map((track) => {
        track.album = album;
        return track;
    });

    if (browser) {
        page.subscribe(async (page) => {
            if (page.url.pathname.startsWith('/albums/') && page.url.pathname != `/albums/${album.id}`) {
                const response = await fetch(`${import.meta.env.VITE_API_URL}/albums/${page.params.id}`, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });
                album = await response.json();
                album.tracks = album.tracks.slice().map((track) => {
                    track.album = album;
                    return track;
                });
            }
        });
    }

    let tracksTable;

    function likeAlbum() {
        fetch(`${import.meta.env.VITE_API_URL}/albums/${album.id}/like${album.liked ? '/delete' : ''}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        album.liked = !album.liked;
    }
</script>

<svelte:head>
    <title>{album.title} - Albums - BassieMusic</title>
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
        <div class="box p-0 has-image-tags">
            <img style="aspect-ratio: 1;" src={album.large_cover} alt="Cover of album {album.title}" loading="lazy" />

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
                    <span class="tag is-danger">E</span>
                {/if}
            </div>
        </div>
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title mb-3">{album.title}</h2>
        <p class="mb-3">{album.released_at.split('T')[0]}</p>
        {#if album.genres != undefined}
            <p class="mb-3">
                {#each album.genres as genre}
                    <a href="/genres/{genre.id}" class="mr-2">{genre.name}</a>
                {/each}
            </p>
        {/if}
        <p class="mb-4">
            {#each album.artists as artist}
                <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
            {/each}
        </p>

        <div class="buttons">
            <button class="button is-large" on:click={tracksTable.playFirstTrack}>
                <svg class="icon" viewBox="0 0 24 24">
                    <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                </svg>
            </button>

            <button class="button is-large" on:click={likeAlbum}>
                {#if album.liked}
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

<h3 class="title is-4">Tracks</h3>
{#if album.tracks.length > 0}
    <TracksTable bind:this={tracksTable} {token} tracks={album.tracks} showAlbum={false} />
{:else}
    <p><i>This album doens't have any tracks</i></p>
{/if}
