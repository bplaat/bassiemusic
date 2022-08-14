<script setup>
import AlbumCard from '../../components/AlbumCard.vue';
</script>

<script>
import { API_URL } from '../../config.js';

export default {
    data() {
        return {
            artist: undefined
        }
    },

    created() {
        this.fetchArtist(this.$route.params.id);
    },

    methods: {
        async fetchArtist(id) {
            const response = await fetch(`${API_URL}/artists/${id}`);
            this.artist = await response.json();
        }
    }
};
</script>

<template>
    <div>
        <div v-if="artist != undefined">
            <p class="mb-3">
                <router-link to="/artists" style="font-weight: bold;">&lt;&lt; Back to artists</router-link>
            </p>
            <div class="columns">
                <div class="column is-one-quarter mr-5">
                    <div class="box" style="padding: 0; overflow: hidden;">
                        <img :src="artist.image" style="display: block;">
                    </div>
                </div>
                <div class="column" style="display: flex; flex-direction: column; justify-content: center;">
                    <h2 class="title mb-3">{{ artist.name }}</h2>
                    <p>Albums: {{ artist.albums != undefined ? artist.albums.length : 0 }}</p>
                </div>
            </div>

            <h3 class="title is-4">Albums</h3>
            <div v-if="artist.albums != undefined && artist.albums.length > 0" class="columns is-multiline">
                <div class="column is-one-fifth" v-for="album in artist.albums" :key="album.id">
                    <album-card :album="album" />
                </div>
            </div>
            <p v-else><i>This artist has no albums of it's own</i></p>
        </div>
    </div>
</template>
