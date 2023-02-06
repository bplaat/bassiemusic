import { browser } from '$app/environment';
import { writable } from 'svelte/store';

export const trackAutoplay = writable(false);
export const trackPosition = writable(0);
export const playingQueue = writable([]);
export const playingTrack = writable(null);

export const audioVolume = writable(browser ? (localStorage.getItem('audio_volume') ?? 1) : 1);
audioVolume.subscribe(audioVolume => {
    if (browser) localStorage.setItem('audio_volume', audioVolume);
});
