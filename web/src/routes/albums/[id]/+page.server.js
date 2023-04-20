import { error } from '@sveltejs/kit';

export async function load({ locals, fetch, cookies, params }) {
    // Fetch album
    const response = await fetch(`${import.meta.env.VITE_API_URL}/albums/${params.id}`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    if (response.status === 404) {
        throw error(404, 'Not Found');
    }
    const album = await response.json();
    album.tracks = album.tracks.slice().map((track) => {
        track.album = album;
        return track;
    });

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        album,
    };
}
