<script>
    import { goto } from '$app/navigation';
    import TracksTable from '../../../components/tracks-table.svelte';
    import ImageEditButton from '../../../components/buttons/image-edit-button.svelte';
    import LikeButton from '../../../components/buttons/like-button.svelte';
    import EditModal from '../../../components/modals/artists/edit-modal.svelte';
    import DeleteModal from '../../../components/modals/delete-modal.svelte';
    import AlbumCard from '../../../components/cards/album-card.svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: '$1 - Artists - BassieMusic',
            back: 'Go back one page',
            play: 'Play artist top tracks',
            edit: 'Edit artist',
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
            title: '$1 - Artiesten - BassieMusic',
            back: 'Ga een pagina terug',
            play: 'Speel artiest top tracks',
            edit: 'Verander artiest',
            artist: 'artiest',
            delete: 'Verwijder artiest',
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
    let editModal;
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
        <ImageEditButton
            token={data.token}
            item={data.artist}
            itemRoute="artists"
            editable={data.authUser.role === 'admin'}
        />
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title">{data.artist.name}</h2>

        <div class="buttons">
            <button class="button is-large" on:click={topTracksTable.playFirstTrack} title={t('play')}>
                <svg class="icon" viewBox="0 0 24 24">
                    <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                </svg>
            </button>

            <LikeButton
                token={data.token}
                item={data.artist}
                itemRoute="artists"
                itemLabel={t('artist')}
                isLarge={true}
            />

            {#if data.authUser.role === 'admin'}
                <button class="button is-large" on:click={() => editModal.open()} title={t('edit')}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path
                            d="M20.71,7.04C21.1,6.65 21.1,6 20.71,5.63L18.37,3.29C18,2.9 17.35,2.9 16.96,3.29L15.12,5.12L18.87,8.87M3,17.25V21H6.75L17.81,9.93L14.06,6.18L3,17.25Z"
                        />
                    </svg>
                </button>

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
            <li class:is-active={data.filterAlbumsBy === 'all'}>
                <a href="?albums_filter=all">{t('album_type_all')}</a>
            </li>
            <li class:is-active={data.filterAlbumsBy === 'album'}>
                <a href="?albums_filter=album">{t('album_type_album')}</a>
            </li>
            <li class:is-active={data.filterAlbumsBy === 'ep'}>
                <a href="?albums_filter=ep">{t('album_type_ep')}</a>
            </li>
            <li class:is-active={data.filterAlbumsBy === 'single'}>
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

{#if data.authUser.role === 'admin'}
    <EditModal
        bind:this={editModal}
        token={data.token}
        artist={data.artist}
        on:update={(event) => {
            data.artist = event.detail.artist;
        }}
    />

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
