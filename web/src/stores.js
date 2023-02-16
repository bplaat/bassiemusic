import { writable } from 'svelte/store';

export const musicPlayer = writable();

export const musicState = writable({
    queue: [],
    track: undefined,
});
