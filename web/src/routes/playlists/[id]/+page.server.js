import { error } from '@sveltejs/kit';

export async function load({ locals, fetch, cookies, params }) {
    // Fetch playlist
    const response = await fetch(`${import.meta.env.VITE_API_URL}/playlists/${params.id}`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    if (response.status !== 200) {
        throw error(404, 'Not Found');
    }
    const playlist = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        playlist,
    };
}
