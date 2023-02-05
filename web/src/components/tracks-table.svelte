<script>
    export let tracks;
    import { playingTrack, playingQueue } from "../stores.js";
    import { formatDuration } from "../filters.js";

    function playTrack(track) {
        const index = tracks.indexOf(track);
        playingQueue.set(tracks);
        playingTrack.set(index);
    }
</script>

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
        <tr on:dblclick|preventDefault={()=> playTrack(track)}
            class:has-background-light={$playingQueue.length > 0 &&
            $playingQueue[$playingTrack].id == track.id}
            >
            <td>{index + 1}</td>
            <td style="display: flex;">
                <div class="box mr-4 mb-0" style="width: 64px; height: 64px; padding: 0; overflow: hidden;">
                    <img src={track.album.cover} alt="{track.title} album's cover" style="display: block;" />
                </div>
                <div style="flex: 1; display: flex; flex-direction: column; justify-content: center;">
                    <p>
                        <a href="/albums/{track.album.id}" style="font-weight: 500;">{track.title}</a>
                    </p>
                    <p>
                        {#if track.explicit}
                        <span class="tag is-danger mr-1">E</span>
                        {/if}
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
