<script>
    import {
        trackAutoplay,
        trackPosition,
        playingQueue,
        playingTrack,
    } from "../stores.js";
    import { formatDuration } from "../filters.js";

    export let token;
    export let tracks;
    export let showAlbum = true;

    function playTrack(track) {
        const index = tracks.indexOf(track);
        trackAutoplay.set(true);
        trackPosition.set(0);
        playingQueue.set(tracks);
        playingTrack.set(index);
    }

    export function playFirstTrack() {
        playTrack(tracks[0]);
    }

    function likeTrack(track) {
        fetch(
            `${import.meta.env.VITE_API_URL}/tracks/${track.id}/like${
                track.liked ? "/delete" : ""
            }`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        track.liked = !track.liked;
        tracks = tracks;
    }
</script>

<table class="table" style="width: 100%;">
    <thead>
        <th style="width: 10%;">#</th>
        <th style="width: 30%;">Title</th>
        {#if showAlbum}
            <th style="width: 30%;">Album</th>
        {/if}
        <th style="width: 15%;">Duration</th>
        <th style="width: 14%;">Plays</th>
        <th style="width: 1%;" />
    </thead>
    <tbody>
        {#each tracks as track, index}
            <tr
                on:dblclick|preventDefault={() => playTrack(track)}
                class:has-background-light={$playingQueue.length > 0 &&
                    $playingQueue[$playingTrack].id == track.id}
            >
                <td>{index + 1}</td>
                {#if showAlbum}
                    <td style="display: flex;">
                        <div
                            class="box mr-4 mb-0"
                            style="width: 64px; height: 64px; padding: 0; overflow: hidden;"
                        >
                            <img
                                src={track.album.small_cover}
                                alt="{track.title} album's cover"
                                style="display: block;"
                            />
                        </div>
                        <div
                            style="flex: 1; display: flex; flex-direction: column; justify-content: center;"
                        >
                            <p>
                                <a
                                    href="/albums/{track.album.id}"
                                    style="font-weight: 500;">{track.title}</a
                                >
                            </p>
                            <p>
                                {#if track.explicit}
                                    <span class="tag is-danger mr-1">E</span>
                                {/if}
                                {#each track.artists as artist}
                                    <a href="/artists/{artist.id}" class="mr-2"
                                        >{artist.name}</a
                                    >
                                {/each}
                            </p>
                        </div>
                    </td>
                {:else}
                    <td>
                        <p style="font-weight: 500;">{track.title}</p>
                        <p>
                            {#if track.explicit}
                                <span class="tag is-danger mr-1">E</span>
                            {/if}
                            {#each track.artists as artist}
                                <a href="/artists/{artist.id}" class="mr-2"
                                    >{artist.name}</a
                                >
                            {/each}
                        </p>
                    </td>
                {/if}
                {#if showAlbum}
                    <td
                        ><a href="/albums/{track.album.id}"
                            >{track.album.title}</a
                        ></td
                    >
                {/if}
                <td>{formatDuration(track.duration)}</td>
                <td>{track.plays}</td>
                <td>
                    <div class="buttons">
                        <button
                            class="button"
                            on:click={() => likeTrack(track)}
                        >
                            {#if track.liked}
                                <svg
                                    class="icon is-colored"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        fill="#f14668"
                                        d="M12,21.35L10.55,20.03C5.4,15.36 2,12.27 2,8.5C2,5.41 4.42,3 7.5,3C9.24,3 10.91,3.81 12,5.08C13.09,3.81 14.76,3 16.5,3C19.58,3 22,5.41 22,8.5C22,12.27 18.6,15.36 13.45,20.03L12,21.35Z"
                                    />
                                </svg>
                            {:else}
                                <svg class="icon" viewBox="0 0 24 24">
                                    <path
                                        d="M12.1,18.55L12,18.65L11.89,18.55C7.14,14.24 4,11.39 4,8.5C4,6.5 5.5,5 7.5,5C9.04,5 10.54,6 11.07,7.36H12.93C13.46,6 14.96,5 16.5,5C18.5,5 20,6.5 20,8.5C20,11.39 16.86,14.24 12.1,18.55M16.5,3C14.76,3 13.09,3.81 12,5.08C10.91,3.81 9.24,3 7.5,3C4.42,3 2,5.41 2,8.5C2,12.27 5.4,15.36 10.55,20.03L12,21.35L13.45,20.03C18.6,15.36 22,12.27 22,8.5C22,5.41 19.58,3 16.5,3Z"
                                    />
                                </svg>
                            {/if}
                        </button>
                    </div>
                </td>
            </tr>
        {/each}
    </tbody>
</table>
