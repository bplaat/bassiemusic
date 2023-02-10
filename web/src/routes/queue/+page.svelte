<script>
    import TracksTable from '../../components/tracks-table.svelte';
    import { musicPlayer } from '../../stores';

    export let data;
    const { token } = data;
</script>

<svelte:head>
    <title>Play Queue - BassieMusic</title>
</svelte:head>

<div class="buttons">
    <button class="button" on:click={() => history.back()}>
        <svg class="icon" viewBox="0 0 24 24">
            <path d="M20,11V13H8L13.5,18.5L12.08,19.92L4.16,12L12.08,4.08L13.5,5.5L8,11H20Z" />
        </svg>
    </button>
</div>

<h2 class="title">Play Queue</h2>
{#if $musicPlayer.queue.length > 0}
    <TracksTable {token} isMusicQueue={true} tracks={$musicPlayer.queue} />
{:else}
    <p><i>The music player queue is empty</i></p>
{/if}
