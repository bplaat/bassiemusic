<script>
    import { page } from '$app/stores';
    import { browser } from '$app/environment';
    import AlbumCard from '../../../components/album-card.svelte';

    export let data;
    let { token, genre } = data;

    if (browser) {
        page.subscribe(async (page) => {
            if (page.url.pathname.startsWith('/genres/') && page.url.pathname != `/genres/${genre.id}`) {
                const response = await fetch(`${import.meta.env.VITE_API_URL}/genres/${page.params.id}`, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });
                genre = await response.json();
            }
        });
    }
</script>

<svelte:head>
    <title>{genre.name} - Genres - BassieMusic</title>
</svelte:head>

<div class="buttons">
    <button class="button" on:click={() => history.back()}>
        <svg class="icon" viewBox="0 0 24 24">
            <path d="M20,11V13H8L13.5,18.5L12.08,19.92L4.16,12L12.08,4.08L13.5,5.5L8,11H20Z" />
        </svg>
    </button>
</div>

<div class="columns">
    <div class="column is-one-quarter mr-5 mr-0-mobile">
        <div class="box m-0 p-0">
            <img style="aspect-ratio: 1;" src={genre.large_image} alt="Image of genre {genre.name}" loading="lazy" />
        </div>
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title">{genre.name}</h2>
    </div>
</div>

<h2 class="title">Albums</h2>
{#if genre.albums != undefined}
    <div class="columns is-multiline is-mobile">
        {#each genre.albums as album}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <AlbumCard {album} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>This genre has no albums</i></p>
{/if}
