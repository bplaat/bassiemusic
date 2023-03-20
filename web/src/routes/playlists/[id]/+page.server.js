import { error } from '@sveltejs/kit';
import { isAuthedMiddleware } from '../../../middlewares/auth.js';

export async function load({ url, fetch, cookies, params }) {
    const authUser = await isAuthedMiddleware({ url, fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/playlists/${params.id}`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    if (response.status != 200) {
        throw error(404, 'Not Found');
    }
    const playlist = await response.json();

    return { token: cookies.get('token'), authUser, playlist };
}