<script>
    import Cookies from "js-cookie";
    import AlbumCard from "../../components/album-card.svelte";

    export let data;
    const { albums } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/albums?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${Cookies.get("token")}`,
                },
            }
        );
        const { data: newAlbums, pagination } = await response.json();
        albums.push(...newAlbums);
        if (albums.length != pagination.total) {
            fetchPage(page + 1);
        }
    }
    if (albums.length != data.total) {
        fetchPage(2);
    }
</script>

<svelte:head>
    <title>Albums - BassieMusic</title>
</svelte:head>

<h2 class="title">Albums</h2>

<div class="columns is-multiline">
    {#each albums as album}
        <div class="column is-one-fifth">
            <AlbumCard {album} />
        </div>
    {/each}
</div>
