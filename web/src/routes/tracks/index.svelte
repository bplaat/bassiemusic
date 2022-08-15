<script context="module">
import { API_URL } from '../../config.js';

export async function load({ fetch }) {
    const response = await fetch(`${API_URL}/tracks?limit=50`);
    return {
        status: response.status,
        props: {
            tracks: response.ok && (await response.json())
        }
    };
};
</script>

<script>
import { playingTrack, playingQueue } from '../../stores.js';
import { formatDuration } from '../../filters.js';

export let tracks;

async function fetchPage(page) {
    const response = await fetch(`${API_URL}/tracks?page=${page}&limit=50`);
    const newTracks = await response.json();
    if (newTracks.length > 0) {
        tracks.push(...newTracks);
        tracks = tracks;
        fetchPage(page + 1);
    }
}
fetchPage(2);

function playTrack(track) {
    const index = tracks.indexOf(track);
    playingQueue.set(tracks);
    playingTrack.set(index);
}
</script>

<svelte:head>
    <title>Tracks - BassieMusic</title>
</svelte:head>

<h2 class="title">Tracks</h2>

<table class="table" style="width: 100%;">
    <thead>
        <th style="width: 10%;">#</th>
        <th style="width: 30%;">Title</th>
        <th style="width: 30%;">Album</th>
        <th style="width: 15%;">Duration</th>
        <th style="width: 15%;">Plays</th>
    </thead>
    <tbody>
        {#each tracks as track, index}
            <tr on:dblclick|preventDefault={playTrack(track)}
                class:has-background-light="{$playingQueue.length > 0 && $playingQueue[$playingTrack].id == track.id}">
                <td>{index + 1}</td>
                <td style="display: flex;">
                    <div class="box mr-4 mb-0" style="width: 64px; height: 64px; padding: 0; overflow: hidden;">
                        <img src="{track.album.cover}" alt="{track.title} album's cover" style="display: block;">
                    </div>
                    <div style="flex: 1; display: flex; flex-direction: column; justify-content: center;">
                        <p><a href="/albums/{track.album.id}" style="font-weight: bold;">{track.title}</a></p>
                        <p>
                            {#each track.artists as artist}
                                <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
                            {/each}
                        </p>
                    </div>
                </td>
                <td><a href="/albums/{track.album.id}">{track.album.title}</a></td>
                <td>{formatDuration(track.duration)}</td>
                <td>{track.plays}</td>
            </tr>
        {/each}
    </tbody>
</table>
