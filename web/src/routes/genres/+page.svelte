<script>
    import SortByDropdown from '../../components/buttons/sort-by-dropdown.svelte';
    import GenreCard from '../../components/cards/genre-card.svelte';
    import { lazyLoader } from '../../utils.js';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Genres - BassieMusic',
            header: 'Genres',
            sort_by_name: 'Name (A - Z)',
            sort_by_name_desc: 'Name (Z - A)',
            sort_by_created_at_desc: 'Downloaded at (new - old)',
            sort_by_created_at: 'Downloaded at (old - new)',
            empty: "You don't have added any genres",
        },
        nl: {
            title: 'Genres - BassieMusic',
            header: 'Genres',
            sort_by_name: 'Naam (A - Z)',
            sort_by_name_desc: 'Naam (Z - A)',
            sort_by_created_at_desc: 'Gedownload op (nieuw - oud)',
            sort_by_created_at: 'Gedownload op (oud - nieuw)',
            empty: 'Je hebt nog geen enkele genre toegevoegd',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;

    // Lazy loader
    lazyLoader(
        data.total,
        () => data.genres.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/genres?${new URLSearchParams({
                    page,
                    sort_by: data.sortBy,
                })}`,
                {
                    headers: {
                        Authorization: `Bearer ${data.token}`,
                    },
                }
            );
            const { data: newGenres } = await response.json();
            data.genres = [...data.genres, ...newGenres];
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

{#if data.genres.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each data.genres as genre}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <GenreCard {genre} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>{t('empty')}</i></p>
{/if}
