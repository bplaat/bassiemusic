<script>
    import { onMount, onDestroy } from 'svelte';
    import TracksTable from '../../components/tracks-table.svelte';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Tracks - BassieMusic',
            header: 'Tracks',
            empty: "You don't have added any tracks",
        },
        nl: {
            title: 'Tracks - BassieMusic',
            header: 'Tracks',
            empty: 'Je hebt nog geen enkele track toegevoegd',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;
    let { token, authUser, tracks } = data;

    // Page fetcher
    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/tracks?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newTracks } = await response.json();
        tracks.push(...newTracks);
        tracks = tracks;
    }

    let bottom;
    if (tracks.length != data.total) {
        let observer;
        onMount(() => {
            let page = 2;
            observer = new IntersectionObserver(
                (entries, observer) => {
                    for (const entry of entries) {
                        if (tracks.length >= data.total) {
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
{#if tracks.length > 0}
    <TracksTable {token} {authUser} {tracks} />
{:else}
    <p><i>{t('empty')}</i></p>
{/if}

<div bind:this={bottom} />
