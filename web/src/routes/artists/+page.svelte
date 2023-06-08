<script>
    import { getContext } from 'svelte';
    import SortByDropdown from '../../components/buttons/sort-by-dropdown.svelte';
    import ArtistCard from '../../components/cards/artist-card.svelte';
    import { newLazyLoader } from '../../utils.js';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Artists - BassieMusic',
            header: 'Artists',
            sort_by_name: 'Name (A - Z)',
            sort_by_name_desc: 'Name (Z - A)',
            sort_by_sync: 'Sync (synced - not synced)',
            sort_by_sync_desc: 'Sync (not synced - synced)',
            sort_by_created_at_desc: 'Downloaded at (new - old)',
            sort_by_created_at: 'Downloaded at (old - new)',
            empty: "You don't have added any artists",
        },
        nl: {
            title: 'Artiesten - BassieMusic',
            header: 'Artiesten',
            sort_by_name: 'Naam (A - Z)',
            sort_by_name_desc: 'Naam (Z - A)',
            sort_by_sync: 'Gesynced (gesynced - niet gesynced)',
            sort_by_sync_desc: 'Gesynced (niet gesynced - gesynced)',
            sort_by_created_at_desc: 'Gedownload op (nieuw - oud)',
            sort_by_created_at: 'Gedownload op (oud - nieuw)',
            empty: 'Je hebt nog geen enkele artiest toegevoegd',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;
    const token = getContext('token');
    const authUser = getContext('authUser');

    // Lazy loader
    const limit = 20;
    let artists = new Array(limit).fill(null);
    let total = null;
    newLazyLoader({
        getTotal: () => total,
        getCount: () => artists.length,
        async loadPage(page) {
            if (page == 1 && artists.length > 0) artists = new Array(limit).fill(null);
            if (page != 1) artists = [...artists, ...new Array(limit).fill(null)];
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/artists?${new URLSearchParams({ page, limit, sort_by: data.sortBy })}`,
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                },
            );
            const { data: newArtists, pagination } = await response.json();
            if (page == 1) total = pagination.total;
            artists.splice((page - 1) * limit, limit, ...newArtists);
            artists = artists;
        },
    });
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

{#if artists.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each artists as artist}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <ArtistCard {artist} {token} {authUser} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>{t('empty')}</i></p>
{/if}
