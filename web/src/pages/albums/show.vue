<script>
import { API_URL } from '../../config.js';

export default {
    data() {
        return {
            album: undefined
        }
    },

    created() {
        this.fetchAlbum(this.$route.params.id);
    },

    methods: {
        async fetchAlbum(id) {
            const response = await fetch(`${API_URL}/albums/${id}`);
            this.album = await response.json();
            this.album.tracks = this.album.tracks.map(track => {
                track.album = this.album;
                return track;
            });
        },

        playTrack(id) {
            const track = this.album.tracks.find(track => track.id == id);
            track.album = this.album;
            this.$root.playTrack(track);
        }
    }
};
</script>

<template>
    <div>
        <div v-if="album != undefined">
            <p class="mb-3">
                <router-link to="/albums" style="font-weight: bold;">&lt;&lt; Back to albums</router-link>
            </p>
            <div class="columns">
                <div class="column is-one-quarter mr-5">
                    <div class="box" style="padding: 0; overflow: hidden;">
                        <img :src="album.cover" style="display: block;">
                    </div>
                </div>
                <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
                    <h2 class="title mb-3">{{ album.title }}</h2>
                    <p class="mb-3">{{ album.released_at.split('T')[0] }}</p>
                    <p>
                        <router-link v-for="artist in album.artists" :key="artist.id" class="mr-2"
                            :to="`/artists/${artist.id}`">{{ artist.name }}</router-link>
                    </p>
                </div>
            </div>

            <h3 class="title is-4">Tracks</h3>
            <table class="table is-fullwidth">
                <thead>
                    <tr>
                        <th style="width: 10%;">#</th>
                        <th style="width: 30%;">Artist</th>
                        <th style="width: 30%;">Title</th>
                        <th style="width: 15%;">Duration</th>
                        <th style="width: 15%;">Plays</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="track in album.tracks" :key="track.id" @dblclick.prevent="playTrack(track.id)"
                        :class="{ 'has-background-light': $root.playingTrack != undefined && track.id == $root.playingTrack.id }">
                        <td>{{ track.position }}</td>
                        <td>
                            <router-link v-for="artist in track.artists" :key="artist.id" class="mr-2"
                                :to="`/artists/${artist.id}`">{{ artist.name }}</router-link>
                        </td>
                        <td>{{ track.title }}</td>
                        <td>{{ $filters.formatDuration(track.duration) }}</td>
                        <td>{{ track.plays }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>
