<script>
    import { musicPlayer, musicState } from '../stores.js';
    import { formatDuration } from '../filters.js';

    export let token;
    export let tracks;
    export let showAlbum = true;
    export let isMusicQueue = false;

    function playTrack(track) {
        $musicPlayer.playTracks(tracks.slice(), track);
    }

    export function playFirstTrack() {
        playTrack(tracks[0]);
    }

    function likeTrack(track) {
        fetch(`${import.meta.env.VITE_API_URL}/tracks/${track.id}/like${track.liked ? '/delete' : ''}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        track.liked = !track.liked;
        tracks = tracks;
    }
</script>

<table class="table" style="width: 100%; table-layout: fixed;">
    <thead>
        {#if showAlbum}
            <th style="width: 10%;"><div class="track-index">#</div></th>
            <th style="width: calc(64px + 1.5em);">Title</th>
            <th class="track-title" />
            <th style="width: 30%;" class="is-hidden-mobile">Album</th>
            <th style="width: 15%;">Duration</th>
            <th style="width: 15%;" class="is-hidden-mobile">Plays</th>
            <th style="width: calc(40px + .75em);" />
            <th style="width: calc(40px + .75em);" class:is-hidden-mobile={!isMusicQueue} />
        {:else}
            <th style="width: 10%;"><div class="track-index">#</div></th>
            <th style="width: 50%;">Title</th>
            <th style="width: 20%;">Duration</th>
            <th class="is-hidden-mobile">Plays</th>
            <th style="width: calc(40px + .75em);" />
            <th style="width: calc(40px + .75em);" class:is-hidden-mobile={!isMusicQueue} />
        {/if}
    </thead>
    <tbody>
        {#each tracks as track, index}
            <tr
                class="track-container"
                on:dblclick|preventDefault={() => playTrack(track)}
                class:has-background-light={$musicState.track != undefined && $musicState.track.id == track.id}
            >
                <td>
                    <div class="track-index">{index + 1}</div>
                    <button class="button is-small track-play" on:click={() => playTrack(track)}>
                        <svg class="icon" viewBox="0 0 24 24">
                            <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                        </svg>
                    </button>
                </td>
                {#if showAlbum}
                    <td>
                        <div class="box has-image m-0 p-0" style="width: 64px; height: 64px;">
                            <img src={track.album.small_cover} alt="Cover of album {track.album}" loading="lazy" />
                        </div>
                    </td>
                {/if}
                <td>
                    <p class="ellipsis mb-1" style="font-weight: 500;">
                        <a href="/albums/{track.album.id}">{track.title}</a>
                    </p>
                    <p class="ellipsis">
                        {#if track.explicit}
                            <span class="tag is-danger mr-1">E</span>
                        {/if}
                        {#each track.artists as artist}
                            <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
                        {/each}
                    </p>
                </td>
                {#if showAlbum}
                    <td class="ellipsis is-hidden-mobile"><a href="/albums/{track.album.id}">{track.album.title}</a></td
                    >
                {/if}
                <td>{formatDuration(track.duration)}</td>
                <td class="is-hidden-mobile">{track.plays}</td>
                <td class="px-0">
                    <button class="button" on:click={() => likeTrack(track)}>
                        {#if track.liked}
                            <svg class="icon is-colored" viewBox="0 0 24 24">
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
                </td>
                <td class="pl-0" class:is-hidden-mobile={!isMusicQueue}>
                    {#if isMusicQueue}
                        <button
                            class="button"
                            on:click={() => $musicPlayer.removeTrack(track)}
                            disabled={$musicState.track != undefined && $musicState.track.id == track.id}
                        >
                            <svg class="icon" viewBox="0 0 24 24">
                                <path
                                    d="M14 10H3V12H14V10M14 6H3V8H14V6M3 16H10V14H3V16M14.4 22L17 19.4L19.6 22L21 20.6L18.4 18L21 15.4L19.6 14L17 16.6L14.4 14L13 15.4L15.6 18L13 20.6L14.4 22Z"
                                />
                            </svg>
                        </button>
                    {:else}
                        <button class="button" on:click={() => $musicPlayer.addTrack(track)}>
                            <svg class="icon" viewBox="0 0 24 24">
                                <path
                                    d="M3 16H10V14H3M18 14V10H16V14H12V16H16V20H18V16H22V14M14 6H3V8H14M14 10H3V12H14V10Z"
                                />
                            </svg>
                        </button>
                    {/if}
                </td>
            </tr>
        {/each}
    </tbody>
</table>

<style>
    .track-title {
        width: 50%;
    }
    @media (min-width: 768px) {
        .track-title {
            width: 30%;
        }
    }

    .track-index {
        width: 30px;
        text-align: center;
    }
    .track-play,
    .track-container:hover .track-index {
        display: none;
    }
    .track-index,
    .track-container:hover .track-play {
        display: block;
    }
</style>
