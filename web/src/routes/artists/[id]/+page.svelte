<script>
    import { goto } from '$app/navigation';
    import TracksTable from '../../../components/tracks-table.svelte';
    import LikeButton from '../../../components/like-button.svelte';
    import DeleteModal from '../../../components/modals/delete-modal.svelte';
    import AlbumCard from '../../../components/cards/album-card.svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: '$1 - Artists - BassieMusic',
            back: 'Go back one page',
            image_alt: 'Image of artist $1',
            sync: 'This aritist is synced, we will download automatic new albums',
            play: 'Play artist top tracks',
            artist: 'artist',
            delete: 'Delete artist',
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
            sync: 'Deze artiest is gesynchroniseerd, we zullen automatisch nieuwe albums downloaden',
            play: 'Speel artist top tracks',
            artist: 'artist',
            delete: 'Verwijder artist',
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
    let deleteModal;
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
        <div class="box has-image has-image-tags p-0">
            <figure class="image is-1by1">
                <img
                    src={data.artist.large_image || '/images/avatar-default.svg'}
                    alt={t('image_alt', data.artist.name)}
                />
            </figure>
            <div class="image-tags">
                {#if data.artist.sync}
                    <span class="tag px-2 py-1" style="height: auto;" title={t('sync')}>
                        <svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                            <path
                                d="M12,18A6,6 0 0,1 6,12C6,11 6.25,10.03 6.7,9.2L5.24,7.74C4.46,8.97 4,10.43 4,12A8,8 0 0,0 12,20V23L16,19L12,15M12,4V1L8,5L12,9V6A6,6 0 0,1 18,12C18,13 17.75,13.97 17.3,14.8L18.76,16.26C19.54,15.03 20,13.57 20,12A8,8 0 0,0 12,4Z"
                            />
                        </svg>
                    </span>
                {/if}
            </div>
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

            <LikeButton token={data.token} item={data.artist} itemRoute="artists" itemLabel={t('artist')} isLarge={true} />

            {#if data.authUser.role == 'admin'}
                <button class="button is-large" on:click={() => deleteModal.open()} title={t('delete')}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path d="M19,4H15.5L14.5,3H9.5L8.5,4H5V6H19M6,19A2,2 0 0,0 8,21H16A2,2 0 0,0 18,19V7H6V19Z" />
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
            <li class:is-active={data.filterAlbumsBy == 'all'}>
                <a href="?albums_filter=all">{t('album_type_all')}</a>
            </li>
            <li class:is-active={data.filterAlbumsBy == 'album'}>
                <a href="?albums_filter=album">{t('album_type_album')}</a>
            </li>
            <li class:is-active={data.filterAlbumsBy == 'ep'}>
                <a href="?albums_filter=ep">{t('album_type_ep')}</a>
            </li>
            <li class:is-active={data.filterAlbumsBy == 'single'}>
                <a href="?albums_filter=single">{t('album_type_single')}</a>
            </li>
        </ul>
    </div>

    {#if data.filteredAlbums.length > 0}
        <div class="columns is-multiline is-mobile">
            {#each data.filteredAlbums as album}
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

{#if data.authUser.role == 'admin'}
    <DeleteModal
        bind:this={deleteModal}
        token={data.token}
        item={data.artist}
        itemRoute="artists"
        itemLabel={t('artist')}
        on:delete={() => {
            goto('/artists');
        }}
    />
{/if}
