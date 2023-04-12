<script>
    import { goto } from '$app/navigation';
    import LikeButton from '../../../components/like-button.svelte';
    import EditModal from '../../../components/modals/playlists/edit-modal.svelte';
    import DeleteModal from '../../../components/modals/delete-modal.svelte';
    import TracksTable from '../../../components/tracks-table.svelte';
    import { sidebar, language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: '$1 - Playlists - BassieMusic',
            back: 'Go back one page',
            image_alt: 'Image of playlist $1',
            public: 'Public playlist',
            made_by: 'Made by $1',
            play: 'Play playlist tracks',
            playlist: 'playlist',
            edit: 'Edit playlist',
            delete: 'Delete playlist',
            tracks: 'Tracks',
            tracks_empty: "This playlist doesn't have any tracks",
        },
        nl: {
            title: '$1 - Afspeellijsten - BassieMusic',
            back: 'Ga een pagina terug',
            image_alt: 'Afbeelding van afspeellijst $1',
            public: 'Publieke afspeellijst',
            made_by: 'Gemaakt door $1',
            play: 'Speel afspeellijst tracks',
            playlist: 'afspeellijst',
            edit: 'Verander afspeellijst',
            delete: 'Verwijder afspeellijst',
            tracks: 'Tracks',
            tracks_empty: 'Deze afspeellijst heeft geen enkele track',
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
    <title>{t('title', data.playlist.name)}</title>
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
                    src={data.playlist.medium_image || '/images/album-default.svg'}
                    alt={t('image_alt', data.playlist.name)}
                />
            </figure>
            <div class="image-tags">
                {#if data.playlist.public}
                    <span class="tag px-2 py-1" style="height: auto;" title={t('public')}>
                        <svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                            <path
                                d="M16.36,14C16.44,13.34 16.5,12.68 16.5,12C16.5,11.32 16.44,10.66 16.36,10H19.74C19.9,10.64 20,11.31 20,12C20,12.69 19.9,13.36 19.74,14M14.59,19.56C15.19,18.45 15.65,17.25 15.97,16H18.92C17.96,17.65 16.43,18.93 14.59,19.56M14.34,14H9.66C9.56,13.34 9.5,12.68 9.5,12C9.5,11.32 9.56,10.65 9.66,10H14.34C14.43,10.65 14.5,11.32 14.5,12C14.5,12.68 14.43,13.34 14.34,14M12,19.96C11.17,18.76 10.5,17.43 10.09,16H13.91C13.5,17.43 12.83,18.76 12,19.96M8,8H5.08C6.03,6.34 7.57,5.06 9.4,4.44C8.8,5.55 8.35,6.75 8,8M5.08,16H8C8.35,17.25 8.8,18.45 9.4,19.56C7.57,18.93 6.03,17.65 5.08,16M4.26,14C4.1,13.36 4,12.69 4,12C4,11.31 4.1,10.64 4.26,10H7.64C7.56,10.66 7.5,11.32 7.5,12C7.5,12.68 7.56,13.34 7.64,14M12,4.03C12.83,5.23 13.5,6.57 13.91,8H10.09C10.5,6.57 11.17,5.23 12,4.03M18.92,8H15.97C15.65,6.75 15.19,5.55 14.59,4.44C16.43,5.07 17.96,6.34 18.92,8M12,2C6.47,2 2,6.5 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z"
                            />
                        </svg>
                    </span>
                {/if}
            </div>
        </div>
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title">{data.playlist.name}</h2>
        <p class="mb-5">{t('made_by', data.playlist.user.username)}</p>

        <div class="buttons">
            <button class="button is-large" on:click={tracksTable.playFirstTrack} title={t('play')}>
                <svg class="icon" viewBox="0 0 24 24">
                    <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                </svg>
            </button>

            <LikeButton token={data.token} item={data.playlist} itemRoute="playlists" itemLabel={t('playlist')} isLarge={true} />

            {#if data.playlist.user.id == data.authUser.id || data.authUser.role == 'admin'}
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
{#if data.playlist.tracks.length != 0}
    <TracksTable
        bind:this={tracksTable}
        token={data.token}
        authUser={data.authUser}
        tracks={data.playlist.tracks}
        inPlaylist={data.playlist}
    />
{:else}
    <p><i>{t('tracks_empty')}</i></p>
{/if}

<EditModal
    bind:this={editModal}
    token={data.token}
    playlist={data.playlist}
    on:updatePlaylist={(event) => {
        $sidebar.updateLastPlaylists();
        data.playlist = event.detail.playlist;
    }}
/>

<DeleteModal
    bind:this={deleteModal}
    token={data.token}
    item={data.playlist}
    itemRoute="playlists"
    itemLabel={t('playlist')}
    on:delete={() => {
        $sidebar.updateLastPlaylists();
        goto('/your_playlists');
    }}
/>
