<script>
    import AlbumCard from '../../../components/cards/album-card.svelte';
    import { lazyLoader } from '../../../utils.js';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: '$1 - Genres - BassieMusic',
            back: 'Go back one page',
            image_alt: 'Image of genre $1',
            albums: 'Albums',
            empty: "This genre doesn't have any albums",
        },
        nl: {
            title: '$1 - Genres - BassieMusic',
            back: 'Ga een pagina terug',
            image_alt: 'Afbeelding van genre $1',
            albums: 'Albums',
            empty: 'Dit genre heeft geen albums',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // State
    export let data;

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
        <div class="box has-image m-0 p-0">
            <figure class="image is-1by1">
                <img
                    src={data.genre.large_image || '/images/album-default.svg'}
                    alt={t('image_alt', data.genre.name)}
                />
            </figure>
        </div>
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title">{data.genre.name}</h2>
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
