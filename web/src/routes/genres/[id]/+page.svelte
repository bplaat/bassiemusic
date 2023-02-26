<script>
    import { page } from '$app/stores';
    import { onMount, onDestroy } from 'svelte';
    import AlbumCard from '../../../components/cards/album-card.svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: '$1 - Genres - BassieMusic',
            back: 'Go back one page',
            image_alt: 'Image of genre $1',
            albums: 'Albums',
            empty: "This genre doesn't have any albums",
        },
        nl: {
            title: '$1 - Genres - BassieMusic',
            back: 'Ga een pagina terug',
            image_alt: 'Afbeelding van genre $1',
            albums: 'Albums',
            empty: "Dit genre heeft geen albums",
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // State
    export let data;
    let { token, authUser, genre } = data;

    // Reresh genre on page change to same page
    let unsubscribe;
    onMount(() => {
        unsubscribe = page.subscribe(async (page) => {
            if (page.url.pathname.startsWith('/genres/') && page.url.pathname != `/genres/${genre.id}`) {
                const response = await fetch(`${import.meta.env.VITE_API_URL}/genres/${page.params.id}`, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });
                genre = await response.json();
            }
        });
    });
    onDestroy(() => {
        if (unsubscribe) {
            unsubscribe();
        }
    });

    // Genre albums
    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/genres/${genre.id}/albums?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newAlbums } = await response.json();
        genre.albums.push(...newAlbums);
        genre = genre;
    }

    let bottom;
    if (genre.albums.length != data.albumsTotal) {
        let observer;
        onMount(() => {
            let page = 2;
            observer = new IntersectionObserver(
                (entries, observer) => {
                    for (const entry of entries) {
                        if (genre.albums.length >= data.albumsTotal) {
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
    <title>{t('title', genre.name)}</title>
</svelte:head>

<div class="buttons">
    <button class="button" on:click={() => history.back()} title={t('back')}>
        <svg class="icon" viewBox="0 0 24 24">
            <path d="M20,11V13H8L13.5,18.5L12.08,19.92L4.16,12L12.08,4.08L13.5,5.5L8,11H20Z" />
        </svg>
    </button>
</div>

<div class="columns">
    <div class="column is-one-quarter mr-5 mr-0-mobile">
        <div class="box has-image m-0 p-0" style="aspect-ratio: 1;">
            <img src={genre.large_image} alt={t('image_alt', genre.name)} />
        </div>
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title">{genre.name}</h2>
    </div>
</div>

<h2 class="title">{t('albums')}</h2>
{#if genre.albums != undefined}
    <div class="columns is-multiline is-mobile">
        {#each genre.albums as album}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <AlbumCard {album} {token} {authUser} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>{t('empty')}</i></p>
{/if}

<div bind:this={bottom} />
