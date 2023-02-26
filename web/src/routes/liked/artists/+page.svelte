<script>
    import { onMount, onDestroy } from 'svelte';
    import ArtistCard from '../../../components/cards/artist-card.svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Artists - Liked - BassieMusic',
            artists: 'Artists',
            albums: 'Albums',
            tracks: 'Tracks',
            header: 'Liked Artists',
            empty: 'You have not liked any artists',
        },
        nl: {
            title: 'Artisten - Geliked - BassieMusic',
            artists: 'Artisten',
            albums: 'Albums',
            tracks: 'Tracks',
            header: 'Gelikede Artisten',
            empty: 'Je hebt geen artist geliked',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;
    let { token, authUser, artists } = data;

    // Page fetcher
    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${authUser.id}/liked_artists?${new URLSearchParams({
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

    onMount(() => {
        localStorage.setItem('liked-tab', 'artists');
    });

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

<div class="tabs is-toggle">
    <ul>
        <li class="is-active"><a href="/liked/artists">{t('artists')}</a></li>
        <li><a href="/liked/albums">{t('albums')}</a></li>
        <li><a href="/liked/tracks">{t('tracks')}</a></li>
    </ul>
</div>

<h1 class="title">{t('header')}</h1>

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
