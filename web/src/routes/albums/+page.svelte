<script>
    import { onMount, onDestroy } from 'svelte';
    import AlbumCard from '../../components/cards/album-card.svelte';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Albums - BassieMusic',
            header: 'Albums',
            empty: "You don't have added any albums",
        },
        nl: {
            title: 'Albums - BassieMusic',
            header: 'Albums',
            empty: 'Je hebt nog geen enkele album toegevoegd',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;
    let { token, authUser, albums } = data;

    // Page fetcher
    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/albums?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newAlbums } = await response.json();
        albums.push(...newAlbums);
        albums = albums;
    }

    let bottom;
    if (albums.length != data.total) {
        let observer;
        onMount(() => {
            let page = 2;
            observer = new IntersectionObserver(
                (entries, observer) => {
                    for (const entry of entries) {
                        if (albums.length >= data.total) {
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

{#if albums.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each albums as album}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <AlbumCard {album} {token} {authUser} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>{t('empty')}</i></p>
{/if}

<div bind:this={bottom} />
