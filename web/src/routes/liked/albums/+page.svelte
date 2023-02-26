<script>
    import { onMount, onDestroy } from 'svelte';
    import AlbumCard from '../../../components/cards/album-card.svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Albums - Liked - BassieMusic',
            artists: 'Artists',
            albums: 'Albums',
            tracks: 'Tracks',
            header: 'Liked Albums',
            empty: 'You have not liked any albums',
        },
        nl: {
            title: 'Albums - Geliked - BassieMusic',
            artists: 'Artisten',
            albums: 'Albums',
            tracks: 'Tracks',
            header: 'Gelikede Albums',
            empty: 'Je hebt geen albums geliked',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;
    let { token, authUser, albums } = data;

    // Page fetcher
    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${authUser.id}/liked_albums?${new URLSearchParams({
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

    onMount(() => {
        localStorage.setItem('liked-tab', 'albums');
    });

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

<div class="tabs is-toggle">
    <ul>
        <li><a href="/liked/artists">{t('artists')}</a></li>
        <li class="is-active"><a href="/liked/albums">{t('albums')}</a></li>
        <li><a href="/liked/tracks">{t('tracks')}</a></li>
    </ul>
</div>

<h1 class="title">{t('header')}</h1>

{#if albums.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each albums as album}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <AlbumCard {album} {token} {authUser} />
            </div>
        {/each}
    </div>
{:else}
    <p>{t('empty')}</p>
{/if}

<div bind:this={bottom} />
