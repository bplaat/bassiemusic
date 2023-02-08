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
            track_id: lastTrack.id,
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
            <style>
                /* slider dark mode */
                .slider-container {position:relative;margin-top:9px;width:100%;height:5px;background-color:lightgray;}
                .slider {position:absolute;width:100%;height: 100%;background-color:#0099ff;}
                .slider-thumb {position:absolute;top:-5px;width:15px;height:15px;background-color:white;border-radius:50%;cursor: pointer;}
            </style>
        {:else}
            <style>
                /* slider light mode */
                .slider-container {position:relative;margin-top:9px;width:100%;height:5px;background-color:lightgray;}
                .slider {position:absolute;width:100%;height: 100%;background-color:#0099ff;}
                .slider-thumb {position:absolute;top:-5px;width:15px;height:15px;background-color:#999999;border-radius:50%;cursor: pointer;}
            </style>
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
