<script setup>
import AlbumCard from '../../components/AlbumCard.vue';
</script>

<script>
import { API_URL } from '../../config.js';

export default {
    data() {
        return {
            albums: []
        }
    },

    created() {
        this.fetchPage(1);
    },

    methods: {
        async fetchPage(page) {
            const response = await fetch(`${API_URL}/albums?page=${page}`);
            const newAlbums = await response.json();
            if (newAlbums.length > 0) {
                this.albums.push(...newAlbums);
                this.fetchPage(page + 1);
            }
        }
    }
};
</script>

<template>
    <div>
        <h2 class="title">Albums</h2>
        <div class="columns is-multiline">
            <div class="column is-one-fifth" v-for="album in albums" :key="album.id">
                <album-card :album="album" />
            </div>
        </div>
    </div>
</template>
