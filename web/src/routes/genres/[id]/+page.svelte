<script>
    import { goto } from '$app/navigation';
    import ImageEditButton from '../../../components/buttons/image-edit-button.svelte';
    import LikeButton from '../../../components/buttons/like-button.svelte';
    import EditModal from '../../../components/modals/genres/edit-modal.svelte';
    import DeleteModal from '../../../components/modals/delete-modal.svelte';
    import AlbumCard from '../../../components/cards/album-card.svelte';
    import { lazyLoader } from '../../../utils.js';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: '$1 - Genres - BassieMusic',
            back: 'Go back one page',
            edit: 'Edit genre',
            genre: 'genre',
            delete: 'Delete genre',
            albums: 'Albums',
            empty: "This genre doesn't have any albums",
        },
        nl: {
            title: '$1 - Genres - BassieMusic',
            back: 'Ga een pagina terug',
            edit: 'Verander genre',
            genre: 'genre',
            delete: 'Verwijder genre',
            albums: 'Albums',
            empty: 'Dit genre heeft geen albums',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // State
    export let data;
    let editModal;
    let deleteModal;

    // Lazy loader
    lazyLoader(
        data.albumsTotal,
        () => data.genre.albums.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/genres/${data.genre.id}/albums?${new URLSearchParams({
                    page,
                })}`,
                {
                    headers: {
                        Authorization: `Bearer ${data.token}`,
                    },
                }
            );
            const { data: newAlbums } = await response.json();
            data.genre.albums = [...data.genre.albums, ...newAlbums];
        }
    );
</script>

<svelte:head>
    <title>{t('title', data.genre.name)}</title>
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
            item={data.genre}
            itemRoute="genres"
            editable={data.authUser.role === 'admin'}
        />
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title">{data.genre.name}</h2>

        <div class="buttons">
            <LikeButton token={data.token} item={data.genre} itemRoute="genres" itemLabel={t('genre')} isLarge={true} />

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

<h2 class="title">{t('albums')}</h2>
{#if data.genre.albums.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each data.genre.albums as album}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <AlbumCard {album} token={data.token} authUser={data.authUser} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>{t('empty')}</i></p>
{/if}

{#if data.authUser.role === 'admin'}
    <EditModal
        bind:this={editModal}
        token={data.token}
        genre={data.genre}
        on:update={(event) => {
            data.genre = event.detail.genre;
        }}
    />

    <DeleteModal
        bind:this={deleteModal}
        token={data.token}
        item={data.genre}
        itemRoute="genres"
        itemLabel={t('genre')}
        on:delete={() => {
            goto('/genres');
        }}
    />
{/if}
