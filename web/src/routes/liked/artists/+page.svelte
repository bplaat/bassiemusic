<script>
    import { onMount, onDestroy } from 'svelte';
    import ArtistCard from '../../../components/artist-card.svelte';

    export let data;
    let { token, authUser, artists } = data;

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
    <title>Artists - Liked - BassieMusic</title>
</svelte:head>

<div class="tabs is-toggle">
    <ul>
        <li class="is-active"><a href="/liked/artists">Artists</a></li>
        <li><a href="/liked/albums">Albums</a></li>
        <li><a href="/liked/tracks">Tracks</a></li>
    </ul>
</div>

<h1 class="title">Liked Artists</h1>

{#if artists.length > 0}
    <div class="columns is-multiline is-mobile">
        {#each artists as artist}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <ArtistCard {artist} />
            </div>
        {/each}
    </div>
{:else}
    <p>You have not liked any artists</p>
{/if}

<div bind:this={bottom} />
