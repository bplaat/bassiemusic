import { error } from '@sveltejs/kit';
import { isAuthedMiddleware } from '../../../middlewares/auth.js';

export async function load({ url, fetch, cookies, params }) {
    const authUser = await isAuthedMiddleware({ url, fetch, cookies });

    // Get genre
    const response = await fetch(`${import.meta.env.VITE_API_URL}/genres/${params.id}`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    if (response.status == 404) {
        throw error(404, 'Not Found');
    }
    const genre = await response.json();

    // Get genre first albums page
    const albumsResponse = await fetch(`${import.meta.env.VITE_API_URL}/genres/${params.id}/albums`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data, pagination } = await albumsResponse.json();
    genre.albums = data;

    return { token: cookies.get('token'), authUser, genre, albumsTotal: pagination.total, albumsPage: 2 };
}
