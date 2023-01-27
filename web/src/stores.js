import { browser } from '$app/environment';
import { writable } from 'svelte/store';

// TODO
export const token = writable(browser ? (localStorage.token ?? undefined) : undefined);
token.subscribe(token => {
    if (browser) localStorage.token = token;
});

export const playingTrack = writable(undefined);
export const playingQueue = writable([]);
export const audioVolume = writable(browser ? (localStorage.audioVolume ?? 1) : 1);
audioVolume.subscribe(audioVolume => {
    if (browser) {
        localStorage.audioVolume = audioVolume;
    }
});
