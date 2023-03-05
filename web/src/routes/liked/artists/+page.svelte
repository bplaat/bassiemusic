<script>
    import { onMount } from 'svelte';
    import SortByDropdown from '../../../components/sort-by-dropdown.svelte';
    import ArtistCard from '../../../components/cards/artist-card.svelte';
    import { lazyLoader } from '../../../utils.js';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Artists - Liked - BassieMusic',
            artists: 'Artists',
            albums: 'Albums',
            tracks: 'Tracks',
            header: 'Liked Artists',
            sort_by_liked_at_desc: 'Liked at (new - old)',
            sort_by_liked_at: 'Liked at (old - new)',
            sort_by_name: 'Name (A - Z)',
            sort_by_name_desc: 'Name (Z - A)',
            sort_by_sync: 'Sync (synced - not synced)',
            sort_by_sync_desc: 'Sync (not synced - synced)',
            sort_by_created_at_desc: 'Downloaded at (new - old)',
            sort_by_created_at: 'Downloaded at (old - new)',
            empty: 'You have not liked any artists',
        },
        nl: {
            title: 'Artisten - Geliked - BassieMusic',
            artists: 'Artisten',
            albums: 'Albums',
            tracks: 'Tracks',
            header: 'Gelikede Artisten',
            sort_by_liked_at_desc: 'Geliked op (nieuw - oud)',
            sort_by_liked_at: 'Geliked op (oud - nieuw)',
            sort_by_name: 'Naam (A - Z)',
            sort_by_name_desc: 'Naam (Z - A)',
            sort_by_sync: 'Gesynced (gesynced - niet gesynced)',
            sort_by_sync_desc: 'Gesynced (niet gesynced - gesynced)',
            sort_by_created_at_desc: 'Gedownload op (nieuw - oud)',
            sort_by_created_at: 'Gedownload op (oud - nieuw)',
            empty: 'Je hebt geen artist geliked',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;

    onMount(() => {
        localStorage.setItem('liked-tab', 'artists');
    });

    // Lazy loader
    lazyLoader(
        data.total,
        () => data.artists.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/users/${data.authUser.id}/liked_artists?${new URLSearchParams({
                    page,
                    sort_by: data.sortBy,
                })}`,
                {
                    headers: {
                        Authorization: `Bearer ${data.token}`,
                    },
                }
            );
            const { data: newArtists } = await response.json();
            data.artists = [...data.artists, ...newArtists];
        }
    );
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<div class="tabs is-toggle">
    <ul>
        <li class="is-active"><a href="/liked/artists">{t('artists')}</a></li>
        <li><a href="/liked/albums">{t('albums')}</a></li>
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
                name: t('sort_by_name'),
                name_desc: t('sort_by_name_desc'),
                sync: t('sort_by_sync'),
                sync_desc: t('sort_by_sync_desc'),
                created_at_desc: t('sort_by_created_at_desc'),
                created_at: t('sort_by_created_at'),
            }}
        />
    </div>
</div>

{#if data.artists.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each data.artists as artist}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <ArtistCard {artist} token={data.token} authUser={data.authUser} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>{t('empty')}</i></p>
{/if}
