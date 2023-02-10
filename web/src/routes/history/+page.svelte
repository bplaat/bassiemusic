<script>
    import { onMount,  onDestroy } from "svelte";
    import TracksTable from "../../components/tracks-table.svelte";

    export let data;
    let { token, authUser, tracks } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${authUser.id}/played_tracks?${new URLSearchParams({
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
            observer.unobserve(bottom);
        });
    }
</script>

<svelte:head>
    <title>Play History - BassieMusic</title>
</svelte:head>

<h1 class="title">Play History</h1>

{#if tracks.length > 0}
    <TracksTable {token} {tracks} />
{:else}
    <p>You have not listened to any tracks</p>
{/if}

<div bind:this={bottom} />
