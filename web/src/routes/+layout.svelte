<script>
    import Sidebar from "../components/sidebar.svelte";
    import MusicPlayer from "../components/music-player.svelte";

    export let data;
    const { authUser } = data;
</script>

<Sidebar {authUser} />

<div class="section">
    <slot />
</div>

<MusicPlayer />
