<script>
    import { browser } from "$app/environment";
    import {
        PLAYER_UPDATE_UI_TIMEOUT,
        PLAYER_UPDATE_SERVER_TIMEOUT,
        PLAYER_SEEK_TIME,
    } from "../consts.js";
    import { musicPlayer, audioVolume } from "../stores.js";
    import { formatDuration } from "../filters.js";
    import Slider from "./slider.svelte";

    export let token;
    let isPlaying = false,
        audio,
        audioDuration,
        audioCurrentTime,
        updateUiTimeout,
        updateServerTimeout,
        musicSlider,
        volumeSlider,
        ignoreMusicPlayerUpdate = false;

    $: track = $musicPlayer.queue[$musicPlayer.index];

    if (browser) {
        musicPlayer.subscribe((musicPlayer) => {
            if (musicPlayer.queue.length == 0) return;
            if (ignoreMusicPlayerUpdate) return;
            const track = musicPlayer.queue[musicPlayer.index];

            if (audio != undefined) {
                audio.pause();
            }
            if (updateUiTimeout != undefined) {
                clearTimeout(updateUiTimeout);
            }
            if (updateServerTimeout != undefined) {
                clearTimeout(updateServerTimeout);
            }

            document.body.classList.add("is-playing");

            audio = new Audio(track.music);
            audio.volume = $audioVolume;
            audio.onloadedmetadata = () => {
                audio.currentTime =
                    musicPlayer.action == "init" ? musicPlayer.position : 0;
                audioDuration = audio.duration;
                audioCurrentTime = audio.currentTime;

                //will crash when hot save ¯\_(ツ)_/¯
                try {
                    musicSlider.seekToValue(audioCurrentTime, audioDuration)
                    volumeSlider.seekToValue($audioVolume)
                } catch (e) {}

                if (musicPlayer.action == "play") {
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
        if ("mediaSession" in navigator) {
            navigator.mediaSession.metadata = new MediaMetadata({
                title: track.title,
                artist: track.artists.map((artist) => artist.name).join(", "),
                album: track.album.title,
                artwork: [
                    {
                        type: "image/jpeg",
                        src: track.album.large_cover,
                        sizes: "1024x1024",
                    },
                ],
            });

            navigator.mediaSession.setActionHandler("play", play);
            navigator.mediaSession.setActionHandler("pause", pause);
            navigator.mediaSession.setActionHandler("stop", pause);
            navigator.mediaSession.setActionHandler(
                "seekbackward",
                seekBackward
            );
            navigator.mediaSession.setActionHandler("seekforward", seekForward);
            navigator.mediaSession.setActionHandler("seekto", seekTo);
            navigator.mediaSession.setActionHandler(
                "previoustrack",
                previousTrack
            );
            navigator.mediaSession.setActionHandler("nexttrack", nextTrack);
        }
    }

    function updatePositionState() {
        audioCurrentTime = audio.currentTime;
        musicSlider.seekToValue(audioCurrentTime, audioDuration)

        if ("mediaSession" in navigator) {
            navigator.mediaSession.setPositionState({
                duration: audio.duration,
                playbackRate: audio.playbackRate,
                position: audio.currentTime,
            });
        }
    }

    function updateUiLoop() {
        audioCurrentTime = audio.currentTime;
        musicSlider.seekToValue(audioCurrentTime, audioDuration)

        if (isPlaying) {
            updateUiTimeout = setTimeout(
                updateUiLoop,
                PLAYER_UPDATE_UI_TIMEOUT
            );
        }
    }

    let isSendingTrackPlay = false;

    async function sendTrackPlay() {
        if (isSendingTrackPlay) return;
        isSendingTrackPlay = true;
        await fetch(
            `${import.meta.env.VITE_API_URL}/tracks/${
                track.id
            }/play?${new URLSearchParams({
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
            updateServerTimeout = setTimeout(
                updateServerLoop,
                PLAYER_UPDATE_SERVER_TIMEOUT
            );
        }
    }

    function seekTo(event) {
        if (!isPlaying) play();
        audio.currentTime = event.target.value;;
        sendTrackPlay();
        updatePositionState();
    }

    function sliderSeek(event){
        if (!isPlaying) play();
        audio.currentTime = event.detail.value;
        sendTrackPlay();
        updatePositionState();
    }

    function previousTrack() {
        musicPlayer.update((musicPlayer) => {
            musicPlayer.index =
                musicPlayer.index - 1 >= 0
                    ? musicPlayer.index - 1
                    : musicPlayer.queue.length - 1;
            return musicPlayer;
        });
    }

    function seekBackward(details) {
        if (!isPlaying) play();
        audio.currentTime = Math.max(
            0,
            audio.currentTime - (details.seekOffset || PLAYER_SEEK_TIME)
        );
        sendTrackPlay();
        updatePositionState();
    }

    function play() {
        ignoreMusicPlayerUpdate = true;
        musicPlayer.update((musicPlayer) => {
            musicPlayer.action = "play";
            return musicPlayer;
        });
        ignoreMusicPlayerUpdate = false;

        audio.play();
        if ("mediaSession" in navigator) {
            navigator.mediaSession.playbackState = "playing";
        }
        isPlaying = true;
        updateUiLoop();
        updateServerLoop();
        updatePositionState();
    }

    function pause() {
        audio.pause();
        if ("mediaSession" in navigator) {
            navigator.mediaSession.playbackState = "paused";
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
        audio.currentTime = Math.min(
            audio.duration,
            audio.currentTime + (details.seekOffset || PLAYER_SEEK_TIME)
        );
        sendTrackPlay();
        updatePositionState();
    }

    function nextTrack() {
        musicPlayer.update((musicPlayer) => {
            musicPlayer.index =
                musicPlayer.index + 1 <= musicPlayer.queue.length - 1
                    ? musicPlayer.index + 1
                    : 0;
            return musicPlayer;
        });
    }

    // Like
    function likeTrack() {
        const track = $musicPlayer.queue[$musicPlayer.index];
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
        ignoreMusicPlayerUpdate = true;
        track.liked = !track.liked;
        $musicPlayer = $musicPlayer;
        ignoreMusicPlayerUpdate = false;
    }

    // Volume
    let oldAudioVolume;

    audioVolume.subscribe((audioVolume) => {
        if (audio != undefined) {
            audio.volume = audioVolume;
        }
    });

    function volumeSeek(event){
        audioVolume.set(event.detail.value);
    }

    function toggleVolume() {
        if ($audioVolume > 0) {
            oldAudioVolume = $audioVolume;
            audioVolume.set(0);
            volumeSlider.seekToValue(0)
        } else {
            if (oldAudioVolume != undefined) {
                audioVolume.set(oldAudioVolume);
                volumeSlider.seekToValue(oldAudioVolume)
                oldAudioVolume = undefined;
            } else {
                audioVolume.set(1);
                volumeSlider.seekToValue(1)
            }
        }
    }
</script>

{#if $musicPlayer.queue.length > 0}
    <div class="player-controls box has-background-white-bis m-0">
        <div style="display: flex; align-items: center;">
            <div
                class="box is-image mr-4 mb-0"
                style="width: 64px; height: 64px; min-width: 64px; background-image: url({track
                    .album.small_cover});"
            />

            <div class="mr-5" style="width: 10rem">
                <p class="ellipsis">
                    <a href="/albums/{track.album.id}" style="font-weight: 500;"
                        >{track.title}</a
                    >
                </p>
                <p class="ellipsis">
                    {#each track.artists as artist}
                        <a href="/artists/{artist.id}" class="mr-2"
                            >{artist.name}</a
                        >
                    {/each}
                </p>
            </div>

            <button class="button mr-3" on:click={likeTrack}>
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

            <div class="field has-addons mb-0">
                <p class="control">
                    <button class="button" on:click={previousTrack}>
                        <svg class="icon" viewBox="0 0 24 24">
                            <path d="M6,18V6H8V18H6M9.5,12L18,6V18L9.5,12Z" />
                        </svg>
                    </button>
                </p>
                <p class="control">
                    <button class="button" on:click={seekBackward}>
                        <svg class="icon" viewBox="0 0 24 24">
                            <path d="M11.5,12L20,18V6M11,18V6L2.5,12L11,18Z" />
                        </svg>
                    </button>
                </p>
                <p class="control">
                    <button class="button" on:click={playPause}>
                        <svg class="icon" viewBox="0 0 24 24">
                            {#if isPlaying}
                                <path d="M14,19H18V5H14M6,19H10V5H6V19Z" />
                            {:else}
                                <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                            {/if}
                        </svg>
                    </button>
                </p>
                <p class="control">
                    <button class="button" on:click={seekForward}>
                        <svg class="icon" viewBox="0 0 24 24">
                            <path d="M13,6V18L21.5,12M4,18L12.5,12L4,6V18Z" />
                        </svg>
                    </button>
                </p>
                <p class="control">
                    <button class="button" on:click={nextTrack}>
                        <svg class="icon" viewBox="0 0 24 24">
                            <path d="M16,18H18V6H16M6,18L14.5,12L6,6V18Z" />
                        </svg>
                    </button>
                </p>
            </div>

            <div style="flex: 4; display: flex;">
                <span class="mr-3" style="width: 4rem; text-align: right;"
                    >{formatDuration(audioCurrentTime)}</span
                >
                <Slider maxValue={audioDuration} bind:this={musicSlider} on:newValue={sliderSeek} />
                <span class="ml-3" style="width: 4rem;"
                    >-{formatDuration(audioDuration - audioCurrentTime)}</span
                >
            </div>

            <div style="flex: 1; display: flex;">
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
                <div style="margin-top: 8px; flex: 1;">
                    <Slider bind:this={volumeSlider} on:newValue={volumeSeek} maxValue="1" />
                </div>
            </div>
        </div>
    </div>
{/if}
