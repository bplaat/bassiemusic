<script>
    import { page } from '$app/stores';
    import { browser } from '$app/environment';
    import { PLAYER_UPDATE_UI_TIMEOUT, PLAYER_UPDATE_SERVER_TIMEOUT, PLAYER_SEEK_TIME } from '../consts.js';
    import { musicPlayer, audioVolume } from '../stores.js';
    import { formatDuration } from '../filters.js';
    import Slider from './slider.svelte';

    export let token;
    let playingMusicPlayerTrackId,
        isPlaying = false,
        audio,
        audioDuration,
        audioCurrentTime,
        updateUiTimeout,
        updateServerTimeout,
        musicSlider,
        volumeSlider;

    $: track = $musicPlayer.queue.find((track) => track.id == $musicPlayer.track_id);

    if (browser) {
        musicPlayer.subscribe((musicPlayer) => {
            if (musicPlayer.queue.length == 0) return;
            if (playingMusicPlayerTrackId == musicPlayer.track_id) return;
            const track = musicPlayer.queue.find((track) => track.id == musicPlayer.track_id);
            playingMusicPlayerTrackId = musicPlayer.track_id;

            if (audio != undefined) {
                audio.pause();
            }
            if (updateUiTimeout != undefined) {
                clearTimeout(updateUiTimeout);
            }
            if (updateServerTimeout != undefined) {
                clearTimeout(updateServerTimeout);
            }

            document.body.classList.add('is-playing');

            audio = new Audio(track.music);
            audio.volume = $audioVolume;
            audio.onloadedmetadata = () => {
                audio.currentTime = musicPlayer.action == 'init' ? musicPlayer.position : 0;
                audioDuration = audio.duration;
                audioCurrentTime = audio.currentTime;

                if (musicSlider) {
                    musicSlider.seekToValue(audioCurrentTime, audioDuration);
                }
                if (volumeSlider) {
                    volumeSlider.seekToValue($audioVolume);
                }

                if (musicPlayer.action == 'play') {
                    play();
                }
            };
            audio.onratechange = () => {
                updatePositionState();
            };
            audio.onended = nextTrack;

            setMediaSession(track);
        });
    }

    function setMediaSession(track) {
        if ('mediaSession' in navigator) {
            navigator.mediaSession.metadata = new MediaMetadata({
                title: track.title,
                artist: track.artists.map((artist) => artist.name).join(', '),
                album: track.album.title,
                artwork: [
                    {
                        type: 'image/jpeg',
                        src: track.album.large_cover,
                        sizes: '1024x1024',
                    },
                ],
            });

            navigator.mediaSession.setActionHandler('play', play);
            navigator.mediaSession.setActionHandler('pause', pause);
            navigator.mediaSession.setActionHandler('stop', pause);
            navigator.mediaSession.setActionHandler('seekbackward', seekBackward);
            navigator.mediaSession.setActionHandler('seekforward', seekForward);
            navigator.mediaSession.setActionHandler('seekto', seekTo);
            navigator.mediaSession.setActionHandler('previoustrack', previousTrack);
            navigator.mediaSession.setActionHandler('nexttrack', nextTrack);
        }
    }

    function updatePositionState() {
        audioCurrentTime = audio.currentTime;
        musicSlider.seekToValue(audioCurrentTime, audioDuration);
        if ('mediaSession' in navigator && audio.readyState >= 1) {
            navigator.mediaSession.setPositionState({
                duration: audio.duration,
                playbackRate: audio.playbackRate,
                position: audio.currentTime,
            });
        }
    }

    function updateUiLoop() {
        audioCurrentTime = audio.currentTime;
        musicSlider.seekToValue(audioCurrentTime, audioDuration);

        if (isPlaying) {
            updateUiTimeout = setTimeout(updateUiLoop, PLAYER_UPDATE_UI_TIMEOUT);
        }
    }

    let isSendingTrackPlay = false;

    async function sendTrackPlay() {
        if (isSendingTrackPlay) return;
        isSendingTrackPlay = true;
        await fetch(
            `${import.meta.env.VITE_API_URL}/tracks/${track.id}/play?${new URLSearchParams({
                position: audio.currentTime,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        isSendingTrackPlay = false;
    }

    async function updateServerLoop() {
        await sendTrackPlay();
        if (isPlaying) {
            updateServerTimeout = setTimeout(updateServerLoop, PLAYER_UPDATE_SERVER_TIMEOUT);
        }
    }

    function seekTo(event) {
        if (!isPlaying) play();
        audio.currentTime = event.detail.value;
        sendTrackPlay();
        updatePositionState();
    }

    function previousTrack() {
        playingMusicPlayerTrackId = undefined;
        musicPlayer.update((musicPlayer) => {
            const track = musicPlayer.queue.find((track) => track.id == musicPlayer.track_id);
            const index = musicPlayer.queue.indexOf(track);
            musicPlayer.track_id = musicPlayer.queue[index - 1 >= 0 ? index - 1 : musicPlayer.queue.length - 1].id;
            return musicPlayer;
        });
    }

    function seekBackward(details) {
        if (!isPlaying) play();
        audio.currentTime = Math.max(0, audio.currentTime - (details.seekOffset || PLAYER_SEEK_TIME));
        sendTrackPlay();
        updatePositionState();
    }

    function play() {
        musicPlayer.update((musicPlayer) => {
            musicPlayer.action = 'play';
            return musicPlayer;
        });

        audio.play();
        if ('mediaSession' in navigator) {
            navigator.mediaSession.playbackState = 'playing';
        }
        isPlaying = true;
        updateUiLoop();
        updateServerLoop();
        updatePositionState();
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
        sendTrackPlay();
        updatePositionState();
    }

    function nextTrack() {
        playingMusicPlayerTrackId = undefined;
        musicPlayer.update((musicPlayer) => {
            const track = musicPlayer.queue.find((track) => track.id == musicPlayer.track_id);
            const index = musicPlayer.queue.indexOf(track);
            musicPlayer.track_id = musicPlayer.queue[index + 1 <= musicPlayer.queue.length - 1 ? index + 1 : 0].id;
            return musicPlayer;
        });
    }

    // Like
    function likeTrack() {
        const track = $musicPlayer.queue.find((track) => track.id == $musicPlayer.track_id);
        fetch(`${import.meta.env.VITE_API_URL}/tracks/${track.id}/like${track.liked ? '/delete' : ''}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        track.liked = !track.liked;
        $musicPlayer = $musicPlayer;
    }

    // Volume
    let oldAudioVolume;

    audioVolume.subscribe((audioVolume) => {
        if (audio != undefined) {
            audio.volume = audioVolume;
        }
    });

    function volumeSeek(event) {
        audioVolume.set(event.detail.value);
    }

    function toggleVolume() {
        if ($audioVolume > 0) {
            oldAudioVolume = $audioVolume;
            audioVolume.set(0);
            volumeSlider.seekToValue(0);
        } else {
            if (oldAudioVolume != undefined) {
                audioVolume.set(oldAudioVolume);
                volumeSlider.seekToValue(oldAudioVolume);
                oldAudioVolume = undefined;
            } else {
                audioVolume.set(1);
                volumeSlider.seekToValue(1);
            }
        }
    }
</script>

{#if $musicPlayer.queue.length > 0}
    <div class="player-controls box m-0 p-0 pl-4 pr-5 has-background-white-bis">
        <div class="box m-0 p-0 mr-4" style="width: 64px; height: 64px;">
            <img src={track.album.small_cover} alt="Cover of album {track.album}" loading="lazy" />
        </div>

        <div class="mr-4" style="width: calc(13.5rem - 64px);">
            <p class="ellipsis">
                <a href="/albums/{track.album.id}" style="font-weight: 500;">{track.title}</a>
            </p>
            <p class="ellipsis">
                {#each track.artists as artist}
                    <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
                {/each}
            </p>
        </div>

        <button class="button mr-4" on:click={likeTrack}>
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

        <span class="mr-3" style="width: 4rem; text-align: right;">{formatDuration(audioCurrentTime)}</span>
        <Slider style="flex: 1;" maxValue={audioDuration} bind:this={musicSlider} on:newValue={seekTo} />
        <span class="ml-3" style="width: 4rem;">-{formatDuration(audioDuration - audioCurrentTime)}</span>

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
                {#if $audioVolume == 0}
                    <path
                        d="M12,4L9.91,6.09L12,8.18M4.27,3L3,4.27L7.73,9H3V15H7L12,20V13.27L16.25,17.53C15.58,18.04 14.83,18.46 14,18.7V20.77C15.38,20.45 16.63,19.82 17.68,18.96L19.73,21L21,19.73L12,10.73M19,12C19,12.94 18.8,13.82 18.46,14.64L19.97,16.15C20.62,14.91 21,13.5 21,12C21,7.72 18,4.14 14,3.23V5.29C16.89,6.15 19,8.83 19,12M16.5,12C16.5,10.23 15.5,8.71 14,7.97V10.18L16.45,12.63C16.5,12.43 16.5,12.21 16.5,12Z"
                    />
                {/if}
                {#if $audioVolume > 0 && $audioVolume < 0.33}
                    <path d="M7,9V15H11L16,20V4L11,9H7Z" />
                {/if}
                {#if $audioVolume >= 0.33 && $audioVolume < 0.67}
                    <path
                        d="M5,9V15H9L14,20V4L9,9M18.5,12C18.5,10.23 17.5,8.71 16,7.97V16C17.5,15.29 18.5,13.76 18.5,12Z"
                    />
                {/if}
                {#if $audioVolume >= 0.67}
                    <path
                        d="M14,3.23V5.29C16.89,6.15 19,8.83 19,12C19,15.17 16.89,17.84 14,18.7V20.77C18,19.86 21,16.28 21,12C21,7.72 18,4.14 14,3.23M16.5,12C16.5,10.23 15.5,8.71 14,7.97V16C15.5,15.29 16.5,13.76 16.5,12M3,9V15H7L12,20V4L7,9H3Z"
                    />
                {/if}
            </svg>
        </button>

        <Slider style="width: 8rem;" bind:this={volumeSlider} on:newValue={volumeSeek} maxValue="1" />
    </div>
{/if}
