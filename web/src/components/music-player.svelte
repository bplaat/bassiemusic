<script>
    import { page } from '$app/stores';
    import { musicState } from '../stores.js';
    import {
        WEBSOCKET_RECONNECT_TIMEOUT,
        PLAYER_PREVIOUS_RESET_TIMEOUT,
        PLAYER_UPDATE_UI_TIMEOUT,
        PLAYER_UPDATE_SERVER_TIMEOUT,
        PLAYER_SEEK_TIME,
    } from '../consts.js';
    import { formatDuration } from '../filters.js';
    import { onMount, onDestroy } from 'svelte';

    // Utils
    function rand(min, max) {
        return Math.floor(Math.random() * (max - min + 1)) + min;
    }

    // Interface
    export let token = undefined;
    export let queue = [];
    export let track = undefined;
    export let position = 0;
    export let duration = 0;

    export function playTracks(_queue, _track) {
        queue = _queue;
        track = _track;
        position = 0;
        musicState.set({ queue, track });
        loadAndPlayTrack(true);
    }
    export function addTrack(track) {
        queue = [...queue, track];
        musicState.update((musicState) => {
            musicState.queue = queue;
            return musicState;
        });
    }
    export function removeTrack(track) {
        queue = queue.filter((otherTrack) => otherTrack.id != track.id);
        musicState.update((musicState) => {
            musicState.queue = queue;
            return musicState;
        });
    }

    // Websocket connection
    let ws;
    let connected = false;
    function websocketConnect() {
        ws = new WebSocket(import.meta.env.VITE_WEBSOCKET_URL);
        ws.onopen = () => {
            connected = true;
            ws.send(JSON.stringify({ type: 'auth', token }));
        };
        ws.disconnect = () => {
            connected = false;
            setTimeout(websocketConnect, WEBSOCKET_RECONNECT_TIMEOUT);
        };
    }
    onMount(websocketConnect);
    onDestroy(() => {
        if (!connected) return;
        ws.close();
    });

    // Music player
    let isPlaying = false,
        audio,
        updateUiTimeout,
        updateServerTimeout,
        isShuffling = false;

    function loadAndPlayTrack(autoplay) {
        if (audio != undefined) {
            audio.pause();
        }
        if (updateUiTimeout != undefined) {
            clearTimeout(updateUiTimeout);
        }
        if (updateServerTimeout != undefined) {
            clearTimeout(updateServerTimeout);
        }

        document.querySelector('.app').classList.add('is-playing');

        audio = new Audio(track.music);
        audio.volume = volume;
        audio.onloadedmetadata = () => {
            audio.currentTime = position;
            duration = audio.duration;
            if (autoplay) play();
        };
        audio.onratechange = () => {
            updatePlaybackState();
        };
        audio.onended = nextTrack;

        if ('mediaSession' in navigator) {
            navigator.mediaSession.metadata = new MediaMetadata({
                title: track.title,
                artist: track.artists.map((artist) => artist.name).join(', '),
                album: track.album.title,
                artwork: [
                    {
                        type: 'image/jpeg',
                        src: track.album.small_cover,
                        sizes: '256x256',
                    },
                    {
                        type: 'image/jpeg',
                        src: track.album.medium_cover,
                        sizes: '512x512',
                    },
                    {
                        type: 'image/jpeg',
                        src: track.album.large_cover,
                        sizes: '1024x1024',
                    },
                ],
            });
        }
    }
    onMount(() => {
        musicState.set({ queue, track });

        isShuffling = localStorage.getItem('player-shuffling') == 'true';

        if (track != undefined) {
            loadAndPlayTrack(false);
        }
    });

    function setIsShuffling(newIsShuffling) {
        isShuffling = newIsShuffling;
        localStorage.setItem('player-shuffling', newIsShuffling);
    }

    function updatePlaybackState() {
        position = audio.currentTime;
        if ('mediaSession' in navigator && audio.readyState >= 1) {
            navigator.mediaSession.setPositionState({
                duration: audio.duration,
                playbackRate: audio.playbackRate,
                position: audio.currentTime,
            });
        }
    }
    function updateUiLoop() {
        position = audio.currentTime;
        if (isPlaying) {
            updateUiTimeout = setTimeout(updateUiLoop, PLAYER_UPDATE_UI_TIMEOUT);
        }
    }

    function sendTrackPlay() {
        if (!connected) return;
        ws.send(JSON.stringify({ type: 'track_play', track_id: track.id, position: audio.currentTime }));
    }
    async function updateServerLoop() {
        sendTrackPlay();
        if (isPlaying) {
            updateServerTimeout = setTimeout(updateServerLoop, PLAYER_UPDATE_SERVER_TIMEOUT);
        }
    }

    function seekTo({ seekTime }) {
        if (!isPlaying) play();
        audio.currentTime = seekTime;
        updatePlaybackState();
        sendTrackPlay();
    }

    function previousTrack() {
        if (audio.currentTime * 1000 > PLAYER_PREVIOUS_RESET_TIMEOUT) {
            seekTo({ seekTime: 0 });
        } else {
            if (isShuffling) {
                track = queue[rand(0, queue.length - 1)];
            } else {
                const index = queue.indexOf(track);
                track = queue[index - 1 >= 0 ? index - 1 : queue.length - 1];
            }
            musicState.update((musicState) => {
                musicState.track = track;
                return musicState;
            });
            position = 0;
            loadAndPlayTrack(true);
        }
    }

    function seekBackward(details) {
        if (!isPlaying) play();
        audio.currentTime = Math.max(0, audio.currentTime - (details.seekOffset || PLAYER_SEEK_TIME));
        updatePlaybackState();
        sendTrackPlay();
    }

    function play() {
        audio.play();
        if ('mediaSession' in navigator) {
            navigator.mediaSession.playbackState = 'playing';
        }
        isPlaying = true;
        updateUiLoop();
        updateServerLoop();
    }

    function pause() {
        audio.pause();
        if ('mediaSession' in navigator) {
            navigator.mediaSession.playbackState = 'paused';
        }
        isPlaying = false;
    }

    function playPause() {
        if (isPlaying) {
            pause();
        } else {
            play();
        }
    }

    function seekForward(details) {
        if (!isPlaying) play();
        audio.currentTime = Math.min(audio.duration, audio.currentTime + (details.seekOffset || PLAYER_SEEK_TIME));
        updatePlaybackState();
        sendTrackPlay();
    }

    function nextTrack() {
        if (isShuffling) {
            track = queue[rand(0, queue.length - 1)];
        } else {
            const index = queue.indexOf(track);
            track = queue[index + 1 <= queue.length - 1 ? index + 1 : 0];
        }
        musicState.update((musicState) => {
            musicState.track = track;
            return musicState;
        });
        position = 0;
        loadAndPlayTrack(true);
    }

    onMount(() => {
        if ('mediaSession' in navigator) {
            navigator.mediaSession.setActionHandler('play', play);
            navigator.mediaSession.setActionHandler('pause', pause);
            navigator.mediaSession.setActionHandler('stop', pause);
            navigator.mediaSession.setActionHandler('seekbackward', seekBackward);
            navigator.mediaSession.setActionHandler('seekforward', seekForward);
            navigator.mediaSession.setActionHandler('seekto', seekTo);
            navigator.mediaSession.setActionHandler('previoustrack', previousTrack);
            navigator.mediaSession.setActionHandler('nexttrack', nextTrack);
        }
    });

    // Like
    function likeTrack() {
        fetch(`${import.meta.env.VITE_API_URL}/tracks/${track.id}/like${track.liked ? '/delete' : ''}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        track.liked = !track.liked;
    }

    // Volume
    let volume = 1;
    onMount(() => {
        volume = localStorage.getItem('player-volume') ?? 1;
    });
    function setVolume(newVolume) {
        localStorage.setItem('player-volume', newVolume);
        volume = newVolume;
        if (audio != undefined) {
            audio.volume = newVolume;
        }
    }

    let oldVolume = undefined;
    function toggleVolume() {
        if (volume > 0) {
            oldVolume = volume;
            setVolume(0);
        } else {
            if (oldVolume != undefined) {
                setVolume(oldVolume);
                oldVolume = undefined;
            } else {
                setVolume(1);
            }
        }
    }
</script>

{#if track != undefined}
    <div class="music-player box m-0 p-0 py-2 has-background-white-bis">
        <div class="media px-4 py-2 is-hidden-desktop">
            <div class="media-left">
                <div class="box m-0 p-0" style="width: 48px; height: 48px;">
                    <img src={track.album.small_cover} alt="Cover of album {track.album}" loading="lazy" />
                </div>
            </div>
            <div class="media-content">
                <p class="ellipsis">
                    <a href="/albums/{track.album.id}" style="font-weight: 500;">{track.title}</a>
                </p>
                <p class="ellipsis">
                    {#each track.artists as artist}
                        <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
                    {/each}
                </p>
            </div>

            <button class="button ml-3" on:click={likeTrack}>
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
        </div>

        <div class="box m-0 mx-3 p-0 is-hidden-touch" style="width: 64px; height: 64px;">
            <img src={track.album.small_cover} alt="Cover of album {track.album}" loading="lazy" />
        </div>

        <div style="width: 10rem;" class="is-hidden-touch">
            <p class="ellipsis">
                <a href="/albums/{track.album.id}" style="font-weight: 500;">{track.title}</a>
            </p>
            <p class="ellipsis">
                {#each track.artists as artist}
                    <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
                {/each}
            </p>
        </div>

        <button class="button is-hidden-touch ml-4" on:click={likeTrack}>
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

        <button
            class="button is-hidden-touch ml-4"
            class:is-link={isShuffling}
            on:click={() => setIsShuffling(!isShuffling)}
        >
            <svg class="icon" viewBox="0 0 24 24">
                <path
                    d="M14.83,13.41L13.42,14.82L16.55,17.95L14.5,20H20V14.5L17.96,16.54L14.83,13.41M14.5,4L16.54,6.04L4,18.59L5.41,20L17.96,7.46L20,9.5V4M10.59,9.17L5.41,4L4,5.41L9.17,10.58L10.59,9.17Z"
                />
            </svg>
        </button>

        <div class="music-player-controls px-4 py-2">
            <span class="is-hidden-desktop" style="width: 2.5rem;">{formatDuration(position)}</span>
            <div class="flex" />

            <div class="buttons has-addons m-0">
                <button class="button m-0" on:click={previousTrack}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path d="M6,18V6H8V18H6M9.5,12L18,6V18L9.5,12Z" />
                    </svg>
                </button>
                <button class="button m-0" on:click={seekBackward}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path d="M11.5,12L20,18V6M11,18V6L2.5,12L11,18Z" />
                    </svg>
                </button>
                <button class="button m-0" on:click={playPause}>
                    <svg class="icon" viewBox="0 0 24 24">
                        {#if isPlaying}
                            <path d="M14,19H18V5H14M6,19H10V5H6V19Z" />
                        {:else}
                            <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                        {/if}
                    </svg>
                </button>
                <button class="button m-0" on:click={seekForward}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path d="M13,6V18L21.5,12M4,18L12.5,12L4,6V18Z" />
                    </svg>
                </button>
                <button class="button m-0" on:click={nextTrack}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path d="M16,18H18V6H16M6,18L14.5,12L6,6V18Z" />
                    </svg>
                </button>
            </div>

            <div class="flex" />
            <span class="is-hidden-desktop" style="width: 2.5rem; text-align: right;"
                >-{formatDuration(duration - position)}</span
            >
        </div>

        <div class="music-player-slider px-4">
            <span class="is-hidden-touch mr-4" style="width: 3rem; text-align: right;">{formatDuration(position)}</span>
            <input
                class="range"
                type="range"
                bind:value={position}
                on:input={(event) => seekTo({ seekTime: event.target.value })}
                max={duration}
                style="width: 100%; background-size: {(position / duration) * 100}% 100%;"
            />
            <span class="is-hidden-touch ml-4" style="width: 3rem;">-{formatDuration(duration - position)}</span>
        </div>

        <div class="music-player-volume px-4 is-hidden-touch">
            {#if $page.url.pathname == '/queue'}
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a class="button mr-4 is-link" href="#" on:click|preventDefault={() => history.back()}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path
                            d="M15,6H3V8H15V6M15,10H3V12H15V10M3,16H11V14H3V16M17,6V14.18C16.69,14.07 16.35,14 16,14A3,3 0 0,0 13,17A3,3 0 0,0 16,20A3,3 0 0,0 19,17V8H22V6H17Z"
                        />
                    </svg>
                </a>
            {:else}
                <a class="button mr-4" href="/queue">
                    <svg class="icon" viewBox="0 0 24 24">
                        <path
                            d="M15,6H3V8H15V6M15,10H3V12H15V10M3,16H11V14H3V16M17,6V14.18C16.69,14.07 16.35,14 16,14A3,3 0 0,0 13,17A3,3 0 0,0 16,20A3,3 0 0,0 19,17V8H22V6H17Z"
                        />
                    </svg>
                </a>
            {/if}

            <button class="button mr-3" on:click={toggleVolume}>
                <svg class="icon" viewBox="0 0 24 24">
                    {#if volume == 0}
                        <path
                            d="M12,4L9.91,6.09L12,8.18M4.27,3L3,4.27L7.73,9H3V15H7L12,20V13.27L16.25,17.53C15.58,18.04 14.83,18.46 14,18.7V20.77C15.38,20.45 16.63,19.82 17.68,18.96L19.73,21L21,19.73L12,10.73M19,12C19,12.94 18.8,13.82 18.46,14.64L19.97,16.15C20.62,14.91 21,13.5 21,12C21,7.72 18,4.14 14,3.23V5.29C16.89,6.15 19,8.83 19,12M16.5,12C16.5,10.23 15.5,8.71 14,7.97V10.18L16.45,12.63C16.5,12.43 16.5,12.21 16.5,12Z"
                        />
                    {/if}
                    {#if volume > 0 && volume < 0.33}
                        <path d="M7,9V15H11L16,20V4L11,9H7Z" />
                    {/if}
                    {#if volume >= 0.33 && volume < 0.67}
                        <path
                            d="M5,9V15H9L14,20V4L9,9M18.5,12C18.5,10.23 17.5,8.71 16,7.97V16C17.5,15.29 18.5,13.76 18.5,12Z"
                        />
                    {/if}
                    {#if volume >= 0.67}
                        <path
                            d="M14,3.23V5.29C16.89,6.15 19,8.83 19,12C19,15.17 16.89,17.84 14,18.7V20.77C18,19.86 21,16.28 21,12C21,7.72 18,4.14 14,3.23M16.5,12C16.5,10.23 15.5,8.71 14,7.97V16C15.5,15.29 16.5,13.76 16.5,12M3,9V15H7L12,20V4L7,9H3Z"
                        />
                    {/if}
                </svg>
            </button>

            <input
                class="range"
                type="range"
                bind:value={volume}
                on:input={(event) => setVolume(event.target.value)}
                max="1"
                step="0.01"
                style="width: 8rem; background-size: {volume * 100}% 100%;"
            />
        </div>
    </div>
{/if}

<style>
    .music-player {
        position: fixed;
        left: 0;
        bottom: 0;
        width: 100%;
        height: 10rem;
        z-index: 100;
        border-radius: 0;
        display: flex;
        flex-direction: column;
    }

    .music-player-controls,
    .music-player-volume {
        display: flex;
        align-items: center;
    }

    @media (min-width: 1024px) {
        .music-player {
            flex-direction: row;
            align-items: center;
            height: 6rem;
            z-index: 300;
        }

        .music-player-slider {
            display: flex;
            align-items: center;
            flex: 1;
            padding: 0 !important;
        }
        .music-player-slider > .range {
            flex: 1;
        }
    }

    /* Music player input range */
    .range {
        -webkit-appearance: none;
        background-color: #ccc;
        background-image: linear-gradient(#2196f3, #2196f3);
        background-repeat: no-repeat;
        height: 4px;
        border-radius: 4px;
    }
    .range:focus {
        outline: none;
    }
    .range::-webkit-slider-runnable-track {
        -webkit-appearance: none;
        border: none;
        background-color: transparent;
    }
    .range::-moz-range-track {
        border: none;
        background-color: transparent;
    }
    .range::-webkit-slider-thumb {
        -webkit-appearance: none;
        height: 16px;
        width: 16px;
        background-color: #999;
        border-radius: 50%;
    }
    .range::-moz-range-thumb {
        width: 16px;
        height: 16px;
        background-color: #999;
        border-radius: 50%;
        border: 0;
    }
    @media (prefers-color-scheme: dark) {
        .range::-webkit-slider-thumb {
            background-color: #fff;
        }
        .range::-moz-range-thumb {
            background-color: #fff;
        }
    }

    :global(.app.is-light .range::-webkit-slider-thumb) {
        background-color: #999 !important;
    }
    :global(.app.is-light .range::-moz-range-thumb) {
        background-color: #999 !important;
    }

    :global(.app.is-dark .range::-webkit-slider-thumb) {
        background-color: #fff !important;
    }
    :global(.app.is-dark .range::-moz-range-thumb) {
        background-color: #fff !important;
    }
</style>
