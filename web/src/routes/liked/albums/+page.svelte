<script>
    import AlbumCard from "../../../components/album-card.svelte";

    export let data;
    let { token, authUser, albums } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${
                authUser.id
            }/liked_albums?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newAlbums, pagination } = await response.json();
        albums.push(...newAlbums);
        albums = albums;
        if (albums.length != pagination.total) {
            fetchPage(page + 1);
        }
    }
    if (albums.length != data.total) {
        fetchPage(2);
    }
</script>

<svelte:head>
    <title>Albums - Liked - BassieMusic</title>
</svelte:head>

<div class="tabs">
    <ul>
        <li><a href="/liked/artists">Artists</a></li>
        <li class="is-active"><a href="/liked/albums">Albums</a></li>
        <li><a href="/liked/tracks">Tracks</a></li>
    </ul>
</div>

<div class="content">
    <h1 class="title">Liked Albums</h1>

    {#if albums.length > 0}
        <div class="columns is-multiline">
            {#each albums as album}
                <div class="column is-one-fifth">
                    <AlbumCard {album} />
                </div>
            {/each}
        </div>
    {:else}
        <p>You have not liked any albums</p>
    {/if}
</div>
