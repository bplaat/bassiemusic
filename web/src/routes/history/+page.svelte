<script>
    import TracksTable from "../../components/tracks-table.svelte";

    export let data;
    let { token, authUser, tracks } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${
                authUser.id
            }/played_tracks?${new URLSearchParams({
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
    <title>Play History - BassieMusic</title>
</svelte:head>

<h1 class="title">Play History</h1>

{#if tracks.length > 0}
    <TracksTable {token} {tracks} />
{:else}
    <p>You have not listened to any tracks</p>
{/if}
