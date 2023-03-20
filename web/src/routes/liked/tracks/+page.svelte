<script>
    import { onMount } from 'svelte';
    import SortByDropdown from '../../../components/sort-by-dropdown.svelte';
    import TracksTable from '../../../components/tracks-table.svelte';
    import { lazyLoader } from '../../../utils.js';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Tracks - Liked - BassieMusic',
            artists: 'Artists',
            albums: 'Albums',
            tracks: 'Tracks',
            playlists: 'Playlists',
            header: 'Liked Tracks',
            sort_by_liked_at_desc: 'Liked at (new - old)',
            sort_by_liked_at: 'Liked at (old - new)',
            sort_by_plays_desc: 'Plays (high - low)',
            sort_by_plays: 'Plays (low - high)',
            sort_by_title: 'Title (A - Z)',
            sort_by_title_desc: 'Title (Z - A)',
            sort_by_created_at_desc: 'Downloaded at (new - old)',
            sort_by_created_at: 'Downloaded at (old - new)',
            empty: 'You have not liked any tracks',
        },
        nl: {
            title: 'Tracks - Geliked - BassieMusic',
            artists: 'Artisten',
            albums: 'Albums',
            tracks: 'Tracks',
            playlists: 'Afspeellijsten',
            header: 'Gelikede Tracks',
            sort_by_liked_at_desc: 'Geliked op (nieuw - oud)',
            sort_by_liked_at: 'Geliked op (oud - nieuw)',
            sort_by_plays_desc: 'Plays (hoog - laag)',
            sort_by_plays: 'Plays (laag - hoog)',
            sort_by_title: 'Titel (A - Z)',
            sort_by_title_desc: 'Titel (Z - A)',
            sort_by_created_at_desc: 'Gedownload op (nieuw - oud)',
            sort_by_created_at: 'Gedownload op (oud - nieuw)',
            empty: 'Je hebt geen track geliked',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;

    onMount(() => {
        localStorage.setItem('liked-tab', 'tracks');
    });

    // Lazy loader
    lazyLoader(
        data.total,
        () => data.tracks.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/users/${data.authUser.id}/liked_tracks?${new URLSearchParams({
                    page,
                    sort_by: data.sortBy,
                })}`,
                {
                    headers: {
                        Authorization: `Bearer ${data.token}`,
                    },
                }
            );
            const { data: newTracks } = await response.json();
            data.tracks = [...data.tracks, ...newTracks];
        }
    );
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<div class="tabs is-toggle">
    <ul>
        <li><a href="/liked/artists">{t('artists')}</a></li>
        <li><a href="/liked/albums">{t('albums')}</a></li>
        <li class="is-active"><a href="/liked/tracks">{t('tracks')}</a></li>
        <li><a href="/liked/playlists">{t('playlists')}</a></li>
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
                plays_desc: t('sort_by_plays_desc'),
                plays: t('sort_by_plays'),
                title: t('sort_by_title'),
                title_desc: t('sort_by_title_desc'),
                created_at_desc: t('sort_by_created_at_desc'),
                created_at: t('sort_by_created_at'),
            }}
        />
    </div>
</div>

{#if data.tracks.length > 0}
    <TracksTable token={data.token} authUser={data.authUser} tracks={data.tracks} />
{:else}
    <p><i>{t('empty')}</i></p>
{/if}
