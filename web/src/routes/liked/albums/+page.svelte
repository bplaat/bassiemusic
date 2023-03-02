<script>
    import { onMount } from 'svelte';
    import SortByDropdown from '../../../components/sort-by-dropdown.svelte';
    import AlbumCard from '../../../components/cards/album-card.svelte';
    import { lazyLoader } from '../../../utils.js';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Albums - Liked - BassieMusic',
            artists: 'Artists',
            albums: 'Albums',
            tracks: 'Tracks',
            header: 'Liked Albums',
            sort_by_liked_at_desc: 'Liked at (new - old)',
            sort_by_liked_at: 'Liked at (old - new)',
            sort_by_title: 'Title (A - Z)',
            sort_by_title_desc: 'Title (Z - A)',
            sort_by_released_at_desc: 'Released at (new - old)',
            sort_by_released_at: 'Released at (old - new)',
            sort_by_created_at_desc: 'Downloaded at (new - old)',
            sort_by_created_at: 'Downloaded at (old - new)',
            empty: 'You have not liked any albums',
        },
        nl: {
            title: 'Albums - Geliked - BassieMusic',
            artists: 'Artisten',
            albums: 'Albums',
            tracks: 'Tracks',
            header: 'Gelikede Albums',
            sort_by_liked_at_desc: 'Geliked op (nieuw - oud)',
            sort_by_liked_at: 'Geliked op (oud - nieuw)',
            sort_by_title: 'Titel (A - Z)',
            sort_by_title_desc: 'Titel (Z - A)',
            sort_by_released_at_desc: 'Verschenen op (nieuw - oud)',
            sort_by_released_at: 'Verschenen op (oud - nieuw)',
            sort_by_created_at_desc: 'Gedownload op (nieuw - oud)',
            sort_by_created_at: 'Gedownload op (oud - nieuw)',
            empty: 'Je hebt geen albums geliked',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;

    onMount(() => {
        localStorage.setItem('liked-tab', 'albums');
    });

    // Lazy loader
    lazyLoader(
        data.total,
        () => data.albums.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/users/${data.authUser.id}/liked_albums?${new URLSearchParams({
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

<div class="tabs is-toggle">
    <ul>
        <li><a href="/liked/artists">{t('artists')}</a></li>
        <li class="is-active"><a href="/liked/albums">{t('albums')}</a></li>
        <li><a href="/liked/tracks">{t('tracks')}</a></li>
    </ul>
</div>

<div class="columns">
    <div class="column">
        <h2 class="title">{t('header')}</h2>
    </div>
    <div class="column">
        <SortByDropdown
            sortBy={data.sortBy}
            options={{
                liked_at_desc: t('sort_by_liked_at_desc'),
                liked_at: t('sort_by_liked_at'),
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
    <p>{t('empty')}</p>
{/if}
