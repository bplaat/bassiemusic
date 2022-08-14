<script>
import { API_URL, UPDATE_TIMEOUT, SEEK_TIME } from './config.js';

export default {
    data() {
        return {
            playingTrack: undefined,
            isPlaying: false,
            audio: undefined,
            audioCurrentTime: 0,
            audioDuration: 0,
            audioVolume: 1,
            oldAudioVolume: undefined
        };
    },

    created() {
        if (localStorage.getItem('volume') != null) {
            this.audioVolume = parseFloat(localStorage.getItem('volume'));
        }
        if (localStorage.getItem('oldVolume') != null) {
            this.oldAudioVolume = parseFloat(localStorage.getItem('oldVolume'));
        }
    },

    watch: {
        audioVolume(newAudioVolume) {
            if (this.audio != undefined) {
                this.audio.volume = newAudioVolume;
            }
            localStorage.setItem('volume', newAudioVolume);
        },

        oldAudioVolume(newOldAudioVolume) {
            if (newOldAudioVolume != undefined) {
                localStorage.setItem('oldVolume', newOldAudioVolume);
            } else {
                localStorage.removeItem('oldVolume');
            }
        }
    },

    methods: {
        playTrack(track) {
            this.playingTrack = track;

            fetch(`${API_URL}/tracks/${track.id}/play`);

            if (this.audio != undefined) {
                this.audio.pause();
            }

            this.audio = new Audio(track.music);
            this.audio.volume = this.audioVolume;

            if ('mediaSession' in navigator) {
                navigator.mediaSession.metadata = new MediaMetadata({
                    title: track.title,
                    artist: track.artists.map(artist => artist.name).join(', '),
                    album: track.album.title,
                    artwork: [
                        { type: 'image/jpeg', src: track.album.cover, sizes: '1024x1024' }
                    ]
                });

                navigator.mediaSession.setActionHandler('play', this.play.bind(this));
                navigator.mediaSession.setActionHandler('pause', this.pause.bind(this));
                navigator.mediaSession.setActionHandler('stop', this.pause.bind(this));
                navigator.mediaSession.setActionHandler('seekbackward', this.seekBackward.bind(this));
                navigator.mediaSession.setActionHandler('seekforward', this.seekForward.bind(this));
                navigator.mediaSession.setActionHandler('seekto', this.seekTo.bind(this));
                navigator.mediaSession.setActionHandler('previoustrack', this.previousTrack.bind(this));
                navigator.mediaSession.setActionHandler('nexttrack', this.nextTrack.bind(this));
            }

            this.audio.onloadedmetadata = () => {
                this.audioDuration = this.audio.duration;
                this.audioCurrentTime = this.audio.currentTime;
                this.play();
                this.updatePositionState();
            };
            this.audio.onratechange = () => {
                this.updatePositionState();
            };
        },

        updateCurrentTime() {
            this.audioCurrentTime = this.audio.currentTime;
            if (this.isPlaying) {
                setTimeout(this.updateCurrentTime.bind(this), UPDATE_TIMEOUT);
            }
        },

        updatePositionState() {
            this.audioCurrentTime = this.audio.currentTime;

            if ('setPositionState' in navigator.mediaSession) {
                navigator.mediaSession.setPositionState({
                    duration: this.audio.duration,
                    playbackRate: this.audio.playbackRate,
                    position: this.audio.currentTime,
                });
            }
        },

        previousTrack() {

        },

        seekBackward(details) {
            if (this.audio.paused) {
                this.play();
            }
            this.audio.currentTime = Math.max(0, this.audio.currentTime - (details.seekOffset || SEEK_TIME));
            this.updatePositionState();
        },

        play() {
            this.audio.play();
            navigator.mediaSession.playbackState = 'playing';
            this.isPlaying = true;
            this.updateCurrentTime();
            this.updatePositionState();
        },

        pause() {
            this.audio.pause();
            navigator.mediaSession.playbackState = 'paused';
            this.isPlaying = false;
        },

        playPause() {
            if (this.isPlaying) {
                this.pause();
            } else {
                this.play();
            }
        },

        seekTo(details) {
            if (this.audio.paused) {
                this.play();
            }
            this.audio.currentTime = details.seekTime;
            this.updatePositionState();
        },

        seekToInput(event) {
            if (this.audio.paused) {
                this.play();
            }
            this.audio.currentTime = event.target.value;
            this.updatePositionState();
        },

        seekForward(details) {
            if (this.audio.paused) {
                this.play();
            }
            this.audio.currentTime = Math.min(this.audio.duration, this.audio.currentTime + (details.seekOffset || SEEK_TIME));
            this.updatePositionState();
        },

        nextTrack() {

        },

        toggleVolume() {
            if (this.audioVolume > 0) {
                this.oldAudioVolume = this.audioVolume;
                this.audioVolume = 0;
            } else {
                if (this.oldAudioVolume != undefined) {
                    this.audioVolume = this.oldAudioVolume;
                    this.oldAudioVolume = undefined;
                } else {
                    this.audioVolume = 1;
                }
            }
        }
    }
};
</script>

<template>
    <div class="sidebar box has-background-white-bis m-0">
        <h1 class="title is-4 mb-5"><router-link to="/">BassieMusic</router-link></h1>
        <div class="menu">
            <p class="menu-label">Library</p>
            <ul class="menu-list mb-5">
                <li><router-link :class="{ 'is-active': $route.path == '/' }" to="/">🏠 &nbsp;Home</router-link></li>
                <li><router-link :class="{ 'is-active': $route.path.startsWith('/artists') }" to="/artists">🎤 &nbsp;Artists</router-link></li>
                <li><router-link :class="{ 'is-active': $route.path.startsWith('/albums') }" to="/albums">💿 &nbsp;Albums</router-link></li>
                <li><router-link :class="{ 'is-active': $route.path.startsWith('/tracks') }" to="/tracks">🎵 &nbsp;Tracks</router-link></li>
            </ul>
            <p class="menu-label">Playlists</p>
            <ul class="menu-list">
                <li>Comming soon...</li>
            </ul>
        </div>
        <div class="flex"></div>
        <p>
            Made with ❤️ by<br/>
            <a href="https://bplaat.nl" target="_blank">Bastiaan van der Plaat</a>
        </p>
    </div>

    <div class="section">
        <router-view></router-view>
    </div>

    <div class="player-controls box has-background-white-bis m-0">
        <div v-if="audio != undefined" style="display: flex; align-items: center;">
            <div class="box mr-4 mb-0" style="padding: 0; overflow: hidden; width: 64px; height: 64px;">
                <img :src="playingTrack.album.cover" style="display: block;">
            </div>

            <div class="mr-5" style="min-width: 10rem">
                <p><router-link :to="`/albums/${playingTrack.album.id}`" style="font-weight: bold;">{{ playingTrack.title }}</router-link></p>
                <p>
                    <router-link v-for="artist in playingTrack.artists" :key="artist.id" class="mr-2"
                        :to="`/artists/${artist.id}`">{{ artist.name }}</router-link>
                </p>
            </div>

            <div class="field has-addons mb-0">
                <p class="control">
                    <button class="button" @click.prevent="previousTrack">
                        <svg class="icon" viewBox="0 0 24 24"><path d="M6,18V6H8V18H6M9.5,12L18,6V18L9.5,12Z" /></svg>
                    </button>
                </p>
                <p class="control">
                    <button class="button" @click.prevent="seekBackward">
                        <svg class="icon" viewBox="0 0 24 24"><path d="M11.5,12L20,18V6M11,18V6L2.5,12L11,18Z" /></svg>
                    </button>
                </p>
                <p class="control">
                    <button class="button" @click.prevent="playPause">
                        <svg v-if="isPlaying" class="icon" viewBox="0 0 24 24"><path d="M14,19H18V5H14M6,19H10V5H6V19Z" /></svg>
                        <svg v-else class="icon" viewBox="0 0 24 24"><path d="M8,5.14V19.14L19,12.14L8,5.14Z" /></svg>
                    </button>
                </p>
                <p class="control">
                    <button class="button" @click.prevent="seekForward">
                        <svg class="icon" viewBox="0 0 24 24"><path d="M13,6V18L21.5,12M4,18L12.5,12L4,6V18Z" /></svg>
                    </button>
                </p>
                <p class="control">
                    <button class="button" @click.prevent="nextTrack">
                        <svg class="icon" viewBox="0 0 24 24"><path d="M16,18H18V6H16M6,18L14.5,12L6,6V18Z" /></svg>
                    </button>
                </p>
            </div>

            <div style="flex: 4; display: flex;">
                <span class="mr-3" style="width: 4rem; text-align: right;">{{ $filters.formatDuration(audioCurrentTime) }}</span>
                <input type="range" style="flex: 1;" :value="audioCurrentTime" @input.prevent="seekToInput" min="0" :max="audioDuration">
                <span class="ml-3" style="width: 4rem;">-{{ $filters.formatDuration(audioDuration - audioCurrentTime) }}</span>
            </div>

            <div style="flex: 1; display: flex;">
                <button class="button mr-3" @click.prevent="toggleVolume">
                    <svg v-if="audioVolume == 0" class="icon" viewBox="0 0 24 24"><path d="M12,4L9.91,6.09L12,8.18M4.27,3L3,4.27L7.73,9H3V15H7L12,20V13.27L16.25,17.53C15.58,18.04 14.83,18.46 14,18.7V20.77C15.38,20.45 16.63,19.82 17.68,18.96L19.73,21L21,19.73L12,10.73M19,12C19,12.94 18.8,13.82 18.46,14.64L19.97,16.15C20.62,14.91 21,13.5 21,12C21,7.72 18,4.14 14,3.23V5.29C16.89,6.15 19,8.83 19,12M16.5,12C16.5,10.23 15.5,8.71 14,7.97V10.18L16.45,12.63C16.5,12.43 16.5,12.21 16.5,12Z" /></svg>
                    <svg v-if="audioVolume > 0 && audioVolume < 0.33" class="icon" viewBox="0 0 24 24"><path d="M7,9V15H11L16,20V4L11,9H7Z" /></svg>
                    <svg v-if="audioVolume >= 0.33 && audioVolume < 0.67" class="icon" viewBox="0 0 24 24"><path d="M5,9V15H9L14,20V4L9,9M18.5,12C18.5,10.23 17.5,8.71 16,7.97V16C17.5,15.29 18.5,13.76 18.5,12Z" /></svg>
                    <svg v-if="audioVolume >= 0.67" class="icon" viewBox="0 0 24 24"><path d="M14,3.23V5.29C16.89,6.15 19,8.83 19,12C19,15.17 16.89,17.84 14,18.7V20.77C18,19.86 21,16.28 21,12C21,7.72 18,4.14 14,3.23M16.5,12C16.5,10.23 15.5,8.71 14,7.97V16C15.5,15.29 16.5,13.76 16.5,12M3,9V15H7L12,20V4L7,9H3Z" /></svg>
                </button>
                <input type="range" style="flex: 1;" v-model="audioVolume" min="0" max="1" step="0.01">
            </div>
        </div>
    </div>
</template>
