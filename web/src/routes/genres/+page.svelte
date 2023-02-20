<script>
    import { onDestroy, onMount } from 'svelte';
    import GenreCard from '../../components/cards/genre-card.svelte';

    export let data;
    let { token, genres } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/genres?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newGenres } = await response.json();
        genres.push(...newGenres);
        genres = genres;
    }

    let bottom;
    if (genres.length != data.total) {
        let observer;
        onMount(() => {
            let page = 2;
            observer = new IntersectionObserver(
                (entries, observer) => {
                    for (const entry of entries) {
                        if (genres.length >= data.total) {
                            observer.unobserve(entry.target);
                        } else {
                            fetchPage(page++);
                        }
                    }
                },
                {
                    root: document.body,
                }
            );
            observer.observe(bottom);
        });
        onDestroy(() => {
            if (observer) observer.unobserve(bottom);
        });
    }
</script>

<svelte:head>
    <title>Genres - BassieMusic</title>
</svelte:head>

<h2 class="title">Genres</h2>

{#if genres.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each genres as genre}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <GenreCard {genre} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>You don't have added any genres</i></p>
{/if}

<div bind:this={bottom} />
