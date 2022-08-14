<script>
export default {
    props: ['album']
};
</script>

<template>
    <router-link class="card" :to="`/albums/${album.id}`">
        <div class="card-image" :style="{ backgroundImage: `url(${album.cover})` }"></div>
        <div class="card-content">
            <h3 class="title is-6 mb-2">{{ album.title }}</h3>
            <p>
                <router-link v-for="artist in album.artists" :key="artist.id" class="mr-2"
                    :to="'/artists/' + artist.id">{{ artist.name }}</router-link>
            </p>
        </div>
    </router-link>
</template>
