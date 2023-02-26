<script>
    import { onMount, onDestroy } from 'svelte';
    import ArtistCard from '../../components/cards/artist-card.svelte';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Artists - BassieMusic',
            header: 'Artists',
            empty: "You don't have added any artists",
        },
        nl: {
            title: 'Artists - BassieMusic',
            header: 'Artists',
            empty: 'Je hebt nog geen enkele artist toegevoegd',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;
    let { token, artists } = data;

    // Page fetcher
    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/artists?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newArtists } = await response.json();
        artists.push(...newArtists);
        artists = artists;
    }

    let bottom;
    if (artists.length != data.total) {
        let observer;
        onMount(() => {
            let page = 2;
            observer = new IntersectionObserver(
                (entries, observer) => {
                    for (const entry of entries) {
                        if (artists.length >= data.total) {
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
    <title>{t('title')}</title>
</svelte:head>

<h2 class="title">{t('header')}</h2>

{#if artists.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each artists as artist}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <ArtistCard {artist} {token} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>{t('empty')}</i></p>
{/if}

<div bind:this={bottom} />
