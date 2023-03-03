<script>
    import { goto } from '$app/navigation';
    import GenreCard from '../../components/cards/genre-card.svelte';
    import AlbumCard from '../../components/cards/album-card.svelte';
    import ArtistCard from '../../components/cards/artist-card.svelte';
    import TracksTable from '../../components/tracks-table.svelte';
    import { lazyLoader } from '../../utils.js';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Search - BassieMusic',
            header: 'Search',
            query_placeholder: 'Find a track, album, artist or genre...',
            search: 'Search',
            empty: "Can't find anything with your search query",
            tracks: 'Tracks',
            artists: 'Artists',
            albums: 'Albums',
            genres: 'Genres',
        },
        nl: {
            title: 'Zoeken - BassieMusic',
            header: 'Zoeken',
            query_placeholder: 'Zoek een nummer, album, artiest of genre...',
            search: 'Zoeken',
            empty: 'Kan niets vinden met je zoekopdracht',
            tracks: 'Tracks',
            artists: 'Artisten',
            albums: 'Albums',
            genres: 'Genres',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;

    // Lazy loader
    lazyLoader(
        data.genresTotal,
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

    // Methods
    async function search() {
        if (data.query != '') {
            const newUrl = new URL(window.location.href);
            newUrl.searchParams.set('q', data.query);
            goto(newUrl);
        } else {
            goto('/search');
        }
    }
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<h2 class="title">{t('header')}</h2>

<form on:submit|preventDefault={search} class="field has-addons mb-5">
    <div class="control" style="width: 100%;">
        <!-- svelte-ignore a11y-autofocus -->
        <input
            class="input"
            type="text"
            placeholder={t('query_placeholder')}
            bind:value={data.query}
            autofocus={data.query == ''}
        />
    </div>
    <div class="control">
        <button type="submit" class="button is-link">{t('search')}</button>
    </div>
</form>

{#if data.searchResult != null}
    {#if data.searchResult.tracks.length == 0 && data.searchResult.artists.length == 0 && data.searchResult.albums.length == 0 && data.searchResult.genres.length == 0}
        <p><i>{t('empty')}</i></p>
    {:else}
        {#if data.searchResult.tracks.length > 0}
            <h2 class="title is-5">{t('tracks')}</h2>
            <TracksTable token={data.token} authUser={data.authUser} tracks={data.searchResult.tracks.slice(0, 5)} />
        {/if}

        {#if data.searchResult.artists.length > 0}
            <h2 class="title is-5 mt-5">{t('artists')}</h2>
            <div class="columns is-multiline mb-5">
                {#each data.searchResult.artists as artist}
                    <div
                        class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen"
                    >
                        <ArtistCard {artist} token={data.token} authUser={data.authUser} />
                    </div>
                {/each}
            </div>
        {/if}

        {#if data.searchResult.albums.length > 0}
            <h2 class="title is-5">{t('albums')}</h2>
            <div class="columns is-multiline mb-5">
                {#each data.searchResult.albums as album}
                    <div
                        class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen"
                    >
                        <AlbumCard {album} token={data.token} authUser={data.authUser} />
                    </div>
                {/each}
            </div>
        {/if}

        {#if data.searchResult.genres.length > 0}
            <h2 class="title is-5">{t('genres')}</h2>
            <div class="columns is-multiline is-mobile">
                {#each data.searchResult.genres as genre}
                    <div
                        class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen"
                    >
                        <GenreCard {genre} />
                    </div>
                {/each}
            </div>
        {/if}
    {/if}
{:else}
    <h2 class="title is-5">{t('genres')}</h2>
    <div class="columns is-multiline is-mobile">
        {#each data.genres as genre}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <GenreCard {genre} />
            </div>
        {/each}
    </div>
{/if}
