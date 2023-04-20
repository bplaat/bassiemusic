<script>
    import { formatBytes } from '../../../filters.js';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Downloader - Admin - BassieMusic',
            header: 'Admin Downloader',

            storage_size: 'Storage folder size',
            storage_used: 'Used: $1',
            storage_max: 'Max: $1',

            search_header: 'Search and download albums and artists',
            query_placeholder: 'Find an album or artist...',
            search: 'Search',
            albums: 'Albums',
            cover_alt: 'Cover of album $1',
            add_album: 'Add album to BassieMusic',
            albums_empty: "Can't find any albums on Deezer",
            artists: 'Artists',
            image_alt: 'Image of artist $1',
            add_artist: 'Add artist to BassieMusic',
            artists_empty: "Can't find any artists on Deezer",
        },
        nl: {
            title: 'Downloader - Admin - BassieMusic',
            header: 'Admin Downloader',

            storage_size: 'Storage folder groote',
            storage_used: 'Gebruikt: $1',
            storage_max: 'Max: $1',

            search_header: 'Zoek en download albums en artisten',
            query_placeholder: 'Vind een album of artist...',
            search: 'Zoeken',
            albums: 'Albums',
            cover_alt: 'Hoes van album $1',
            add_album: 'Voeg album toe aan BassieMusic',
            albums_empty: 'Kan geen albums vinden op Deezer',
            artists: 'Artiesten',
            image_alt: 'Afbeelding van artist $1',
            add_artist: 'Voeg artist toe aan BassieMusic',
            artists_empty: 'Kan geen artisten vinden op Deezer',
        },
    };
    const t = (key, p1 = '') => lang[$language][key].replace('$1', p1);

    // State
    export let data;
    let query = '';
    let results = false;
    let albums = [];
    let artists = [];

    // Methods
    async function search() {
        if (query === '') {
            results = false;
            albums = [];
            artists = [];
            return;
        }

        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/deezer_search?${new URLSearchParams({
                q: query,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${data.token}`,
                },
            }
        );
        const result = await response.json();
        results = true;
        albums = result.albums;
        artists = result.artists;
    }

    async function downloadAlbum(album) {
        await fetch(`${import.meta.env.VITE_API_URL}/download/album`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${data.token}`,
            },
            body: new URLSearchParams({
                deezer_id: album.id,
            }),
        });
        albums = albums.filter((otherAlbum) => otherAlbum.id !== album.id);
    }

    async function downloadArtist(artist) {
        await fetch(`${import.meta.env.VITE_API_URL}/download/artist`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${data.token}`,
            },
            body: new URLSearchParams({
                deezer_id: artist.id,
            }),
        });
        artists = artists.filter((otherArtist) => otherArtist.id !== artist.id);
    }
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<h1 class="title">{t('header')}</h1>

<div class="box">
    <h2 class="title is-4">{t('storage_size')}</h2>
    <progress class="progress is-link" value={data.storage.used} max={data.storage.max}>
        {((data.storage.used / data.storage.max) * 100).toFixed(2)}%
    </progress>
    <p>
        <span class="mr-3">{t('storage_used', formatBytes(data.storage.used))}</span>
        <span>{t('storage_max', formatBytes(data.storage.max))}</span>
    </p>
</div>

<div class="box">
    <h2 class="title is-4">{t('search_header')}</h2>

    <form on:submit|preventDefault={search} class="field has-addons">
        <div class="control" style="width: 100%;">
            <input class="input" type="text" bind:value={query} placeholder={t('query_placeholder')} />
        </div>
        <div class="control">
            <button type="submit" class="button is-link">{t('search')}</button>
        </div>
    </form>

    {#if results}
        <div class="columns mt-5">
            <div class="column is-half">
                <h2 class="title is-4">{t('albums')}</h2>
                {#each albums as album}
                    <div class="media">
                        <div class="media-left">
                            <div class="box m-0 p-0" style="width: 48px; height: 48px;">
                                <img src={album.cover_medium} alt={t('cover_alt', album.title)} />
                            </div>
                        </div>
                        <div class="media-content" style="min-width: 0;">
                            <p class="ellipsis" style="font-weight: 500;">{album.title}</p>
                            <p class="ellipsis">{album.artist.name}</p>
                        </div>
                        <div class="media-right">
                            <button class="button is-link" on:click={() => downloadAlbum(album)} title={t('add_album')}>
                                <svg class="icon" viewBox="0 0 24 24">
                                    <path d="M19,13H13V19H11V13H5V11H11V5H13V11H19V13Z" />
                                </svg>
                            </button>
                        </div>
                    </div>
                {/each}
                {#if albums.length === 0}
                    <p><i>{t('albums_empty')}</i></p>
                {/if}
            </div>

            <div class="column is-half">
                <h2 class="title is-4">{t('artists')}</h2>
                {#each artists as artist}
                    <div class="media">
                        <div class="media-left">
                            <div class="box m-0 p-0" style="width: 48px; height: 48px;">
                                <img src={artist.picture_medium} alt={t('image_alt', artist.name)} />
                            </div>
                        </div>
                        <div class="media-content">
                            <p class="ellipsis" style="font-weight: 500;">{artist.name}</p>
                        </div>
                        <button class="button is-link" on:click={() => downloadArtist(artist)} title={t('add_artist')}>
                            <svg class="icon" viewBox="0 0 24 24">
                                <path d="M19,13H13V19H11V13H5V11H11V5H13V11H19V13Z" />
                            </svg>
                        </button>
                    </div>
                {/each}
                {#if artists.length === 0}
                    <p><i>{t('artists_empty')}</i></p>
                {/if}
            </div>
        </div>
    {/if}
</div>
