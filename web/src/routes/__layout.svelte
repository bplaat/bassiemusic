<script>
import Sidebar from '../components/Sidebar.svelte';
import MusicPlayer from '../components/MusicPlayer.svelte';
</script>

<Sidebar />

<div class="section">
    <slot />
</div>

<MusicPlayer />
