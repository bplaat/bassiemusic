<script>
    import { onMount,  onDestroy } from "svelte";
    import AlbumCard from "../../components/album-card.svelte";

    export let data;
    let { token, albums } = data;

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
            observer.unobserve(bottom);
        });
    }
</script>

<svelte:head>
    <title>Albums - BassieMusic</title>
</svelte:head>

<h2 class="title">Albums</h2>

<div class="columns is-multiline is-mobile">
    {#each albums as album}
        <div
            class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen"
        >
            <AlbumCard {album} />
        </div>
    {/each}
</div>

<div bind:this={bottom} />
