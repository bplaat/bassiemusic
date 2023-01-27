<script>
    import Cookies from "js-cookie";

    export let data;
    const { artists } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/artists?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${Cookies.get("token")}`,
                },
            }
        );
        const { data: newArtists, pagination } = await response.json();
        artists.push(...newArtists);
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
            <a class="card" href="/artists/{artist.id}">
                <div
                    class="card-image"
                    style="background-image: url({artist.image});"
                />
                <div class="card-content">
                    <h3 class="title is-6">{artist.name}</h3>
                </div>
            </a>
        </div>
    {/each}
</div>
