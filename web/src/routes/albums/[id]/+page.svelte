<script>
    import { goto } from '$app/navigation';
    import DeleteModal from '../../../components/modals/delete-modal.svelte';
    import TracksTable from '../../../components/tracks-table.svelte';
    import LikeButton from '../../../components/like-button.svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: '$1 - Albums - BassieMusic',
            back: 'Go back one page',
            cover_alt: 'Cover of album $1',
            explicit: 'Explicit lyrics',
            play: 'Play album',
            album: 'album',
            delete: 'Delete album',
            tracks: 'Tracks',
            tracks_empty: "This album doesn't have any tracks",
        },
        nl: {
            title: '$1 - Albums - BassieMusic',
            back: 'Ga een pagina terug',
            cover_alt: 'Hoes van album $1',
            explicit: 'Expliciete songtekst',
            play: 'Speel album',
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
        <div class="box has-image p-0 has-image-tags">
            <figure class="image is-1by1">
                <img
                    src={data.album.large_cover || '/images/album-default.svg'}
                    alt={t('cover_alt', data.album.title)}
                />
            </figure>
            <div class="image-tags">
                {#if data.album.type == 'album'}
                    <span class="tag">ALBUM</span>
                {/if}
                {#if data.album.type == 'ep'}
                    <span class="tag">EP</span>
                {/if}
                {#if data.album.type == 'single'}
                    <span class="tag">SINGLE</span>
                {/if}
                {#if data.album.explicit}
                    <span class="tag is-danger" title={t('explicit')}>E</span>
                {/if}
            </div>
        </div>
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

{#if data.authUser.role == 'admin'}
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
