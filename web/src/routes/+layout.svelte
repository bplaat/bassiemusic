<script>
    import { browser } from "$app/environment";
    import {  musicPlayer } from "../stores.js";
    import Sidebar from "../components/sidebar.svelte";
    import MusicPlayer from "../components/music-player.svelte";

    export let data;
    const { token, authUser, lastTrack, lastTrackPosition } = data;

    if (browser && lastTrack) {
        musicPlayer.set({
            action: 'init',
            queue: [lastTrack],
            index: 0,
            position: lastTrackPosition
        });
    }
</script>

<svelte:head>
    {#if authUser}
        {#if authUser.theme == "system"}
            <link rel="stylesheet" href="/css/bulma-system.min.css" />
        {/if}
        {#if authUser.theme == "light"}
            <link rel="stylesheet" href="/css/bulma-light.min.css" />
        {/if}
        {#if authUser.theme == "dark"}
            <link rel="stylesheet" href="/css/bulma-dark.min.css" />
        {/if}
    {:else}
        <link rel="stylesheet" href="/css/bulma-system.min.css" />
    {/if}
</svelte:head>

<Sidebar {token} {authUser} />

<div class="section">
    <slot />
</div>

<MusicPlayer {token} />
