<script>
    import { goto } from '$app/navigation';
    import ImageEditButton from '../../../components/buttons/image-edit-button.svelte';
    import LikeButton from '../../../components/buttons/like-button.svelte';
    import EditModal from '../../../components/modals/albums/edit-modal.svelte';
    import DeleteModal from '../../../components/modals/delete-modal.svelte';
    import TracksTable from '../../../components/tracks-table.svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: '$1 - Albums - BassieMusic',
            back: 'Go back one page',
            play: 'Play album',
            edit: 'Edit album',
            album: 'album',
            delete: 'Delete album',
            tracks: 'Tracks',
            tracks_empty: "This album doesn't have any tracks",
        },
        nl: {
            title: '$1 - Albums - BassieMusic',
            back: 'Ga een pagina terug',
            play: 'Speel album',
            edit: 'Verander album',
            album: 'album',
            delete: 'Verwijder album',
            tracks: 'Tracks',
            tracks_empty: 'Dit album heeft geen enkele track',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // State
    export let data;
    let tracksTable;
    let editModal;
    let deleteModal;
</script>

<svelte:head>
    <title>{t('title', data.album.title)}</title>
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
            item={data.album}
            itemRoute="albums"
            editable={data.authUser.role === 'admin'}
        />
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title mb-3">{data.album.title}</h2>
        <p class="mb-3">{data.album.released_at.split('T')[0]}</p>
        {#if data.album.genres.length > 0}
            <p class="mb-3">
                {#each data.album.genres as genre}
                    <a href="/genres/{genre.id}" class="mr-2">{genre.name}</a>
                {/each}
            </p>
        {/if}
        <p class="mb-4">
            {#each data.album.artists as artist}
                <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
            {/each}
        </p>

        <div class="buttons">
            <button class="button is-large" on:click={tracksTable.playFirstTrack} title={t('play')}>
                <svg class="icon" viewBox="0 0 24 24">
                    <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                </svg>
            </button>

            <LikeButton token={data.token} item={data.album} itemRoute="albums" itemLabel={t('album')} isLarge={true} />

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

<h3 class="title is-4">{t('tracks')}</h3>
{#if data.album.tracks.length > 0}
    <TracksTable
        bind:this={tracksTable}
        token={data.token}
        authUser={data.authUser}
        tracks={data.album.tracks}
        isAlbum={true}
    />
{:else}
    <p><i>{t('tracks_empty')}</i></p>
{/if}

{#if data.authUser.role === 'admin'}
    <EditModal
        bind:this={editModal}
        token={data.token}
        album={data.album}
        on:update={(event) => {
            data.album = event.detail.album;
        }}
    />

    <DeleteModal
        bind:this={deleteModal}
        token={data.token}
        item={data.album}
        itemRoute="albums"
        itemLabel={t('album')}
        on:delete={() => {
            goto('/albums');
        }}
    />
{/if}
