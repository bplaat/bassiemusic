<script>
    import TracksTable from '../../components/tracks-table.svelte';
    import { musicState } from '../../stores';

    export let data;
    const { token, authUser } = data;
</script>

<svelte:head>
    <title>Play Queue - BassieMusic</title>
</svelte:head>

<div class="buttons">
    <button class="button" on:click={() => history.back()} title="Go back one page">
        <svg class="icon" viewBox="0 0 24 24">
            <path d="M20,11V13H8L13.5,18.5L12.08,19.92L4.16,12L12.08,4.08L13.5,5.5L8,11H20Z" />
        </svg>
    </button>
</div>

<h2 class="title">Play Queue</h2>
{#if $musicState.queue.length > 0}
    <TracksTable {token} {authUser} isMusicQueue={true} tracks={$musicState.queue} />
{:else}
    <p><i>The music player queue is empty</i></p>
{/if}
