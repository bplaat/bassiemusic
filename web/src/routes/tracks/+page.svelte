<script>
    import SortByDropdown from '../../components/buttons/sort-by-dropdown.svelte';
    import TracksTable from '../../components/tracks-table.svelte';
    import { lazyLoader } from '../../utils.js';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Tracks - BassieMusic',
            header: 'Tracks',
            sort_by_plays_desc: 'Plays (high - low)',
            sort_by_plays: 'Plays (low - high)',
            sort_by_title: 'Title (A - Z)',
            sort_by_title_desc: 'Title (Z - A)',
            sort_by_created_at_desc: 'Downloaded at (new - old)',
            sort_by_created_at: 'Downloaded at (old - new)',
            empty: "You don't have added any tracks",
        },
        nl: {
            title: 'Tracks - BassieMusic',
            header: 'Tracks',
            sort_by_plays_desc: 'Plays (hoog - laag)',
            sort_by_plays: 'Plays (laag - hoog)',
            sort_by_title: 'Titel (A - Z)',
            sort_by_title_desc: 'Titel (Z - A)',
            sort_by_created_at_desc: 'Gedownload op (nieuw - oud)',
            sort_by_created_at: 'Gedownload op (oud - nieuw)',
            empty: 'Je hebt nog geen enkele track toegevoegd',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;

    // Lazy loader
    lazyLoader(
        data.total,
        () => data.tracks.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/tracks?${new URLSearchParams({
                    page: page,
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

<div class="columns">
    <div class="column">
        <h2 class="title">{t('header')}</h2>
    </div>
    <div class="column">
        <SortByDropdown
            sortBy={data.sortBy}
            options={{
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
