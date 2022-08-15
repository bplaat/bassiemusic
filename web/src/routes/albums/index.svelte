<script context="module">
import { API_URL } from '../../config.js';

export async function load({ fetch }) {
    const response = await fetch(`${API_URL}/albums`);
    return {
        status: response.status,
        props: {
            albums: response.ok && (await response.json())
        }
    };
};
</script>

<script>
import AlbumCard from '../../components/AlbumCard.svelte';

export let albums;

async function fetchPage(page) {
    const response = await fetch(`${API_URL}/albums?page=${page}`);
    const newAlbums = await response.json();
    if (newAlbums.length > 0) {
        albums.push(...newAlbums);
        albums = albums;
        fetchPage(page + 1);
    }
}
fetchPage(2);
</script>

<svelte:head>
    <title>Albums - BassieMusic</title>
</svelte:head>

<h2 class="title">Albums</h2>

<div class="columns is-multiline">
    {#each albums as album}
        <div class="column is-one-fifth">
            <AlbumCard album={album} />
        </div>
    {/each}
</div>
