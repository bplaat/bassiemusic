<script>
    import TracksTable from '../../../components/tracks-table.svelte';
    import AlbumCard from '../../../components/cards/album-card.svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: '$1 - Artists - BassieMusic',
            back: 'Go back one page',
            image_alt: 'Image of artist $1',
            play: 'Play artist top tracks',
            like: 'Like artist',
            remove_like: 'Remove artist like',
            top_tracks: 'Top Tracks',
            top_tracks_empty: "This artist doesn't have any top tracks",
            albums: 'Albums',
            album_type_all: 'All',
            album_type_album: 'Albums',
            album_type_ep: 'EPs',
            album_type_single: 'Singles',
            albums_type_empty: 'This artist has no albums of the selected type',
            albums_empty: 'This artist has no albums',
        },
        nl: {
            title: '$1 - Artisten - BassieMusic',
            back: 'Ga een pagina terug',
            image_alt: 'Afbeelding van artist $1',
            play: 'Speel artist top tracks',
            like: 'Like artist',
            remove_like: 'Verwijder artist like',
            top_tracks: 'Top Tracks',
            top_tracks_empty: 'Deze artiest heeft geen topnummers',
            albums: 'Albums',
            album_type_all: 'Alles',
            album_type_album: 'Albums',
            album_type_ep: 'EPs',
            album_type_single: 'Singles',
            albums_type_empty: 'Deze artiest heeft geen albums van het geselecteerde type',
            albums_empty: 'Deze artiest heeft geen albums',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // State
    export let data;
    let topTracksTable;
    let albumType = 'all';
    $: filteredAlbums = data.artist.albums.filter((album) => {
        if (albumType == 'all') return true;
        if (albumType == 'album') return album.type == 'album';
        if (albumType == 'ep') return album.type == 'ep';
        if (albumType == 'single') return album.type == 'single';
    });

    // Methods
    function likeArtist() {
        fetch(`${import.meta.env.VITE_API_URL}/artists/${data.artist.id}/like`, {
            method: data.artist.liked ? 'DELETE' : 'PUT',
            headers: {
                Authorization: `Bearer ${data.token}`,
            },
        });
        data.artist.liked = !data.artist.liked;
    }
</script>

<svelte:head>
    <title>{t('title', data.artist.name)}</title>
</svelte:head>

<div class="buttons">
    <button class="button" on:click={() => history.back()} title={t('back')}>
        <svg class="icon" viewBox="0 0 24 24">
            <path d="M20,11V13H8L13.5,18.5L12.08,19.92L4.16,12L12.08,4.08L13.5,5.5L8,11H20Z" />
        </svg>
    </button>
</div>

<div class="columns">
    <div class="column is-one-quarter mr-5 mr-0-mobile">
        <div class="box has-image p-0" style="aspect-ratio: 1;">
            <img src={data.artist.large_image} alt={t('image_alt', data.artist.name)} />
        </div>
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title">{data.artist.name}</h2>

        <div class="buttons">
            <button class="button is-large" on:click={topTracksTable.playFirstTrack} title={t('play')}>
                <svg class="icon" viewBox="0 0 24 24">
                    <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                </svg>
            </button>

            {#if !data.artist.liked}
                <button class="button is-large" on:click={likeArtist} title={t('like')}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path
                            d="M12.1,18.55L12,18.65L11.89,18.55C7.14,14.24 4,11.39 4,8.5C4,6.5 5.5,5 7.5,5C9.04,5 10.54,6 11.07,7.36H12.93C13.46,6 14.96,5 16.5,5C18.5,5 20,6.5 20,8.5C20,11.39 16.86,14.24 12.1,18.55M16.5,3C14.76,3 13.09,3.81 12,5.08C10.91,3.81 9.24,3 7.5,3C4.42,3 2,5.41 2,8.5C2,12.27 5.4,15.36 10.55,20.03L12,21.35L13.45,20.03C18.6,15.36 22,12.27 22,8.5C22,5.41 19.58,3 16.5,3Z"
                        />
                    </svg>
                </button>
            {:else}
                <button class="button is-large" on:click={likeArtist} title={t('remove_like')}>
                    <svg class="icon is-colored" viewBox="0 0 24 24">
                        <path
                            fill="#f14668"
                            d="M12,21.35L10.55,20.03C5.4,15.36 2,12.27 2,8.5C2,5.41 4.42,3 7.5,3C9.24,3 10.91,3.81 12,5.08C13.09,3.81 14.76,3 16.5,3C19.58,3 22,5.41 22,8.5C22,12.27 18.6,15.36 13.45,20.03L12,21.35Z"
                        />
                    </svg>
                </button>
            {/if}
        </div>
    </div>
</div>

<h2 class="title mt-5">{t('top_tracks')}</h2>
{#if data.artist.top_tracks.length > 0}
    <TracksTable
        bind:this={topTracksTable}
        token={data.token}
        authUser={data.authUser}
        tracks={data.artist.top_tracks}
    />
{:else}
    <p><i>{t('top_tracks_empty')}</i></p>
{/if}

<h2 class="title mt-5">{t('albums')}</h2>
{#if data.artist.albums.length > 0}
    <div class="tabs is-toggle">
        <ul>
            <li class:is-active={albumType == 'all'}>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a href="#" on:click|preventDefault={() => (albumType = 'all')}>{t('album_type_all')}</a>
            </li>
            <li class:is-active={albumType == 'album'}>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a href="#" on:click|preventDefault={() => (albumType = 'album')}>{t('album_type_album')}</a>
            </li>
            <li class:is-active={albumType == 'ep'}>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a href="#" on:click|preventDefault={() => (albumType = 'ep')}>{t('album_type_ep')}</a>
            </li>
            <li class:is-active={albumType == 'single'}>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a href="#" on:click|preventDefault={() => (albumType = 'single')}>{t('album_type_single')}</a>
            </li>
        </ul>
    </div>

    {#if filteredAlbums.length > 0}
        <div class="columns is-multiline is-mobile">
            {#each filteredAlbums as album}
                <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                    <AlbumCard {album} token={data.token} authUser={data.authUser} />
                </div>
            {/each}
        </div>
    {:else}
        <p><i>{t('albums_type_empty')}</i></p>
    {/if}
{:else}
    <p><i>{t('albums_empty')}</i></p>
{/if}
