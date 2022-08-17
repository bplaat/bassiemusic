<script context="module">
import { API_URL } from '../../config.js';

export async function load({ fetch }) {
    const response = await fetch(`${API_URL}/genres`, {
        headers: {
            Authorization: 'Bearer ' + localStorage.token
        }
    });
    return {
        status: response.status,
        props: {
            genres: response.ok && (await response.json())
        }
    };
};
</script>

<script>
export let genres;

async function fetchPage(page) {
    const response = await fetch(`${API_URL}/genres?page=${page}`, {
        headers: {
            Authorization: 'Bearer ' + localStorage.token
        }
    });
    const newGenres = await response.json();
    if (newGenres.length > 0) {
        genres.push(...newGenres);
        genres = genres;
        fetchPage(page + 1);
    }
}
fetchPage(2);
</script>

<svelte:head>
    <title>Genres - BassieMusic</title>
</svelte:head>

<h2 class="title">Genres</h2>

<div class="columns is-multiline">
    {#each genres as genre}
        <div class="column is-one-fifth">
            <a class="card" href="/genres/{genre.id}">
                <div class="card-image" style="background-image: url({genre.image});"></div>
                <div class="card-content">
                    <h3 class="title is-6">{genre.name}</h3>
                </div>
            </a>
        </div>
    {/each}
</div>
