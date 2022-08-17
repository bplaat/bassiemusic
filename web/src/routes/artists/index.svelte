<script context="module">
import { API_URL } from '../../config.js';

export async function load({ fetch }) {
    const response = await fetch(`${API_URL}/artists`, {
        headers: {
            Authorization: 'Bearer ' + localStorage.token
        }
    });
    return {
        status: response.status,
        props: {
            artists: response.ok && (await response.json())
        }
    };
};
</script>

<script>
export let artists;

async function fetchPage(page) {
    const response = await fetch(`${API_URL}/artists?page=${page}`, {
        headers: {
            Authorization: 'Bearer ' + localStorage.token
        }
    });
    const newArtists = await response.json();
    if (newArtists.length > 0) {
        artists.push(...newArtists);
        artists = artists;
        fetchPage(page + 1);
    }
}
fetchPage(2);
</script>

<svelte:head>
    <title>Artists - BassieMusic</title>
</svelte:head>

<h2 class="title">Artists</h2>

<div class="columns is-multiline">
    {#each artists as artist}
        <div class="column is-one-fifth">
            <a class="card" href="/artists/{artist.id}">
                <div class="card-image" style="background-image: url({artist.image});"></div>
                <div class="card-content">
                    <h3 class="title is-6">{artist.name}</h3>
                </div>
            </a>
        </div>
    {/each}
</div>
