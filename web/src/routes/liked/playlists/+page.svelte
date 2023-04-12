<script>
    import { onMount } from 'svelte';
    import SortByDropdown from '../../../components/sort-by-dropdown.svelte';
    import PlaylistCard from '../../../components/cards/playlist-card.svelte';
    import { lazyLoader } from '../../../utils.js';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Playlists - Liked - BassieMusic',
            artists: 'Artists',
            genres: 'Genres',
            albums: 'Albums',
            tracks: 'Tracks',
            playlists: 'Playlists',
            header: 'Liked Playlists',
            sort_by_liked_at_desc: 'Liked at (new - old)',
            sort_by_liked_at: 'Liked at (old - new)',
            sort_by_title: 'Title (A - Z)',
            sort_by_title_desc: 'Title (Z - A)',
            sort_by_public: 'Public (public - private)',
            sort_by_public_desc: 'Public (private - public)',
            sort_by_released_at_desc: 'Released at (new - old)',
            sort_by_released_at: 'Released at (old - new)',
            sort_by_created_at_desc: 'Downloaded at (new - old)',
            sort_by_created_at: 'Downloaded at (old - new)',
            empty: 'You have not liked any playlists',
        },
        nl: {
            title: 'Afspeellijsten - Geliked - BassieMusic',
            artists: 'Artisten',
            genres: 'Genres',
            albums: 'Albums',
            tracks: 'Tracks',
            playlists: 'Afspeellijsten',
            header: 'Gelikede Afspeellijsten',
            sort_by_liked_at_desc: 'Geliked op (nieuw - oud)',
            sort_by_liked_at: 'Geliked op (oud - nieuw)',
            sort_by_title: 'Titel (A - Z)',
            sort_by_title_desc: 'Titel (Z - A)',
            sort_by_public: 'Publiekelijk (publiekelijk - persoonlijk)',
            sort_by_public_desc: 'Publiekelijk (persoonlijk - publiekelijk)',
            sort_by_released_at_desc: 'Verschenen op (nieuw - oud)',
            sort_by_released_at: 'Verschenen op (oud - nieuw)',
            sort_by_created_at_desc: 'Gedownload op (nieuw - oud)',
            sort_by_created_at: 'Gedownload op (oud - nieuw)',
            empty: 'Je hebt geen playlists geliked',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;

    onMount(() => {
        localStorage.setItem('liked-tab', 'playlists');
    });

    // Lazy loader
    lazyLoader(
        data.total,
        () => data.playlists.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/users/${data.authUser.id}/liked_playlists?${new URLSearchParams({
                    page,
                    sort_by: data.sortBy,
                })}`,
                {
                    headers: {
                        Authorization: `Bearer ${data.token}`,
                    },
                }
            );
            const { data: newPlaylists } = await response.json();
            data.playlists = [...data.playlists, ...newPlaylists];
        }
    );
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<div class="tabs is-toggle">
    <ul>
        <li><a href="/liked/artists">{t('artists')}</a></li>
        <li><a href="/liked/genres">{t('genres')}</a></li>
        <li><a href="/liked/albums">{t('albums')}</a></li>
        <li><a href="/liked/tracks">{t('tracks')}</a></li>
        <li class="is-active"><a href="/liked/playlists">{t('playlists')}</a></li>
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
                public: t('sort_by_public'),
                public_desc: t('sort_by_public_desc'),
                released_at_desc: t('sort_by_released_at_desc'),
                released_at: t('sort_by_released_at'),
                created_at_desc: t('sort_by_created_at_desc'),
                created_at: t('sort_by_created_at'),
            }}
        />
    </div>
</div>

{#if data.playlists.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each data.playlists as playlist}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <PlaylistCard {playlist} token={data.token} authUser={data.authUser} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>{t('empty')}</i></p>
{/if}
