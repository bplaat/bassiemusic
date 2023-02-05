<script>
    import Cookies from "js-cookie";
    import TracksTable from "../../components/tracks-table.svelte";

    export let data;
    const { tracks } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/tracks?${new URLSearchParams({
                page,
                limit: 50,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${Cookies.get("token")}`,
                },
            }
        );
        const { data: newTracks, pagination } = await response.json();
        tracks.push(...newTracks);
        if (tracks.length != pagination.total) {
            fetchPage(page + 1);
        }
    }
    if (tracks.length != data.total) {
        fetchPage(2);
    }
</script>

<svelte:head>
    <title>Tracks - BassieMusic</title>
</svelte:head>

<h2 class="title">Tracks</h2>
<TracksTable {tracks} />
