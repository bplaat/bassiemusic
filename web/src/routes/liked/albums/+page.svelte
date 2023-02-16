<script>
    import { onMount, onDestroy } from 'svelte';
    import AlbumCard from '../../../components/album-card.svelte';

    export let data;
    let { token, authUser, albums } = data;

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
    <title>Albums - Liked - BassieMusic</title>
</svelte:head>

<div class="tabs is-toggle">
    <ul>
        <li><a href="/liked/artists">Artists</a></li>
        <li class="is-active"><a href="/liked/albums">Albums</a></li>
        <li><a href="/liked/tracks">Tracks</a></li>
    </ul>
</div>

<h1 class="title">Liked Albums</h1>

{#if albums.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each albums as album}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <AlbumCard {album} {token} {authUser} />
            </div>
        {/each}
    </div>
{:else}
    <p>You have not liked any albums</p>
{/if}

<div bind:this={bottom} />
