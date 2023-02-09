<script>
    import { browser } from "$app/environment";
    import {  musicPlayer } from "../stores.js";
    import Sidebar from "../components/sidebar.svelte";
    import MusicPlayer from "../components/music-player.svelte";
    import { not_equal } from "svelte/internal";

    export let data;
    const { token, authUser, agent, lastTrack, lastTrackPosition } = data;

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
    <!-- Themes -->
    {#if authUser == undefined || (authUser != undefined && authUser.theme == "system")}
        <link rel="stylesheet" href="/css/bulma-system.min.css" />
        <style>
            .slider-container{position:relative;width:100%;height:5px;background-color:lightgray}
            .slider{position:absolute;width:100%;height:100%;background-color:#0099ff}
            .slider-thumb{position:absolute;top:-5px;width:15px;height:15px;background-color:#999999;border-radius:50%;cursor:pointer}

            ::-webkit-scrollbar-thumb{background-color:rgba(0,0,0,.3)}

            @media (prefers-color-scheme: dark) {
                .slider-container{position:relative;height:5px;background-color:lightgray}
                .slider{position:absolute;width:100%;height:100%;background-color:#0099ff}
                .slider-thumb{position:absolute;top:-5px;width:15px;height:15px;background-color:white;border-radius:50%;cursor:pointer}

                ::-webkit-scrollbar-thumb{background-color:rgba(255,255,255,.3)}
            }
        </style>
    {/if}

    {#if authUser != undefined && authUser.theme == "light"}
        <link rel="stylesheet" href="/css/bulma-light.min.css" />
        <style>
            .slider-container{position:relative;width:100%;height:5px;background-color:lightgray}
            .slider{position:absolute;width:100%;height:100%;background-color:#0099ff}
            .slider-thumb{position:absolute;top:-5px;width:15px;height:15px;background-color:#999999;border-radius:50%;cursor:pointer}

            ::-webkit-scrollbar-thumb{background-color:rgba(0,0,0,.3)}
        </style>
    {/if}

    {#if authUser != undefined && authUser.theme == "dark"}
        <link rel="stylesheet" href="/css/bulma-dark.min.css" />
        <style>
            .slider-container{position:relative;height:5px;background-color:lightgray}
            .slider{position:absolute;width:100%;height:100%;background-color:#0099ff}
            .slider-thumb{position:absolute;top:-5px;width:15px;height:15px;background-color:white;border-radius:50%;cursor:pointer}

            ::-webkit-scrollbar-thumb{background-color:rgba(255,255,255,.3)}
        </style>
    {/if}

    <!-- Custom CSS for BassieMusic apps -->
    {#if agent.os == "macOS" && agent.name == "BassieMusic App"}
        <style>
            .sidebar{padding-top:calc(28px + 1.25rem)!important}
            .macos-is-fullscreen .sidebar{padding-top:calc(1.25rem)!important}
        </style>
    {/if}
</svelte:head>

<Sidebar {token} {authUser} />

<div class="section">
    <slot />
</div>

<MusicPlayer {token} />
