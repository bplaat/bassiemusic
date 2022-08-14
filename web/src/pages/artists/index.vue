<script>
import { API_URL } from '../../config.js';

export default {
    data() {
        return {
            artists: []
        }
    },

    created() {
        this.fetchPage(1);
    },

    methods: {
        async fetchPage(page) {
            const response = await fetch(`${API_URL}/artists?page=${page}`);
            const newArtists = await response.json();
            if (newArtists.length > 0) {
                this.artists.push(...newArtists);
                this.fetchPage(page + 1);
            }
        }
    }
};
</script>

<template>
    <div>
        <h2 class="title">Artists</h2>
        <div class="columns is-multiline">
            <div class="column is-one-fifth" v-for="artist in artists" :key="artist.id">
                <router-link class="card" :to="'/artists/' + artist.id">
                    <div class="card-image" :style="{ backgroundImage: 'url(' + artist.image + ')' }"></div>
                    <div class="card-content">
                        <h3 class="title is-6">{{ artist.name }}</h3>
                    </div>
                </router-link>
            </div>
        </div>
    </div>
</template>
