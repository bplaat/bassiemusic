<script>
import { API_URL } from '../../config.js';

export default {
    data() {
        return {
            tracks: []
        }
    },

    created() {
        this.fetchPage(1);
    },

    methods: {
        async fetchPage(page) {
            const response = await fetch(`${API_URL}/tracks?page=${page}`);
            const newTracks = await response.json();
            if (newTracks.length > 0) {
                this.tracks.push(...newTracks);
                this.fetchPage(page + 1);
            }
        },

        playTrack(id) {
            const track = this.tracks.find(track => track.id == id);
            this.$root.playTrack(track);
        }
    }
};
</script>

<template>
    <div>
        <h2 class="title">Tracks</h2>
        <table class="table is-fullwidth">
            <thead>
                <tr>
                    <th style="width: 5%;">#</th>
                    <th style="width: 30%;">Title</th>
                    <th style="width: 30%;">Album</th>
                    <th style="width: 20%;">Duration</th>
                    <th style="width: 15%;">Plays</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="track, index in tracks" :key="track.id" @dblclick.prevent="playTrack(track.id)"
                    :class="{ 'has-background-light': $root.playingTrack != undefined && track.id == $root.playingTrack.id }">
                    <td style="vertical-align: middle;">{{ index + 1 }}</td>
                    <td style="display: flex;">
                        <div class="box mr-4 mb-0" style="padding: 0; overflow: hidden; width: 64px; height: 64px;">
                            <img :src="track.album.cover" style="display: block;">
                        </div>
                        <div style="flex: 1; display: flex; flex-direction: column; justify-content: center;">
                            <p><router-link :to="'/albums/' + track.album.id" style="font-weight: bold;">{{ track.title }}</router-link></p>
                            <p><router-link v-for="artist in track.artists" :key="artist.id" class="mr-2"
                                :to="'/artists/' + artist.id">{{ artist.name }}</router-link></p>
                        </div>
                    </td>
                    <td style="vertical-align: middle;"><router-link :to="'/albums/' + track.album.id">{{ track.album.title }}</router-link></td>
                    <td style="vertical-align: middle;">{{ $filters.formatDuration(track.duration) }}</td>
                    <td style="vertical-align: middle;">{{ track.plays }}</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>
