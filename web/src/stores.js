import { browser } from '$app/environment';
import { writable } from 'svelte/store';

export const musicPlayer = writable({
    action: 'init',
    queue: []
});

export const audioVolume = writable(browser ? (localStorage.getItem('audio_volume') ?? 1) : 1);
audioVolume.subscribe(audioVolume => {
    if (browser) localStorage.setItem('audio_volume', audioVolume);
});
