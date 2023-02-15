<script>
    import { onMount, onDestroy } from 'svelte';
    import TracksTable from '../../../components/tracks-table.svelte';

    export let data;
    let { token, authUser, tracks } = data;

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
    <title>Tracks - Liked - BassieMusic</title>
</svelte:head>

<div class="tabs is-toggle">
    <ul>
        <li><a href="/liked/artists">Artists</a></li>
        <li><a href="/liked/albums">Albums</a></li>
        <li class="is-active"><a href="/liked/tracks">Tracks</a></li>
    </ul>
</div>

<h1 class="title">Liked Tracks</h1>

{#if tracks.length > 0}
    <TracksTable {token} {tracks} />
{:else}
    <p><i>You have not liked any tracks</i></p>
{/if}

<div bind:this={bottom} />
