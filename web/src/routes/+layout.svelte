<script>
    import { browser } from '$app/environment';
    import { afterNavigate } from '$app/navigation';
    import Sidebar from '../components/sidebar.svelte';
    import MusicPlayer from '../components/music-player.svelte';
    import { musicPlayer } from '../stores.js';

    export let data;
    const { token, authUser, agent, lastTrack, lastTrackPosition } = data;

    // Init last played track
    if (browser && lastTrack) {
        musicPlayer.set({
            action: 'init',
            queue: [lastTrack],
            track_id: lastTrack.id,
            position: lastTrackPosition,
        });
    }

    // Sidebar
    let sidebar;
    afterNavigate(() => {
        document.body.scrollTop = 0;
        sidebar.close();
    });

    // Window is-resizing
    let windowResizeTimeout;
    function windowResize() {
        document.body.classList.add('is-resizing');
        if (windowResizeTimeout) clearTimeout(windowResizeTimeout);
        windowResizeTimeout = setTimeout(() => {
            windowResizeTimeout = undefined;
            document.body.classList.remove('is-resizing');
        }, 100);
    }
</script>

<svelte:window on:contextmenu|preventDefault={() => {}} on:resize={windowResize} />

<svelte:head>
    {#if authUser == undefined || (authUser != undefined && authUser.theme == 'system')}
        <link rel="stylesheet" href="/css/bulma-light.min.css" media="(prefers-color-scheme: light)" />
        <link rel="stylesheet" href="/css/bulma-dark.min.css" media="(prefers-color-scheme: dark)" />
    {/if}
    {#if authUser != undefined && authUser.theme == 'light'}
        <link rel="stylesheet" href="/css/bulma-light.min.css" />
    {/if}
    {#if authUser != undefined && authUser.theme == 'dark'}
        <link rel="stylesheet" href="/css/bulma-dark.min.css" />
    {/if}
    <link rel="stylesheet" href="/css/app.css" />
</svelte:head>

<div
    class="app"
    class:is-macos-app={agent.os == 'macOS' && agent.name == 'BassieMusic App'}
    class:is-windows-app={agent.os == 'Windows' && agent.name == 'BassieMusic App'}
    class:is-linux-app={agent.os == 'Linux' && agent.name == 'BassieMusic App'}
    class:is-light={authUser != undefined && authUser.theme == 'light'}
    class:is-dark={authUser != undefined && authUser.theme == 'dark'}
>
    <nav class="navbar has-background-white-bis is-fixed-top is-hidden-desktop">
        <div class="navbar-brand">
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a href="#" class="navbar-burger ml-0" on:click|preventDefault={() => sidebar.open()}>
                <span />
                <span />
                <span />
            </a>
            <div class="navbar-item" style="font-weight: 500;">BassieMusic</div>
        </div>
    </nav>

    <Sidebar bind:this={sidebar} {token} {authUser} />

    <div class="section">
        <slot />
    </div>

    <MusicPlayer {token} />
</div>
