import { browser } from '$app/env';
import { writable } from 'svelte/store';

export const playingTrack = writable(undefined);
export const playingQueue = writable([]);
export const audioVolume = writable(browser ? (localStorage.audioVolume ?? 1) : 1);
audioVolume.subscribe(audioVolume => {
    if (browser) {
        localStorage.audioVolume = audioVolume;
    }
});
