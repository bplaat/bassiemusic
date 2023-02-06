<script>
    import AlbumCard from "../../../components/album-card.svelte";

    export let data;
    const { genre } = data;
</script>

<svelte:head>
    <title>{genre.name} - Genres - BassieMusic</title>
</svelte:head>

<div class="buttons">
    <button class="button" on:click={() => history.back()}>
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
                src={genre.large_image}
                alt="{genre.large_image}'s image"
                style="display: block;"
            />
        </div>
    </div>

    <div
        class="column"
        style="display: flex; flex-direction: column; justify-content: center;"
    >
        <h2 class="title">{genre.name}</h2>
    </div>
</div>

<h2 class="title">Albums</h2>
{#if genre.albums != undefined}
    <div class="columns is-multiline">
        {#each genre.albums as album}
            <div class="column is-one-fifth">
                <AlbumCard {album} />
            </div>
        {/each}
    </div>
{:else}
    <p><i>This genre has no albums</i></p>
{/if}
