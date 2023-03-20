import { writable } from 'svelte/store';

export const language = writable('en');

export const sidebar = writable();

export const musicPlayer = writable();

export const musicState = writable({
    queue: [],
    track: undefined,
});
