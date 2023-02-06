<script>
    import ArtistCard from "../../components/artist-card.svelte";

    export let data;
    let { token, artists } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/artists?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newArtists, pagination } = await response.json();
        artists.push(...newArtists);
        artists = artists;
        if (artists.length != pagination.total) {
            fetchPage(page + 1);
        }
    }
    if (artists.length != data.total) {
        fetchPage(2);
    }
</script>

<svelte:head>
    <title>Artists - BassieMusic</title>
</svelte:head>

<h2 class="title">Artists</h2>

<div class="columns is-multiline">
    {#each artists as artist}
        <div class="column is-one-fifth">
            <ArtistCard {artist} />
        </div>
    {/each}
</div>
