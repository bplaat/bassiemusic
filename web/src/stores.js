import { browser } from '$app/environment';
import { writable } from 'svelte/store';

export const autoplay = writable(false);

export const playingTrack = writable(null);

export const playingQueue = writable([]);

export const audioVolume = writable(browser ? (localStorage.getItem('audio_volume') ?? 1) : 1);
audioVolume.subscribe(audioVolume => {
    if (browser) localStorage.setItem('audio_volume', audioVolume);
});
