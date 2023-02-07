<script>
    import ArtistCard from "../../../components/artist-card.svelte";

    export let data;
    let { token, authUser, artists } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${
                authUser.id
            }/liked_artists?${new URLSearchParams({
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
    <title>Artists - Liked - BassieMusic</title>
</svelte:head>

<div class="tabs">
    <ul>
        <li class="is-active"><a href="/liked/artists">Artists</a></li>
        <li><a href="/liked/albums">Albums</a></li>
        <li><a href="/liked/tracks">Tracks</a></li>
    </ul>
</div>

<h1 class="title">Liked Artists</h1>

{#if artists.length > 0}
    <div class="columns is-multiline">
        {#each artists as artist}
            <div class="column is-one-fifth">
                <ArtistCard {artist} />
            </div>
        {/each}
    </div>
{:else}
    <p>You have not liked any artists</p>
{/if}
