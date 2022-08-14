import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import './style.css';
import App from './App.vue';

import Home from './pages/home.vue';
import ArtistsIndex from './pages/artists/index.vue';
import ArtistsShow from './pages/artists/show.vue';
import AlbumsIndex from './pages/albums/index.vue';
import AlbumsShow from './pages/albums/show.vue';
import TracksIndex from './pages/tracks/index.vue';

const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: '/', component: Home },
        { path: '/artists', component: ArtistsIndex },
        { path: '/artists/:id', component: ArtistsShow },
        { path: '/albums', component: AlbumsIndex },
        { path: '/albums/:id', component: AlbumsShow },
        { path: '/tracks', component: TracksIndex }
    ]
});

const app = createApp(App);

app.config.globalProperties.$filters = {
    formatDuration(totalSeconds) {
        totalSeconds = Math.floor(totalSeconds);
        const hours = Math.floor(totalSeconds / 3600);
        const minutes = Math.floor((totalSeconds % 3600) / 60);
        const seconds = totalSeconds % 60;
        if (totalSeconds > 3600) {
            return hours + (minutes < 10 ? '0' + minutes : minutes) + ':' +
                (seconds < 10 ? '0' + seconds : seconds);
        }
        return minutes + ':' + (seconds < 10 ? '0' + seconds : seconds);
    }
};

app.use(router);
app.mount('#app');
