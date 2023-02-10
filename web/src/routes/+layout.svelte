<script>
    import { browser } from "$app/environment";
    import { afterNavigate } from "$app/navigation";
    import Sidebar from "../components/sidebar.svelte";
    import MusicPlayer from "../components/music-player.svelte";
    import { musicPlayer } from "../stores.js";

    export let data;
    const { token, authUser, agent, lastTrack, lastTrackPosition } = data;

    if (browser && lastTrack) {
        musicPlayer.set({
            action: "init",
            queue: [lastTrack],
            track_id: lastTrack.id,
            position: lastTrackPosition,
        });
    }

    let isSidebarOpen = false;
    afterNavigate(() => {
        isSidebarOpen = false;
        document.body.scrollTop = 0;
    });
</script>

<svelte:window on:contextmenu|preventDefault={() => {}} />

<svelte:head>
    <!-- Themes -->
    {#if authUser == undefined || (authUser != undefined && authUser.theme == "system")}
        <link
            rel="stylesheet"
            href="/css/bulma-light.min.css"
            media="(prefers-color-scheme: light)"
        />
        <link
            rel="stylesheet"
            href="/css/bulma-dark.min.css"
            media="(prefers-color-scheme: dark)"
        />
        <style>
            .card-image > img {
                background-color: #ccc;
            }
            .slider-container {
                position: relative;
                width: 100%;
                height: 5px;
                background-color: lightgray;
            }
            .slider {
                position: absolute;
                width: 100%;
                height: 100%;
                background-color: #0099ff;
            }
            .slider-thumb {
                position: absolute;
                top: -5px;
                width: 15px;
                height: 15px;
                background-color: #999999;
                border-radius: 50%;
                cursor: pointer;
            }
            ::-webkit-scrollbar-thumb {
                background-color: rgba(0, 0, 0, 0.3);
            }

            @media (prefers-color-scheme: dark) {
                .card-image > img {
                    background-color: #333;
                }
                .slider-container {
                    background-color: lightgray;
                }
                .slider-thumb {
                    background-color: white;
                }
                ::-webkit-scrollbar-thumb {
                    background-color: rgba(255, 255, 255, 0.3);
                }
            }
        </style>
    {/if}

    {#if authUser != undefined && authUser.theme == "light"}
        <link rel="stylesheet" href="/css/bulma-light.min.css" />
        <style>
            .card-image > img {
                background-color: #ccc;
            }
            .slider-container {
                position: relative;
                width: 100%;
                height: 5px;
                background-color: lightgray;
            }
            .slider {
                position: absolute;
                width: 100%;
                height: 100%;
                background-color: #0099ff;
            }
            .slider-thumb {
                position: absolute;
                top: -5px;
                width: 15px;
                height: 15px;
                background-color: #999999;
                border-radius: 50%;
                cursor: pointer;
            }
            ::-webkit-scrollbar-thumb {
                background-color: rgba(0, 0, 0, 0.3);
            }
        </style>
    {/if}

    {#if authUser != undefined && authUser.theme == "dark"}
        <link rel="stylesheet" href="/css/bulma-dark.min.css" />
        <style>
            .card-image > img {
                background-color: #333;
            }
            .slider-container {
                position: relative;
                height: 5px;
                background-color: lightgray;
            }
            .slider {
                position: absolute;
                width: 100%;
                height: 100%;
                background-color: #0099ff;
            }
            .slider-thumb {
                position: absolute;
                top: -5px;
                width: 15px;
                height: 15px;
                background-color: white;
                border-radius: 50%;
                cursor: pointer;
            }
            ::-webkit-scrollbar-thumb {
                background-color: rgba(255, 255, 255, 0.3);
            }
        </style>
    {/if}

    <!-- Custom CSS for BassieMusic apps -->
    {#if agent.os == "macOS" && agent.name == "BassieMusic App"}
        <style>
            .navbar {
                padding-top: 28px !important;
            }
            @media (max-width: 1024px) {
                body {
                    margin-top: calc(80px - 1.5rem) !important;
                }
            }
            .sidebar {
                padding-top: calc(28px + 1.25rem) !important;
            }
            .macos-is-fullscreen .sidebar {
                padding-top: calc(1.25rem) !important;
            }
        </style>
    {/if}
</svelte:head>

<nav class="navbar is-light is-fixed-top is-hidden-desktop">
    <div class="navbar-brand">
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a
            href="#"
            class="navbar-burger ml-0"
            on:click|preventDefault={() => (isSidebarOpen = true)}
        >
            <span />
            <span />
            <span />
        </a>
        <a class="navbar-item" href="/" style="font-weight: 500;">BassieMusic</a
        >
    </div>
</nav>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
    class="modal-background"
    class:is-hidden={!isSidebarOpen}
    on:click={() => (isSidebarOpen = false)}
    style="z-index: 200;"
/>
<Sidebar open={isSidebarOpen} {token} {authUser} />

<div class="section">
    <slot />
</div>

<MusicPlayer {token} />
