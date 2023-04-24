<script>
    import { page } from '$app/stores';
    import { onMount, onDestroy } from 'svelte';
    import LikeButton from './buttons/like-button.svelte';
    import { musicState, language } from '../stores.js';
    import {
        WEBSOCKET_RECONNECT_TIMEOUT,
        PLAYER_PREVIOUS_RESET_TIMEOUT,
        PLAYER_UPDATE_UI_TIMEOUT,
        PLAYER_UPDATE_SERVER_TIMEOUT,
        PLAYER_SEEK_TIME,
    } from '../consts.js';
    import { formatDuration } from '../filters.js';
    import { rand } from '../utils.js';

    // Language strings
    const lang = {
        en: {
            cover_alt: 'Cover of album $1',
            track: 'track',
            shuffle: 'Start shuffling play queue',
            stop_shuffle: 'Stop shuffling play queue',
            previous: 'Previous',
            seek_backward: 'Seek backward',
            play: 'Play',
            pause: 'Pause',
            seek_forward: 'Seek forward',
            next: 'Next',
            open_queue: 'Open play queue',
            close_queue: 'Close play queue',
            mute_volume: 'Mute volume',
            restore_volume: 'Restore volume',
        },
        nl: {
            cover_alt: 'Hoes van album $1',
            track: 'track',
            shuffle: 'Start willekeurige wachtrij',
            stop_shuffle: 'Stop willekeurige wachtrij',
            previous: 'Vorige',
            seek_backward: 'Spoel terug',
            play: 'Speel',
            pause: 'Pauze',
            seek_forward: 'Spel naar voren',
            next: 'Volgende',
            open_queue: 'Open wachtrij',
            close_queue: 'Sluit wachtrij',
            mute_volume: 'Demp volume',
            restore_volume: 'Herstel volume',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // Props
    export let token;
    export let queue = [];
    export let track = null;
    export let position = 0;
    export let duration = 0;

    // Methods
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
        queue = queue.filter((otherTrack) => otherTrack.id !== track.id);
        musicState.update((musicState) => {
            musicState.queue = queue;
            return musicState;
        });
    }

    // Websocket connection
    let ws;
    let connecting = false;
    let connected = false;
    function websocketConnect() {
        connecting = true;
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
    onDestroy(() => {
        if (!connected) return;
        ws.close();
    });

    // Music player
    let isPlaying = false,
        audio = null,
        updateUiTimeout = null,
        updateServerTimeout = null,
        isShuffling = false;

    function loadAndPlayTrack(autoplay) {
        if (audio !== null) {
            audio.pause();
        }
        if (updateUiTimeout !== null) {
            clearTimeout(updateUiTimeout);
        }
        if (updateServerTimeout !== null) {
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
                artwork:
                    track.album.small_cover !== null
                        ? [
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
                          ]
                        : [
                              {
                                  type: 'image/svg+xml',
                                  src: '/images/album-default.svg',
                                  sizes: '1024x1024',
                              },
                          ],
            });
        }
    }
    onMount(() => {
        musicState.set({ queue, track });

        isShuffling = localStorage.getItem('player-shuffling') === 'true';

        if (track !== null) {
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
        if (!connected && !connecting) {
            websocketConnect();
            return;
        }
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

    // Volume
    let volume = 1;
    onMount(() => {
        volume = localStorage.getItem('player-volume') ?? 1;
    });
    function setVolume(newVolume) {
        localStorage.setItem('player-volume', newVolume);
        volume = newVolume;
        if (audio !== null) {
            audio.volume = newVolume;
        }
    }

    let oldVolume = null;
    function toggleVolume() {
        if (volume > 0) {
            oldVolume = volume;
            setVolume(0);
        } else {
            if (oldVolume !== null) {
                setVolume(oldVolume);
                oldVolume = null;
            } else {
                setVolume(1);
            }
        }
    }
</script>

{#if track !== null}
    <div class="music-player box m-0 p-0 py-2 has-background-white-bis">
        <div class="media px-4 py-2">
            <div class="media-left">
                <a href="/albums/{track.album.id}" class="music-player-album-cover box m-0 p-0">
                    <img
                        src={track.album.small_cover || '/images/album-default.svg'}
                        alt={t('cover_alt', track.album)}
                    />
                </a>
            </div>
            <div class="media-content" style="width: 10rem; min-width: 0;">
                <p class="ellipsis">
                    <a href="/albums/{track.album.id}#{track.disk}-{track.position}" style="font-weight: 500;">
                        {track.title}
                    </a>
                </p>
                <p class="ellipsis">
                    {#each track.artists as artist}
                        <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
                    {/each}
                </p>
            </div>
            <div class="media-right">
                <LikeButton {token} item={track} itemRoute="tracks" itemLabel={t('track')} />
            </div>
        </div>

        <button
            class="button is-hidden-touch"
            class:is-link={isShuffling}
            on:click={() => setIsShuffling(!isShuffling)}
            title={!isShuffling ? t('shuffle') : t('stop_shuffle')}
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
                <button class="button m-0" on:click={previousTrack} title={t('previous')}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path d="M6,18V6H8V18H6M9.5,12L18,6V18L9.5,12Z" />
                    </svg>
                </button>
                <button class="button m-0" on:click={seekBackward} title={t('seek_backward')}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path d="M11.5,12L20,18V6M11,18V6L2.5,12L11,18Z" />
                    </svg>
                </button>
                <button class="button m-0" on:click={playPause} title={!isPlaying ? t('play') : t('pause')}>
                    <svg class="icon" viewBox="0 0 24 24">
                        {#if isPlaying}
                            <path d="M14,19H18V5H14M6,19H10V5H6V19Z" />
                        {:else}
                            <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                        {/if}
                    </svg>
                </button>
                <button class="button m-0" on:click={seekForward} title={t('seek_forward')}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path d="M13,6V18L21.5,12M4,18L12.5,12L4,6V18Z" />
                    </svg>
                </button>
                <button class="button m-0" on:click={nextTrack} title={t('next')}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path d="M16,18H18V6H16M6,18L14.5,12L6,6V18Z" />
                    </svg>
                </button>
            </div>

            <div class="flex" />
            <span class="is-hidden-desktop" style="width: 2.5rem; text-align: right;">
                -{formatDuration(duration - position)}
            </span>
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
            {#if $page.url.pathname !== '/queue'}
                <a class="button mr-4" href="/queue" title={t('open_queue')}>
                    <svg class="icon" viewBox="0 0 24 24">
                        <path
                            d="M15,6H3V8H15V6M15,10H3V12H15V10M3,16H11V14H3V16M17,6V14.18C16.69,14.07 16.35,14 16,14A3,3 0 0,0 13,17A3,3 0 0,0 16,20A3,3 0 0,0 19,17V8H22V6H17Z"
                        />
                    </svg>
                </a>
            {:else}
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a
                    class="button mr-4 is-link"
                    href="#"
                    on:click|preventDefault={() => history.back()}
                    title={t('close_queue')}
                >
                    <svg class="icon" viewBox="0 0 24 24">
                        <path
                            d="M15,6H3V8H15V6M15,10H3V12H15V10M3,16H11V14H3V16M17,6V14.18C16.69,14.07 16.35,14 16,14A3,3 0 0,0 13,17A3,3 0 0,0 16,20A3,3 0 0,0 19,17V8H22V6H17Z"
                        />
                    </svg>
                </a>
            {/if}

            <div class="dropdown is-up is-right is-hoverable">
                <div class="dropdown-trigger">
                    <button
                        class="button"
                        on:click={toggleVolume}
                        title={volume > 0 ? t('mute_volume') : t('restore_volume')}
                    >
                        <svg class="icon" viewBox="0 0 24 24">
                            {#if volume === 0}
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
                </div>
                <div class="dropdown-menu">
                    <div class="dropdown-content px-5 py-4">
                        <input
                            class="range"
                            type="range"
                            bind:value={volume}
                            on:input={(event) => setVolume(event.target.value)}
                            max="1"
                            step="0.01"
                            style="width: 100%; background-size: {volume * 100}% 100%;"
                        />
                    </div>
                </div>
            </div>
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
        overflow: visible;
        border-radius: 0;
        display: flex;
        flex-direction: column;
    }

    .music-player-album-cover {
        width: 48px;
        height: 48px;
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

        .music-player-album-cover {
            width: 64px;
            height: 64px;
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
        appearance: none;
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
