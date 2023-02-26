<script>
    import { onMount, onDestroy } from 'svelte';
    import TracksTable from '../../../components/tracks-table.svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Tracks - Liked - BassieMusic',
            artists: 'Artists',
            albums: 'Albums',
            tracks: 'Tracks',
            header: 'Liked Tracks',
            empty: 'You have not liked any tracks',
        },
        nl: {
            title: 'Tracks - Geliked - BassieMusic',
            artists: 'Artisten',
            albums: 'Albums',
            tracks: 'Tracks',
            header: 'Gelikede Tracks',
            empty: 'Je hebt geen track geliked',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;
    let { token, authUser, tracks } = data;

    // Page fetcher
    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${authUser.id}/liked_tracks?${new URLSearchParams({
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

    onMount(() => {
        localStorage.setItem('liked-tab', 'tracks');
    });

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

<div class="tabs is-toggle">
    <ul>
        <li><a href="/liked/artists">{t('artists')}</a></li>
        <li><a href="/liked/albums">{t('albums')}</a></li>
        <li class="is-active"><a href="/liked/tracks">{t('tracks')}</a></li>
    </ul>
</div>

<h1 class="title">{t('header')}</h1>

{#if tracks.length > 0}
    <TracksTable {token} {authUser} {tracks} />
{:else}
    <p><i>{t('empty')}</i></p>
{/if}

<div bind:this={bottom} />
