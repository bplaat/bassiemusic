<script>
    import TracksTable from "../../../components/tracks-table.svelte";
    import AlbumCard from "../../../components/album-card.svelte";

    export let data;
    const { artist } = data;
</script>

<svelte:head>
    <title>{artist.name} - Artists - BassieMusic</title>
</svelte:head>

<div class="buttons">
    <button on:click={() => history.back()}>
        <svg class="icon" viewBox="0 0 24 24">
            <path
                d="M20,11V13H8L13.5,18.5L12.08,19.92L4.16,12L12.08,4.08L13.5,5.5L8,11H20Z"
            />
        </svg>
    </button>
</div>

<div class="columns">
    <div class="column is-one-quarter mr-5">
        <div class="box" style="padding: 0; overflow: hidden;">
            <img
                src={artist.image}
                alt="{artist.image}'s image"
                style="display: block;"
            />
        </div>
    </div>

    <div
        class="column"
        style="display: flex; flex-direction: column; justify-content: center;"
    >
        <h2 class="title">{artist.name}</h2>
    </div>
</div>

<h2 class="title">Top Tracks</h2>
<TracksTable tracks={artist.top_tracks} />

<h2 class="title" style="margin-top: 3rem;">Albums</h2>
{#if artist.albums != undefined}
    <div class="columns is-multiline">
        {#each artist.albums as album}
            <div class="column is-one-fifth">
                <AlbumCard {album} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>This artist has no albums</i></p>
{/if}
