<script context="module">
import { API_URL } from '../../config.js';

export async function load({ params, fetch }) {
    const response = await fetch(`${API_URL}/albums/${params.id}`, {
        headers: {
            Authorization: 'Bearer ' + localStorage.token
        }
    });
    const album = await response.json();
    for (const track of album.tracks) {
        track.album = album;
    }
    return {
        status: response.status,
        props: { album }
    };
};
</script>

<script>
import { playingTrack, playingQueue } from '../../stores.js';
import { formatDuration } from '../../filters.js';

export let album;

function playTrack(track) {
    const index = album.tracks.indexOf(track);
    playingQueue.set(album.tracks.slice());
    playingTrack.set(index);
}
</script>

<svelte:head>
    <title>{album.title} - Albums - BassieMusic</title>
</svelte:head>

<div class="buttons">
    <button on:click={() => history.back()}>
        <svg class="icon" viewBox="0 0 24 24">
            <path d="M20,11V13H8L13.5,18.5L12.08,19.92L4.16,12L12.08,4.08L13.5,5.5L8,11H20Z" />
        </svg>
    </button>
</div>

<div class="columns">
    <div class="column is-one-quarter mr-5">
        <div class="box" style="position: relative; padding: 0; overflow: hidden;">
            <img src="{album.cover}" alt="{album.title}'s cover" style="display: block;">
            <div class="card-image-tags">
                {#if album.type == 'album'}
                    <span class="tag">ALBUM</span>
                {/if}
                {#if album.type == 'ep'}
                    <span class="tag">EP</span>
                {/if}
                {#if album.type == 'single'}
                    <span class="tag">SINGLE</span>
                {/if}
                {#if album.explicit}
                    <span class="tag is-danger">E</span>
                {/if}
            </div>
        </div>
    </div>

    <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
        <h2 class="title mb-3">{album.title}</h2>
        <p class="mb-3">{album.released_at.split('T')[0]}</p>
        {#if album.genres != undefined}
            <p class="mb-3">
                {#each album.genres as genre}
                    <a href="/genres/{genre.id}" class="mr-2">{genre.name}</a>
                {/each}
            </p>
        {/if}
        <p>
            {#each album.artists as artist}
                <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
            {/each}
        </p>
    </div>
</div>

<h3 class="title is-4">Tracks</h3>

<table class="table is-fullwidth">
    <thead>
        <th style="width: 20%;">#</th>
        <th style="width: 60%;">Title</th>
        <th style="width: 20%;">Duration</th>
        <th style="width: 20%;">Plays</th>
    </thead>
    <tbody>
        {#each album.tracks as track, index}
            <tr on:dblclick|preventDefault={playTrack(track)}
                class:has-background-light="{$playingQueue.length > 0 && $playingQueue[$playingTrack].id == track.id}">
                <td>{index + 1}</td>
                <td>
                    <p style="font-weight: 500;">{track.title}</p>
                    <p>
                        {#if track.explicit}
                            <span class="tag is-danger mr-1">E</span>
                        {/if}
                        {#each track.artists as artist}
                            <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
                        {/each}
                    </p>
                </td>
                <td>{formatDuration(track.duration)}</td>
                <td>{track.plays}</td>
            </tr>
        {/each}
    </tbody>
</table>
