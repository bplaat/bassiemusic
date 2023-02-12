<script>
    import AlbumCard from '../components/album-card.svelte';
    import ArtistCard from '../components/artist-card.svelte';
    import TracksTable from '../components/tracks-table.svelte';

    export let data;
    const { token, authUser, lastPlayedTracks } = data;

    function uniques(items) {
        const uniques = {};
        for (const item of items) {
            if (!uniques[item.id]) {
                uniques[item.id] = item;
            }
        }
        return Object.values(uniques);
    }

    $: lastPlayedAlbums = authUser != undefined ? uniques(lastPlayedTracks.map((track) => track.album)).slice(0, 5) : [];
    $: lastPlayedArtists =
        authUser != undefined ? uniques(lastPlayedTracks.map((track) => track.artists).flat()).slice(0, 5) : [];
</script>

<svelte:head>
    <title>Home - BassieMusic</title>
</svelte:head>

{#if authUser != undefined}
    <h1 class="title">Hey, {authUser.username}!</h1>

    {#if lastPlayedTracks.length > 0}
        <h2 class="title is-4">Last played albums</h2>
        <div class="columns is-multiline is-mobile">
            {#each lastPlayedAlbums as album}
                <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                    <AlbumCard {album} />
                </div>
            {/each}
        </div>

        <h2 class="title is-4 mt-5">Last played artists</h2>
        <div class="columns is-multiline is-mobile">
            {#each lastPlayedArtists as artist}
                <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                    <ArtistCard {artist} />
                </div>
            {/each}
        </div>

        <h2 class="title is-4 mt-5">Last played tracks</h2>
        <TracksTable {token} tracks={lastPlayedTracks.slice(0, 5)} />
    {:else}
        <p>You haven't listened to any tracks yet, use the sidebar to find something you like</p>
    {/if}
{:else}
    <h1 class="title">Welcome to BassieMusic</h1>
    <p>Login with your account to start listening to music</p>
{/if}
