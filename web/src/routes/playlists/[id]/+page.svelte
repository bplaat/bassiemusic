<script>
    import { goto } from '$app/navigation';
    import ImageEditButton from '../../../components/buttons/image-edit-button.svelte';
    import LikeButton from '../../../components/buttons/like-button.svelte';
    import EditModal from '../../../components/modals/playlists/edit-modal.svelte';
    import DeleteModal from '../../../components/modals/delete-modal.svelte';
    import TracksTable from '../../../components/tracks-table.svelte';
    import { sidebar, language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: '$1 - Playlists - BassieMusic',
            back: 'Go back one page',
            made_by: 'Made by',
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
            made_by: 'Gemaakt door',
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

    $: isOwner =
        data.authUser.role === 'admin' ||
        data.playlist.owners.map((owner) => owner.username).indexOf(data.authUser.username) !== -1;
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
        <ImageEditButton
            token={data.token}
            item={data.playlist}
            itemRoute="playlists"
            editable={isOwner}
            on:update={() => {
                $sidebar.updateLastPlaylists();
            }}
        />
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title">{data.playlist.name}</h2>
        <p class="mb-5">
            {t('made_by')}
            {#each data.playlist.owners as owner}
                <span class="mr-2">{owner.username}</span>
            {/each}
        </p>

        <div class="buttons">
            <button class="button is-large" on:click={tracksTable.playFirstTrack} title={t('play')}>
                <svg class="icon" viewBox="0 0 24 24">
                    <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                </svg>
            </button>

            <LikeButton
                token={data.token}
                item={data.playlist}
                itemRoute="playlists"
                itemLabel={t('playlist')}
                isLarge={true}
            />

            {#if isOwner}
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
{#if data.playlist.tracks.length !== 0}
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

{#if isOwner}
    <EditModal
        bind:this={editModal}
        token={data.token}
        playlist={data.playlist}
        on:update={(event) => {
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
{/if}
