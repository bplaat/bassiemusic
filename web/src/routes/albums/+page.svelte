<script>
    import SortByDropdown from '../../components/sort-by-dropdown.svelte';
    import AlbumCard from '../../components/cards/album-card.svelte';
    import { lazyLoader } from '../../utils.js';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Albums - BassieMusic',
            header: 'Albums',
            sort_by_title: 'Title (A - Z)',
            sort_by_title_desc: 'Title (Z - A)',
            sort_by_released_at_desc: 'Released at (new - old)',
            sort_by_released_at: 'Released at (old - new)',
            sort_by_created_at_desc: 'Downloaded at (new - old)',
            sort_by_created_at: 'Downloaded at (old - new)',
            empty: "You don't have added any albums",
        },
        nl: {
            title: 'Albums - BassieMusic',
            header: 'Albums',
            sort_by_title: 'Titel (A - Z)',
            sort_by_title_desc: 'Titel (Z - A)',
            sort_by_released_at_desc: 'Verschenen op (nieuw - oud)',
            sort_by_released_at: 'Verschenen op (oud - nieuw)',
            sort_by_created_at_desc: 'Gedownload op (nieuw - oud)',
            sort_by_created_at: 'Gedownload op (oud - nieuw)',
            empty: 'Je hebt nog geen enkele album toegevoegd',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;

    // Lazy loader
    lazyLoader(
        data.total,
        () => data.albums.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/albums?${new URLSearchParams({
                    page,
                    sort_by: data.sortBy,
                })}`,
                {
                    headers: {
                        Authorization: `Bearer ${data.token}`,
                    },
                }
            );
            const { data: newAlbums } = await response.json();
            data.albums = [...data.albums, ...newAlbums];
        }
    );
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<div class="columns">
    <div class="column">
        <h2 class="title">{t('header')}</h2>
    </div>
    <div class="column">
        <SortByDropdown
            sortBy={data.sortBy}
            options={{
                title: t('sort_by_title'),
                title_desc: t('sort_by_title_desc'),
                released_at_desc: t('sort_by_released_at_desc'),
                released_at: t('sort_by_released_at'),
                created_at_desc: t('sort_by_created_at_desc'),
                created_at: t('sort_by_created_at'),
            }}
        />
    </div>
</div>

{#if data.albums.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each data.albums as album}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <AlbumCard {album} token={data.token} authUser={data.authUser} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>{t('empty')}</i></p>
{/if}
