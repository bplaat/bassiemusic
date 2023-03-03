<script>
    import SortByDropdown from '../../components/sort-by-dropdown.svelte';
    import ArtistCard from '../../components/cards/artist-card.svelte';
    import { lazyLoader } from '../../utils.js';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Artists - BassieMusic',
            header: 'Artists',
            sort_by_name: 'Name (A - Z)',
            sort_by_name_desc: 'Name (Z - A)',
            sort_by_created_at_desc: 'Downloaded at (new - old)',
            sort_by_created_at: 'Downloaded at (old - new)',
            empty: "You don't have added any artists",
        },
        nl: {
            title: 'Artists - BassieMusic',
            header: 'Artists',
            sort_by_name: 'Naam (A - Z)',
            sort_by_name_desc: 'Naam (Z - A)',
            sort_by_created_at_desc: 'Gedownload op (nieuw - oud)',
            sort_by_created_at: 'Gedownload op (oud - nieuw)',
            empty: 'Je hebt nog geen enkele artist toegevoegd',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;

    // Lazy loader
    lazyLoader(
        data.total,
        () => data.artists.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/artists?${new URLSearchParams({
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
