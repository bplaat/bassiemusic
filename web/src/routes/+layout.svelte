<script>
    import { musicPlayer } from '../stores.js';
    import { afterNavigate } from '$app/navigation';
    import Sidebar from '../components/sidebar.svelte';
    import MusicPlayer from '../components/music-player.svelte';

    export let data;
    const { token, authUser, agent, lastTrack, lastTrackPosition } = data;

    // Sidebar
    let sidebar;
    afterNavigate(() => {
        document.body.scrollTop = 0;
        sidebar.close();
    });

    // App is-resizing
    let resizing = false;
    let resizeTimeout;
    function windowResize() {
        resizing = true;
        if (resizeTimeout) clearTimeout(resizeTimeout);
        resizeTimeout = setTimeout(() => {
            resizeTimeout = undefined;
            resizing = false;
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
    class:is-playing={lastTrack != undefined}
    class:is-resizing={resizing}
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

    {#if lastTrack != undefined}
        <MusicPlayer
            bind:this={$musicPlayer}
            {token}
            queue={[lastTrack]}
            track={lastTrack}
            position={lastTrackPosition}
            duration={lastTrack.duration}
        />
    {:else}
        <MusicPlayer bind:this={$musicPlayer} {token} />
    {/if}
</div>

<style>
    .app {
        overflow-y: auto;
        margin-top: 52px;
        height: calc(100% - 52px);
    }
    .app.is-playing {
        margin-bottom: 10rem;
        height: calc(100% - 52px - 10rem);
    }
    @media (max-width: 1024px) {
        .section {
            padding: 1.5rem;
        }
    }
    @media (min-width: 1024px) {
        .app {
            margin-top: 0;
            margin-left: 16.5rem;
            height: 100%;
        }
        .app.app.is-playing {
            margin-bottom: 6rem;
            height: calc(100% - 6rem);
        }
    }

    /* macOS app */
    .app.is-macos-app .navbar {
        padding-top: 28px !important;
    }
    @media (max-width: 1024px) {
        .app.is-macos-app {
            margin-top: 80px !important;
        }
    }
</style>
