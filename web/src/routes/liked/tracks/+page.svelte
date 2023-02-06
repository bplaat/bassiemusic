<script>
    import TracksTable from "../../../components/tracks-table.svelte";

    export let data;
    let { token, authUser, tracks } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${
                authUser.id
            }/liked_tracks?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newTracks, pagination } = await response.json();
        tracks.push(...newTracks);
        tracks = tracks;
        if (tracks.length != pagination.total) {
            fetchPage(page + 1);
        }
    }
    if (tracks.length != data.total) {
        fetchPage(2);
    }
</script>

<svelte:head>
    <title>Tracks - Liked - BassieMusic</title>
</svelte:head>

<div class="tabs">
    <ul>
        <li><a href="/liked/artists">Artists</a></li>
        <li><a href="/liked/albums">Albums</a></li>
        <li class="is-active"><a href="/liked/tracks">Tracks</a></li>
    </ul>
</div>

<div class="content">
    <h1 class="title">Liked Tracks</h1>

    {#if tracks.length > 0}
        <TracksTable {token} {tracks} />
    {:else}
        <p>You have not liked any tracks</p>
    {/if}
</div>
